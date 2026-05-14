import { defineStore } from 'pinia';
import Cookies from 'js-cookie';
import { getUserInfoAPI } from '@/api/user';
import { statusCode } from '@/utils/status-code';
import { storageData } from '@/utils/storage-data';
import { getAuthMeAPI, logoutAPI } from '@/api/auth';
import { authDebug } from '@/utils/auth-debug';

export type AuthStatus = 'guest' | 'auth';

const getRedirectUrl = () => {
  if (typeof window === 'undefined') return '/';
  return `${window.location.pathname}${window.location.search}${window.location.hash}`;
};

// 统一凭证写入，避免多处重复
const saveCredentials = (data: { token?: string; refreshToken?: string; userId?: number | string | null }) => {
  if (data.token) {
    storageData.set('token', data.token, 60);
    Cookies.set('token', data.token, { expires: 1 }); // 供 SSR 读取
  }
  if (data.refreshToken) storageData.set('refreshToken', data.refreshToken, 7 * 24 * 60);
  if (data.userId != null) Cookies.set('user_id', String(data.userId));
};

// 清除本地凭证
const clearCredentials = () => {
  storageData.remove('token');
  storageData.remove('refreshToken');
  Cookies.remove('token');
  Cookies.remove('user_id');
};

export { saveCredentials, clearCredentials };

export const useAuthStore = defineStore('auth', {
  state: () => ({
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
        hasToken: Boolean(storageData.get('token')),
        hasRefreshToken: Boolean(storageData.get('refreshToken')),
        hasUserIdCookie: Boolean(Cookies.get('user_id')),
      });

      const hasReadableCred =
        Boolean(storageData.get('token') || storageData.get('refreshToken') || Cookies.get('user_id'));

      this.lastAuthError = '';
      try {
        if (hasReadableCred) {
          const res = await getUserInfoAPI();
          if (res.data.code === statusCode.OK) {
            this.user = res.data.data.userInfo;
            this.status = 'auth';
            authDebug('[auth] fetchMe ok', { status: this.status, uid: this.user?.uid });
            return;
          }
          this.user = null;
          this.status = 'guest';
          this.lastAuthError = res.data.msg || '获取用户信息失败';
          authDebug('[auth] fetchMe guest', { code: res.data.code, msg: res.data.msg });
          return;
        }

        // 无本地可读凭证时仍可能具备 HttpOnly refresh Cookie：与 /auth/me 对齐会话
        const meRes = await getAuthMeAPI();
        if (meRes.data.code === statusCode.OK && meRes.data.data?.userInfo) {
          const d = meRes.data.data;
          saveCredentials({ token: d.token, refreshToken: d.refreshToken, userId: d.userId });
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
        const rt = storageData.get('refreshToken');
        await logoutAPI(rt || undefined);
      } catch {
        // 退出登录失败不阻止本地清理
      }

      clearCredentials();
      this.markGuest();
      authDebug('[auth] logout end', { status: this.status });
    },
  },
});

