import request, { baseURL } from '@/utils/request';
import { storageData } from '@/utils/storage-data';

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

// 获取视频支持的分辨率
export const getResourceQualityApi = async (resourceId: number | string) => {
  return request.get(`v1/video/getResourceQualityManage?resourceId=${resourceId}`)
}

// 获取视频文件URL
export const getVideoFileUrl = (resourceId: number | string, quality: string) => {
  return `${baseURL}/api/v1/video/getVideoFileManage?resourceId=${resourceId}&quality=${quality}`
}

// 获取视频文件URL
export const getVideoFileAPI = (src:string) => {
  return request.get(src);
}
