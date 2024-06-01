<template>
  <div class="recommend-list">
    <div class="video-card" v-for="item in videoList">
      <div class="card-box">
        <nuxt-link class="cover-box" :to="`/video/${item.vid}`">
          <img :src="getResourceUrl(item.cover)" alt="封面" />
          <span class="duration">{{ toDuration(item.duration) }}</span>
        </nuxt-link>
        <div class="info">
          <nuxt-link class="title" :to="`/video/${item.vid}`">{{ item.title }}</nuxt-link>
          <div class="up-name">
            <el-icon class="icon" :size="16">
              <up-icon></up-icon>
            </el-icon>
            <nuxt-link class="name" :to="`/user/${item.author.uid}`">{{ item.author.name }}</nuxt-link>
          </div>
          <div class="play-info">
            <el-icon class="icon" :size="16">
              <play-count-icon></play-count-icon>
            </el-icon>
            <span class="val">{{ item.clicks }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import UpIcon from '@/components/icons/UpIcon.vue';
import { asyncGetRelatedVideoList } from "@/api/video";
import PlayCountIcon from '@/components/icons/PlayCountIcon.vue';

const props = defineProps<{
  vid: number
}>();

const videoList = ref<VideoType[]>([])
const { data } = await asyncGetRelatedVideoList(props.vid);
if ((data.value as any).code === statusCode.OK) {
  videoList.value = (data.value as any).data.videos;
}
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
        background-color: #c9ccd0;

        img {
          width: 100%;
          height: 100%;
          border-radius: 6px;
        }

        .duration {
          position: absolute;
          bottom: 6px;
          right: 6px;
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

      .info {
        margin-left: 10px;
        flex: 1;

        .title {
          margin: 0;
          min-height: 38px;
          color: #18191c;
          font-size: 15px;
          line-height: 19px;
          transition: color 0.3s;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
        }

        .up-name {
          width: 100%;
          height: 100%;
          display: inline-flex;
          align-items: center;
          color: #9499a0;
          font-size: 13px;
          cursor: pointer;
          margin: 2px 0;
          height: 20px;
          transition: color 0.3s;

          .icon {
            margin-right: 4px;
          }

          .name {
            color: #9499a0;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
        }

        .play-info {
          color: #9499a0;
          display: inline-flex;
          align-items: center;

          .icon {
            width: 18px;
            height: 18px;
            margin-right: 4px;
          }

          .val {
            font-size: 13px;
          }
        }
      }
    }
  }
}
</style>