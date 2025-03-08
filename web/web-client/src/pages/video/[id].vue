<template>
  <div class="video">
    <header-bar class="header"></header-bar>
    <div class="video-main">
      <div class="mian-content">
        <div class="left-column" :style="{ width: `${videoMainWidth}px` }">
          <div class="video-player">
            <client-only>
              <video-player v-if="videoInfo" :video-info="videoInfo" :part="currentPart"></video-player>
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
          <div class="danmaku-list-container">
            <div class="danmaku-header" @click="toggleDanmakuList">
              <div class="header-left">
                <span class="title">弹幕列表</span>
                <span class="count">({{ danmakuList.length }})</span>
              </div>
              <div class="header-right">
                <el-icon :class="{ 'is-fold': !showDanmakuList }">
                  <arrow-down />
                </el-icon>
              </div>
            </div>
            
            <div class="danmaku-content" v-show="showDanmakuList">
              <!-- 表头 -->
              <div class="danmaku-header-row">
                <div class="time">时间</div>
                <div class="text">弹幕内容</div>
                <div class="send-time">发送时间</div>
              </div>
              
              <!-- 弹幕列表 -->
              <div class="danmaku-item" v-for="item in danmakuList" :key="`${item.time}-${item.text}`">
                <div class="time">{{ formatDanmakuTime(item.time) }}</div>
                <div class="text">{{ item.text }}</div>
                <div class="send-time">{{ dayjs(item.createdAt).format('MM-DD HH:mm') }}</div>
              </div>
            </div>
          </div>
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
import { ref, onMounted, onBeforeUnmount, watch } from "vue";
import { ElIcon } from "element-plus";
import { Forbid as ForbidIcon, ArrowDown } from "@icon-park/vue-next";
import { formatTime } from "@/utils/format";
import PartList from "./components/PartList.vue";
import AuthorCard from './components/AuthorCard.vue';
import ArchiveInfo from './components/ArchiveInfo.vue';
import CommentList from "./components/CommentList.vue";
import HeaderBar from "@/components/header-bar/index.vue";
import VideoPlayer from "@/components/video-player/index.vue";
import RecommendList from "./components/RecommendList.vue";
import { asyncGetVideoInfoAPI } from "@/api/video";
import { createUUID } from "@/utils/uuid";
import { getDanmakuAPI } from "@/api/danmaku";
import dayjs from "dayjs";

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

const videoMainWidth = ref(0);
const handelResize = () => {
  // w = (16 / 9) *  (屏幕高度 - marginTop - 96px - headerBar高度)
  videoMainWidth.value = (16 / 9) * (window.innerHeight - 170);
}

// 视频分集
const currentPart = ref(1);
const changePart = (target: number) => {
  if (videoInfo.value?.resources[target - 1]) {
    currentPart.value = target;
  }
  router.replace({ query: { p: currentPart.value } });
}

// 简介部分
const foldDesc = ref(true); // 是否折叠简介
const descRef = ref<HTMLElement>();
const showPlayer = ref(false);
const showFoldBtn = ref(false); // 是否显示展开和折叠按钮
const foldDescHeight = ref('auto'); // 折叠状态下简介的最大高度
onMounted(() => {
  if (descRef.value!.clientHeight >= 80) {
    showFoldBtn.value = true;
    foldDescHeight.value = '80px';
  } else {
    showFoldBtn.value = false;
    foldDescHeight.value = 'auto';
  }

  handelResize();
  window.addEventListener("resize", handelResize);

  nextTick(() => {
    showPlayer.value = true;
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

onBeforeMount(()=>{
  initWebSocket();
})

onBeforeUnmount(() => {
  window.removeEventListener("resize", handelResize);
})

// 添加弹幕列表相关的代码
const showDanmakuList = ref(false);
const danmakuList = ref<DanmakuType[]>([]);

// 格式化弹幕时间
const formatDanmakuTime = (seconds: number) => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = Math.floor(seconds % 60);
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
};

// 切换弹幕列表显示状态
const toggleDanmakuList = () => {
  showDanmakuList.value = !showDanmakuList.value;
};

// 获取弹幕列表
const getDanmakuList = async (vid: number, part: number) => {
  const res = await getDanmakuAPI(vid, part);
  if (res.data.code === statusCode.OK) {
    danmakuList.value = res.data.data.danmaku || [];
  }
};

// 监听分P变化
watch(() => currentPart.value, (newPart) => {
  if (videoInfo.value) {
    getDanmakuList(videoInfo.value.vid, newPart);
  }
});

// 组件挂载时获取弹幕
onMounted(() => {
  if (videoInfo.value) {
    getDanmakuList(videoInfo.value.vid, currentPart.value);
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
  width: 90%; /* 宽度占满 90% */
  min-width: 1200px; /* 保持最小宽度为 1200px */
  display: flex;
  justify-content: center; /* 确保子元素水平居中 */
  align-items: center; /* 确保子元素垂直居中 */
}

.mian-content {
  display: flex;
  width: 100%; /* 让子元素宽度占满父容器 */
  max-width: calc(100% - 100px); /* 根据父容器计算宽度 */
  margin: auto 50px; /* 确保子元素水平居中，并有50px的左右边距 */
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

.danmaku-list-container {
  margin: 20px 0;
  background: #fff;
  border: 1px solid #e3e5e7;
  border-radius: 8px;
  
  .danmaku-header {
    padding: 12px 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    user-select: none;
    
    &:hover {
      background-color: #f6f7f8;
    }
    
    .header-left {
      display: flex;
      align-items: center;
      
      .title {
        font-size: 14px;
        font-weight: 500;
        color: #18191c;
      }
      
      .count {
        margin-left: 8px;
        color: #9499a0;
        font-size: 12px;
      }
    }

    .header-right {
      .el-icon {
        transition: transform 0.3s;
        
        &.is-fold {
          transform: rotate(-180deg);
        }
      }
    }
  }
  
  .danmaku-content {
    max-height: 300px;
    overflow-y: auto;
    
    .danmaku-header-row {
      padding: 8px 16px;
      display: flex;
      align-items: center;
      background-color: #f6f7f8;
      border-top: 1px solid #f1f2f3;
      font-size: 12px;
      color: #9499a0;
      font-weight: 500;

      .time {
        width: 45px;
      }
      
      .text {
        flex: 1;
        margin: 0 12px;
      }
      
      .send-time {
        width: 85px;
        text-align: right;
      }
    }
    
    .danmaku-item {
      padding: 8px 16px;
      display: flex;
      align-items: center;
      border-top: 1px solid #f1f2f3;
      
      .time {
        width: 45px;
        color: #9499a0;
        font-size: 12px;
        flex-shrink: 0;
      }
      
      .text {
        flex: 1;
        margin: 0 12px;
        color: #18191c;
        font-size: 13px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
      
      .send-time {
        width: 85px;
        flex-shrink: 0;
        text-align: right;
        color: #9499a0;
        font-size: 12px;
      }

      &:hover {
        background-color: #f6f7f8;
      }
    }
    
    &::-webkit-scrollbar {
      width: 4px;
    }
    
    &::-webkit-scrollbar-thumb {
      background-color: #c9ccd0;
      border-radius: 2px;
    }
    
    &::-webkit-scrollbar-track {
      background-color: #f1f2f3;
    }
  }
}
</style>