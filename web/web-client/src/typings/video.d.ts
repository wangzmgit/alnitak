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
  copyright: boolean;
  partitionId: number;
}

// 视频状态
interface VideoStatusType extends BaseVideoType {
  status: number;
  copyright: boolean;
  partition: number;
  resources: ResourceType[];
}

// 投稿视频信息
interface ManuscriptVideoType extends BaseVideoType {
  status: number;
  clicks: number;
}

// 视频信息
// 视频信息
interface VideoType extends BaseVideoType {
  tags: string;
  clicks: number;
  copyright: boolean;
  author: UserInfoType;
  resources: ResourceType[];
}