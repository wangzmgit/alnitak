<template>
  <div v-title :data-title="`${globalConfig.title} - 首页`">
    <header-bar class="header"></header-bar>
    <div class="partition">
      <nav-bar @partition-change="partitionChange"></nav-bar>
    </div>
    <div ref="scrollBox" class="content" @scroll="scrollList">
      <!-- 只有当前分区为首页时显示轮播图 -->
      <div v-if="currentPartition === 0" class="carousel">
        <home-carousel></home-carousel>
      </div>
      <video-list :videos="videoList"></video-list>
      <div v-show="loading" class="loading">
        <n-spin></n-spin>
      </div>
      <footer-bar></footer-bar>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount } from 'vue';
import HomeCarousel from './components/HomeCarousel.vue';
import HeaderBar from '@/components/header-bar/index.vue';
import NavBar from './components/NavBar.vue';
import FooterBar from '@/components/footer-bar/index.vue';
import VideoList from '@/components/video-list/index.vue';
import { NSpin, useMessage } from 'naive-ui';
import { globalConfig } from "@/utils/global-config";
import { getHotVideoAPI, getVideoByPartitionAPI } from "@/api/video";
import { statusCode } from '@/utils/status-code';

const message = useMessage();

//视频列表
let page = 1;
const pageSize = 10;
const noMore = ref(false);
const loading = ref(false);
const currentPartition = ref(0);
const videoList = ref(Array<VideoType>());
const getVideoList = async (partition: number) => {
  let res: any = null;
  loading.value = true;
  if (partition === 0) {
    res = await getHotVideoAPI(page, pageSize);
  } else {
    res = await getVideoByPartitionAPI(pageSize, partition);
  }
  if (res.data.code === statusCode.OK) {
    if (res.data.data.videos) {
      videoList.value = videoList.value.concat(res.data.data.videos);
    } else {
      noMore.value = true;
    }
  }
  loading.value = false;
}

const partitionChange = (partition: number) => {
  page = 1;
  noMore.value = false;
  videoList.value = [];
  currentPartition.value = partition;

  getVideoList(partition);
}

const scrollBox = ref<HTMLElement | null>(null);
const scrollList = (() => {
  //节流
  var timer: number | null = null;
  return () => {
    if (timer) {
      return
    }

    timer = setTimeout(() => {
      const scrollTop = scrollBox.value?.scrollTop || 0;
      const clientHeight = scrollBox.value?.clientHeight || 0;
      const scrollHeight = scrollBox.value?.scrollHeight || 0;
      if (scrollTop + clientHeight >= scrollHeight - 50) {
        if (!noMore.value && !loading.value) {
          page++;
          getVideoList(currentPartition.value);
        }
      }
      timer = null;
    }, 500);
  }
})();

onBeforeMount(() => {
  // getVideoList();

})
</script>

<style lang="scss" scoped>
.header {
  top: 0;
  z-index: 999;
  position: fixed;
  height: 50px;
}

.partition {
  position: fixed;
  top: 50px;
  width: 100%;
  z-index: 999;


}

.content {
  margin-top: 100px;
  overflow-y: scroll;
  height: calc(100vh - 100px);

  .carousel {
    height: 220px;
  }

  .loading {
    padding: 20px 0;
    text-align: center;
  }
}
</style>