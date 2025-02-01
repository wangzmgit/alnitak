interface AddHistoryType {
  vid: number,
  part: number,
  time: number
}

interface HistoryVideoType {
  vid: number;
  uid: number;
  title: string;
  cover: string;
  desc: string;
  time: number;
  updatedAt: string;

  // 观看日期，后端不返回
  viewingDate?:string;
}	