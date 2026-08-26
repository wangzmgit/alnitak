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

// 刷新 token：可传本地 refresh；不传时由后端从 HttpOnly Cookie 读取（阶段 B）
export function updateTokenAPI(refreshToken?: string) {
  return request.post('v1/auth/updateToken', refreshToken ? { refreshToken } : {});
}

// 当前会话（Cookie / Authorization 双栈，供 fetchMe 与 SSR 对齐）
export function getAuthMeAPI() {
  return request.get('v1/auth/me');
}

// 退出登录：可传本地 refresh；不传则仅靠 Cookie 即可完成吊销与清 Cookie
export function logoutAPI(refreshToken?: string) {
  return request.post('v1/auth/logout', refreshToken ? { refreshToken } : {});
}

// 重置密码验证
export const resetpwdCheckAPI = (data: ModifyPwdCheckType) => {
  return request.post('v1/auth/resetpwdCheck', data);
}

// 重置密码
export const mpdifyPwdAPI = (data: ModifyPwdType) => {
  return request.post('v1/auth/modifyPwd', data);
}

// 修改密码（已登录用户，需旧密码）
export const changePasswordAPI = (data: { oldPassword: string; newPassword: string }) => {
  return request.post('v1/auth/changePassword', data);
}

