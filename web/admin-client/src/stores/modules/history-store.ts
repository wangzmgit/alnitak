// 路由历史记录
import { defineStore } from "pinia";
import { ref } from "vue";
import type { VNode } from "vue";
import router from "@/router"
import type { RouteLocationNormalized } from "vue-router";

/** 路由历史记录数组元素类型 */
export interface routerHistoryItem extends RouteLocationNormalized {
  iconNode?: () => VNode;
}

const useHistoryStore = defineStore("history", () => {
  /** 路由历史记录 */
  const routerHistory = ref<routerHistoryItem[]>([]);
  /** keep-alive include字段，所需要keep-alive的路由name数组集合 */
  const keepAliveInclude = ref<string[]>([]);

  /**
  * 存储路由历史记录
  */
  const setRouterHistory = async (history: RouteLocationNormalized) => {
    const index = routerHistory.value.findIndex(item => item.path === history.path);
    const data = {
      ...history,
    };
    if (!~index) {
      routerHistory.value.push(data);
    } else {
      routerHistory.value[index] = data;
    }
  }
  /**
   * 删除单个路由历史记录
   */
  const removeRouterHistory = (index: number) => {
    const { path: curPath } = router.currentRoute.value;
    const finalIndex = routerHistory.value.length - 1;
    // 是否为当前页面
    if (curPath === routerHistory.value[index].path) {
      // 少于等于一个元素，删除后列表为空，跳转到首页
      if (routerHistory.value.length <= 1) {
        router.push("/");
      }
      // 如果为最后一个元素，删除后跳转到前一个元素页面
      else if (index === finalIndex) {
        router.push(routerHistory.value[index - 1].path);
      }
      // 不为最后一个元素，即删除后后面还有元素，跳转到后一个元素页面
      else {
        router.push(routerHistory.value[index + 1].path);
      }
    }
    routerHistory.value.splice(index, 1);
  }

  /**
   * 清除路由历史记录
   */
  const clearRouterHistory = () => {
    routerHistory.value = [];
  }

  /** 添加keep-alive include字段 */
  const addKeepAliveInclude = (routeName: string) => {
    if (!keepAliveInclude.value.includes(routeName)) {
      keepAliveInclude.value.push(routeName);
    }
  }

  return {
    routerHistory,
    keepAliveInclude,
    setRouterHistory,
    removeRouterHistory,
    clearRouterHistory,
    addKeepAliveInclude
  }
})


export default useHistoryStore;
