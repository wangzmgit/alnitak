<template>
  <ul class="video-list">
    <li class="video-item" v-for="item in videoList">
      <nuxt-link class="cover" :to="`/video/${item.vid}`">
        <img class="img" :src="getResourceUrl(item.cover)" />
      </nuxt-link>
      <nuxt-link class="title" :to="`/video/${item.vid}`">{{ item.title }}</nuxt-link>
      <div class="meta">
        <div class="play-count">
          <span class="time">看到 {{ toDuration(item.time) }}</span>
        </div>
        <div class="date">{{ formatDate(item.updatedAt) }}</div>
      </div>
    </li>
  </ul>
</template>

<script setup  lang="ts">
import { toDuration } from "@/utils/format";

const props = defineProps<{
  videoList: HistoryVideoType[];
}>()

</script>

<style lang="scss" scoped>
.video-list {
  margin: 0;
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

        .time {
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
</style>