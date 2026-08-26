<template>
    <!-- 播放器容器和弹幕发送区 -->
  <div class="player-container">
    <div class="player" id="dplayer"></div>
    <!-- 音轨选择器（仅多音轨时显示） -->
    <div v-if="audioTracks.length > 1" class="audio-track-selector"
         @mouseenter="showAudioMenu = true" @mouseleave="showAudioMenu = false">
      <button class="audio-track-btn" :title="`音轨: ${currentAudioLang || '默认'}`">
        <svg class="audio-track-icon" viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
          <path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3A4.5 4.5 0 0 0 14 8.5v7a4.5 4.5 0 0 0 2.5-3.5zM14 3.23v2.06a7.5 7.5 0 0 1 0 13.42v2.06A9.5 9.5 0 0 0 14 3.23z"/>
        </svg>
        <span class="audio-track-text">{{ currentAudioLang || '音轨' }}</span>
      </button>
      <transition name="audio-fade">
        <div v-if="showAudioMenu" class="audio-track-dropdown">
          <div class="audio-track-dropdown-title">音轨切换</div>
          <div v-for="track in audioTracks" :key="track.language"
               class="audio-track-option"
               :class="{ selected: track.language === currentAudioLang }"
               @click="switchAudioTrack(track.language)">
            <span class="audio-track-option-label">{{ track.title || track.language }}</span>
            <span v-if="track.isDefault" class="audio-track-option-badge">默认</span>
          </div>
        </div>
      </transition>
    </div>
    <div class="danmaku-send">
      <danmaku-send ref="danmakuSendRef" @send="sendDanmaku" @change-show="changeShow" @opacity-change="opacityChange"
        @set-filter="filterDanmaku" :is-logged-in="isLoggedIn"></danmaku-send>
    </div>
  </div>
</template>

<script setup lang="ts">
// ===== 依赖与类型定义 =====
import Hls from "hls.js";
import * as dashjs from "dashjs";
import Wplayer from 'wplayer-next';
import { ref, shallowRef, onBeforeMount, watch, onMounted, onBeforeUnmount, computed } from 'vue';
import { sendDanmakuAPI } from "@/api/danmaku";
import DanmakuSend from "./components/DanmakuSend.vue";
import { getResourceQualityApi, getVideoFileUrl, getVideoFileUrlDash, getVideoFileUrlDashUnified, postPlayGrantAPI, getPlayUrlsAPI, getAudioTracksAPI } from "@/api/video";
import { addHistoryAPI } from "@/api/history";
import { useAuthStore } from "@/stores/auth-store";
import {
  createHlsPlayer,
  destroyHlsPlayer,
  getSavedPlaybackState,
  restorePlaybackState,
  setupVolumePersistence,
  getSavedVolumeState,
  type HlsPlayerState,
  type PlaybackState,
} from "@/utils/hls-player";
import { fetchAndApplySubtitles } from "@/utils/subtitle-tracks";


// ===== 组件属性定义 =====
const props = withDefaults(defineProps<{
  videoInfo: VideoType;
  part: number;
  progress: number | null;
  episodePickerList?: Array<{ label: string; index: number; vid?: string; part?: number; rid?: string; epId?: number }>;
  episodePickerActiveIndex?: number;
  episodePickerType?: 'none' | 'parts' | 'collection' | 'pgc';
}>(), {
  part: 1,
  progress: null,
  episodePickerList: () => [],
  episodePickerActiveIndex: 0,
  episodePickerType: 'none',
})

const emit = defineEmits<{
  danmakuSent: [];
  episodePick: [item: { vid?: string; part?: number; rid?: string; epId?: number }];
}>()

// 获取当前分P的资源ShortID
const getCurrentResourceShortId = () => {
  const resource = props.videoInfo?.resources?.[props.part - 1];
  return resource?.shortId;
}

/** 指定分 P 的 resourceShortId（用于字幕与历史 rid；无 shortId 时用数字 id 兼容后端 ParseResourceID） */
const getResourceShortIdForPart = (partNum: number): string | undefined => {
  const r = props.videoInfo?.resources?.[partNum - 1];
  if (!r) return undefined;
  if (r.shortId) return String(r.shortId);
  if (r.id != null) return String(r.id);
  return undefined;
}

// ===== 播放器与弹幕相关变量 =====
let player: any = null;
let dashPlayer: any = null;
let loadingPart = false;
const defaultQuality = ref('');
const hlsPlayerState: HlsPlayerState = { instance: null, videoElement: null, playPromise: null };
const dash = shallowRef<any>(null);
const hasEnded = ref(false);

// ===== PlayURL 授权 & 备用 OSS URL =====
const playGrantToken = ref<string>('');
const playGrantExpires = ref<number>(0);
const backupVideoUrl = ref<string>('');
const backupAudioUrl = ref<string>('');
const selectedLineLabel = ref<'primary' | 'backup'>('primary');

// ===== 多音轨支持 =====
const audioTracks = ref<AudioTrackInfo[]>([]);
const currentAudioLang = ref<string>('');
const showAudioMenu = ref(false);
/** dash.js 原生 AudioTrack 对象引用，用于切换 */
let dashNativeAudioTracks: any[] = [];

// ===== DASH 统一 MPD 模式状态 =====
let dashUnifiedMode = false;
let dashQualityMap: Map<string, number> = new Map(); // 清晰度显示名 → dash.js Representation 索引

/** DASH 片尾兜底延迟（ms），给 ended 事件缓冲 */
const DASH_END_FALLBACK_DELAY_MS = 220;
/** DASH MPD 提前刷新时间（ms）：OSS 签名 URL 有效期 24h，提前 30min 刷新 */
const DASH_MPD_REFRESH_AHEAD_MS = 30 * 60 * 1000;
/** DASH MPD OSS 签名 URL 有效期（ms）：24h */
const DASH_OSS_URL_TTL_MS = 24 * 60 * 60 * 1000;
let dashMpdRefreshTimer: ReturnType<typeof setTimeout> | null = null;

const cancelDashMpdRefresh = () => {
  if (dashMpdRefreshTimer) {
    clearTimeout(dashMpdRefreshTimer);
    dashMpdRefreshTimer = null;
  }
};

/** DASH 容灾标记：下载失败后尝试续签 MPD 一次，续签后仍失败才切备用 OSS */
let dashTokenRefreshed = false;

// ===== HLS 清晰度切换时保存播放状态（HLS 模式下 Wplayer 仍会创建新 video 元素） =====
let lastPlaybackState: { time: number; playing: boolean } = { time: 0, playing: false };
const danmakuSendRef = ref<InstanceType<typeof DanmakuSend> | null>(null);

// ===== 全屏选集状态 =====
let pickerBtnEl: HTMLElement | null = null;
let pickerOverlayEl: HTMLElement | null = null;
let pickerCleanup: (() => void) | null = null;
const auth = useAuthStore();
const isLoggedIn = computed(() => auth.isLoggedIn);
const options: PlayerOptionsType = {
  container: null,
  autoplay: localStorage.getItem('wplayer-autoplay') !== '0',
  setting: true,
  lang: 'zh-cn',
  video: {
    quality: [],
    defaultQuality: 0,
    pic: '',
    type: 'customHls',
    // wplayer-next：初始为空，分 P 加载后由 player.updateSubtitles 写入（见 vendor/wplayer-next/src/js/subtitle.js）
    subtitles: [],
    customType: {
      customHls: function (video: HTMLVideoElement) {
        const savedVolumeState = getSavedVolumeState();
        const playbackState = getSavedPlaybackState(video);
        const volumeState = {
          volume: playbackState.currentTime > 0 ? playbackState.volume : savedVolumeState.volume,
          muted: playbackState.currentTime > 0 ? playbackState.muted : savedVolumeState.muted,
        };
        const savedTime = video.currentTime; // 备份切换前保存进度

        setupVolumePersistence(video);

        if (Hls.isSupported()) {
          let hlsBackupRetried = false; // 是否已切备用 OSS（跨 video 元素重置）
          const origUrl = video.src;

          createHlsPlayer(
            video,
            video.src,
            hlsPlayerState,
            { ...playbackState, volume: volumeState.volume, muted: volumeState.muted },
            {
              maxBufferLength: 30,
              maxMaxBufferLength: 60,
                onTokenExpired: () => {
                  // JWT 过期：重新获取 manifest（带时间戳防缓存），恢复播放
                  const rid = getCurrentResourceShortId();
                  if (!rid) return;
                  const curQuality = options.video.quality[options.video.defaultQuality]?.name || '';
                  const newUrl = getVideoFileUrl(rid, curQuality, Date.now());
                  const savedTime = video.currentTime;
                  const wasPlaying = !video.paused;
                  if (hlsPlayerState.instance) {
                    hlsPlayerState.instance.loadSource(newUrl);
                  }
                  video.currentTime = savedTime;
                  if (wasPlaying) video.play().catch(() => {});
                  console.log('[HLS] Token 过期续签，currentTime:', savedTime);
                },
                onError: (_event, data) => {
                  if (!hlsBackupRetried && data.fatal && data.type === Hls.ErrorTypes.NETWORK_ERROR) {
                    hlsBackupRetried = true;
                    selectedLineLabel.value = 'backup';

                    const backupUrl = origUrl + (origUrl.includes('?') ? '&' : '?') + 'backup=true';
                    console.log('[HLS] 主 OSS 网络错误，切 &backup=true:', backupUrl);
                    if (hlsPlayerState.instance) {
                      hlsPlayerState.instance.loadSource(backupUrl);
                    }
                    if (savedTime > 0) {
                      video.currentTime = savedTime;
                    }
                  }
                },
            }
          );
        } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
          // Safari/iOS 原生 HLS 播放
          video.src = video.src;
          if (playbackState.currentTime > 0) {
            video.currentTime = playbackState.currentTime;
          }
          video.volume = volumeState.volume;
          video.muted = volumeState.muted;
        }
      },
      // DASH 播放（统一 MPD 模式：所有清晰度在一个 MPD 内，通过 dash.js API 无缝切换）
      customDash: function (video: HTMLVideoElement) {
        const savedVolume = localStorage.getItem('wplayer-volume');
        const savedMuted = localStorage.getItem('wplayer-muted');
        const prevVolume = savedVolume !== null ? parseFloat(savedVolume) : 1;
        const prevMuted = savedMuted === '1';

        // 销毁旧实例
        if (dash.value) {
          dash.value.reset();
          dash.value = null;
        }

        dash.value = dashjs.MediaPlayer().create();
        dash.value.updateSettings({
          streaming: {
            buffer: {
              bufferTimeDefault: 12,
              bufferTimeAtTopQuality: 30,
              bufferTimeAtTopQualityLongForm: 60,
              bufferPruningInterval: 10,
              bufferToKeep: 20,
            },
            abr: {
              autoSwitchBitrate: { video: false, audio: false },
            },
          },
          debug: { logLevel: import.meta.dev ? 3 : 0 },
        });
        dash.value.initialize(video, video.src, false);

        // 流初始化完成：恢复音量 + 设置初始清晰度
        dash.value.on('streamInitialized', () => {
          video.volume = prevVolume;
          video.muted = prevMuted;
          video.dispatchEvent(new Event('loadedmetadata'));

          // 设置用户偏好的初始清晰度（forceReplace=true：此时尚未播放，立即切到目标清晰度）
          if (dashUnifiedMode) {
            const dashIndex = dashQualityMap.get(defaultQuality.value);
            if (import.meta.dev) {
              console.log('[DASH] 初始清晰度设置:', {
                defaultQuality: defaultQuality.value,
                dashIndex,
                dashQualityMap: Object.fromEntries(dashQualityMap),
                representations: dash.value.getRepresentationsByType?.('video'),
              });
            }
            if (dashIndex !== undefined) {
              dash.value.setRepresentationForTypeByIndex('video', dashIndex, true);
            }
          }

          // 读取 dash.js 原生音频轨（MPD 中多个 audio AdaptationSet 对应多音轨）
          try {
            const nativeTracks = dash.value.getTracksFor('audio');
            if (nativeTracks && nativeTracks.length > 1) {
              dashNativeAudioTracks = nativeTracks;
              // 映射为 AudioTrackInfo 提供给 UI
              const mapped: AudioTrackInfo[] = nativeTracks.map((t: any) => ({
                language: t.lang || '',
                title: t.label || t.lang || (t.roles && t.roles[0]) || '',
                isDefault: false,
              }));
              // 标记默认轨
              const defaultTrack = nativeTracks.find((t: any) => t.defaultSelected);
              if (defaultTrack) {
                currentAudioLang.value = defaultTrack.lang || '';
              }
              audioTracks.value = mapped;
            }
          } catch (e) {
            console.warn('[DASH] 读取音轨信息失败:', e);
          }

          // DASH MPD 定时续签：提前 30min 刷新 MPD，让 dash.js 无感切换到新 OSS URL
          cancelDashMpdRefresh();
          const refreshMs = DASH_OSS_URL_TTL_MS - DASH_MPD_REFRESH_AHEAD_MS;
          dashMpdRefreshTimer = setTimeout(() => {
            const rid = getCurrentResourceShortId();
            if (!rid || !dash.value) return;
            const savedTime = video.currentTime;
            const wasPlaying = !video.paused;
            const curQuality = options.video.quality[options.video.defaultQuality]?.name || '';
            const freshUrl = (dashUnifiedMode
              ? getVideoFileUrlDashUnified(rid, Date.now())
              : getVideoFileUrlDash(rid, curQuality, Date.now()));
            // dash.js CmcdController 要求绝对 URL
            const absUrl = freshUrl.startsWith('http') ? freshUrl : window.location.origin + freshUrl;
            console.log('[DASH] 定时续签 MPD, currentTime:', savedTime);
            dash.value.attachSource(absUrl);
            if (savedTime > 0) {
              dash.value.seek(savedTime);
            }
            if (wasPlaying) video.play().catch(() => {});
            // 续签后重置 error handler 的 dashTokenRefreshed 标记
            dashTokenRefreshed = false;
          }, refreshMs);
        });

        dash.value.on('playbackMetaDataLoaded', () => {
          video.dispatchEvent(new Event('durationchange'));
        });

        // DASH 播放结束处理：
        // 1. SegmentBase 模式下原生 ended 可能不触发，需要 playbackEnded 兜底
        // 2. 循环播放时必须通过 dash.js seek(0) + play() 重置内部状态，
        //    仅设置 video.currentTime = 0 不会让 dash.js 重新调度 segment 请求
        // 3. 拦截原生 ended 事件防止 Wplayer 的 pause() 打断 dash.js 的重播
        // DASH 循环重播辅助方法（防重入：ended + playbackEnded 可能都触发）
        let dashLoopReplaying = false;
        const dashLoopReplay = () => {
          if (dashLoopReplaying) return;
          dashLoopReplaying = true;
          dash.value.seek(0);
          dash.value.play();
          if (player && player.danmaku) player.danmaku.danIndex = 0;
          setTimeout(() => { dashLoopReplaying = false; }, 500);
        };

        // 与 embed-player 一致：用标记区分「原生 ended 已走过」与「仅 dash playbackEnded」，
        // 避免漏派发（原先依赖 video.ended 在部分 MPD/时序下会误判导致 Wplayer 收不到 ended）
        let dashEndedHandled = false;
        video.addEventListener('ended', () => {
          dashEndedHandled = true;
        });

        // 极少数情况下 native ended 与 playbackEnded 都不驱动业务层：用片尾 timeupdate 再补一次
        let dashEndedFallbackTicking = false;
        const resetDashEndedFallback = () => {
          dashEndedFallbackTicking = false;
        };
        video.addEventListener('seeked', () => {
          const d = video.duration;
          if (Number.isFinite(d) && d > 0 && video.currentTime < d - 1) {
            resetDashEndedFallback();
            dashEndedHandled = false;
          }
        });

        // 拦截原生 ended：循环模式下阻止 Wplayer 的 pause()，改用 dash.js API 重播
        video.addEventListener('ended', (e) => {
          if (player && player.setting && player.setting.loop) {
            e.stopImmediatePropagation();
            dashLoopReplay();
          }
        }, true); // capture 阶段先于 Wplayer handler

		// playbackEnded 兜底（SegmentBase 等可能不触发原生 ended）
		dash.value.on('playbackEnded', () => {
			if (player && player.setting && player.setting.loop) {
				dashLoopReplay();
				return;
			}
			if (dashEndedHandled) {
				dashEndedHandled = false;
				return;
			}
			video.dispatchEvent(new Event('ended'));
		});

        video.addEventListener('timeupdate', () => {
          if (player?.setting?.loop || hasEnded.value || dashEndedFallbackTicking) return;
          // 兜底：video.ended 已为 true 但原生 ended 事件未触发（极少数 DASH SegmentBase 场景）
          if (!video.ended) return;
          dashEndedFallbackTicking = true;
          window.setTimeout(() => {
            dashEndedFallbackTicking = false;
            if (player?.setting?.loop || hasEnded.value) return;
            // 双重确认 video 确实已结束
            if (!video.ended) return;
            video.dispatchEvent(new Event('ended'));
          }, DASH_END_FALLBACK_DELAY_MS);
        });

        // DASH 容灾：下载失败（含 JWT 过期 403）时刷新 MPD（新 ts 防缓存）
        dash.value.on('error', (e: any) => {
          console.error('[DASH] 播放错误:', e);
          const isDownloadError = e?.error === 'download';
          if (!isDownloadError) return;

          const rid = getCurrentResourceShortId();
          if (!rid) return;

          const savedDashTime = video.currentTime;
          const wasPlaying = !video.paused;

          if (!dashTokenRefreshed) {
            // 首次下载失败：先尝试用新 ts 刷新 MPD（解决 JWT 过期）
            dashTokenRefreshed = true;
            const curQuality = options.video.quality[options.video.defaultQuality]?.name || '';
            const freshUrl = (dashUnifiedMode
              ? getVideoFileUrlDashUnified(rid, Date.now())
              : getVideoFileUrlDash(rid, curQuality, Date.now()));
            const absUrl = freshUrl.startsWith('http') ? freshUrl : window.location.origin + freshUrl;
            console.log('[DASH] 下载失败，刷新 MPD 续签, currentTime:', savedDashTime);
            dash.value.attachSource(absUrl);
            if (savedDashTime > 0) {
              dash.value.seek(savedDashTime);
            }
            if (wasPlaying) video.play().catch(() => {});
            return;
          }

          // 续签后仍失败：切备用 OSS（新 ts 防缓存）
          selectedLineLabel.value = 'backup';
          const curQ = options.video.quality[options.video.defaultQuality]?.name || '';
          const baseBackupUrl = (dashUnifiedMode
            ? getVideoFileUrlDashUnified(rid, Date.now())
            : getVideoFileUrlDash(rid, curQ, Date.now()));
          const backupRaw = baseBackupUrl + (baseBackupUrl.includes('?') ? '&' : '?') + 'backup=true';
          const backupUrl = backupRaw.startsWith('http') ? backupRaw : window.location.origin + backupRaw;
          console.log('[DASH] 续签后仍失败，切 &backup=true MPD:', backupUrl);
          dash.value.attachSource(backupUrl);
          if (savedDashTime > 0) {
            dash.value.seek(savedDashTime);
          }
        });
      },
    },
  },
  danmaku: {
    data: [],
    bottom: '52px',
  }
}

// ===== 弹幕过滤配置 =====
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

// ===== 进度续播相关 =====
let pendingSeek: number | null = null;
// 标记当前 player 实例是否已经完成 loadedmetadata（可安全 seek）
let playerReady = false;

const doSeek = (time: number) => {
  if (!player) return false;
  try {
    player.seek(time);
    return true;
  } catch (e) {
    console.warn('[video-player] seek 失败:', e);
    return false;
  }
};

/** 已知 duration 时将续播秒数限制在 [0, duration-ε]，避免 seek 越界 */
const clampResumeSeconds = (seconds: number, video: HTMLVideoElement): number => {
  const d = video.duration;
  if (!Number.isFinite(d) || d <= 0) return Math.max(0, seconds);
  const safeEnd = Math.max(0, d - 0.35);
  return Math.min(Math.max(0, seconds), safeEnd);
};

/** 每次 Wplayer 重新创建后：用当前 props 续播进度入队（progress 未变时 watch 不会触发） */
const queueProgressRestoreForNewPlayer = () => {
  const p = props.progress;
  if (p != null && p > 0) {
    pendingSeek = p;
  } else {
    pendingSeek = null;
  }
};

const attachPlayerReadyAndProgressFlush = (partNum: number) => {
  if (!player) return;
  /** loadedmetadata / canplay 都可能触发 flush，避免重复 fetchAndApply（第二次会 revoke 第一次的 blob，多轨字幕全挂） */
  let subtitlesHookFired = false;
  const flushPendingSeek = () => {
    playerReady = true;
    onReadyCallbacks.forEach(cb => cb());
    onReadyCallbacks.length = 0;
    if (pendingSeek != null && pendingSeek > 0 && player.video) {
      doSeek(clampResumeSeconds(pendingSeek, player.video));
    }
    pendingSeek = null;
    if (!subtitlesHookFired) {
      subtitlesHookFired = true;
      void fetchAndApplySubtitles(getResourceShortIdForPart(partNum), player);
    }
  };
  player.on('loadedmetadata', () => {
    flushPendingSeek();
  });
  player.on('canplay', () => {
    if (!playerReady || pendingSeek != null) {
      flushPendingSeek();
    }
  });
};

// ===== 监听 progress 属性变化，自动 seek =====
watch(
  () => props.progress,
  (val) => {
    if (val == null || val <= 0) {
      pendingSeek = null;
      return;
    }
    if (playerReady) {
      if (player?.video) {
        doSeek(clampResumeSeconds(val, player.video));
      } else {
        doSeek(val);
      }
      pendingSeek = null;
    } else {
      pendingSeek = val;
    }
  },
  { immediate: true }
);
// ===== 工具函数 =====
// 将进度转换为整数秒
const toBizSecond = (v: number | null | undefined): number => {
  if (typeof v !== 'number' || !isFinite(v) || v <= 0) return 0;
  return Math.floor(v);
};


// ===== 播放器 ready 回调与定时上报历史 =====
let timer: number | null = null;
let hasReportedWatched = false; // 是否已上报过“已看完”
const onReadyCallbacks: Array<() => void> = [];
const setOnReady = (cb: () => void) => {
  onReadyCallbacks.length = 0;
  onReadyCallbacks.push(cb);
};

// ===== 本地已看完标记工具函数 =====
const getWatchedKey = () => `video-watched-${props.videoInfo.vid}-${props.part}`;
const getWatchedKeyForPart = (partNum: number) =>
  `video-watched-${props.videoInfo.vid}-${partNum}`;
const isWatched = () => localStorage.getItem(getWatchedKey()) === '1';
const setWatched = () => localStorage.setItem(getWatchedKey(), '1');
const clearWatched = () => localStorage.removeItem(getWatchedKey());

/** 片尾容差：离开末尾超过该秒数视为重新播放，恢复进度上报 */
const REPLAY_LEAVE_END_SEC = 0.55;

const videoLeftEndAfterWatched = (v: HTMLVideoElement) => {
  const d = v.duration;
  if (!Number.isFinite(d) || d <= 0) return false;
  return v.currentTime < d - REPLAY_LEAVE_END_SEC;
};

/** 已标记看完的本轮播放结束后，用户不刷新再次播放时需恢复上报 */
const resetReportingAfterReplay = (v: HTMLVideoElement) => {
  if (!videoLeftEndAfterWatched(v)) return;
  if (!hasEnded.value && !isWatched() && !hasReportedWatched) return;
  hasEnded.value = false;
  hasReportedWatched = false;
  clearWatched();
};

// ===== 分集切换与播放器实例化 =====
// 添加播放结束回调
const onEndedCallback = ref<(() => void) | null>(null);

const setOnEnded = (callback: () => void) => {
  onEndedCallback.value = callback;
};

const loadPart = async (part: number) => {
  if (loadingPart) return;
  loadingPart = true;
  cancelDashMpdRefresh();
  try {
  // 重置播放结束标记
  hasEnded.value = false;
  // 新实例未 ready，允许下次 pendingSeek 消费
  playerReady = false;

  const el = document.getElementById('dplayer');
  if (el) {
    await loadResource(part);
    if (!options.video.quality.length) {
      console.error('[video-player] 无可用清晰度，无法播放');
      return;
    }
    /* === 播放器销毁与重建实例化片段 start === */
    if (player) {
      // 清理选集注入的 DOM
      if (pickerBtnEl && pickerBtnEl.parentNode) {
        pickerBtnEl.parentNode.removeChild(pickerBtnEl);
      }
      if (pickerOverlayEl && pickerOverlayEl.parentNode) {
        pickerOverlayEl.parentNode.removeChild(pickerOverlayEl);
      }
      if (pickerCleanup) pickerCleanup();
      pickerBtnEl = null;
      pickerOverlayEl = null;
      pickerCleanup = null;

      // 全屏保护：Wplayer.fullScreen.destroy() 会调用 document.exitFullscreen()，
      // 分P切换时会强制退出浏览器全屏。在销毁前暂替 fullScreen.destroy 为无退出版本，
      // 保留容器(#dplayer)的全屏状态，新实例创建后自动继承。
      const wasFullscreen = !!(document.fullscreenElement || (document as any).webkitFullscreenElement);
      if (wasFullscreen && (player as any).fullScreen) {
        const origFS = (player as any).fullScreen;
        (player as any).fullScreen = {
          destroy() {
            // 只清理事件监听器，不调 exitFullscreen → 容器保持全屏
            document.removeEventListener('fullscreenchange', origFS.handler);
            document.removeEventListener('mozfullscreenchange', origFS.handler);
            document.removeEventListener('webkitfullscreenchange', origFS.handler);
            (origFS as any).handler = null;
            (origFS as any).isFullscreen = true;
          },
        };
      }

      player.destroy();
    }
    // 复用同一 options 时上一分 P 的 subtitles 会污染新实例，导致轨加载失败、CC 一直 disabled
    options.video.subtitles = [];
    options.container = el;
    options.autoplay = localStorage.getItem('wplayer-autoplay') !== '0';
    player = new Wplayer(options);
    /* === 播放器销毁与重建实例化片段 end === */
    hasReportedWatched = false;
    clearWatched();

    // 统一 DASH 模式：拦截 Wplayer 的清晰度切换，改用 dash.js API 无缝切换
    if (dashUnifiedMode) {
      player.switchQuality = function (index: number | string) {
        const idx = typeof index === 'string' ? parseInt(index) : index;
        if (idx === player.qualityIndex) return;

        const quality = player.options.video.quality[idx];
        if (!quality) return;

        // 更新 Wplayer 内部状态（不创建新 video 元素）
        player.qualityIndex = idx;
        player.quality = quality;
        const qualityText = player.template?.qualityButton?.querySelector('.wplayer-quality-text');
        if (qualityText) qualityText.textContent = quality.name;

        // 通过 dash.js 切换 Representation（不强制替换缓冲，等当前缓冲段播完后自然衔接）
        const dashIndex = dashQualityMap.get(quality.name);
        if (dashIndex !== undefined && dash.value) {
          dash.value.setRepresentationForTypeByIndex('video', dashIndex, false);
          player.notice(`切换至 ${quality.name}`, 1000, undefined, 'switch-quality');
        }

        // 保存清晰度偏好 & 触发事件
        localStorage.setItem('default-video-quality', quality.name);
        player.events.trigger('quality_start', quality);
      };
    } else {
      // 非统一模式（HLS 或降级 DASH）：用户主动切换清晰度时保存偏好
      // 不依赖 quality_start（初始化时也会触发，会误存降级后的清晰度）
      {
        const origSwitchQuality = player.switchQuality;
        if (origSwitchQuality) {
          player.switchQuality = function (index: number | string) {
            const idx = typeof index === 'string' ? parseInt(index) : index;
            const quality = player.options.video.quality[idx];
            if (quality) {
              localStorage.setItem('default-video-quality', quality.name);
            }
            return origSwitchQuality.call(this, index);
          };
        }
      }
      player.on('quality_start', (quality: PlayerQualityType) => {
        // 切清晰度会更换 video 元素，需用 WPlayer.updateSubtitles 重新挂载 data-wplayer-subtitle 轨
        void fetchAndApplySubtitles(getResourceShortIdForPart(part), player);
      });

      // HLS 模式下 Wplayer 仍会创建新 video 元素，需要保存播放状态用于恢复
      player.on('timeupdate', () => {
        if (player?.video && player.video.currentTime > 0) {
          lastPlaybackState = { time: player.video.currentTime, playing: !player.video.paused };
        }
      });
      player.on('pause', () => {
        if (player?.video && player.video.currentTime > 0) {
          lastPlaybackState = { time: player.video.currentTime, playing: false };
        }
      });
    }
    // DASH 播放结束后 MediaSource 关闭，原生 video.play() 无法自动重启。
    // 拦截 play 方法，在视频已结束时先通过 dash.js seek(0) 重置 MediaSource 再播放。
    if (dash.value) {
      const origPlay = player.play.bind(player);
      player.play = function (fromNative?: boolean) {
        if (dash.value && player.video?.ended) {
          // seek(0) 重新打开 MediaSource，使后续原生 play() 能正常恢复播放
          dash.value.seek(0);
        }
        origPlay(fromNative);
      };
    }
    filterDanmaku({ disableLeave, disableType });

    if (player && typeof player.play === 'function' && options.autoplay) {
      player.play();
    }

    // 监听播放完成事件，上报已看完并终止定时上报
    player.on('ended', async () => {
      if (hasEnded.value) return;
      hasEnded.value = true; // 标记为已结束

      try {
        // ✅ 业务层统一使用“整数秒”
        const duration = Math.floor(player.video.duration || 0);

        await addHistoryAPI({
          vid: props.videoInfo.shortId || String(props.videoInfo.vid),
          part: props.part,
          time: -1,        // 已看完统一用 -1
          duration,        // 整数秒
          rid: getCurrentResourceShortId(),
        });
      } catch (error) {
        console.error('上报播放完成失败:', error);
      }

      hasReportedWatched = true;
      setWatched();

      if (onEndedCallback.value) {
        onEndedCallback.value();
      }
    });


    // 监听进度条大跨度跳转
    let lastSeekTime = 0;
    player.on('seeked', () => {
      const v = player.video;
      if (v) resetReportingAfterReplay(v);
      const current = player.video.currentTime;
      if (Math.abs(current - lastSeekTime) > 10 && !isWatched() && !hasEnded.value) {
        const current = Math.floor(player.video.currentTime || 0);
        const duration = Math.floor(player.video.duration || 0);
        addHistoryAPI({ vid: props.videoInfo.shortId || String(props.videoInfo.vid), part: props.part, time: current, duration, rid: getCurrentResourceShortId() });
      }
      lastSeekTime = current;
    });

    player.on('play', () => {
      const v = player?.video;
      if (v) resetReportingAfterReplay(v);
    });

    renderEpisodePicker();
    queueProgressRestoreForNewPlayer();
    attachPlayerReadyAndProgressFlush(part);

    // Bug 01 fix: 连播切换分P时重新从 localStorage 读取清晰度偏好并强制应用
    const savedQuality = localStorage.getItem('default-video-quality');
    if (savedQuality && player && player.options && player.options.video) {
      const qIdx = player.options.video.quality.findIndex((q: any) => q.name === savedQuality);
      if (qIdx >= 0) {
        // 同步 ref，防止 DASH 模式下 streamInitialized 异步事件用 stale defaultQuality 覆盖
        defaultQuality.value = savedQuality;
        if (qIdx !== player.qualityIndex) {
          player.switchQuality(qIdx);
        }
      } else {
        // 当前视频不包含该清晰度，用最高档但不覆盖用户偏好
        defaultQuality.value = player.options.video.quality[0]?.name || '720p';
        // 不写 localStorage，保留用户选择的清晰度给下一个视频用
      }
    }

    // Bug 03 fix: 全屏模式下连播重建播放器后保持控制栏可见
    if (document.fullscreenElement || (document as any).webkitFullscreenElement) {
      setTimeout(() => {
        if (player && player.container) {
          const evt = new MouseEvent('mousemove', { bubbles: true });
          player.container.dispatchEvent(evt);
        }
      }, 100);
    }
  }
  } finally {
    loadingPart = false;
  }
}

// ===== 全屏选集：注入按钮 + 下拉弹窗到 Wplayer 控制栏 =====
const renderEpisodePicker = () => {
  // 清理上一轮注入的 DOM
  if (pickerBtnEl && pickerBtnEl.parentNode) {
    pickerBtnEl.parentNode.removeChild(pickerBtnEl);
  }
  if (pickerOverlayEl && pickerOverlayEl.parentNode) {
    pickerOverlayEl.parentNode.removeChild(pickerOverlayEl);
  }
  if (pickerCleanup) pickerCleanup();
  pickerBtnEl = null;
  pickerOverlayEl = null;
  pickerCleanup = null;

  const list = props.episodePickerList;
  if (!list || list.length === 0) return;

  const container = player?.container;
  if (!container) return;
  // 注入到右侧控制栏
  const rightIcons = container.querySelector('.wplayer-icons-right');
  if (!rightIcons) return;

  // ----- 按钮（直接用 wplayer-icon，避免包裹层导致对齐问题） -----
  const btn = document.createElement('button');
  btn.className = 'wplayer-icon';
  btn.style.cssText = 'position:relative;';
  btn.setAttribute('data-balloon', '选集');
  btn.setAttribute('data-balloon-pos', 'up');
  btn.innerHTML = `<span class="wplayer-icon-content">
      <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor">
        <path d="M3 3h7v7H3V3zm0 11h7v7H3v-7zm11-11h7v7h-7V3zm4 15.5L14 16h2.5v-2h-5v5h2v-2.5L17 19l2-1.5z"/>
      </svg>
    </span>`;

  // ----- 弹窗(挂到按钮上，向上弹出) -----
  const overlay = document.createElement('div');
  overlay.className = 'wplayer-episode-picker-overlay';
  overlay.style.cssText = 'display:none;position:absolute;bottom:calc(100% + 4px);right:0;z-index:999;background:rgba(28,28,32,0.96);border-radius:6px;border:1px solid rgba(255,255,255,0.1);max-height:320px;overflow-y:auto;min-width:200px;box-shadow:0 4px 20px rgba(0,0,0,0.5);';

  const type = props.episodePickerType;
  const activeIdx = props.episodePickerActiveIndex || 0;

  list.forEach((item, i) => {
    const row = document.createElement('div');
    const isActive = item.index === activeIdx;
    row.style.cssText = `display:flex;align-items:center;padding:7px 14px;cursor:pointer;font-size:13px;gap:8px;color:${isActive ? '#00a1d6' : '#ddd'};background:${isActive ? 'rgba(0,161,214,0.08)' : 'transparent'};border-bottom:1px solid rgba(255,255,255,0.04);transition:background 0.15s;`;

    // 序号
    const idxSpan = document.createElement('span');
    idxSpan.style.cssText = 'color:#888;font-size:12px;min-width:28px;text-align:center;flex-shrink:0;';
    idxSpan.textContent = String(item.index);
    row.appendChild(idxSpan);

    // 标题
    const labelSpan = document.createElement('span');
    labelSpan.style.cssText = 'overflow:hidden;text-overflow:ellipsis;white-space:nowrap;flex:1;';
    labelSpan.textContent = item.label;
    row.appendChild(labelSpan);

    // "正在播放"标记
    if (isActive) {
      const tag = document.createElement('span');
      tag.style.cssText = 'font-size:11px;color:#00a1d6;flex-shrink:0;';
      tag.textContent = '正在播放';
      row.appendChild(tag);
    }

    row.addEventListener('mouseenter', () => {
      if (!isActive) row.style.background = 'rgba(255,255,255,0.06)';
    });
    row.addEventListener('mouseleave', () => {
      if (!isActive) row.style.background = 'transparent';
    });
    row.addEventListener('click', (e) => {
      e.stopPropagation();
      overlay.style.display = 'none';
      emit('episodePick', { vid: item.vid, part: item.part, rid: item.rid, epId: item.epId });
    });
    overlay.appendChild(row);
  });

  // 悬停展开/收起（带延迟，避免移到弹窗时误关）
  let hoverTimer = null;
  const showOverlay = () => { clearTimeout(hoverTimer); overlay.style.display = 'block'; };
  const scheduleHide = () => { hoverTimer = setTimeout(() => { overlay.style.display = 'none'; }, 300); };
  btn.addEventListener('mouseenter', showOverlay);
  btn.addEventListener('mouseleave', scheduleHide);
  overlay.addEventListener('mouseenter', showOverlay);
  overlay.addEventListener('mouseleave', scheduleHide);

  btn.style.position = 'relative';
  // 插到右侧图标最前面（全屏按钮等之前）
  rightIcons.insertBefore(btn, rightIcons.firstChild);
  btn.appendChild(overlay);

  // 控制栏隐藏时自动关闭弹窗
  const hideObserver = new MutationObserver(() => {
    if (container.classList.contains('wplayer-hide-controller')) {
      overlay.style.display = 'none';
    }
  });
  hideObserver.observe(container, { attributeFilter: ['class'] });

  // 只在全屏时显示按钮
  const updateFs = () => {
    const isFs = player.fullScreen && (player.fullScreen.isFullScreen('browser') || player.fullScreen.isFullScreen('web'));
    btn.style.display = isFs ? '' : 'none';
    if (!isFs) overlay.style.display = 'none';
  };
  player.on('fullscreen', updateFs);
  player.on('fullscreen_cancel', updateFs);
  updateFs();

  pickerBtnEl = btn;
  pickerOverlayEl = overlay;
  pickerCleanup = () => {
    hideObserver.disconnect();
  };
};

// ===== 清晰度映射表与资源加载 =====
const resourceNameMap: Record<string, string> = {
  "640x360_1000k_30": "360p",
  "854x480_1500k_30": "480p",
  "1280x720_3000k_30": "720p",
  "1920x1080_6000k_30": "1080p",
  "1920x1080_8000k_60": "1080p60",
}

/**
 * 根据清晰度字符串动态生成显示名称
 * 格式: "宽x高_码率k_帧率" 例如 "854x480_900k_30" 或 "1920x1080_6000k_60"
 */
const getQualityDisplayName = (qualityStr: string): string => {
  // 如果映射表中存在，直接返回
  if (resourceNameMap[qualityStr]) {
    return resourceNameMap[qualityStr];
  }

  try {
    // 解析格式: "宽x高_码率k_帧率"
    const parts = qualityStr.split('_');
    const resolution = parts[0]; // "854x480"
    const fpsStr = parts[parts.length - 1]; // "30"、"60"、"24"、"50" 等任意帧率值
    const fps = parseInt(fpsStr, 10); // 转换为数字

    if (resolution.includes('x')) {
      const [width, height] = resolution.split('x').map(Number);
      // YouTube风格：取短边作为分辨率标签，竖屏视频不会显示为超高分辨率
      const shortSide = Math.min(width, height);
      const fpsSuffix = fps > 30 ? fps.toString() : '';

      if (shortSide <= 360) {
        return fpsSuffix ? `360p${fpsSuffix}` : '360p';
      } else if (shortSide <= 480) {
        return fpsSuffix ? `480p${fpsSuffix}` : '480p';
      } else if (shortSide <= 720) {
        return fpsSuffix ? `720p${fpsSuffix}` : '720p';
      } else if (shortSide <= 1080) {
        return fpsSuffix ? `1080p${fpsSuffix}` : '1080p';
      } else if (shortSide <= 1440) {
        return fpsSuffix ? `1440p${fpsSuffix}` : '1440p';
      } else if (shortSide <= 2160) {
        return fpsSuffix ? `4K${fpsSuffix}` : '4K';
      } else {
        return fpsSuffix ? `${shortSide}p${fpsSuffix}` : `${shortSide}p`;
      }
    }
  } catch (error) {
    console.warn('Failed to parse quality string:', qualityStr, error);
  }

  return qualityStr.split('_')[0] || qualityStr;
}



// 视频播放信息缓存
const videoPlayInfoCache = new Map<string, any>();

const loadResource = async (part: number) => {
  // 防御性检查
  if (!props.videoInfo?.resources?.length) {
    console.warn('[video-player] videoInfo.resources is empty or undefined');
    return;
  }

  const resource = props.videoInfo.resources[part - 1]
  if (!resource?.id) {
    console.warn('[video-player] resource not found for part:', part);
    return;
  }

  const rid = resource.shortId || resource.id;  // 优先使用 shortId
  const requestTs = Date.now();

  const res = await getResourceQualityApi(rid)
  if (res.data.code === statusCode.OK && res.data.data.quality?.length > 0) {
    // 复制并根据分辨率宽度 & 帧率从高到低排序
    const qualities = [...res.data.data.quality] as string[]
    qualities.sort((a, b) => {
      const wa = parseInt(a.split('x')[0], 10)
      const wb = parseInt(b.split('x')[0], 10)
      if (wb !== wa) return wb - wa
      const fpsA = parseInt(a.split('_').pop() || '0', 10)
      const fpsB = parseInt(b.split('_').pop() || '0', 10)
      return fpsB - fpsA
    })

    // 必须浏览器支持且服务器资源支持才使用 DASH
    const serverSupportsDash = res.data.data.supportsDash === true
    const useDash = supportsDashJs() && serverSupportsDash
    const qualityOrderFromServer = (res.data.data.qualityOrder as string[]) || []

    // 当前视频不包含用户偏好清晰度时，用最高档但不覆盖 defaultQuality.value
    const qualityNames = qualities.map((q) => getQualityDisplayName(q))
    // 不更新 defaultQuality.value，保持用户偏好给下一个视频用
    // options.video.defaultQuality 自动回退到最高档

    if (useDash && qualityOrderFromServer.length > 0) {
      // ===== 统一 DASH MPD 模式：所有清晰度在一个 MPD 内 =====
      dashUnifiedMode = true
      dashQualityMap = new Map()
      qualityOrderFromServer.forEach((q, index) => {
        dashQualityMap.set(getQualityDisplayName(q), index)
      })

      const unifiedMpdUrl = getVideoFileUrlDashUnified(rid, requestTs)
      options.video.quality = qualities.map((item, index) => {
        const name = getQualityDisplayName(item)
        if (name === defaultQuality.value) {
          options.video.defaultQuality = index
        }
        return { name, url: unifiedMpdUrl }
      })
      options.video.type = 'customDash'
    } else if (useDash) {
      // 降级：逐清晰度 DASH（后端未返回 qualityOrder 时）
      dashUnifiedMode = false
      options.video.quality = qualities.map((item, index) => {
        const name = getQualityDisplayName(item)
        if (name === defaultQuality.value) options.video.defaultQuality = index
        return { name, url: getVideoFileUrlDash(rid, item, requestTs) }
      })
      options.video.type = 'customDash'
    } else {
      // HLS 模式（Safari/iOS）
      dashUnifiedMode = false
      options.video.quality = qualities.map((item, index) => {
        const name = getQualityDisplayName(item)
        if (name === defaultQuality.value) options.video.defaultQuality = index
        return { name, url: getVideoFileUrl(rid, item, requestTs) }
      })
      options.video.type = 'customHls'
    }

    // ===== PlayURL 授权 + 备用 OSS URL（B站风格多源容灾） =====
    const resourceShortId = getCurrentResourceShortId()
    if (resourceShortId) {
      try {
        const grantRes = await postPlayGrantAPI(resourceShortId)
        if (grantRes.data.code === statusCode.OK) {
          const token = grantRes.data.data.token
          playGrantToken.value = token
          playGrantExpires.value = grantRes.data.data.expires

          const playRes = await getPlayUrlsAPI(resourceShortId, token)
          if (playRes.data.code === statusCode.OK) {
            backupVideoUrl.value = playRes.data.data.backupVideo || ''
            backupAudioUrl.value = playRes.data.data.backupAudio || ''
            // 读取后端返回的音轨列表（仅非 DASH 模式使用，DASH 模式从 dash.js 直接读取）
            if (playRes.data.data.audioTracks && playRes.data.data.audioTracks.length > 0) {
              // 非 DASH 模式：从 API 获取音轨列表
              if (!dashUnifiedMode && audioTracks.value.length === 0) {
                audioTracks.value = playRes.data.data.audioTracks
                const def = playRes.data.data.audioTracks.find((t: AudioTrackInfo) => t.isDefault)
                if (def) currentAudioLang.value = def.language
              }
              if (import.meta.dev) {
                console.log('[video-player] 可用音轨:', playRes.data.data.audioTracks)
              }
            }
            if (import.meta.dev) {
              console.log('[video-player] 备用 OSS URL:', {
                backupVideo: backupVideoUrl.value,
                backupAudio: backupAudioUrl.value,
              })
            }
          }
        }
      } catch (e) {
        console.warn('[video-player] PlayURL grant 获取失败（不影响主播放）:', e)
      }
    }
  }
}

// 检测是否为 Safari 或 iOS 设备（它们不完整支持 MSE / dashjs）
const isSafariOrIOS = (): boolean => {
  const ua = navigator.userAgent;
  if (/iPad|iPhone|iPod/.test(ua)) return true;
  if (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1) return true;
  if (/Safari/.test(ua) && !/Chrome|CriOS|FxiOS|Edg/.test(ua)) return true;
  return false;
};

// 检测是否支持 dash.js
const supportsDashJs = (): boolean => {
  if (isSafariOrIOS()) return false;
  return !!(
    (window as any).MediaSource ||
    (window as any).ManagedMediaSource
  );
}

// ===== 音轨切换 =====
const switchAudioTrack = (lang: string) => {
  showAudioMenu.value = false
  if (lang === currentAudioLang.value) return

  // DASH 模式：通过 dash.js 原生 API 无缝切换
  if (dash.value && dashNativeAudioTracks.length > 0) {
    const target = dashNativeAudioTracks.find((t: any) => t.lang === lang)
    if (target) {
      try {
        dash.value.setCurrentTrack(target)
        currentAudioLang.value = lang
        if (player) {
          player.notice(`已切换音轨: ${target.label || lang}`, 2000, undefined, 'switch-audio')
        }
        if (import.meta.dev) {
          console.log('[DASH] 切换音轨:', lang, target)
        }
      } catch (e) {
        console.warn('[DASH] 切换音轨失败:', e)
      }
      return
    }
  }

  // HLS / 降级模式：暂不支持无刷新切换（dash.js 未就绪时）
  console.warn('[video-player] 非 DASH 模式暂不支持无刷新切换音轨')
}

// ===== 弹幕相关方法 =====
const originalDanmaku = shallowRef<DanmakuType[]>([]);

/** 将 originalDanmaku 按当前屏蔽规则过滤后同步到 WPlayer 弹幕渲染引擎 */
const syncDanmakuToPlayer = () => {
  if (!player || !player.danmaku) return;
  const filtered = originalDanmaku.value.filter((item: DanmakuType) => {
    return !isDisableType(item, disableType) && (Math.floor(Math.random() * 10) + 1) > disableLeave;
  }).map((d: DanmakuType) => ({ ...d }));
  player.danmaku.update(filtered);
};

const setDanmaku = (data: DanmakuType[]) => {
  originalDanmaku.value = data;
  // 更新弹幕数量统计（总数，不含屏蔽）
  danmakuSendRef.value?.updateDanmakuCount(data.length);
  // 按当前屏蔽规则过滤后同步到播放器渲染引擎
  try {
    syncDanmakuToPlayer();
  } catch (e) {
    console.error('[video-player] setDanmaku sync error:', e);
  }
}
// 本地刚发出、正在等 ws 回广播的弹幕键（避免自己的弹幕被当成他人弹幕二次渲染）
const recentlySent = new Map<string, number>();
// 把 time 按 0.1s 取整，规避后端 float32 往返精度差异，保证发送端构造的 key 能匹配 ws 回广播的 key
const makeDanmakuKey = (d: { time?: number; text?: string; color?: string; type?: number | string }) =>
  `${Math.round((d.time ?? 0) * 10)}|${d.text ?? ''}|${d.color ?? ''}|${d.type ?? ''}`;

// 追加单条弹幕到播放器（不触发 wplayer 的 reload/seek，避免卡顿 + 丢掉在飞的弹幕）
const addDanmaku = (danmaku: DanmakuType) => {
  // 1. 同步本地数据源与计数
  originalDanmaku.value = [...originalDanmaku.value, danmaku];
  danmakuSendRef.value?.updateDanmakuCount(originalDanmaku.value.length);

  if (!player || !player.danmaku) return;

  // 2. 同步 wplayer 的 options.data，确保后续 reload/resize 不会漏掉
  const dataRef = player.danmaku.options && player.danmaku.options.data;
  if (Array.isArray(dataRef)) dataRef.push(danmaku);

  // 3. 若是自己刚发出的那条 ws 回广播：player.danmaku.send 已经插入 dan 并绘制过了，跳过渲染
  const key = makeDanmakuKey(danmaku);
  if (recentlySent.has(key)) {
    recentlySent.delete(key);
    return;
  }

  // 4. 他人的弹幕：只绘制接近当前时刻的，过时的直接丢弃（防止跑完又冒一条）
  const dan = player.danmaku.dan;
  if (!Array.isArray(dan)) return;
  const nowT = player.video?.currentTime ?? 0;
  const dTime = danmaku.time ?? 0;
  // 已错过 > 0.5s 的历史弹幕：只留在 options.data 供后续 reload 用，不再补绘
  if (dTime < nowT - 0.5) return;

  const idx = player.danmaku.danIndex ?? 0;
  let i = idx;
  while (i < dan.length && (dan[i]?.time ?? 0) <= dTime) i++;
  dan.splice(i, 0, danmaku);
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

  player.danmaku.send(danmakuForm, async (danmaku: AddDanmakuType) => {
    danmaku.vid = props.videoInfo.vid;
    danmaku.part = props.part;
    // 附带 rid 用于精准绑定
    const currentRid = props.videoInfo.resources?.[props.part - 1]?.shortId;
    if (currentRid) {
      danmaku.rid = currentRid;
    }
    // 后端要求 vid 为字符串
    const danmakuData = {
      ...danmaku,
      vid: String(danmaku.vid)
    };

    // 记录本地已绘制的 key，让 ws 回广播到这条时跳过重复渲染（30s 后自动过期清理）
    const echoKey = makeDanmakuKey(danmaku);
    recentlySent.set(echoKey, Date.now());
    setTimeout(() => recentlySent.delete(echoKey), 30000);

    const res = await sendDanmakuAPI(danmakuData);
    if (res.data.code !== statusCode.OK) {
      ElMessage.error(res.data.msg);
      // 发送失败：ws 不会回广播，主动清理占位
      recentlySent.delete(echoKey);
    }
  });
}

//过滤弹幕
const filterDanmaku = (filter: FilterDanmakuType) => {
  disableType = filter.disableType;
  disableLeave = filter.disableLeave;
  localStorage.setItem('danmaku-disable-type', filter.disableType.toString());
  localStorage.setItem('danmaku-disable-leave', filter.disableLeave.toString());

  syncDanmakuToPlayer();

  const onLoaded = () => {
    console.log("danmaku_load_end");
  };
  player.on('danmaku_loaded', onLoaded);
}

//是否为屏蔽类型
const isDisableType = (item: DanmakuType, disableType: Array<number>) => {
  if (disableType.includes(item.type))
    return true;
  if (disableType.includes(3) && (item.color !== '#fff' && item.color !== '#ffffff'))
    return true;

  return false;
}

// ===== 历史记录上报 =====
/** 串行化分 P 切换：避免连续切换时在上一个 loadPart 未完成时用错画面做快照 */
let partSwitchTail: Promise<void> = Promise.resolve();

const flushHistoryBeforePartChange = async (previousPart: number) => {
  if (!props.videoInfo?.resources?.length) return;
  if (!player?.video || typeof player.video.currentTime !== 'number') return;
  if (localStorage.getItem(getWatchedKeyForPart(previousPart)) === '1') return;

  const v = player.video;
  const snapshotTime = Math.floor(v.currentTime);
  const snapshotDuration = Math.floor(v.duration || 0);
  const rid = props.videoInfo.resources[previousPart - 1]?.shortId;

  try {
    await addHistoryAPI({
      vid: props.videoInfo.shortId || String(props.videoInfo.vid),
      part: previousPart,
      time: snapshotDuration > 0 && snapshotTime >= snapshotDuration ? -1 : snapshotTime,
      duration: snapshotDuration,
      ...(rid ? { rid } : {}),
    });
  } catch (e) {
    console.error('[video-player] 分P切换前进度上报失败:', e);
  }
};

const uploadHistory = async () => {
  // 如果视频已播放结束，不再上报进度
  if (hasEnded.value) {
    console.log('视频已播放结束，跳过进度上报');
    return;
  }

  const duration = Math.floor(player.video.duration || 0); // 总时长取整
  const currentTime = Math.floor(player.video.currentTime); // 当前进度取整

  await addHistoryAPI({
    vid: props.videoInfo.shortId || String(props.videoInfo.vid),
    part: props.part,
    time: currentTime >= duration ? -1 : currentTime, // 播放完了就上报 -1
    duration,
    rid: getCurrentResourceShortId(),
  });
}


// ===== 分集切换监听（快照上一 P 进度后再 loadPart，链式串行防竞态） =====
watch(
  () => props.part,
  (newPart, oldPart) => {
    if (newPart === oldPart) return;
    const prev = oldPart;
    const next = newPart;
    partSwitchTail = partSwitchTail
      .catch(() => {})
      .then(async () => {
        if (prev !== undefined) {
          await flushHistoryBeforePartChange(prev);
        }
        await loadPart(next);
      });
  }
);

// ===== 选集列表变化时重新渲染 =====
watch(
  () => [props.episodePickerList, props.episodePickerActiveIndex, props.episodePickerType],
  () => {
    if (player) renderEpisodePicker();
  },
  { deep: true }
);

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

  // loadedmetadata / canplay 与续播 flush 已在 loadPart 内按实例绑定

  // 定时上报历史进度，若已看完则停止上报
  timer = window.setInterval(() => {
    if (!hasReportedWatched && !isWatched()) {
      // 检查视频是否正在播放，如果暂停或未播放则跳过上报
      if (player && player.video) {
        const isPlaying = !player.video.paused && !player.video.ended && player.video.currentTime > 0;
        if (isPlaying) {
          uploadHistory(); // 只有在播放时才上报
        } else {
          console.log('视频暂停或未播放，跳过定时上报');
        }
      } else {
        console.log('播放器未初始化，跳过定时上报');
      }
    }
  }, 10000)
})

// 分P/视频切换：原地重载播放器，保持容器不卸载（全屏状态由浏览器保留）
watch([() => props.part, () => props.videoInfo?.vid], ([newPart, newVid], [oldPart, oldVid]) => {
  if ((newPart !== oldPart || newVid !== oldVid) && newPart > 0 && !loadingPart) {
    void loadPart(newPart);
  }
});

// ===== 页面关闭/离开时上报进度，已看完则不再上报 =====
const reportOnLeave = () => {
  if (player && player.video && typeof player.video.currentTime === 'number' && !isWatched()) {
    const duration = Math.floor(player.video.duration || 0); // 总时长取整
    const currentTime = Math.floor(player.video.currentTime); // 当前进度取整
    addHistoryAPI({ vid: props.videoInfo.shortId || String(props.videoInfo.vid), part: props.part, time: currentTime >= duration ? -1 : currentTime, duration, rid: getCurrentResourceShortId() });
  }
};
if (typeof window !== 'undefined') {
  window.addEventListener('beforeunload', reportOnLeave);
}
onBeforeUnmount(() => {
  if (timer) clearInterval(timer);
  cancelDashMpdRefresh();
  reportOnLeave();
  if (typeof window !== 'undefined') {
    window.removeEventListener('beforeunload', reportOnLeave);
  }
  if (player) {
    // 清理选集注入
    if (pickerBtnEl && pickerBtnEl.parentNode) pickerBtnEl.parentNode.removeChild(pickerBtnEl);
    if (pickerOverlayEl && pickerOverlayEl.parentNode) pickerOverlayEl.parentNode.removeChild(pickerOverlayEl);
    if (pickerCleanup) pickerCleanup();
    pickerBtnEl = null;
    pickerOverlayEl = null;
    pickerCleanup = null;
    player.destroy();
    player = null;
  }
  destroyHlsPlayer(hlsPlayerState);
  if (dash.value) {
    dash.value.reset();
    dash.value = null;
  }
});

// ===== 对外暴露方法 =====
const seek = (time: number) => {
  if (player) {
    player.seek(time);
  }
};

const getCurrentTime = () => player?.video?.currentTime ?? 0;
const getDuration = () => player?.video?.duration ?? 0;
/** 通过闭包暴露循环状态（player 为 let 变量，defineExpose 直接暴露始终为 null） */
const isLoopEnabled = () => player?.setting?.loop ?? false;

defineExpose({
  seek,
  getCurrentTime,
  getDuration,
  setOnReady,
  uploadHistory,
  setDanmaku,
  addDanmaku,
  setOnEnded,
  isLoopEnabled,
  player,
  // 备用 OSS URL（供父组件/embed-player 使用）
  backupVideoUrl,
  backupAudioUrl,
  playGrantToken,
  // 当前选择的播放线路（'primary' | 'backup'）
  selectedLineLabel,
  // 多音轨
  audioTracks,
  currentAudioLang,
  switchAudioTrack,
})
</script>

<style lang="scss" scoped>
// ===== 播放器与弹幕样式 =====
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

    &.wplayer-fulled {
      position: fixed;
      top: 0;
      left: 0;
      width: 100vw;
      height: 100vh;
      z-index: 9999;
    }

    /* CC：无轨时为 disabled（默认很淡）；有轨后 wplayer 会去掉 disabled 并可点 */
    :deep(.wplayer-subtitles-quick.wplayer-subtitles-quick-disabled) {
      opacity: 0.72 !important;
    }
  }


  .danmaku-send {
    position: absolute;
    width: 100%;
    bottom: -40px;

    .player-container.wplayer-fulled & {
      display: none;
    }
  }
}

// ===== 音轨选择器 =====
.audio-track-selector {
  position: absolute;
  bottom: 48px;
  right: 80px;
  z-index: 25;
  user-select: none;

  .audio-track-btn {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 4px 10px;
    background: rgba(0, 0, 0, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 4px;
    color: #fff;
    font-size: 12px;
    cursor: pointer;
    transition: background 0.2s;
    white-space: nowrap;

    &:hover {
      background: rgba(0, 0, 0, 0.8);
    }

    .audio-track-icon {
      flex-shrink: 0;
    }

    .audio-track-text {
      max-width: 60px;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }

  .audio-track-dropdown {
    position: absolute;
    bottom: 100%;
    right: 0;
    margin-bottom: 6px;
    min-width: 140px;
    background: rgba(30, 30, 30, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: 6px;
    padding: 4px 0;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4);

    .audio-track-dropdown-title {
      padding: 6px 14px;
      font-size: 11px;
      color: rgba(255, 255, 255, 0.5);
      border-bottom: 1px solid rgba(255, 255, 255, 0.08);
      margin-bottom: 2px;
    }

    .audio-track-option {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 8px 14px;
      font-size: 13px;
      color: #ddd;
      cursor: pointer;
      transition: background 0.15s;

      &:hover {
        background: rgba(255, 255, 255, 0.1);
        color: #fff;
      }

      &.selected {
        color: var(--wplayer-theme, #00a1d6);
        font-weight: 500;
      }

      .audio-track-option-label {
        flex: 1;
      }

      .audio-track-option-badge {
        font-size: 10px;
        opacity: 0.5;
        margin-left: 8px;
        padding: 1px 5px;
        border: 1px solid currentColor;
        border-radius: 3px;
      }
    }
  }
}

.audio-fade-enter-active,
.audio-fade-leave-active {
  transition: opacity 0.15s ease;
}
.audio-fade-enter-from,
.audio-fade-leave-to {
  opacity: 0;
}
</style>
