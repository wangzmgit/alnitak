import request from '@/utils/request';

// 获取轮播图列表
export const getCarouselListAPI = (data: CarouselListParam) => {
  return request.post("v1/carousel/getCarouselList", data);
}

// 新增轮播图
export const addCarouselAPI = (data: AddOrEditCarouselType) => {
  return request.post("v1/carousel/addCarousel", data);
}

// 编辑轮播图
export const editCarouselAPI = (data: AddOrEditCarouselType) => {
  return request.put("v1/carousel/editCarousel", data);
}

// 删除轮播图
export const deleteCarouselAPI = (id: string | number) => {
  return request.delete(`v1/carousel/deleteCarousel/${id}`);
}
