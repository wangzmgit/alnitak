<template>
    <div class="system-config">
  <n-card class="card">
    <n-tabs type="line" default-value="email" animated>
      <n-tab-pane name="email" tab="邮箱">
        <email-config></email-config>
      </n-tab-pane>
      <n-tab-pane name="storage" tab="文件存储">
        <storage-config></storage-config>
      </n-tab-pane>
      <n-tab-pane name="transcoding" tab="转码">
        <transcoding-config></transcoding-config>
      </n-tab-pane>
      <n-tab-pane name="other" tab="其他配置">
        <other-config></other-config>
      </n-tab-pane>
      <n-tab-pane name="cleanup" tab="资源清理">
        <cleanup-config></cleanup-config>
      </n-tab-pane>
    </n-tabs>
  </n-card>
</div>
</template>

<script setup lang="ts">
import EmailConfig from './components/EmailConfig.vue';
import StorageConfig from './components/StorageConfig.vue';
import OtherConfig from './components/OtherConfig.vue';
import CleanupConfig from './components/CleanupConfig.vue';
import TranscodingConfig from './components/TranscodingConfig.vue';

import { NCard, NTabs, NTabPane } from "naive-ui";
</script>

<style lang="scss" scoped>

.system-config {
  height: 100%;

  .card {
    height: 100%;
  }
}

:deep(.n-card__content) {
  padding-top: 6px !important;
}
</style>