import { ref, shallowRef } from 'vue';

// 引入主题
import defaultTheme from '@/theme/default';
import darkTheme from '@/theme/dark';

// 定义在全局的样式变量
const themeName = ref<string>('default');
const theme = shallowRef<ThemeType>(defaultTheme);

export function useTheme() {
  const setDefaultTheme = () => {
    theme.value = defaultTheme;
    themeName.value = 'default';
    localStorage.setItem('theme', 'default');
  };

  const setDarkTheme = () => {
    theme.value = darkTheme;
    themeName.value = 'dark';
    localStorage.setItem('theme', 'dark');
  };

  return {
    theme,
    themeName,
    setDefaultTheme,
    setDarkTheme,
  };
}