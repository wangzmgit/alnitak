<template>
  <div class="search-page">
    <header-bar></header-bar>
    <div class="search">
      <div class="search-form">
        <input class="input" v-model="searchParam.keywords" @keydown.enter="searchVideo()">
        <button class="btn" @click="searchVideo()">
          <search-icon class="icon" size="16" />
        </button>
      </div>
    </div>
    <div class="card-list">
      <video-item v-for="item in videoList" :info="item" :keywords="searchParam.keywords"></video-item>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, onBeforeUnmount, ref } from "vue";
import { searchVideoAPI } from "@/api/video";
import HeaderBar from "@/components/header-bar/index.vue";
import VideoItem from "./components/VideoItem.vue";
import { Search as SearchIcon } from '@icon-park/vue-next';

const route = useRoute();

const searchParam = reactive<SearchVideoType>({
  page: 1,
  pageSize: 15,
  keywords: "",
})

let noMore = false;
let loading = false;
const videoList = ref<VideoType[]>([]);
const searchVideo = async (init = false) => {
  loading = true;
  if (init) {
    searchParam.page = 1;
    videoList.value = [];
    noMore = false;
  }
  const res = await searchVideoAPI(searchParam);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.videos) {
      videoList.value.push(...res.data.data.videos);
      if (res.data.data.videos.length < 15) {
        noMore = true;
      }
    } else {
      noMore = true;
      ElMessage.error("获取失败");
    }
  }
  loading = false;
}

//加载更多
const lazyLoading = () => {
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
  const clientHeight = document.documentElement.clientHeight;
  const scrollHeight = document.documentElement.scrollHeight;
  if (scrollTop + clientHeight >= (scrollHeight - 10)) {
    if (!noMore && !loading) {
      searchParam.page++;
      searchVideo();
    }
  }
}

onBeforeMount(() => {
  searchParam.keywords = route.params.keywords.toString();
  searchVideo();
  window.addEventListener('scroll', lazyLoading, true);
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', lazyLoading, true);
})
</script>

<style lang="scss" scoped>
.search {
  width: 100%;
  margin: 30px auto;
  width: 500px;

  .search-form {
    position: relative;

    .input {
      border: 1px solid #e5e5e5;
      outline: none;
      padding: 8px 30px 8px 11px;
      height: 36px;
      font-size: 12px;
      line-height: 14px;
      border-radius: 18px;
      width: 500px;
      vertical-align: top;
      color: var(--font-primary-1);
      box-sizing: border-box;
    }

    .btn {
      position: absolute;
      top: 0;
      right: 10px;
      border: none;
      width: 20px;
      height: 36px;
      line-height: 36px;
      font-size: 14px;
      vertical-align: top;
      background: transparent;
      padding: 0;
      cursor: pointer;

      .icon {
        display: block;
        margin-top: 3px;
        width: 20px;
        height: 36px;
        color: var(--font-primary-5);
      }
    }
  }
}

.card-list {
  display: grid;
  // row-gap: 16px;
  column-gap: 16px;
  grid-template-columns: repeat(5, 1fr);
  width: 90%;
  margin: 20px auto;
}
</style>