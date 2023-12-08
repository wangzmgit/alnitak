// 菜单列表
import { ref } from "vue";
import router from "@/router";
import { defineStore } from "pinia";
import type { MenuOption } from "naive-ui";
import { renderRouterLink, renderIconStr } from "@/utils/render";
import { getUserMenuAPI } from "@/api/menu";
import { statusCode } from "@/utils/status-code";
import { asyncRouterHandle } from "@/utils/async-router";

// 获取路由数据
const getMenu = async () => {
  const res = await getUserMenuAPI();
  if (res.data.code === statusCode.OK) {
    return res.data.data.menus;
  }
  return [];
}

const handelMenu = async (list: MenuType[]) => {
  const finalList: MenuOption[] = [];

  for (let i = 0; i < list.length; i++) {
    const item = list[i];
    if (!item.meta.hidden) {
      const finalItem: MenuOption = {
        key: item.name,
        label: item.meta.title,
        disabled: item.meta.disabled,
        icon: await renderIconStr(item.meta.icon)
      };

      if (item.children) { // 如果有子路由则外层路由不做跳转使用
        const childList: MenuOption[] = [];
        for (let j = 0; j < item.children.length; j++) {
          const child = item.children[j];
          if (!child.meta.hidden) {
            const childItem: MenuOption = {
              key: child.name,
              label: () => renderRouterLink(child.name, child.meta.title),
              disabled: child.meta.disabled,
              icon: await renderIconStr(child.meta.icon)
            };

            childList.push(childItem);
          }
        }

        if (childList.length) {
          finalItem.children = childList;
        }
      } else {
        finalItem.label = () => renderRouterLink(item.name, item.meta.title);
      }

      finalList.push(finalItem);
    }
  }

  return finalList;
}

/** 菜单信息仓库 */
const useMenuStore: any = defineStore('menu', () => {
  /** 菜单是否加载完成 */
  const isMenuLoaded = ref(false);
  /** 菜单列表，用于naive的menu */
  const list = ref<MenuOption[]>([]);
  /** 根据routes创建菜单列表 */
  const createMenuList = async () => {
    const menus = await getMenu();
    // 处理路由
    const routerList = asyncRouterHandle(menus);
    routerList.forEach((item) => {
      router.addRoute('layout', item);
    })

    // 处理菜单数据
    // @ts-ignore
    list.value = await handelMenu(menus);
    isMenuLoaded.value = true;
  }



  return {
    list,
    isMenuLoaded,
    createMenuList
  }
});

export default useMenuStore;
