interface ReviewListParam {
  page: number;
  pageSize: number;
}

interface ReviewArticleListParam {
  page: number;
  pageSize: number;
}

interface ReviewVideoParam {
  vid: number;
  status?: number;
  remark?: string;
}

interface ReviewArticleParam {
  aid: number;
  status?: number;
  remark?: string;
}