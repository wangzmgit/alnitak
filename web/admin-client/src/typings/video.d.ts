interface VideoListParam {
  page: number;
  pageSize: number;
}

// 视频信息
interface VideoType  {
  vid: number;
  title: string;
  cover: string;
  desc: string;
  createdAt: string;
  tags: string;
  clicks: number;
  copyright: boolean;
  author: UserInfoType;
  resources: ResourceType[];
}