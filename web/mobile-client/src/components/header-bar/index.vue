<template>
  <div class="header-bar">
    <nuxt-link class="header-left" to="/">
      <h1 class="title">{{ globalConfig.title }}</h1>
    </nuxt-link>
    <div v-show="isSearchPage" class="header-center">
      <div class="search-form">
        <input class="input" v-model="keywords" @keydown.enter="handelSearch">
        <button class="btn" @click="handelSearch">
          <search-icon class="icon" size="16" />
        </button>
      </div>
    </div>
    <div v-if="!loading" class="header-right">
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
        <nuxt-link class="login-btn" to="/login">登录</nuxt-link>
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
        <upload-icon class="upload-icon"></upload-icon>
        <span>投稿</span>
      </nuxt-link>
    </div>
    <div v-else class="header-right"></div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import {
  Search as SearchIcon, Message as MessageIcon,
  Upload as UploadIcon, History as HistoryIcon,
  User as UserIcon, Logout as LogoutIcon,
  Right as RightIcon, FolderFocusOne as CollectIcon
} from '@icon-park/vue-next';
import { logoutAPI } from '@/api/auth';
import { getUserInfoAPI } from '@/api/user';
import CommonAvatar from '@/components/common-avatar/index.vue';


const route = useRoute();
const isSearchPage = ref(route.name !== 'search-keywords');

const keywords = ref('');
const handelSearch = () => {
  if (!keywords.value) {
    ElMessage.warning("请先输入搜索内容");
    return;
  }

  navigateTo(`/search/${keywords.value}`, {
    open: { target: '_blank' }
  })
}

const loading = ref(true);
const isLoggedIn = ref(false);
const userInfo = ref<UserInfoType>()
const getUserInfo = async () => {
  const res = await getUserInfoAPI();
  if (res.data.code === statusCode.OK) {
    userInfo.value = res.data.data.userInfo;
    isLoggedIn.value = true;
  }

  loading.value = false;
}

const logout = async () => {
  await logoutAPI(storageData.get('refreshToken'));

  storageData.remove("token");
  storageData.remove('refreshToken');
  isLoggedIn.value = false;
}

onBeforeMount(() => {
  getUserInfo();
})

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
  background-color: #fff;
  -webkit-box-shadow: 0px 0px 3px #c8c8c8;
  -moz-box-shadow: 0px 0px 3px #c8c8c8;
  box-shadow: 0px 0px 3px #c8c8c8;

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
        border: 1px solid #e5e5e5;
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
    left: -110px;
    z-index: 999;

    .transition {
      height: 10px;
    }

    .header-menu {
      box-sizing: border-box;
      padding: 12px 12px 6px;
      background-color: #fff;
      border-radius: 10px;
      animation: menu .3s ease-in;
      box-shadow: 0 0 30px rgba(0, 0, 0, .1);

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
        background-color: #e5e5e5;
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

        &:hover {
          background-color: #e3e5e7;
        }
      }
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