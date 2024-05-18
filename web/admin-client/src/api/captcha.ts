import request from "@/utils/request";

// 获取滑块
export const getSliderAPI = (captchaId: string) => {
  return request.get(`v1/verify/captcha/get?captchaId=${captchaId}`);
}

// 验证滑块
export const validateSliderAPI = (captchaId: string, x: number) => {
  return request.post("v1/verify/captcha/validate", { captchaId, x });
}