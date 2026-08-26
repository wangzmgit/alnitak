import axios from "axios";
import type { AxiosInstance, AxiosError } from "axios";
import { updateTokenAPI } from "@/api/auth";
import { statusCode } from "./status-code";
import { globalConfig as config, } from "./global-config";
import { useAuthStore, getBrowserToken, saveCredentials, clearCredentials } from "@/stores/auth-store";

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

// 开发环境仅在"浏览器端"通过前端代理，SSR/生产环境直连后端
const isBrowser = typeof window !== 'undefined';
export const baseURL = (process.dev && isBrowser)
  ? ''
  : (config.domain ? `http${config.https ? 's' : ''}://${config.domain}` : '');

const service: AxiosInstance = axios.create({
  baseURL: `${baseURL}/api/`,
  withCredentials: true, // 跨域请求时发送 Cookie（含 HttpOnly refresh_token）
  timeout: 5000,
  headers: {
    'X-Requested-With': 'XMLHttpRequest',
  },
});

// 刷新 token 的统一函数（依赖 HttpOnly Cookie，不从前端存储读取 refreshToken）
const refreshToken = async (): Promise<string> => {
  if (refreshPromise) {
    return refreshPromise;
  }

  refreshPromise = (async () => {
    try {
      // 不传 refreshToken，后端从 HttpOnly Cookie 中读取
      const tokenRes = await updateTokenAPI();
      if (tokenRes.data.code === statusCode.OK) {
        const token = tokenRes.data.data.token;

        // 更新内存中的 token
        saveCredentials({ token, userId: tokenRes.data.data.userId });

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

  // 确保只在客户端处理 Authorization
  if (process.client) {
    const token = getBrowserToken();
    if (token) {
      config.headers.Authorization = token;
    }
    // 无 token 时不设置 Authorization 头，后端通过 Cookie 鉴权或返回 401
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
        // token 过期：通过 HttpOnly Cookie 尝试刷新
        if (!isRefreshing) {
          isRefreshing = true;
          try {
            const token = await refreshToken();
            res.config.headers.Authorization = token;
            return service.request(res.config); // 重新发起请求
          } catch (error) {
            console.warn('[request] Token refresh in response interceptor failed:', error);
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
        // 仅清理本地凭证和切换游客态。
        // 不调用 closeLoginModal()——用户可能正主动打开登录弹窗，
        // 关闭它会打断登录流程。弹窗由用户操作控制，不由后端响应接管。
        clearCredentials();
        try {
          const auth = useAuthStore();
          auth.token = '';
          auth.markGuest();
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

  // 重试耗尽后，给用户一个友好提示
  if (isRetryableError(error)) {
    ElMessage.error('网络连接失败，请检查服务器状态');
  }

  return Promise.reject(error);
});

export default service;
