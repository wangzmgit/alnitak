interface BaseVideoType {
  vid: number;
  title: string;
  cover: string;
  desc: string;
  createdAt: string;
}

// 视频信息
interface VideoType extends BaseVideoType {
  tags: string;
  clicks: number;
  status: number;
  copyright: boolean;
  duration: number;
  author: UserInfoType;
  resources: ResourceType[];
}

// 搜索视频
interface SearchVideoType {
  page: number;
  pageSize: number;
  keywords: string;
}