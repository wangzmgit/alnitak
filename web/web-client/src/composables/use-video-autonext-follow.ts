import { ref, watch, effectScope, onMounted } from 'vue';

/** 合集列表与右侧「接下来播放」共用：播完后是否自动接播侧栏队列中的下一则 */
export const VIDEO_AUTONEXT_FOLLOW_KEY = 'video-autonext-follow';

const LEGACY_COLLECTION = 'video-autonext-collection';
const LEGACY_RECOMMEND = 'video-autonext-recommend';

function readInitial(): boolean {
  if (typeof localStorage === 'undefined') return false;
  try {
    const unified = localStorage.getItem(VIDEO_AUTONEXT_FOLLOW_KEY);
    if (unified !== null) return unified === 'true';
    const legacyOn =
      localStorage.getItem(LEGACY_COLLECTION) === 'true' ||
      localStorage.getItem(LEGACY_RECOMMEND) === 'true';
    localStorage.setItem(VIDEO_AUTONEXT_FOLLOW_KEY, legacyOn ? 'true' : 'false');
    return legacyOn;
  } catch {
    return false;
  }
}

function persistFollow(v: boolean) {
  try {
    localStorage.setItem(VIDEO_AUTONEXT_FOLLOW_KEY, String(v));
  } catch {
    /* 隐私模式 / 配额等 */
  }
}

const followAutonext = ref(false);
const persistScope = effectScope(true);
let persistWatchRegistered = false;
/** 只从本地恢复一次；必须用 onMounted，不可 queueMicrotask（RecommendList 等 async setup 会在 await 期间先跑微任务，导致水合 class 错位） */
let followStorageSynced = false;

function syncFollowFromStorageOnce() {
  if (followStorageSynced) return;
  followStorageSynced = true;
  followAutonext.value = readInitial();
}

/**
 * 全局单例：同一页面内合集开关与推荐区开关始终保持一致。
 */
export function useVideoAutonextFollow() {
  if (typeof window !== 'undefined') {
    if (!persistWatchRegistered) {
      persistWatchRegistered = true;
      persistScope.run(() => {
        watch(followAutonext, persistFollow);
      });
    }
    onMounted(() => {
      syncFollowFromStorageOnce();
    });
  }
  return followAutonext;
}
