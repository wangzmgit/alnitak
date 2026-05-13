<template>
  <div class="playlist-video">
    <div class="header">
      <p class="title">合集视频管理</p>
      <div class="header-actions">
        <el-button @click="navigateTo('/upload/playlist-manage')">返回</el-button>
        <el-button type="primary" @click="showAddDialog">添加视频</el-button>
      </div>
    </div>
    <div class="video-box">
      <el-scrollbar>
        <div v-if="videoList.length === 0 && !loading" class="empty-tip">暂无视频，点击"添加视频"按钮添加</div>
        <ul v-else class="video-list">
          <li v-for="(element, index) in videoList" :key="element.vid" class="video-item"
            draggable="true"
            :class="{ 'drag-over': dragOverIndex === index, 'dragging': dragIndex === index }"
            @dragstart="onDragStart(index, $event)"
            @dragover.prevent="onDragOver(index)"
            @dragleave="onDragLeave(index)"
            @drop.prevent="onDrop(index)"
            @dragend="onDragEnd">
            <div class="drag-handle">⠿</div>
            <div class="item-index">P{{ index + 1 }}</div>
            <div class="item-cover">
              <img v-if="element.cover" :src="getResourceUrl(element.cover)" alt="封面">
            </div>
            <div class="item-info">
              <span class="item-title">{{ element.title }}</span>
              <span class="item-desc">{{ element.desc }}</span>
            </div>
            <div class="item-actions">
              <el-button type="danger" text size="small" @click="removeVideo(element.vid, index)">移除</el-button>
            </div>
          </li>
        </ul>
      </el-scrollbar>
    </div>

    <!-- 添加视频弹窗 -->
    <client-only>
      <el-dialog v-model="addDialogVisible" title="添加视频到合集" width="600">
        <div class="add-video-list">
          <div v-if="myVideos.length === 0" class="empty-tip">暂无可添加的视频</div>
          <div v-for="v in myVideos" :key="v.vid" class="add-video-item" :class="{ 'is-disabled': v.inOtherPlaylist }">
            <el-checkbox v-model="v.checked" :disabled="v.inPlaylist || v.inOtherPlaylist"></el-checkbox>
            <div class="add-cover">
              <img v-if="v.cover" :src="getResourceUrl(v.cover)" alt="封面">
            </div>
            <div class="add-info">
              <span class="add-title">{{ v.title }}</span>
              <span v-if="v.inPlaylist" class="in-playlist">已在合集中</span>
              <span v-else-if="v.inOtherPlaylist" class="in-other-playlist">已在其他合集中</span>
            </div>
          </div>
        </div>
        <template #footer>
          <el-button @click="addDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmAdd" :loading="adding">确认添加</el-button>
        </template>
      </el-dialog>
    </client-only>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount } from 'vue';
import {
  getPlaylistVideoListAPI,
  addPlaylistVideoAPI,
  delPlaylistVideoAPI,
  sortPlaylistVideoAPI,
  getMyPlaylistVideoIdsAPI,
} from '@/api/playlist';
import { getAllVideoAPI } from '@/api/video';
import { statusCode } from '@/utils/status-code';
import { getResourceUrl } from '@/utils/resource';

const route = useRoute();
const playlistId = ref(0);

interface PlaylistVideoItem {
  vid: number;
  title: string;
  cover: string;
  desc: string;
  duration: number;
  clicks: number;
  createdAt: string;
}

interface MyVideoItem {
  vid: number;
  title: string;
  cover: string;
  checked: boolean;
  inPlaylist: boolean;
  inOtherPlaylist: boolean; // 是否已在其他合集中
}

const loading = ref(true);
const videoList = ref<PlaylistVideoItem[]>([]);

const getVideoList = async () => {
  const res = await getPlaylistVideoListAPI(playlistId.value, 1, 200);
  loading.value = false;
  if (res.data.code === statusCode.OK) {
    videoList.value = (res.data.data.videos || []).map((v: any) => ({
      vid: Number(v.shortId) || 0,
      title: v.title,
      cover: v.cover,
      desc: v.desc,
      duration: v.duration,
      clicks: v.clicks,
      createdAt: v.createdAt,
    }));
  }
}

const saveSort = async () => {
  const vids = videoList.value.map(v => v.vid);
  await sortPlaylistVideoAPI({ playlistId: playlistId.value, vids });
}

// 拖拽排序
const dragIndex = ref(-1);
const dragOverIndex = ref(-1);

const onDragStart = (index: number, e: DragEvent) => {
  dragIndex.value = index;
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move';
  }
}

const onDragOver = (index: number) => {
  if (dragIndex.value === index) return;
  dragOverIndex.value = index;
}

const onDragLeave = (index: number) => {
  if (dragOverIndex.value === index) {
    dragOverIndex.value = -1;
  }
}

const onDrop = (index: number) => {
  if (dragIndex.value < 0 || dragIndex.value === index) return;
  const list = [...videoList.value];
  const [item] = list.splice(dragIndex.value, 1);
  list.splice(index, 0, item);
  videoList.value = list;
  dragOverIndex.value = -1;
  dragIndex.value = -1;
  saveSort();
}

const onDragEnd = () => {
  dragIndex.value = -1;
  dragOverIndex.value = -1;
}

const removeVideo = async (vid: number, index: number) => {
  const res = await delPlaylistVideoAPI({ playlistId: playlistId.value, vids: [vid] });
  if (res.data.code === statusCode.OK) {
    videoList.value.splice(index, 1);
    ElMessage.success('已移除');
  }
}

// 添加视频
const addDialogVisible = ref(false);
const adding = ref(false);
const myVideos = ref<MyVideoItem[]>([]);

const showAddDialog = async () => {
  const [videoRes, mapRes] = await Promise.all([
    getAllVideoAPI(),
    getMyPlaylistVideoIdsAPI(),
  ]);
  if (videoRes.data.code === statusCode.OK) {
    const videos = videoRes.data.data.videos || [];
    const existingVids = new Set(videoList.value.map(v => v.vid));
    // videoPlaylistMap: vid -> playlistId（该视频已在哪个合集中）
    const videoPlaylistMap: Record<number, number> = mapRes.data.code === statusCode.OK
      ? (mapRes.data.data.videoPlaylistMap || {})
      : {};
    myVideos.value = videos.map((v: any) => {
      const vid = Number(v.shortId) || 0;
      const belongsToPlaylistId = videoPlaylistMap[vid];
      const inCurrentPlaylist = existingVids.has(vid);
      // 已在其他合集中（不是当前合集）
      const inOtherPlaylist = belongsToPlaylistId !== undefined && belongsToPlaylistId !== playlistId.value && !inCurrentPlaylist;
      return {
        vid,
        title: v.title,
        cover: v.cover,
        checked: false,
        inPlaylist: inCurrentPlaylist,
        inOtherPlaylist,
      };
    });
  }
  addDialogVisible.value = true;
}

const confirmAdd = async () => {
  const selectedVids = myVideos.value.filter(v => v.checked && !v.inPlaylist).map(v => v.vid);
  if (selectedVids.length === 0) {
    ElMessage.warning('请选择要添加的视频');
    return;
  }

  adding.value = true;
  try {
    const res = await addPlaylistVideoAPI({ playlistId: playlistId.value, vids: selectedVids });
    if (res.data.code === statusCode.OK) {
      ElMessage.success('添加成功');
      addDialogVisible.value = false;
      getVideoList();
    }
  } finally {
    adding.value = false;
  }
}

onBeforeMount(() => {
  const id = Number(route.query.id);
  if (id) {
    playlistId.value = id;
    getVideoList();
  } else {
    navigateTo('/upload/playlist-manage');
  }
})
</script>

<style lang="scss" scoped>
.playlist-video {
  padding: 0 18px;
  height: 100%;
  box-sizing: border-box;
  background-color: var(--bg-elev-1);

  .header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 0 12px;

    .title {
      font-size: 18px;
      margin: 0;
    }
  }

  .video-box {
    height: calc(100% - 70px);
  }

  .empty-tip {
    text-align: center;
    color: var(--font-primary-3);
    padding: 40px 0;
  }

  .video-list {
    list-style: none;
    margin: 0;
    padding: 0;

    .video-item {
      display: flex;
      align-items: center;
      padding: 10px 0;
      border-bottom: 1px solid var(--border-color);
      transition: background-color 0.15s;

      &.dragging {
        opacity: 0.4;
      }

      &.drag-over {
        background-color: var(--primary-color-light, rgba(64, 158, 255, 0.08));
        border-top: 2px solid var(--primary-color);
      }

      .drag-handle {
        cursor: grab;
        padding: 0 8px;
        font-size: 18px;
        color: var(--font-primary-3);
        user-select: none;

        &:active {
          cursor: grabbing;
        }
      }

      .item-index {
        width: 30px;
        font-size: 13px;
        color: var(--font-primary-3);
        text-align: center;
      }

      .item-cover {
        width: 100px;
        height: 60px;
        margin: 0 10px;
        background-color: var(--fill-1, #f1f2f3);
        border-radius: 2px;

        img {
          width: 100%;
          height: 100%;
          object-fit: contain;
          border-radius: 2px;
        }
      }

      .item-info {
        flex: 1;
        display: flex;
        flex-direction: column;

        .item-title {
          font-size: 14px;
          color: var(--font-primary-1);
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .item-desc {
          font-size: 12px;
          color: var(--font-primary-3);
          margin-top: 4px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }

      .item-actions {
        padding: 0 8px;
      }
    }
  }
}

.add-video-list {
  max-height: 400px;
  overflow-y: auto;

  .empty-tip {
    text-align: center;
    color: var(--font-primary-3);
    padding: 20px;
  }

  .add-video-item {
    display: flex;
    align-items: center;
    padding: 8px 0;
    border-bottom: 1px solid var(--border-color);

    &.is-disabled {
      opacity: 0.5;
    }

    .add-cover {
      width: 80px;
      height: 48px;
      margin: 0 10px;
      background-color: var(--fill-1, #f1f2f3);
      border-radius: 2px;

      img {
        width: 100%;
        height: 100%;
        object-fit: contain;
        border-radius: 2px;
      }
    }

    .add-info {
      flex: 1;

      .add-title {
        font-size: 13px;
        color: var(--font-primary-1);
      }

      .in-playlist {
        font-size: 12px;
        color: var(--font-primary-3);
        margin-left: 8px;
      }

      .in-other-playlist {
        font-size: 12px;
        color: #E6A23C;
        margin-left: 8px;
      }
    }
  }
}
</style>
