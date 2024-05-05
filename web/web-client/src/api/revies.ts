import request from '@/utils/request';

// 获取审核记录
export const getVideoReviewRecordAPI = (videoId: number) => {
  return request.get(`v1/review/getVideoReviewRecord?vid=${videoId}`);
}

// 获取审核记录
export const getArticleReviewRecordAPI = (articleId: number) => {
  return request.get(`v1/review/getArticleReviewRecord?aid=${articleId}`);
}
