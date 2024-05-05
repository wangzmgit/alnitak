import request from '@/utils/request';

// 审核通过
export const reviewVideoApprovedAPI = (data: ReviewVideoParam) => {
  return request.post("v1/review/reviewVideoApproved", data)
}

// 审核不通过
export const reviewVideoFailedAPI = (data: ReviewVideoParam) => {
  return request.post("v1/review/reviewVideoFailed", data)
}

// 审核通过
export const reviewArticleApprovedAPI = (data: ReviewArticleParam) => {
  return request.post("v1/review/reviewArticleApproved", data)
}

// 审核不通过
export const reviewArticleFailedAPI = (data: ReviewArticleParam) => {
  return request.post("v1/review/reviewArticleFailed", data)
}


