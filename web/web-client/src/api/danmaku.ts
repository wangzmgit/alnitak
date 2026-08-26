import request from "@/utils/request";

//发送弹幕
export const sendDanmakuAPI = (danmaku: DanmakuType) => {
  return request.post('v1/danmaku/sendDanmaku', danmaku);
}

//获取弹幕
export const getDanmakuAPI = (vid: string | number, part?: number, rid?: string) => {
  let url = `v1/danmaku/getDanmaku?vid=${vid}`;
  if (rid) {
    url += `&rid=${rid}`;
  } else if (part) {
    url += `&part=${part}`;
  }
  return request.get(url);
}
