interface AddDanmakuType {
  vid: string | number,
  time: number,
  color: string,
  type: number,
  text: string,
  part: number,
  rid?: string,
}

interface DanmakuType {
  vid: string | number;
  time: number,
  color: string,
  type: number,
  text: string,
  part?: number,
  rid?: string,
  createdAt?: number;
}

interface DrawDanmakuType {
  color: string,
  type: number,
  text: string,
}

interface FilterDanmakuType {
  disableType: Array<number>,
  disableLeave: number,
}