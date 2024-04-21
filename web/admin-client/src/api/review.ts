import request from '@/utils/request';

// 审核通过
export const reviewVideoApprovedAPI = (data: ReviewParam) => {
  return request.post("v1/review/reviewVideoApproved", data)
}

// 审核不通过
export const reviewVideoFailedAPI = (data: ReviewParam) => {
  return request.post("v1/review/reviewVideoFailed", data)
}

