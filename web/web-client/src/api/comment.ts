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
export const getVideoCommentAPI = (vid: number | string, page: number, pageSize: number) => {
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
export const getVideoCommentListAPI = (id: string | number, page: number, pageSize: number) => {
  let idParam = '';
  if (id) {
    idParam = typeof id === 'number' ? `vid=${id}&` : `shortId=${id}&`;
  }
  return request.get(`v1/comment/video/getCommentList?${idParam}page=${page}&pageSize=${pageSize}`);
}

// 获取文章评论
export const getArticleCommentListAPI = (id: string | number, page: number, pageSize: number) => {
  let idParam = '';
  if (id) {
    idParam = typeof id === 'number' ? `aid=${id}&` : `shortId=${id}&`;
  }
  return request.get(`v1/comment/article/getCommentList?${idParam}page=${page}&pageSize=${pageSize}`);
}

// 点赞视频评论
export const likeVideoCommentAPI = (commentId: number) => {
  return request.post(`v1/comment/video/like/${commentId}`);
}

// 取消点赞视频评论
export const unlikeVideoCommentAPI = (commentId: number) => {
  return request.delete(`v1/comment/video/like/${commentId}`);
}

// 点赞文章评论
export const likeArticleCommentAPI = (commentId: number) => {
  return request.post(`v1/comment/article/like/${commentId}`);
}

// 取消点赞文章评论
export const unlikeArticleCommentAPI = (commentId: number) => {
  return request.delete(`v1/comment/article/like/${commentId}`);
}

// 点踩视频评论
export const dislikeVideoCommentAPI = (commentId: number) => {
  return request.post(`v1/comment/video/dislike/${commentId}`);
}

// 取消点踩视频评论
export const undislikeVideoCommentAPI = (commentId: number) => {
  return request.delete(`v1/comment/video/dislike/${commentId}`);
}

// 点踩文章评论
export const dislikeArticleCommentAPI = (commentId: number) => {
  return request.post(`v1/comment/article/dislike/${commentId}`);
}

// 取消点踩文章评论
export const undislikeArticleCommentAPI = (commentId: number) => {
  return request.delete(`v1/comment/article/dislike/${commentId}`);
}

