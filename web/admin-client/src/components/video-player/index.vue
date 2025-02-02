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
  "640x360": "360p",
  "854x480": "480p",
  "1280x720": "720p",
  "1920x1080": "1080p",
};

const loadResource = async (resourceId: number) => {
  const res = await getResourceQualityApi(resourceId);

  if (res.data.code === statusCode.OK) {
    let videoQuality: PlayerQualityType[] = [];
    let has60fps = false; // 标记是否存在60fps

    res.data.data.quality.forEach((item: string) => {
      // 使用正则表达式提取分辨率和帧率
      const resolutionMatch = item.match(/^(\d{3,4}x\d{3,4})(_\d+k)?(_\d+)?$/);
      
      if (resolutionMatch) {
        const resolution = resolutionMatch[1]; // 分辨率
        const frameRate = resolutionMatch[3]; // 帧率（如果存在）

        // 根据是否有帧率来判断应该显示什么分辨率
        let qualityName = resourceNameMap[resolution as keyof typeof resourceNameMap] || 'Unknown';

        // 如果有帧率并且是60fps，设置为1080p60
        if (frameRate === '_60') {
          qualityName = '1080p60';
          has60fps = true; // 如果有60fps标记为true
        }

        // 将质量名称和对应的视频 URL 添加到数组中
        videoQuality.push({
          name: qualityName,
          url: getVideoFileUrl(resourceId, item),
        });
      }
    });

    // 如果没有60fps（即所有视频清晰度都没有包含60fps），则不考虑1080p60
    if (!has60fps) {
      // 排序：首先按分辨率从高到低（1080p > 720p > 480p > 360p）
      videoQuality.sort((a, b) => {
        const aRes = parseInt(a.name.replace('p', '').replace('60', ''));
        const bRes = parseInt(b.name.replace('p', '').replace('60', ''));

        return bRes - aRes; // 从高到低排序
      });
    } else {
      // 如果有60fps帧率，按分辨率和帧率从高到低排序
      videoQuality.sort((a, b) => {
        const aRes = parseInt(a.name.replace('p', '').replace('60', ''));
        const bRes = parseInt(b.name.replace('p', '').replace('60', ''));

        // 优先处理帧率：60fps 优先
        if (a.name.includes("60") && !b.name.includes("60")) {
          return -1; // 1080p60 比 1080p 高
        } else if (!a.name.includes("60") && b.name.includes("60")) {
          return 1; // 1080p 比 1080p60 低
        }

        // 若分辨率相同，按数字排序（例如：1080p vs 720p vs 480p vs 360p）
        return bRes - aRes; // 从高到低排序
      });
    }

    // 确保播放器使用排序后的清晰度，避免选择默认加载项
    options.video.quality = videoQuality;
    options.video.defaultQuality = 0;  // 选择排序后的第一项作为默认清晰度

    console.log("Sorted Video Qualities:", videoQuality);  // 调试输出排序后的质量
    console.log("Default Quality Set:", videoQuality[0].name);  // 调试输出默认质量
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