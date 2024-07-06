<template>
  <div class="upload-cover">
    <el-upload multirple :show-file-list="false" :before-upload="beforeUploadCover" @change="fileChange">
      <img v-if="currentCover" :src="getResourceUrl(currentCover)" class="cover" alt="封面" />
      <div v-else class="cover placeholder">
        <div class="tips-icon">
          <add-picture-icon size="22"></add-picture-icon>
        </div>
        <p class="upload-title">添加封面</p>
      </div>
    </el-upload>

    <image-cropper ref="cropperRef">
      <template #content="fileSlot">
        <cover-cropper :file="fileSlot.file" @state-change="changeUpload"></cover-cropper>
      </template>
    </image-cropper>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessage } from "element-plus";
import { isUrl } from "@/utils/verify";
import { globalConfig } from "@/utils/global-config";
import { AddPicture as AddPictureIcon } from "@icon-park/vue-next";
import ImageCropper from "@/components/image-cropper/index.vue";
import CoverCropper from "@/components/image-cropper/components/CoverCropper.vue";

const emits = defineEmits(["finish"]);
const props = defineProps<{
  cover?: string
}>()

const currentCover = ref(props.cover);

//上传之前的回调
const beforeUploadCover = async (file: any) => {
  const isJpgOrPng =
    file.type === "image/jpeg" || file.type === "image/png";
  if (!isJpgOrPng) {
    ElMessage.error("文件只支持jpg/jpeg/png格式");
  }
  const isLtMaxSize = file.file.size / 1024 / 1024 < globalConfig.maxImgSize;
  if (!isLtMaxSize) {
    ElMessage.error(`图片大小不能超过${globalConfig.maxImgSize}M`);
  }

  return isJpgOrPng && isLtMaxSize;
}

const cropperRef = ref<InstanceType<typeof ImageCropper> | null>(null);
const fileChange = (uploadFile: any) => {
  cropperRef.value?.setImgFile(uploadFile.raw);
  cropperRef.value?.open();
}

//上传变化的回调
const changeUpload = (status: string, data: any) => {
  switch (status) {
    case "success":     //更新封面图
      // 如果返回内容本身为url或者配置中不存在域名则直接返回
      if (isUrl(data.data.url) || !globalConfig.domain) {
        currentCover.value = data.data.url;
      } else {
        currentCover.value = `http${globalConfig.https ? 's' : ''}://${globalConfig.domain}${data.data.url}`;
      }
      emits("finish", data.data.url);
      break;
    case "error":
      ElMessage.error('封面上传失败');
      break;
  }
  cropperRef.value?.closeCropper();
}
</script>

<style lang="scss" scoped>
.upload-cover {
  width: 169px;
  height: 127px;
  border-radius: 3px;

  .cover {
    width: 169px;
    height: 127px;
    border-radius: 3px;
  }

  .placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    background-color: #f4f5f7;

    .tips-icon {
      color: #99a2aa;
      margin-bottom: 12px;
    }

    .upload-title {
      color: #99a2aa;
      font-size: 12px;
      margin: 0;
    }
  }
}

:deep(.el-upload-dragger) {
  padding: 0px;
}

:deep(.el-upload-dragger:hover) {
  border-color: var(--primary-color);
}
</style>