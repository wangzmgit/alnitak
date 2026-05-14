<template>
  <div class="upload-video">
    <p class="title">合集管理</p>
    <div class="action-bar">
      <el-button type="primary" @click="navigateTo('/upload/playlist-edit')">创建合集</el-button>
    </div>
    <div class="video-box">
      <el-scrollbar>
        <template v-if="!playlistList.length && !loading">
          <el-empty description="暂无合集" />
        </template>
        <ul class="video-list" v-else>
          <li class="video-item" v-for="(item, index) in playlistList" :key="item.id">
            <div class="item-left">
              <div class="cover">
                <img v-if="item.cover" :src="getResourceUrl(item.cover)" alt="封面">
              </div>
            </div>
            <div class="item-center">
              <span class="item-title unlinked">{{ item.title }}</span>
              <span class="desc">简介：{{ item.desc || '暂无' }}</span>
              <div class="desc">
                <span>视频数：{{ item.videoCount }} | 播放：{{ item.views }}</span>
                <span class="status" v-if="getStatusText(item.status)"
                  :style="`color: ${getStatusTextColor(item.status)}`">{{ getStatusText(item.status) }}</span>
                <span class="status status-btn" v-if="item.status === reviewCode.REVIEW_FAILED"
                  @click="showReason(item.id)">查看原因</span>
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
                    <el-dropdown-item @click="editPlaylist(item.id)">编辑</el-dropdown-item>
                    <el-dropdown-item @click="manageVideos(item.id)">管理视频</el-dropdown-item>
                    <el-dropdown-item @click="deletePlaylist(item, index)">删除合集</el-dropdown-item>
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
        <div class="delete-dialog-title">请输入 <strong>{{ deleteTarget?.title }}</strong> 删除此合集</div>
        <div class="delete-dialog-desc">合集删除后将无法恢复，请谨慎操作</div>
        <el-input class="input" v-model="deleteInputTitle" placeholder="请输入合集标题"></el-input>
        <el-button type="danger" class="delete-btn" @click="submitDelete">确认删除</el-button>
      </el-dialog>
    </client-only>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from 'vue';
import { getMyPlaylistAPI, deletePlaylistAPI, getPlaylistReviewRecordAPI } from '@/api/playlist';
import { MoreOne as MoreIcon } from '@icon-park/vue-next';
import { reviewCode } from '@/utils/review-code';
import { statusCode } from '@/utils/status-code';
import { getResourceUrl } from '@/utils/resource';

interface PlaylistItem {
  id: number;
  title: string;
  cover: string;
  desc: string;
  isOpen: boolean;
  status: number;
  videoCount: number;
  views: number;
  favorites: number;
  createdAt: string;
}

const loading = ref(true);
const playlistList = ref<PlaylistItem[]>([]);

const getMyPlaylists = async () => {
  const res = await getMyPlaylistAPI();
  loading.value = false;
  if (res.data.code === statusCode.OK) {
    playlistList.value = res.data.data.playlists || [];
  }
}

const editPlaylist = (id: number) => {
  navigateTo({ path: '/upload/playlist-edit', query: { id } });
}

const manageVideos = (id: number) => {
  navigateTo({ path: '/upload/playlist-video', query: { id } });
}

// 删除相关
const deleteTarget = ref<PlaylistItem>();
const deleteInputTitle = ref("");
const deleteDialogVisible = ref(false);
const deleteTargetIndex = ref(-1);

const deletePlaylist = (item: PlaylistItem, index: number) => {
  deleteTarget.value = item;
  deleteTargetIndex.value = index;
  deleteDialogVisible.value = true;
}

const beforeClose = () => {
  deleteInputTitle.value = "";
  deleteTarget.value = undefined;
  deleteTargetIndex.value = -1;
  deleteDialogVisible.value = false;
}

const submitDelete = async () => {
  if (deleteInputTitle.value === deleteTarget.value?.title) {
    const res = await deletePlaylistAPI(deleteTarget.value.id);
    if (res.data.code === statusCode.OK) {
      playlistList.value.splice(deleteTargetIndex.value, 1);
    }
    beforeClose();
  } else {
    ElMessage.error("输入标题与原标题不一致");
  }
}

const getStatusText = (status: number) => {
  switch (status) {
    case reviewCode.WAITING_REVIEW:
      return "待审核"
    case reviewCode.REVIEW_FAILED:
      return "审核不通过"
  }
}

const getStatusTextColor = (status: number) => {
  switch (status) {
    case reviewCode.WAITING_REVIEW:
      return "var(--primary-hover-color)"
    case reviewCode.REVIEW_FAILED:
      return "#f56c6c"
  }
}

const showReason = async (id: number) => {
  const res = await getPlaylistReviewRecordAPI(id);
  if (res.data.code === statusCode.OK) {
    const review = res.data.data.review;
    ElMessageBox.alert(review.remark || '审核未通过', '审核不通过原因', {
      confirmButtonText: '确认',
    })
  }
}

onBeforeMount(() => {
  getMyPlaylists();
})
</script>

<style lang="scss" scoped>
.upload-video {
  padding: 0 18px 0;
  height: 100%;
  box-sizing: border-box;
  background-color: var(--bg-elev-1);

  .title {
    font-size: 18px;
    margin: 0;
    padding: 16px 0 10px;
  }

  .action-bar {
    padding-bottom: 12px;
  }

  .video-box {
    height: calc(100% - 100px);
  }

  .video-list {
    list-style: none;
    box-sizing: border-box;
    width: 100%;
    margin: 0;
    padding: 0;

    .empty-tip {
      text-align: center;
      color: var(--font-primary-3);
      padding: 40px 0;
    }

    .video-item {
      display: flex;
      padding: 16px 0;
      width: 100%;
      height: 80px;
      margin-bottom: 12px;
      border-bottom: 1px solid var(--border-color);
      padding-bottom: 12px;

      .item-left {
        width: 120px;
        height: 80px;
        margin-right: 10px;

        .cover {
          border-radius: 2px;
          width: 100%;
          height: 100%;
          background-color: var(--fill-1, #f1f2f3);

          img {
            width: 100%;
            height: 100%;
            border-radius: 2px;
            object-fit: contain;
          }
        }
      }

      .item-center {
        flex: 1;

        .item-title {
          font-size: 14px;
          color: var(--font-primary-1);
          line-height: 18px;
          margin: 0 0 26px;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 1;
          line-clamp: 1;
          -webkit-box-orient: vertical;

          &.unlinked {
            color: var(--font-primary-1);
          }
        }

        .desc {
          font-size: 12px;
          color: var(--font-primary-3);
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 1;
          line-clamp: 1;
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
    color: var(--font-primary-1);
    text-align: center;
    margin: 20px 0;
  }

  .delete-dialog-desc {
    color: var(--font-primary-3);
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
