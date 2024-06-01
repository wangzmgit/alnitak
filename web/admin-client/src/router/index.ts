import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import Layout from "@/layout/index.vue";
import Login from "@/views/login/index.vue";
import Error404 from "@/views/errors/404.vue";
import Error403 from "@/views/errors/403.vue";

export const routes: RouteRecordRaw[] = [
  {
    path: "/",
    name: "layout",
    component: Layout,
    children: [],
  },
  {
    path: "/login",
    name: "login",
    meta: { hidden: true },
    component: Login,
  },
  {
    path: "/404",
    name: "error404",
    meta: { hidden: true },
    component: Error404,
  },
  {
    path: "/403",
    name: "error403",
    meta: { hidden: true },
    component: Error403,
  },
  {
    path: "/error",
    component: Layout,
    name: "errorPage",
    meta: { label: "错误页", icon: "SkullOutline" },
    children: [
      {
        path: "403",
        name: "errorPage403",
        meta: { label: "403", icon: "TrashBinOutline" },
        component: () => import("@/views/errors/403.vue"),
      },
      {
        path: "404",
        name: "errorPage404",
        meta: { label: "404", icon: "TrashBinOutline" },
        component: () => import("@/views/errors/404.vue"),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
