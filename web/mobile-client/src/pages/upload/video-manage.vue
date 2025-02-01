<template>
  <div class="upload-video">
    <p class="title">视频管理</p>
    <div class="video-box">
      <el-scrollbar>
        <ul class="video-list" v-infinite-scroll="scrollLoad">
          <li class="video-item" v-for="(item, index) in videoList" :key="index">
            <div class="item-left">
              <div class="cover">
                <img v-if="item.cover" :src="getResourceUrl(item.cover)" alt="封面">
              </div>
            </div>
            <div class="item-center">
              <nuxt-link class="item-title" :to="`/video/${item.vid}`">{{ item.title }}</nuxt-link>
              <span class="desc">简介：{{ item.desc }}</span>
              <div class="desc">
                <span>创建于：{{ formatTime(item.createdAt) }}</span>
                <span class="status" v-if="getStatusText(item.status)"
                  :style="`color: ${getStatusTextColor(item.status)}`">{{ getStatusText(item.status) }}</span>
                <span class="status status-btn" v-if="item.status === reviewCode.REVIEW_FAILED"
                  @click="showReason(item.vid)">查看原因</span>
              </div>
            </div>
            <div class="item-right">
              <el-dropdown>
                <el-button class="more-btn" plain>
                  <el-icon size="16">
                    <more-icon></more-icon>
                  </el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="modifyVideo(item.vid)">编辑</el-dropdown-item>
                    <el-dropdown-item @click="deleteVideo(item, index)">删除稿件</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </li>
        </ul>
      </el-scrollbar>
    </div>
    <client-only>
      <el-dialog v-model="deleteDialogVisible" class="delete-dialog" width="500" :before-close="beforeClose">
        <div class="delete-dialog-title">请输入 <strong>{{ deleteVideoInfo?.title }}</strong> 删除此视频</div>
        <div class="delete-dialog-desc">视频删除后将无法恢复，请谨慎操作</div>
        <el-input class="input" v-model="deleteVideoTitle" placeholder="请输入视频标题"></el-input>
        <el-button type="danger" class="delete-btn" @click="submitDelete">确认删除</el-button>
      </el-dialog>
    </client-only>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from 'vue';
import { getUploadVideoAPI, deleteVideoAPI } from '@/api/video';
import { MoreOne as MoreIcon } from '@icon-park/vue-next';
import { getVideoReviewRecordAPI } from '@/api/revies';

const page = ref(1);
const total = ref(0);
const pageSize = 8;
const noMore = ref(false);
const loading = ref(false);
const videoList = ref<Array<VideoType>>([]);
const getUploadVideo = async () => {
  if (loading.value || noMore.value) return;
  loading.value = true;
  const res = await getUploadVideoAPI(page.value, pageSize);
  if (res.data.code === statusCode.OK) {
    total.value = res.data.data.total;
    if (res.data.data.videos) {
      videoList.value = videoList.value.concat(res.data.data.videos);
    } else {
      noMore.value = true;
    }
  }
  loading.value = false;
}

const scrollLoad = () => {
  if (!loading.value) {
    page.value++;
    getUploadVideo();
  }
}

const deleteVideoIndex = ref(-1);
const deleteVideoTitle = ref("");
const deleteDialogVisible = ref(false);
const deleteVideoInfo = ref<VideoType>();
const deleteVideo = async (video: VideoType, index: number) => {
  deleteVideoInfo.value = video;
  deleteVideoIndex.value = index;
  deleteDialogVisible.value = true;
}

const beforeClose = () => {
  deleteVideoTitle.value = "";
  deleteVideoIndex.value = -1;
  deleteVideoInfo.value = undefined;
  deleteDialogVisible.value = false;
}

const submitDelete = async () => {
  if (deleteVideoTitle.value === deleteVideoInfo.value?.title) {
    const res = await deleteVideoAPI(deleteVideoInfo.value.vid);
    if (res.data.code === statusCode.OK) {
      videoList.value.splice(deleteVideoIndex.value, 1);
    }

    deleteVideoTitle.value = "";
    deleteVideoIndex.value = -1;
    deleteVideoInfo.value = undefined;
    deleteDialogVisible.value = false;
  } else {
    ElMessage.error("输入标题与原标题不一致");
  }
}

const getStatusText = (status: number) => {
  switch (status) {
    // case reviewCode.CREATED_VIDEO:
    //   return "未提交"
    case reviewCode.SUBMIT_REVIEW:
    case reviewCode.WAITING_REVIEW:
      return "审核中"
    case reviewCode.REVIEW_FAILED:
      return "审核不通过"
    case reviewCode.PROCESSING_FAIL:
      return "视频处理失败"
  }
}

const getStatusTextColor = (status: number) => {
  switch (status) {
    case reviewCode.CREATED_VIDEO:
      return "#999"
    case reviewCode.SUBMIT_REVIEW:
    case reviewCode.WAITING_REVIEW:
      return "var(--primary-hover-color)"
    case reviewCode.REVIEW_FAILED:
      return "#f56c6c"
    case reviewCode.PROCESSING_FAIL:
      return "#f56c6c"
  }
}

const showReason = async (vid: number) => {
  const res = await getVideoReviewRecordAPI(vid);
  if (res.data.code === statusCode.OK) {
    ElMessageBox.alert(res.data.data.review.remark, '', {
      confirmButtonText: '确认',
    })
  }
}

//前往修改视频
const modifyVideo = (vid: number) => {
  navigateTo({ name: "upload-video", query: { vid: vid } });
}

onBeforeMount(() => {
  getUploadVideo();
})
</script>

<style lang="scss" scoped>
.upload-video {
  padding: 0 18px 0;
  height: 100%;
  box-sizing: border-box;
  background-color: #fff;

  .title {
    font-size: 18px;
    margin: 0;
    padding: 16px 0 10px;
  }

  .video-box {
    height: calc(100% - 60px);
  }

  .video-list {
    list-style: none;
    box-sizing: border-box;
    width: 100%;
    margin: 0;
    padding: 0;

    .video-item {
      display: flex;
      padding: 16px 0;
      width: 100%;
      height: 80px;
      margin-bottom: 12px;
      border-bottom: 1px solid #e6e6e6;
      padding-bottom: 12px;

      .item-left {
        width: 120px;
        height: 80px;
        margin-right: 10px;

        .cover {
          border-radius: 2px;
          width: 100%;
          height: 100%;
          background-color: #f1f2f3;

          img {
            width: 100%;
            height: 100%;
            border-radius: 2px;
          }
        }
      }

      .item-center {
        flex: 1;

        .item-title {
          font-size: 14px;
          color: #212121;
          line-height: 18px;
          margin: 0 0 26px;
          cursor: pointer;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 1;
          -webkit-box-orient: vertical;

          &:hover {
            color: var(--primary-hover-color);
          }
        }

        .desc {
          font-size: 12px;
          color: #999;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 1;
          -webkit-box-orient: vertical;
        }

        .status {
          margin-left: 12px;
          color: var(--primary-hover-color);
        }

        .status-btn {
          cursor: pointer;
        }
      }

      .item-right {
        width: 90px;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }
  }
}

.delete-dialog {

  .delete-dialog-title {
    font-size: 16px;
    color: #1f2328;
    text-align: center;
    margin: 20px 0;
  }

  .delete-dialog-desc {
    color: #666;
    font-size: 13px;
    text-align: center;

  }

  .input {
    margin: 20px 0;

  }

  .delete-btn {
    width: 100%;
    color: #d03050;
    border: none;
    font-family: inherit;
    background-color: rgba(208, 48, 80, 0.16);

    &:hover {
      background-color: rgba(208, 48, 80, 0.22);
    }
  }
}
</style>