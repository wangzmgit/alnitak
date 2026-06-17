interface PlayerOptionsType {
  container: HTMLElement | null;
  video: PlayerVideoOptionsType;
  danmaku: PlayerDanmakuType;
}

interface PlayerVideoOptionsType {
  quality: PlayerQualityType[];
  defaultQuality?: number;
  pic?:string;
  type?: string;
  customType: any;
}

interface PlayerQualityType {
  name: string;
  url:string;
}

interface PlayerDanmakuType {
  data?: DanmakuType[];
}

/** 后端返回的音频轨信息（用于多音轨切换） */
interface AudioTrackInfo {
  language: string;
  title: string;
  isDefault: boolean;
}