<template>
  <n-config-provider :theme-overrides="themeOverrides" :theme="nTheme" :locale="zhCN" :date-locale="dateZhCN">
    <n-notification-provider>
      <n-message-provider>
        <router-view></router-view>
      </n-message-provider>
    </n-notification-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { zhCN, dateZhCN, darkTheme } from "naive-ui";
import { useTheme } from '@/hooks/theme-hooks';
import type { GlobalThemeOverrides } from 'naive-ui';
import { NNotificationProvider, NConfigProvider, NMessageProvider } from 'naive-ui';

const { theme, themeName } = useTheme();

const nTheme = computed(() => (themeName.value === "dark" ? darkTheme : null));

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: theme.value.primaryColor
  },
  Button: {
    // default
    textColor: theme.value.primaryColor,
    textColorFocus: theme.value.primaryHoverColor,
    borderFocus: `1px solid ${theme.value.primaryColor}`,
    textColorHover: theme.value.primaryHoverColor,
    borderHover: `1px solid ${theme.value.primaryHoverColor}`,
    textColorPressed: theme.value.primaryHoverColor,
    borderPressed: `1px solid ${theme.value.primaryHoverColor}`,


    // primary
    borderHoverPrimary: `1px solid ${theme.value.primaryHoverColor}`,
    colorHoverPrimary: theme.value.primaryHoverColor,
    colorFocusPrimary: theme.value.primaryHoverColor,
    borderFocusPrimary: `1px solid ${theme.value.primaryHoverColor}`,
    colorPressedPrimary: theme.value.primaryHoverColor,
    borderPressedPrimary: `1px solid ${theme.value.primaryHoverColor}`,

    // primary-ghost
    textColorGhostHoverPrimary: theme.value.primaryHoverColor,
    textColorGhostPressedPrimary: theme.value.primaryHoverColor,

    //text-default
    textColorText: '#999',
    textColorTextHover: theme.value.primaryColor,
    textColorTextPressed: theme.value.primaryColor,

    // text-primary
    textColorTextHoverPrimary: theme.value.primaryHoverColor,
  },
  Select: {
    peers: {
      InternalSelection: {
        textColor: theme.value.primaryColor,
        borderFocus: `1px solid ${theme.value.primaryColor}`,
        borderHover: `1px solid ${theme.value.primaryColor}`
      },
    },
  },
  Input: {
    borderFocus: `1px solid ${theme.value.primaryColor}`,
    borderHover: `1px solid ${theme.value.primaryColor}`,
  },
  Menu: {
    itemTextColorHoverHorizontal: theme.value.primaryHoverColor,
  }
}
</script>

<style>
@import url(@/assets/styles/global.scss);

body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
    "Helvetica Neue", Arial, "Noto Sans", sans-serif, "Apple Color Emoji",
    "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
}
</style>
