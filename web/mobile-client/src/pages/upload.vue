<template>
  <div class="upload">
    <header-bar class="header-bar"></header-bar>
    <div class="upload-container">
      <div class="upload-menu-container">
        <nuxt-link v-for="item in menuList" :to="item.to" class="menu-item" v-show="item.show"
          :class="route.name === item.key ? 'menu-item-active' : ''">
          <span class="menu-icon">
            <component :is="item.icon" size="18" :strokeWidth="3"></component>
          </span>
          <span class="menu-title">{{ item.name }}</span>
        </nuxt-link>
      </div>
      <div class="upload-router">
        <NuxtPage />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { h, ref, onBeforeMount } from "vue";
import HeaderBar from '@/components/header-bar/index.vue';
import VideoIcon from "@/components/icons/VideoIcon.vue";
import CommentIcon from "@/components/icons/CommentIcon.vue";
import UploadIcon from "@/components/icons/UploadIcon.vue";
import { Editor as EditorIcon, Inbox } from '@icon-park/vue-next';

useHead({
  title: `创作中心 - ${globalConfig.title}`
})

definePageMeta({
  middleware: ['auth', (to) => {
    if (to.name === 'upload') {
      return navigateTo(`/upload/video`);
    }
  }],
})

const route = useRoute();
const menuList = [
  {
    key: 'upload-video',
    name: '视频投稿',
    to: '/upload/video',
    show: true,
    icon: UploadIcon,
  },
  {
    key: 'upload-article',
    name: '专栏投稿',
    to: '/upload/article',
    show: globalConfig.article,
    icon: EditorIcon,
  },
  {
    key: 'upload-video-manage',
    name: '视频管理',
    to: '/upload/video-manage',
    show: true,
    icon: VideoIcon,
  },
  {
    key: 'upload-article-manage',
    name: '专栏管理',
    to: '/upload/article-manage',
    show: globalConfig.article,
    icon: Inbox,
  },
  {
    key: 'upload-comment-manage',
    name: '评论管理',
    to: '/upload/comment-manage',
    show: true,
    icon: CommentIcon,
  },
]
</script>

<style lang="scss" scoped>
.upload {
  top: 0;
  bottom: 0;
  width: 100%;
  height: 100%;
  position: fixed;
  overflow-y: scroll;
  background: #f7f7f7;

  .header-bar {
    position: fixed;
    z-index: 999;
  }

  .upload-container {
    margin-top: 70px;
    height: calc(100% - 80px);

    .upload-menu-container {
      position: fixed;
      height: calc(100% - 80px);
      width: 220px;
      margin-left: 60px;
      background-color: #fff;

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

    .upload-router {
      width: calc(100% - 280px);
      height: 100%;
      background: #fff;
      margin-left: 280px;
      box-sizing: border-box;
      padding: 0 60px 0 20px;
      background: #fafafa;
    }
  }
}
</style>