// 渲染组件相关函数
import { NIcon } from "naive-ui";
import { h } from "vue";
import type { AnchorHTMLAttributes, Component } from "vue";
import { RouterLink, } from "vue-router";
import type { RouteRecordName } from "vue-router";

/** 渲染图标组件方法（不传入参数返回undefined） */
export const renderIcon = (icon?: Component) => {
  return icon ? () => h(NIcon, null, { default: () => h(icon) }) : undefined;
};

/** 渲染图标组件方法（传入值为 @vicons/ionicons5 的图标名称字符串） */
export const renderIconStr = async (icon?: keyof typeof import("@vicons/ionicons5")) => {
  if (icon) {
    const { [icon]: iconNode } = await import("@vicons/ionicons5");
    return () => h(NIcon, null, { default: () => h(iconNode) });
  } else {
    return undefined;
  }
};

/**
 * 渲染router-link组件方法
 * @param   name    对应routes的跳转name
 * @param   label   菜单名称
 */
export const renderRouterLink = (name: RouteRecordName, label?: string) => {
  return h(RouterLink, { to: { name } }, { default: () => label });
};

/**
 * 渲染a标签方法
 * @param   content     a标签显示内容 <a>content</a>
 * @param   options     a标签对应的属性，如：href，target
 */
export const renderATag = (content: string, options: AnchorHTMLAttributes) => {
  return h("a", { target: "_blank", ...options } as AnchorHTMLAttributes, {
    default: () => content,
  });
};
