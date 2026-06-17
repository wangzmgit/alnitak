<template>
  <div class="player-container">
    <div class="player" id="wplayer"></div>
  </div>
</template>

<script setup lang="ts">
import Hls from "hls.js";
import * as dashjs from "dashjs";
import Wplayer from 'wplayer-next';
import { statusCode } from "@/utils/status-code";
import { ref, shallowRef, onBeforeUnmount } from 'vue';
import { getResourceQualityApi, getVideoFileUrl, getVideoFileUrlDash, getVideoFileUrlDashUnified, getVideoFileAPI } from "@/api/video";
import { useMessage } from "naive-ui";
import { getResourceUrl } from "@/utils/resource";
import { storageData } from "@/utils/storage-data";
import { baseURL } from "@/utils/request";

const message = useMessage();

let player: any = null;
let mpdBlobUrl: string | null = null;
const defaultQuality = ref('');
const hls = shallowRef<Hls | null>(null);
const dash = shallowRef<any>(null);

// DASH 统一 MPD 模式状态
let dashUnifiedMode = false;
let dashQualityMap: Map<string, number> = new Map();

const options: any = {
  container: null,
  video: {
    quality: [],
    defaultQuality: 0,
    type: 'customHls',
    customType: {
      // HLS 播放逻辑
      customHls: (video: HTMLVideoElement) => {
        if (Hls.isSupported()) {
          getVideoFileAPI(video.src).then((res) => {
            if (!res.data) return;
            if (hls.value) hls.value.destroy();
            hls.value = new Hls();
            const indexFile = res.data.split('\n').map((line: string) => {
              return line.includes(".ts") ? getResourceUrl(line) : line;
            });
            const blob = new Blob([indexFile.join('\n')], { type: 'text/plain' });
            const blobUrl = URL.createObjectURL(blob);
            hls.value.loadSource(blobUrl);
            hls.value.attachMedia(video);
          }).catch((err) => {
            console.error('[HLS] 获取索引文件失败:', err);
            message.error('视频资源加载失败，请稍后重试');
          });
        } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
          video.load();
        }
      },

      // DASH 播放逻辑
      customDash: (video: HTMLVideoElement) => {
        const token = storageData.get('token');

        // 预先保存 URL，清掉 video.src 防止浏览器预加载跨域 MPD 触发 CORB
        const mpdUrl = video.src;
        video.src = '';

        // 1. 实例复用检查
        if (dash.value) {
          try {
            if (dash.value.getSource() === mpdUrl) return;
          } catch (e) {
            console.warn('[DASH] 读取当前播放源失败，重建播放器实例:', e);
          }
          dash.value.reset();
          dash.value = null;
        }

        // 2. 创建实例
        const dp = dashjs.MediaPlayer().create();
        
        // 3. 配置参数 (dash.js v4.x)
        // 注意: dashjs@4.7.4 的 protection 不可通过 updateSettings 配置
        // EME 检测警告无害，只是 protection module 检查 DRM 支持的日志
        dp.updateSettings({
          streaming: {
            buffer: {
              stableBufferTime: 12,
              bufferTimeAtTopQuality: 30,
            },
            abr: {
              autoSwitchBitrate: { video: false, audio: false },
            },
          },
          debug: { logLevel: 3 },
        });

        // 4. 【关键修复】无论什么模式，只要有 token 就必须注入鉴权
        // 否则后端返回 JSON 错误导致 dash.js 报 "parsing failed"
        // 注意：仅对后端 API 请求加 Authorization，外部 R2 分片请求不加
        // 否则 R2 presigned URL 触发 CORS preflight 导致请求被拦截
        if (token) {
          let _currentReqUrl = '';
          dp.extend('RequestModifier', function () {
            return {
              modifyRequestURL: function (url: string) {
                _currentReqUrl = url;
                return url;
              },
              modifyRequestHeader: function (xhr: XMLHttpRequest) {
                if (_currentReqUrl.indexOf(baseURL) === 0) {
                  xhr.setRequestHeader('Authorization', token);
                }
                return xhr;
              },
            };
          }, true);
        }

        // 5. 注册事件处理器（在 initialize 之前注册，确保不丢失事件）
        dp.on('streamInitialized', () => {
          const quality = dp.getQualityFor('video');
          console.log('[DASH] 初始清晰度索引:', quality, '总清晰度:', options.video.quality.length);
          if (options.video.quality[quality]) {
            console.log('[DASH] 初始清晰度名称:', options.video.quality[quality].name);
          }
          const maxQualityIndex = options.video.quality.length - 1;
          dp.setQualityFor('video', maxQualityIndex);
          console.log('[DASH] 切换到最高清晰度索引:', maxQualityIndex);
          video.dispatchEvent(new Event('loadedmetadata'));
        });

        dp.on('qualitySwitch', (e: any) => {
          console.log('[DASH] 切换清晰度:', e['0'], '->', e['1']);
        });

        let endedHandled = false;
        video.addEventListener('ended', () => { endedHandled = true; });
        dp.on('playbackEnded', () => {
          if (endedHandled) { endedHandled = false; return; }
          video.dispatchEvent(new Event('ended'));
        });

        dp.on('error', (e: any) => {
          console.error('[DASH] 播放错误:', e);
        });

        dash.value = dp;

        // 6. 【关键修复】Wplayer customType 必须为同步函数（不能 async）
        // 因此用 IIFE 包裹异步的 MPD 拉取逻辑
        (async () => {
          let manifestUrl = mpdUrl;
          try {
            const res = await getVideoFileAPI(mpdUrl);
            console.log('[DASH] MPD 拉取完成, 状态:', res.status, '数据长度:', res.data?.length);
            if (res.data) {
              const blob = new Blob([res.data], { type: 'application/dash+xml' });
              manifestUrl = URL.createObjectURL(blob);
              if (mpdBlobUrl) URL.revokeObjectURL(mpdBlobUrl);
              mpdBlobUrl = manifestUrl;
            } else {
              console.warn('[DASH] MPD 返回数据为空，使用直连 URL');
            }
          } catch (e) {
            console.warn('[DASH] MPD 拉取失败，回退直连:', e);
          }
          dp.initialize(video, manifestUrl, false);
        })();
      },
    },
  },
  danmaku: {}
}

const loadVideo = async (resourceId: number) => {
  const el = document.getElementById('wplayer');
  if (el) {
    // 销毁旧播放器和解码器实例
    if (player) player.destroy();
    if (dash.value) { dash.value.reset(); dash.value = null; }
    if (hls.value) { hls.value.destroy(); hls.value = null; }
    if (mpdBlobUrl) { URL.revokeObjectURL(mpdBlobUrl); mpdBlobUrl = null; }

    const resourceReady = await loadResource(resourceId);
    if (!resourceReady) return;

    options.container = el;
    player = new Wplayer(options);

    player.on('quality_start', (quality: any) => {
      localStorage.setItem('default-video-quality', quality.name);
    })

    // 统一 DASH 模式下的无缝切换逻辑
    if (dashUnifiedMode) {
      player.switchQuality = function (index: number | string) {
        const idx = typeof index === 'string' ? parseInt(index) : index;
        console.log('[切换] 目标索引:', idx, '当前索引:', player.qualityIndex);
        if (idx === player.qualityIndex) return;

        const quality = player.options.video.quality[idx];
        if (!quality) {
          console.log('[切换] 未找到清晰度:', idx);
          return;
        }

        console.log('[切换] 切换到:', quality.name, 'dash索引:', dashQualityMap.get(quality.name));
        
        player.qualityIndex = idx;
        player.quality = quality;
        
        const qualityText = player.template?.qualityButton?.querySelector('.wplayer-quality-text');
        console.log('[切换] qualityButton 元素:', player.template?.qualityButton, 'qualityText:', qualityText);
        if (qualityText) {
          qualityText.textContent = quality.name;
        } else if (player.template?.qualityButton) {
          player.template.qualityButton.textContent = quality.name;
        }

        const dashIndex = dashQualityMap.get(quality.name);
        if (dashIndex !== undefined && dash.value) {
          // dash.js v4 使用 setQualityFor
          dash.value.setQualityFor('video', dashIndex);
          player.notice(`切换至 ${quality.name}`, 1000, undefined, 'switch-quality');
        }

        localStorage.setItem('default-video-quality', quality.name);
        player.events.trigger('quality_start', quality);
      };
    }
  }
}

const resourceNameMap: Record<string, string> = {
  "640x360_500k_30": "360p",
  "854x480_900k_30": "480p",
  "1280x720_2000k_30": "720p",
  "1920x1080_3000k_30": "1080p",
  "1920x1080_6000k_60": "1080p60",
}

const getQualityDisplayName = (qualityStr: string): string => {
  if (resourceNameMap[qualityStr]) return resourceNameMap[qualityStr];
  const parts = qualityStr.split('_');
  const resolution = parts[0];
  const fps = parseInt(parts[parts.length - 1], 10);
  if (resolution.includes('x')) {
    const [width, height] = resolution.split('x').map(Number);
    if (!Number.isFinite(width) || !Number.isFinite(height)) {
      return qualityStr.split('_')[0] || qualityStr;
    }

    const shortSide = Math.min(width, height);
    const fpsSuffix = fps > 30 ? fps.toString() : '';

    if (shortSide <= 360) return `360p${fpsSuffix}`;
    if (shortSide <= 480) return `480p${fpsSuffix}`;
    if (shortSide <= 720) return `720p${fpsSuffix}`;
    if (shortSide <= 1080) return `1080p${fpsSuffix}`;
    if (shortSide <= 1440) return `1440p${fpsSuffix}`;
    if (shortSide <= 2160) return `4K${fpsSuffix}`;
    return `${shortSide}p${fpsSuffix}`;
  }
  return qualityStr.split('_')[0] || qualityStr;
}

const isSafariOrIOS = (): boolean => {
  const ua = navigator.userAgent;
  return /iPad|iPhone|iPod/.test(ua) || 
         (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1) ||
         (/Safari/.test(ua) && !/Chrome|CriOS|FxiOS|Edg/.test(ua));
};

const supportsDashJs = (): boolean => {
  if (isSafariOrIOS()) return false;
  return !!((window as any).MediaSource || (window as any).ManagedMediaSource);
}

const loadResource = async (resourceId: number): Promise<boolean> => {
  options.video.quality = [];
  options.video.defaultQuality = 0;
  dashUnifiedMode = false;
  dashQualityMap = new Map();

  try {
    const res = await getResourceQualityApi(resourceId);
    console.log('[清晰度] API 返回:', res.data.data);
    if (!(res.data.code === statusCode.OK && res.data.data.quality?.length > 0)) {
      message.error(res.data.msg || '未获取到可播放的视频清晰度');
      return false;
    }

    const serverSupportsDash = res.data.data.supportsDash === true;
    const useDash = supportsDashJs() && serverSupportsDash;
    const qualityOrderFromServer = (res.data.data.qualityOrder as string[]) || [];

    // 按分辨率降序排列（1080p -> 720p -> 480p -> 360p）
    const qualities = qualityOrderFromServer.length > 0 
      ? [...qualityOrderFromServer].reverse()
      : [...res.data.data.quality].sort((a, b) => {
          const wa = parseInt(a.split('x')[0], 10);
          const wb = parseInt(b.split('x')[0], 10);
          return wb - wa;
        });
    console.log('[清晰度] qualities 数组 (降序):', qualities);

    if (useDash && qualityOrderFromServer.length > 0) {
      dashUnifiedMode = true;
      dashQualityMap = new Map();
      qualityOrderFromServer.forEach((q, index) => {
        dashQualityMap.set(getQualityDisplayName(q), index);
      });

      const unifiedMpdUrl = getVideoFileUrlDashUnified(resourceId);
      options.video.quality = qualities.map((item, index) => {
        const name = getQualityDisplayName(item);
        console.log(`[清晰度] ${item} -> ${name}`);
        if (name === defaultQuality.value) options.video.defaultQuality = index;
        return { name, url: unifiedMpdUrl };
      });
      console.log('[清晰度] 最终 quality 数组:', options.video.quality);
      options.video.type = 'customDash';
    } else {
      dashUnifiedMode = false;
      options.video.quality = qualities.map((item, index) => {
        const name = getQualityDisplayName(item);
        console.log(`[清晰度] ${item} -> ${name}`);
        if (name === defaultQuality.value) options.video.defaultQuality = index;
        return {
          name,
          url: useDash ? getVideoFileUrlDash(resourceId, item) : getVideoFileUrl(resourceId, item),
        };
      });
      console.log('[清晰度] 最终 quality 数组:', options.video.quality);
      options.video.type = useDash ? 'customDash' : 'customHls';
    }
    return options.video.quality.length > 0;
  } catch (err) {
    console.error('[清晰度] 获取失败:', err);
    message.error('获取视频清晰度失败，请稍后重试');
    return false;
  }
}

defineExpose({ loadVideo });

onBeforeUnmount(() => {
  if (player) player.destroy();
  if (hls.value) hls.value.destroy();
  if (dash.value) dash.value.reset();
  if (mpdBlobUrl) { URL.revokeObjectURL(mpdBlobUrl); mpdBlobUrl = null; }
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