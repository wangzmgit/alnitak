<template>
  <div id="artplayer" ref="playerContainer"></div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref, watch, nextTick } from 'vue';
import Artplayer from 'artplayer';
import Hls from 'hls.js';
import * as dashjs from 'dashjs';
import artplayerPluginDanmuku from 'artplayer-plugin-danmuku';
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
import { fetchSubtitleTracksForArtplayer } from '@/utils/subtitle-tracks';
import { pickSubtitleTrackIndexByPreference, writeStoredSubtitlePreference } from '@/utils/subtitle-preference';

const props = defineProps<{
  videoInfo: VideoType;
  part: number;
  progress: number | null;
}>();

const playerContainer = ref<HTMLElement | null>(null);
let player: any = null;
let dashPlayer: any = null;
let hlsPlayerState: HlsPlayerState = { instance: null, videoElement: null, playPromise: null };
let originalDanmaku: DanmakuType[] = [];

// DASH 统一 MPD 模式状态
let dashUnifiedMode = false;
let dashQualityMap: Map<string, number> = new Map();

/** 外链 Artplayer blob 字幕撤销 */
let embedSubtitleRevoke: (() => void) | null = null;

/** CC/字幕底线栏快捷键：是否与上次保持一致 */
const ARTPLAYER_CC_LS = 'artplayer-cc-visible';

/** 字幕叠加层字号/字重（Artplayer 写入 $subtitle 行内样式） */
const EMBED_SUBTITLE_STYLE: Record<string, string> = {
  color: '#fff',
  'font-size': 'clamp(25px, 1.5vw + 17px, 37px)',
  'font-weight': '600',
  'line-height': '1.45',
  'text-shadow': [
    '-3px 0 0 #000',
    '3px 0 0 #000',
    '0 -3px 0 #000',
    '0 3px 0 #000',
    '-3px -3px 0 #000',
    '3px -3px 0 #000',
    '-3px 3px 0 #000',
    '3px 3px 0 #000',
    '-2px -2px 0 #000',
    '2px -2px 0 #000',
    '-2px 2px 0 #000',
    '2px 2px 0 #000',
    '-1px 0 0 #000',
    '1px 0 0 #000',
    '0 -1px 0 #000',
    '0 1px 0 #000',
    '-1px -1px 0 #000',
    '1px -1px 0 #000',
    '-1px 1px 0 #000',
    '1px 1px 0 #000',
    '-2px 0 0 #000',
    '2px 0 0 #000',
    '0 -2px 0 #000',
    '0 2px 0 #000',
    '0 4px 9px rgba(0,0,0,0.6)',
  ].join(','),
  'font-family':
    "system-ui,-apple-system,'Segoe UI','PingFang SC','PingFang TC','Microsoft YaHei UI','Microsoft YaHei',sans-serif",
  '-webkit-font-smoothing': 'antialiased',
};

/** Same paths as wplayer-next vendor subtitles.svg / subtitles-on.svg */
const WPLAYER_SUB_SVG_OFF =
  '<path d="M21.2 3.01L21 3H3l-.21.01c-.49.05-.95.28-1.28.64-.33.37-.52.85-.51 1.35v14l.01.2c.04.46.25.88.57 1.21.33.32.75.53 1.21.58L3 21h18l.2-.02c.46-.04.88-.25 1.21-.57.32-.33.53-.75.58-1.21l.01-.14V5c0-.5-.19-.98-.52-1.35-.33-.36-.79-.59-1.28-.64zM3 19V5h18v14H3zm5-8H6c-.27 0-.52.1-.71.29a1 1 0 000 1.42c.19.19.44.29.71.29h2c.26 0 .51-.11.7-.29a1 1 0 000-1.42A.97.97 0 008 11zm10 0h-6c-.27 0-.52.1-.71.29a1 1 0 000 1.42c.19.19.44.29.71.29h6c.26 0 .51-.11.7-.29a1 1 0 000-1.42c-.19-.19-.44-.29-.7-.29zm0 4h-2c-.27 0-.52.1-.71.29a1 1 0 000 1.42c.19.19.44.29.71.29h2c.26 0 .51-.11.7-.29a1 1 0 000-1.42c-.19-.19-.44-.29-.7-.29zm-6 0H6c-.27 0-.52.1-.71.29a1 1 0 000 1.42c.19.19.44.29.71.29h6c.26 0 .51-.11.7-.29a1 1 0 000-1.42c-.19-.19-.44-.29-.7-.29z"/>';
const WPLAYER_SUB_SVG_ON =
  '<path d="M21 3H3C2.46 3 1.96 3.21 1.58 3.58C1.21 3.96 1 4.46 1 5V19C1 19.53 1.21 20.03 1.58 20.41C1.96 20.78 2.46 21 3 21H21C21.53 21 22.03 20.78 22.41 20.41C22.78 20.03 23 19.53 23 19V5C23 4.46 22.78 3.96 22.41 3.58C22.03 3.21 21.53 3 21 3ZM6 11H8C8.26 11 8.51 11.10 8.70 11.29C8.89 11.48 9 11.73 9 12C9 12.26 8.89 12.51 8.70 12.70C8.51 12.89 8.26 13 8 13H6C5.73 13 5.48 12.89 5.29 12.70C5.10 12.51 5 12.26 5 12C5 11.73 5.10 11.48 5.29 11.29C5.48 11.10 5.73 11 6 11ZM12 11H18C18.26 11 18.51 11.10 18.70 11.29C18.89 11.48 19 11.73 19 12C19 12.26 18.89 12.51 18.70 12.70C18.51 12.89 18.26 13 18 13H12C11.73 13 11.48 12.89 11.29 12.70C11.10 12.51 11 12.26 11 12C11 11.73 11.10 11.48 11.29 11.29C11.48 11.10 11.73 11 12 11ZM16 15H18C18.26 15 18.51 15.10 18.70 15.29C18.89 15.48 19 15.73 19 16C19 16.26 18.89 16.51 18.70 16.70C18.51 16.89 18.26 17 18 17H16C15.73 17 15.48 16.89 15.29 16.70C15.10 16.51 15 16.26 15 16C15 15.73 15.10 15.48 15.29 15.29C15.48 15.10 15.73 15 16 15ZM6 15H12C12.26 15 12.51 15.10 12.70 15.29C12.89 15.48 13 15.73 13 16C13 16.26 12.89 16.51 12.70 16.70C12.51 16.89 12.26 17 12 17H6C5.73 17 5.48 16.89 5.29 16.70C5.10 16.51 5 16.26 5 16C5 15.73 5.10 15.48 5.29 15.29C5.48 15.10 5.73 15 6 15Z"/>';

function embedCcButtonHtml(active: boolean) {
  const fill = active ? '#2196f3' : 'rgba(255,255,255,0.85)';
  const inner = active ? WPLAYER_SUB_SVG_ON : WPLAYER_SUB_SVG_OFF;
  const svg =
    '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="' +
    fill +
    '" width="23" height="23" aria-hidden="true">' +
    inner +
    '</svg>';
  return (
    '<span style="display:inline-flex;align-items:center;justify-content:center;width:32px;height:100%;box-sizing:border-box;line-height:0">' +
    svg +
    '</span>'
  );
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
    const fpsStr = parts[parts.length - 1]!; // "30"、"60"、"24"、"50" 等任意帧率值
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

function guessType(url: string, qualityItem?: any) {
  if (qualityItem?.type === 'dash') {
    return 'dash';
  }
  if (
    url.includes('/api/v1/video/getVideoFile') ||
    (qualityItem && qualityItem.type === 'hls')
  ) {
    return 'm3u8';
  }
  if (url.endsWith('.m3u8')) return 'm3u8';
  if (url.endsWith('.mp4')) return 'mp4';
  if (url.endsWith('.flv')) return 'flv';
  return 'mp4';
}

const getQualities = (qualityList: string[], resourceId: number | string, serverSupportsDash: boolean, qualityOrderFromServer: string[] = []) => {
  const sorted = [...qualityList].sort((a, b) => {
    const wa = parseInt(a.split('x')[0], 10);
    const wb = parseInt(b.split('x')[0], 10);
    if (wb !== wa) return wb - wa;
    const fpsA = parseInt(a.split('_').pop() || '0', 10);
    const fpsB = parseInt(b.split('_').pop() || '0', 10);
    return fpsB - fpsA;
  });

  // 必须浏览器支持且服务器资源支持才使用 DASH
  const useDash = supportsDashJs() && serverSupportsDash;
  // 读取保存的清晰度偏好
  const savedQuality = localStorage.getItem('artplayer-quality');

  if (useDash && qualityOrderFromServer.length > 0) {
    // 统一 DASH MPD 模式
    dashUnifiedMode = true;
    dashQualityMap = new Map();
    qualityOrderFromServer.forEach((q, index) => {
      dashQualityMap.set(getQualityDisplayName(q), index);
    });

    const unifiedMpdUrl = getVideoFileUrlDashUnified(resourceId);
    return sorted.map((item, idx) => {
      const displayName = getQualityDisplayName(item);
      return {
        default: savedQuality ? displayName === savedQuality : idx === 0,
        html: displayName,
        url: unifiedMpdUrl,
        type: 'dash',
      };
    });
  }

  dashUnifiedMode = false;
  return sorted.map((item, idx) => {
    const displayName = getQualityDisplayName(item);
    return {
      default: savedQuality ? displayName === savedQuality : idx === 0,
      html: displayName,
      url: useDash ? getVideoFileUrlDash(resourceId, item) : getVideoFileUrl(resourceId, item),
      type: useDash ? 'dash' : 'm3u8',
    };
  });
};

// 检测是否为 Safari 或 iOS 设备（它们不完整支持 MSE / dashjs）
const isSafariOrIOS = (): boolean => {
  const ua = navigator.userAgent;
  if (/iPad|iPhone|iPod/.test(ua)) return true;
  if (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1) return true;
  if (/Safari/.test(ua) && !/Chrome|CriOS|FxiOS|Edg/.test(ua)) return true;
  return false;
};

const supportsDashJs = (): boolean => {
  if (isSafariOrIOS()) return false;
  return !!(
    (window as any).MediaSource ||
    (window as any).ManagedMediaSource
  );
};

const loadDanmaku = async () => {
  const vid = props.videoInfo.vid;
  const part = props.part;
  const res = await getDanmakuAPI(vid, part);
  originalDanmaku = res.data.code === 200 && Array.isArray(res.data.data.danmaku)
    ? res.data.data.danmaku.map((d: any) => ({
        ...d,
        mode: d.type,
      }))
    : [];
};

const initPlayer = async () => {
  const container = playerContainer.value;
  if (!container) {
    return;
  }
  if (player) {
    return;
  }

  if (!props.videoInfo?.resources?.length) {
    console.warn('[art.vue] videoInfo.resources is empty or undefined');
    return;
  }

  const resource = props.videoInfo.resources[props.part - 1];
  if (!resource?.id) {
    return;
  }
  const rid = resource.shortId || resource.id;
  const subPrep = await fetchSubtitleTracksForArtplayer(String(rid));
  const res = await getResourceQualityApi(rid);
  let qualities = [];
  if (res.data.code === 200 && res.data.data.quality.length > 0) {
    const serverSupportsDash = res.data.data.supportsDash === true;
    const qualityOrderFromServer = (res.data.data.qualityOrder as string[]) || [];
    qualities = getQualities(res.data.data.quality, rid, serverSupportsDash, serverSupportsDash ? qualityOrderFromServer : []);
  } else {
    qualities = [{ default: true, html: '默认', url: resource.url, type: 'm3u8' }];
  }

  await loadDanmaku();

  const embedSubTracks = subPrep.tracks;
  const subtitleBarToggle = {
    visible: embedSubTracks.length === 0 || localStorage.getItem(ARTPLAYER_CC_LS) !== '0',
  };

  function applyPreferredEmbedSubtitleWhenShowing() {
    if (!embedSubTracks.length || !player?.subtitle?.switch) return;
    const idx = pickSubtitleTrackIndexByPreference(embedSubTracks);
    const st = embedSubTracks[idx];
    if (!st?.url) return;
    player.subtitle.switch(st.url, { name: st.label });
  }

  function syncSubtitleQuickControl() {
    if (!embedSubTracks.length || !player) return;
    const html = embedCcButtonHtml(subtitleBarToggle.visible);
    const tooltip = subtitleBarToggle.visible ? '隐藏字幕' : '显示字幕';
    try {
      if (typeof player.controls?.update === 'function') {
        player.controls.update({
          name: 'subtitle-quick',
          html,
          tooltip,
        });
        return;
      }
    } catch {
      /* 继续 DOM 兜底 */
    }
    try {
      const root = player.controls['subtitle-quick'] as HTMLElement | undefined;
      const innerBtn = root?.querySelector?.('button') ?? root?.firstElementChild;
      const target = (innerBtn ?? root) as HTMLElement | undefined;
      if (target) {
        target.innerHTML = html;
        target.setAttribute('aria-label', tooltip);
        target.setAttribute('aria-pressed', subtitleBarToggle.visible ? 'true' : 'false');
      }
      const tipHost = root?.closest?.('[data-balloon]') ?? root?.parentElement;
      if (tooltip && tipHost) {
        tipHost.setAttribute('data-balloon', tooltip);
        tipHost.setAttribute('aria-label', tooltip);
      }
    } catch {
      /* noop */
    }
  }

  function syncGearSubtitleToggle() {
    if (!player || !embedSubTracks.length) return;
    const idx = pickSubtitleTrackIndexByPreference(embedSubTracks);
    try {
      player.setting.update({
        name: 'subtitle-setting',
        width: 200,
        html: '字幕',
        tooltip: embedSubTracks[idx]?.label ?? '',
        icon:
          '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#2196f3" width="20" height="20">' +
          WPLAYER_SUB_SVG_OFF +
          '</svg>',
        selector: [
          {
            html: '显示',
            tooltip: subtitleBarToggle.visible ? '隐藏' : '显示',
            switch: subtitleBarToggle.visible,
          },
          ...embedSubTracks.map((st, i) => ({
            html: st.label,
            url: st.url,
            default: i === idx,
          })),
        ],
        onSelect(item: { html?: string; url?: string }) {
          if (item.url && player) {
            player.subtitle.switch(item.url, { name: item.html ?? '' });
            const st = embedSubTracks.find((t) => t.url === item.url);
            writeStoredSubtitlePreference({
              label: ((item.html || st?.label) ?? '').trim(),
              lang: (st?.srclang ?? '').trim(),
            });
          }
          return item.html ?? '';
        },
      });
    } catch {
      /* setting.update may not be available */
    }
  }

  const subtitleQuickControl =
    embedSubTracks.length > 0
      ? {
          name: 'subtitle-quick',
          index: 8,
          position: 'right',
          tooltip: subtitleBarToggle.visible ? '隐藏字幕' : '显示字幕',
          html: embedCcButtonHtml(subtitleBarToggle.visible),
          click() {
            if (!player?.subtitle) return;
            const next = !subtitleBarToggle.visible;
            subtitleBarToggle.visible = next;
            player.subtitle.show = next;
            localStorage.setItem(ARTPLAYER_CC_LS, subtitleBarToggle.visible ? '1' : '0');
            if (next) applyPreferredEmbedSubtitleWhenShowing();
            syncSubtitleQuickControl();
            syncGearSubtitleToggle();
          },
        }
      : null;

  let artSubtitleMenu: Record<string, unknown> | undefined;
  if (embedSubTracks.length) {
    const prefIdx = pickSubtitleTrackIndexByPreference(embedSubTracks);
    const active = embedSubTracks[prefIdx]!;
    artSubtitleMenu = {
      name: 'subtitle-setting',
      width: 200,
      html: '字幕',
      tooltip: active.label,
      icon:
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#2196f3" width="20" height="20">' +
        WPLAYER_SUB_SVG_OFF +
        '</svg>',
      selector: [
        {
          html: '显示',
          tooltip: subtitleBarToggle.visible ? '隐藏' : '显示',
          switch: subtitleBarToggle.visible,
          onSwitch(item: { tooltip?: string; switch: boolean }) {
            if (!player) return item.switch;
            item.tooltip = item.switch ? '隐藏' : '显示';
            const show = !item.switch;
            player.subtitle.show = show;
            subtitleBarToggle.visible = show;
            localStorage.setItem(ARTPLAYER_CC_LS, subtitleBarToggle.visible ? '1' : '0');
            if (show) applyPreferredEmbedSubtitleWhenShowing();
            syncSubtitleQuickControl();
            return !item.switch;
          },
        },
        ...embedSubTracks.map((st, idx) => ({
          html: st.label,
          url: st.url,
          default: idx === prefIdx,
        })),
      ],
      onSelect(item: { html?: string; url?: string }) {
        if (item.url && player) {
          player.subtitle.switch(item.url, { name: item.html ?? '' });
          const st = embedSubTracks.find((t) => t.url === item.url);
          writeStoredSubtitlePreference({
            label: ((item.html || st?.label) ?? '').trim(),
            lang: (st?.srclang ?? '').trim(),
          });
        }
        return item.html ?? '';
      },
    };
  }

  const type = guessType(qualities[0].url, qualities[0])!;
  const isDash = type === 'dash';

  // 读取本地循环播放初始状态
  const loopInit = localStorage.getItem('artplayer-loop') === '1';

  try {
  player = new Artplayer({
    container,
    url: qualities[0].url,
    quality: qualities,
    type,
    isLive: false,
    autoplay: localStorage.getItem('artplayer-autoplay') !== '0',
    muted: false,
    volume: 0.8,
    fullscreen: true,
    setting: true,
    playbackRate: true,
    aspectRatio: true,
    autoPlayback: true,
    screenshot: true,
    hotkey: true,
    pip: true,
    controls: subtitleQuickControl ? [subtitleQuickControl] : [],

    subtitleOffset: false,

    ...(embedSubTracks.length
      ? { subtitle: { style: EMBED_SUBTITLE_STYLE } }
      : {}),

    theme: '#2196f3',

    // DASH 模式下禁用 artplayer 的 loop
    loop: isDash ? false : loopInit,
  settings: [
    {
      html: '自动播放',
      icon: '<svg viewBox="0 0 24 24" width="20" height="20" fill="#2196f3"><path d="M8 5v14l11-7z"/></svg>',
      switch: localStorage.getItem('artplayer-autoplay') !== '0',
      onSwitch(item: any): boolean {
        const newValue = !item.switch;
        localStorage.setItem('artplayer-autoplay', newValue ? '1' : '0');
        if (player) {
          player.option.autoplay = newValue;
          if (newValue && player.video?.paused) {
            player.play();
          }
        }
        return newValue;
      },
      name: 'autoplay-setting',
    },
    {
      html: '循环播放',
      icon: '<svg viewBox="0 0 1024 1024" width="20" height="20"><path d="M512 64C264.6 64 64 264.6 64 512h64c0-211.7 172.3-384 384-384s384 172.3 384 384-172.3 384-384 384c-70.7 0-137.2-19.2-194.1-52.6l90.1-90.1c12.5-12.5 12.5-32.8 0-45.3s-32.8-12.5-45.3 0l-144 144c-12.5 12.5-12.5 32.8 0 45.3l144 144c12.5 12.5 32.8 12.5 45.3 0s12.5-32.8 0-45.3l-90.1-90.1C374.8 924.8 443.2 944 512 944c247.4 0 448-200.6 448-448S759.4 64 512 64z" fill="#2196f3"/></svg>',
      switch: loopInit,
      onSwitch(item: any): boolean {
        const newLoop = !item.switch;
        localStorage.setItem('artplayer-loop', newLoop ? '1' : '0');

        if (player) {
          if (isDash) {
            // DASH 模式：切换手动循环标记
            dashLoopEnabled = newLoop;
          } else {
            // HLS 模式：用 artplayer 自身的 loop
            player.option.loop = newLoop;
          }
        }

        return newLoop;
      },
      name: 'loop-setting', // 保留 name：主题/存档用
    },
    ...(artSubtitleMenu ? [artSubtitleMenu] : []),
  ],
    layers: [
      // ...Layer 配置...
    ],
customType: {
      m3u8: function (video: HTMLVideoElement, url: string) {
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
            url,
            hlsPlayerState,
            { ...playbackState, volume: volumeState.volume, muted: volumeState.muted },
            {
              maxBufferLength: 60,
              maxMaxBufferLength: 120,
            }
          );
        } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
          video.src = url;
          if (playbackState.currentTime > 0) {
            video.currentTime = playbackState.currentTime;
          }
        }
      },
      dash: function (video: HTMLVideoElement, url: string, art: any) {
        // 统一 MPD 模式下，如果 dashPlayer 已存在且 URL 相同，跳过重复初始化
        if (dashPlayer) {
          try {
            const prevUrl = dashPlayer.getManifestUrl();
            if (prevUrl === url) {
              return;
            }
          } catch {}
          dashPlayer.reset();
          dashPlayer = null;
        }

        if (dashjs && dashjs.MediaPlayer) {
          dashPlayer = dashjs.MediaPlayer().create();
          dashPlayer.updateSettings({
            streaming: {
              buffer: {
                bufferTimeDefault: 20,
                bufferTimeAtTopQuality: 40,
                bufferTimeAtTopQualityLongForm: 90,
                bufferPruningInterval: 15,
                bufferToKeep: 40,
              },
              abr: {
                autoSwitchBitrate: { video: false, audio: false },
              },
            },
            debug: {
              logLevel: 0,
            },
          });
          dashPlayer.initialize(video, url, false);

          // 保存到 art.dash（官方推荐模式）
          art.dash = dashPlayer;

          // playbackEnded 兜底触发 ended 事件（DASH SegmentBase 可能不触发原生 ended）
          let endedHandled = false;
          video.addEventListener('ended', () => { endedHandled = true; });
          dashPlayer.on('playbackEnded', () => {
            if (endedHandled) { endedHandled = false; return; }
            video.dispatchEvent(new Event('ended'));
          });

          dashPlayer.on('error', (e: any) => {
            console.error('[art.vue] DASH 播放错误:', e);
          });

          // 统一 DASH 模式：初始化完成后设置到默认选中的清晰度
          if (dashUnifiedMode) {
            dashPlayer.on('streamInitialized', () => {
              const defaultItem = art.option.quality?.find((q: any) => q.default);
              const defaultName = defaultItem?.html;
              const dashIndex = dashQualityMap.get(defaultName);
              if (dashIndex !== undefined) {
                dashPlayer.setRepresentationForTypeByIndex('video', dashIndex, true);
              }
            });
          }
        } else if (video.canPlayType('application/dash+xml')) {
          video.src = url;
        }
      },
    },
    plugins: [
      artplayerPluginDanmuku({
        danmuku: originalDanmaku,
        speed: 5,
        margin: [10, '25%'],
        opacity: 1,
        mode: 0,
        modes: [0, 1, 2],
        fontSize: 25,
        antiOverlap: true,
        synchronousPlayback: false,
        heatmap: false,
        width: 512,
        points: [],
        filter: (danmu: any) => (danmu.text || danmu.content || '').length <= 100,
        beforeVisible: () => true,
        visible: true,
        emitter: false,
        maxLength: 200,
        lockTime: 5,
        theme: 'dark',
        beforeEmit() {
          return Promise.resolve(true);
        },
      }),
    ],
  });
  embedSubtitleRevoke = subPrep.revoke;
  } catch (e) {
    subPrep.revoke();
    embedSubtitleRevoke = null;
    console.error('[art.vue] Artplayer init failed:', e);
    throw e;
  }

//弹幕输入框
//  player.on('ready', () => {
//    const danmakuInput = player.container.querySelector('.apd-input');
//    if (danmakuInput) {
//      danmakuInput.setAttribute('id', 'artplayer-danmaku-input');
//      danmakuInput.setAttribute('name', 'artplayerDanmaku');
//    }
//  });

  // DASH 模式下手动处理循环播放
  // artplayer 的 loop 会 seek+play 产生竞态，dashjs 的 SegmentBase 也不支持原生 video.loop
  // 所以直接监听 video 的 ended 事件，用 dashjs seek 后延时 play
  let dashLoopEnabled = isDash && loopInit;
  if (isDash) {
    player.on('video:ended', () => {
      if (!dashLoopEnabled || !dashPlayer) return;
      dashPlayer.seek(0);
      setTimeout(() => {
        player?.video?.play().catch(() => {});
      }, 150);
    });
  }

  // 信息面板打开时点击任意位置关闭（标题/UP 浮层 z-index 9999 会盖住 [x] 按钮）
  player.on('info', (show: boolean) => {
    if (!show) return;
    const handler = () => { if (player) player.info.show = false; document.removeEventListener('click', handler); };
    setTimeout(() => document.addEventListener('click', handler), 100);
  });
  player.on('ready', () => {
    if (embedSubTracks.length && player?.subtitle !== undefined) {
      player.subtitle.show = subtitleBarToggle.visible;
      if (subtitleBarToggle.visible) {
        applyPreferredEmbedSubtitleWhenShowing();
      }
    }
    syncSubtitleQuickControl();
    if (player.option.autoplay) {
      const playPromise = player.play();
      if (playPromise !== undefined) {
        playPromise.catch(() => {});
      }
    }
  });

  player.on('quality', (item: any) => {
    // 保存清晰度选择
    if (item.html) {
      localStorage.setItem('artplayer-quality', item.html);
    }

    if (dashUnifiedMode && dashPlayer) {
      // 统一 MPD 模式：通过 dash.js API 无缝切换，不重新加载 MPD
      const dashIndex = dashQualityMap.get(item.html);
      if (dashIndex !== undefined) {
        dashPlayer.setRepresentationForTypeByIndex('video', dashIndex, true);
      }
    } else {
      // HLS 或旧 DASH 模式：使用 switchQuality 保持进度
      if (player && player.switchQuality) {
        player.switchQuality(item.url);
      }
    }
  });
};

onMounted(() => {
  nextTick(() => {
    if (props.videoInfo?.resources?.length) {
      initPlayer();
    }
  });
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

// 组件卸载时清理资源
onBeforeUnmount(() => {
  if (embedSubtitleRevoke) {
    embedSubtitleRevoke();
    embedSubtitleRevoke = null;
  }
  if (player) {
    player.destroy();
    player = null;
  }
  if (dashPlayer) {
    dashPlayer.reset();
    dashPlayer = null;
  }
  destroyHlsPlayer(hlsPlayerState);
});
</script>

<style scoped>
#artplayer {
  height: 100vh;
  width: 100vw;
  margin: 0;
  padding: 0;
}

#artplayer :deep(.art-subtitle) {
  bottom: 20px !important;
  transition: none !important;
}
</style>
