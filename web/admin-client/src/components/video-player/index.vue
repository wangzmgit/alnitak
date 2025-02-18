<template>
  <div class="player-container">
    <div class="player" id="wplayer"></div>
  </div>
</template>
    
<script setup lang="ts">
import Hls from "hls.js";
import Wplayer from 'wplayer-next';
import { statusCode } from "@/utils/status-code";
import { ref, shallowRef, onBeforeUnmount } from 'vue';
import { getResourceQualityApi, getVideoFileUrl, getVideoFileAPI } from "@/api/video";
import { useMessage } from "naive-ui";
import { getResourceUrl } from "@/utils/resource";

const message = useMessage();

let player: any = null;
const defaultQuality = ref('');
const hls = shallowRef<Hls | null>(null);
const options: PlayerOptionsType = {
  container: null,
  video: {
    quality: [],
    defaultQuality: 0,
    type: 'customHls',
    customType: {
      customHls: (video: HTMLVideoElement) => {
        getVideoFileAPI(video.src).then((res) => {
          if (!res.data) return;
          if (!hls.value) hls.value = new Hls();
          const indexFile = res.data.split('\n').map((line: string) => {
            if (line.includes(".ts")) {
              return getResourceUrl(line)
            } else {
              return line
            }
          })
          var blob = new Blob([indexFile.join('\n')], { type: 'text/plain' });
          var blobUrl = URL.createObjectURL(blob);
          hls.value.loadSource(blobUrl);
          hls.value.attachMedia(video);
          hls.value.on(Hls.Events.ERROR, () => {
            console.error("资源加载失败");
          });
        })
      },
    },
  },
  danmaku: {}
}

const loadVideo = async (resourceId: number) => {
  const el = document.getElementById('wplayer');
  if (el) {
    await loadResource(resourceId);

    if (player) player.destroy();

    options.container = el;
    player = new Wplayer(options);
    player.on('quality_start', (quality: PlayerQualityType) => {
      localStorage.setItem('default-video-quality', quality.name);
    })
  }
}

const resourceNameMap = {
  "640x360_500k_30": "360p",
  "854x480_900k_30": "480p",
  "1080x720_2000k_30": "720p", // 兼容之前的错误
  "1280x720_2000k_30": "720p",
  "1920x1080_3000k_30": "1080p",
  "1920x1080_6000k_60": "1080p60",
}

const loadResource = async (resourceId: number) => {
  const res = await getResourceQualityApi(resourceId);
  if (res.data.code === statusCode.OK) {
    options.video.quality = [];
    for (let i = 0; i < res.data.data.quality.length; i++) {
      const item = res.data.data.quality[i] as (keyof typeof resourceNameMap);
      if (resourceNameMap[item] === defaultQuality.value) {
        options.video.defaultQuality = i;
      }

      options.video.quality[i] = {
        name: resourceNameMap[item],
        url: getVideoFileUrl(resourceId, item),
      }
    }
  }
}

defineExpose({
  loadVideo
})

onBeforeUnmount(() => {
  if (player) player.destroy();
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