import request from "@/utils/request";

// 发表评论/回复
export const addVideoCommentAPI = (addComment: AddCommentType) => {
  return request.post('v1/comment/addVideoComment', addComment);
}

// 获取评论
export const getVideoCommentAPI = (vid: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/getVideoComment?vid=${vid}&page=${page}&pageSize=${pageSize}`);
}

// 获取回复
export const getVideoReplyAPI = (commentId: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/getVideoReply?commentId=${commentId}&page=${page}&pageSize=${pageSize}`);
}

// 获取评论/回复
export const deleteVideoCommentAPI = (commentId: number) => {
  return request.delete(`v1/comment/deleteVideoComment/${commentId}`);
}

// 获取视频评论
export const getVideoCommentListAPI = (vid: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/getVideoCommentList?vid=${vid}&page=${page}&pageSize=${pageSize}`);
}

// 获取文章评论
export const getArticleCommentListAPI = (aid: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/getArticleCommentList?aid=${aid}&page=${page}&pageSize=${pageSize}`);
}