interface BaseVideoType {
  vid: number;
  // 短 ID（后端返回 shortId，可选）
  shortId?: string;
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
  /** 后端可为 string[] 或历史 CSV 字符串 */
  tags?: string | string[];
  status: number;
  copyright: boolean;
  partitionId: number;
  resources: ResourceType[];
}

// 投稿视频信息
interface ManuscriptVideoType extends BaseVideoType {
  status: number;
  clicks: number;
  /** 总时长（秒，整数） */
  duration?: number;
  transcodingProgress?: number;
  transcodingDetails?: TranscodingProgressDetail[];
}

interface TranscodingProgressDetail {
  resourceId: number;
  resourceTitle: string;
  quality: string;
  progress: number;
  status: 'processing' | 'success' | 'fail' | string;
  upload?: UploadProgressInfo;
}

// 视频信息
interface VideoType extends BaseVideoType {
  /** 后端 JSON 多为 string[]；历史/缓存可能为 CSV 字符串 */
  tags?: string | string[];
  clicks: number;
  /** 详情页：作者粉丝数（与 author 分离） */
  fans?: number;
  status?: number;
  copyright: boolean;
  /** 总时长（秒，整数） */
  duration: number;
  author: UserInfoType;
  resources: ResourceType[];
}

// 全部视频列表
interface AllVideoType {
  shortId: string;
  title: string;
}

// 搜索视频
interface SearchVideoType {
  page: number;
  pageSize: number;
  keywords: string;
}