import request from '@/utils/request';

//获取分区
export const getPartitionAPI = (type?: "video" | "article") => {
  const typeParam = `?type=${type === 'article' ? 1 : 0}`;
  return request.get(`v1/partition/getPartitionList` + typeParam);
}
