<template>
  <div class="video">
    <header-bar class="header"></header-bar>
    <div class="video-main">
      <div class="mian-content">
        <div class="left-column">
          <div class="video-player" ref="playerContainerRef">
            <client-only>
              <video-player v-if="videoInfo && playerReady" ref="playerRef" :video-info="videoInfo" :part="currentPart"
                :progress="pendingProgress" :key="videoInfo?.vid + '-' + currentPart"></video-player>
            </client-only>
            <div v-if="!showPlayer" class="skeleton"></div>
          </div>
          <!-- 标题和版权信息 -->
          <div class="video-title-box">
            <p class="video-title">{{ pgcInfo?.title || videoInfo?.title }}</p>
            <p v-show="videoInfo?.copyright" class="copyright">
              <el-icon class="icon" color='#fd6d6f'>
                <forbid-icon></forbid-icon>
              </el-icon>
              <span>未经作者授权，禁止转载</span>
            </p>
          </div>
          <!-- 点赞收藏等数据 -->
          <div class="video-toolbar">
            <div class="toolbar-left" v-if="!isPGCPage">
              <archive-info v-if="videoInfo" :vid="videoInfo.vid" :short-id="videoInfo.shortId"></archive-info>
            </div>
            <div class="toolbar-right">
              <span>{{ onlineCount }} 人在看</span>
              <span>{{ videoInfo?.clicks }} 播放</span>
              <span>{{ videoInfo ? formatTime(videoInfo.createdAt) : '' }}</span>
            </div>
          </div>
          <div class="pgc-info-card" v-if="isPGCPage && pgcInfo">
            <img class="pgc-cover" :src="getResourceUrl(pgcInfo.cover)" alt="封面" />
            <div class="pgc-meta">
              <div class="pgc-name">{{ pgcInfo.title }}</div>
              <div class="pgc-sub">
                <span v-if="pgcInfo.year">{{ pgcInfo.year }}</span>
                <span v-if="pgcInfo.area"> · {{ formatAreaName(pgcInfo.area) }}</span>
                <span v-if="pgcInfo.current_episodes"> · 全{{ pgcInfo.current_episodes }}话</span>
              </div>
              <div class="pgc-rating" v-if="pgcInfo.rating">评分 {{ pgcInfo.rating }}</div>
              <div class="pgc-desc">{{ pgcInfo.desc || '暂无简介' }}</div>
            </div>
          </div>
          <!-- 简介部分 -->
          <div class="video-desc-container" v-if="!isPGCPage">
            <div ref="descRef" class="basic-desc-info" :style="`height: ${foldDesc ? foldDescHeight : 'auto'};`">
              <span class="desc-info-text">{{ videoInfo?.desc }}</span>
            </div>
            <div class="toggle-btn" v-show="showFoldBtn" @click="foldDesc = !foldDesc">
              <span class="toggle-btn-text">{{ foldDesc ? '展开更多' : '收起' }}</span>
            </div>
          </div>
          <!-- 标签部分 -->
          <div class="tags-box">
            <div class="tag" v-for="item in videoTagList" :key="item">{{ item }}</div>
          </div>
          <!-- 评论区 -->
          <comment-list v-if="videoInfo" :vid="videoInfo.vid" :short-id="videoInfo.shortId" @seek-time="handleSeekTime"></comment-list>
        </div>
        <div class="right-column" :class="{ 'pgc-mode': isPGCPage }">
<!-- 作者信息 -->
          <author-card v-if="videoInfo && !isPGCPage" :info="videoInfo.author"></author-card>
          <!-- 添加弹幕列表 -->
          <div class="danmaku-list-container">
            <danmaku-list ref="danmakuListRef" :height="danmakuListHeight"></danmaku-list>
          </div>
<!-- 合并的分P和合集列表 / PGC正片列表 -->
          <PGCSeasonPanel
            ref="collectionRef"
            v-if="videoInfo && isPGCPage"
            :vid="videoInfo.shortId || videoInfo.vid"
            :initial-seasons="pgcPanel.seasons"
            :initial-episodes="pgcPanel.episodes"
            :initial-active-season-id="pgcPanel.activeSeasonId"
          ></PGCSeasonPanel>
<video-collection 
            ref="collectionRef" 
            v-else-if="videoInfo" 
            :vid="videoInfo.shortId || videoInfo.vid"
            :resources="videoInfo.resources"
            :current-part="currentPart"
            @change-part="changePart"
          ></video-collection>
<!-- 相关推荐 -->
          <PGCRecommendList
            ref="recommendListRef"
            v-if="videoInfo && isPGCPage"
            :vid="videoInfo.shortId || videoInfo.vid"
          ></PGCRecommendList>
          <recommend-list
            ref="recommendListRef"
            v-else-if="videoInfo"
            :vid="videoInfo.shortId || videoInfo.vid"
            :show-autoplay-control="!videoInfo || (videoInfo.resources.length <= 1 && !hasCollection)"
          ></recommend-list>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick, type ComponentPublicInstance } from "vue";
import { ElIcon } from "element-plus";
import { Forbid as ForbidIcon } from "@icon-park/vue-next";
import { formatTime } from "@/utils/format";
import PartList from "./video/components/PartList.vue";
import AuthorCard from './video/components/AuthorCard.vue';
import ArchiveInfo from './video/components/ArchiveInfo.vue';
import VideoCollection from "./video/components/VideoCollection.vue";
import PGCSeasonPanel from "./video/components/PGCSeasonPanel.vue";
import CommentList from "./video/components/CommentList.vue";
import DanmakuList from "./video/components/DanmakuList.vue";
import HeaderBar from "@/components/header-bar/index.vue";
import VideoPlayer from "@/components/video-player/index.vue";
import RecommendList from "./video/components/RecommendList.vue";
import PGCRecommendList from "./video/components/PGCRecommendList.vue";
import { asyncGetVideoInfoAPI } from "@/api/video";
import { getPGCPlayPanelByVideoAPI } from "@/api/pgc";
import { resolveWatchVideoIdForInitialLoad, resolveWatchVideoIdOnQueryChange } from "@/utils/watch-route";
import { normalizeVideoTags } from "@/utils/video-tags";
import { getResourceUrl } from "@/utils/resource";
import { createUUID } from "@/utils/uuid";
import { getDanmakuAPI } from "@/api/danmaku";
import { getHistoryProgressAPI, addHistoryAPI } from "@/api/history";
import { globalConfig } from '@/utils/global-config';
import { updateTokenAPI } from '@/api/auth';
import { storageData } from '@/utils/storage-data';
import { useAuthStore, saveCredentials } from '@/stores/auth-store';
import { statusCode } from '@/utils/status-code';

const route = useRoute();
const router = useRouter();

// 路由兼容：/watch?ep= 或 /watch?v=（逻辑见 utils/watch-route）
const videoId = await resolveWatchVideoIdForInitialLoad(route.query);

// 获取视频信息
const videoInfo = ref<VideoType>();
const videoTagList = computed(() => normalizeVideoTags(videoInfo.value?.tags));
const pgcInfo = ref<any>(null);
const pgcPanel = ref<{ seasons: any[]; episodes: any[]; activeSeasonId: string }>({
  seasons: [],
  episodes: [],
  activeSeasonId: '',
});
const areaCodeMap: Record<string, string> = {
  CN: '中国大陆',
  JP: '日本',
  HK: '中国香港',
  TW: '中国台湾',
  KR: '韩国',
  US: '美国',
};
const formatAreaName = (raw: unknown) => {
  const code = String(raw ?? '').trim().toUpperCase();
  return areaCodeMap[code] || String(raw ?? '').trim();
};
const isPGCBound = ref(false);
const routePGCMode = computed(() => {
  const raw = route.query.mode;
  const mode = Array.isArray(raw) ? raw[0] : raw;
  return String(mode || '').trim().toLowerCase() === 'pgc';
});
const isPGCPage = computed(() => routePGCMode.value);
const currentWatchVQuery = computed(() => {
  const fromRoute = Array.isArray(route.query.v) ? route.query.v[0] : route.query.v;
  if (fromRoute && String(fromRoute).trim() !== '') return String(fromRoute);
  return String(videoInfo.value?.shortId || videoInfo.value?.vid || '');
});

const { data } = await asyncGetVideoInfoAPI(videoId);
if ((data.value as any).code === statusCode.OK) {
  videoInfo.value = (data.value as any).data.video as VideoType;
} else {
  await navigateTo('/404');
  throw new Error('video not found');
}

const loadPGCBinding = async (vid: number | string) => {
  try {
    const res = await getPGCPlayPanelByVideoAPI(vid);
    if (res?.data?.code === statusCode.OK && res?.data?.data?.current) {
      isPGCBound.value = true;
      pgcInfo.value = res.data.data.current;
      pgcPanel.value = {
        seasons: res.data.data.seasons || [],
        episodes: res.data.data.episodes || [],
        activeSeasonId: String(res.data.data.active_season_id || ''),
      };
    } else {
      isPGCBound.value = false;
      pgcInfo.value = null;
      pgcPanel.value = { seasons: [], episodes: [], activeSeasonId: '' };
    }
  } catch {
    isPGCBound.value = false;
    pgcInfo.value = null;
    pgcPanel.value = { seasons: [], episodes: [], activeSeasonId: '' };
  }
};
if (videoInfo.value?.vid) {
  const vid = videoInfo.value.shortId || videoInfo.value.vid;
  await loadPGCBinding(vid);
}

const playerContainerRef = ref<HTMLElement | null>(null)
const danmakuListHeight = ref(300);
const playerRef = ref<ComponentPublicInstance<{
  seek: (time: number) => void;
  uploadHistory: () => void;
  setDanmaku: (data: any[]) => void;
  setOnReady: (cb: () => void) => void;
  setOnEnded: (cb: () => void) => void;
}> | null>(null);

const handelResize = () => {
  nextTick(() => {
    danmakuListHeight.value = ((playerContainerRef.value?.clientWidth || 730) * 0.5625) + 40 - 104;
  })
}

// 计算当前分P：rid 优先（不受分P排序影响），否则用 p，最后兜底 1
const resolveInitialPart = (): number => {
  if (route.query.rid) {
    const rid = String(route.query.rid);
    const idx = videoInfo.value?.resources.findIndex(r => r.shortId === rid) ?? -1;
    if (idx >= 0) return idx + 1;
  }
  return Number(route.query.p) || 1;
};
const currentPart = ref(resolveInitialPart());

// URL 冗余/非法参数清理，仅在 client 端做以避免 SSR 期间副作用
if (process.client) {
  if (route.query.rid && route.query.p) {
    const { p, ...queryWithoutP } = route.query;
    router.replace({ path: '/watch', query: queryWithoutP });
  } else if (route.query.p && !route.query.rid && Number(route.query.p) > videoInfo.value!.resources.length) {
    router.replace({ path: '/watch', query: { ...route.query, v: currentWatchVQuery.value, p: undefined } });
  }
}

const pendingProgress = ref<number | null>(null);

// 获取组件引用
const recommendListRef = ref<ComponentPublicInstance<{
  autonext?: boolean;
  getNextVideo?: () => any;
  resetPlayIndex?: (vid: number) => void;
}> | null>(null);
const partListRef = ref<InstanceType<typeof PartList> | null>(null);
const collectionRef = ref<ComponentPublicInstance<{
  autonext?: boolean;
  getNextVideo?: () => any;
  hasPlaylist?: boolean;
}> | null>(null);
const hasCollection = computed(() => !!collectionRef.value?.hasPlaylist);

// 视频播放结束时的自动连播逻辑
const onVideoEnded = () => {
  // PGC模式不执行自动连播
  if (isPGCPage.value) return;
  
  // 检查用户是否开启了自动连播
  if (localStorage.getItem('video-autonext') !== 'true') return;

  // 优先检查当前视频是否有分P
  const hasParts = (videoInfo.value?.resources?.length || 0) > 1
  const curPart = currentPart.value
  
  if (hasParts) {
    // 当前视频有分P，检查是否还有下一个分P
    if (curPart < (videoInfo.value?.resources?.length || 0)) {
      // 切换到下一个分P
      setTimeout(() => {
        changePart(curPart + 1)
      }, 1000)
      return
    }
  }
  
  // 没有分P或已是最后一个分P，检查合集自动连播
  const collectionRefVal = collectionRef.value
  if (collectionRefVal?.autonext) {
    const nextVideo = collectionRefVal.getNextVideo?.()
    if (nextVideo) {
      setTimeout(() => {
        const nextEp = Number((nextVideo as any).epId || 0);
        if (isPGCPage.value && Number.isFinite(nextEp) && nextEp > 0) {
          navigateTo(`/watch?ep=${Math.floor(nextEp)}&mode=pgc`);
          return;
        }
        const nextV = (nextVideo as any).shortId || String((nextVideo as any).vid);
        const pgcMode = isPGCPage.value ? '&mode=pgc' : '';
        const nextShortId = (nextVideo as any).shortId;
        if (nextShortId) {
          navigateTo(`/watch?v=${nextV}&rid=${nextShortId}${pgcMode}`)
        } else {
          navigateTo(`/watch?v=${nextV}${pgcMode}`)
        }
      }, 1000)
      return
    }
  }
  
  // 合集没有下一个或未开启，检查推荐
  checkRecommendAutoplay()
};

// 检查推荐视频自动连播
const checkRecommendAutoplay = () => {
  const recommendRef = recommendListRef.value
  if (!recommendRef?.autonext) return;
  const nextVideo = recommendRef.getNextVideo?.();
  if (!nextVideo) return;
  setTimeout(() => {
    const nextEp = Number((nextVideo as any).epId || 0);
    if (isPGCPage.value && Number.isFinite(nextEp) && nextEp > 0) {
      navigateTo(`/watch?ep=${Math.floor(nextEp)}&mode=pgc`);
      return;
    }
    const nextV = nextVideo.shortId || String(nextVideo.vid);
    const pgcMode = isPGCPage.value ? '&mode=pgc' : '';
    const nextShortId = nextVideo.shortId;
    if (nextShortId) {
      navigateTo(`/watch?v=${nextV}&rid=${nextShortId}${pgcMode}`)
    } else {
      navigateTo(`/watch?v=${nextV}${pgcMode}`)
    }
  }, 3000);
};

const onPlayerReady = () => {
  // 延迟执行 seek，确保播放器完全准备好
  setTimeout(() => {
    // 原有的进度恢复逻辑
    if (pendingProgress.value === -1 && playerRef.value && playerRef.value.seek) {
      playerRef.value.seek(0);
      pendingProgress.value = null;
      return;
    }
    if (pendingProgress.value !== null && playerRef.value && playerRef.value.seek) {
      playerRef.value.seek(pendingProgress.value);
      pendingProgress.value = null;
    }

    // 新增：绑定播放结束事件
    if (playerRef.value && playerRef.value.setOnEnded) {
      playerRef.value.setOnEnded(onVideoEnded);
    }
  }, 100);
};

watch(playerRef, (val) => {
  if (val && val.setOnReady) {
    val.setOnReady(onPlayerReady);
  }
});

// 处理评论区时间跳转
const handleSeekTime = (seconds: number) => {
  if (playerRef.value && playerRef.value.seek) {
    playerRef.value.seek(seconds);
  }
};

// 获取弹幕列表
const danmakuListRef = ref<InstanceType<typeof DanmakuList> | null>(null);

const getDanmakuList = async (vid: string | number, part?: number, rid?: string) => {
  const res = await getDanmakuAPI(vid, part, rid);
  if (res.data.code === statusCode.OK) {
    const danmakus = res.data.data.danmaku || [];
    nextTick(() => {
      playerRef.value?.setDanmaku(danmakus)
      danmakuListRef.value?.setDanmaku(danmakus)
    })
  }
};

// 加载某分P的进度与弹幕；rid 存在则优先走 rid 精准定位
const refreshProgressAndDanmaku = async (partNum: number) => {
  if (!videoInfo.value) return;
  const vid = videoInfo.value.shortId || videoInfo.value.vid;
  const rid = videoInfo.value.resources[partNum - 1]?.shortId;
  try {
    const res = rid
      ? await getHistoryProgressAPI(vid, undefined, rid)
      : await getHistoryProgressAPI(vid, partNum);
    const progress = res?.data?.code === 200 ? res.data.data?.progress : null;
    pendingProgress.value = typeof progress === 'number' && progress !== 0 ? progress : null;
  } catch {
    pendingProgress.value = null;
  }
  if (rid) {
    getDanmakuList(vid, undefined, rid);
  } else {
    getDanmakuList(vid, partNum);
  }
};

const changePart = async (target: number) => {
  if (!videoInfo.value?.resources[target - 1]) return;
  currentPart.value = target;
  const targetRid = videoInfo.value.resources[target - 1]?.shortId;
  if (targetRid) {
    router.replace({ path: '/watch', query: { ...route.query, rid: targetRid, p: undefined } });
  } else {
    router.replace({ path: '/watch', query: { ...route.query, v: currentWatchVQuery.value, p: currentPart.value } });
  }
  await refreshProgressAndDanmaku(target);
};

// 简介部分
const foldDesc = ref(true);
const descRef = ref<HTMLElement>();
const showPlayer = ref(false);
const showFoldBtn = ref(false);
const foldDescHeight = ref('auto');
const playerReady = ref(false);
onMounted(async () => {
  if (descRef.value && descRef.value.clientHeight >= 80) {
    showFoldBtn.value = true;
    foldDescHeight.value = '80px';
  } else {
    showFoldBtn.value = false;
    foldDescHeight.value = 'auto';
  }

  if (videoInfo.value) {
    await refreshProgressAndDanmaku(currentPart.value);
  }

  handelResize();
  window.addEventListener("resize", handelResize);

  nextTick(() => {
    showPlayer.value = true;
    playerReady.value = true;
  })

  initWebSocket();
  document.addEventListener('visibilitychange', handleVisibilityChange);
})

//websocket
const onlineCount = ref(0);
let SocketURL = "";
let websocket: WebSocket | null = null;
let reconnectTimer: number | null = null;
let reconnectAttempts = 0;
const MAX_RECONNECT_ATTEMPTS = 10;
let isManualClose = false;
let heartbeatTimer: number | null = null;
let lastMessageTime = 0;

const closeWebSocket = () => {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer);
    reconnectTimer = null;
  }
  if (heartbeatTimer) {
    clearInterval(heartbeatTimer);
    heartbeatTimer = null;
  }
  if (websocket) {
    isManualClose = true;
    websocket.close();
    websocket = null;
  }
  reconnectAttempts = 0;
  onlineCount.value = 0;
}

const reconnectWebSocket = () => {
  closeWebSocket();
  isManualClose = false;
  initWebSocket();
}

const initWebSocket = () => {
  let clientId = localStorage.getItem("ws-client-id");
  if (!clientId) {
    clientId = createUUID();
    localStorage.setItem("ws-client-id", clientId);
  }

  if (process.dev) {
    const wsProtocol = globalConfig.https ? 'wss://' : 'ws://';
    SocketURL = `${wsProtocol}${globalConfig.domain}/api/v1/online/video?vid=${videoId}&clientId=${clientId}`;
  } else {
    const wsProtocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
    SocketURL = `${wsProtocol}${window.location.host}/api/v1/online/video?vid=${videoId}&clientId=${clientId}`;
  }

  if (heartbeatTimer) {
    clearInterval(heartbeatTimer);
    heartbeatTimer = null;
  }

  try {
    console.log('[WebSocket] 开始连接:', SocketURL);
    websocket = new WebSocket(SocketURL);

    websocket.onopen = () => {
      console.log('[WebSocket] 连接成功');
      reconnectAttempts = 0;
      lastMessageTime = Date.now();
      startHeartbeat();
    };

    websocket.onmessage = websocketOnmessage;

    websocket.onerror = (error) => {
      console.error('[WebSocket] 连接错误:', error);
    };

    websocket.onclose = (event) => {
      console.log('[WebSocket] 连接关闭 - Code:', event.code, 'Reason:', event.reason, 'wasClean:', event.wasClean);
      websocket = null;

      if (heartbeatTimer) {
        clearInterval(heartbeatTimer);
        heartbeatTimer = null;
      }

      if (!isManualClose) {
        if (reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
          reconnectAttempts++;
          const delay = Math.min(1000 * Math.pow(2, reconnectAttempts - 1), 10000);
          console.log(`[WebSocket] ${delay}ms 后尝试第 ${reconnectAttempts} 次重连...`);
          reconnectTimer = window.setTimeout(() => {
            initWebSocket();
          }, delay);
        } else {
          console.error('[WebSocket] 已达到最大重连次数,停止重连');
        }
      }
    };
  } catch (error) {
    console.error('[WebSocket] 创建连接失败:', error);
  }
}

const startHeartbeat = () => {
  heartbeatTimer = window.setInterval(() => {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      try {
        websocket.send('ping');
        lastMessageTime = Date.now();
      } catch {
        console.warn('[WebSocket] 发送心跳失败');
        websocket.close();
      }
    } else if (websocket && websocket.readyState !== WebSocket.CONNECTING) {
      console.warn('[WebSocket] 检测到连接异常,状态:', websocket.readyState);
      websocket.close();
    }
  }, 25000);
}

const handleVisibilityChange = async () => {
  if (document.hidden) {
    console.log('[WebSocket] 页面进入后台');
  } else {
    console.log('[WebSocket] 页面回到前台');
    // 重连 WebSocket
    if (!websocket || websocket.readyState !== WebSocket.OPEN) {
      console.log('[WebSocket] 页面恢复时检测到连接断开,尝试重连');
      reconnectAttempts = 0;
      initWebSocket();
    }
    // 同步登录状态：watch 页的 request interceptor 跳过了 token 自动刷新，
    // 所以在页面回到前台时手动尝试续签并同步 auth 状态
    try {
      const localRefreshToken = storageData.get('refreshToken');
      const token = storageData.get('token');
      if (!token && localRefreshToken) {
        const res = await updateTokenAPI(localRefreshToken);
        if (res.data.code === statusCode.OK) {
          saveCredentials({
            token: res.data.data.token,
            refreshToken: res.data.data.refreshToken,
            userId: res.data.data.userId,
          });
        }
      }
      const auth = useAuthStore();
      await auth.fetchMe(true);
    } catch {
      // 不阻塞其他逻辑
    }
  }
}

const websocketOnmessage = (e: any) => {
  try {
    lastMessageTime = Date.now();
    const res = JSON.parse(e.data);
    if (typeof res.number === 'number') {
      onlineCount.value = res.number;
      console.log('[WebSocket] 更新在线人数:', res.number);
    }
  } catch (error) {
    console.error('[WebSocket] 解析消息失败:', error);
  }
}


onBeforeUnmount(() => {
  window.removeEventListener("resize", handelResize);
  document.removeEventListener('visibilitychange', handleVisibilityChange);
  closeWebSocket();
  playerRef.value = null;
  recommendListRef.value = null;
  partListRef.value = null;
  collectionRef.value = null;
  danmakuListRef.value = null;
  playerContainerRef.value = null;
  descRef.value = undefined;
})

// 监听 v/ep/p/rid 变化，重新拉取视频信息和重置状态
watch(
  () => [route.query.v, route.query.ep, route.query.p, route.query.rid],
  async ([newV, newEp, newP, newRid], [oldV, oldEp, oldP, oldRid]) => {
    // 视频主体未变（v/ep/rid 都没动），仅 p 变：不重拉 videoInfo，只切分P
    if (newV === oldV && newEp === oldEp && newRid === oldRid) {
      if (newP !== oldP && videoInfo.value) {
        const partNum = Number(newP) || 1;
        if (videoInfo.value.resources[partNum - 1]) {
          currentPart.value = partNum;
          await refreshProgressAndDanmaku(partNum);
        } else {
          router.replace({ path: '/watch', query: { ...route.query, v: currentWatchVQuery.value, p: 1 } });
        }
      }
      return;
    }
    const resolved = await resolveWatchVideoIdOnQueryChange(route.query);
    if (!resolved.ok) return;
    const newVideoId = resolved.videoId;
    const { data } = await asyncGetVideoInfoAPI(newVideoId);
    if ((data.value as any).code === statusCode.OK) {
      videoInfo.value = (data.value as any).data.video as VideoType;
      const vid = videoInfo.value.shortId || videoInfo.value.vid;
      await loadPGCBinding(vid);
      // 初始化分P
      if (route.query.rid) {
        const rid = String(route.query.rid);
        const idx = videoInfo.value.resources.findIndex(r => r.shortId === rid);
        currentPart.value = idx >= 0 ? idx + 1 : 1;
      } else {
        currentPart.value = Number(route.query.p) || 1;
      }
      await refreshProgressAndDanmaku(currentPart.value);
      if (process.client) {
        reconnectWebSocket();
      }
    } else {
      navigateTo('/404');
    }
  }
);

useHead({
  title: () => {
    const pageTitle = (isPGCPage.value ? (pgcInfo.value?.title || '') : '') || videoInfo.value?.title || ''
    return pageTitle ? `${pageTitle} - ${globalConfig.title}` : globalConfig.title
  }
})
</script>

<style lang="scss" scoped>
.header {
  position: fixed;
}

.video-main {
  padding-top: 80px;
  margin: 0 auto;
  min-width: 1200px;
}

.mian-content {
  display: flex;
  justify-content: center;
  width: 100%;
  max-width: calc(100% - 100px);
  margin: 0 auto;
  position: relative;
}

.left-column {
  flex: 1;
  max-width: 1200px;
  margin-top: 20px;

  .video-player {
    position: relative;
    margin: 0 auto;
    width: 100%;
    min-width: 680px;
    min-height: 382px;

    .skeleton {
      width: 100%;
      padding-bottom: 56.25%;
      background-color: var(--bg-elev-1);
      border: 1px solid var(--border-color);
      position: relative;
      overflow: hidden;
    }

    .skeleton::after {
      content: '';
      position: absolute;
      inset: 0;
      background: linear-gradient(90deg,
          transparent 0%,
          rgba(255, 255, 255, 0.06) 50%,
          transparent 100%);
      animation: skeleton-shimmer 1.2s infinite;
    }
  }

  .video-title-box {
    width: 100%;
    height: 54px;
    display: flex;

    .video-title {
      width: calc(100% - 160px);
      font-weight: 500;
      line-height: 28px;
      margin: 13px 0;
      font-size: 20px;
      color: var(--font-primary-1);
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }

    .copyright {
      width: 180px;
      display: flex;
      align-items: center;
      justify-content: flex-end;
      font-size: 13px;
      color: var(--font-primary-3);

      .icon {
        padding: 0 6px;
      }
    }
  }

  .video-toolbar {
    color: var(--font-primary-3);
    font-size: 13px;
    padding-bottom: 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid var(--border-color);

    .toolbar-right {
      display: inline-block;

      span {
        margin-left: 20px;
      }
    }
  }

  .video-desc-container {
    margin: 16px 0;

    .basic-desc-info {
      white-space: pre-line;
      letter-spacing: 0;
      color: var(--font-primary-1);
      font-size: 15px;
      line-height: 24px;
      overflow: hidden;

      .desc-info-text {
        white-space: pre-line;
      }
    }

    .toggle-btn {
      margin-top: 10px;
      font-size: 13px;
      line-height: 18px;

      .toggle-btn-text {
        cursor: pointer;
        color: var(--font-primary-2);

        &:hover {
          color: var(--primary-hover-color);
        }
      }
    }
  }

  .tags-box {
    padding-bottom: 6px;
    margin: 16px 0 20px 0;
    border-bottom: 1px solid var(--border-color);

    .tag {
      color: var(--font-primary-2);
      background: var(--border-color);
      height: 28px;
      line-height: 28px;
      border-radius: 14px;
      font-size: 13px;
      padding: 0 12px;
      box-sizing: border-box;
      transition: all .3s;
      display: inline-flex;
      align-items: center;
      cursor: pointer;
      margin: 0 12px 8px 0;

      &:hover {
        background: var(--hover-bg);
        color: var(--font-primary-1);
      }
    }
  }

  .pgc-info-card {
    margin: 16px 0;
    padding: 14px;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    background: var(--bg-elev-1);
    display: flex;
    gap: 12px;

    .pgc-cover {
      width: 120px;
      height: 160px;
      border-radius: 6px;
      object-fit: cover;
      flex-shrink: 0;
    }

    .pgc-meta {
      min-width: 0;

      .pgc-name {
        font-size: 18px;
        color: var(--font-primary-1);
        margin-bottom: 6px;
      }

      .pgc-sub {
        font-size: 13px;
        color: var(--font-primary-3);
        margin-bottom: 8px;
      }

      .pgc-rating {
        color: #f5a623;
        font-size: 16px;
        margin-bottom: 8px;
      }

      .pgc-desc {
        color: var(--font-primary-2);
        font-size: 14px;
        line-height: 22px;
        white-space: pre-line;
      }
    }
  }
}

@keyframes skeleton-shimmer {
  0% {
    transform: translateX(-100%);
  }

  100% {
    transform: translateX(100%);
  }
}

.right-column {
  width: 340px;
  margin-left: 30px;
  z-index: 1;

  .danmaku-list-container {
    margin-bottom: 18px;
  }

  &.pgc-mode {
    margin-top: 20px;
  }
}
</style>
