import request from '@/utils/request';

// 获取用户信息
export const getUserInfoAPI = () => {
  return request.get('v1/user/getUserInfo');
}