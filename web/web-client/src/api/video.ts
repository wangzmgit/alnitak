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
export const getAllVideoAPI = (page: number = 1, pageSize: number = 10) => {
  return request.get(`v1/video/getAllVideoList?page=${page}&pageSize=${pageSize}`);
}

// 删除视频
export const deleteVideoAPI = (id: number) => {
  return request.delete(`v1/video/deleteVideo/${id}`);
}

// 获取稿件列表
export const getUploadVideoAPI = (page: number, pageSize: number, category: string = 'all') => {
  return request.get(`v1/video/getUploadVideo?page=${page}&pageSize=${pageSize}&category=${category}`);
}

//获取我的视频
export const getVideoByUser = (userId: number, page: number, pageSize: number) => {
  return request.get(`v1/video/getVideoByUser?userId=${userId}&page=${page}&pageSize=${pageSize}`);
}

export const asyncGetVideoByUser = async (userId: number, page: number, pageSize: number) => {
  return await useAsyncData(`video-by-user-${userId}-${page}-${pageSize}`, () => $fetch(`${baseURL}/api/v1/video/getVideoByUser?userId=${userId}&page=${page}&pageSize=${pageSize}`));
}

// 获取视频信息
export const getVideoInfoAPI = (videoId: number | string) => {
  return request.get(`v1/video/getVideoById?vid=${videoId}`);
}

// 获取视频信息（SSR）
export const asyncGetVideoInfoAPI = async (videoId: number | string) => {
  return await useAsyncData(`video-info-${videoId}`, () => $fetch(`${baseURL}/api/v1/video/getVideoById?vid=${videoId}`));
}

// 获取视频支持的分辨率
export const getResourceQualityApi = async (resourceId: number | string) => {
  return request.get(`v1/video/getResourceQuality?resourceId=${resourceId}`)
}

// 获取视频文件URL（HLS格式）
export const getVideoFileUrl = (resourceId: number | string, quality: string, ts?: number) => {
  const cacheBust = ts ? `&ts=${ts}` : '';
  return `${baseURL}/api/v1/video/getVideoFile?resourceId=${resourceId}&quality=${quality}&format=m3u8${cacheBust}`;
}

// 获取视频文件URL（DASH格式 - 返回MPD XML）
export const getVideoFileUrlDash = (resourceId: number | string, quality: string, ts?: number) => {
  const cacheBust = ts ? `&ts=${ts}` : '';
  return `${baseURL}/api/v1/video/getVideoFile?resourceId=${resourceId}&quality=${quality}&format=dash${cacheBust}`;
}

// 获取统一DASH MPD URL（所有清晰度合并到一个MPD，用于无缝切换）
export const getVideoFileUrlDashUnified = (resourceId: number | string, ts?: number) => {
  const cacheBust = ts ? `&ts=${ts}` : '';
  return `${baseURL}/api/v1/video/getVideoFile?resourceId=${resourceId}&format=dash-unified${cacheBust}`;
}

// 获取视频播放信息（JSON格式，类似B站）
export const getVideoPlayInfo = (resourceId: number | string, quality: string) => {
  return request.get(`v1/video/getVideoFile?resourceId=${resourceId}&quality=${quality}`);
}

// 获取热门视频
export const asyncGetHotVideoAPI = async (page: number, pageSize: number) => {
  return await useAsyncData(`hot-video-${page}-${pageSize}`, () => $fetch(`${baseURL}/api/v1/video/getHotVideo?page=${page}&pageSize=${pageSize}`));
}

// 获取热门视频
export const getHotVideoAPI = (page: number, pageSize: number) => {
  return request.get(`v1/video/getHotVideo?page=${page}&pageSize=${pageSize}`);
}

// 获取最近上传的视频（SSR）
export const asyncGetLatestVideoAPI = async (page: number, pageSize: number) => {
  return await useAsyncData(`latest-video-${page}-${pageSize}`, () => $fetch(`${baseURL}/api/v1/video/getLatestVideo?page=${page}&pageSize=${pageSize}`));
}

// 获取最近上传的视频
export const getLatestVideoAPI = (page: number, pageSize: number) => {
  return request.get(`v1/video/getLatestVideo?page=${page}&pageSize=${pageSize}`);
}

// 获取分区视频
export const asyncGetVideoByPartitionAPI = async (size: number, partitionId: number | string) => {
  return await useAsyncData(`partition-video-${partitionId}-${size}`, () => $fetch(`${baseURL}/api/v1/video/getVideoListByPartition?size=${size}&partitionId=${partitionId}`));
}

// 获取分区视频
export const getVideoByPartitionAPI = (size: number, partitionId: number | string) => {
  return request.get(`v1/video/getVideoListByPartition?size=${size}&partitionId=${partitionId}`);
}

// 获取相关推荐视频
export const asyncGetRelatedVideoList = async (videoId: number) => {
  return await useAsyncData(`related-video-${videoId}`, () => $fetch(`${baseURL}/api/v1/video/getRelatedVideoList?vid=${videoId}`));
}

// 搜索视频
export const searchVideoAPI = (data: SearchVideoType) => {
  return request.post("v1/video/searchVideo", data);
}

/** 播放授权：换发 grant，在 expires 前续期即可无感 */
export const postPlayGrantAPI = (resourceShortId: string) => {
  return request.post('v1/play/grant', { resourceShortId });
}

/** 换取当前分 P 的音视频直链（query 需带 token） */
export const getPlayUrlsAPI = (resourceShortId: string, token: string, quality?: string) => {
  return request.get(`v1/play/${encodeURIComponent(resourceShortId)}`, {
    params: { token, ...(quality ? { quality } : {}) },
  });
}