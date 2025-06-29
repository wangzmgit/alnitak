<template>
  <div class="video">
    <header-bar class="header"></header-bar>
    <div class="video-main">
      <div class="mian-content">
        <div class="left-column">
          <div class="video-player" ref="playerContainerRef">
            <client-only>
              <video-player v-if="videoInfo && playerReady" ref="playerRef" :video-info="videoInfo" :part="currentPart" :progress="pendingProgress" :key="currentPart"></video-player>
            </client-only>
            <div v-if="!showPlayer" class="skeleton"></div>
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
              <archive-info v-if="videoInfo" :vid="videoInfo.vid"></archive-info>
            </div>
            <div class="toolbar-right">
              <span>{{ onlineCount }} 人在看</span>
              <span>{{ videoInfo?.clicks }} 播放</span>
              <span>{{ videoInfo ? formatTime(videoInfo.createdAt) : '' }}</span>
            </div>
          </div>
          <!-- 简介部分 -->
          <div class="video-desc-container">
            <div ref="descRef" class="basic-desc-info" :style="`height: ${foldDesc ? foldDescHeight : 'auto'};`">
              <span class="desc-info-text">{{ videoInfo?.desc }}</span>
            </div>
            <div class="toggle-btn" v-show="showFoldBtn" @click="foldDesc = !foldDesc">
              <span class="toggle-btn-text">{{ foldDesc ? '展开更多' : '收起' }}</span>
            </div>
          </div>
          <!-- 标签部分 -->
          <div class="tags-box">
            <div class="tag" v-for="item in videoInfo?.tags.split(',')">{{ item }}</div>
          </div>
          <!-- 评论区 -->
          <comment-list v-if="videoInfo" :vid="videoInfo.vid"></comment-list>
        </div>
        <div class="right-column">
          <!-- 作者信息 -->
          <author-card v-if="videoInfo" :info="videoInfo.author"></author-card>
          <!-- 添加弹幕列表 -->
          <danmaku-list ref="danmakuListRef" :height="danmakuListHeight"></danmaku-list>
          <!-- 视频分集 -->
          <div v-if="videoInfo && videoInfo.resources.length > 1">
            <part-list :resources="videoInfo.resources" :active="currentPart" @change="changePart"></part-list>
          </div>
          <!-- 相关推荐 -->
          <recommend-list v-if="videoInfo" :vid="videoInfo.vid"></recommend-list>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, type ComponentPublicInstance } from "vue";
import { ElIcon } from "element-plus";
import { Forbid as ForbidIcon } from "@icon-park/vue-next";
import { formatTime } from "@/utils/format";
import PartList from "./components/PartList.vue";
import AuthorCard from './components/AuthorCard.vue';
import ArchiveInfo from './components/ArchiveInfo.vue';
import CommentList from "./components/CommentList.vue";
import DanmakuList from "./components/DanmakuList.vue";
import HeaderBar from "@/components/header-bar/index.vue";
import VideoPlayer from "@/components/video-player/index.vue";
import RecommendList from "./components/RecommendList.vue";
import { asyncGetVideoInfoAPI } from "@/api/video";
import { createUUID } from "@/utils/uuid";
import { getDanmakuAPI } from "@/api/danmaku";
import { getHistoryProgressAPI, addHistoryAPI } from "@/api/history";

const route = useRoute();
const router = useRouter();

// 获取视频信息
const videoInfo = ref<VideoType>();
const videoId = route.params.id.toString();
const { data } = await asyncGetVideoInfoAPI(videoId);
if ((data.value as any).code === statusCode.OK) {
  videoInfo.value = (data.value as any).data.video as VideoType;
} else {
  // 处理视频信息不存在
  navigateTo('/404');
}


const playerContainerRef = ref<HTMLElement | null>(null)
const danmakuListHeight = ref(300);
const playerRef = ref<ComponentPublicInstance<{ seek: (time: number) => void; uploadHistory: () => void; setDanmaku: (data: any[]) => void; setOnReady: (cb: () => void) => void; }> | null>(null);

const handelResize = () => {
  nextTick(() => {
    danmakuListHeight.value = ((playerContainerRef.value?.clientWidth || 730) * 0.5625) + 40 - 104;
  })
}

// 视频分集
// 校验 p 参数有效性，无效则重定向到 p1
if (route.query.p && Number(route.query.p) > videoInfo.value!.resources.length) {
  router.replace({ path: `/video/${videoId}`, query: { p: 1 } });
}
const currentPart = ref(Number(route.query.p) || 1);
const pendingProgress = ref<number | null>(null);

const onPlayerReady = () => {
  if (pendingProgress.value !== null && playerRef.value && playerRef.value.seek) {
    playerRef.value.seek(pendingProgress.value);
    pendingProgress.value = null;
  }
};

watch(playerRef, (val) => {
  if (val && val.setOnReady) {
    val.setOnReady(onPlayerReady);
  }
});

let needReportAfterSwitch = false;

const changePart = async (target: number) => {
  // 切换分集前先上报历史
  if (playerRef.value && playerRef.value.uploadHistory) {
    await playerRef.value.uploadHistory();
  }
  if (videoInfo.value?.resources[target - 1]) {
    currentPart.value = target;
  }
  router.replace({ query: { p: currentPart.value } });
  needReportAfterSwitch = true;

  // 主动请求新分集进度
  if (videoInfo.value) {
    const res = await getHistoryProgressAPI(videoInfo.value.vid, currentPart.value);
    if (res.data.code === 200 && res.data.data && typeof res.data.data.progress === 'number' && res.data.data.progress > 0) {
      pendingProgress.value = res.data.data.progress;
    } else {
      pendingProgress.value = null;
    }
  }
}

// 监听路由参数 p，自动切换分P和弹幕
watch(() => route.query.p, (newP) => {
  const partNum = Number(newP) || 1;
  // 如果分P有效，则切换；否则重定向到 p1
  if (videoInfo.value?.resources[partNum - 1]) {
    currentPart.value = partNum;
    getDanmakuList(videoInfo.value.vid, partNum);
  } else {
    router.replace({ path: `/video/${videoId}`, query: { p: 1 } });
  }
});

// 获取弹幕列表
const danmakuListRef = ref<InstanceType<typeof DanmakuList> | null>(null);

const getDanmakuList = async (vid: number, part: number) => {
  const res = await getDanmakuAPI(vid, part);
  if (res.data.code === statusCode.OK) {
    const danmakus = res.data.data.danmaku || [];
    nextTick(() => {
      playerRef.value?.setDanmaku(danmakus)
      danmakuListRef.value?.setDanmaku(danmakus)
    })
  }
};

// 简介部分
const foldDesc = ref(true); // 是否折叠简介
const descRef = ref<HTMLElement>();
const showPlayer = ref(false);
const showFoldBtn = ref(false); // 是否显示展开和折叠按钮
const foldDescHeight = ref('auto'); // 折叠状态下简介的最大高度
const playerReady = ref(false);
onMounted(async () => {
  if (descRef.value!.clientHeight >= 80) {
    showFoldBtn.value = true;
    foldDescHeight.value = '80px';
  } else {
    showFoldBtn.value = false;
    foldDescHeight.value = 'auto';
  }

  if (videoInfo.value) {
    try {
      // 带 p 参数时也请求 getProgress?vid=xx&part=xx
      let res;
      if (!route.query.p) {
        res = await getHistoryProgressAPI(videoInfo.value.vid);
      } else {
        res = await getHistoryProgressAPI(videoInfo.value.vid, Number(route.query.p));
      }
      if (res && res.data.code === 200 && res.data.data) {
        const { part, progress } = res.data.data;
        if (part && part !== currentPart.value && videoInfo.value.resources[part - 1]) {
          currentPart.value = part;
          router.replace({ query: { p: part } });
          await nextTick();
        }
        getDanmakuList(videoInfo.value.vid, part || currentPart.value);
        if (typeof progress === 'number' && progress > 0) {
          console.log('[id.vue] pendingProgress赋值:', progress);
          pendingProgress.value = progress;
        } else {
          console.log('[id.vue] pendingProgress赋值: null');
          pendingProgress.value = null;
        }
      } else {
        getDanmakuList(videoInfo.value.vid, currentPart.value);
        console.log('[id.vue] pendingProgress赋值: null');
        pendingProgress.value = null;
      }
    } catch (e) {
      getDanmakuList(videoInfo.value.vid, currentPart.value);
    }
  }

  handelResize();
  window.addEventListener("resize", handelResize);

  nextTick(() => {
    showPlayer.value = true;
    playerReady.value = true;
  })
})

//websocket
const onlineCount = ref(1);//在线人数
let SocketURL = "";
let websocket: WebSocket | null = null;
//初始化weosocket
const initWebSocket = () => {
  let clientId = localStorage.getItem("ws-client-id");
  if (!clientId) {
    clientId = createUUID();
    localStorage.setItem("ws-client-id", clientId);
  }
  const wsProtocol = window.location.protocol === 'http:' ? 'ws://' : 'wss://';
  const domain = globalConfig.domain || window.location.host;
  SocketURL = wsProtocol + domain + `/api/v1/online/video?vid=${videoId}&clientId=${clientId}`;

  websocket = new WebSocket(SocketURL);
  websocket.onmessage = websocketOnmessage;
}

//数据接收
const websocketOnmessage = (e: any) => {
  const res = JSON.parse(e.data);
  onlineCount.value = res.number;
}

onBeforeMount(() => {
  initWebSocket();
})

onBeforeUnmount(() => {
  window.removeEventListener("resize", handelResize);
  if (websocket) {
    websocket.close();
    websocket = null;
  }
})

watch(pendingProgress, (val) => {
  console.log('[id.vue] pendingProgress变化:', val);
  console.log('[id.vue] 实时传递给video-player的progress:', pendingProgress.value);
  if (needReportAfterSwitch && typeof val === 'number' && val > 0) {
    needReportAfterSwitch = false;
    if (videoInfo.value) {
      addHistoryAPI({ vid: videoInfo.value.vid, part: currentPart.value, time: val });
      console.log('[id.vue] 上报切换分集后的历史记录:', { vid: videoInfo.value.vid, part: currentPart.value, time: val });
    }
  }
});
</script>

<style lang="scss" scoped>
.header {
  position: fixed;
}

.video-main {
  padding-top: 80px;
  margin: 0 auto;
  min-width: 1200px;
  /* 保持最小宽度为 1200px */
}

.mian-content {
  display: flex;
  width: 85%;
  max-width: calc(100% - 100px);
  margin: 0 auto;
  position: relative;
}

.left-column {
  flex: 1;

  .video-player {
    position: relative;
    margin: 0 auto;
    width: 100%;
    /*16:9*/
    min-width: 680px;
    min-height: 382px;

    .skeleton {
      width: 100%;
      padding-bottom: 56.25%;
      background-color: #f0f2f5;
    }
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

      .desc-info-text {
        white-space: pre-line;
      }
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

  // 标签部分
  .tags-box {
    padding-bottom: 6px;
    margin: 16px 0 20px 0;
    border-bottom: 1px solid #E3E5E7;

    .tag {
      color: #61666d;
      background: #f1f2f3;
      height: 28px;
      line-height: 28px;
      border-radius: 14px;
      font-size: 13px;
      padding: 0 12px;
      box-sizing: border-box;
      transition: all .3s;
      display: inline-flex;
      align-items: center;
      cursor: pointer;
      margin: 0 12px 8px 0;
    }
  }
}

.right-column {
  width: 340px;
  margin-left: 30px;
  z-index: 1;
}
</style>
