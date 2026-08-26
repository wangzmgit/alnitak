import request from '@/utils/request';

// 获取在线 Worker 列表
export const getWorkersAPI = () => {
  return request.get('v1/admin/workers');
}
