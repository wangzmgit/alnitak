import request from '@/utils/request';

// 获取公告
export const getAnnounceListAPI = (page: number, pageSize: number) => {
  return request.get(`v1/message/getAnnounce?page=${page}&pageSize=${pageSize}`);
}

// 添加公告
export const addAnnounceAPI = (addAnnounce: AddAnnounceType) => {
  return request.post('v1/message/addAnnounce', addAnnounce);
}

// 删除公告
export const deleteAnnounceAPI = (id: number) => {
  return request.delete(`v1/message/deleteAnnounce/${id}`);
}