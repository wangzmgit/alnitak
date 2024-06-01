interface CarouselListParam {
  page: number;
  pageSize: number;
}

interface CarouselType {
  id: number;
  img: string;
  title: string;
  url?: string;
  color: string;
  use: boolean;
  createdAt: string;
  partitionId: number;
}

interface AddOrEditCarouselType {
  id?: number;
  img: string;
  title: string;
  url?: string;
  color: string;
  use: boolean;
  partitionId: number;
}
