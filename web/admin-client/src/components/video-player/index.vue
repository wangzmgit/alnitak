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
    autoplay: true,  // 确保播放器设置为自动播放
    muted: true,     // 设置为静音
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
  "640x360_500k": "360p",
  "854x480_900k": "480p",
  "1080x720_2000k": "720p", // 兼容之前的错误
  "1280x720_2000k": "720p",
  "1920x1080_3000k": "1080p",
};

const loadResource = async (resourceId: number) => {
  const res = await getResourceQualityApi(resourceId);

  if (res.data.code === statusCode.OK) {
    options.video.quality = [];

    res.data.data.quality.forEach((item: string, i: number) => {
      // 使用正则表达式移除分辨率和码率后的小数帧率部分
      const normalizedItem = item.replace(/(_\d+)(?:\.\d+)?$/, ""); // 去掉可能出现的帧率部分
      const qualityName = resourceNameMap[normalizedItem];

      if (qualityName === defaultQuality.value) {
        options.video.defaultQuality = i;
      }

      options.video.quality.push({
        name: qualityName || "Unknown", // 如果未找到匹配的名称，显示 "Unknown"
        url: getVideoFileUrl(resourceId, item),
      });
    });
  }
};

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