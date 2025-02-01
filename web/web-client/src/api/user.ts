import request from '@/utils/request';
import { baseURL } from '@/utils/request';
import { useAsyncData } from 'nuxt/app';

// 获取用户信息
export const getUserInfoAPI = () => {
  return request.get('v1/user/getUserInfo');
}

// 编辑用户信息
export const editUserInfoAPI = (data: EditUserInfoType) => {
  return request.put("v1/user/editUserInfo", data);
}

// 获取用户基本信息
export const getUserBaseInfoAPI = (userId: number) => {
  return request.get(`v1/user/getUserBaseInfo?userId=${userId}`);
}

// 获取用户基本信息
export const asyncGetUserBaseInfoAPI = async (userId: number | string) => {
  return await useAsyncData(() => $fetch(`${baseURL}/api/v1/user/getUserBaseInfo?userId=${userId}`));
}