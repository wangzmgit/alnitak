import request from '@/utils/request';

// 视频是否点赞
export const getLikeVideoStatusAPI = (vid: number) => {
  return request.get(`v1/archive/video/hasLike?vid=${vid}`);
}

// 视频点赞
export const likeVideoAPI = (vid: number) => {
  return request.post('v1/archive/video/like', { vid });
}

// 视频取消赞
export const cancelLikeVideoAPI = (vid: number) => {
  return request.post('v1/archive/video/cancelLike', { vid })
}

// 文章是否点赞
export const getLikeArticleStatusAPI = (aid: number) => {
  return request.get(`v1/archive/article/hasLike?aid=${aid}`);
}

// 文章点赞
export const likeArticleAPI = (aid: number) => {
  return request.post('v1/archive/article/like', { aid });
}

// 文章取消赞
export const cancelLikeArticleAPI = (aid: number) => {
  return request.post('v1/archive/article/cancelLike', { aid })
}
