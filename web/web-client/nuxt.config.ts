// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@element-plus/nuxt',
    '@pinia/nuxt',
  ],
  css: [
    'element-plus/dist/index.css',
    '~/assets/styles/element.scss'
  ],
  devtools: { enabled: true },
  srcDir: 'src/',
})
