import { defineNuxtPlugin } from '#app';
import { Editor, Toolbar } from '@wangeditor/editor-for-vue';

export default defineNuxtPlugin(nuxtApp => {
  nuxtApp.vueApp.component('WangEditor', Editor)
  nuxtApp.vueApp.component('WangToolbar', Toolbar)
})