<template>
  <div class="alnitak-editor">
    <WangEditor v-model="valueHtml" :defaultConfig="editorConfig" :mode="mode" @onCreated="handleCreated" />
  </div>
</template>


<script  setup lang="ts">
import '@wangeditor/editor/dist/css/style.css' // 引入 css
import { onBeforeUnmount, ref, shallowRef } from 'vue';

const props = withDefaults(defineProps<{
  content: string; //弹窗可见性
}>(), {
  content: "",
})

const mode = 'default';
const editorRef = shallowRef();
const editorConfig = { placeholder: '', readOnly: true }

const handleCreated = (editor: any) => {
  editorRef.value = editor // 记录 editor 实例，重要！
};

// 内容 HTML
const valueHtml = ref(props.content);

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  if (editorRef.value) {
    editorRef.value.destroy()
    editorRef.value = null;
  }
})
</script> 

<style lang="scss" scoped>
.alnitak-editor {
  margin-top: 16px;
  border-radius: 3px;
}
</style>
