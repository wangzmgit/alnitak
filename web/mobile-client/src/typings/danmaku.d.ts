interface AddDanmakuType {
  vid: number,
  time: number,
  color: string,
  type: number,
  text: string,
  part: number,
}

interface DanmakuType {
  vid: number;
  time: number,
  color: string,
  type: number,
  text: string,
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