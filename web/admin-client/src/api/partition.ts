import request from '@/utils/request';

// 获取分区
export const getPartitionAPI = () => {
  return request.get(`v1/partition/getPartitionList`);
}

// 新增分区
export const addPartitionAPI = (data: AddPartitionType) => {
  return request.post(`v1/partition/addPartition`, data);
}

// 删除分区
export const deletePartitionAPI = (id: number) => {
  return request.delete(`v1/partition/deletePartition/${id}`);
}
