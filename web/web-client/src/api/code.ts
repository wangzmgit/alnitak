import request from '@/utils/request';

// 发送验证码
export const sendEmailCodeAPI = (data: SendCodeType) => {
  return request.post("v1/verify/getEmailCode", data)
}