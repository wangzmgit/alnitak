<template>
  <div class="alnitak-editor">
    <WangToolbar class="toolbar" :editor="editorRef" :defaultConfig="toolbarConfig" :mode="mode" />
    <WangEditor class="editor" v-model="valueHtml" :defaultConfig="editorConfig" :mode="mode" @onCreated="handleCreated"
      @onChange="handleChange" @onDestroyed="handleDestroyed" @onFocus="handleFocus" @onBlur="handleBlur"
      @customAlert="customAlert" @customPaste="customPaste" />
  </div>
</template>


<script  setup lang="ts">
import { uploadFileAPI } from '@/api/upload';
import '@wangeditor/editor/dist/css/style.css' // 引入 css
import { onBeforeUnmount, ref, shallowRef, onMounted } from 'vue';

const emit = defineEmits(['update:content']);
const props = withDefaults(defineProps<{
  content: string; //弹窗可见性
}>(), {
  content: "",
})

const mode = 'default';
// 编辑器实例，必须用 shallowRef
const editorRef = shallowRef();

const toolbarConfig = {
  excludeKeys: [
    "blockquote",
    "bgColor",
    "fontSize",
    "fontFamily",
    "lineHeight",
    "todo",
    "group-indent",
    "emotion",
    "insertTable",
    "group-image",
    "group-video",
    "codeBlock",
    "fullScreen"
  ],
  insertKeys: {
    index: 20,
    keys: "uploadImage"
  }
}
const editorConfig = {
  height: 200,
  placeholder: '请输入文章内容...',
  MENU_CONF: {
    uploadImage: {
      async customUpload(file: File, insertFn: any) {
        uploadFileAPI({
          name: "image",
          action: "v1/upload/image",
          file: file,
          onProgress: () => { },
          onError: () => {
            ElMessage.error("上传失败");
          },
          onFinish: (data?: any) => {
            insertFn(getResourceUrl(data.data.url), "图片", "");
          },
        })
      }
    }
  }
}

const handleCreated = (editor: any) => {
  editorRef.value = editor // 记录 editor 实例，重要！
};

// 内容 HTML
const valueHtml = ref(props.content);
const handleChange = (editor: any) => {
  emit("update:content", editor.getHtml())
};

const handleDestroyed = (editor: any) => {
  // console.log('destroyed', editor);
};

const handleFocus = (editor: any) => {
  // console.log('focus', editor);
};

const handleBlur = (editor: any) => {
  // console.log('blur', editor)
}

const customAlert = (info: any, type: any) => { alert(`【自定义提示】${type} - ${info}`) }
const customPaste = (editor: any, event: any, callback: any) => {
  console.log('ClipboardEvent 粘贴事件对象', event)
  // const html = event.clipboardData.getData('text/html') // 获取粘贴的 html
  const text = event.clipboardData.getData('text/plain') // 获取粘贴的纯文本
  // const rtf = event.clipboardData.getData('text/rtf') // 获取 rtf 数据（如从 word wsp 复制粘贴）

  // 自定义插入内容
  editor.insertText(text)

  // 返回 false ，阻止默认粘贴行为
  event.preventDefault()
  callback(false) // 返回值（注意，vue 事件的返回值，不能用 return）

  // 返回 true ，继续默认的粘贴行为
  // callback(true)
}

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null) return
  editor.destroy()
})
</script> 

<style lang="scss" scoped>
.alnitak-editor {
  border: 1px solid #ccc;

  .toolbar {
    border-bottom: 1px solid #ccc;
  }
}

:deep(.w-e-text-container) {
  min-height: 300px !important;
}
</style>
