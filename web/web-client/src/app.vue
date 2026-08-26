<template>
<el-config-provider :locale="zhCn">
  <NuxtPage />
  <client-only>
    <login-dialog
      v-if="auth.loginModalOpen"
      @close="auth.closeLoginModal"
      @success="handleLoginSuccess"
    />
  </client-only>
</el-config-provider>
</template>

<script setup lang="ts">
import { ElConfigProvider } from 'element-plus';
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import { useAuthStore } from '@/stores/auth-store';
import { initNetworkLine } from '@/utils/network-line';

const auth = useAuthStore();

// 应用启动时检测网络线路，为整个会话选择最优 OSS
initNetworkLine();

const handleLoginSuccess = async () => {
  auth.closeLoginModal();
  await auth.fetchMe();

  const redirect = auth.redirectAfterLogin;
  auth.redirectAfterLogin = '';
  if (redirect) {
    await navigateTo(redirect);
  }
};
</script>

<style lang="scss">
@import url(@/assets/styles/global.scss);

body {
  margin: 0;
  padding: 0;
  background-color: var(--bg-page);
  color: var(--font-primary-1);
  font-family: PingFang SC, HarmonyOS Regular, Helvetica Neue, Microsoft YaHei, sans-serif;
}

@font-face {
  font-family: "HarmonyOS Regular";
  src: url("@/assets/fonts/HarmonyOS_Sans_SC_Regular.ttf");
  font-display: swap;
}

a {
  text-decoration: none;
}
</style>