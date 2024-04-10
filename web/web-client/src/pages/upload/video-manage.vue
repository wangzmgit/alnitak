<template>
  <div class="upload-video">
    <p class="title">视频管理</p>
    <div class="video-box">
      <el-scrollbar>
        <ul class="video-list" v-infinite-scroll="scrollLoad">
          <li class="video-item" v-for="(item, index) in videoList" :key="index">
            <div class="item-left">
              <div class="cover">
                <img v-if="item.cover" :src="getResourceUrl(item.cover)" alt="收藏夹封面">
              </div>
            </div>
            <div class="item-center">
              <nuxt-link class="item-title" :to="`/video/${item.vid}`">{{ item.title }}</nuxt-link>
              <span class="desc">简介：{{ item.desc }}</span>
              <span class="desc">创建于：{{ formatTime(item.createdAt) }}</span>
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
                    <el-dropdown-item @click="modifyVideo(item.vid, 'info')">编辑信息</el-dropdown-item>
                    <el-dropdown-item @click="modifyVideo(item.vid, 'video')">修改视频</el-dropdown-item>
                    <el-dropdown-item @click="deleteVideo(item.vid)">删除稿件</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </li>
        </ul>
      </el-scrollbar>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from 'vue';
import { getUploadVideoAPI, deleteVideoAPI } from '@/api/video';
import { MoreOne as MoreIcon } from '@icon-park/vue-next';

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

const deleteVideo = async (vid: number) => {
  // TODO: 删除提醒
  const res = await deleteVideoAPI(vid);
  if (res.data.code === statusCode.OK) {
    getUploadVideo();
  }
}

//前往修改视频
const modifyVideo = (vid: number, status: string) => {
  navigateTo({ name: "upload-video", query: { vid: vid, modify: status } });
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
</style>