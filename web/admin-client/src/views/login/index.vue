<template>
  <div class="login">
    <div class="login-left">
      <img class="illustration" src="@/assets/images/login-illustration.png" />
    </div>
    <div class="login-right">
      <div class="login-form">
        <span class="title">{{ globalConfig.title }}后台管理系统</span>
        <login-form @success="loginSuccess"></login-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import LoginForm from './component/LoginForm.vue';
import useUserStore from '@/stores/modules/user-store';
import { useTheme } from '@/hooks/theme-hooks';
import { globalConfig } from '@/utils/global-config';

const { theme } = useTheme();
const { login } = useUserStore();

const loginSuccess = () => {
  login();
}
</script>

<style lang="scss" scoped>
.login {
  display: flex;
  width: 100%;
  height: 100vh;

  .login-left {
    width: 50%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: v-bind('theme.primaryColor');

    .illustration {
      width: 80%;
    }
  }

  .login-right {
    width: 50%;
    height: 100%;
  }
}

.login-form {
  width: 60%;
  margin: 220px auto 0;

  .title {
    text-align: center;
    display: block;
    width: 100%;
    margin: 20px 0;
    color: var(--title-color);
    font-size: 28px;
  }
}
</style>