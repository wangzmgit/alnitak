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
}
