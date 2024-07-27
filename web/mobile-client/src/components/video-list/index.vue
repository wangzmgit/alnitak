<template>
  <div class="video-list-box">
    <div class="video-list">
      <div class="v-card" v-for="(item, index) in props.videos" @click="goToVideo(item.vid)">
        <div class="card">
          <img class="cover" :src="getResourceUrl(item.cover)" />
        </div>
        <p class="title">{{ item.title }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { getResourceUrl } from '@/utils/resource';

const props = defineProps<{
  videos: Array<VideoType>
}>();

const router = useRouter();
const goToVideo = (vid: number) => {
  router.push({ name: "Video", params: { vid: vid } });
}
</script>

<style lang="scss" scoped>
.video-list-box {
  padding: 0 1.33333vmin;
  margin-bottom: 5.33333vmin;
}

.video-list {
  display: flex;
  flex-wrap: wrap;

  .v-card {
    width: 50%;
    box-sizing: border-box;
    padding: 2.13333vmin 1.33333vmin;
    display: inline-block;

    .card {
      position: relative;
      background-color: #f3f3f3;
      border-radius: 0.53333vmin;
      overflow: hidden;

      &::before {
        box-sizing: border-box;
        display: block;
        width: 100%;
        padding-bottom: 56.25%;
        content: "";
      }

      .cover {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
      }
    }

    .title {
      margin: 0;
      font-size: 3.2vmin;
      color: #212121;
      margin-top: 1.6vmin;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }
  }
}
</style>