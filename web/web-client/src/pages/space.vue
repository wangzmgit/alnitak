<template>
  <div class="space">
    <header-bar class="header-bar "></header-bar>
    <div class="space-container">
      <div class="space-header">
        <button class="upload-btn">上传封面图片</button>
        <img class="cover" v-if="userInfo?.spacecover" :src="userInfo.spacecover" alt="用户封面图" />
        <div class="header-inner">
          <common-avatar :url="userInfo?.avatar" :size="60" :iconSize="36"></common-avatar>
          <div class="header-info">
            <div class="name">
              <span>{{ userInfo?.name }}</span>
              <span class="gender-icon">
                <male v-if="userInfo?.gender === 1" fill="#1890ff"></male>
                <female v-else-if="userInfo?.gender === 2" fill="#eb2f96"></female>
              </span>
            </div>
            <div class="sign">{{ userInfo?.sign }}</div>
          </div>
          <!--关注粉丝信息-->
          <div class="user-data">
            <div>
              <p class="data-title">稿件</p>
              <p>{{ videoCount }}</p>
            </div>
            <div>
              <p class="data-title">关注</p>
              <p class="data-content">1</p>
            </div>
            <div>
              <p class="data-title">粉丝</p>
              <p class="data-content">1</p>
            </div>
          </div>
        </div>
      </div>
      <div class="space-content">
        <div class="space-menu">
          <nuxt-link v-for="item in menuList" :to="item.to" class="menu-item"
            :class="route.name === item.key ? 'menu-item-active' : ''">
            <span class="menu-icon">
              <component :is="item.icon" size="18"></component>
            </span>
            <span class="menu-title">{{ item.name }}</span>
          </nuxt-link>
        </div>
        <div class="space-router">
          <NuxtPage />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import HeaderBar from '@/components/header-bar/index.vue';
import CommonAvatar from '@/components/common-avatar/index.vue';
import {
  Male, Female, VideoTwo, FolderFocus, Upload as UploadIcon,
  History as HistoryIcon, Message as MessageIcon, Config
} from '@icon-park/vue-next';
import { useVideoCountStore } from "@/composables/video-count-store";
import { getUserBaseInfoAPI, asyncGetUserBaseInfoAPI } from '@/api/user';
import { storeToRefs } from 'pinia';

definePageMeta({
  middleware: ['auth', (to) => {
    if (to.name === 'space') {
      return navigateTo(`/space/manuscript`);
    }
  }]
})

const route = useRoute();
const menuList = [
  {
    key: 'space-manuscript',
    name: '稿件',
    to: '/space/manuscript',
    icon: VideoTwo,
  },
  {
    key: 'space-collection',
    name: '收藏夹',
    to: '/space/collection',
    icon: FolderFocus,
  },
  {
    key: 'space-history',
    name: '历史记录',
    to: '/space/history',
    icon: HistoryIcon,
  },
  {
    key: 'upload-video',
    name: '创作中心',
    to: '/upload/video',
    icon: UploadIcon,
  },
  {
    key: 'space-message',
    name: '消息',
    to: '/message/announce',
    icon: MessageIcon,
  },
  {
    key: 'space-setting',
    name: '设置',
    to: '/space/setting',
    icon: Config,
  },
]

const videoCountStore = useVideoCountStore();
const { videoCount } = storeToRefs(videoCountStore);

const userId = useCookie('user_id');
const userInfo = ref<UserInfoType>();
const { data } = await asyncGetUserBaseInfoAPI(userId.value!);
if ((data.value as any).code === statusCode.OK) {
  userInfo.value = (data.value as any).data.userInfo;
}
useHead({
  title: `${userInfo.value?.name}的个人中心`
})

onBeforeMount(() => {
  console.log('userInfo.value', userInfo.value)
})
</script>

<style lang="scss" scoped>
.space {
  width: 100%;

  .header-bar {
    position: fixed;
    top: 0;
  }

  .space-container {
    width: 1300px;
    margin: 70px auto 0;
  }
}

.space-header {
  position: relative;
  width: 100%;
  height: 220px;
  background-color: #9196a1;

  .upload-btn {
    position: absolute;
    top: 12px;
    right: 20px;
    width: 120px;
    height: 32px;
    color: #d3d3d3;
    background-color: transparent;
    cursor: pointer;
    border-radius: 3px;
    border: 1px solid rgba(255, 255, 255, .2);

    &:hover {
      border-color: #d3d3d3;
    }
  }

  .cover {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .header-inner {
    width: 100%;
    height: 100px;
    position: absolute;
    bottom: 0;
    display: flex;
    align-items: center;
    box-sizing: border-box;
    padding: 0 20px;

    .header-info {
      flex: 1;
      padding-left: 20px;

      .name {
        color: var(--font-primary-6);
        display: inline-block;
        margin-right: 5px;
        font-weight: 600;
        line-height: 18px;
        font-size: 18px;
        vertical-align: middle;
      }

      .gender-icon {
        display: inline-block;
        width: 18px;
        height: 18px;
        margin: 0 6px;
        vertical-align: middle;
      }

      .sign {
        font-size: 12px;
        color: #ccd0d7;
        margin: 6px 0 0 0;
      }
    }

    .user-data {
      width: 236px;
      display: flex;

      div {
        color: #fff;
        width: 78px;
        text-align: center;

        .data-title {
          margin-bottom: 6px;
          font-weight: bold;
        }

        .data-content:hover {
          cursor: pointer;
        }
      }
    }
  }
}

.space-content {
  display: flex;
  margin: 0 auto 30px;

  .space-menu {
    width: 200px;
    background-color: #f8f8f8;

    .menu-item {
      position: relative;
      height: 42px;
      padding-left: 24px;
      cursor: pointer;
      border-radius: 2px;
      margin: 6px 5px 0;
      box-sizing: border-box;
      display: flex;
      align-items: center;
      text-decoration: none;
      color: var(--font-primary-1);

      &:hover {
        background-color: #f3f3f5;
      }

      .menu-icon {
        width: 18px;
        height: 18px;
        margin-right: 8px;
      }
    }

    .menu-item-active {
      color: var(--primary-color) !important;
      background-color: var(--primary-color-active) !important;
    }
  }

  .space-router {
    margin-left: 10px;
    width: calc(100% - 210px);
    min-height: 630px;
    background-color: #fff;
  }
}

@media (max-width: 1400px) {
  .space {
    .space-container {
      width: 1100px;
    }
  }

  .space-header {
    height: 260px;
  }
}
</style>