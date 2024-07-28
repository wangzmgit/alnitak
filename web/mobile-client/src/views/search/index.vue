<template>
  <div v-title :data-title="`${globalConfig.title} - 搜索`">
    <header-bar class="header"></header-bar>
    <div ref="scrollBox" class="content" @scroll="scrollList">
      <video-list :videos="videoDataList"></video-list>
      <div v-show="loading" class="loading">
        <n-spin></n-spin>
      </div>
      <!-- <footer-bar></footer-bar> -->
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount, reactive, watch } from 'vue';
import { useRoute } from "vue-router";
import { statusCode } from "@/utils/status-code";
import { globalConfig } from "@/utils/global-config";
import HeaderBar from "@/components/header-bar/index.vue";
import FooterBar from '@/components/footer-bar/index.vue';
import VideoList from '@/components/video-list/index.vue'
import { searchVideoAPI } from "@/api/video";
import { NSpin, useMessage } from 'naive-ui';

const message = useMessage();
const route = useRoute();

const searchParam = reactive<SearchVideoType>({
  page: 1,
  pageSize: 15,
  keywords: "",
})

//视频列表
let page = 1;
const pageSize = 10;
let noMore = false;
const loading = ref(false);
const videoDataList = ref(Array<VideoType>());
const searchVideo = async (init: boolean = false) => {
  loading.value = true;
  if (init) {
    searchParam.page = 1;
    videoDataList.value = [];
    noMore = false;
  }
  const res = await searchVideoAPI(searchParam);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.videos) {
      videoDataList.value.push(...res.data.data.videos);
      if (res.data.data.videos.length < 15) {
        noMore = true;
      }
    } else {
      noMore = true;
      message.error("获取失败");
    }
  }
  loading.value = false;
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
        if (!noMore && !loading.value) {
          page++;
          searchVideo();
        }
      }
      timer = null;
    }, 500);
  }
})();

watch(() => route.query.keywords, (newValue) => {
  searchParam.keywords = (newValue || "").toString();
  searchVideo(true);
})

onBeforeMount(() => {
  searchParam.keywords = (route.query.keywords || "").toString();
  searchVideo();
})
</script>

<style lang="scss" scoped>
.header {
  top: 0;
  z-index: 999;
  position: fixed;
  height: 50px;
}
.content {
  overflow-y: scroll;
  margin-top: 50px;
  height: calc(100vh - 50px);

  .loading {
    padding: 20px 0;
    text-align: center;
  }
}
</style>