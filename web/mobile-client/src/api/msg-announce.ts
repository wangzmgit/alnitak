import request from '@/utils/request';

//获取公告
export const getAnnounceAPI = (page: number, pageSize: number) => {
  return request.get(`v1/message/getAnnounce?page=${page}&pageSize=${pageSize}`);
}