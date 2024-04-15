import request from '@/utils/request';

// 获取视频列表
export const getVideoListAPI = (data: VideoListParam) => {
  return request.post("v1/video/getVideoListManage", data);
}

// 获取审核列表
export const getReviewListAPI = (data: ReviewListParam) => {
  return request.post("v1/video/getReviewList", data);
}

// 获取审核资源列表
export const getReviewResourceListAPI = (vid: number | string) => {
  return request.get(`v1/video/getReviewResourceList?vid=${vid}`);
}

// 删除视频
export const deleteVideoAPI = (id: number) => {
  return request.delete(`v1/video/deleteVideoManage/${id}`);
}