import request from '@/utils/request';
import { baseURL } from '@/utils/request';
import { useAsyncData } from 'nuxt/app';

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

// 获取关注数据
export const getFollowDataAPI = (uid: number | string) => {
  return request.get(`v1/relation/getFollowCount?userId=${uid}`)
}

// 获取关注数据
export const asyncetFollowDataAPI = async (userId: number | string) => {
  return await useAsyncData(() => $fetch(`${baseURL}/api/v1/relation/getFollowCount?userId=${userId}`));
}

//获取关注列表
export const getFollowingAPI = (userId: number, page: number, pageSize: number) => {
  return request.get(`v1/relation/getFollowings?userId=${userId}&page=${page}&pageSize=${pageSize}`);
}

//获取粉丝列表
export const getFollowersAPI = (userId: number, page: number, pageSize: number) => {
  return request.get(`v1/relation/getFollowers?userId=${userId}&page=${page}&pageSize=${pageSize}`);
}