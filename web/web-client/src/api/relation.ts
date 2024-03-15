import request from '@/utils/request';

// 获取用户关系
export const getUserRelationAPI = (userId: number) => {
  return request.get(`v1/relation/getUserRelation?userId=${userId}`);
}

// 关注
export const followAPI = (userId: number) => {
  return request.post("v1/relation/follow", { id: userId });
}

// 取消关注
export const unfollowAPI = (userId: number) => {
  return request.post("v1/relation/unfollow", { id: userId });
}