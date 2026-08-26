<template>
  <div class="video">
    <header-bar class="header"></header-bar>
    <div class="video-main">
      <div class="main-content">
        <div class="left-column">
          <div class="video-player" ref="playerContainerRef">
<client-only>
              <video-player v-if="videoInfo && playerReady" ref="playerRef" :video-info="videoInfo" :part="currentPart"
                :progress="pendingProgress" @danmaku-sent="handleDanmakuSent"
                :episode-picker-list="episodePickerList" :episode-picker-active-index="episodePickerActiveIndex"
                :episode-picker-type="episodePickerType" @episode-pick="handleEpisodePick"></video-player>
            </client-only>
            <div v-if="!showPlayer" class="skeleton"></div>
          </div>
          <!-- 标题和版权信息 -->
          <div class="video-title-box">
            <p class="video-title">{{ pgcInfo?.title || videoInfo?.title }}</p>
            <p v-show="videoInfo.copyright === 1" class="copyright">
              <el-icon class="icon" color='#fd6d6f'>
                <forbid-icon></forbid-icon>
              </el-icon>
              <span>未经作者授权，禁止转载</span>
            </p>
          </div>
          <!-- 点赞收藏等数据 -->
          <div class="video-toolbar">
            <div class="toolbar-left" v-if="!isPGCPage">
              <archive-info v-if="videoInfo" :vid="videoInfo.vid" :short-id="videoInfo.shortId" :rid="videoInfo.resources?.[currentPart - 1]?.shortId"></archive-info>
            </div>
            <div class="toolbar-right">
              <span>{{ onlineCount }} 人在看</span>
              <span>{{ videoInfo?.clicks }} 播放</span>
              <span>{{ videoInfo ? formatTime(videoInfo.createdAt) : '' }}</span>
            </div>
          </div>
          <div class="pgc-info-card" v-if="isPGCPage && pgcInfo">
            <oss-image class="pgc-cover" :src="pgcInfo.cover" alt="封面" />
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
            <danmaku-list ref="danmakuListRef" :height="danmakuListHeight" @seek-time="handleSeekTime"></danmaku-list>
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
            :show-autoplay-control="(videoInfo.resources?.length ?? 0) <= 1 && !hasCollection"
          ></recommend-list>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, unref, onMounted, onBeforeUnmount, watch, nextTick, type ComponentPublicInstance } from "vue";
import { ElIcon } from "element-plus";
import { Forbid as ForbidIcon } from "@icon-park/vue-next";
import { formatTime } from "@/utils/format";
import AuthorCard from './video/components/AuthorCard.vue';
import ArchiveInfo from './video/components/ArchiveInfo.vue';
import VideoCollection from "./video/components/VideoCollection.vue";
import PGCSeasonPanel from "./video/components/PGCSeasonPanel.vue";
import CommentList from "./video/components/CommentList.vue";
import DanmakuList from "./video/components/DanmakuList.vue";

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
import { statusCode } from '@/utils/status-code';
import { useAuthStore } from '@/stores/auth-store';

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

const { data: videoApiData } = await asyncGetVideoInfoAPI(videoId);
if (videoApiData.value) {
  if (videoApiData.value.code === statusCode.OK) {
    videoInfo.value = videoApiData.value.data.video as VideoType;
  } else {
    if (process.client) {
      await navigateTo('/404');
    }
    throw new Error('video not found');
  }
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
// PGC 绑定：放在 watch 里客户端执行，避免 SSR 水合不匹配
watch(videoInfo, (val) => {
  if (val?.vid && process.client) {
    const vid = val.shortId || val.vid;
    loadPGCBinding(vid);
  }
}, { immediate: true });

const playerContainerRef = ref<HTMLElement | null>(null)
const danmakuListHeight = ref(300);
const playerRef = ref<ComponentPublicInstance<{
  seek: (time: number) => void;
  uploadHistory: () => void;
  setDanmaku: (data: any[]) => void;
  setOnReady: (cb: () => void) => void;
  setOnEnded: (cb: () => void) => void;
  getCurrentTime: () => number;
  getDuration: () => number;
}> | null>(null);

const handelResize = () => {
  nextTick(() => {
    danmakuListHeight.value = ((playerContainerRef.value?.clientWidth || 730) * 0.5625) + 40 - 104;
  })
}

// 非历史页进入且有多分P时，查播放历史恢复上次观看的分P
let historyPartOverride = 0;
if (videoInfo.value && videoInfo.value.resources.length > 1 && !route.query.rid && !route.query.p && process.client) {
  try {
    const res = await getHistoryProgressAPI(videoInfo.value.shortId || videoInfo.value.vid);
    if (res?.data?.code === 200) {
      const hp = res.data.data?.part;
      if (hp && hp > 1 && hp <= videoInfo.value.resources.length) {
        historyPartOverride = hp;
      }
    }
  } catch { /* 忽略，兜底 P1 */ }
}

// 计算当前分P：rid 优先（不受分P排序影响），其次 p，其次播放历史，最后兜底 1
const resolveInitialPart = (): number => {
  if (route.query.rid) {
    const rid = String(route.query.rid);
    const idx = videoInfo.value?.resources.findIndex(r => r.shortId === rid) ?? -1;
    if (idx >= 0) return idx + 1;
  }
  if (route.query.p) return Number(route.query.p);
  if (historyPartOverride > 1) return historyPartOverride;
  return 1;
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

// 同步当前分P到 URL（使用 rid，其是稳定的资源标识符）。
// 解决场景：用户从历史记录恢复播放非 P1 时，地址栏仍只有 ?v=xxx，分享链接只能回到 P1。
if (process.client) {
  const targetPart = currentPart.value;
  const resources = videoInfo.value?.resources;
  const currentResource = resources?.[targetPart - 1];
  // 单分P情况不自动添加 rid 到地址栏，避免冲路径
  const isSinglePart = resources && resources.length <= 1;
  if (currentResource?.shortId && !route.query.rid && !isSinglePart) {
    // 使用 history.replaceState 避免触发 Vue Router 的 watcher 导致重复拉取视频信息
    const newQuery: Record<string, string> = { v: currentWatchVQuery.value, rid: currentResource.shortId };
    if (route.query.ep) newQuery.ep = String(route.query.ep);
    if (route.query.mode) newQuery.mode = String(route.query.mode);
    const search = Object.entries(newQuery)
      .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
      .join('&');
    window.history.replaceState(null, '', `/watch?${search}`);
  }
}

const pendingProgress = ref<number | null>(null);

// 获取组件引用
const recommendListRef = ref<ComponentPublicInstance<{
  autonext?: boolean;
  getNextVideo?: () => any;
  resetPlayIndex?: (vid: number) => void;
}> | null>(null);
const collectionRef = ref<ComponentPublicInstance<{
  autonext?: boolean;
  getNextVideo?: () => any;
  hasPlaylist?: boolean;
  mergedList?: any[];
  listType?: string;
  currentIndex?: number;
  orderedEpisodes?: any[];
  isActiveEpisode?: (ep: any) => boolean;
}> | null>(null);
const hasCollection = computed(() => !!collectionRef.value?.hasPlaylist);

// ===== 全屏选集数据 =====
const episodePickerType = computed<'none' | 'parts' | 'collection' | 'pgc'>(() => {
  if (isPGCPage.value) return 'pgc';
  const listType = collectionRef.value?.listType;
  if (listType === 'collection') return 'collection';
  if (videoInfo.value?.resources && videoInfo.value.resources.length > 1) return 'parts';
  return 'none';
});

const episodePickerList = computed(() => {
  if (!videoInfo.value) return [];
  if (episodePickerType.value === 'none') return [];

  if (episodePickerType.value === 'pgc') {
    const eps = collectionRef.value?.orderedEpisodes;
    if (!eps || eps.length <= 1) return [];
    return eps.map((ep: any) => ({
      label: `第${ep.episode_number}话${ep.title ? ' ' + ep.title : ''}`,
      index: ep.episode_number,
      vid: ep.vid,
      epId: Number(ep.ep_id || ep.id || 0) || undefined,
    }));
  }

  if (episodePickerType.value === 'collection') {
    const merged = collectionRef.value?.mergedList;
    if (!merged || merged.length <= 1) return [];
    return merged.map((item: any, idx: number) => ({
      label: item.partTitle || item.title || `P${item.p || idx + 1}`,
      index: idx + 1,
      vid: item.vid,
      part: item.p || 1,
      rid: item.resourceRid,
    }));
  }

  // parts mode
  const resources = videoInfo.value.resources;
  if (!resources || resources.length <= 1) return [];
  return resources.map((r: any, idx: number) => ({
    label: r.title || `P${idx + 1}`,
    index: idx + 1,
    part: idx + 1,
    rid: r.shortId,
  }));
});

const episodePickerActiveIndex = computed(() => {
  const type = episodePickerType.value;
  if (type === 'none') return 0;

  if (type === 'pgc') {
    const eps = collectionRef.value?.orderedEpisodes;
    if (!eps) return 0;
    const active = eps.find((ep: any) => {
      const check = collectionRef.value?.isActiveEpisode;
      return check ? check(ep) : false;
    });
    return active?.episode_number || 0;
  }

  if (type === 'collection') {
    return collectionRef.value?.currentIndex || currentPart.value;
  }

  return currentPart.value;
});

const handleEpisodePick = (item: { vid?: string; part?: number; rid?: string; epId?: number }) => {
  if (isPGCPage.value && item.epId) {
    void navigateTo(`/watch?ep=${item.epId}&mode=pgc`);
    return;
  }
  if (isPGCPage.value && item.vid) {
    void navigateTo(`/watch?v=${item.vid}&mode=pgc`);
    return;
  }
  // 同一稿件不同分P
  if (!item.vid && item.part) {
    changePart(item.part);
    return;
  }
  // 合集不同稿件
  if (item.vid) {
    const v = String(item.vid);
    if (item.rid) {
      void navigateTo(`/watch?v=${v}&rid=${encodeURIComponent(item.rid)}`);
    } else if (item.part && item.part > 1) {
      void navigateTo(`/watch?v=${v}&p=${item.part}`);
    } else {
      void navigateTo(`/watch?v=${v}`);
    }
  }
};

/** 合集与推荐共用 useVideoAutonextFollow；defineExpose 的 autonext 可能是 Ref / ComputedRef，任一侧有实例即可，只 unref 一次 */
const isFollowAutonextOn = () => {
  const inst = collectionRef.value ?? recommendListRef.value;
  return Boolean(inst && unref((inst as { autonext?: unknown }).autonext as any));
};

// ===== 自动连播倒计时提示 =====
const advanceCancelled = ref(false);
const advancePending = ref<{ label: string; countdown: number; handler: () => void } | null>(null);
let advanceToastEl: HTMLElement | null = null;
let advanceCountTimer: number | null = null;
let nearEndPollTimer: number | null = null;

/** 用户已取消推荐列表连播（本次会话内不再弹出） */
const recommendCancelled = ref(false);

/** 计算下一个自动连播项 */
const computeNextAdvance = (): { label: string; handler: () => void } | null => {
  if (isPGCPage.value) return null;

  // Bug 02: 循环播放开启时不显示连播倒计时，也不自动跳到下一集
  if (playerRef.value?.isLoopEnabled?.()) return null;

  const partCount = videoInfo.value?.resources?.length || 0;
  const hasParts = partCount > 1;
  const curPart = currentPart.value;

  if (hasParts && localStorage.getItem('video-autonext-parts') === 'true') {
    if (curPart < partCount) {
      const nextPart = curPart + 1;
      const label = videoInfo.value?.resources[nextPart - 1]?.title || `P${nextPart}`;
      return { label, handler: () => changePart(nextPart) };
    }
  }

  const cv = collectionRef.value;
  if (isFollowAutonextOn() && cv) {
    const nextVideo = (cv as any).getNextVideo?.();
    const cur = videoInfo.value;
    if (nextVideo) {
      console.log('[autonext] computeNextAdvance collection', {
        currentVid: cur?.vid,
        currentShortId: cur?.shortId,
        currentTitle: cur?.title,
        nextVid: nextVideo.vid,
        nextShortId: nextVideo.shortId,
        nextResourceRid: nextVideo.resourceRid,
        nextTitle: nextVideo.title,
        nextPartTitle: nextVideo.partTitle,
        nextP: nextVideo.p,
      });
      const label = nextVideo.partTitle || nextVideo.title || `P${nextVideo.p || 1}`;
      return { label, handler: () => scheduleNavigateToWatchNext(nextVideo, 0) };
    }

    // 没有下一项（已到列表末尾）
    console.log('[autonext] computeNextAdvance collection: no more items', {
      currentVid: cur?.vid,
      currentShortId: cur?.shortId,
      currentTitle: cur?.title,
    });
  }

  // 合集没有下一个，检查推荐列表
  const rv = recommendListRef.value;
  if (isFollowAutonextOn() && rv) {
    const nextVideo = rv.getNextVideo?.();
    if (nextVideo) {
      const item = nextVideo as any;
      return { label: item.title || '推荐视频', handler: () => scheduleNavigateToWatchNext(item, 0) };
    }
  }

  return null;
};

/** 注入/更新 Toast 到 #dplayer（全屏可见） */
const showAdvanceToast = (label: string, countdown: number) => {
  const player = document.getElementById('dplayer');
  if (!player) return;
  hideAdvanceToast();

  const toast = document.createElement('div');
  toast.className = 'wplayer-advance-toast';
  toast.style.cssText = [
    'position:absolute;bottom:50px;left:50%;transform:translateX(-50%);z-index:9999;',
    'background:rgba(0,0,0,0.85);border-radius:6px;padding:10px 20px;',
    'display:flex;align-items:center;gap:14px;font-size:14px;color:#fff;white-space:nowrap;',
    'backdrop-filter:blur(8px);border:1px solid rgba(255,255,255,0.08);',
    'box-shadow:0 4px 20px rgba(0,0,0,0.5);',
  ].join('');

  toast.innerHTML = `
    <span style="color:#aaa;flex-shrink:0;">即将播放</span>
    <span class="wplayer-advance-label" style="max-width:260px;overflow:hidden;text-overflow:ellipsis;flex-shrink:1;"></span>
    <span class="wplayer-advance-countdown" style="color:var(--wplayer-theme,#00a1d6);font-weight:700;min-width:28px;text-align:center;font-size:16px;flex-shrink:0;">${countdown}</span>
    <span style="color:#888;font-size:12px;flex-shrink:0;">秒</span>
    <button class="wplayer-advance-cancel-btn" style="background:rgba(255,255,255,0.08);border:none;border-radius:4px;color:#ddd;padding:6px 14px;cursor:pointer;font-size:13px;transition:background 0.15s;flex-shrink:0;">取消</button>
  `;
  // textContent 安全设置（防 XSS），同时确保 title 属性一致
  const labelSpan = toast.querySelector('.wplayer-advance-label')! as HTMLElement;
  labelSpan.textContent = label;
  labelSpan.title = label;

  toast.querySelector('.wplayer-advance-cancel-btn')!.addEventListener('click', (e) => {
    e.stopPropagation();
    advanceCancelled.value = true;
    hideAdvanceToast();
  });

  // Hover effect on cancel button
  const cancelBtn = toast.querySelector('.wplayer-advance-cancel-btn')! as HTMLElement;
  cancelBtn.addEventListener('mouseenter', () => { cancelBtn.style.background = 'rgba(255,255,255,0.15)'; });
  cancelBtn.addEventListener('mouseleave', () => { cancelBtn.style.background = 'rgba(255,255,255,0.08)'; });

  player.appendChild(toast);
  advanceToastEl = toast;
};

const hideAdvanceToast = () => {
  if (advanceToastEl && advanceToastEl.parentNode) {
    advanceToastEl.parentNode.removeChild(advanceToastEl);
  }
  advanceToastEl = null;
  if (advanceCountTimer !== null) {
    clearInterval(advanceCountTimer);
    advanceCountTimer = null;
  }
};

/** 启动轮询检测接近片尾 */
const startNearEndPolling = () => {
  stopNearEndPolling();
  nearEndPollTimer = window.setInterval(() => {
    const p = playerRef.value as any;
    if (!p || typeof p.getCurrentTime !== 'function') return;
    const ct = p.getCurrentTime();
    const dur = p.getDuration();
    if (!dur || dur <= 0 || ct == null || ct <= 0) return;
    const remaining = dur - ct;

    if (advanceCancelled.value) {
      hideAdvanceToast();
      return;
    }

    // 已存在待定连播 → 更新倒计时数字
    if (advancePending.value) {
      const cd = Math.ceil(Math.max(0, remaining));
      if (cd <= 2) {
        advancePending.value.countdown = cd;
        if (advanceToastEl) {
          const cdSpan = advanceToastEl.querySelector('.wplayer-advance-countdown');
          if (cdSpan) cdSpan.textContent = String(cd);
        }
      }
      if (remaining <= 0) {
        // ended 事件应已触发，清理兜底
        hideAdvanceToast();
        advancePending.value = null;
      }
      return;
    }

    // 接近片尾（剩余 <= 2s）
    if (remaining <= 2 && remaining > 0) {
      // 用户已取消推荐连播，本次会话不再弹出
      if (recommendCancelled.value) return;

      const next = computeNextAdvance();
      if (next) {
        advanceCancelled.value = false;
        advancePending.value = { label: next.label, countdown: Math.ceil(remaining), handler: next.handler };
        showAdvanceToast(next.label, Math.ceil(remaining));
      }
    }
  }, 500);
};

const stopNearEndPolling = () => {
  if (nearEndPollTimer !== null) {
    clearInterval(nearEndPollTimer);
    nearEndPollTimer = null;
  }
};

const resetAdvanceState = () => {
  advanceCancelled.value = false;
  advancePending.value = null;
  hideAdvanceToast();
};

/** 侧栏/推荐「下一则」：与 VideoCollection.goVideo 一致：多分 P 用 resourceRid 作查询 rid，勿把稿件 shortId 当成资源 rid */
function scheduleNavigateToWatchNext(
  item: { shortId?: string; vid?: number | string; resourceRid?: string; p?: number },
  delayMs: number,
) {
  const v = String(item.shortId ?? item.vid ?? '').trim();
  if (!v) {
    console.log('[autonext] scheduleNavigateToWatchNext: empty vid, abort');
    return;
  }
  const resourceRid = item.resourceRid;
  const p = item.p;
  window.setTimeout(() => {
    let url: string;
    if (resourceRid) {
      url = `/watch?v=${v}&rid=${encodeURIComponent(resourceRid)}`;
    } else if (typeof p === 'number' && p > 1) {
      url = `/watch?v=${v}&p=${p}`;
    } else {
      url = `/watch?v=${v}`;
    }
    console.log('[autonext] navigateTo', { url, fromVid: v, resourceRid, p });
    void navigateTo(url);
  }, delayMs);
}

// 视频播放结束时的自动连播逻辑（由 advance 倒计时驱动，结束时立即切换）
const onVideoEnded = () => {
  // PGC模式不执行自动连播
  if (isPGCPage.value) return;

  // 用户取消了连播
  if (advanceCancelled.value) {
    // 检测是否推荐列表场景（无分P且无合集下一项）→ 本次会话不再弹出
    const cv = collectionRef.value;
    const noMoreParts = !((videoInfo.value?.resources?.length || 0) > 1 && currentPart.value < (videoInfo.value?.resources?.length || 0));
    const noMoreCollection = !(cv?.hasPlaylist && cv.getNextVideo?.());
    if (noMoreParts && noMoreCollection) {
      recommendCancelled.value = true;
    }
    resetAdvanceState();
    return;
  }

  // 清除 polling 阶段缓存的 advancePending（仅作 UI 倒计时用），
  // 导航目标始终在结束时重新计算，避免闭包捕获与数据过渡期的时间差
  advancePending.value = null;
  hideAdvanceToast();

  // 分P续播（单P合集不命中）
  const partCount = videoInfo.value?.resources?.length || 0;
  const hasParts = partCount > 1;
  const curPart = currentPart.value;

  if (hasParts && localStorage.getItem('video-autonext-parts') === 'true') {
    if (curPart < partCount) {
      changePart(curPart + 1);
      return;
    }
  }

  // 合集续播
  const collectionRefVal = collectionRef.value;
  if (isFollowAutonextOn() && collectionRefVal) {
    const nextVideo = (collectionRefVal as any).getNextVideo?.();
    const cur = videoInfo.value;
    console.log('[autonext] onVideoEnded collection', {
      currentVid: cur?.vid,
      currentShortId: cur?.shortId,
      currentTitle: cur?.title,
      hasAdvancePending: !!advancePending.value,
    });
    if (nextVideo) {
      console.log('[autonext] navigating to next', {
        vid: nextVideo.vid,
        shortId: nextVideo.shortId,
        resourceRid: nextVideo.resourceRid,
        title: nextVideo.title,
        partTitle: nextVideo.partTitle,
      });
      scheduleNavigateToWatchNext(nextVideo, 0);
      return;
    }
    console.log('[autonext] onVideoEnded collection: no next video');
  }

  // 推荐列表续播
  const rv = recommendListRef.value;
  if (isFollowAutonextOn() && rv) {
    const nextVideo = rv.getNextVideo?.();
    if (nextVideo) {
      scheduleNavigateToWatchNext(nextVideo as { shortId?: string; vid?: number | string }, 0);
    }
  }
};

const onPlayerReady = () => {
  // 进度恢复由 video-player 内部通过 props.progress + pendingSeek + loadedmetadata/canplay 统一处理
  // 这里只做"已看完即重新开始"的场景兜底 + 绑定播放结束事件
  if (pendingProgress.value === -1 && playerRef.value?.seek) {
    playerRef.value.seek(0);
  }
  pendingProgress.value = null;
  if (playerRef.value?.setOnEnded) {
    playerRef.value.setOnEnded(onVideoEnded);
  }
  // 启动片尾倒计时轮询
  startNearEndPolling();
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

// 弹幕发送成功后无需额外处理：server 会通过 ws 回广播自己发的弹幕，走 websocketOnmessage 分支单条插入
const handleDanmakuSent = () => {};

// 加载某分P的进度与弹幕；rid 存在则优先走 rid 精准定位
const refreshProgressAndDanmaku = async (partNum: number) => {
  if (!videoInfo.value) return;
  const vid = videoInfo.value.shortId || videoInfo.value.vid;
  const rid = videoInfo.value.resources[partNum - 1]?.shortId;
  try {
    const res = rid
      ? await getHistoryProgressAPI(vid, partNum, rid)
      : await getHistoryProgressAPI(vid, partNum);
    const progress = res?.data?.code === 200 ? res.data.data?.progress : null;
    // -1 = 已看完：不做续播（由 onPlayerReady 兜底 seek(0)）；其他 0 或非数字也视为无续播
    if (progress === -1) {
      pendingProgress.value = -1;
    } else if (typeof progress === 'number' && progress > 0) {
      pendingProgress.value = progress;
    } else {
      pendingProgress.value = null;
    }
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
  resetAdvanceState();
  currentPart.value = target;
  const targetRid = videoInfo.value.resources[target - 1]?.shortId;
  if (targetRid) {
    router.replace({ path: '/watch', query: { ...route.query, rid: targetRid, p: undefined } });
  } else {
    router.replace({ path: '/watch', query: { ...route.query, v: currentWatchVQuery.value, p: currentPart.value } });
  }
  await refreshProgressAndDanmaku(target);
  // 切换分P后重新连接WebSocket以使用新的rid
  reconnectWebSocket();
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

  // 实时读取当前 vid，避免 SPA 路由切换后用到初始 setup 阶段捕获的旧值
  const currentVid = videoInfo.value?.shortId || videoInfo.value?.vid || videoId;
  if (!currentVid) {
    console.warn('[WebSocket] 当前 vid 为空，放弃建立连接');
    return;
  }
  // 获取当前分P的rid
  const rid = videoInfo.value?.resources?.[currentPart.value - 1]?.shortId || '';
  const ridParam = rid ? `&rid=${rid}` : '';

  const wsProtocol = globalConfig.https ? 'wss://' : 'ws://';
  SocketURL = `${wsProtocol}${globalConfig.domain}/api/v1/online/video?vid=${currentVid}&clientId=${clientId}${ridParam}`;

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
    // 同步登录状态：通过 HttpOnly Cookie 续签（不再依赖 localStorage token）
    try {
      const auth = useAuthStore();
      if (!auth.token) {
        // 无内存 token 时直接由 fetchMe 通过 Cookie 获取新 token
        await auth.fetchMe(true);
      }
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
    // 处理弹幕消息：单条插入，不再全量刷新列表
    if (res.type === 'danmaku' && res.danmaku) {
      const currentRid = videoInfo.value?.resources?.[currentPart.value - 1]?.shortId;
      const danmakuWithMeta = {
        ...res.danmaku,
        vid: videoInfo.value?.vid,
        part: currentPart.value,
        rid: currentRid
      };
      playerRef.value?.addDanmaku(danmakuWithMeta);
      danmakuListRef.value?.addDanmaku(danmakuWithMeta);
    }
  } catch (error) {
    console.error('[WebSocket] 解析消息失败:', error);
  }
}


onBeforeUnmount(() => {
  window.removeEventListener("resize", handelResize);
  document.removeEventListener('visibilitychange', handleVisibilityChange);
  closeWebSocket();
  stopNearEndPolling();
  hideAdvanceToast();
  playerRef.value = null;
  recommendListRef.value = null;
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
    resetAdvanceState();
    const resolved = await resolveWatchVideoIdOnQueryChange(route.query);
    if (!resolved.ok) return;
    const newVideoId = resolved.videoId;
    const { data } = await asyncGetVideoInfoAPI(newVideoId);
    if ((data.value as any).code === statusCode.OK) {
      videoInfo.value = (data.value as any).data.video as VideoType;
      const vid = videoInfo.value.shortId || videoInfo.value.vid;
      // 先设置 currentPart，确保与 videoInfo 在同一微任务批次内完成，
      // 避免播放器先用旧 currentPart 渲染一次导致错加载音频/进度
      if (route.query.rid) {
        const rid = String(route.query.rid);
        const idx = videoInfo.value.resources.findIndex(r => r.shortId === rid);
        currentPart.value = idx >= 0 ? idx + 1 : 1;
      } else {
        currentPart.value = Number(route.query.p) || 1;
      }
      await loadPGCBinding(vid);
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

.main-content {
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
