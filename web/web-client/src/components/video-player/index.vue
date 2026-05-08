<template>
  <!-- 播放器容器和弹幕发送区 -->
  <div class="player-container">
    <div class="player" id="dplayer"></div>
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
import { getResourceQualityApi, getVideoFileUrl, getVideoFileUrlDash, getVideoFileUrlDashUnified } from "@/api/video";
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

// ===== 组件属性定义 =====
const props = withDefaults(defineProps<{
  videoInfo: VideoType;
  part: number;
  progress: number | null;
}>(), {
  part: 1,
  progress: null
})

const emit = defineEmits<{
  danmakuSent: []
}>()

// 获取当前分P的资源ShortID
const getCurrentResourceShortId = () => {
  const resource = props.videoInfo?.resources?.[props.part - 1];
  return resource?.shortId;
}

// ===== 播放器与弹幕相关变量 =====
let player: any = null;
let dashPlayer: any = null;
const defaultQuality = ref('');
const hlsPlayerState: HlsPlayerState = { instance: null, videoElement: null, playPromise: null };
const dash = shallowRef<any>(null);
const hasEnded = ref(false);

// ===== DASH 统一 MPD 模式状态 =====
let dashUnifiedMode = false;
let dashQualityMap: Map<string, number> = new Map(); // 清晰度显示名 → dash.js Representation 索引

/** DASH 片尾兜底：与 dash 调度/缓冲的秒级容差，过小易误判未完结 */
const DASH_NEAR_END_SEC = 0.45;
const DASH_VERIFY_END_SEC = 0.35;
const DASH_END_FALLBACK_DELAY_MS = 220;

// ===== HLS 清晰度切换时保存播放状态（HLS 模式下 Wplayer 仍会创建新 video 元素） =====
let lastPlaybackState: { time: number; playing: boolean } = { time: 0, playing: false };
const danmakuSendRef = ref<InstanceType<typeof DanmakuSend> | null>(null);
const auth = useAuthStore();
const isLoggedIn = computed(() => auth.isLoggedIn);
const options: PlayerOptionsType = {
  container: null,
  video: {
    quality: [],
    defaultQuality: 0,
    pic: '',
    type: 'customHls',
    customType: {
      customHls: function (video: HTMLVideoElement) {
        const savedVolumeState = getSavedVolumeState();
        const playbackState = getSavedPlaybackState(video);
        const volumeState = {
          volume: playbackState.currentTime > 0 ? playbackState.volume : savedVolumeState.volume,
          muted: playbackState.currentTime > 0 ? playbackState.muted : savedVolumeState.muted,
        };

        setupVolumePersistence(video);

        if (Hls.isSupported()) {
          createHlsPlayer(
            video,
            video.src,
            hlsPlayerState,
            { ...playbackState, volume: volumeState.volume, muted: volumeState.muted },
            {
              maxBufferLength: 30,
              maxMaxBufferLength: 60,
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
          const d = video.duration;
          if (!Number.isFinite(d) || d <= 0) return;
          if (video.currentTime < d - DASH_NEAR_END_SEC) return;
          dashEndedFallbackTicking = true;
          window.setTimeout(() => {
            dashEndedFallbackTicking = false;
            if (player?.setting?.loop || hasEnded.value) return;
            const dur = video.duration;
            if (!Number.isFinite(dur) || dur <= 0) return;
            if (video.currentTime < dur - DASH_VERIFY_END_SEC) return;
            video.dispatchEvent(new Event('ended'));
          }, DASH_END_FALLBACK_DELAY_MS);
        });

        dash.value.on('error', (e: any) => {
          console.error('[DASH] 播放错误:', e);
        });
      },
    },
  },
  danmaku: {
    data: [],
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

/** 每次 Wplayer 重新创建后：用当前 props 续播进度入队（progress 未变时 watch 不会触发） */
const queueProgressRestoreForNewPlayer = () => {
  const p = props.progress;
  if (p != null && p > 0) {
    pendingSeek = p;
  } else {
    pendingSeek = null;
  }
};

const attachPlayerReadyAndProgressFlush = () => {
  if (!player) return;
  const flushPendingSeek = () => {
    playerReady = true;
    onReadyCallbacks.forEach(cb => cb());
    onReadyCallbacks.length = 0;
    if (pendingSeek != null && pendingSeek > 0) {
      doSeek(pendingSeek);
    }
    pendingSeek = null;
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
      doSeek(val);
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
  onReadyCallbacks.push(cb);
};

// ===== 本地已看完标记工具函数 =====
const getWatchedKey = () => `video-watched-${props.videoInfo.vid}-${props.part}`;
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
    if (player) player.destroy();
    options.container = el;
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
      // 非统一模式（HLS 或降级 DASH）：保存清晰度偏好
      player.on('quality_start', (quality: PlayerQualityType) => {
        localStorage.setItem('default-video-quality', quality.name);
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
    filterDanmaku({ disableLeave, disableType });

    if (player && typeof player.play === 'function') {
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
          vid: props.videoInfo.vid,
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
        addHistoryAPI({ vid: props.videoInfo.vid, part: props.part, time: current, duration, rid: getCurrentResourceShortId() });
      }
      lastSeekTime = current;
    });

    player.on('play', () => {
      const v = player?.video;
      if (v) resetReportingAfterReplay(v);
    });

    queueProgressRestoreForNewPlayer();
    attachPlayerReadyAndProgressFlush();
  }
}

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

    // 当前视频不包含上次保存的清晰度名时，回退到本视频的最高档（不记住选择，避免影响其他视频）
    const qualityNames = qualities.map((q) => getQualityDisplayName(q))
    if (!qualityNames.includes(defaultQuality.value)) {
      const highestName = qualityNames[0] || '720p'
      defaultQuality.value = highestName
      // 不更新 localStorage，避免影响其他视频
    }

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

// ===== 弹幕相关方法 =====
const originalDanmaku = shallowRef<DanmakuType[]>([]);
const setDanmaku = (data: DanmakuType[]) => {
  originalDanmaku.value = data;
  // 更新弹幕数量统计
  danmakuSendRef.value?.updateDanmakuCount(data.length);
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
  localStorage.setItem('danmaku-disable-type', filter.disableType.toString());
  localStorage.setItem('danmaku-disable-leave', filter.disableLeave.toString());

  const data = originalDanmaku.value.filter((item: DanmakuType) => {
    return !isDisableType(item, filter.disableType) && (Math.floor(Math.random() * 10) + 1) > filter.disableLeave;
  }).map((d: DanmakuType) => { return { ...d } });

  player.danmaku.update(data);

  player.on('danmaku_loaded', () => {
    console.log("danmaku_load_end")
  })

  // 更新弹幕数量
  danmakuSendRef.value?.updateDanmakuCount(data.length);
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
const uploadHistory = async () => {
  // 如果视频已播放结束，不再上报进度
  if (hasEnded.value) {
    console.log('视频已播放结束，跳过进度上报');
    return;
  }

  const duration = Math.floor(player.video.duration || 0); // 总时长取整
  const currentTime = Math.floor(player.video.currentTime); // 当前进度取整

  await addHistoryAPI({
    vid: props.videoInfo.vid,
    part: props.part,
    time: currentTime >= duration ? -1 : currentTime, // 播放完了就上报 -1
    duration,
    rid: getCurrentResourceShortId(),
  });
}


// ===== 分集切换监听 =====
watch(() => props.part, (newPart, oldPart) => {
  if (newPart !== oldPart) {
    // 切换前上报当前进度（如果未播放完）
    if (!hasEnded.value && !isWatched()) {
      uploadHistory();
    }
    // 加载新分集
    loadPart(newPart);
  }
});

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

// ===== 页面关闭/离开时上报进度，已看完则不再上报 =====
const reportOnLeave = () => {
  if (player && player.video && typeof player.video.currentTime === 'number' && !isWatched()) {
    const duration = Math.floor(player.video.duration || 0); // 总时长取整
    const currentTime = Math.floor(player.video.currentTime); // 当前进度取整
    addHistoryAPI({ vid: props.videoInfo.vid, part: props.part, time: currentTime >= duration ? -1 : currentTime, duration, rid: getCurrentResourceShortId() });
  }
};
if (typeof window !== 'undefined') {
  window.addEventListener('beforeunload', reportOnLeave);
}
onBeforeUnmount(() => {
  if (timer) clearInterval(timer);
  reportOnLeave();
  if (typeof window !== 'undefined') {
    window.removeEventListener('beforeunload', reportOnLeave);
  }
  if (player) {
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

defineExpose({
  seek,
  setOnReady,
  uploadHistory,
  setDanmaku,
  addDanmaku,
  setOnEnded
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
</style>
