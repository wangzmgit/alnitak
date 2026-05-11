import { useAuthStore } from '@/stores/auth-store';
import { statusCode } from '@/utils/status-code';
import { globalConfig } from '@/utils/global-config';

export default defineNuxtPlugin(async () => {
  if (process.client) return;

  const auth = useAuthStore();
  const headers = useRequestHeaders(['cookie']);
  
  // SSR阶段直接请求后端（端口9001），避免代理问题
  const domain = globalConfig.domain;
  const protocol = globalConfig.https ? 'https' : 'http';
  const url = `${protocol}://${domain}/api/v1/auth/me`;

  try {
    // NODE_TLS_REJECT_UNAUTHORIZED=0 已经全局禁用证书验证
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

