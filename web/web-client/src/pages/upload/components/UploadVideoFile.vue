<template>
  <div class="title-bar">
    <h2 class="title-text">文件上传</h2>
    <el-upload :show-file-list="false" :before-upload="beforeUploadVideo" @change="handleChange">
      <el-button type="primary" :icon="Plus">添加视频</el-button>
    </el-upload>
  </div>
  <div class="upload-video">
    <draggable v-model="resourceList" item-key="id" handle=".drag-handle" @end="onDragEnd"
      :animation="200" ghost-class="drag-ghost">
      <template #item="{ element: item, index }">
        <div class="video-item">
          <div class="drag-handle" title="拖拽排序">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <circle cx="5" cy="3" r="1.5"/><circle cx="11" cy="3" r="1.5"/>
              <circle cx="5" cy="8" r="1.5"/><circle cx="11" cy="8" r="1.5"/>
              <circle cx="5" cy="13" r="1.5"/><circle cx="11" cy="13" r="1.5"/>
            </svg>
          </div>
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
                <el-input v-else ref="titleInput" v-model="modifyForm.title" maxlength="100" show-word-limit
                  @blur="modifyTitle(item)" />
              </div>
              <div class="action-btns">
                <el-upload :show-file-list="false" :before-upload="beforeUploadVideo"
                  @change="(file: any) => handleReplace(file, item, index)">
                  <span class="replace-btn" v-if="!item.uploading">替换</span>
                </el-upload>
                <client-only>
                  <el-popconfirm v-if="resourceList.length > 1" title="是否移除该条视频？" confirm-button-text="确认" cancel-button-text="取消"
                    @confirm="deleteResource(item.id, index)">
                    <template #reference>
                      <span class="remove-btn">移除</span>
                    </template>
                  </el-popconfirm>
                </client-only>
              </div>
            </div>
            <div class="progress-box">
              <span class="upload-status">{{ item.uploading ? `上传中 ${item.percent}%` : getTagText(item.status) }}</span>
              <div class="progress-bar">
                <div class="progress" :style="`width: ${item.uploading ? item.percent : 100}%`"></div>
              </div>
            </div>
            <div v-if="canManageSubtitles(item)" class="subtitle-action">
              <el-button link type="primary" size="small" @click="openSubtitleDialog(subtitleResourceKey(item))">字幕管理</el-button>
            </div>
          </div>
        </div>
      </template>
    </draggable>
  </div>
  <el-dialog v-model="subtitleDialogVisible" title="字幕管理" width="600" :destroy-on-close="true">
    <resource-subtitle-editor
      v-if="subtitleDialogResourceId"
      :resource-short-id="subtitleDialogResourceId"
      :key="subtitleDialogResourceId"
    />
  </el-dialog>
</template>

<script setup lang="ts">
import { reactive, nextTick, ref, watch } from "vue";
import { Plus } from '@icon-park/vue-next';
import { reviewCode } from '@/utils/review-code';
import { ElIcon, ElButton, ElInput, ElPopconfirm } from "element-plus";
import MonitorIcon from "@/components/icons/MonitorIcon.vue";
import draggable from "vuedraggable";
import { submitReviewAPI, getVideoStatusAPI } from "@/api/video";
import { deleteResourceAPI, modifyTitleAPI, replaceResourceAPI, checkReplaceResourceAPI, reorderResourceAPI } from "@/api/resource";
import { uploadFileChunkAPI } from "@/api/upload";
import { getFileMD5 } from '@/utils/md5';
import ResourceSubtitleEditor from './ResourceSubtitleEditor.vue';

const subtitleDialogVisible = ref(false);
const subtitleDialogResourceId = ref('');

const openSubtitleDialog = (resourceId: string) => {
  subtitleDialogResourceId.value = resourceId;
  subtitleDialogVisible.value = true;
};

const emit = defineEmits(["review"]);
const props = defineProps<{
  vid: number,
  resources: Array<ResourceType> | null
}>();

const resourceList = ref<Array<ResourceType | UploadResourceType>>(props.resources ?? []);

/** 分 P 已有正式 ID 且未在上传中时可管理字幕 */
const canManageSubtitles = (item: ResourceType | UploadResourceType) => {
  if (item.id <= 0) return false;
  if ('uploading' in item && item.uploading) return false;
  return true;
};

const subtitleResourceKey = (item: ResourceType | UploadResourceType): string => {
  const r = item as ResourceType;
  if (r.shortId) return String(r.shortId);
  return String(item.id);
};

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
const titleInput = ref<Array<InstanceType<typeof ElInput>>>([]);
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
    const inst = titleInput.value?.[index];
    // element-plus ElInput 实例有 focus 方法，但在渲染/拖拽重排时可能短暂为 undefined
    inst?.focus?.();
  });
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

// 上传临时ID计数器，用于唯一标识每个上传中的项
let uploadIdCounter = 0;

//上传变化的回调
const handleChange = (uploadFile: any) => {
  if (!uploadFile.raw) return;
  if (!Array.isArray(resourceList.value)) {
    resourceList.value = [];
  }

  const uploadKey = --uploadIdCounter; // 负数ID，避免和真实资源ID冲突

  const uploadData: UploadResourceType = {
    id: uploadKey,
    status: -1,
    title: "",
    percent: 0,
    uploading: true,
  }

  resourceList.value.push(uploadData);

  const findIndex = () => resourceList.value.findIndex(r => r.id === uploadKey);

  uploadFileChunkAPI({
    name: "video",
    action: props.vid ? `v1/upload/video/${props.vid}` : `v1/upload/video`,
    file: uploadFile.raw,
    onProgress: (val: any) => {
      const idx = findIndex();
      if (idx === -1) return;
      uploadData.percent = val;
      resourceList.value[idx] = JSON.parse(JSON.stringify(uploadData));
    },
    onError: () => {
      const idx = findIndex();
      if (idx !== -1) resourceList.value.splice(idx, 1);
    },
    onFinish: (data?: any) => {
      const idx = findIndex();
      if (idx !== -1) resourceList.value[idx] = data.data.resource;
    },
  })
}

// 替换资源的回调
const handleReplace = async (uploadFile: any, resource: ResourceType | UploadResourceType, index: number) => {
  if (!uploadFile.raw) return;

  // 保存原始资源信息
  const originalResource = { ...resource };
  const resourceId = resource.id;

  // 先计算新文件的hash
  const hash = await getFileMD5(uploadFile.raw);

  // 先检查hash是否相同
  const checkRes = await checkReplaceResourceAPI(resourceId, hash);
  if (checkRes.data.code !== statusCode.OK) {
    // hash相同或其他错误，无需上传
    ElMessage.info(checkRes.data.msg || '无需替换');
    return;
  }

  // hash不同，需要上传新文件
  // 设置上传状态
  const uploadData: UploadResourceType = {
    id: resourceId,
    status: -1,
    title: originalResource.title,
    percent: 0,
    uploading: true,
  }
  resourceList.value[index] = uploadData;

  uploadFileChunkAPI({
    name: "video",
    action: `v1/resource/replaceResource?resourceId=${resourceId}`,
    file: uploadFile.raw,
    onProgress: (val: any) => {
      uploadData.percent = val;
      resourceList.value[index] = JSON.parse(JSON.stringify(uploadData));
    },
    onError: () => {
      // 恢复原始资源
      resourceList.value[index] = originalResource;
      ElMessage.error('上传失败');
    },
    onFinish: async (data?: any) => {
      if (data.code === statusCode.OK) {
        resourceList.value[index] = data.data.resource;
        ElMessage.success('替换成功，正在转码中');
      } else {
        resourceList.value[index] = originalResource;
        ElMessage.error(data.msg || '替换失败');
      }
    },
  })
}

// 拖拽排序结束
const onDragEnd = async () => {
  const ids = resourceList.value.filter(r => r.id > 0).map(r => r.id);
  if (ids.length === 0 || !props.vid) return;
  const res = await reorderResourceAPI(props.vid, ids);
  if (res.data.code !== statusCode.OK) {
    ElMessage.error(res.data.msg || '排序失败');
  }
}

watch(() => props.resources, (newVal) => {
  resourceList.value = newVal ?? [];
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
    color: var(--font-primary-1);
    font-weight: 600;
    line-height: 50px;
  }
}

.upload-video {
  width: 100%;
  margin: 0 auto;
  padding: 0 20px;
  box-sizing: border-box;

  .drag-ghost {
    opacity: 0.4;
    background-color: var(--hover-bg);
    border-radius: 4px;
  }

  .video-item {
    display: flex;
    align-items: center;
    width: 100%;
    padding: 20px 0;

    &:first-child {
      padding: 10px 0 20px 0;
    }

    &:hover {
      .drag-handle {
        opacity: 1;
      }

      .info-box {
        .file-info {
          .replace-btn,
          .remove-btn {
            display: block;
          }
        }
      }
    }

    .drag-handle {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 20px;
      cursor: grab;
      color: var(--font-primary-3);
      opacity: 0;
      transition: opacity 0.2s;
      flex-shrink: 0;

      &:active {
        cursor: grabbing;
      }

      &:hover {
        color: var(--font-primary-2);
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

        .action-btns {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        .replace-btn,
        .remove-btn {
          display: none;
          font-size: 12px;
          color: var(--font-primary-3);
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
          color: var(--font-primary-3);
          padding-bottom: 2px;
          font-size: 12px;
        }

        .progress-bar {
          position: relative;
          width: 100%;
          height: 4px;
          background-color: var(--fill-1, #efeff5);

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