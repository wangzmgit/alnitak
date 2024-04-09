<template>
  <div class="setting">
    <base-tabs class="tabs" :current="route.name?.toString()" :tabs="tabsOptions" @tab-change="tabChange" />
    <NuxtPage class="router" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import BaseTabs from "@/components/base-tabs/index.vue";

definePageMeta({
  middleware: ['auth', (to) => {
    if (to.name === 'space-setting') {
      return navigateTo("/space/setting/info");
    }
  }]
})

const route = useRoute();
const tabsOptions = [
  { key: 'space-setting-info', label: '基本信息' },
  { key: 'space-setting-security', label: '账号安全' }
];

const tabChange = (key: string) => {
  navigateTo({ name: key });
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