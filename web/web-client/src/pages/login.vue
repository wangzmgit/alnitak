<template>
  <div class="login-bg">
    <login-card @success="loginSuccess"></login-card>
  </div>
</template>

<script setup lang="ts" >
import { globalConfig } from '@/utils/global-config';
import { useAuthStore } from '@/stores/auth-store';

const route = useRoute();
const auth = useAuthStore();

useHead({
  title: `${globalConfig.title} - 登录`
})

const loginSuccess = async () => {
  // 登录页属于"独立入口"，这里需要显式同步全局登录态，避免其它页面仍显示未登录。
  await auth.fetchMe();
  const redirect = route.query.redirect?.toString();
  // 防开放重定向：只允许站内路径（以 / 开头且非 // 协议相对路径）
  if (redirect && /^\/(?!\/)/.test(redirect)) {
    await navigateTo(redirect);
  } else {
    await navigateTo('/');
  }
}
</script>

<style lang="scss" scoped>
.login-bg {
  width: 100%;
  height: 100vh;
  background-image: url(@/assets/images/login/login-bg.png);
  background-size: 100% 100%;
}
</style>