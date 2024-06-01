import request from '@/utils/request';

//获取回复消息
export const getReplyMsgAPI = (page: number, pageSize: number) => {
  return request.get(`v1/message/getReplyMsg?page=${page}&pageSize=${pageSize}`);
}