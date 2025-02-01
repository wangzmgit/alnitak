<template>
  <div class="video-item">
    <nuxt-link class="img" :to="`/video/${info.vid}`" target="_blank">
      <img :src="getResourceUrl(info.cover)" alt="封面" />
      <span class="duration">{{ toDuration(info.duration) }}</span>
    </nuxt-link>
    <div class="video-info">
      <nuxt-link v-if="!props.keywords" class="title" :to="`/video/${info.vid}`" target="_blank"
        v-html="keyHighlight(info.title)"></nuxt-link>
      <nuxt-link v-else class="title" :to="`/video/${info.vid}`">{{ info.title }}</nuxt-link>
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
import DOMPurify from 'isomorphic-dompurify';

const props = defineProps<{
  info: VideoType,
  keywords?: string,
}>()

//关键词高亮
const keyHighlight = (title: string) => {
  let res = '';
  let indexArr: Array<number> = []; // 需要标红的字的下标数组
  const keywordsArray = props.keywords!.split(" ");
  const getReplaceStr = (str: string) => `<font color="var(--primary-color)">${str}</font>`;
  keywordsArray.forEach((keyword) => {
    let filterStr = title;
    let stopFlag = false;
    while (!stopFlag && filterStr && keyword) {
      const index = filterStr.indexOf(keyword); // 返回匹配的第一个字符的下标
      if (index === -1) stopFlag = true;
      else {
        keyword.split("").forEach((s, i) => {
          indexArr.push(index + Number(i));
        });
        filterStr = `${filterStr.substring(0, index)} ${filterStr.substring(index + 1)}`;
      }
    }
  });
  indexArr = Array.from(new Set(indexArr)); // 去重
  title.split("").forEach((char, charIndex) => {
    res += indexArr.includes(charIndex) ? getReplaceStr(char) : char;
  });

  return DOMPurify.sanitize(res);
}
</script>

<style lang="scss" scoped>
.video-item {
  width: 100%;
  height: 260px;

  .img {
    position: relative;
    display: block;
    width: 100%;
    height: 160px;
    border-radius: 10px;
    overflow: hidden;
    cursor: pointer;
    background-color: rgba(0, 0, 0, .2);

    img {
      width: 100%;
      height: 100%;
    }

    .duration {
      position: absolute;
      bottom: 10px;
      right: 12px;
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
      -webkit-box-orient: vertical;
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