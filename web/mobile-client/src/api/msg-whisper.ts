import request from "@/utils/request";

// 获取私信列表
export const getWhisperListAPI = () => {
  return request.get('v1/message/getWhisperList')
}

// 获取私信列表
export const getWhisperDetailsAPI = (fid: number, page: number, pageSize: number) => {
  return request.get(`v1/message/getWhisperDetails?fid=${fid}&page=${page}&pageSize=${pageSize}`)
}

// 发送私信
export const sendWhisperAPI = (msg: WhisperType) => {
  return request.post('v1/message/sendWhisper', msg)
}

// 已读私信
export const readWhisperAPI = (id: number) => {
  return request.post('v1/message/readWhisper', { id })
}