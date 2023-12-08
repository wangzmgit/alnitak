import axios from "axios";
import type { AxiosInstance } from "axios";
import { updateTokenAPI } from "@/api/auth";
import { statusCode } from "./status-code";
import { storageData as storage } from "./storage-data";
import { globalConfig as config } from "./global-config";

const baseURL = config.domain ? `http${config.https ? 's' : ''}://${config.domain}` : '';

const service: AxiosInstance = axios.create({
  baseURL: `${baseURL}/api/`,
  withCredentials: true, // 跨域请求时发送Cookie
  timeout: 5000,
  headers: {},
});


//请求拦截器
service.interceptors.request.use(async (config) => {
  //如果为刷新token的请求则不拦截
  if (config.url === "v1/auth/updateToken") return config;
  if (storage.get('token')) {
    config.headers.Authorization = storage.get('token');
  } else {
    //如果没有accessToken且有refreshTokenoken
    if (storage.get('refreshToken')) {
      const res = await updateTokenAPI(storage.get('refreshToken'));
      if (res.data.code === statusCode.OK) {
        const token = res.data.data.token;
        const refreshToken = res.data.data.refreshToken;

        storage.set("token", token, 5);
        storage.set("refreshToken", refreshToken, 14 * 24 * 60);
        config.headers.Authorization = token;
      }
    }
  }

  return config;
}), (error: any) => {
  return Promise.reject(error);
}

//响应拦截器
service.interceptors.response.use(async (res) => {
  // token 过期
  if (res.data.code === statusCode.TOKEN_EXPRIED) {
    if (storage.get('refreshToken')) {
      const res = await updateTokenAPI(storage.get('refreshToken'));
      if (res.data.code === statusCode.OK) {
        const token = res.data.data.token;
        const refreshToken = res.data.data.refreshToken;

        storage.set("token", token, 5);
        storage.set("refreshToken", refreshToken, 14 * 24 * 60);
        res.config.headers.Authorization = token;
        return service.request(res.config);
      }
    }
  }
  return res;
}, (error) => {
  return Promise.reject(error);
});

export default service;