import request from '@/utils/request';

// 创建收藏夹
export const addCollectionAPI = (name: string) => {
  return request.post('v1/collection/addCollection', { name });
}

// 获取收藏夹列表
export const getCollectionListAPI = () => {
  return request.get('v1/collection/getCollectionList');
}

// 获取收藏夹
export const getCollectionInfoAPI = (id: number) => {
  return request.get(`v1/collection/getCollectionInfo?id=${id}`);
}

//获取收藏夹内的视频
export const getCollectVideoAPI = (id: number, page: number, page_size: number) => {
  return request.get(`v1/collection/getVideoList?cid=${id}&page=${page}&pageSize=${page_size}`);
}

// 编辑收藏夹
export const editCollectionAPI = (collection: EditCollectionType) => {
  return request.put('v1/collection/editCollection', collection);
}

// 删除收藏夹
export const deleteCollectionAPI = (id: number) => {
  return request.delete(`v1/collection/deleteCollection/${id}`);
}

