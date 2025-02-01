import request from '@/utils/request';

// 获取@消息
export const getAtMsgAPI = (page: number, pageSize: number) => {
  return request.get(`v1/message/getAtMsg?page=${page}&pageSize=${pageSize}`);
}