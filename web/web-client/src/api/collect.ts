import request from '@/utils/request';

// 是否收藏
export const getCollectStatusAPI = (vid: number) => {
  return request.get(`v1/archive/hasCollect?vid=${vid}`);
}

// 已收藏的文件夹
export const getCollectInfoAPI = (vid: number) => {
  return request.get(`v1/archive/getCollectInfo?vid=${vid}`);
}

//收藏
export const collectAPI = (collect: CollectType) => {
  return request.post('v1/archive/collect', collect);
}
