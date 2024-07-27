import request from '@/utils/request';
import { baseURL } from '@/utils/request';

// 获取热门视频
export const getHotVideoAPI = (page: number, pageSize: number) => {
  return request.get(`v1/video/getHotVideo?page=${page}&pageSize=${pageSize}`);
}

// 获取分区视频
export const getVideoByPartitionAPI = (size: number, partitionId: number | string) => {
  return request.get(`v1/video/getVideoListByPartition?size=${size}&partitionId=${partitionId}`);
}

// 获取稿件列表
export const getUploadVideoAPI = (page: number, pageSize: number) => {
  return request.get(`v1/video/getUploadVideo?page=${page}&pageSize=${pageSize}`);
}

// 搜索视频
export const searchVideoAPI = (data: SearchVideoType) => {
  return request.post("v1/video/searchVideo", data);
}

// 获取视频信息
export const getVideoInfoAPI = async (videoId: number | string) => {
  return request.get(`v1/video/getVideoById?vid=${videoId}`);
}

// 获取视频支持的分辨率
export const getResourceQualityApi = async (resourceId: number | string) => {
  return request.get(`v1/video/getResourceQuality?resourceId=${resourceId}`)
}

// 获取视频文件URL
export const getVideoFileUrl = (resourceId: number, quality: string) => {
  return `${baseURL}/api/v1/video/getVideoFile?resourceId=${resourceId}&quality=${quality}`;
}
