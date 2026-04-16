<template>
  <ul class="video-list">
    <li class="video-item" v-for="item in videoList">
      <nuxt-link class="cover" :to="watchLink(item)">
        <img class="img" :src="getResourceUrl(item.cover)" />
      </nuxt-link>
      <nuxt-link class="title" :to="watchLink(item)">
        <template v-if="item.pgcAttached && item.pgcTitle">
          <span class="pgc-series">{{ item.pgcTitle }}</span>
          <span v-if="item.episodeNumber > 0 || item.episodeTitle" class="pgc-ep">
            <template v-if="item.episodeNumber > 0">第{{ item.episodeNumber }}话</template>
            <template v-if="item.episodeTitle">
              <template v-if="item.episodeNumber > 0"> </template>{{ item.episodeTitle }}
            </template>
          </span>
        </template>
        <template v-else>{{ item.title }}</template>
      </nuxt-link>
      <div class="meta">
        <div class="play-count">
          <span class="time">
            <template v-if="item.time === -1">已看完</template>
            <template v-else>看到 {{ toDuration(item.time) }}</template>
          </span>
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

const watchLink = (item: HistoryVideoType) => {
  const p = item.part ?? 1;
  if (item.pgcAttached && item.epId && item.epId > 0) {
    if (p > 1) {
      return `/watch?ep=${item.epId}&mode=pgc&p=${p}`;
    }
    return `/watch?ep=${item.epId}&mode=pgc`;
  }
  // 仅当分P > 1 时使用 rid 精准跳转，分P=1 时不带 rid
  if (item.rid && p > 1) {
    return `/watch?v=${item.shortId || String(item.vid)}&rid=${item.rid}`;
  }
  if (p > 1) {
    return `/watch?v=${item.shortId || String(item.vid)}&p=${p}`;
  }
  return `/watch?v=${item.shortId || String(item.vid)}`;
};
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
      background-color: var(--fill-1, #f9f9f9);

      .img {
        width: 100%;
        height: 100%;
        object-fit: contain;
      }
    }

    .title {
      font-size: 12px;
      color: var(--font-primary-1);
      display: block;
      line-height: 20px;
      height: 38px;
      margin-top: 6px;
      overflow: hidden;
      cursor: pointer;

      .pgc-series {
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        font-weight: 500;
      }

      .pgc-ep {
        display: block;
        margin-top: 2px;
        font-size: 11px;
        color: var(--font-primary-3);
        line-height: 16px;
        white-space: nowrap;
        text-overflow: ellipsis;
        overflow: hidden;
      }

      &:hover {
        color: var(--primary-hover-color);
      }
    }

    .meta {
      display: flex;
      align-items: center;
      color: var(--font-primary-3);
      white-space: nowrap;
      margin-top: 5px;
      height: 16px;
      line-height: 16px;

      .play-count {
        width: 92px;
        display: flex;
        align-items: center;

        .time {
          color: var(--font-primary-3);
          font-size: 12px;
        }
      }

      .date {
        color: var(--font-primary-3);
        font-size: 12px;
      }
    }
  }
}
</style>
