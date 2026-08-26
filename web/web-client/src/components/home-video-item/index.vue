<template>
  <div class="video-item">
    <nuxt-link class="img" :to="`/watch?v=${info.shortId || String(info.vid)}`" target="_blank">
      <oss-image :src="info.cover" alt="封面" />
      <span class="duration">{{ toDuration(info.duration) }}</span>
    </nuxt-link>
    <div class="video-info">
      <nuxt-link class="title" :to="`/watch?v=${info.shortId || String(info.vid)}`" target="_blank">{{ info.title }}</nuxt-link>
      <div class="author">
        <div class="avatar">
          <common-avatar :url="info.author.avatar" :size="26" :iconsize="16"></common-avatar>
        </div>
        <div class="name-date">
          <nuxt-link class="name" :to="`/user/${info.author.uid}`" target="_blank">{{ info.author.name }}</nuxt-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { toDuration } from "@/utils/format";

const props = defineProps<{
  info: VideoType;
}>()
</script>

<style lang="scss" scoped>
.video-item {
  width: 100%;

  .img {
    position: relative;
    display: block;
    aspect-ratio: 16 / 9;
    max-height: 200px;
    border-radius: 9px;
    overflow: hidden;
    cursor: pointer;
    background-color: rgba(0, 0, 0, .2);

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .duration {
      position: absolute;
      right: 12px;
      bottom: 10px;
      color: var(--font-primary-6);
      height: 20px;
      line-height: 20px;
      transition: opacity 0.3s;
      z-index: 5;
      font-size: 13px;
      background-color: rgba(0, 0, 0, 0.4);
      border-radius: 2px;
      padding: 0 4px;
    }
  }

  .video-info {
    margin-top: 8px;

    .title {
      height: 44px;
      color: var(--font-primary-1);
      padding-right: 12px;
      font-size: 15px;
      line-height: 22px;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
      word-break: break-all;
      font-weight: 500;
      cursor: pointer;

      &:hover {
        color: var(--primary-color);
      }
    }

    .author {
      display: flex;
      align-items: center;
      font-size: 13px;
      color: var(--font-primary-3);
      //UP信息位置调整
      margin-top: 2px;

      .avatar {
        width: 26px;
        height: 26px;
        border-radius: 50%;
        background-color: rgba(0, 0, 0, .2);
      }

      .name-date {
        margin-left: 10px;

        .name {
          color: var(--font-primary-3);

          &:hover {
            color: var(--primary-color);
          }
        }
      }

    }
  }
}
</style>
