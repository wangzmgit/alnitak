import request from '@/utils/request';

// 获取用户信息
export const getUserInfoAPI = () => {
  return request.get('v1/user/getUserInfo');
}

// 获取用户基本信息
export const getUserBaseInfoAPI = (userId: number) => {
  return request.get(`v1/user/getUserBaseInfo?userId=${userId}`);
}

