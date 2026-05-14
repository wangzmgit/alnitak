<template>
  <div class="header-bar">
    <div class="sidebar-title">
      <el-icon class="menu-icon-btn" size="22" @click="foldClick">
        <hamburger-button></hamburger-button>
      </el-icon>
      <nuxt-link class="title" to="/">{{ globalConfig.title }}</nuxt-link>
    </div>
    <div class="header-search">
      <div class="search-input">
        <input class="input" v-model="keywords" @keydown.enter="handelSearch" name="search">
        <button class="btn" @click="handelSearch">
          <search-icon class="icon" size="16" />
        </button>
      </div>
    </div>
    <div class="header-right">
      <!-- 用户头像 -->
      <div v-if="isLoggedIn" class="avatar-box">
        <nuxt-link to="/space">
          <common-avatar :url="userInfo?.avatar" :size="40" :iconSize="22"></common-avatar>
        </nuxt-link>
        <div class="menu-container">
          <div class="transition"></div>
          <div class="header-menu">
            <div class="menu-info">
              <div class="name-box">
                <span class="name">{{ userInfo?.name }}</span>
                <span class="sign">{{ userInfo?.sign }}</span>
              </div>
            </div>
            <div class="divider disabled-select"></div>
            <!-- 主题入口（与通用 Header 保持一致，带二级菜单） -->
            <div class="menu-item disabled-select theme">
              <div class="link-title">
                <theme-icon class="icon"></theme-icon>
                <span>主题：</span>
                <span class="theme-label">{{ themeLabel }}</span>
              </div>
              <right-icon class="right-icon"></right-icon>
              <!-- 二级悬浮选项 -->
              <div class="submenu">
                <div class="submenu-item" @click.stop.prevent="setTheme('light')">
                  <sun-icon class="icon" />
                  <span>浅色</span>
                </div>
                <div class="submenu-item" @click.stop.prevent="setTheme('dark')">
                  <moon-icon class="icon" />
                  <span>深色</span>
                </div>
              </div>
            </div>
            <nuxt-link class="menu-item disabled-select" to="/space">
              <div class="link-title">
                <user-icon class="icon"></user-icon>
                <span>个人中心</span>
              </div>
              <right-icon class="right-icon"></right-icon>
            </nuxt-link>
            <div class="menu-item disabled-select" @click="logout">
              <div class="link-title">
                <logout-icon class="icon"></logout-icon>
                <span>退出登录</span>
              </div>
              <right-icon class="right-icon"></right-icon>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="avatar-box">
        <div class="login-btn" @click="auth.openLoginModal()">登录</div>
        <div class="menu-container">
          <div class="transition"></div>
          <div class="header-menu">
            <div class="menu-item disabled-select theme">
              <div class="link-title">
                <theme-icon class="icon"></theme-icon>
                <span>主题：</span>
                <span class="theme-label">{{ themeLabel }}</span>
              </div>
              <right-icon class="right-icon"></right-icon>
              <div class="submenu">
                <div class="submenu-item" @click.stop.prevent="setTheme('light')">
                  <sun-icon class="icon" />
                  <span>浅色</span>
                </div>
                <div class="submenu-item" @click.stop.prevent="setTheme('dark')">
                  <moon-icon class="icon" />
                  <span>深色</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <!-- 图形按钮 -->
      <nuxt-link class="icon-btn" to="/message/announce">
        <message-icon class="icon"></message-icon>
        <div class="icon-text">消息</div>
      </nuxt-link>
      <nuxt-link class="icon-btn" to="/space/history">
        <history-icon class="icon"></history-icon>
        <div class="icon-text">历史</div>
      </nuxt-link>
      <nuxt-link class="icon-btn" to="/space/collection">
        <collect-icon class="icon"></collect-icon>
        <div class="icon-text">收藏</div>
      </nuxt-link>
      <!-- 投稿按钮 -->
      <nuxt-link class="upload-btn disabled-select" to="/upload/video">
        <client-only>
          <upload-icon class="upload-icon"></upload-icon>
          <template #fallback>
            <span class="upload-icon"></span>
          </template>
        </client-only>
        <span>投稿</span>
      </nuxt-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue';
import {
  HamburgerButton, Upload as UploadIcon, Search as SearchIcon,
  Message as MessageIcon, History as HistoryIcon,
  FolderFocusOne as CollectIcon,
  User as UserIcon, Logout as LogoutIcon, Right as RightIcon,
  Theme as ThemeIcon, SunOne as SunIcon, Moon as MoonIcon
} from '@icon-park/vue-next';
import { useAuthStore } from '@/stores/auth-store';

const emits = defineEmits(["changeFold"]);

const auth = useAuthStore();
const isLoggedIn = computed(() => auth.isLoggedIn);
const userInfo = computed(() => auth.user);

// 左侧菜单（用 cookie 持久化，SSR 阶段即可读取，避免首屏闪烁）
const menuFoldCookie = useCookie<boolean>('menu-fold-state', {
  default: () => false,
  maxAge: 60 * 60 * 24 * 365,
  path: '/',
});
const menuFold = ref(menuFoldCookie.value);

const foldClick = () => {
  menuFold.value = !menuFold.value;
  menuFoldCookie.value = menuFold.value;
  emits("changeFold", menuFold.value);
}

// 退出登录
const logout = async () => {
  await auth.logout();
}

//搜索功能
const keywords = ref("");
const handelSearch = () => {
  if (!keywords.value) {
    ElMessage.warning("请先输入搜索内容");
    return;
  }

  navigateTo(`/search/${keywords.value}`, {
    open: { target: '_blank' }
  })
}

// 主题切换（与通用 Header 保持一致，使用相同存储键）
type ThemeMode = 'light' | 'dark';
const THEME_KEY = 'ui-theme-mode';
const themeLabel = ref('浅色');

const applyTheme = (mode: ThemeMode) => {
  if (typeof window === 'undefined') return;
  const root = document.documentElement;
  root.setAttribute('data-theme', mode);
  if (mode === 'dark') root.classList.add('dark');
  else root.classList.remove('dark');
};

const setTheme = (mode: ThemeMode) => {
  themeLabel.value = mode === 'light' ? '浅色' : '深色';
  applyTheme(mode);
  try { localStorage.setItem(THEME_KEY, mode); } catch {}
};

// 初始化时通知父组件当前折叠状态
if (menuFold.value) {
  emits("changeFold", true);
}

onMounted(() => {
  // 读取主题设置
  let mode: ThemeMode = 'light';
  try {
    const saved = localStorage.getItem(THEME_KEY) as ThemeMode | null;
    if (saved === 'dark' || saved === 'light') mode = saved;
  } catch {}
  themeLabel.value = mode === 'light' ? '浅色' : '深色';
  applyTheme(mode);
});
</script>

<style lang="scss" scoped>
.header-bar {
  display: flex;
  align-items: center;
  width: 54;
  height: 60px;
  display: flex;
  align-items: center;
  box-shadow: inset 0 -1px var(--border-color);

  .sidebar-title {
    display: flex;
    align-items: center;
    height: 100%;
    width: 220px;

    .menu-icon-btn {
      margin: 0 6px;
      padding: 6px;
      border-radius: 50%;
      cursor: pointer;

      span {
        height: 22px;
      }

      &:hover { background-color: var(--hover-bg); }
    }

    .title {
      color: var(--font-primary-1);
      overflow: hidden;
      white-space: nowrap;
    }
  }


  .header-search {
    flex: 1;

    .search-input {
      width: 500px;
      position: relative;
      margin: 0 auto;

      .input {
        border: 1px solid var(--border-color);
        outline: none;
        padding: 8px 30px 8px 11px;
        height: 36px;
        font-size: 12px;
        line-height: 14px;
        border-radius: 18px;
        width: 500px;
        vertical-align: top;
        color: var(--font-primary-1);
        box-sizing: border-box;
      }

      .btn {
        position: absolute;
        top: 0;
        right: 10px;
        border: none;
        width: 20px;
        height: 36px;
        line-height: 36px;
        font-size: 14px;
        vertical-align: top;
        background: transparent;
        padding: 0;
        cursor: pointer;

        .icon {
          display: block;
          margin-top: 3px;
          width: 20px;
          height: 36px;
          color: var(--font-primary-5);
        }
      }
    }
  }

  .header-right {
    width: 340px;
    display: flex;
    align-items: center;

    .avatar-box {
      margin-right: 10px;
    }

    .login-btn {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      color: var(--primary-text-color);
      text-align: center;
      line-height: 40px;
      font-size: 14px;
      font-weight: 500;
      background-color: var(--primary-hover-color);
    }

    .upload-btn {
      color: var(--primary-text-color);
      display: flex;
      align-items: center;
      justify-content: center;
      background-color: var(--primary-color);
      margin-left: 10px;
      width: 100px;
      height: 36px;
      border-radius: 5px;
      text-align: center;
      cursor: pointer;
      text-decoration: none;
      transition: background-color .3s;

      .upload-icon {
        width: 16px;
        height: 16px;
        margin-right: 5px;
      }

      &:hover {
        background-color: var(--primary-hover-color);
      }
    }
  }
}

.avatar-box {
  position: relative;
  cursor: pointer;
  margin: 0 10px;

  &:hover {
    .menu-container {
      display: block;
    }
  }

  .menu-container {
    display: none;
    position: absolute;
    width: 230px;
    top: 40px;
    left: -110px;
    z-index: 999;

    .transition {
      height: 10px;
    }

    .header-menu {
      box-sizing: border-box;
      padding: 12px 12px 6px;
      background-color: var(--bg-elev-1);
      border-radius: 10px;
      animation: menu .3s ease-in;
      box-shadow: 0 0 30px var(--shadow-weak);

      .menu-info {
        display: flex;

        .name-box {
          width: calc(100% - 52px);
          padding-left: 12px;
          display: flex;
          flex-direction: column;
          justify-content: space-between;

          .name {
            font-size: 14px;
            color: var(--font-primary-2);
          }

          .sign {
            font-size: 12px;
            color: var(--font-primary-3);
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
          }
        }
      }

      .divider {
        height: 1px;
        width: 100%;
        margin: 12px 0 6px;
        background-color: var(--border-color);
      }

      .menu-item {
        width: 100%;
        height: 36px;
        color: var(--font-primary-2);
        text-decoration: none;
        box-sizing: border-box;
        padding: 0 6px;
        border-radius: 3px;
        display: flex;
        align-items: center;
        justify-content: space-between;

        .link-title {
          display: flex;
          align-items: center;

          .icon {
            width: 18px;
            height: 18px;
            margin-right: 10px;
          }
        }

        .right-icon {
          width: 18px;
          height: 18px;
        }

      &:hover { background-color: var(--hover-bg); }
        &:hover { background-color: var(--hover-bg); }
      }

      /* 主题二级菜单 */
      .menu-item.theme { position: relative; }
      .menu-item.theme:hover .submenu { display: block; }
      .submenu {
        display: none;
        position: absolute;
        top: 0;
        right: -160px;
        width: 150px;
        background: var(--bg-elev-1);
        border-radius: 8px;
        box-shadow: 0 0 20px var(--shadow-weak);
        padding: 6px;
      }
      .submenu-item {
        display: flex;
        align-items: center;
        height: 32px;
        border-radius: 4px;
        padding: 0 8px;
        color: var(--font-primary-2);
        cursor: pointer;
      }
      .submenu-item .icon { width: 18px; height: 18px; margin-right: 8px; }
      .submenu-item:hover { background: var(--hover-bg); }
    }
  }
}

.icon-btn {
  color: var(--font-primary-3);
  position: relative;
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  flex-direction: column;
  justify-content: center;
  text-decoration: none;

  .icon {
    width: 18px;
    height: 18px;
    margin-bottom: 2px;
  }

  .icon-text {
    font-size: 12px;
    line-height: 14px;
  }
}


@keyframes menu {
  0% {
    opacity: 0;
    transform: translateY(-10px);
  }

  100% {
    opacity: 1;
  }
}
</style>