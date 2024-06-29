<template>
  <div class="home">
    <home-header class="home-header" @change-fold="changeMenuFold"></home-header>
    <div class="home-content">
      <div class="home-left" :class="menuFold ? 'home-left-fold' : ''">
        <home-sidebar class="home-sidebar" :fold="menuFold"></home-sidebar>
      </div>
      <div class="home-right" :style="`margin-left: ${menuFold ? '50px' : '220px'};`">
        <div class="home-recommended" :class="menuFold ? 'recommended-fold' : ''">
          <div class="recommended-carousel">
            <client-only>
              <HomeCarousel :partition-id="partitionId"></HomeCarousel>
            </client-only>
          </div>
          <video-item v-for="item in videoList" :info="item"></video-item>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import { onMounted, onBeforeUnmount, ref } from 'vue';
import VideoItem from '@/components/home-video-item/index.vue';
import HomeSidebar from '@/components/home-sidebar/index.vue';
import HomeHeader from "@/components/home-header/index.vue";
import HomeCarousel from '@/components/alnitak-carousel/index.vue';
import { asyncGetVideoByPartitionAPI, getVideoByPartitionAPI } from "@/api/video";

const route = useRoute();

const partitionId = route.params.partition.toString();

const menuFold = ref(false);
const changeMenuFold = (val: boolean) => {
  menuFold.value = val;
}

// 获取分区
const size = 10;
const videoList = ref<VideoType[]>([])
const { data } = await asyncGetVideoByPartitionAPI(size, partitionId);
if ((data.value as any).code === statusCode.OK) {
  videoList.value = (data.value as any).data.videos;
}

const noMore = ref(false);
const loading = ref(false);
const getViedeoList = async () => {
  loading.value = true;
  const res = await getVideoByPartitionAPI(size, partitionId);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.videos) {
      videoList.value = videoList.value.concat(res.data.data.videos);
    }
  }
  loading.value = false;
}

const lazyLoading = (e: Event) => {
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
  if (scrollTop === 0) return;

  const clientHeight = document.documentElement.clientHeight;
  const scrollHeight = document.documentElement.scrollHeight;
  if (scrollTop + clientHeight >= scrollHeight) {
    if (!loading.value) {
      getViedeoList();
    }
  }
}

onMounted(() => {
  window.addEventListener('scroll', lazyLoading, true);
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', lazyLoading, true);
})
</script>

<style lang="scss" scoped>
.home {
  width: 100%;
  min-width: 1200px;
  overflow: hidden;

  .home-header {
    position: fixed;
    top: 0;
    width: 100%;
    z-index: 999;
    background-color: #fff;
  }
}

.home-content {
  display: flex;
  margin-top: 60px;

  .home-left {
    position: fixed;
    height: 100%;
    width: 220px;
    z-index: 1;
    background-color: #fff;
    transition: width .25s;
  }

  .home-left-fold {
    width: 50px;
  }

  .home-right {
    flex: 1;
    margin-top: 6px;
  }
}

.home-recommended {
  display: grid;
  margin-left: 20px;
  width: calc(100% - 50px);
  gap: 0 16px;
  grid-template-columns: repeat(4, 1fr);
  overflow: hidden;

  .recommended-carousel {
    height: 420px;
    grid-row: 1/ span 2;
    grid-column: 1/ span 2;
  }
}

.recommended-fold {
  max-height: 780px;
  grid-template-columns: repeat(5, 1fr);
}
</style>