import request from "@/utils/request";

// 发表评论/回复
export const addVideoCommentAPI = (addComment: AddCommentType) => {
  return request.post('v1/comment/video/addComment', addComment);
}

// 发表文章评论/回复
export const addArticleCommentAPI = (addComment: AddCommentType) => {
  return request.post('v1/comment/article/addComment', addComment);
}

// 获取评论
export const getVideoCommentAPI = (vid: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/video/getComment?vid=${vid}&page=${page}&pageSize=${pageSize}`);
}

// 获取文章评论
export const getArticleCommentAPI = (aid: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/article/getComment?aid=${aid}&page=${page}&pageSize=${pageSize}`);
}

// 获取回复
export const getVideoReplyAPI = (commentId: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/video/getReply?commentId=${commentId}&page=${page}&pageSize=${pageSize}`);
}

// 获取文章回复
export const getArticleReplyAPI = (commentId: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/article/getReply?commentId=${commentId}&page=${page}&pageSize=${pageSize}`);
}

// 删除评论/回复
export const deleteVideoCommentAPI = (commentId: number) => {
  return request.delete(`v1/comment/video/deleteComment/${commentId}`);
}

// 删除文章评论/回复
export const deleteArticleCommentAPI = (commentId: number) => {
  return request.delete(`v1/comment/article/deleteComment/${commentId}`);
}

// 获取视频评论
export const getVideoCommentListAPI = (vid: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/video/getCommentList?vid=${vid}&page=${page}&pageSize=${pageSize}`);
}

// 获取文章评论
export const getArticleCommentListAPI = (aid: number, page: number, pageSize: number) => {
  return request.get(`v1/comment/article/getCommentList?aid=${aid}&page=${page}&pageSize=${pageSize}`);
}