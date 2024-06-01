import request from '@/utils/request';
import { baseURL } from '@/utils/request';
import { useAsyncData } from 'nuxt/app';

// 上传视频信息
export const uploadVideoInfoAPI = (uploadVideo: UploadVideoType) => {
  return request.post('v1/video/uploadVideoInfo', uploadVideo);
}

// 获取上传视频信息
export const getVideoStatusAPI = (vid: number) => {
  return request.get(`v1/video/getVideoStatus?vid=${vid}`);
}

// 提交审核
export const submitReviewAPI = (id: number) => {
  return request.post('v1/video/submitReview', { id });
}

// 编辑视频
export const editVideoAPI = (editVideo: EditVideoType) => {
  return request.put("v1/video/editVideoInfo", editVideo);
}

// 获取所有视频列表
export const getAllVideoAPI = () => {
  return request.get("v1/video/getAllVideoList");
}

// 删除视频
export const deleteVideoAPI = (id: number) => {
  return request.delete(`v1/video/deleteVideo/${id}`);
}

// 获取稿件列表
export const getUploadVideoAPI = (page: number, pageSize: number) => {
  return request.get(`v1/video/getUploadVideo?page=${page}&pageSize=${pageSize}`);
}

//获取我的视频
export const getVideoByUser = (userId: number, page: number, pageSize: number) => {
  return request.get(`v1/video/getVideoByUser?userId=${userId}&page=${page}&pageSize=${pageSize}`);
}

export const asyncGetVideoByUser = async (userId: number, page: number, pageSize: number) => {
  return await useAsyncData(() => $fetch(`${baseURL}/api/v1/video/getVideoByUser?userId=${userId}&page=${page}&pageSize=${pageSize}`));
}

// 获取视频信息
export const asyncGetVideoInfoAPI = async (videoId: number | string) => {
  return await useAsyncData(() => $fetch(`${baseURL}/api/v1/video/getVideoById?vid=${videoId}`));
}

// 获取视频支持的分辨率
export const getResourceQualityApi = async (resourceId: number | string) => {
  return request.get(`v1/video/getResourceQuality?resourceId=${resourceId}`)
}

// 获取视频文件URL
export const getVideoFileUrl = (resourceId: number, quality: string) => {
  return `${baseURL}/api/v1/video/getVideoFile?resourceId=${resourceId}&quality=${quality}`;
}

// 获取热门视频
export const asyncGetHotVideoAPI = async (page: number, pageSize: number) => {
  return await useAsyncData(() => $fetch(`${baseURL}/api/v1/video/getHotVideo?page=${page}&pageSize=${pageSize}`));
}

// 获取热门视频
export const getHotVideoAPI = (page: number, pageSize: number) => {
  return request.get(`v1/video/getHotVideo?page=${page}&pageSize=${pageSize}`);
}

// 获取分区视频
export const asyncGetVideoByPartitionAPI = async (size: number, partitionId: number | string) => {
  return await useAsyncData(() => $fetch(`${baseURL}/api/v1/video/getVideoListByPartition?size=${size}&partitionId=${partitionId}`));
}

// 获取分区视频
export const getVideoByPartitionAPI = (size: number, partitionId: number | string) => {
  return request.get(`v1/video/getVideoListByPartition?size=${size}&partitionId=${partitionId}`);
}

// 获取相关推荐视频
export const asyncGetRelatedVideoList = async (videoId: number) => {
  return await useAsyncData(() => $fetch(`${baseURL}/api/v1/video/getRelatedVideoList?vid=${videoId}`));
}

// 搜索视频
export const searchVideoAPI = (data: SearchVideoType) => {
  return request.post("v1/video/searchVideo", data);
}