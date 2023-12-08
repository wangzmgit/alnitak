import { ref } from "vue";
import { defineStore } from "pinia";
import { useTheme } from "@/hooks/theme-hooks";

const { setDarkTheme, setDefaultTheme } = useTheme();

const useThemeStore = defineStore('theme', () => {

  const theme = ref('default');

  const setThemeConfig = async () => {
    if(theme.value === 'default') {
      setDefaultTheme();
    } else {
      setDarkTheme();
    }
  }

  return {
    theme,
    setThemeConfig
  }
})

export default useThemeStore;
