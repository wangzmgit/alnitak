import request from '@/utils/request';
import { baseURL } from '@/utils/request';
import { useAsyncData } from 'nuxt/app';

// 发布文章
export const uploadArticleInfoAPI = (uploadArticle: UploadArticleType) => {
  return request.post('v1/article/uploadArticleInfo', uploadArticle);
}

// 编辑文章
export const editArticleAPI = (editArticle: EditArticleType) => {
  return request.put("v1/article/editArticleInfo", editArticle);
}

// 获取稿件列表
export const getUploadArticleAPI = (page: number, pageSize: number) => {
  return request.get(`v1/article/getUploadArticle?page=${page}&pageSize=${pageSize}`);
}

// 获取上传文章信息
export const getArticleStatusAPI = (aid: number) => {
  return request.get(`v1/article/getArticleStatus?aid=${aid}`);
}

// 获取所有文章列表
export const getAllArticleAPI = () => {
  return request.get("v1/article/getAllArticleList");
}

// 获取文章信息
export const asyncGetArticleInfoAPI = async (articleId: number | string) => {
  return await useAsyncData(() => $fetch(`${baseURL}/api/v1/article/getArticleById?aid=${articleId}`));
}

// 获取随机文章
export const asyncGetRandomArticleAPI = async (size: number) => {
  return await useAsyncData(() => $fetch(`${baseURL}/api/v1/article/getRandomArticleList?size=${size}`));
}

// 获取随机文章
export const getRandomArticleAPI = async (size: number) => {
  return request.get(`v1/article/getRandomArticleList?size=${size}`);
}

// 删除文章
export const deleteArticleAPI = (id: number) => {
  return request.delete(`v1/article/deleteArticle/${id}`);
}
