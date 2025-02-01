<template>
  <div class="video" v-title :data-title="`视频播放-${globalConfig.title}`">
    <header-bar class="header-bar"></header-bar>
    <div v-if="!loading" class="video-main">
      <div class="video-player">
        <video-player v-if="videoInfo" :video-info="videoInfo" :part="currentPart"></video-player>
      </div>
      <!-- 视频信息 -->
      <div class="video-info-box">
        <div class="video-title-box">
          <h1 class="title-text" :class="spread ? 'title-spread' : ''">{{ videoInfo?.title }}</h1>
          <div class="spread" @click="spread = !spread">
            <n-icon size="22" color='#999'>
              <arrow-down v-if="spread" />
              <arrow-up v-else />
            </n-icon>
          </div>
        </div>
        <div class="data-box">
          <div class="data-itme">
            <span class="online">1 人在看</span>
            <span>{{ videoInfo?.clicks }} 播放</span>
          </div>
          <div class="data-itme">
            <span>{{ videoInfo ? formatTime(videoInfo.createdAt) : '' }}</span>
          </div>
        </div>
        <n-collapse-transition :show="!spread">
                  <!-- 视频分集 -->
        <div v-if="videoInfo && videoInfo.resources.length > 1">
          <part-list :resources="videoInfo.resources" :active="currentPart" @change="changePart"></part-list>
        </div>
          <p class="info-item" v-show="videoInfo?.copyright">
            <n-icon color='#fd6d6f'>
              <forbid-icon></forbid-icon>
            </n-icon>
            <span style="padding-left: 6px;">未经作者授权，禁止转载</span>
          </p>
          
          <p class="info-item">{{ videoInfo?.desc }}</p>
        </n-collapse-transition>
        <div class="author">
          <div class="author-left">
            <common-avatar :size="30" :url="videoInfo?.author.avatar"></common-avatar>
            <span class="name">{{ videoInfo?.author.name }}</span>
          </div>
          <div class="author-right">
            <archive-info v-if="videoInfo" :vid="videoInfo.vid"></archive-info>
          </div>
        </div>
      </div>
      <!--发表评论-->
      <comment-list :vid="videoInfo!.vid"></comment-list>
    </div>
    <div v-else class="player-skeleton">
      <n-skeleton height="200px" />
      <n-skeleton text :repeat="2" />
      <n-skeleton text width=" 60%" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { NIcon, NTime, NSkeleton, NCollapseTransition } from 'naive-ui';
import { onBeforeMount, ref } from 'vue';
import { formatTime } from "@/utils/format";

import { Forbid as ForbidIcon } from "@icon-park/vue-next";
import ArrowUp from "@/components/icons/ArrowUpIcon.vue";
import ArrowDown from "@/components/icons/ArrowDownIcon.vue";
import { statusCode } from "@/utils/status-code";
import { globalConfig } from "@/utils/global-config";
import PartList from './components/PartList.vue';
import ArchiveInfo from './components/ArchiveInfo.vue';
import CommentList from './components/CommentList.vue';
import HeaderBar from '@/components/header-bar/index.vue';
import VideoPlayer from '@/components/video-player/index.vue';
import CommonAvatar from "@/components/common-avatar/index.vue";

// import { CommonAvatar, VideoPlayer, IosVideoPlayer } from '@leaf/components';
// import { Forbid, ArrowDown, ArrowUp } from "@leaf/icons";

import { getVideoInfoAPI } from '@/api/video';

const route = useRoute();
const router = useRouter();
const vid = parseInt(route.params.vid.toString());

const ios = ref(navigator.userAgent.includes("iPhone") ? 1 : 0);
const currentPart = ref(1);//当前分集
const spread = ref(true);// 展开

//获取视频信息
const loading = ref(true);
const resources = ref([]);
const videoInfo = ref<VideoType | null>(null);
const getVideoInfo = async (vid: number) => {
  const res = await getVideoInfoAPI(vid);
  if (res.data.code === statusCode.OK) {
    videoInfo.value = res.data.data.video;
    //设置播放的资源
    if (!resources.value[currentPart.value - 1]) {
      currentPart.value = 1;
    }

    //修改网站标题
    document.title = `${res.data.data.video.title}-${globalConfig.title}`
    loading.value = false;
  }
  // getVideoInfoAPI(vid, ios).then((res) => {
  //   if (res.data.code === statusCode.OK) {
  //     videoInfo.value = res.data.data.video;
  //     resources.value = res.data.data.video.resources;
  //     //设置播放的资源
  //     if (!resources.value[part.value - 1]) {
  //       part.value = 1;
  //     }

  //     //修改网站标题
  //     document.title = `${res.data.data.video.title}-${window.$title || globalConfig.title}`
  //     loading.value = false;
  //   }
  // })
}

const changePart = (target: number) => {
  if (resources.value[target - 1]) {
    currentPart.value = target;
  }
  router.replace({ query: { p: currentPart.value } });
}

onBeforeMount(() => {
  getVideoInfo(vid);
  if (route.query.p) {
    currentPart.value = Number(route.query.p);
  }
})
</script>

<style lang="scss" scoped>
.video {
  height: 100%;
  width: 100%;
}

.header-bar {
  position: fixed;
  top: 0;
  z-index: 1000;
}

.video-main {
  width: 100%;
  background: #fff;
  margin-top: 50px;
  display: flex;
  flex-direction: column;
  // flex-wrap: nowrap;

  .video-player {
    width: 100%;
  }

  //骨架占位
  .player-skeleton {
    margin: 0 auto;
    width: 100%;
  }
}



// .video-info {
//   margin: 50px 0 10px;
//   padding: 0 10px 16px;
//   border-bottom: 1px solid #efeff5;

//   .title-wrapper {
//     display: flex;
//     align-items: center;
//     justify-content: space-between;

//     .title {
//       font-size: 4.26667vmin;
//       font-weight: 400;
//     }

//     .fold {
//       width: 50px;
//       padding-top: 6px;
//       text-align: right;
//     }
//   }

//   .info-item {
//     color: #666;
//     margin: 4px 0 2px 0;
//   }

//   .author {
//     margin-top: 12px;
//     display: flex;
//     align-items: center;
//     justify-content: space-between;

//     .info-fold {
//       display: flex;
//       align-items: center;

//       .name {
//         font-size: 16px;
//         line-height: 30px;
//         padding-left: 6px;
//       }
//     }
//   }
// }

.video-info-box {
  margin-top: 1.86667vmin;
  margin-bottom: 3.73333vmin;
  padding: 0 3.2vmin;
  border-bottom: 1px solid #f0f2f5;

  .video-title-box {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .title-text {
      width: calc(100% - 22px);
      margin: 0;
      padding: 0;
      font-size: 4.26667vmin;
      font-weight: 400;
      line-height: 5.86667vmin;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }

    .spread {
      position: absolute;
      width: 22px;
      right: 0;
      top: 0;
      height: 5.86667vmin;
    }

    .title-spread {
      white-space: normal;
    }
  }

  .data-box {
    color: #999;
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 6px;

    .data-itme {
      .online {
        padding-right: 6px;

      }
    }
  }

  .info-item {
    color: #666;
    margin: 4px 0 0 0;
  }

  .author {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 3.46667vmin;
    padding-bottom: 3.73333vmin;

    .author-left {
      display: flex;
      align-items: center;

      .name {
        color: #18191c;
        margin-left: 2.13333vmin;
      }
    }
  }
}
</style>