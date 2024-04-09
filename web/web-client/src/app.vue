<template>
  <el-config-provider :locale="zhCn">
    <NuxtPage />
  </el-config-provider>
</template>

<script setup lang="ts">
import { ElConfigProvider } from 'element-plus';
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
</script>

<style lang="scss">
@import url(@/assets/styles/global.scss);

body {
  margin: 0;
  padding: 0;
  font-family: PingFang SC, HarmonyOS Regular, Helvetica Neue, Microsoft YaHei, sans-serif;
}

@font-face {
  font-family: "HarmonyOS Regular";
  src: url("@/assets/fonts/HarmonyOS_Sans_SC_Regular.ttf");
}

a {
  text-decoration: none;
}
</style>