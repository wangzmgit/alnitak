import request from '@/utils/request';

// 是否点赞
export const getLikeVideoStatusAPI = (vid: number) => {
  return request.get(`v1/archive/video/hasLike?vid=${vid}`);
}

//点赞
export const likeVideoAPI = (vid: number) => {
  return request.post('v1/archive/video/like', { vid });
}

//取消赞
export const cancelLikeVideoAPI = (vid: number) => {
  return request.post('v1/archive/video/cancelLike', { vid })
}
