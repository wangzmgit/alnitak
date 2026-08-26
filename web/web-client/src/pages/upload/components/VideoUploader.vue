<template>
  <div class="video-uploader">
    <el-upload :key="uploaderKey" drag multiple :show-file-list="false" :before-upload="beforeUploadVideo"
      @change="handleChange">
      <el-progress v-if="uploading" type="circle" :percentage="percent" />
      <div class="upload-tips" v-else>
        <div class="tips-icon">
          <client-only>
            <upload-icon size="48"></upload-icon>
          </client-only>
        </div>
        <p class="upload-title">点击或拖拽视频到此处上传视频</p>
        <p class="upload-limit">上传文件大小需小于{{ globalConfig.maxVideoSize }}M,仅支持.mp4格式文件</p>
        <el-button class="btn" type="primary">上传视频</el-button>
      </div>
    </el-upload>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessage } from "element-plus";
import { uploadFileChunkAPI } from "@/api/upload";
import { globalConfig } from "@/utils/global-config";
import { Upload as UploadIcon } from "@icon-park/vue-next";

const emits = defineEmits<{
  (e: 'finish', data: any, coverURL?: string | null): void
}>();
const props = defineProps<{
  vid?: number
}>();


const percent = ref(0);//上传百分比
const uploading = ref(false);//是否上传中

/** MP4 文件魔数（文件头前 4 字节：ftyp 类型） */
const MP4_MAGIC_PREFIXES = [
  new Uint8Array([0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70]),  // mp42 / isom
  new Uint8Array([0x00, 0x00, 0x00, 0x1C, 0x66, 0x74, 0x79, 0x70]),  // mp42 variant
  new Uint8Array([0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70]),  // dash / iso5
];

/** 读取文件头字节并与已知魔数比对 */
const isMp4ByMagicBytes = async (file: File): Promise<boolean> => {
  try {
    const header = await file.slice(0, 8).arrayBuffer();
    const view = new Uint8Array(header);
    return MP4_MAGIC_PREFIXES.some(prefix =>
      prefix.every((b, i) => b === view[i])
    );
  } catch {
    return false;
  }
};

//上传之前的回调
const beforeUploadVideo = async (options: any) => {
  const file = options.file;
  const isMp4Type = file.type === "video/mp4";
  const isMp4Magic = await isMp4ByMagicBytes(file.file);
  if (!isMp4Type || !isMp4Magic) {
    ElMessage.error("文件只支持mp4格式");
  }
  const isLtMaxSize = file.file.size / 1024 / 1024 < globalConfig.maxVideoSize;
  if (!isLtMaxSize) {
    ElMessage.error(`视频大小不能超过${globalConfig.maxVideoSize}M`);
  }
  return (isMp4Type && isMp4Magic) && isLtMaxSize;
}

//上传变化的回调
const handleChange = (uploadFile: any) => {
  if (!uploadFile.raw) return;
  uploading.value = true;
  uploadFileChunkAPI({
    name: "video",
    action: props.vid ? `v1/upload/video/${props.vid}` : `v1/upload/video`,
    file: uploadFile.raw,
    onProgress: (val: any) => {
      changeUpload("uploading", val);
    },
    onError: () => {
      changeUpload("error");
      uploading.value = false;
    },
    onFinish: (data?: any) => {
      changeUpload("success", data);
      uploading.value = false;
    },
  })
}

//上传变化的回调
const uploaderKey = ref(0);
const changeUpload = (status: string, data?: any) => {
  switch (status) {
    case "success":
      emits("finish", data);
      ElMessage.success('上传完成');
      uploaderKey.value = Date.now();
      break;
    case "uploading":
      percent.value = data;
      break;
    case "error":
      ElMessage.error('文件上传失败');
      uploaderKey.value = Date.now();
      break;
  }
}
</script>

<style lang="scss" scoped>
.video-uploader {
  width: 80%;
  margin: 0 auto;

  .upload-tips {
    padding: 24px;

    .tips-icon {
      color: var(--font-primary-2);
      opacity: 0.6;
      margin-bottom: 12px;
    }

    .upload-title {
      color: var(--font-primary-1);
      font-size: 16px
    }

    .upload-limit {
      color: var(--font-primary-3);
      font-size: 14px;
      margin: 8px 0 0 0;
    }

    .btn {
      margin-top: 26px;
      width: 200px;
      height: 40px;
    }
  }
}
</style>