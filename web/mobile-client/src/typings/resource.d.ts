interface BaseResourceType {
  id: number;
  title: string;
}

// 分P
interface ResourceType extends BaseResourceType {
  url: string;
  duration: number;
  status: number;
  quality: number;
  uploading?: boolean;
  percent?: number;
}

// 分P
interface UploadResourceType extends BaseResourceType {
  status: number;
  uploading: boolean;
  percent: number;
}
