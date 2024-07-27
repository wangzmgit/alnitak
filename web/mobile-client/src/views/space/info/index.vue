<template>
  <div v-title data-title="个人中心">
    <header-bar></header-bar>
    <div class="user-card">
      <div class="card-avatar">
        <common-avatar :url="userInfo?.avatar" :size="90"></common-avatar>
      </div>
      <!--关注粉丝信息-->
      <div class="data-info-box">
        <div class="data-info">
          <span class="data-item">
            <span class="data-num">{{ userData.videoCount }}</span>
            <br>
            <span class="data-title">投稿</span>
          </span>
          <span class="data-item">
            <span class="data-num">{{ userData.followingCount }}</span>
            <br>
            <span class="data-title">关注</span>
          </span>
          <span class="data-item">
            <span class="data-num">{{ userData.followerCount }}</span>
            <br>
            <span class="data-title">粉丝</span>
          </span>
        </div>
        <n-button type="primary" ghost @click="navigateTo('EditInfo')">编辑信息</n-button>
      </div>
    </div>
    <div class="name-box">
      <div class="name-wrap">
        <span class="name">{{ userInfo?.name }}</span>
      </div>
      <div class="sign-wrap">
        <span class="sign" :class="spread ? 'spread' : ''">{{ userInfo?.sign }}</span>
        <div class="spread-btn" @click="handleSpread">{{ spread ? '收起' : '展开' }}</div>
      </div>
    </div>
    <video-list :videos="videoList"></video-list>
    <div v-show="loading" class="loading">
      <n-spin></n-spin>
    </div>
    <footer-bar></footer-bar>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { ref, reactive, onMounted, onBeforeUnmount } from "vue";
import { NButton, NSpin, useMessage } from 'naive-ui';
import { statusCode } from "@/utils/status-code";
import { getUserInfoAPI } from "@/api/user";
import { getUploadVideoAPI } from "@/api/video";
import { getFollowDataAPI } from "@/api/relation";
import HeaderBar from "@/components/header-bar/index.vue";
import FooterBar from '@/components/footer-bar/index.vue';
import VideoList from '@/components/video-list/index.vue'
import CommonAvatar from "@/components/common-avatar/index.vue";

const router = useRouter();
const message = useMessage();

const userData = reactive({
  videoCount: 0,
  followingCount: 0,
  followerCount: 0
})

const spread = ref(false);
const handleSpread = () => {
  spread.value = !spread.value;
}

//前往关注和粉丝页面
const navigateTo = (name: string) => {
  router.push({ name: name });
}

//获取关注数和粉丝数
const getFollowData = async (id: number) => {
  const res = await getFollowDataAPI(id);
  if (res.data.code === statusCode.OK) {
    userData.followerCount = res.data.data.follower;
    userData.followingCount = res.data.data.following;
  }
}

//视频相关
let page = 1;
const pageSize = 10;
let noMore = false;
const loading = ref(false);
const videoList = ref<Array<VideoType>>([]);
const lazyLoading = (e: Event) => {
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
  if (scrollTop === 0) return;

  const clientHeight = document.documentElement.clientHeight;
  const scrollHeight = document.documentElement.scrollHeight;
  if (scrollTop + clientHeight >= scrollHeight) {
    if (!loading.value && !noMore) {
      page++;
      getUploadVideo();
    }
  }
}


//获取我的视频
const getUploadVideo = async () => {
  loading.value = true;
  const res = await getUploadVideoAPI(page, pageSize);
  if (res.data.code === statusCode.OK) {
    // total.value = res.data.data.total;
    if (res.data.data.videos) {
      videoList.value.push(...res.data.data.videos);
      if (res.data.data.videos.length < pageSize) {
        noMore = true;
        message.info("没有更多了");
      }
    } else {
      noMore = true;
      message.info("没有更多了");
    }
  }
  loading.value = false;
}

const userInfo = ref<UserInfoType>();
const getUserBaseInfo = async () => {
  const res = await getUserInfoAPI();
  if (res.data.code === statusCode.OK) {
    userInfo.value = res.data.data.userInfo;
  }
}

onMounted(async () => {
  await getUserBaseInfo();
  getUploadVideo();
  getFollowData(userInfo.value!.uid);
  window.addEventListener('scroll', lazyLoading, true);

})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', lazyLoading, true);
})
</script>

<style lang="scss" scoped>
.user-card {
  display: flex;
  padding: 3.2vmin 3.2vmin 0;
  align-items: flex-start;

  .card-avatar {
    width: 160px;

    span {
      margin: 16px 30px;
    }
  }

  .data-info-box {
    width: 58.66667vmin;

    .data-info {
      display: flex;
      align-items: center;
      justify-content: space-around;

      .data-item {
        width: 19.2vmin;
        display: inline-block;
        text-align: center;
        vertical-align: middle;

        .data-title {
          font-size: 2.93333vmin;
          color: #999;
        }

        .data-num {
          font-size: 3.73333vmin;
          color: #212121;
          line-height: 4.26667vmin;
        }
      }
    }

    button {
      width: 100%;
      height: 8vmin;
      margin-top: 3.2vmin;
    }
  }
}

.name-box {
  clear: both;
  padding: 3.2vmin;
  border-bottom: 1px solid #e7e7e7;

  .name-wrap {
    .name {
      font-size: 4.8vmin;
      color: #212121;
      max-width: 62.66667vmin;
      vertical-align: middle;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      display: inline-block;
    }
  }

  .sign-wrap {
    position: relative;
    clear: both;

    .sign {
      display: block;
      margin-top: 2.13333vmin;
      font-size: 3.46667vmin;
      color: #999;
      line-height: 4.53333vmin;
      width: 83.46667vmin;
      height: 4.53333vmin;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .spread-btn {
      position: absolute;
      top: 0;
      right: 0;
      line-height: 4.53333vmin;
      font-size: 3.46667vmin;
      color: var(--primary-color);
    }
  }
}

// .video-list {
//   height: calc(100vh - 200px);
//   overflow: scroll;

//   .loading {
//     padding: 20px 0;
//     text-align: center;
//   }
// }

.loading {
  padding: 20px 0;
  text-align: center;
}

.spread {
  height: auto !important;
  white-space: normal !important;
  word-break: break-all;
}
</style>