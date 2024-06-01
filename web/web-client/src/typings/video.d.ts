interface BaseVideoType {
  vid: number;
  title: string;
  cover: string;
  desc: string;
  createdAt: string;
}

// 上传视频
interface UploadVideoType {
  title: string;
  cover: string;
  desc: string;
  tags: string;
  copyright: boolean;
  partitionId: number;
}

// 编辑视频
interface EditVideoType {
  vid: number;
  title: string;
  cover: string;
  desc: string;
  tags: string;
}

// 视频状态
interface VideoStatusType extends BaseVideoType {
  tags: string;
  status: number;
  copyright: boolean;
  partitionId: number;
  resources: ResourceType[];
}

// 投稿视频信息
interface ManuscriptVideoType extends BaseVideoType {
  status: number;
  clicks: number;
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

// 全部视频列表
interface AllVideoType {
  vid: number;
  title: string;
}

// 搜索视频
interface SearchVideoType {
  page: number;
  pageSize: number;
  keywords: string;
}