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

// 重置密码验证
export const resetpwdCheckAPI = (data: ModifyPwdCheckType) => {
  return request.post('v1/auth/resetpwdCheck', data);
}

// 重置密码
export const mpdifyPwdAPI = (data: ModifyPwdType) => {
  return request.post('v1/auth/modifyPwd', data);
}

