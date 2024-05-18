<template>
  <div class="login-bg">
    <login-card @success="loginSuccess"></login-card>
  </div>
</template>

<script setup lang="ts" >
import { globalConfig } from '@/utils/global-config';
import LoginCard from '@/components/login-card/index.vue';

const route = useRoute();

useHead({
  title: `${globalConfig.title} - 登录`
})

const loginSuccess = () => {
  const redirect = route.query.redirect?.toString();
  if (redirect) {
    navigateTo(redirect);
  } else {
    navigateTo('/');
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