import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { storageData } from '@/utils/storage-data';
import HomeView from '../views/home/index.vue';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: HomeView
  },
  {
    path: '/login',
    name: "Login",
    component: () => import("../views/login/index.vue"),
  },
  {
    path: '/register',
    name: "Register",
    component: () => import("../views/register/index.vue"),
  },
  {
    path: '/space',
    name: 'Space',
    meta: { auth: true },
    redirect: '/space/info',
    children: [
      {
        path: '/space/info',
        name: "SpaceInfo",
        meta: { auth: true },
        component: () => import("../views/space/info/index.vue"),
      },
      {
        path: '/space/edit',
        name: "EditInfo",
        meta: { auth: true },
        component: () => import("../views/space/edit/index.vue"),
      },
    ]
  },
  {
    path: '/search',
    name: "Search",
    component: () => import("../views/search/index.vue")
  },
  {
    path: '/video/:vid',
    name: "Video",
    component: () => import("../views/video/index.vue")
  },
  {
    path: '/message',
    name: 'Message',
    meta: { auth: true },
    component: () => import("../views/message/index.vue"),
  },
  {
    path: '/announce',
    name: "Announce",
    meta: { auth: true },
    component: () => import("../views/message/announce/index.vue"),
  },
  {
    path: '/message/whisper',
    name: 'Whisper',
    meta: { auth: true },
    component: () => import("../views/message/whisper/index.vue"),
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes
})

router.beforeEach((to, from, next) => {
  // 是否需要登录
  if (to.meta.auth && !storageData.get('token')) {
    if (from.name !== "Home") {
      router.push({ name: 'Login' });
      next();
    }
  } else {
    next();
  }
})

export default router