<template>
  <div class="player-container">
    <div class="player" id="dplayer"></div>
    <div class="danmaku-send">
      <danmaku-send ref="danmakuSendRef" @send="sendDanmaku" @change-show="changeShow" @opacity-change="opacityChange"
        @set-filter="filterDanmaku"></danmaku-send>
    </div>
  </div>
</template>

<script setup lang="ts">
import Hls from "hls.js";
import Wplayer from 'wplayer-next';
import { ref, onBeforeMount, watch, onMounted } from 'vue';
import { getDanmakuAPI, sendDanmakuAPI } from "@/api/danmaku";
import DanmakuSend from "./components/DanmakuSend.vue";
import { getResourceQualityApi, getVideoFileUrl } from "@/api/video";
import { addHistoryAPI, getHistoryProgressAPI } from "@/api/history";

const props = withDefaults(defineProps<{
  videoInfo: VideoType;
  part: number
}>(), {
  part: 1
})

let player: any = null;
const defaultQuality = ref('');
const hls = shallowRef<Hls | null>(null);
const danmakuSendRef = ref<InstanceType<typeof DanmakuSend> | null>(null);
const options: PlayerOptionsType = {
  container: null,
  video: {
    quality: [],
    defaultQuality: 0,
    pic: 'demo.png',
    type: 'customHls',
    customType: {
      // TODO: 处理IOS系统中的hls视频播放
      customHls: function (video: HTMLVideoElement) {
        if (!hls.value) hls.value = new Hls();
        hls.value.loadSource(video.src);
        hls.value.attachMedia(video);
        hls.value.on(Hls.Events.ERROR, () => {
          console.error("资源加载失败");
        });
      },
    },
  },
  danmaku: {}
}

let disableLeave = 0;
let disableType: number[] = [];
const initFilterConfig = () => {
  const disableTypeConfig = localStorage.getItem('danmaku-disable-type');
  if (disableTypeConfig) {
    disableType = disableTypeConfig.split(',').map((item) => parseInt(item));
  }

  const disableLeaveConfig = localStorage.getItem('danmaku-disable-leave');
  if (disableLeaveConfig) {
    disableLeave = parseInt(disableLeaveConfig);
  }
}

const loadPart = async (part: number) => {
  const el = document.getElementById('dplayer');
  if (el) {
    await loadResource(part);
    await getDanmaku(part);

    if (player) player.destroy();

    options.container = el;
    player = new Wplayer(options);
    player.on('quality_start', (quality: PlayerQualityType) => {
      localStorage.setItem('default-video-quality', quality.name);
    })
    filterDanmaku({ disableLeave, disableType });

  }
}

const resourceNameMap = {
  "640x360_500k_30": "360p",
  "854x480_900k_30": "480p",
  "1080x720_2000k_30": "720p",// 兼容之前的错误
  "1280x720_2000k_30": "720p",
  "1920x1080_3000k_30": "1080p",
  "1920x1080_6000k_60": "1080p60",
}

const loadResource = async (part: number) => {
  const resource = props.videoInfo.resources[part - 1];
  const res = await getResourceQualityApi(resource.id);
  if (res.data.code === statusCode.OK) {
    options.video.quality = res.data.data.quality.map((item: keyof typeof resourceNameMap, index: number) => {
      if (resourceNameMap[item] === defaultQuality.value) {
        options.video.defaultQuality = index;
      }
      return {
        name: resourceNameMap[item],
        url: getVideoFileUrl(resource.id, item),
      }
    })
  }
}

let originalDanmaku: DanmakuType[] = [];
const getDanmaku = async (part: number) => {
  const res = await getDanmakuAPI(props.videoInfo.vid, part);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.danmaku) {
      originalDanmaku = res.data.data.danmaku;
    }
  }
}

// 弹幕显示改变
const changeShow = (val: boolean) => {
  if (val) {
    player.danmaku.show();
  } else {
    player.danmaku.hide();
  }
}

const opacityChange = (val: number) => {
  player.danmaku.opacity(val);
}

const sendDanmaku = (danmakuForm: DrawDanmakuType) => {
  if (danmakuForm.text == "") {
    player.notice("弹幕内容不能为空")
    return;
  }

  player.danmaku.send(danmakuForm, (danmaku: AddDanmakuType) => {
    danmaku.vid = props.videoInfo.vid;
    danmaku.part = props.part;
    sendDanmakuAPI(danmaku);
  })
}

//过滤弹幕
const filterDanmaku = (filter: FilterDanmakuType) => {
  localStorage.setItem('danmaku-disable-type', filter.disableType.toString());
  localStorage.setItem('danmaku-disable-leave', filter.disableLeave.toString());

  const data = originalDanmaku.filter((item) => {
    return !isDisableType(item, filter.disableType) && (Math.floor(Math.random() * 10) + 1) > filter.disableLeave;
  }).map((d) => { return { ...d } });

  player.danmaku.update(data);

  player.on('danmaku_loaded', () => {
    console.log("danmaku_load_end")
  })

  // 更新弹幕数量
  danmakuSendRef.value?.updateDanmakuCount(data.length);
}

//是否为屏蔽类型
const isDisableType = (item: DanmakuType, disableType: Array<number>) => {
  if (disableType.includes(item.type))
    return true;
  if (disableType.includes(3) && (item.color !== '#fff' && item.color !== '#ffffff'))
    return true;

  return false;
}

// 上传历史记录
const uploadHistory = async () => {
  await addHistoryAPI({ vid: props.videoInfo.vid, part: props.part, time: player.video.currentTime });
}

// 获取历史记录
const getHistoryProgress = async () => {
  const res = await getHistoryProgressAPI(props.videoInfo.vid, props.part);
  if (res.data.code === statusCode.OK) {
    player.seek(res.data.data.progress)
  } else {
    uploadHistory();
  }
}

watch(() => props.part, (val) => {
  loadPart(val);
});

let timer: number | null = null;
onMounted(async () => {
  const quality = localStorage.getItem('default-video-quality');
  if (quality) {
    defaultQuality.value = quality;
  } else {
    defaultQuality.value = '720p';
    localStorage.setItem('default-video-quality', '720p');
  }

  initFilterConfig();
  await loadPart(props.part);

  // 获取当前播放进度（如果没有就保存历史记录）
  await getHistoryProgress();

  timer = window.setInterval(() => {
    uploadHistory();
  }, 10000)
})

onBeforeUnmount(() => {
  if (timer) clearInterval(timer);
})
</script>

<style lang="scss" scoped>
.player-container {
  height: 0;
  width: 100%;
  padding-bottom: 56.25%;
  position: relative;
  margin-bottom: 40px;

  .player {
    width: 100%;
    height: 100%;
    position: absolute;
    background-color: black;
  }


  .danmaku-send {
    position: absolute;
    width: 100%;
    bottom: -40px;
  }
}
</style>