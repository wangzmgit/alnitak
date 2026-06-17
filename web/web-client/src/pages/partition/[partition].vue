<template>
  <div class="home">
    <home-header class="home-header" @change-fold="changeMenuFold"></home-header>
    <div class="home-content">
      <div class="home-left" :class="menuFold ? 'home-left-fold' : ''">
        <home-sidebar class="home-sidebar" :fold="menuFold"></home-sidebar>
      </div>
      <div class="home-right" :style="`margin-left: ${menuFold ? '50px' : '220px'};`">
        <div class="home-recommended" :class="menuFold ? 'recommended-fold' : ''">
          <div v-if="carouselList.length" class="recommended-carousel">
            <div class="recommended-carousel-inner">
              <AlnitakCarousel :list="carouselList"></AlnitakCarousel>
            </div>
          </div>
          <video-item v-for="item in videoList" :key="item.vid" :info="item"></video-item>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import { onMounted, onBeforeUnmount, ref } from 'vue';
import VideoItem from '@/components/home-video-item/index.vue';
import { asyncGetVideoByPartitionAPI, getVideoByPartitionAPI } from "@/api/video";
import { getCarouselAPI } from "@/api/carousel";
import { throttle } from "@/utils/debounce";

const route = useRoute();

const partitionId = route.params.partition.toString();

const menuFoldCookie = useCookie<boolean>('menu-fold-state', { default: () => false });
const menuFold = ref(menuFoldCookie.value);
const changeMenuFold = (val: boolean) => {
  menuFold.value = val;
}

// 客户端挂载后同步折叠状态
onMounted(() => {
  try {
    const saved = localStorage.getItem('menu-fold-state');
    if (saved === 'true') {
      menuFold.value = true;
    }
  } catch {}
});

// 获取分区
const size = 10;
const videoList = ref<VideoType[]>([])
const [partitionResult, carouselResult] = await Promise.all([
  asyncGetVideoByPartitionAPI(size, partitionId),
  getCarouselAPI(partitionId),
]);
const { data } = partitionResult;
if ((data.value as any).code === statusCode.OK) {
  videoList.value = (data.value as any).data.videos;
}
const carouselList = ref<CarouselType[]>([]);
if (carouselResult.data.code === statusCode.OK && carouselResult.data.data.carousels) {
  carouselList.value = carouselResult.data.data.carousels;
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

const throttledLoading = throttle(lazyLoading, 150);

onMounted(() => {
  window.addEventListener('scroll', throttledLoading, true);
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', throttledLoading, true);
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
    background-color: var(--bg-elev-1);
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
    background-color: var(--bg-elev-1);
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
    /* 行业常用：响应式高度 + 上限，避免轮播过高 */
    min-height: 0;
    grid-row: 1/ span 2;
    grid-column: 1/ span 2;
    align-self: stretch;
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .recommended-carousel-inner {
    height: clamp(220px, 26vw, 890px);
    position: relative;
    border-radius: 9px;
    overflow: hidden;
  }

  .recommended-carousel-inner :deep(.carousel-area) {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
  }
}

.recommended-fold {
  grid-template-columns: repeat(5, 1fr);
}
</style>