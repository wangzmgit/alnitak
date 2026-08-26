<template>
  <div class="embed-player-wrapper">
    <div id="wplayer" ref="playerContainer"></div>
    <!-- 音轨选择器（仅多音轨时显示） -->
    <div v-if="audioTracks.length > 1" class="embed-audio-track-selector"
         @mouseenter="showAudioMenu = true" @mouseleave="showAudioMenu = false">
      <button class="audio-track-btn" :title="`音轨: ${currentAudioLang || '默认'}`">
        <svg class="audio-track-icon" viewBox="0 0 24 24" width="14" height="14" fill="currentColor">
          <path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3A4.5 4.5 0 0 0 14 8.5v7a4.5 4.5 0 0 0 2.5-3.5zM14 3.23v2.06a7.5 7.5 0 0 1 0 13.42v2.06A9.5 9.5 0 0 0 14 3.23z"/>
        </svg>
        <span class="audio-track-text">{{ currentAudioLang || '音轨' }}</span>
      </button>
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch, nextTick, onBeforeUnmount } from 'vue';
import Hls from 'hls.js';
import * as dashjs from 'dashjs';
import Wplayer from 'wplayer-next';
import { getResourceQualityApi, getVideoFileUrl, getVideoFileUrlDash, getVideoFileUrlDashUnified } from '@/api/video';
import { getDanmakuAPI } from '@/api/danmaku';
import {
  createHlsPlayer,
  destroyHlsPlayer,
  getSavedPlaybackState,
  setupVolumePersistence,
  getSavedVolumeState,
  type HlsPlayerState,
} from '@/utils/hls-player';
import { fetchAndApplySubtitles } from '@/utils/subtitle-tracks';

const props = defineProps<{
  videoInfo: VideoType;
  part: number;
  progress: number | null;
}>();

const playerContainer = ref<HTMLElement | null>(null);
let player: any = null;
let dashPlayer: any = null;
let originalDanmaku: DanmakuType[] = [];
const hlsPlayerState: HlsPlayerState = { instance: null, videoElement: null, playPromise: null };

// DASH 统一 MPD 模式状态
let dashUnifiedMode = false;
let dashQualityMap: Map<string, number> = new Map();

// HLS 清晰度切换时保存播放状态
let lastPlaybackState: { time: number; playing: boolean } = { time: 0, playing: false };

// ===== 多音轨支持 =====
const audioTracks = ref<AudioTrackInfo[]>([]);
const currentAudioLang = ref<string>('');
const showAudioMenu = ref(false);
let dashNativeAudioTracks: any[] = [];

const setDanmaku = (data: DanmakuType[]) => {
  originalDanmaku = Array.isArray(data) ? data : [];
}

const injectDanmaku = () => {
  if (player && player.danmaku) {
    player.danmaku.update(Array.isArray(originalDanmaku) ? originalDanmaku : []);
    player.danmaku.show();
    console.log('[embed-player] danmaku injected:', originalDanmaku.length);
  }
}

const resourceNameMap: Record<string, string> = {
  "640x360_1000k_30": "360p",
  "854x480_1500k_30": "480p",
  "1280x720_3000k_30": "720p",
  "1920x1080_6000k_30": "1080p",
  "1920x1080_8000k_60": "1080p60",
};

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
};

const getQualities = (qualityList: string[], resourceId: number | string, qualityOrderFromServer: string[] = []) => {
  // 主站同款排序
  const sorted = [...qualityList].sort((a, b) => {
    const wa = parseInt(a.split('x')[0], 10);
    const wb = parseInt(b.split('x')[0], 10);
    if (wb !== wa) return wb - wa;
    const fpsA = parseInt(a.split('_').pop() || '0', 10);
    const fpsB = parseInt(b.split('_').pop() || '0', 10);
    return fpsB - fpsA;
  });

  const supportDash = supportsDashJs();

  if (supportDash && qualityOrderFromServer.length > 0) {
    // 统一 DASH MPD 模式
    dashUnifiedMode = true;
    dashQualityMap = new Map();
    qualityOrderFromServer.forEach((q, index) => {
      dashQualityMap.set(getQualityDisplayName(q), index);
    });

    const unifiedMpdUrl = getVideoFileUrlDashUnified(resourceId);
    const mapped = sorted.map((item) => ({
      name: getQualityDisplayName(item),
      url: unifiedMpdUrl,
    }));
    return { qualities: mapped, supportDash: true };
  }

  dashUnifiedMode = false;
  const mapped = sorted.map((item) => ({
    name: getQualityDisplayName(item),
    url: supportDash ? getVideoFileUrlDash(resourceId, item) : getVideoFileUrl(resourceId, item),
  }));
  return { qualities: mapped, supportDash };
};

// 检测是否支持 DASH 播放
const supportsDashJs = (): boolean => {
  const ua = navigator.userAgent
  // Safari / iOS 不使用 DASH
  if (/iPad|iPhone|iPod/.test(ua)) return false
  if (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1) return false
  if (/Safari/.test(ua) && !/Chrome|CriOS|FxiOS|Edg/.test(ua)) return false
  return !!((window as any).MediaSource || (window as any).ManagedMediaSource)
}

const loadDanmaku = async () => {
  const vid = props.videoInfo.vid;
  const part = props.part;
  const res = await getDanmakuAPI(vid, part);
  setDanmaku(res.data.code === 200 && Array.isArray(res.data.data.danmaku) ? res.data.data.danmaku : []);
  injectDanmaku();
}

// 新增：获取 URL 参数
function getQueryParam(name: string): string | null {
  const url = window.location.href;
  name = name.replace(/[\[\]]/g, '\\$&');
  const regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)');
  const results = regex.exec(url);
  if (!results) return null;
  if (!results[2]) return '';
  return decodeURIComponent(results[2].replace(/\+/g, ' '));
}

const autoplayParam = getQueryParam('autoplay');
const mutedParam = getQueryParam('muted');
const shouldAutoplay = autoplayParam === '1' || autoplayParam === 'true';
const shouldMuted = mutedParam === '1' || mutedParam === 'true';

const initPlayer = async () => {
  const container = playerContainer.value;
  if (!container) return;
  if (player) return;

  // 防御性检查：确保 videoInfo 和 resources 存在
  if (!props.videoInfo?.resources?.length) {
    console.warn('[embed-player] videoInfo.resources is empty or undefined');
    return;
  }

  const resource = props.videoInfo.resources[props.part - 1];
  if (!resource?.id) {
    console.warn('[embed-player] resource not found for part:', props.part);
    return;
  }

  // 确保 Hls.js 在全局可用
  if (!(window as any).Hls) {
    (window as any).Hls = Hls;
  }

  const rid = resource.shortId || resource.id;
  const res = await getResourceQualityApi(rid);
  let qualities: any[] = [];
  let supportDash = false;
  if (res.data.code === 200 && res.data.data.quality?.length > 0) {
    const qualityOrderFromServer = (res.data.data.qualityOrder as string[]) || [];
    const serverSupportsDash = res.data.data.supportsDash === true;
    const result = getQualities(res.data.data.quality, rid, serverSupportsDash ? qualityOrderFromServer : []);
    qualities = result.qualities;
    supportDash = result.supportDash;
  } else {
    qualities = [{ name: '默认', url: resource.url }];
    supportDash = false;
  }

  /* === 播放器实例化片段 start === */
  player = new Wplayer({
    container,
    setting: true,
    lang: 'zh-cn',
    video: {
      quality: qualities,
      defaultQuality: 0,
      autoplay: shouldAutoplay,
      controls: ["play", "progress", "volume", "quality", "fullscreen"],
      subtitles: [],
      type: supportDash ? 'customDash' : 'customHls',
      customType: {
        customHls: function (video: HTMLVideoElement) {
          const savedVolumeState = getSavedVolumeState();
          const playbackState = getSavedPlaybackState(video);
          const volumeState = {
            volume: playbackState.currentTime > 0 ? playbackState.volume : savedVolumeState.volume,
            muted: playbackState.currentTime > 0 ? playbackState.muted : savedVolumeState.muted,
          };

          setupVolumePersistence(video);

          createHlsPlayer(
            video,
            video.src,
            hlsPlayerState,
            { ...playbackState, volume: volumeState.volume, muted: volumeState.muted },
            {
              maxBufferLength: 30,
              maxMaxBufferLength: 60,
              onError: (event, data) => {
                console.error('[embed-player] HLS 错误:', event, data);
              },
            }
          );
        },
        // DASH 播放（统一 MPD 模式：所有清晰度在一个 MPD 内，通过 dash.js API 无缝切换）
        customDash: function (video: HTMLVideoElement) {
          // 销毁旧实例
          if (dashPlayer) {
            dashPlayer.reset();
            dashPlayer = null;
          }

          dashPlayer = dashjs.MediaPlayer().create();
          dashPlayer.updateSettings({
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
            debug: { logLevel: 0 },
          });
          dashPlayer.initialize(video, video.src, false);

          // playbackEnded 兜底触发 ended 事件
          let endedHandled = false;
          video.addEventListener('ended', () => { endedHandled = true; });
          dashPlayer.on('playbackEnded', () => {
            if (endedHandled) { endedHandled = false; return; }
            video.dispatchEvent(new Event('ended'));
          });

          dashPlayer.on('error', (e: any) => {
            console.error('[embed-player] DASH 播放错误:', e);
          });

          // 统一 DASH 模式：初始化完成后设置到默认选中的清晰度
          if (dashUnifiedMode) {
            dashPlayer.on('streamInitialized', () => {
              const defaultQualityName = qualities[0]?.name;
              const dashIndex = dashQualityMap.get(defaultQualityName);
              if (dashIndex !== undefined) {
                dashPlayer.setRepresentationForTypeByIndex('video', dashIndex, true);
              }
            });
          }

          // 读取 dash.js 原生音频轨
          dashPlayer.on('streamInitialized', () => {
            try {
              const nativeTracks = dashPlayer.getTracksFor('audio');
              if (nativeTracks && nativeTracks.length > 1) {
                dashNativeAudioTracks = nativeTracks;
                const mapped: AudioTrackInfo[] = nativeTracks.map((t: any) => ({
                  language: t.lang || '',
                  title: t.label || t.lang || (t.roles && t.roles[0]) || '',
                  isDefault: false,
                }));
                const defaultTrack = nativeTracks.find((t: any) => t.defaultSelected);
                if (defaultTrack) {
                  currentAudioLang.value = defaultTrack.lang || '';
                }
                audioTracks.value = mapped;
              }
            } catch (e) {
              console.warn('[embed-player] 读取音轨信息失败:', e);
            }
          });
        },
      },
    },
    danmaku: { show: true, bottom: '52px' },
    preload: "auto",
    volume: shouldMuted ? 0 : 0.8,
    muted: shouldMuted,
  });
  /* === 播放器实例化片段 end === */

  // 统一 DASH 模式：拦截 Wplayer 的清晰度切换，改用 dash.js API 无缝切换
  if (dashUnifiedMode) {
    player.switchQuality = function (index: number | string) {
      const idx = typeof index === 'string' ? parseInt(index) : index;
      if (idx === player.qualityIndex) return;

      const quality = player.options.video.quality[idx];
      if (!quality) return;

      player.qualityIndex = idx;
      player.quality = quality;
      const qualityText = player.template?.qualityButton?.querySelector('.wplayer-quality-text');
      if (qualityText) qualityText.textContent = quality.name;

      const dashIndex = dashQualityMap.get(quality.name);
      if (dashIndex !== undefined && dashPlayer) {
        dashPlayer.setRepresentationForTypeByIndex('video', dashIndex, true);
        player.notice(`切换至 ${quality.name}`, 1000, undefined, 'switch-quality');
      }

      player.events.trigger('quality_start', quality);
    };
  } else {
    // HLS 模式下 Wplayer 仍会创建新 video 元素，需要保存播放状态
    player.on('timeupdate', () => {
      if (player?.video && player.video.currentTime > 0) {
        lastPlaybackState = { time: player.video.currentTime, playing: !player.video.paused };
      }
    });
    player.on('quality_start', () => {
      void fetchAndApplySubtitles(String(rid), player);
    });
  }

  // 强制设置 video 元素属性，确保自动静音播放生效
  setTimeout(() => {
    const videoEl = player?.video;
    if (videoEl) {
      videoEl.muted = shouldMuted;
      videoEl.volume = shouldMuted ? 0 : 0.8;
      if (shouldAutoplay && typeof videoEl.play === 'function') {
        videoEl.play();
      }
    }
  }, 300);

  player.on('loadedmetadata', () => {
    console.log('[embed-player] player loadedmetadata');
    injectDanmaku();
    void fetchAndApplySubtitles(String(rid), player);
  });

  // 加载弹幕
  await loadDanmaku();
};

// ===== 音轨切换 =====
const switchAudioTrack = (lang: string) => {
  showAudioMenu.value = false
  if (lang === currentAudioLang.value || !dashPlayer) return

  if (dashNativeAudioTracks.length > 0) {
    const target = dashNativeAudioTracks.find((t: any) => t.lang === lang)
    if (target) {
      try {
        dashPlayer.setCurrentTrack(target)
        currentAudioLang.value = lang
        if (player) {
          player.notice(`已切换音轨: ${target.label || lang}`, 2000, undefined, 'switch-audio')
        }
      } catch (e) {
        console.warn('[embed-player] 切换音轨失败:', e)
      }
    }
  }
}

onMounted(() => {
  nextTick(() => {
    if (props.videoInfo?.resources?.length) {
      initPlayer();
    }
  });
});

// 组件卸载时清理
onBeforeUnmount(() => {
  if (player) {
    player.destroy();
    player = null;
  }
  destroyHlsPlayer(hlsPlayerState);
  if (dashPlayer) {
    dashPlayer.reset();
    dashPlayer = null;
  }
});

// 监听 videoInfo 变化，当数据加载完成后初始化播放器
watch(
  () => props.videoInfo,
  (newVal) => {
    if (newVal?.resources?.length && !player && playerContainer.value) {
      nextTick(() => {
        initPlayer();
      });
    }
  },
  { immediate: true, deep: true }
);
</script>

<style scoped>
.embed-player-wrapper {
  position: relative;
  width: 100vw;
  height: 100vh;
}

#wplayer {
  height: 100vh;
  width: 100vw;
  margin: 0;
  padding: 0;

  :deep(.wplayer-subtitles-quick.wplayer-subtitles-quick-disabled) {
    opacity: 0.72 !important;
  }
}

// ===== 音轨选择器 =====
.embed-audio-track-selector {
  position: absolute;
  bottom: 52px;
  right: 80px;
  z-index: 25;
  user-select: none;

  .audio-track-btn {
    display: flex;
    align-items: center;
    gap: 3px;
    padding: 3px 8px;
    background: rgba(0, 0, 0, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 4px;
    color: #fff;
    font-size: 11px;
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
      max-width: 50px;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }

  .audio-track-dropdown {
    position: absolute;
    bottom: 100%;
    right: 0;
    margin-bottom: 6px;
    min-width: 130px;
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
</style> 
