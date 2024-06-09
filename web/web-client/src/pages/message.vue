<template>
  <div class="msg">
    <header-bar class="header-bar"></header-bar>
    <div class="msg-container">
      <div class="msg-menu-container">
        <nuxt-link v-for="item in menuList" :to="item.to" class="menu-item"
          :class="route.name === item.key ? 'menu-item-active' : ''">
          <span class="menu-title">{{ item.name }}</span>
        </nuxt-link>
      </div>
      <div class="msg-router">
        <NuxtPage />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount, computed } from "vue";
import HeaderBar from "@/components/header-bar/index.vue";

definePageMeta({
  middleware: ['auth', (to) => {
    if (to.name === 'message') {
      return navigateTo("/message/announce");
    }
  }]
})

useHead({
  title: `消息中心 - ${globalConfig.title}`
})

const defaultOption = ref('');//默认激活菜单
const route = useRoute();
const menuList = [
  {
    name: '站内公告',
    key: "message-announce",
    to: '/message/announce',
  },
  {
    name: '收到的赞',
    key: "message-like",
    to: '/message/like',

  },
  {
    name: '回复我的',
    key: "message-reply",
    to: '/message/reply',

  },
  {
    name: '@我的',
    key: "message-at",
    to: '/message/at',

  },
  {
    name: '私信',
    key: "message-whisper",
    to: '/message/whisper',
  }
]

const changeMenu = (name: string) => {
  // router.push({ name: name });
  // defaultOption.value = name;
}

onBeforeMount(() => {
  // defaultOption.value = route.name as string;
})
</script>

<style lang="scss" scoped>
.msg {
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

  .msg-container {
    margin-top: 70px;
    height: calc(100% - 80px);

    .msg-menu-container {
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

    .msg-router {
      width: calc(100% - 280px);
      height: 100%;
      margin-left: 280px;
      box-sizing: border-box;
      padding: 0 60px 0 10px;
      background: #fafafa;
    }
  }
}
</style>