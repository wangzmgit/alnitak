<template>
  <div class="video-box">
    <p class="video-title">我的视频</p>
    <ul ref="videoListRef" class="video-list">
      <li class="video-item" v-for="item in videoList ">
        <nuxt-link class="cover" :to="item.status === reviewCode.AUDIT_APPROVED ? `/video/${item.vid}` : ''">
          <img class="img" :src="getResourceUrl(item.cover)" />
        </nuxt-link>
        <nuxt-link class="title" :to="item.status === reviewCode.AUDIT_APPROVED ? `/video/${item.vid}` : ''">
          {{ item.title }}
        </nuxt-link>
        <div class="meta">
          <div class="play-count">
            <el-icon size="16" :style="{ marginRight: '4px' }">
              <play-count-icon></play-count-icon>
            </el-icon>
            <span class="clicks">{{ item.clicks }}</span>
          </div>
          <div class="date">{{ formatDate(item.createdAt) }}</div>
        </div>
      </li>
    </ul>
    <div v-show="showPagination" class="pagination">
      <el-pagination background layout="prev, pager, next" :page="page" :page-size="pageSize" :total="total"
        @current-change="pageChange" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { getUploadVideoAPI } from "@/api/video";
import { ElIcon } from 'element-plus';
import PlayCountIcon from "@/components/icons/PlayCountIcon.vue";
import { useVideoCountStore } from '@/composables/video-count-store';

//获取我的视频
const page = ref(1);
const total = ref(0);
const pageSize = ref(8);
const showPagination = ref(false);
const videoCountStore = useVideoCountStore();
const videoList = ref<ManuscriptVideoType[]>([]);
const videoListRef = ref<HTMLElement>();
const getUploadVideo = async () => {
  const res = await getUploadVideoAPI(page.value, pageSize.value);
  if (res.data.code === statusCode.OK) {
    total.value = res.data.data.total;
    videoList.value = res.data.data.videos;
    videoCountStore.setVideoCountState(total.value);

    showPagination.value = total.value > pageSize.value;
  }
}

//页码改变
const pageChange = (target: number) => {
  page.value = target;
  getUploadVideo();
}

const sizeChange = () => {
  if (videoListRef.value!.clientWidth === 1076) {
    pageSize.value = 18;
  } else {
    pageSize.value = 15;
  }
}

onMounted(() => {
  sizeChange();
  getUploadVideo();
  window.addEventListener("resize", sizeChange);
})

onBeforeUnmount(() => {
  window.removeEventListener("resize", sizeChange);
})
</script>

<style lang="scss" scoped>
.video-box {
  padding-left: 12px;

  .video-title {
    font-size: 18px;
    margin-top: 20px;
    padding-left: 10px;
  }

  .pagination {
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

.video-list {
  width: 1076px;
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  padding: 0;

  .video-item {
    display: block;
    width: 160px;
    position: relative;
    margin: 0 0 3px;
    padding: 10px;

    .cover {
      background-size: cover;
      border-radius: 4px;
      display: block;
      width: 160px;
      height: 100px;
      overflow: hidden;
      position: relative;
      cursor: pointer;
      background-color: #f9f9f9;

      .img {
        width: 100%;
        height: 100%;
      }
    }

    .title {
      font-size: 12px;
      color: #222222;
      display: block;
      line-height: 20px;
      height: 38px;
      margin-top: 6px;
      overflow: hidden;
      cursor: pointer;

      &:hover {
        color: var(--primary-hover-color);
      }
    }

    .meta {
      display: flex;
      align-items: center;
      color: #999;
      white-space: nowrap;
      margin-top: 5px;
      height: 16px;
      line-height: 16px;

      .play-count {
        width: 92px;
        display: flex;
        align-items: center;

        .clicks {
          color: #999;
          font-size: 12px;
        }
      }

      .date {
        color: #999;
        font-size: 12px;
      }
    }
  }
}

@media (max-width: 1400px) {
  .video-list {
    width: 876px;
    grid-template-columns: repeat(5, 1fr);
  }
}
</style>