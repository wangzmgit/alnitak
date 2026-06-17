import request, { baseURL } from '@/utils/request';

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

// 获取视频文件URL (HLS m3u8)
export const getVideoFileUrl = (resourceId: number | string, quality: string) => {
  return `${baseURL}/api/v1/video/getVideoFileManage?resourceId=${resourceId}&quality=${quality}`
}

// 获取视频文件URL (DASH mpd)
export const getVideoFileUrlDash = (resourceId: number | string, quality: string) => {
  return `${baseURL}/api/v1/video/getVideoFileManage?resourceId=${resourceId}&quality=${quality}&format=mpd`
}

// 获取统一DASH MPD URL（所有清晰度合并到一个MPD，用于无缝切换）
export const getVideoFileUrlDashUnified = (resourceId: number | string) => {
  return `${baseURL}/api/v1/video/getVideoFileManage?resourceId=${resourceId}&format=dash-unified`
}

// 获取视频文件URL
export const getVideoFileAPI = (src:string) => {
  return request.get(src, { responseType: 'text', transformResponse: [(data: any) => data] });
}

// 获取处理失败的视频列表
export const getFailedVideoListAPI = (data: VideoListParam) => {
  return request.post("v1/video/getFailedVideoList", data);
}

// 获取处理中视频列表
export const getProcessingVideoListAPI = (data: VideoListParam) => {
  return request.post("v1/video/getProcessingVideoList", data);
}

// 重新转码视频
export const reTranscodeVideoAPI = (vid: number, resourceId?: number) => {
	const query = new URLSearchParams();
	query.set('vid', String(vid));
	if (typeof resourceId === 'number') query.set('resourceId', String(resourceId));
	return request.post(`v1/video/reTranscodeVideo?${query.toString()}`);
}

// 重新上传OSS（转码成功但上传失败时重试）
export const reUploadVideoAPI = (vid: number) => {
	const query = new URLSearchParams();
	query.set('vid', String(vid));
	return request.post(`v1/video/reUploadVideo?${query.toString()}`);
}
