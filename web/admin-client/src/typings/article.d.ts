interface ArticleListParam {
  page: number;
  pageSize: number;
}

interface ArticleType {
  aid: number;
  title: string;
  cover: string;
  content: string;
  tags: string;
  clicks: number;
  status: number;
  copyright: boolean;
  partitionId: number;
  author: UserInfoType;
  createdAt: string;
}