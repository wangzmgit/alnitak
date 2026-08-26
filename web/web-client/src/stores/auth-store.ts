import { defineStore } from 'pinia';
import Cookies from 'js-cookie';
import { getUserInfoAPI } from '@/api/user';
import { statusCode } from '@/utils/status-code';
import { getAuthMeAPI, logoutAPI } from '@/api/auth';
import { authDebug } from '@/utils/auth-debug';

export type AuthStatus = 'guest' | 'auth';

const getRedirectUrl = () => {
  if (typeof window === 'undefined') return '/';
  return `${window.location.pathname}${window.location.search}${window.location.hash}`;
};

const AUTH_SYNC_KEY = 'auth_sync';

// 向同源其他标签页广播登录态变更
// 注意：只从用户主动登录/退出的地方调用（LoginForm/logout），
// 不要在 clearCredentials/saveCredentials 里广播，否则 LOGIN_AGAIN 拦截器
// 会触发跨标签无限循环：Tab A 广播 → Tab B fetchMe → LOGIN_AGAIN → 广播 → Tab A ...
export const broadcastAuthChange = () => {
  if (typeof window !== 'undefined') {
    localStorage.setItem(AUTH_SYNC_KEY, Date.now().toString());
  }
};

// 浏览器端可读的 store token 引用，供 request.ts 等无法直接 useAuthStore 的场景使用。
// 只在浏览器端 set，SSR 期间保持空。
let _browserToken = '';

/** 获取当前内存中的 accessToken（浏览器端），SSR 返回空 */
export const getBrowserToken = (): string => _browserToken;

// 统一凭证写入（仅存内存 + user_id Cookie，不再使用 localStorage 存 token）
const saveCredentials = (data: { token?: string; refreshToken?: string; userId?: number | string | null }) => {
  if (data.token) {
    _browserToken = data.token;
    // 同步更新 Pinia store（若已初始化）
    try {
      const auth = useAuthStore();
      auth.token = data.token;
    } catch {
      // Pinia 尚未初始化（SSR 场景），忽略
    }
  }
  if (data.userId != null) Cookies.set('user_id', String(data.userId));
};

// 清除本地凭证
const clearCredentials = () => {
  _browserToken = '';
  Cookies.remove('user_id');
};

export { saveCredentials, clearCredentials };

export const useAuthStore = defineStore('auth', {
  state: () => ({
    // 内存中的 accessToken（刷新页面后丢失，由 fetchMe 通过 HttpOnly Cookie 重新获取）
    token: '' as string,
    // 默认以游客态渲染，确保 SSR 与首屏 hydration DOM 结构一致，避免 hydration mismatch。
    // 若 SSR 已严格校验登录态，将通过 initFromSSR() 覆盖该默认值。
    status: 'guest' as AuthStatus,
    user: null as UserInfoType | null,
    loginModalOpen: false,
    redirectAfterLogin: '' as string,
    lastAuthError: '' as string,
    _fetchPromise: null as Promise<void> | null,
  }),
  getters: {
    isLoggedIn: (s) => s.status === 'auth',
    // 有有效内存 token 或已登录状态都可视为有凭证
    hasToken: (s) => Boolean(s.token),
  },
  actions: {
    initFromSSR(payload?: { status: AuthStatus; user?: UserInfoType | null }) {
      if (!payload) return;
      this.status = payload.status;
      this.user = payload.user ?? null;
    },

    // force=true 时取消复用进行中的请求，强制发起新请求（跨标签同步场景）
    async fetchMe(force = false) {
      if (this._fetchPromise && !force) return this._fetchPromise;
      this._fetchPromise = this._doFetchMe();
      try {
        await this._fetchPromise;
      } finally {
        this._fetchPromise = null;
      }
    },

    async _doFetchMe() {
      authDebug('[auth] fetchMe start', {
        prevStatus: this.status,
        hasToken: Boolean(this.token),
        hasUserIdCookie: Boolean(Cookies.get('user_id')),
      });

      this.lastAuthError = '';
      try {
        // 优先尝试用已有 accessToken 获取用户信息
        if (this.token) {
          const res = await getUserInfoAPI();
          if (res.data.code === statusCode.OK) {
            this.user = res.data.data.userInfo;
            this.status = 'auth';
            authDebug('[auth] fetchMe ok', { status: this.status, uid: this.user?.uid });
            return;
          }
          // token 失效，清除本地状态，继续尝试 Cookie 方式
          this.token = '';
        }

        // 无有效 accessToken 时，尝试通过 HttpOnly refresh Cookie 续签（/auth/me 支持 Cookie 鉴权）
        const meRes = await getAuthMeAPI();
        if (meRes.data.code === statusCode.OK && meRes.data.data?.userInfo) {
          const d = meRes.data.data;
          // 后端返回了新 token，写入内存
          if (d.token) {
            this.token = d.token;
            _browserToken = d.token;
          }
          if (d.userId != null) Cookies.set('user_id', String(d.userId));
          this.user = d.userInfo;
          this.status = 'auth';
          authDebug('[auth] fetchMe ok (cookie me)', { status: this.status, uid: this.user?.uid });
          return;
        }
        this.user = null;
        this.status = 'guest';
        this.lastAuthError = meRes.data.msg || '';
        authDebug('[auth] fetchMe guest (me)', { code: meRes.data.code, msg: meRes.data.msg });
      } catch (e: any) {
        this.user = null;
        this.status = 'guest';
        this.lastAuthError = e?.message || '网络错误';
        authDebug('[auth] fetchMe error', { message: this.lastAuthError });
      } finally {
        authDebug('[auth] fetchMe end', { status: this.status });
      }
    },

    // 由外部（login callback、refreshToken 成功等）通知 store 更新 token
    setToken(token: string) {
      this.token = token;
      _browserToken = token;
    },

    openLoginModal(opts?: { redirect?: string; reason?: string }) {
      const redirect = opts?.redirect || getRedirectUrl();
      this.redirectAfterLogin = redirect;
      this.loginModalOpen = true;
      if (opts?.reason) this.lastAuthError = opts.reason;
      authDebug('[auth] openLoginModal', { redirect, reason: opts?.reason });
    },

    closeLoginModal() {
      this.loginModalOpen = false;
    },

    markGuest() {
      this.status = 'guest';
      this.user = null;
    },

    async logout() {
      authDebug('[auth] logout start');
      try {
        // 后端从 HttpOnly Cookie 读取 refreshToken 进行吊销
        await logoutAPI();
      } catch {
        // 退出登录失败不阻止本地清理
      }

      clearCredentials();
      this.token = '';
      this.markGuest();
      broadcastAuthChange();
      authDebug('[auth] logout end', { status: this.status });
    },
  },
});
