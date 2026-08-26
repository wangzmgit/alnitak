declare module 'wplayer-next' {
  /** 最小声明：实际以包内 API 为准，详见 vendor/wplayer-next/docs/guide/subtitles.md */
  export default class WPlayer {
    constructor(options: unknown);
    video: HTMLVideoElement;
    options: { video?: { subtitles?: WPlayerSubtitleConfigItem[]; [k: string]: unknown }; [k: string]: unknown };
    updateSubtitles(config: WPlayerSubtitleConfigItem[] | WPlayerSubtitleConfigItem): void;
    destroy(): void;
    play(): void;
    seek(time: number): void;
    on(event: string, handler: (...args: unknown[]) => void): void;
    [k: string]: unknown;
  }
}
declare module 'element-plus/dist/locale/zh-cn.mjs';