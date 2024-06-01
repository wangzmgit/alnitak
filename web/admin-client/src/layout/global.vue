<template>
  <router-view />
</template>

<script lang="ts" setup>
import { useDialog, useMessage } from "naive-ui";

// 挂载window下方法组件
window.$message = useMessage();
window.$dialog = useDialog();
</script>

<style lang="scss" scoped>

</style>
