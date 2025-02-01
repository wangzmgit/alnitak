<template>
  <n-space class="header-container" align="center" justify="space-between">
    <header-breadcrumb />
    <n-space align="center">
      <header-controller />
      <header-user />
    </n-space>
  </n-space>
  <n-divider />
  <div class="bread-container">
    <view-tags />
  </div>
</template>

<script lang="ts" setup>
import { NSpace, NDivider } from "naive-ui";
import HeaderController from "../components/header-controller/index.vue";
import HeaderUser from "../components/header-user/index.vue";
import HeaderBreadcrumb from "../components/header-breadcrumb/index.vue";
import ViewTags from "../components/view-tags/index.vue";
</script>

<style lang="scss" scoped>
.header-container {
  padding: 0 15px;
  height: 50px;
}

.bread-container {
  height: 40px;
  padding: 0 10px;
}

.n-divider:not(.n-divider--vertical) {
  margin: 0;
}
</style>
