import request from '@/utils/request';

// 是否点赞
export const getLikeStatusAPI = (vid: number) => {
  return request.get(`v1/archive/hasLike?vid=${vid}`);
}

//点赞
export const likeAPI = (vid: number) => {
  return request.post('v1/archive/like', { vid });
}

//取消赞
export const cancelLikeAPI = (vid: number) => {
  return request.post('v1/archive/cancelLike', { vid })
}
