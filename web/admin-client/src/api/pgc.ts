import request from '@/utils/request';

// ---- PGC 审核 ----
export const getPGCReviewListAPI = (data: PGCReviewListParam) => {
  return request.post('v1/pgc/getReviewList', data);
};

export const reviewPGCApprovedAPI = (data: PGCReviewActionParam) => {
  return request.post('v1/pgc/reviewApproved', data);
};

export const reviewPGCFailedAPI = (data: PGCReviewActionParam) => {
  return request.post('v1/pgc/reviewFailed', data);
};

// ---- PGC 内容管理 ----
export const getPGCManageListAPI = (data: PGCManageListParam) => {
  return request.post('v1/pgc/getManageList', data);
};

export const adminUpdatePGCStatusAPI = (data: PGCUpdateStatusParam) => {
  return request.post('v1/pgc/adminUpdateStatus', data);
};

export const adminDeletePGCAPI = (pgcId: string) => {
  return request.delete(`v1/pgc/adminDelete/${pgcId}`);
};
