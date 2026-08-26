<template>
  <div class="embed-video-container">
    <div class="player-box" @mouseenter="showOverlay" @mouseleave="hideOverlay">
      <client-only>
        <div class="player-info-overlay" v-if="videoInfo" v-show="infoOverlayVisible" @mouseenter="showOverlay" @mouseleave="hideOverlay">
          <a class="title-link" :href="videoLink" target="_blank">{{ displayedTitle }}</a>
          <div class="up-info">
            <oss-image class="avatar" :src="videoInfo.author.avatar" :alt="videoInfo.author.name" />
            <a class="up-name" :href="userLink" target="_blank">{{ videoInfo.author.name }}</a>
          </div>
        </div>
        <embed-player v-if="videoInfo" :video-info="videoInfo" :part="currentPart" :progress="pendingProgress" />
      </client-only>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import EmbedPlayer from '@/components/embed-player/art.vue';
import { asyncGetVideoInfoAPI } from '@/api/video';
import { getDanmakuAPI } from '@/api/danmaku';
import { getResourceUrl } from '@/utils/resource';
import { globalConfig } from '@/utils/global-config';

const route = useRoute();
const router = useRouter();
const videoInfo = ref<VideoType>();

const v = route.query.v;
if (!v || typeof v !== 'string') {
  await navigateTo('/404');
  throw new Error('missing v');
}
const videoId = v;

const currentPart = ref(1);
const pendingProgress = ref<number | null>(null);
const playerRef = ref();
const playerReady = ref(false);
const infoOverlayVisible = ref(false);
let hoverCount = 0;

const displayedTitle = computed(() => {
  return videoInfo.value?.resources?.[currentPart.value - 1]?.title || videoInfo.value?.title;
});

const pageTitle = computed(() => {
  const videoTitle = videoInfo.value?.title || '';
  const partTitle = videoInfo.value?.resources?.[currentPart.value - 1]?.title;

  if (partTitle && partTitle !== videoTitle) {
    return `${partTitle} - ${videoTitle} - ${globalConfig.title}`;
  }

  return `${videoTitle} - ${globalConfig.title}`;
});

const videoLink = computed(() => {
  const vVal = videoInfo.value?.shortId || String(videoInfo.value?.vid ?? '');
  const currentRid = videoInfo.value?.resources?.[currentPart.value - 1]?.shortId;
  if (currentRid) {
    return `/watch?v=${vVal}&rid=${currentRid}`;
  }
  return `/watch?v=${vVal}${currentPart.value > 1 ? `&p=${currentPart.value}` : ''}`;
});

const userLink = computed(() => {
  return `/user/${videoInfo.value?.author.uid}`;
});

// 获取视频信息
const fetchVideoInfo = async () => {
  const { data } = await asyncGetVideoInfoAPI(videoId);
  if ((data.value as any).code === 200) {
    videoInfo.value = (data.value as any).data.video as VideoType;
  } else {
    router.replace('/404');
  }
};

// 校验分集参数：rid 优先，不存在则退回到 p
const validatePartQuery = () => {
  if (!videoInfo.value) return;
  if (route.query.rid) {
    const rid = String(route.query.rid);
    const idx = videoInfo.value.resources.findIndex(r => r.shortId === rid);
    currentPart.value = idx >= 0 ? idx + 1 : 1;
    return;
  }
  const partNum = Number(route.query.p) || 1;
  if (videoInfo.value.resources.length && partNum > videoInfo.value.resources.length) {
    router.replace({ path: '/embed/watch', query: { v: route.query.v as string, p: 1 } });
  } else {
    currentPart.value = partNum;
  }
};

// 加载弹幕：rid 优先
const loadDanmaku = async () => {
  if (!videoInfo.value) return;
  const vid = videoInfo.value.shortId || videoInfo.value.vid;
  const rid = videoInfo.value.resources?.[currentPart.value - 1]?.shortId;
  const res = rid
    ? await getDanmakuAPI(vid, undefined, rid)
    : await getDanmakuAPI(vid, currentPart.value);
  if (res.data.code === 200) {
    playerRef.value?.setDanmaku(res.data.data);
  }
};

// 播放器 ready 回调
const onPlayerReady = () => {
  if (pendingProgress.value !== null && playerRef.value?.seek) {
    playerRef.value.seek(pendingProgress.value);
    pendingProgress.value = null;
  }
  playerReady.value = true;
  loadDanmaku();
};

// 路由监听：rid 或 p 变化时都要重算分P
watch(() => [route.query.p, route.query.rid], () => {
  validatePartQuery();
  loadDanmaku();
});

// 播放器 ref 变更
watch(playerRef, (val) => {
  if (val?.setOnReady) {
    val.setOnReady(onPlayerReady);
  }
});

// 弹幕显隐绑定 infoOverlay（可选）
watch(() => playerRef.value?.controlsVisible, (controlsVisible) => {
  infoOverlayVisible.value = controlsVisible !== false;
});

// 鼠标控制
const showOverlay = () => {
  hoverCount++;
  infoOverlayVisible.value = true;
};
const hideOverlay = () => {
  hoverCount--;
  if (hoverCount <= 0) {
    infoOverlayVisible.value = false;
    hoverCount = 0;
  }
};

// 初次加载 - 使用 await 确保数据加载完成
await fetchVideoInfo();
validatePartQuery();

// 更新页面标题
useHead({
  title: pageTitle
});
</script>

<style scoped>
html, body, #__nuxt, .embed-video-container, .player-box, .player-container, .player, #dplayer {
  width: 100vw;
  height: 100vh;
  min-width: 0;
  min-height: 0;
  margin: 0;
  padding: 0;
  background: #000;
  overflow: hidden;
  position: relative;
}
.embed-video-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
.player-box {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100vw;
  height: 100vh;
}
.player-container, .player, #dplayer {
  width: 100%;
  height: 100%;
  min-width: 0;
  min-height: 0;
}
video {
  width: 100% !important;
  height: 100% !important;
  object-fit: contain;
  background: #000;
  display: block;
}
.player-info-overlay {
  position: absolute;
  top: 18px;
  left: 18px;
  background: transparent !important;
  border-radius: 8px;
  padding: 8px 16px 8px 12px;
  display: flex;
  align-items: center;
  gap: 16px;
  z-index: 9999;
  transition: opacity 0.3s;
  box-shadow: none !important;
}
.player-info-overlay .title-link {
  color: #fff;
  font-size: 15px;
  font-weight: 700;
  max-width: 320px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-decoration: none;
  margin-right: 16px;
  transition: color 0.2s, opacity 0.2s;
  text-shadow: none !important;
  opacity: 0.85;
}
.player-info-overlay .title-link:hover {
  color: #fff;
  opacity: 1;
  text-shadow: none;
  text-decoration: none;
}
.player-info-overlay .up-info {
  display: flex;
  align-items: center;
  gap: 8px;
}
.player-info-overlay .avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
  background: #eee;
}
.player-info-overlay .up-name {
  color: #fff;
  font-size: 13px;
  font-weight: 400;
  text-decoration: none;
  opacity: 0.85;
}
.player-info-overlay .up-name:hover {
  text-decoration: underline;
  opacity: 1;
}
</style>
