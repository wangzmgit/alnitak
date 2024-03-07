import request from '@/utils/request';
import { baseURL } from '@/utils/request';
import { useAsyncData } from 'nuxt/app';

//上传视频信息
export const uploadVideoInfoAPI = (uploadVideo: UploadVideoType) => {
  return request.post('v1/video/uploadVideoInfo', uploadVideo);
}

//获取上传视频信息
export const getVideoStatusAPI = (vid: number) => {
  return request.get(`v1/video/getVideoStatus?vid=${vid}`);
}

//提交审核
export const submitReviewAPI = (id: number) => {
  return request.post('v1/video/submitReview', { id });
}

//提交审核
export const getUploadVideoAPI = (page: number, pageSize: number) => {
  return request.get(`v1/video/getUploadVideo?page=${page}&pageSize=${pageSize}`);
}

// 获取视频信息
export const asyncGetVideoInfoAPI = async (videoId: number | string) => {
  return await useAsyncData(() => $fetch(`${baseURL}/api/v1/video/getVideoById?vid=${videoId}`));
}

