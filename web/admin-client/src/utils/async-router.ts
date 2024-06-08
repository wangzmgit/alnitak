import type { RouteRecordRaw } from "vue-router";
import { globalConfig } from "@/utils/global-config";

const viewModules = import.meta.glob('../views/**/*.vue');

export const asyncRouterHandle = (menuList: MenuType[]) => {
  const routeList: RouteRecordRaw[] = [];

  menuList.forEach((item) => {
    if (item.children) {
      item.children.forEach((child) => {
        routeList.push({
          name: child.name,
          path: `/${globalConfig.baseUrl}/layout/${child.path}`,
          meta: {
            label: child.meta.title,
            keepAlive: child.meta.keepAlive,
          },
          component: dynamicImport(child.component)
        })
      })
    } else if (item.component) { // 没有子路由
      routeList.push({
        name: item.name,
        path: item.path,
        meta: {
          label: item.meta.title,
          keepAlive: item.meta.keepAlive,
        },
        component: dynamicImport(item.component)
      })
    }
  })

  return routeList;
}


function dynamicImport(component: string) {
  const keys = Object.keys(viewModules);
  const matchKeys = keys.filter((key) => {
    const k = key.replace('../', '')
    return k === component;
  })
  const matchKey = matchKeys[0]

  return viewModules[matchKey]
}
