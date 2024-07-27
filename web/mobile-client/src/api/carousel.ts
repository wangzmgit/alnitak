import request from '@/utils/request';

// 获取轮播图
export const getCarouselAPI = (partitionId: number | string) => {
  return request.get(`v1/carousel/getCarousel?partitionId=${partitionId}`)
}