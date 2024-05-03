<template>
  <div class="upload-cover">
    <n-upload multiple directory-dnd :show-file-list="false" @before-upload="beforeUploadCover"
      :custom-request="fileChange">
      <img v-if="currentCover" class="cover" :src="getResourceUrl(currentCover)" alt="封面" />
      <n-upload-dragger v-else>
        <div v-if="!uploading">
          <div style="margin-bottom: 12px">
            <n-icon size="48" :depth="3">
              <upload-icon></upload-icon>
            </n-icon>
          </div>
          <n-text style="font-size: 16px">点击或拖拽图片到此处上传封面</n-text>
          <n-p depth="3" style="margin: 8px 0 0 0">
            上传文件大小需小于{{ globalConfig.maxImgSize }}M,仅支持.jpg .jpeg .png格式文件
          </n-p>
        </div>
        <n-progress v-else type="circle" :percentage="percent" />
      </n-upload-dragger>
    </n-upload>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { uploadFileAPI } from "@/api/upload";
import { getResourceUrl } from "@/utils/resource";
import { globalConfig } from "@/utils/global-config";
import { CloudUpload as UploadIcon } from "@vicons/ionicons5";
import type { UploadCustomRequestOptions } from 'naive-ui';
import { NIcon, NUpload, NUploadDragger, NText, NP, NProgress, useMessage } from 'naive-ui';

const emits = defineEmits(["finish"]);
const props = defineProps<{
  cover?: string
}>()

const currentCover = ref(props.cover);

const percent = ref(0);//上传百分比
const uploading = ref(false);//是否上传中
const message = useMessage();//通知

//上传之前的回调
const beforeUploadCover = async (options: any) => {
  const file = options.file;
  const isJpgOrPng =
    file.type === "image/jpeg" || file.type === "image/png";
  if (!isJpgOrPng) {
    message.error("文件只支持jpg/jpeg/png格式");
  }
  const isLtMaxSize = file.file.size / 1024 / 1024 < globalConfig.maxImgSize;
  if (!isLtMaxSize) {
    message.error(`图片大小不能超过${globalConfig.maxImgSize}M`);
  }
  return isJpgOrPng && isLtMaxSize;
}

const fileChange = (options: UploadCustomRequestOptions) => {
  if (options.file.file) {
    uploadFileAPI({
      name: "image",
      action: "v1/upload/image",
      file: options.file.file,
      onProgress: () => { },
      onError: () => {
        changeUpload("error");
      },
      onFinish: (data?: any) => {
        options.onFinish();
        changeUpload("success", data);
      },
    })
  }
}

//上传变化的回调
const changeUpload = (status: string, data?: any) => {
  switch (status) {
    case "success":
      //更新封面图
      currentCover.value = data.data.url;
      emits("finish", currentCover.value);
      message.success('上传完成');
      break;
    case "error":
      message.error("文件上传失败");
      break;
  }
}
</script>

<style lang="scss" scoped>
.upload-cover {
  width: 350px;
  margin: 20px auto;

  .cover {
    width: 350px;
    height: 200px;
  }

}
</style>