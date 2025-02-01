<template>
  <div class="title-bar">
    <h2 class="title-text">文件上传</h2>
    <el-upload :show-file-list="false" :before-upload="beforeUploadVideo" @change="handleChange">
      <el-button type="primary" :icon="Plus">添加视频</el-button>
    </el-upload>
  </div>
  <div class="upload-video">
    <div class="video-item" v-for="(item, index) in resourceList" :key="index">
      <div class="video-icon-box">
        <el-icon :size="38">
          <monitor-icon></monitor-icon>
        </el-icon>
        <span class="part"> P{{ index + 1 }} </span>
      </div>
      <div class="info-box">
        <div class="file-info">
          <div class="title-box">
            <div class="title" v-if="modifyIndex !== index" @click="titleClick(item, index)">
              <span>{{ item.title || "未命名视频" }}</span>
            </div>
            <el-input v-else ref="titleInput" v-model="modifyForm.title" maxlength="50" show-word-limit
              @blur="modifyTitle(item)" />
          </div>
          <client-only>
            <el-popconfirm title="是否移除该条视频？" confirm-button-text="确认" cancel-button-text="取消"
              @confirm="deleteResource(item.id, index)">
              <template #reference>
                <span class="remove-btn" v-if="resourceList.length > 1">移除</span>
              </template>
            </el-popconfirm>
          </client-only>
        </div>
        <div class="progress-box">
          <span class="upload-status">{{ item.uploading ? `上传中 ${item.percent}%` : getTagText(item.status) }}</span>
          <div class="progress-bar">
            <div class="progress" :style="`width: ${item.uploading ? item.percent : 100}%`"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, nextTick, ref } from "vue";
import { Plus } from '@icon-park/vue-next';
import { reviewCode } from '@/utils/review-code';
import { ElIcon, ElButton, ElInput, ElPopconfirm } from "element-plus";
import MonitorIcon from "@/components/icons/MonitorIcon.vue";
import { submitReviewAPI, getVideoStatusAPI } from "@/api/video";
import { deleteResourceAPI, modifyTitleAPI } from "@/api/resource";
import { uploadFileChunkAPI } from "@/api/upload";

const emit = defineEmits(["review"]);
const props = defineProps<{
  vid: number,
  resources: Array<ResourceType>
}>();

const resourceList = ref<Array<ResourceType | UploadResourceType>>(props.resources);

// 获取标签文本
const getTagText = (state: number) => {
  switch (state) {
    case reviewCode.AUDIT_APPROVED:
      return '审核通过';
    case reviewCode.PROCESSING_FAIL:
      return '处理失败';
    default:
      return '上传成功';
  }
}

// 提交审核
const submitReview = async () => {
  if (resourceList.value.length === 0) {
    ElMessage.error('请先上传视频');
    return;
  }
  const res = await submitReviewAPI(props.vid);
  if (res.data.code === statusCode.OK) {
    emit("review");
  } else {
    ElMessage.error(res.data.msg || '提交失败');
  }
}

const deleteResource = async (id: number, index: number) => {
  const res = await deleteResourceAPI(id);
  if (res.data.code === statusCode.OK) {
    resourceList.value.splice(index, 1);
  } else {
    ElMessage.error(res.data.msg || '删除失败');
  }
}

//修改资源名
const modifyIndex = ref(-1);
const titleInput = ref<Array<InstanceType<typeof ElInput>>>();
const modifyForm = reactive<BaseResourceType>({
  id: 0,
  title: '',
});

//点击标题
const titleClick = (resource: ResourceType | UploadResourceType, index: number) => {
  modifyForm.id = resource.id;
  modifyForm.title = resource.title;
  modifyIndex.value = index;
  nextTick(() => {
    if (titleInput.value) {
      titleInput.value[0].focus();
    }
  })
}

//修改标题
const modifyTitle = async (resource: ResourceType | UploadResourceType) => {
  modifyIndex.value = -1;
  if (!modifyForm.title) return;
  const res = await modifyTitleAPI(modifyForm);
  if (res.data.code === statusCode.OK) {
    resource.title = modifyForm.title;
  } else {
    ElMessage.error('修改失败');
  }
}

//上传之前的回调
const beforeUploadVideo = async (options: any) => {
  const file = options.file;
  const isJpgOrPng = file.type === "video/mp4";
  if (!isJpgOrPng) {
    ElMessage.error("文件只支持mp4格式");
  }
  const isLtMaxSize = file.file.size / 1024 / 1024 < globalConfig.maxVideoSize;
  if (!isLtMaxSize) {
    ElMessage.error(`视频大小不能超过${globalConfig.maxVideoSize}M`);
  }
  return isJpgOrPng && isLtMaxSize;
}

//上传变化的回调
const handleChange = (uploadFile: any) => {
  if (!uploadFile.raw) return;

  const uploadData: UploadResourceType = {
    id: 0,
    status: -1,
    title: "",
    percent: 0,
    uploading: true,
  }

  const index = resourceList.value.push(uploadData) - 1;

  uploadFileChunkAPI({
    name: "video",
    action: props.vid ? `v1/upload/video/${props.vid}` : `v1/upload/video`,
    file: uploadFile.raw,
    onProgress: (val: any) => {
      uploadData.percent = val;
      resourceList.value[index] = JSON.parse(JSON.stringify(uploadData));
    },
    onError: () => {
      // 在上传列表中移除
      resourceList.value.splice(index, 1);
    },
    onFinish: (data?: any) => {
      resourceList.value[index] = data.data.resource;
    },
  })
}

watch(() => props.resources, (newVal) => {
  if (newVal) {
    resourceList.value = newVal;
  }
})
</script>

<style lang="scss" scoped>
.title-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  height: 50px;
  padding: 10px 20px;

  .title-text {
    margin: 0;
    font-size: 16px;
    color: #212121;
    font-weight: 600;
    line-height: 50px;
  }
}

.upload-video {
  width: 100%;
  margin: 0 auto;
  padding: 0 20px;
  box-sizing: border-box;

  .video-item {
    display: flex;
    align-items: center;
    width: 100%;
    padding: 20px 0;

    &:first-child {
      padding: 10px 0 20px 0;
    }

    &:hover {
      .info-box {
        .file-info {
          .remove-btn {
            display: block;
          }
        }
      }
    }

    .video-icon-box {
      position: relative;
      color: var(--primary-hover-color);

      .part {
        display: block;
        position: absolute;
        width: 28px;
        height: 24px;
        left: 2px;
        top: 7px;
        font-size: 12px;
        line-height: 24px;
        text-align: center;
        color: #fff;
      }
    }

    .info-box {
      flex: 1;
      padding: 0 12px;

      .file-info {
        display: flex;
        align-items: center;
        justify-content: space-between;

        .title-box {
          flex: 1;
          height: 100%;
          display: flex;
          align-items: center;

          .title {
            height: 32px;
            display: flex;
            align-items: center;
          }
        }

        .remove-btn {
          display: none;
          font-size: 12px;
          color: #61666d;
          cursor: pointer;

          &:hover {
            color: var(--primary-hover-color);
          }
        }
      }

      .progress-box {
        width: 100%;
        margin-top: 8px;

        .upload-status {
          display: block;
          color: #61666d;
          padding-bottom: 2px;
          font-size: 12px;
        }

        .progress-bar {
          position: relative;
          width: 100%;
          height: 4px;
          background-color: #efeff5;

          .progress {
            position: absolute;
            left: 0;
            top: 0;
            height: 4px;
            transition: all .2s ease;
            background-color: var(--primary-color);
          }
        }
      }
    }
  }
}
</style>