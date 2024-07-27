<template>
  <div class="msg" v-title data-title="私信">
    <div class="msg-header">
      <n-icon :size="22" @click="backMessage()">
        <arrow-left></arrow-left>
      </n-icon>
      <p calss="name">{{ targetUser.name }}</p>
      <span></span>
    </div>
    <div ref="msgBox" class="msg-main" @scroll="lazyLoading">
      <div class="content-box" v-for="(item, index) in msgDetails" :key="index">
        <!--自己发送的-->
        <div v-if="item.fromId == userInfo.uid">
          <common-avatar class="avatar-right" :url="userInfo.avatar" :size="45"></common-avatar>
          <span class="content-right">{{ item.content }}</span>
        </div>
        <!--收到的-->
        <div v-else>
          <common-avatar class="avatar-left" :url="targetUser.avatar" :size="45"></common-avatar>
          <div class="content-left-box">
            <span class="content-left">{{ item.content }}</span>
          </div>
        </div>
        <!-- 解决div无法撑开的问题 -->
        <div style="clear:both;">
      </div>
    </div>

    <div class="msg-input">
      <n-input v-model:value="msgForm.content" placeholder="发个消息呗~" maxlength="255" show-count
        @keydown.enter="sendMsg"></n-input>
      <n-button type="primary" :loading="sendLoading" @click="sendMsg">发送</n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Base64 } from 'js-base64';
import { useRoute, useRouter } from 'vue-router';
import { NInput, NIcon, NButton, useMessage } from 'naive-ui';
import { onBeforeMount, onBeforeUnmount, reactive, ref, computed, nextTick } from 'vue';
import ArrowLeft from '@/components/icons/ArrowLeftIcon.vue';
import CommonAvatar from "@/components/common-avatar/index.vue";
import { statusCode } from '@/utils/status-code';
import { getWhisperDetailsAPI, sendWhisperAPI } from '@/api/msg-whisper';
import { globalConfig } from '@/utils/global-config';
import { storageData } from '@/utils/storage-data';

// import { globalConfig, statusCode, storageData } from '@leaf/utils';
// import { getWhisperDetailsAPI, sendWhisperAPI } from '@leaf/apis';

const msgDetails = ref<WhisperDetailsType[]>([]);
const msgForm = reactive({
  fid: 0,
  content: ''
})

const message = useMessage();//通知
//用户信息
const userInfo = ref<UserInfoType>({
  uid: 0,
  name: "",
  avatar: ""
});

//取消息详情
const page = ref(1);
const loading = ref(false);
const noMore = ref(false);//没有更多
const allowToBottom = ref(true);//是否允许自动跳转到底部
const targetUser = reactive({
  name: "",
  avatar: ""
})

//获取更多消息
const getMsgContent = async (fid: number) => {
  loading.value = true;
  const res = await getWhisperDetailsAPI(fid, page.value, 10);
  if (res.data.code === statusCode.OK) {
    const resMsg = res.data.data.messages;
    if (resMsg.length < 10) {
      noMore.value = true;
    }
    msgDetails.value = resMsg.concat(msgDetails.value);
    nextTick(() => {
      toBottom()
    })
  }
  loading.value = false;
}

//加载更多
const msgBox = ref<HTMLElement | null>(null);
const lazyLoading = () => {
  const scrollTop = msgBox.value?.scrollTop || Infinity;
  if (scrollTop < 30 && !loading.value && !noMore.value) {
    page.value++;
    allowToBottom.value = false;
    getMsgContent(msgForm.fid);
  }
}

//到达底部
const toBottom = () => {
  if (allowToBottom.value) {
    msgBox.value!.scrollTop = msgBox.value!.scrollHeight;
  } else {
    msgBox.value!.scrollTop = 160;
    allowToBottom.value = true;
  }
}

const router = useRouter();
const backMessage = () => {
  router.push({ name: "Message" });
}

//发送消息
const sendLoading = ref(false);
const sendMsg = async () => {
  if (!msgForm.content) {
    message.error('不能发送空消息');
    return;
  }
  sendLoading.value = true;
  const res = await sendWhisperAPI(msgForm);
  if (res.data.code === statusCode.OK) {
    msgDetails.value.push({
      fid: 0,
      fromId: userInfo.value.uid,
      content: msgForm.content,
      createdAt: ""
    });
    msgForm.content = "";
    sendLoading.value = false;
    nextTick(() => {
      toBottom()
    })
  }
}

//websocket
let SocketURL = "";
let websocket: WebSocket | null = null;
//初始化weosocket
const initWebSocket = async () => {
  const wsProtocol = window.location.protocol === 'http:' ? 'ws://' : 'wss://';
  const domain = globalConfig.domain || window.location.host;
  SocketURL = `${wsProtocol}${domain}/api/v1/message/ws?token=${storageData.get("token")}`;
  websocket = new WebSocket(SocketURL);
  websocket.onmessage = websocketOnmessage;
}

//数据接收
const websocketOnmessage = (e: any) => {
  const res = JSON.parse(Base64.decode(e.data));
  if (msgForm.fid === res.fid) {
    msgDetails.value.push({
      fid: msgForm.fid,
      fromId: res.fid,
      content: res.content,
      createdAt: ""
    });
    nextTick(() => {
      toBottom()
    })
  }
}

//加载和卸载页面
const route = useRoute();
onBeforeMount(() => {
  userInfo.value = storageData.get("user_info");
  if (route.query.fid) {
    msgForm.fid = Number(route.query.fid);
    targetUser.name = (route.query.name || "").toString();
    targetUser.avatar = (route.query.avatar || "").toString();
  }
  getMsgContent(msgForm.fid);
})

onBeforeUnmount(() => {
  if (websocket) {
    websocket.close();
  }
})
</script>

<style lang="less" scoped>
.msg {

  .msg-header {
    position: fixed;
    height: 50px;
    width: 100%;
    z-index: 999;
    background: #fff;
    display: flex;
    align-items: center;
    justify-content: space-between;

    span {
      width: 50px;
      text-align: center;
    }
  }

  .msg-main {
    padding-top: 50px;
    height: calc(100vh - 100px);
    background-color: #f7f7f7;
    border-bottom: 1px solid #e1e2e3;
    overflow-y: auto;

    /**修改滚动条样式 */
    &::-webkit-scrollbar {
      width: 8px;
    }

    &::-webkit-scrollbar-thumb {
      /*滚动条里面小方块*/
      border-radius: 2px;
      background: #999999;
    }

    &::-webkit-scrollbar-track {
      /*滚动条里面轨道*/
      border-radius: 5px;
    }
  }

  /**聊天部分 */
  .content-box {
    min-height: 70px;
    margin: 6px 20px;

    &:nth-child(1) {
      margin-top: 10px;
    }

    .avatar-right {
      float: right;
    }

    .content-right {
      float: right;
      max-width: 80%;
      margin-right: 10px;
      margin-top: 6px;
      background-color: #40a9ff;
      color: #fff;
      font-size: 16px;
      border-radius: 3px;
      padding: 5px 10px 5px 10px;
    }

    .avatar-left {
      float: left;
    }

    .content-left-box {
      float: left;
      margin-left: 10px;
      margin-top: 10px;
      max-width: 80%;
      background: #fff;
      padding: 5px 10px 5px 10px;
      border-radius: 3px;
    }

    .content-left {
      font-size: 1rem;
    }
  }

  .msg-input {
    height: 50px;
    display: flex;
    padding: 0 6px;
    align-items: center;
    justify-content: space-between;

    button {
      width: 80px;
    }
  }
}
</style>