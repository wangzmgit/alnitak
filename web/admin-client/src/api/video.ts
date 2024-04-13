import request from '@/utils/request';

// 获取视频列表
export const getVideoListAPI = (data: VideoListParam) => {
  return request.post("v1/video/getVideoListManage", data);
}

// 删除视频
export const deleteVideoAPI = (id: number) => {
  return request.delete(`v1/video/deleteVideoManage/${id}`);
}