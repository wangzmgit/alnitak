<template>
  <div class="recommend-list">
    <!-- 只在单集视频时显示自动连播控制 -->
    <div v-if="showAutoplayControl" class="auto-play-control">
      <p class="next-info">
        接下来播放
        <span class="next-button" @click="autonext = !autonext">
          <span class="txt">自动连播</span>
          <span class="switch-button" :class="{ 'on': autonext }"></span>
        </span>
      </p>
    </div>
    
    <div class="video-card" v-for="item in videoList" :key="item.vid">
      <div class="card-box">
        <nuxt-link class="cover-box" :to="`/watch?v=${item.shortId || String(item.vid)}`">
          <img :src="getResourceUrl(item.cover)" alt="封面" />
          <span class="duration">{{ toDuration(item.duration) }}</span>
        </nuxt-link>
        <div class="info">
          <nuxt-link class="title" :to="`/watch?v=${item.shortId || String(item.vid)}`">{{ item.title }}</nuxt-link>
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
import { useVideoAutonextFollow } from '@/composables/use-video-autonext-follow';

const props = defineProps<{
  vid: number | string;
  showAutoplayControl?: boolean; // 新增：是否显示自动连播控制
}>();

// 与合集列表共用同一开关（localStorage: video-autonext-follow）
const autonext = useVideoAutonextFollow();

const videoList = ref<VideoType[]>([])
const { data } = await asyncGetRelatedVideoList(props.vid);
if ((data.value as any).code === statusCode.OK) {
  videoList.value = (data.value as any).data.videos;
}

// 添加当前播放索引，初始值为 -1
const currentPlayIndex = ref(-1);

// 获取下一个视频
const getNextVideo = () => {
  if (currentPlayIndex.value < videoList.value.length - 1) {
    currentPlayIndex.value++;
    return videoList.value[currentPlayIndex.value];
  }
  return null; // 没有更多视频了
};

// 重置播放索引（当手动切换视频时调用）
const resetPlayIndex = (vid: number) => {
  const index = videoList.value.findIndex(video => video.vid === vid);
  currentPlayIndex.value = index >= 0 ? index : -1;
};

// 暴露给父组件
defineExpose({
  autonext: autonext,
  videoList: readonly(videoList),
  getNextVideo,
  resetPlayIndex
});
</script>

<style lang="scss" scoped>
.recommend-list {
  margin-top: 18px;

  // 添加自动连播控制样式
  .auto-play-control {
    margin-bottom: 16px;
    
    .next-info {
      font-size: 14px;
      color: var(--font-primary-2);
      margin: 0;
      display: flex;
      align-items: center;
      justify-content: space-between;
      
      .next-button {
        display: flex;
        align-items: center;
        cursor: pointer;
        user-select: none;
        
        .txt {
          margin-right: 8px;
          font-size: 13px;
        }
        
        .switch-button {
          width: 32px;
          height: 18px;
          background: var(--border-color);
          border-radius: 9px;
          position: relative;
          transition: background-color 0.3s;
          
          &::after {
            content: '';
            position: absolute;
            top: 2px;
            left: 2px;
            width: 14px;
            height: 14px;
            background: var(--bg-elev-1);
            border-radius: 50%;
            transition: transform 0.3s;
          }
          
          &.on {
            background: var(--primary-hover-color);
            
            &::after {
              transform: translateX(14px);
            }
          }
        }
      }
    }
  }

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

        &.disabled {
          pointer-events: none;
          opacity: .6;
        }

        img {
          width: 100%;
          height: 100%;
          border-radius: 6px;
          object-fit: contain;
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
          color: var(--font-primary-1);
          font-size: 15px;
          line-height: 19px;
          transition: color 0.3s;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          line-clamp: 2;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;

          &.disabled {
            pointer-events: none;
            opacity: .6;
          }
        }

        .up-name {
          width: 100%;
          height: 100%;
          display: inline-flex;
          align-items: center;
          color: var(--font-primary-3);
          font-size: 13px;
          cursor: pointer;
          margin: 2px 0;
          height: 20px;
          transition: color 0.3s;

          .icon {
            margin-right: 4px;
          }

          .name {
            color: var(--font-primary-3);
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
        }

        .play-info {
          color: var(--font-primary-3);
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
