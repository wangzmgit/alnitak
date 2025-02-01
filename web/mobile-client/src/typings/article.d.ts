// 上传视频
interface UploadArticleType {
  title: string;
  cover: string;
  content: string;
  tags: string;
  copyright: boolean;
  partitionId: number;
}

// 编辑视频
interface EditArticleType {
  aid: number;
  title: string;
  cover: string;
  content: string;
  tags: string;
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

// 全部文章列表
interface AllArticleType  {
  aid: number;
  title: string;
}
