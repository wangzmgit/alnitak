import request from "@/utils/request";

//发送弹幕
export const sendDanmakuAPI = (danmaku: DanmakuType) => {
  return request.post('v1/danmaku/sendDanmaku', danmaku);
}

//获取弹幕
export const getDanmakuAPI = (vid: number, part: number) => {
  return request.get(`v1/danmaku/getDanmaku?vid=${vid}&part=${part}`);
}
