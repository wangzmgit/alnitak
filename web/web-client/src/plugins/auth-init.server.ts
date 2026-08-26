import { useAuthStore } from '@/stores/auth-store';
import { statusCode } from '@/utils/status-code';
import { globalConfig } from '@/utils/global-config';

export default defineNuxtPlugin(async () => {
  if (process.client) return;

  const auth = useAuthStore();
  const headers: Record<string, string> = {};

  // 转发 Cookie（含 HttpOnly refresh_token），后端 /auth/me 支持 Cookie 鉴权
  const cookie = useRequestHeaders(['cookie']).cookie;
  if (cookie) {
    headers.cookie = cookie;
  }

  // SSR阶段直接请求后端（端口9001），避免代理问题
  const domain = globalConfig.domain;
  const protocol = globalConfig.https ? 'https' : 'http';
  const url = `${protocol}://${domain}/api/v1/auth/me`;

  try {
    const res: any = await $fetch(url, { headers });
    if (res?.code === statusCode.OK && res?.data?.userInfo) {
      auth.initFromSSR({ status: 'auth', user: res.data.userInfo });
      return;
    }
  } catch {
    // SSR 初始化失败不阻塞页面渲染，回退为游客态
  }

  auth.initFromSSR({ status: 'guest', user: null });
});
