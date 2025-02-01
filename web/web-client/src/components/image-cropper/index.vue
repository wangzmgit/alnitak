<template>
  <div v-show="showCropper" class="cropper-mark" v-loading="!imgFile">
    <div class="cropper-card">
      <input class="upload-input" ref="inputRef" type="file" accept="image/*" @change="imageChange">
      <div class="close-icon" @click="closeCropper">
        <close-icon size="26"></close-icon>
      </div>
      <slot name="content" v-if="imgFile" :file="imgFile"></slot>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { Close as CloseIcon } from "@icon-park/vue-next";

const imgFile = ref<File>();
const showCropper = ref(false);
const inputRef = ref<HTMLInputElement | null>(null);

const imageChange = (e: Event) => {
  const target = (e.target as HTMLInputElement)
  const files = target.files;
  if (files && files[0]) {
    imgFile.value = files[0];
    showCropper.value = true;
  }
}

const open = () => {
  if (!imgFile.value) {
    inputRef.value?.click();
  } else {
    showCropper.value = true;
  }
}

const closeCropper = () => {
  if (inputRef.value) {
    inputRef.value.value = "";
    imgFile.value = undefined;
  }
  showCropper.value = false;
}

const setImgFile = (file: File | null) => {
  if (file) {
    imgFile.value = file;
  }
}

defineExpose({
  open,
  setImgFile,
  closeCropper
})

</script>

<style lang="scss" scoped>
.cropper-mark {
  top: 0;
  left: 0;
  z-index: 999;
  position: fixed;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, .16);
  display: flex;
  align-items: center;
  justify-content: center;
}


.cropper-card {
  position: relative;
  padding: 20px;
  background-color: #fff;
  border-radius: 6px;
  min-width: 300px;
  min-height: 200px;
  animation: fadein .3s ease-in;
  box-shadow: 16px 16px 50px -12px rgba(0, 0, 0, 0.8);

  .upload-input {
    display: none;
  }

  // 关闭按钮
  .close-icon {
    position: absolute;
    color: #adadad;
    top: 16px;
    right: 16px;
    cursor: pointer;

    &:hover {
      color: var(--primary-color);
    }
  }
}
</style>
