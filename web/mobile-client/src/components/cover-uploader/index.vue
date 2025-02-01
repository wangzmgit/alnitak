<template>
  <div class="upload-cover">
    <el-upload drag multiple :show-file-list="false" :before-upload="beforeUploadCover" @change="fileChange">
      <img v-if="currentCover" :src="getResourceUrl(currentCover)" class="cover" alt="封面" />
      <el-progress v-else-if="uploading" type="circle" :percentage="percent" />
      <div class="upload-tips" v-else>
        <div class="tips-icon">
          <upload-icon size="48"></upload-icon>
        </div>
        <p class="upload-title">点击或拖拽图片到此处上传封面</p>
        <p class="upload-limit">上传文件大小需小于{{ globalConfig.maxImgSize }}M,仅支持.jpg .jpeg .png格式文件</p>
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
import { Upload as UploadIcon } from "@icon-park/vue-next";
import ImageCropper from "@/components/image-cropper/index.vue";
import CoverCropper from "@/components/image-cropper/components/CoverCropper.vue";

const emits = defineEmits(["finish"]);
const props = defineProps<{
  cover?: string
}>()

const currentCover = ref(props.cover);

const percent = ref(0);//上传百分比
const uploading = ref(false);//是否上传中

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
  width: 350px;
  height: 200px;
  margin: 20px auto;

  .upload-tips {
    padding: 24px;

    .tips-icon {
      color: #000;
      opacity: 0.38;
      margin-bottom: 12px;
    }

    .upload-title {
      color: #333639;
      font-size: 16px
    }

    .upload-limit {
      color: #767c82;
      font-size: 14px;
      margin: 8px 0 0 0;
    }
  }

  .cover {
    width: 350px;
    height: 200px;
  }
}

:deep(.el-upload-dragger) {
  padding: 0px;
}

:deep(.el-upload-dragger:hover) {
  border-color: var(--primary-color);
}
</style>