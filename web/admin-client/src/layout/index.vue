<template>
  <n-layout class="layout-container" has-sider>
    <!-- 侧栏 -->
    <sider-content @collapsed="handleSiderCollapse" />
    <n-layout>
      <!-- 头部 -->
      <n-layout-header ref="headerRef" class="layout-header" bordered>
        <header-content />
      </n-layout-header>
      <!-- 主内容 -->
      <n-layout-content ref="contentRef" :native-scrollbar="false" embedded class="layout-content-container">
        <div class="layout-content-main">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in" appear>
              <keep-alive :include="keepAliveInclude">
                <component :is="Component" />
              </keep-alive>
            </transition>
          </router-view>
        </div>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script lang="ts" setup>
import { storeToRefs } from "pinia";
import HeaderContent from "./content/header.vue";
import SiderContent from "./content/sider.vue";
import useResizeHook from "@/hooks/resize-hooks";
import useHistoryStore from "@/stores/modules/history-store.js";
import { NLayout, NLayoutContent, NLayoutHeader } from "naive-ui";

const historyStore = useHistoryStore();
const { keepAliveInclude } = storeToRefs(historyStore);

const { style_content_padding, style_content_height, headerRef, contentRef, changeContentWidth } =
  useResizeHook({ contentPadding: 10 });

const handleSiderCollapse = (_collapsed: boolean) => {
  changeContentWidth();
};
</script>

<style lang="scss" scoped>
.layout-container {
  height: 100vh;
}

.layout-content-container {
  height: 100%;
  height: v-bind(style_content_height);

  :deep(.n-scrollbar-content) {
    height: 100%;
  }

  .layout-content-main {
    height: 100%;
    margin: 0;
    box-sizing: border-box;
    padding: v-bind('style_content_padding + "px"');


  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
