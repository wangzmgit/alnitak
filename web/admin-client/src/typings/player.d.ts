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