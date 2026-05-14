<template>
  <div class="recommend-list">
    <div class="video-card" v-for="item in videoList" :key="item.pgc_id">
      <div class="card-box">
        <nuxt-link class="cover-box" :to="toWatchLink(item)">
          <img :src="getResourceUrl(item.cover)" alt="封面" />
          <span class="duration" v-if="item.new_ep?.index_show">{{ item.new_ep.index_show }}</span>
        </nuxt-link>
        <div class="info">
          <nuxt-link class="title" :to="toWatchLink(item)">{{ item.title }}</nuxt-link>
          <div class="up-name">
            <span class="name">{{ item.latest_ep_title || item.desc || '暂无简介' }}</span>
          </div>
          <div class="play-info">
            <span class="val" v-if="item.rating">评分 {{ item.rating }}</span>
            <span class="val" v-if="item.current_episodes"> · {{ item.current_episodes }} 集</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { asyncGetPGCRecommendByVideoAPI } from "@/api/pgc";

const props = defineProps<{
  vid: number | string;
}>();

const videoList = ref<PGCRecommendItem[]>([])
const { data } = await asyncGetPGCRecommendByVideoAPI(props.vid, 1, 12);
if ((data.value as any).code === statusCode.OK) {
  videoList.value = (data.value as any).data.list || [];
}

const cleanVid = (raw: unknown): number => {
  const n = Number(String(raw ?? '').trim());
  return Number.isFinite(n) && n > 0 ? Math.floor(n) : 0;
};

const toWatchLink = (item: PGCRecommendItem) => {
  const epId = Number(item.latest_ep_id || 0);
  if (Number.isFinite(epId) && epId > 0) {
    return `/watch?ep=${Math.floor(epId)}&mode=pgc`;
  }
  const vid = cleanVid(item.latest_vid);
  return vid > 0 ? `/watch?v=${vid}&mode=pgc` : '/404';
};

const currentPlayIndex = ref(-1);

const getNextVideo = () => {
  if (currentPlayIndex.value < videoList.value.length - 1) {
    currentPlayIndex.value++;
    const item = videoList.value[currentPlayIndex.value];
    return { vid: cleanVid(item.latest_vid), epId: Number(item.latest_ep_id || 0) || undefined };
  }
  return null;
};

const resetPlayIndex = (vid: number) => {
  const index = videoList.value.findIndex(video => cleanVid(video.latest_vid) === vid);
  currentPlayIndex.value = index >= 0 ? index : -1;
};

defineExpose({
  videoList: readonly(videoList),
  getNextVideo,
  resetPlayIndex
});
</script>

<style lang="scss" scoped>
.recommend-list {
  margin-top: 18px;

  .video-card {
    margin-bottom: 12px;

    .card-box {
      display: flex;

      .cover-box {
        position: relative;
        width: 140px;
        height: 80px;
        border-radius: 6px;
        cursor: pointer;
        background-color: var(--fill-1, #c9ccd0);

        img {
          width: 100%;
          height: 100%;
          border-radius: 6px;
          object-fit: cover;
        }

        .duration {
          position: absolute;
          bottom: 6px;
          right: 6px;
          color: #fff;
          height: 20px;
          line-height: 20px;
          z-index: 5;
          font-size: 12px;
          background-color: rgba(0, 0, 0, 0.45);
          border-radius: 2px;
          padding: 0 4px;
        }
      }

      .info {
        margin-left: 10px;
        flex: 1;

        .title {
          margin: 0;
          min-height: 38px;
          color: var(--font-primary-1);
          font-size: 15px;
          line-height: 19px;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          line-clamp: 2;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
        }

        .up-name {
          width: 100%;
          display: inline-flex;
          align-items: center;
          color: var(--font-primary-3);
          font-size: 13px;
          margin: 2px 0;
          height: 20px;

          .name {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
        }

        .play-info {
          color: var(--font-primary-3);
          display: inline-flex;
          align-items: center;

          .val {
            font-size: 13px;
          }
        }
      }
    }
  }
}
</style>

