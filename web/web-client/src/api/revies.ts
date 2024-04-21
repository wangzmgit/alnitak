import request from '@/utils/request';

// 获取审核记录
export const getReviewRecordAPI = (videoId: number) => {
  return request.get(`v1/review/getReviewRecord?vid=${videoId}`);
}
