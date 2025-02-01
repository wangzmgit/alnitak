<template>
  <n-layout-sider collapse-mode="width" :collapsed-width="collapsedWidth" :collapsed="collapsed"
    :native-scrollbar="false" class="layout-sider" bordered show-trigger="bar" @collapse="collapsed = true"
    @expand="collapsed = false" @update:collapsed="(collapsed: boolean) => emit('collapsed', collapsed)">
    <n-space :wrap="false" :wrap-item="false" class="logo-container" align="center" justify="center" title="返回首页"
      @click="handleJumpHome">
      <n-image :width="collapsedIconSize" src="/logo.png" preview-disabled />
      <n-text v-if="!collapsed" class="text-nowrap"> {{ globalConfig.title }} </n-text>
    </n-space>
    <n-divider class="divider-style" />
    <n-menu :options="menuOptions" :value="currentMenu" :collapsed="collapsed" :collapsed-width="collapsedWidth"
      :collapsed-icon-size="collapsedIconSize" />
  </n-layout-sider>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue";
import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";
import { globalConfig } from "@/utils/global-config";
import useMenuStore from "@/stores/modules/menu-store.js";
import { NLayoutSider, NImage, NSpace, NText, NDivider, NMenu } from "naive-ui";

const emit = defineEmits<{
  (event: "collapsed", collapsed: boolean): void;
}>();

const router = useRouter();
const menuStore = useMenuStore();

const { list: menuOptions } = storeToRefs(menuStore);

// 部分配置
const collapsedWidth = 64;
const collapsedIconSize = 20;

/** 当前路由对应的菜单name */
const currentMenu = computed(() => {
  return String(router.currentRoute.value.name);
});

const collapsed = ref(false);

const handleJumpHome = () => {
  router.push("/");
};
</script>

<style lang="scss" scoped>
.logo-container {
  height: v-bind('(collapsedIconSize + 10) + "px"');
  padding: 10px 0;
  cursor: pointer;
}

.divider-style {
  margin: 0;
  padding: 0 10px;
}

.text-nowrap {
  font-size: 16px;
  font-weight: 500;
}
</style>
