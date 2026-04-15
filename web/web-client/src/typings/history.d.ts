interface AddHistoryType {
  vid: number,
  part: number,
  time: number,
  duration: number,
  resourceShortId?: string
}

interface HistoryVideoType {
  vid: number;
  shortId?: string;
  uid: number;
  title: string;
  cover: string;
  desc: string;
  time: number;
  duration?: number;
  updatedAt: string;
  part?: number;

  /** PGC：列表展示与续播链接按剧集维度 */
  pgcAttached?: boolean;
  pgcTitle?: string;
  episodeTitle?: string;
  episodeNumber?: number;
  epId?: number;

  // 观看日期，后端不返回
  viewingDate?:string;
}	
