import request from '@/utils/request';

// 视频是否收藏
export const getCollectVideoStatusAPI = (vid: number) => {
  return request.get(`v1/archive/video/hasCollect?vid=${vid}`);
}

// 已收藏的文件夹
export const getCollectVideoInfoAPI = (vid: number) => {
  return request.get(`v1/archive/video/getCollectInfo?vid=${vid}`);
}

// 视频收藏
export const collectVideoAPI = (collect: CollectType) => {
  return request.post('v1/archive/video/collect', collect);
}

// 文章是否收藏
export const getCollectArticleStatusAPI = (aid: number) => {
  return request.get(`v1/archive/article/hasCollect?aid=${aid}`);
}

// 文章收藏
export const collectArticleAPI = (aid: number) => {
  return request.post('v1/archive/article/collect', { aid });
}

// 文章取收藏
export const cancelCollectArticleAPI = (aid: number) => {
  return request.post('v1/archive/article/cancelCollect', { aid })
}
