<template>
  <div class="setting">
    <base-tabs class="tabs" :current="route.name?.toString()" :tabs="tabsOptions" @tab-change="tabChange" />
    <NuxtPage class="router" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
definePageMeta({
  middleware: ['auth', (to, from) => {
    // 只在直接访问 /space/setting 时重定向，避免循环重定向
    if (to.name === 'space-setting' && from.name !== 'space-setting-info') {
      return navigateTo({ name: 'space-setting-info' }, { replace: true });
    }
  }]
})

const route = useRoute();
const tabsOptions = [
  { key: 'space-setting-info', label: '基本信息' },
  { key: 'space-setting-security', label: '账号安全' }
];

const router = useRouter();
const isNavigating = ref(false);

const tabChange = async (key: string) => {
  // 防止重复点击导致的问题
  if (isNavigating.value) {
    return;
  }

  // 如果已经在目标页面，不需要跳转
  if (route.name === key) {
    return;
  }

  try {
    isNavigating.value = true;
    // 使用 router.push 而不是 navigateTo，避免页面刷新
    await router.push({ name: key });
  } catch (error) {
    console.error('路由跳转失败:', error);
  } finally {
    isNavigating.value = false;
  }
}
</script>

<style lang="scss" scoped>
.setting {

  .tabs {
    width: 360px;
    margin: 16px auto 16px;
  }
}
</style>