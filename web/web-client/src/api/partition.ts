import request from '@/utils/request';
import { baseURL } from '@/utils/request';
import { useAsyncData } from 'nuxt/app';

//获取分区
export const getPartitionAPI = (type?: "video" | "article") => {
  const typeParam = `?type=${type === 'article' ? 1 : 0}`;
  return request.get(`v1/partition/getPartitionList` + typeParam);
}

//新增分区
export const addPartitionAPI = (addPartition: AddPartitionType) => {
  return request.post(`v1/partition/addPartition`, addPartition);
}

//删除分区
export const deletePartitionAPI = (id: number) => {
  return request.post(`v1/partition/deletePartition`, { id });
}

export const asyncGetPartition = async () => {
  return await useAsyncData(() => $fetch(`${baseURL}/api/v1/partition/getPartitionList`));
}