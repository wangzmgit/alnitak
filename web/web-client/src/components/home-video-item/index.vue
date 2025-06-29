<template>
  <div class="video-item">
    <nuxt-link class="img" :to="`/video/${info.vid}`" target="_blank">
      <img :src="getResourceUrl(info.cover)" alt="封面" />
      <span class="duration">{{ toDuration(info.duration) }}</span>
    </nuxt-link>
    <div class="video-info">
      <nuxt-link class="title" :to="`/video/${info.vid}`" target="_blank">{{ info.title }}</nuxt-link>
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
  height: 260px;

  .img {
    position: relative;
    display: block;
    height: auto;
    /* 让高度自动调整以保持比例 */
    max-height: 160px;
    /* 最大高度限制 */
    border-radius: 10px;
    overflow: hidden;
    cursor: pointer;
    background-color: rgba(0, 0, 0, .2);

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      /* 保持图片比例，裁剪超出的部分 */
    }

    .duration {
      position: absolute;
      right: 12px;
      bottom: 10px;
      color: #fff;
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
      color: #18191c;
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
      color: #9499a0;
      margin-top: 5px;

      .avatar {
        width: 26px;
        height: 26px;
        border-radius: 50%;
        background-color: rgba(0, 0, 0, .2);
      }

      .name-date {
        margin-left: 10px;

        .name {
          color: #9499a0;

          &:hover {
            color: var(--primary-color);
          }
        }
      }

    }
  }
}
</style>
