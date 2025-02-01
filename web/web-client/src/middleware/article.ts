const articlePages = ["article-list","article-id", "upload-article", "upload-article-manage"]
export default defineNuxtRouteMiddleware((to, from) => {
  if (!globalConfig.article) {
    if (to.name && articlePages.includes(to.name.toString())) {
      return navigateTo("/404");
    }
  }
})

