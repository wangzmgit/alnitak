import Hls from 'hls.js';

export interface HlsPlayerState {
  instance: Hls | null;
  videoElement: HTMLVideoElement | null;
  playPromise: Promise<void> | null;
}

export interface HlsPlayerOptions {
  maxBufferLength?: number;
  maxMaxBufferLength?: number;
  onError?: (event: string, data: any) => void;
  onMediaAttached?: () => void;
  onManifestParsed?: () => void;
  onTokenExpired?: () => void;
}

export interface PlaybackState {
  currentTime: number;
  wasPlaying: boolean;
  volume: number;
  muted: boolean;
}

export function getSavedPlaybackState(video: HTMLVideoElement): PlaybackState {
  return {
    currentTime: video.currentTime > 0 ? video.currentTime : 0,
    wasPlaying: !video.paused,
    volume: video.volume,
    muted: video.muted,
  };
}

export function restorePlaybackState(
  video: HTMLVideoElement,
  playbackState: PlaybackState,
  hlsState: HlsPlayerState,
  options?: { autoPlay?: boolean }
): void {
  if (playbackState.currentTime > 0) {
    video.currentTime = playbackState.currentTime;
  }
  video.volume = playbackState.volume;
  video.muted = playbackState.muted;
}

export function createHlsPlayer(
  video: HTMLVideoElement,
  url: string,
  state: HlsPlayerState,
  playbackState: PlaybackState | null,
  options: HlsPlayerOptions = {}
): Hls {
  if (state.instance) {
    state.instance.destroy();
    state.instance = null;
  }

  state.instance = new Hls({
    maxBufferLength: options.maxBufferLength ?? 30,
    maxMaxBufferLength: options.maxMaxBufferLength ?? 60,
  });

  state.instance.loadSource(url);
  state.instance.attachMedia(video);

  state.instance.on(Hls.Events.MEDIA_ATTACHED, () => {
    if (playbackState) {
      restorePlaybackState(video, playbackState, state);
    }
    options.onMediaAttached?.();
  });

  state.instance.on(Hls.Events.MANIFEST_PARSED, () => {
    options.onManifestParsed?.();
  });

  state.instance.on(Hls.Events.ERROR, (event, data) => {
    // 检测 403 错误（JWT token 过期）
    if (data.type === Hls.ErrorTypes.NETWORK_ERROR && data.details === Hls.ErrorDetails.ERROR_NOTE_TIMEOUT) {
      const response = (data as any).response;
      if (response && response.code === 403) {
        console.warn('[HLS] Token 过期，触发续签');
        options.onTokenExpired?.();
        return;
      }
    }
    
    // 检测 403 HTTP 状态码
    if (data.type === Hls.ErrorTypes.NETWORK_ERROR && data.networkDetails) {
      const xhr = data.networkDetails as XMLHttpRequest;
      if (xhr && xhr.status === 403) {
        console.warn('[HLS] HTTP 403 Token 过期，触发续签');
        options.onTokenExpired?.();
        return;
      }
    }

    if (data.fatal) {
      switch (data.type) {
        case Hls.ErrorTypes.NETWORK_ERROR:
          console.error('[HLS] 网络错误，尝试恢复:', data.details);
          state.instance?.startLoad();
          break;
        case Hls.ErrorTypes.MEDIA_ERROR:
          console.error('[HLS] 媒体错误，尝试恢复:', data.details);
          state.instance?.recoverMediaError();
          break;
        default:
          console.error('[HLS] 致命错误，无法恢复:', data.type, data.details);
          state.instance?.destroy();
          state.instance = null;
          break;
      }
    }
    options.onError?.(event, data);
  });

  return state.instance;
}

export function destroyHlsPlayer(state: HlsPlayerState): void {
  if (state.instance) {
    state.instance.destroy();
    state.instance = null;
  }
}

export function setupVolumePersistence(video: HTMLVideoElement): void {
  if (!(video as any).__volumeChangeBound) {
    (video as any).__volumeChangeBound = true;
    video.addEventListener('volumechange', () => {
      localStorage.setItem('wplayer-volume', video.volume.toString());
      localStorage.setItem('wplayer-muted', video.muted ? '1' : '0');
    });
  }
}

export function getSavedVolumeState(): { volume: number; muted: boolean } {
  const savedVolume = localStorage.getItem('wplayer-volume');
  const savedMuted = localStorage.getItem('wplayer-muted');
  return {
    volume: savedVolume !== null ? parseFloat(savedVolume) : 1,
    muted: savedMuted === '1',
  };
}
