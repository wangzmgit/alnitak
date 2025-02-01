import axios from "axios";
import type { AxiosInstance } from "axios";
import { updateTokenAPI } from "@/api/auth";
import { statusCode } from "./status-code";
import { globalConfig as config, } from "./global-config";
import { storageData as storage } from "./storage-data";

let requests: any[] = [];
let isRefreshing = false;
export const baseURL = config.domain ? `http${config.https ? 's' : ''}://${config.domain}` : '';

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
  if (localStorage) {
    if (storage.get('token')) {
      config.headers.Authorization = storage.get('token');
    } else {
      //如果没有accessToken且有refreshTokenoken
      const localRefreshToken = storage.get('refreshToken')
      if (localRefreshToken) {
        // 只发送一次刷新token请求
        if (!isRefreshing) {
          isRefreshing = true;
          const tokenRes = await updateTokenAPI(localRefreshToken);
          isRefreshing = false;
          if (tokenRes.data.code === statusCode.OK) {
            const token = tokenRes.data.data.token;
            const refreshToken = tokenRes.data.data.refreshToken;

            storage.set("token", token, 60);
            if (refreshToken && refreshToken !== localRefreshToken) {
              storage.set("refreshToken", refreshToken, 7 * 24 * 60);
            }
            config.headers.Authorization = token;

            //token刷新前的401请求队列重试
            if (requests.length > 0) {
              requests.forEach(cb => cb(token))
              requests = []
            }
            return config;
          }
        } else {
          const token = storage.get('token');
          const reqToken = config.headers.Authorization; // 该请求发起时携带的token
          if (token && token !== reqToken) { //token已经更新就重试请求
            config.headers.Authorization = storage.get('token');
            return service(config)
          } else { // token未更新则放在队列中，更新后再释放
            return new Promise((resolve) => {
              requests.push((t: string) => {
                config.headers.Authorization = t
                resolve(config);
              })
            })
          }
        }
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
  if (localStorage) {
    switch (res.data.code) {
      case statusCode.TOKEN_EXPRIED:
        if (storage.get('refreshToken')) {
          if (!isRefreshing) { // 首次收到需要刷新token的响应
            isRefreshing = true;
            const tokenRes = await updateTokenAPI(storage.get('refreshToken'));
            isRefreshing = false;
            if (tokenRes.data.code === statusCode.OK) {
              const token = tokenRes.data.data.token;
              const refreshToken = tokenRes.data.data.refreshToken;

              storage.set("token", token, 60);
              storage.set("refreshToken", refreshToken, 7 * 24 * 60);
              res.config.headers.Authorization = token;

              //token刷新前的401请求队列重试
              if (requests.length > 0) {
                requests.forEach(cb => cb(token))
                requests = []
              }
              return service.request(res.config);
            }
          } else {
            const token = storage.get('token');
            const reqToken = res.config.headers.Authorization; // 该请求发起时携带的token
            if (token && token !== reqToken) { //token已经更新就重试请求
              res.config.headers.Authorization = storage.get('token');
              return service(res.config)
            } else { // token未更新则放在队列中，更新后再释放
              return new Promise((resolve) => {
                requests.push((t: string) => {
                  res.config.headers.Authorization = t
                  resolve(service(res.config));
                })
              })
            }
          }
        }
        break;
      case statusCode.LOGIN_AGAIN:
        // 清理缓存信息
        storageData.remove("token");
        storageData.remove('refreshToken');
        navigateTo({ name: 'login' })
        break;
    }

  }
  return res;
}, (error) => {
  return Promise.reject(error);
});

export default service;