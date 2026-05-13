<template>
  <div class="header-bar">
    <nuxt-link class="header-left" to="/">
      <h1 class="title">{{ globalConfig.title }}</h1>
    </nuxt-link>
    <div v-show="isSearchPage" class="header-center">
      <div class="search-form">
        <input class="input" v-model="keywords" @keydown.enter="handelSearch" id="header-search-input" name="keywords">
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
            <!-- 主题入口（基础占位，不含逻辑） -->
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
            <!-- 仅展示主题切换（未登录也可用） -->
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
import { ref, onMounted, computed } from 'vue';
import {
  Search as SearchIcon, Message as MessageIcon,
  Upload as UploadIcon, History as HistoryIcon,
  User as UserIcon, Logout as LogoutIcon,
  Right as RightIcon, FolderFocusOne as CollectIcon,
  Theme as ThemeIcon, SunOne as SunIcon, Moon as MoonIcon
} from '@icon-park/vue-next';
import { useAuthStore } from '@/stores/auth-store';

const route = useRoute();
const isSearchPage = ref(route.name !== 'search-keywords');

const keywords = ref('');

type ThemeMode = 'light' | 'dark';
const THEME_KEY = 'ui-theme-mode';
const themeLabel = ref('浅色');

const applyTheme = (mode: ThemeMode) => {
  if (typeof window === 'undefined') return;
  const root = document.documentElement;
  // data-theme 提供给自定义样式变量使用
  root.setAttribute('data-theme', mode);
  // Element Plus 深色模式通过 html.dark 触发
  if (mode === 'dark') {
    root.classList.add('dark');
  } else {
    root.classList.remove('dark');
  }
};

const setTheme = (mode: ThemeMode) => {
  themeLabel.value = mode === 'light' ? '浅色' : '深色';
  applyTheme(mode);
  try { localStorage.setItem(THEME_KEY, mode); } catch {}
};

onMounted(() => {
  let mode: ThemeMode = 'light';
  try {
    const saved = localStorage.getItem(THEME_KEY) as ThemeMode | null;
    if (saved === 'dark' || saved === 'light') mode = saved;
  } catch {}
  themeLabel.value = mode === 'light' ? '浅色' : '深色';
  applyTheme(mode);
});
const handelSearch = () => {
  if (!keywords.value) {
    ElMessage.warning("请先输入搜索内容");
    return;
  }

  navigateTo(`/search/${keywords.value}`, {
    open: { target: '_blank' }
  })
}

const auth = useAuthStore();
const isLoggedIn = computed(() => auth.isLoggedIn);
const userInfo = computed(() => auth.user);

const logout = async () => {
  await auth.logout();
}

</script>

<style lang="scss" scoped>
.header-bar {
  width: 100%;
  height: 60px;
  padding: 0 36px;
  z-index: 999;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: space-between;
  /* HeaderBar 组件：顶栏容器背景与阴影（跟随主题） */
  background-color: var(--bg-elev-1);
  -webkit-box-shadow: 0px 0px 3px var(--shadow-weak);
  -moz-box-shadow: 0px 0px 3px var(--shadow-weak);
  box-shadow: 0px 0px 3px var(--shadow-weak);

  .header-left {
    display: flex;
    align-items: center;
    width: 260px;

    .title {
      color: var(--font-primary-1);
      margin: 0;
      font-size: 16px;
      font-weight: 500;
      vertical-align: top;
      cursor: pointer;
      overflow: hidden;
      white-space: nowrap;
    }
  }

  .header-center {
    width: 300px;

    .search-form {
      position: relative;

      .input {
        /* HeaderBar 组件：搜索输入边框颜色（跟随主题） */
        border: 1px solid var(--border-color);
        outline: none;
        padding: 8px 30px 8px 11px;
        height: 36px;
        font-size: 12px;
        line-height: 14px;
        border-radius: 18px;
        width: 300px;
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
    width: 320px;
    display: flex;
    align-items: center;

    .avatar-box {
      margin-right: 10px;
    }

    .login-btn {
      display: block;
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
    left: -95px;
    z-index: 999;

    .transition {
      height: 10px;
    }

      /* HeaderBar 组件：头像下拉菜单容器（背景/阴影跟随主题） */
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

        /* HeaderBar 组件：下拉菜单分割线颜色（跟随主题） */
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

        /* HeaderBar 组件：菜单项悬浮态背景（跟随主题） */
        &:hover {
          background-color: var(--hover-bg);
        }
      }

        /* 主题二级菜单 */
        .menu-item.theme {
          position: relative;
        }
        .menu-item.theme:hover .submenu {
          display: block;
        }
        /* HeaderBar 组件：二级菜单容器（背景/阴影跟随主题） */
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
        /* HeaderBar 组件：二级菜单项悬浮态（跟随主题） */
        .submenu-item:hover { background: var(--hover-bg); }
    }
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