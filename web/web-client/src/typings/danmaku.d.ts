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
  createdAt: number; // 添加 createdAt 属性
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