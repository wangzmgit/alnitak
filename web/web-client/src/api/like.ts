import request from '@/utils/request';

// 是否点赞
export const getLikeStatusAPI = (vid: number) => {
  return request.get(`v1/archive/hasLike?vid=${vid}`);
}

//点赞
export const likeAPI = (id: number) => {
  return request.post('v1/archive/like', { id });
}

//取消赞
export const cancelLikeAPI = (id: number) => {
  return request.post('v1/archive/cancelLike', { id })
}
