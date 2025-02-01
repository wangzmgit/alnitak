<template>
  <div class="avatar-cropper">
    <div class="card-left">
      <picture-cropper ref="cropperRef" :key="props.file?.name" :file="props.file" :minHeight="60" :minWidth="60"
        @change="imgChange"></picture-cropper>
    </div>
    <div class="card-right">
      <span class="title">头像预览</span>
      <img :src="previewURL" alt="图像预览" />
      <el-button type="primary" @click="uploadAvatar">裁剪并上传</el-button>
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
const uploadAvatar = async () => {
  if (cropperRef.value) {
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
}
</script>

<style lang="scss" scoped>
.avatar-cropper {
  box-sizing: border-box;
  padding: 10px;
  display: flex;
  height: 280px;

  .card-left {
    width: 260px;
    text-align: center
  }

  .card-right {
    width: 200px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-direction: column;

    img {
      width: 120px;
      height: 120px;
    }
  }
}
</style>