import request from "@/utils/request";

// 发表评论/回复
export const addCommentAPI = (addComment: AddCommentType) => {
  return request.post('v1/comment/addComment', addComment);
}

// 获取评论
export const getCommentAPI = (vid: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/getComment?vid=${vid}&page=${page}&pageSize=${pageSize}`);
}

// 获取回复
export const getReplyAPI = (commentId: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/getReply?commentId=${commentId}&page=${page}&pageSize=${pageSize}`);
}

// 获取评论/回复
export const deleteCommentAPI = (commentId: number) => {
  return request.delete(`v1/comment/deleteComment/${commentId}`);
}

// 获取评论
export const getCommentListAPI = (vid: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/getCommentList?vid=${vid}&page=${page}&pageSize=${pageSize}`);
}