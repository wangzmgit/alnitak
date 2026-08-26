import { useAuthStore } from '@/stores/auth-store';
import { authDebug } from '@/utils/auth-debug';

const AUTH_SYNC_KEY = 'auth_sync';

export default defineNuxtPlugin(() => {
  if (process.server) return;

  const auth = useAuthStore();

  const nuxtApp = useNuxtApp();
  nuxtApp.hook('app:mounted', () => {
    // 页面加载时通过 HttpOnly Cookie 尝试续签登录态
    auth.fetchMe();
  });

  // 跨标签页登录态同步：监听其他标签页的 localStorage 变更
  window.addEventListener('storage', (e: StorageEvent) => {
    if (e.key === AUTH_SYNC_KEY && e.newValue !== e.oldValue) {
      authDebug('[auth] cross-tab auth sync detected, re-fetching');
      auth.fetchMe(true);
    }
  });
});
