<template>
  <div class="video">
    <header-bar class="header"></header-bar>
    <div class="video-main">
      <div class="mian-content">
        <div class="left-column" :style="{ width: `${videoMainWidth}px` }">
          <div class="video-player">
            <video-player></video-player>
          </div>
          <!-- 标题和版权信息 -->
          <div class="video-title-box">
            <p class="video-title">{{ videoInfo?.title }}</p>
            <p v-show="videoInfo?.copyright" class="copyright">
              <el-icon class="icon" color='#fd6d6f'>
                <forbid-icon></forbid-icon>
              </el-icon>
              <span>未经作者授权，禁止转载</span>
            </p>
          </div>
          <!-- 点赞收藏等数据 -->
          <div class="video-toolbar">
            <div class="toolbar-left">
              <archive-info></archive-info>
            </div>
            <div class="toolbar-right">
              <span>1 人在看</span>
              <span>{{ videoInfo?.clicks }} 播放</span>
              <span>2024-03-03 20:57:17</span>
            </div>
          </div>
          <!-- 简介部分 -->
          <div class="video-desc-container">
            <div class="basic-desc-info" :style="`height: ${foldDesc ? '80px' : 'auto'};`">
              <span class="desc-info-text">{{ videoInfo?.desc }}</span>
            </div>
            <div class="toggle-btn" @click="foldDesc = !foldDesc">
              <span class="toggle-btn-text">{{ foldDesc ? '展开更多' : '收起' }}</span>
            </div>
          </div>
        </div>
        <div class="right-column">
          <!-- 作者信息 -->
          <author-card v-if="videoInfo" :info="videoInfo.author"></author-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";
import { ElIcon } from "element-plus";
import { Forbid as ForbidIcon } from "@icon-park/vue-next";
import AuthorCard from './components/AuthorCard.vue';
import ArchiveInfo from './components/ArchiveInfo.vue';
import HeaderBar from "@/components/header-bar/index.vue";
import VideoPlayer from "@/components/video-player/index.vue";
import { asyncGetVideoInfoAPI } from "@/api/video";

const route = useRoute();

// 获取视频信息
const videoInfo = ref<VideoType>();
const videoId = route.params.id.toString();
const { data } = await asyncGetVideoInfoAPI(videoId);
if ((data.value as any).code === statusCode.OK) {
  videoInfo.value =  (data.value as any).data.video as VideoType;
}


const videoMainWidth = ref(0);
const handelResize = () => {
  // w = (16 / 9) *  (屏幕高度 - marginTop - 96px - headerBar高度)
  videoMainWidth.value = (16 / 9) * (window.innerHeight - 176);
}

// 简介部分
const foldDesc = ref(true); // 是否折叠简介

onMounted(() => {
  handelResize();
  window.addEventListener("resize", handelResize);
})

onBeforeUnmount(() => {
  window.removeEventListener("resize", handelResize);
})

</script>

<style lang="scss" scoped>
.header {
  position: fixed;
}

.video-main {
  padding-top: 80px;
  margin: 0 auto 0;
  min-width: 1200px;

  .mian-content {
    display: flex;
    width: calc(100% - 100px);
    margin: auto 50px;
    position: relative;
    // background-color: antiquewhite;
  }
}

.left-column {
  flex: 1;

  .video-player {
    background-color: aquamarine;
    margin: 0 auto;
    width: 100%;
    /*16:9*/
    min-width: 680px;
    min-height: 382px;
  }

  .video-title-box {
    width: 100%;
    height: 54px;
    display: flex;

    .video-title {
      width: calc(100% - 160px);
      font-weight: 500;
      line-height: 28px;
      margin: 13px 0;
      font-size: 20px;
      color: #18191C;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      font-family: PingFang SC, Helvetica Neue, Microsoft YaHei, sans-serif;
    }

    .copyright {
      width: 180px;
      display: flex;
      align-items: center;
      justify-content: flex-end;
      font-size: 13px;
      color: #9499A0;

      .icon {
        padding: 0 6px;
      }
    }
  }

  .video-toolbar {
    color: #9499A0;
    font-size: 13px;
    padding-bottom: 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid #E3E5E7;

    .toolbar-right {
      display: inline-block;

      span {
        margin-left: 20px;
      }
    }
  }


  // 简介部分
  .video-desc-container {
    margin: 16px 0;

    .basic-desc-info {
      white-space: pre-line;
      letter-spacing: 0;
      color: #18191C;
      font-size: 15px;
      line-height: 24px;
      overflow: hidden;
    }

    .toggle-btn {
      margin-top: 10px;
      font-size: 13px;
      line-height: 18px;

      .toggle-btn-text {
        cursor: pointer;
        color: #61666D;

        &:hover {
          color: var(--primary-hover-color);
        }
      }
    }
  }
}

.right-column {
  width: 340px;
  margin-left: 30px;
  z-index: 1;


}
</style>