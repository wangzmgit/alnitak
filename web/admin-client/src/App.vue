<template>
  <n-config-provider :theme-overrides="themeOverrides" :theme="nTheme" :locale="zhCN" :date-locale="dateZhCN" abstract>
    <n-notification-provider>
      <n-message-provider>
        <n-dialog-provider>
          <global-container />
        </n-dialog-provider>
      </n-message-provider>
    </n-notification-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { zhCN, dateZhCN, darkTheme } from "naive-ui";
import GlobalContainer from "@/layout/global.vue";
import { useTheme } from '@/hooks/theme-hooks';
import type { GlobalThemeOverrides } from 'naive-ui';
import { NNotificationProvider, NConfigProvider, NMessageProvider, NDialogProvider } from 'naive-ui';

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

<style lang="scss">
body,
#app {
  width: 100vw;
  height: 100vh;
}

.n-text.n-text--code {
  margin: 0 3px;
}

.n-result {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>
