import request from '@/utils/request';

// 注册
export const registerAPI = (register: UserRegisterType) => {
  return request.post('v1/auth/register', register);
}

// 登录
export const loginAPI = (login: UserLoginType) => {
  return request.post('v1/auth/login', login);
}

// 邮箱登录
// TODO: 实现邮箱登录
export const emailLoginAPI = (login: UserLoginType) => {
  return request.post('v1/user/login/email', login);
}

// 刷新token
export function updateTokenAPI(refreshToken: string) {
  return request.post('v1/auth/updateToken', { refreshToken });
}
