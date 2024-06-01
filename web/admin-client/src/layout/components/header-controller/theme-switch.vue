<template>
  <n-switch :value="themeSwitchValue"  size="medium" @update:value="handleThemeChange">
    <template #checked-icon>🌙</template>
    <template #unchecked-icon>☀</template>
  </n-switch>
</template>

<script lang="ts" setup>
import { NSwitch } from 'naive-ui';
import { computed } from 'vue';
import useThemeStore from "@/stores/modules/theme-store";

const themeStore = useThemeStore();

// // 主题切换按钮
const themeSwitchValue = computed(() => themeStore.theme === "dark");

const handleThemeChange = (isDark: boolean) => {
  themeStore.theme = isDark ? "dark" : "default";
  themeStore.setThemeConfig();
};
</script>
