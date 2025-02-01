import request from '@/utils/request';

// 获取文章列表
export const getArticleListAPI = (data: ArticleListParam) => {
  return request.post("v1/article/getArticleListManage", data);
}

// 获取文章列表
export const getReviewArticleListAPI = (data: ReviewArticleListParam) => {
  return request.post("v1/article/getReviewArticleList", data);
}

// 删除文章
export const deleteArticleAPI = (id: number) => {
  return request.delete(`v1/article/deleteArticleManage/${id}`);
}