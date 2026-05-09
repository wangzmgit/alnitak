interface PlayerOptionsType {
  container: HTMLElement | null;
  video: PlayerVideoOptionsType;
  danmaku: PlayerDanmakuType;
  /** 显示齿轮设置（内嵌「字幕」轨切换，wplayer-next） */
  setting?: boolean;
  /** 控制 i18n，如 zh-cn */
  lang?: string;
}

interface PlayerVideoOptionsType {
  quality: PlayerQualityType[];
  defaultQuality?: number;
  pic?:string;
  type?: string;
  customType: any;
  /** wplayer-next WebVTT 外挂字幕，运行时用 player.updateSubtitles 更新 */
  subtitles?: WPlayerSubtitleConfigItem[];
}

interface PlayerQualityType {
  name: string;
  url:string;
}

interface PlayerDanmakuType {
  data?: DanmakuType[];
  /** 弹幕区域上移，避免压住 WebVTT 原生 cue */
  bottom?: string;
}