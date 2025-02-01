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
export const emailLoginAPI = (login: UserLoginType) => {
  return request.post('v1/auth/login/email', login);
}

// 刷新token
export function updateTokenAPI(refreshToken: string) {
  return request.post('v1/auth/updateToken', { refreshToken });
}

// 退出登录
export function logoutAPI(refreshToken: string) {
  return request.post('v1/auth/logout', { refreshToken });
}

// 清除Cookie
export function clearCookieAPI() {
  return request.post('v1/auth/clearCookie',);
}
