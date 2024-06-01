<template>
  <n-space align="center" size="small" :wrap-item="false" :wrap="false">
    <theme-switch />
    <fullscreen-btn />
  </n-space>
</template>

<script lang="ts" setup>
import { NSpace } from "naive-ui";
import ThemeSwitch from "./theme-switch.vue";
import FullscreenBtn from "./fullscreen-btn.vue";
</script>
