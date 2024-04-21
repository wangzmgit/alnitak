interface ReviewListParam {
  page: number;
  pageSize: number;
}

interface ReviewParam {
  vid: number;
  status?: number;
  remark?: string;
}