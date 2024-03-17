<template>
  <div class="player-container">
    <w-player v-if="playerKey !== 0" class="player" :key="playerKey" :options="options"></w-player>
  </div>
</template>
    
<script setup lang="ts">
import Hls from "hls.js";
import { WPlayer } from 'vue-wplayer';
import { ref, onBeforeMount, watch } from 'vue';
import { getDanmakuAPI, sendDanmakuAPI } from "@/api/danmaku";
import { getResourceQualityApi, getVideoFileUrl } from "@/api/video";

const props = withDefaults(defineProps<{
  videoInfo: VideoType;
  part: number
}>(), {
  part: 1
})

const playerKey = ref(0);

const hls = shallowRef<Hls | null>(null);
const loadHls = (src: string, player: HTMLVideoElement) => {
  if (!hls.value) hls.value = new Hls();
  hls.value.loadSource(src);
  hls.value.attachMedia(player);
  hls.value.on(Hls.Events.ERROR, () => {
    console.error("资源加载失败");
  });
}

const stopLoadHls = () => {
  if (hls.value) {
    hls.value.stopLoad();
    hls.value.destroy();
    hls.value = null;
  }
}

const options: OptionsType = {
  resource: [],
  type: navigator.userAgent.includes("iPhone") ? "mp4" : "hls",
  customType: (player, src) => {
    loadHls(src, player);
  },
  theme: "var(--primary-color)",
  danmaku: {
    open: true,
    data: [],
    send: (danmaku: AddDanmakuType) => {
      sendDanmaku(danmaku);
    }
  }
}

const loadPart = async (part: number) => {
  await loadResource(part);
  await getDanmaku(part);

  playerKey.value = Date.now();
}

const resourceNameMap = {
  "640x360_500k_30": "360p",
  "854x480_900k_30": "480p",
  "1080x720_2000k_30": "720p",
  "1920x1080_3000k_30": "780p",
}

const loadResource = async (part: number) => {
  const resource = props.videoInfo.resources[part - 1];
  const res = await getResourceQualityApi(resource.id);
  if (res.data.code === statusCode.OK) {
    options.resource = res.data.data.quality.map((item: keyof typeof resourceNameMap) => {
      return {
        name: resourceNameMap[item],
        url: getVideoFileUrl(resource.id, item),
      }
    })
  }
}

const sendDanmaku = (danmakuForm: AddDanmakuType) => {
  danmakuForm.vid = props.videoInfo.vid;
  danmakuForm.part = props.part;
  sendDanmakuAPI(danmakuForm);
}

const getDanmaku = async (part: number) => {
  if (options.danmaku?.data) options.danmaku.data = [];
  const res = await getDanmakuAPI(props.videoInfo.vid, part);
  if (res.data.code === statusCode.OK) {
    const list = res.data.data.danmaku;
    if (list) {
      options.danmaku.data = list;
    }
  }
}

watch(() => props.part, (val) => {
  loadPart(val);
});

onBeforeMount(() => {
  console.log('options', options)
  loadPart(props.part);
});

onBeforeUnmount(() => {
  stopLoadHls();
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
}
</style>