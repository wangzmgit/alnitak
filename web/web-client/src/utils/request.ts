import axios from "axios";
import type { AxiosInstance, AxiosError } from "axios";
import { updateTokenAPI } from "@/api/auth";
import { statusCode } from "./status-code";
import { globalConfig as config, } from "./global-config";
import { storageData as storage } from "./storage-data";
import { useAuthStore, saveCredentials, clearCredentials } from "@/stores/auth-store";

// 重试配置
const MAX_RETRIES = 3;
const RETRY_DELAY = 1000; // 初始延迟 1s
const RETRYABLE_ERRORS = ['ECONNRESET', 'ECONNREFUSED', 'ETIMEDOUT', 'ENOTFOUND', 'ERR_NETWORK'];

// 判断是否为可重试的网络错误
const isRetryableError = (error: AxiosError): boolean => {
  // 网络错误（无响应）
  if (!error.response) {
    const code = error.code || '';
    const message = error.message || '';
    return (
      RETRYABLE_ERRORS.some(e => code.includes(e) || message.includes(e)) ||
      message.includes('Network Error') ||
      message.includes('timeout')
    );
  }
  // 502/503/504 网关错误也可重试
  const status = error.response.status;
  return status === 502 || status === 503 || status === 504;
};

// 延迟函数
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

// Token 刷新队列,使用类型定义
type TokenCallback = (token: string) => void;
let requests: TokenCallback[] = [];
let isRefreshing = false;
let refreshPromise: Promise<string> | null = null;

// 开发环境仅在“浏览器端”通过前端代理，SSR/生产环境直连后端
const isBrowser = typeof window !== 'undefined';
export const baseURL = (process.dev && isBrowser)
  ? ''
  : (config.domain ? `http${config.https ? 's' : ''}://${config.domain}` : '');

const service: AxiosInstance = axios.create({
  baseURL: `${baseURL}/api/`,
  withCredentials: true, // 跨域请求时发送Cookie
  timeout: 5000,
  headers: {},
});

// 刷新 token 的统一函数
const refreshToken = async (): Promise<string> => {
  if (refreshPromise) {
    return refreshPromise;
  }

  refreshPromise = (async () => {
    try {
      const localRefreshToken = storage.get('refreshToken');
      const tokenRes = await updateTokenAPI(localRefreshToken || undefined);
      if (tokenRes.data.code === statusCode.OK) {
        const token = tokenRes.data.data.token;
        const rt = tokenRes.data.data.refreshToken;

        saveCredentials({ token, refreshToken: rt, userId: tokenRes.data.data.userId });

        // 执行队列中的所有等待回调
        requests.forEach(cb => cb(token));
        requests = [];

        return token;
      }
      throw new Error('Token refresh failed');
    } catch (err) {
      // 刷新失败时，reject 队列中所有等待的请求，避免永久挂起
      requests.forEach(cb => cb(''));
      requests = [];
      throw err;
    } finally {
      isRefreshing = false;
      refreshPromise = null;
    }
  })();

  return refreshPromise;
};

// 请求拦截器
service.interceptors.request.use(async (config) => {
  // 如果为刷新 token 的请求则不拦截
  if (config.url === "v1/auth/updateToken") return config;

  // 确保只在客户端处理 token
  if (process.client) {
    const token = storage.get('token');
    if (token) {
      config.headers.Authorization = token;
    } else {
      // 如果没有 accessToken 且有 refreshToken
      const localRefreshToken = storage.get('refreshToken');
      if (localRefreshToken) {
        if (!isRefreshing) {
          isRefreshing = true;
          try {
            // 刷新 token
            const token = await refreshToken();
            config.headers.Authorization = token;
            return config;
          } catch (error) {
            console.error('Token refresh failed:', error);
            // 刷新失败,继续原请求
          }
        } else {
          // 正在刷新中，等待刷新完成
          return new Promise((resolve) => {
            requests.push((token: string) => {
              config.headers.Authorization = token;
              resolve(config);
            });
          });
        }
      }
    }
  }
  return config;
}, (error: any) => {
  return Promise.reject(error);
});

// 响应拦截器
service.interceptors.response.use(async (res) => {
  if (process.client) {
    switch (res.data.code) {
      case statusCode.TOKEN_EXPRIED:
        // token 过期：优先刷新；无本地 refresh 时仍可能通过 HttpOnly Cookie 由服务端续签
        if (!isRefreshing) {
          isRefreshing = true;
          try {
            const token = await refreshToken();
            res.config.headers.Authorization = token;
            return service.request(res.config); // 重新发起请求
          } catch (error) {
            console.error('Token refresh in response interceptor failed:', error);
            return res;
          }
        }
        return new Promise((resolve) => {
          requests.push((token: string) => {
            res.config.headers.Authorization = token;
            resolve(service(res.config));
          });
        });
      case statusCode.LOGIN_AGAIN:
        // 清理缓存信息并切换为游客态。
        // 注意：这里不要自动弹出登录弹窗，否则“跨标签页退出登录”或页面后台轮询时
        // 会在用户无操作的情况下频繁弹窗，影响体验。
        clearCredentials();
        try {
          const auth = useAuthStore();
          auth.markGuest();
          auth.closeLoginModal();
        } catch {
          // 兜底：不阻塞响应流
        }
        break;
    }
  }
  return res;
}, async (error: AxiosError) => {
  const config = error.config;

  // 如果没有 config 或已超过最大重试次数，直接拒绝
  if (!config) {
    return Promise.reject(error);
  }

  // 初始化重试计数
  (config as any).__retryCount = (config as any).__retryCount || 0;

  // 检查是否可重试且未超过最大重试次数
  if (isRetryableError(error) && (config as any).__retryCount < MAX_RETRIES) {
    (config as any).__retryCount += 1;
    const retryCount = (config as any).__retryCount;

    // 指数退避延迟
    const delayMs = RETRY_DELAY * Math.pow(2, retryCount - 1);
    console.warn(`[request] 网络错误，${delayMs}ms 后进行第 ${retryCount} 次重试...`, error.message);

    await delay(delayMs);
    return service.request(config);
  }

  return Promise.reject(error);
});

export default service;