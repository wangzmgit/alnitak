<template>
  <vue-picture-cropper :boxStyle="{
    width: '100%',
    height: '100%',
    backgroundColor: '#f8f8f8',
    margin: 'auto',
  }" :options="cropperOptions" :img="localURL" @cropend="change" @ready="change" />
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from 'vue';
import VuePictureCropper, { cropper } from "vue-picture-cropper";

const emits = defineEmits(["change"]);

const props = withDefaults(defineProps<{
  file?: File
  minWidth: number
  minHeight: number
  aspectRatio?: number
  fillColor?: string
}>(), {
  minWidth: 0,
  minHeight: 0,
  aspectRatio: 1 / 1,
  fillColor: "transparent"
})


const cropperOptions = {
  // 是否可移动or放大图片
  movable: false,
  scalable: false,
  zoomable: false,
  zoomOnTouch: false,
  zoomOnWheel: false,
  // 当点击两次时可以在“crop”和“move”之间切换拖拽模式，
  toggleDragModeOnDblclick: false,
  // 固定裁剪框宽高比
  aspectRatio: 1 / 1,
  minCropBoxWidth: 0,
  minCropBoxHeight: 0,
}

const change = () => {
  emits("change", getDataURL());
}

const localURL = ref("");
const loadImage = () => {
  if (props.file)
    localURL.value = URL.createObjectURL(props.file);
}

const getDataURL = () => {
  return cropper?.getDataURL();
}

const getFile = async () => {
  if(!cropper) return;
  const original = props.file?.name || "";
  const fileName = original.substring(0, original.lastIndexOf("."));
  const blob = await cropper.getBlob({
    fillColor: props.fillColor
  });

  // @ts-ignore
  return new File([blob], `${fileName}.png`, { type: "image/png", lastModified: Date.now() });
}

onBeforeMount(() => {
  cropperOptions.aspectRatio = props.aspectRatio;
  cropperOptions.minCropBoxWidth = props.minWidth;
  cropperOptions.minCropBoxHeight = props.minHeight;
  loadImage();
})

defineExpose({
  getFile,
  getDataURL,
})
</script>