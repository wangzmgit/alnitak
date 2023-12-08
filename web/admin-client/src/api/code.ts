import request from '@/utils/request';

// 发送验证码
export const sendEmailCodeAPI = (email: string) => {
  return request.post("v1/verify/getEmailCode", { email })
}