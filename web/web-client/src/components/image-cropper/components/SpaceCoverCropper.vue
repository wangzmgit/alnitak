<template>
  <div class="space-cover-cropper">
    <div class="cropper-container">
      <picture-cropper ref="cropperRef" :aspectRatio="(5 / 1)" :key="props.file?.name" :file="props.file"
        :minHeight="55" :minWidth="275" fillColor="#e6e6e6" @change="imgChange"></picture-cropper>
    </div>
    <div class="btn-container">
      <button @click="uploadCover">裁剪并上传</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import PictureCropper from './PictureCropper.vue';
import { uploadFileAPI } from "@/api/upload";

const props = defineProps<{
  file?: File
}>();

const cropperRef = ref<InstanceType<typeof PictureCropper> | null>(null);

const previewURL = ref("");
const imgChange = (url: string) => {
  previewURL.value = url;
}

const emits = defineEmits(["stateChange"])
const uploadCover = async () => {
  if (!cropperRef.value) return;
  const file = await cropperRef.value.getFile();
  if (!file) return;
  uploadFileAPI({
    name: "image",
    action: "v1/upload/image",
    file: file,
    onProgress: () => { },
    onError: () => {
      emits("stateChange", "error")
    },
    onFinish: (data?: any) => {
      emits("stateChange", "success", data)
    },
  })
}

</script>

<style lang="scss" scoped>
.space-cover-cropper {
  box-sizing: border-box;
  padding: 10px;
  max-width: 700px;
  max-height: 360px;

  .cropper-container {
    padding-top: 20px;
    width: 100%;
    height: 280px;
  }

  .btn-container {
    height: 60px;
    position: relative;

    button {
      position: absolute;
      right: 0;
      bottom: 6px;
      width: 100px;
      height: 36px;
      border-radius: 5px;
      cursor: pointer;
      border: none;
      color: var(--primary-text-color);
      background-color: var(--primary-color);

      &:hover {
        background-color: var(--primary-hover-color);
      }
    }
  }
}
</style>