<template>
  <div class="video-box">
    <p class="video-title">我的视频</p>
    <ul ref="videoListRef" class="video-list">
      <li class="video-item" v-for="item in videoList" :key="item.vid">
        <nuxt-link class="cover" :to="item.status === reviewCode.AUDIT_APPROVED ? `/watch?v=${item.shortId || String(item.vid)}` : ''">
          <oss-image class="img" :src="item.cover" alt="封面" />
        </nuxt-link>
        <nuxt-link class="title" :to="item.status === reviewCode.AUDIT_APPROVED ? `/watch?v=${item.shortId || String(item.vid)}` : ''">
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
    <div v-if="loading" class="loading-tip">加载中...</div>
    <div v-if="noMore && videoList.length > 0" class="no-more-tip">没有更多了</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { getUploadVideoAPI } from "@/api/video";
import { ElIcon } from 'element-plus';
import PlayCountIcon from "@/components/icons/PlayCountIcon.vue";
import { useVideoCountStore } from '@/composables/video-count-store';
import { throttle } from "@/utils/debounce";

definePageMeta({
  middleware: ['auth']
})

//获取我的视频
const page = ref(1);
const total = ref(0);
const pageSize = ref(18);
const noMore = ref(false);
const loading = ref(false);
const videoCountStore = useVideoCountStore();
const videoList = ref<ManuscriptVideoType[]>([]);
const videoListRef = ref<HTMLElement>();

const getUploadVideo = async () => {
  if (loading.value || noMore.value) return;
  loading.value = true;
  const res = await getUploadVideoAPI(page.value, pageSize.value);
  if (res.data.code === statusCode.OK) {
    const videos = res.data.data.videos || [];
    if (videos.length > 0) {
      const approvedVideos = videos.filter(
        (video: ManuscriptVideoType) => video.status === reviewCode.AUDIT_APPROVED
      );
      videoList.value = videoList.value.concat(approvedVideos);
      total.value = videoList.value.length;
      videoCountStore.setVideoCountState(total.value);

      if (videos.length < pageSize.value) {
        noMore.value = true;
      }
    } else {
      noMore.value = true;
    }
  }
  loading.value = false;
}

// 监听窗口滚动事件
const handleScroll = () => {
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
  const clientHeight = document.documentElement.clientHeight;
  const scrollHeight = document.documentElement.scrollHeight;

  // 距离底部 100px 时触发加载
  if (scrollTop + clientHeight >= scrollHeight - 100) {
    if (!loading.value && !noMore.value) {
      page.value++;
      getUploadVideo();
    }
  }
  loading.value = false;
}

const throttledScroll = throttle(handleScroll, 150);

onMounted(() => {
  getUploadVideo();
  window.addEventListener('scroll', throttledScroll, true);
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', throttledScroll, true);
})
</script>

<style lang="scss" scoped>
.video-box {
  padding-left: 12px;

  .video-title {
    font-size: 18px;
    margin-top: 20px;
    padding-left: 10px;
    color: var(--font-primary-1);
  }

  .loading-tip,
  .no-more-tip {
    text-align: center;
    padding: 20px 0;
    color: var(--font-primary-3);
    font-size: 14px;
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
      color: var(--font-primary-1);
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
      color: var(--font-primary-3);
      white-space: nowrap;
      margin-top: 5px;
      height: 16px;
      line-height: 16px;

      .play-count {
        width: 92px;
        display: flex;
        align-items: center;

        .clicks {
          color: var(--font-primary-3);
          font-size: 12px;
        }
      }

      .date {
        color: var(--font-primary-3);
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