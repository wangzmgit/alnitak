import { ref } from "vue";
import { defineStore } from "pinia";

const useSystemStore = defineStore('system', () => {
  /** 头部容器高度 */
  const headerHeight = ref(0);
  /** 主内容容器高度 */
  const contentHeight = ref(0);
  /** 主内容容器高度（不含padding以外层，border-box） */
  const contentBorderBoxHeight = ref(0);
  /** 主内容容器宽度（不含padding以外层） */
  const contentWidth = ref(0);

  return {
    headerHeight,
    contentHeight,
    contentBorderBoxHeight,
    contentWidth,
  }
})

export default useSystemStore;
