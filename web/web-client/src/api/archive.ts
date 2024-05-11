import request from '@/utils/request';

// 获取视频点赞收藏数据
export const getVideoArchiveStatAPI = (vid: number) => {
  return request.get(`v1/archive/video/stat?vid=${vid}`);
}

// 获取文章点赞收藏数据
export const getArticleArchiveStatAPI = (aid: number | string) => {
  return request.get(`v1/archive/article/stat?aid=${aid}`);
}