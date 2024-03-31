import request from '@/utils/request';

// 获取点赞消息
export const getLikeMsgAPI = (page: number, pageSize: number) => {
  return request.get(`v1/message/getLikeMsg?page=${page}&pageSize=${pageSize}`);
}