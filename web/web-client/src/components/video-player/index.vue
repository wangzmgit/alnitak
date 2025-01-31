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

// 获取路由信息
const route = useRoute();

// 设置网络状态自动选择的清晰度
const getNetworkQuality = () => {
  const connection = (navigator as any).connection || (navigator as any).mozConnection || (navigator as any).webkitConnection;
  if (connection) {
    // 根据网络类型或下行速度选择清晰度
    const effectiveType = connection.effectiveType;
    if (effectiveType === "4g" || effectiveType === "3g") {
      return '1080p';  // 高速网络，选择1080p
    } else if (effectiveType === "2g" || effectiveType === "slow-2g") {
      return '360p';  // 网络较差，选择360p
    }
  }
  return '720p';  // 默认选择720p
};

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
      customHls: function (video: HTMLVideoElement) {
        if (!hls.value) hls.value = new Hls();
        hls.value.loadSource(video.src);
        hls.value.attachMedia(video);
        hls.value.on(Hls.Events.ERROR, () => {
          console.error("资源加载失败");
        });
      },
    },
    //autoplay: true,  // 确保播放器设置为自动播放
    //muted: true,     // 设置为静音
  },
  danmaku: {}
}

let disableLeave = 0;
let disableType: number[] = [];
const initFilterConfig = () => {
  const disableTypeConfig = localStorage.getItem('danmaku-disable-type');
  if (disableTypeConfig) {
    disableType = JSON.parse(disableTypeConfig);
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
	// 模拟点击播放
    simulatePlay();
  }
}

const resourceNameMap = {
  "640x360_500k": "360p",
  "854x480_900k": "480p",
  "1080x720_2000k": "720p", 
  "1280x720_2000k": "720p",
  "1920x1080_3000k": "1080p",
};

const loadResource = async (part: number) => {
  const resource = props.videoInfo.resources[part - 1];
  const res = await getResourceQualityApi(resource.id);

  if (res.data.code === statusCode.OK) {
    options.video.quality = res.data.data.quality.map((item: string, index: number) => {
      // 使用正则表达式移除分辨率和码率后的小数帧率部分
      const normalizedItem = item.replace(/(_\d+)(?:\.\d+)?$/, "") as keyof typeof resourceNameMap; // 去掉可能出现的帧率部分
      const qualityName = resourceNameMap[normalizedItem];

      if (qualityName === defaultQuality.value) {
        options.video.defaultQuality = index;
      }

      return {
        name: qualityName || "Unknown", 
        url: getVideoFileUrl(resource.id, item),
      };
    });
  }
};

// 获取当前清晰度设置
const getDefaultQuality = () => {
  const storedQuality = localStorage.getItem('default-video-quality');
  return storedQuality || getNetworkQuality();
};

let originalDanmaku: DanmakuType[] = [];
const getDanmaku = async (part: number) => {
  const res = await getDanmakuAPI(props.videoInfo.vid, part);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.danmaku) {
      originalDanmaku = res.data.data.danmaku;
    }
  }
}

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
  });
  player.danmaku.update(data);

  player.on('danmaku_loaded', () => {
    console.log("danmaku_load_end")
  })

  danmakuSendRef.value?.updateDanmakuCount(data.length);
}

const isDisableType = (item: DanmakuType, disableType: Array<number>) => {
  if (disableType.includes(item.type))
    return true;
  if (disableType.includes(3) && (item.color !== '#fff' && item.color !== '#ffffff'))
    return true;

  return false;
}

const uploadHistory = async () => {
  await addHistoryAPI({ vid: props.videoInfo.vid, part: props.part, time: player.video.currentTime });
}

const getHistoryProgress = async () => {
  const res = await getHistoryProgressAPI(props.videoInfo.vid,props.part);
  if(res.data.code === statusCode.OK){
    player.seek(res.data.data.progress)
  } else {
    uploadHistory();
  }
}

watch(() => props.part, (val) => {
  loadPart(val);
});

let timer: number | null = null;

// 模拟点击播放
const simulatePlay = () => {
  const videoElement = document.querySelector('video');
  if (videoElement) {
    const playButton = videoElement.querySelector('.play-button'); // 如果有播放按钮，模拟点击
    if (playButton) {
      playButton.dispatchEvent(new MouseEvent('click', { bubbles: true }));
    } else {
      videoElement.play();  // 如果没有按钮，直接播放
    }
  }
}


onMounted(async () => {
  defaultQuality.value = getDefaultQuality(); // 自动选择清晰度
  localStorage.setItem('default-video-quality', defaultQuality.value);

  initFilterConfig();
  await loadPart(props.part);

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
