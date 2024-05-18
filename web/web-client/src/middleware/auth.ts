import { storageData } from "@/utils/storage-data";

export default defineNuxtRouteMiddleware((to, from) => {
  let isServer = process.server;

  if (isServer) {
    const userId = useCookie('user_id')
    if (!userId.value) {
      return navigateTo(`/login?redirect=${to.path}`);
    }
  } else {
    if (!storageData.get("refreshToken")) {
      return navigateTo(`/login?redirect=${to.path}`);
    }
  }
})

