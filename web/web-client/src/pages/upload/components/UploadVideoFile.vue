<template>
  <div class="upload-video">
    <video-uploader :vid="vid" @finish="finishUpload"></video-uploader>
    <div class="video-box">
      <el-scrollbar height="300px">
        <div class="video-item" v-for="(item, index) in resourceList">
          <span class="part"> P{{ index + 1 }} </span>
          <div class="title-box">
            <div class="title" v-if="modifyIndex !== index" @click="titleClick(item, index)">
              <span>{{ item.title }}</span>
              <!---<n-tag class="tag" :type="toTagType(item.status)">{{ toTagText(item.status) }}</n-tag> -->
            </div>
            <el-input v-else ref="titleInput" v-model="modifyForm.title" maxlength="50" show-word-limit
              @blur="modifyTitle(item)" />
          </div>
          <client-only>
            <el-popconfirm title="是否移除该条视频？" confirm-button-text="确认" cancel-button-text="取消"
              @confirm="deleteResource(item.id, index)">
              <template #reference>
                <el-icon class="delete-icon" :size="16">
                  <close-icon />
                </el-icon>
              </template>
            </el-popconfirm>
          </client-only>
        </div>
      </el-scrollbar>
    </div>
    <div class="upload-next-btn">
      <el-button type="primary" @click="submitReview">提交</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, nextTick, ref } from "vue";
import VideoUploader from "./VideoUploader.vue";
import { reviewCode } from '@/utils/review-code';
import { ElIcon, ElButton, ElInput, ElPopconfirm, ElScrollbar } from "element-plus";
import CloseIcon from "@/components/icons/CloseIcon.vue";
import { submitReviewAPI, getVideoStatusAPI } from "@/api/video";
import { deleteResourceAPI, modifyTitleAPI } from "@/api/resource";

const emit = defineEmits(["review"]);
const props = defineProps<{
  vid: number,
  resources: Array<ResourceType>
}>();

const resourceList = ref<Array<ResourceType>>(props.resources);
const finishUpload = async () => {
  const res = await getVideoStatusAPI(props.vid);
  if (res.data.code === statusCode.OK) {
    resourceList.value = res.data.data.video.resources;
  } else {
    ElMessage.error(res.data.msg || '提交失败');
  }
}

// 获取标签类型
const toTagType = (state: number) => {
  switch (state) {
    case reviewCode.AUDIT_APPROVED:
      return "success";
    case reviewCode.VIDEO_PROCESSING: case reviewCode.WAITING_REVIEW:
      return "info";
    default:
      return "error";
  }
}

// 获取标签文本
const toTagText = (state: number) => {
  switch (state) {
    case reviewCode.VIDEO_PROCESSING:
      return '处理中';
    case reviewCode.WAITING_REVIEW:
      return '审核中';
    case reviewCode.AUDIT_APPROVED:
      return '审核通过';
    case reviewCode.WRONG_VIDEO_INFO:
      return '视频信息存在问题';
    case reviewCode.WRONG_VIDEO_CONTENT:
      return '视频内容存在问题';
    case reviewCode.PROCESSING_FAIL:
      return '处理失败';
    default:
      return '未知信息';
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
    ElMessage.error('删除失败');
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
const titleClick = (resource: ResourceType, index: number) => {
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
const modifyTitle = async (resource: ResourceType) => {
  modifyIndex.value = -1;
  if (!modifyForm.title) return;
  const res = await modifyTitleAPI(modifyForm);
  if (res.data.code === statusCode.OK) {
    resource.title = modifyForm.title;
  } else {
    ElMessage.error('修改失败');
  }
}
</script>

<style lang="scss" scoped>
.video-box {
  width: 80%;
  margin: 0 auto;
  padding-bottom: 10px;

  .video-item {
    height: 50px;
    padding: 0 20px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border: 1px solid #efeff5;
    border-radius: 6px;
    margin: 10px 20px;

    .part {
      color: var(--primary-color);
    }

    .title-box {
      flex: 1;
      height: 100%;
      display: flex;
      align-items: center;
      padding: 0 12px;

      .title {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
      }
    }

    .delete-icon {
      cursor: pointer;
      color: #767676;

      &:hover {
        color: #929292;
      }
    }
  }
}

.upload-next-btn {
  width: 80%;
  margin: 0 auto;

  button {
    float: right;
    width: 160px;
    height: 40px;
  }
}

:deep(.el-upload-dragger) {
  padding: 0px;
}

:deep(.el-upload-dragger:hover) {
  border-color: var(--primary-color);
}
</style>