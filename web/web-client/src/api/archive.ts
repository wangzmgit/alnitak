import request from '@/utils/request';

// 获取点赞收藏数据
export const getArchiveStatAPI = (vid: number) => {
  return request.get(`v1/archive/stat?vid=${vid}`);
}