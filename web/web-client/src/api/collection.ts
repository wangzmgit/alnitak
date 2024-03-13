import request from '@/utils/request';

// 创建收藏夹
export const addCollectionAPI = (name: string) => {
  return request.post('v1/collection/addCollection', { name });
}

// 获取收藏夹列表
export const getCollectionListAPI = () => {
  return request.get('v1/collection/getCollectionList');
}


