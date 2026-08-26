interface PGCReviewListParam {
  page: number;
  pageSize: number;
}

interface PGCReviewActionParam {
  pgc_id: string;
}

/** 与后端 formatPGCContent 字段一致（snake_case） */
interface PGCReviewRow {
  id?: number;
  pgc_id: string;
  title: string;
  cover: string;
  desc: string;
  pgc_type: number;
  status: number;
  year: number;
  area: string;
  created_at?: string;
  updated_at?: string;
}

interface PGCManageListParam {
  page: number;
  pageSize: number;
  pgc_type?: number;
  status?: number;
  keyword?: string;
}

interface PGCUpdateStatusParam {
  pgc_id: string;
  status: number;
}
