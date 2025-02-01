// 路由守卫
import router from "./index";
import NProgress from "nprogress";
import "@/assets/styles/nprogress.scss";
import { isEmpty } from "lodash-es";
import { useTheme } from "@/hooks/theme-hooks";
import { globalConfig } from "@/utils/global-config";
import useUserStore from "@/stores/modules/user-store";
import useMenuStore from "@/stores/modules/menu-store";
import useHistoryStore from "@/stores/modules/history-store";

const { theme } = useTheme();

// nprogress配置
NProgress.configure({
  showSpinner: false,
  template: `<div class="bar" style="--base-color: ${theme.value.primaryColor}" role="bar"><div class="peg"></div></div><div class="spinner" role="spinner"><div class="spinner-icon"></div></div>`,
});

/** 白名单 直接跳过的路由路径 */
const whiteList = [`/${globalConfig.baseUrl}/login`, `/${globalConfig.baseUrl}/403`, `/${globalConfig.baseUrl}/404`];

router.beforeEach(async (to, _from, next) => {
  NProgress.start();

  const loginStore = useUserStore();
  /** token */
  const token = loginStore.token;
  /** 是否为白名单 */
  const isWhite = whiteList.some(item => item === to.path);
  if (isWhite) {
    next();
  } else if (token) {     // has token
    // has token, no userInfo
    if (isEmpty(loginStore.userInfo)) {
      await loginStore.getUserInfo();
    }
    // 创建菜单
    const menuStore = useMenuStore();
    if (!menuStore.list.length && !isEmpty(loginStore.userInfo)) {
      await menuStore.createMenuList();
      next({ ...to, replace: true })
    }
    // 历史记录存储
    const historyStore = useHistoryStore();
    if (to.path !== '/') {
      await historyStore.setRouterHistory(to);
    }

    if (to.matched.length) {
      next();
    } else {
      if (router.hasRoute(loginStore.homePage)) {
        next({ name: loginStore.homePage });
      } else {
        next({ name: 'errorPage404' });
      }
    }
  } else { // no token
    const redirect = encodeURIComponent(to.fullPath);
    next({ path: `/${globalConfig.baseUrl}/login`, query: { redirect } });
    window.$message.error("请先登录再进行操作");
  }

  // 修改标题
  document.title = `${globalConfig.title}  ${to.meta.title ? `- ${to.meta.title}` : ""}`;
  // keep-alive include判断
  if (to.meta?.keepAlive) {
    const historyStore = useHistoryStore();
    historyStore.addKeepAliveInclude(to.name as string);
  }
});

router.afterEach(() => {
  NProgress.done();
});
