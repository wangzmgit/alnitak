import request from '@/utils/request';

// 获取视频点赞收藏数据
export const getVideoArchiveStatAPI = (vid: number | string) => {
  return request.get(`v1/archive/video/stat?vid=${vid}`);
}

// 获取文章点赞收藏数据
export const getArticleArchiveStatAPI = (aid: number | string) => {
  return request.get(`v1/archive/article/stat?aid=${aid}`);
}

// 视频分享计数
export const shareVideoAPI = (vid: number) => {
  return request.post(`v1/archive/video/share`, { vid });
}

// 文章分享计数
export const shareArticleAPI = (aid: number) => {
  return request.post(`v1/archive/article/share`, { aid });
}