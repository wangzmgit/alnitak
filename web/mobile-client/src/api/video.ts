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

// 获取视频信息（vid 可为数字 id 或 shortId）
export const getVideoInfoAPI = async (videoId: number | string) => {
  const v = encodeURIComponent(String(videoId).trim())
  return request.get(`v1/video/getVideoById?vid=${v}`)
}

/** 播放授权（与 Web 对齐，后续播放器可接入 signed URL） */
export const postPlayGrantAPI = (resourceShortId: string) => {
  return request.post('v1/play/grant', { resourceShortId })
}

export const getPlayUrlsAPI = (resourceShortId: string, token: string, quality?: string, audioLang?: string) => {
  return request.get(`v1/play/${encodeURIComponent(resourceShortId)}`, {
    params: { token, ...(quality ? { quality } : {}), ...(audioLang ? { audioLang } : {}) },
  })
}

/** 查询指定资源的可用音轨列表 */
export const getAudioTracksAPI = (resourceId: number | string) => {
  return request.get(`v1/audio/tracks/${resourceId}`);
}

// 获取视频支持的分辨率
export const getResourceQualityApi = async (resourceId: number | string) => {
  return request.get(`v1/video/getResourceQuality?resourceId=${resourceId}`)
}

// 获取视频文件URL（HLS格式）
export const getVideoFileUrl = (resourceId: number, quality: string) => {
  return `${baseURL}/api/v1/video/getVideoFile?resourceId=${resourceId}&quality=${quality}&format=m3u8`;
}

// 获取视频文件URL（DASH格式）
export const getVideoFileUrlDash = (resourceId: number, quality: string) => {
  return `${baseURL}/api/v1/video/getVideoFile?resourceId=${resourceId}&quality=${quality}&format=dash`;
}

// 获取统一DASH MPD URL（所有清晰度合并到一个MPD，用于无缝切换）
export const getVideoFileUrlDashUnified = (resourceId: number) => {
  return `${baseURL}/api/v1/video/getVideoFile?resourceId=${resourceId}&format=dash-unified`;
}
