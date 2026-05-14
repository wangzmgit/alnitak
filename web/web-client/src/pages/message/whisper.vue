<template>
  <div class="whisper-wrap">
    <!--左边-->
    <div class="msg-left">
      <div class="left-top"></div>
      <div class="left-scroll">
        <el-empty v-if="!msgList.length && !initLoading" description="暂无私信" />
        <div class="msg-user-item" v-for="(item, index) in msgList" :key="index" @click="getMsgContent(item, index)">
          <common-avatar class="msg-avatar" :url="item.user.avatar" :size="45"></common-avatar>
          <span class="msg-name">{{ item.user.name }}</span>
          <span class="msg-date"> {{ formatRelativeTime(item.createdAt) }}</span>
          <div class="dot" v-if="!item.status"></div>
        </div>
      </div>
    </div>
    <!--右侧-->
    <div class="msg-right">
      <div class="left-top right-top">{{ targetUser.name || "消息内容" }}</div>
      <template v-if="!msgList.length && !initLoading">
        <el-empty class="right-empty" description="暂无私信" />
      </template>
      <div v-else ref="msgBox" class="msg-main" @scroll="lazyLoading">
        <div class="content-container" v-for="(item, index) in msgDetails" :key="index">
          <!-- 自己发送的 -->
          <div v-if="item.fromId == userInfo?.uid">
            <common-avatar class="avatar-right" :url="userInfo.avatar" :size="40"></common-avatar>
            <span class="content-right">{{ item.content }}</span>
          </div>
          <!-- 收到的 -->
          <div v-else>
            <common-avatar class="avatar-left" :url="targetUser.avatar" :size="40"></common-avatar>
            <div class="content-left-box">
              <span class="content-left">{{ item.content }}</span>
            </div>
          </div>
          <!-- 解决div无法撑开的问题 -->
          <div style="clear:both;"></div>
        </div>
      </div>
      <div class="msg-input">
        <el-input v-model="msgForm.content" placeholder="发个消息呗~" maxlength="255" show-word-limit type="textarea" :rows="3" id="whisper-msg-input" name="whisperMsg"
          :autosize="{ minRows: 4, maxRows: 4 }" resize="none" />
        <div class="btn-box">
          <div></div>
          <el-button type="primary" :loading="sendLoading" @click="sendMsg">发送</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Base64 } from 'js-base64';
import { formatRelativeTime } from "@/utils/format";
import { ElMessage } from 'element-plus';
import { getUserBaseInfoAPI, getUserInfoAPI } from '@/api/user';
import { onBeforeMount, onBeforeUnmount, reactive, ref, nextTick } from 'vue';
import { getWhisperDetailsAPI, getWhisperListAPI, readWhisperAPI, sendWhisperAPI } from "@/api/msg-whisper";
import { storageData } from "@/utils/storage-data";
import { globalConfig } from "@/utils/global-config";

definePageMeta({
  middleware: ['auth']
})

const userInfo = ref<UserInfoType>();
const getUserInfo = async () => {
  const res = await getUserInfoAPI();
  if (res.data.code === statusCode.OK) {
    userInfo.value = res.data.data.userInfo;
  }
}

const initLoading = ref(true);
const msgList = ref<WhisperListType[]>([]);
const msgDetails = ref<WhisperDetailsType[]>([]);
const msgForm = reactive({
  fid: 0,
  content: ''
})

//获取消息列表
const getMsgList = async () => {
  const res = await getWhisperListAPI();
  if (res.data.code === statusCode.OK) {
    msgList.value = res.data.data.messages;
    if (msgList.value) {
      if (msgForm.fid !== 0) {
        initSendUser(msgForm.fid);
      }
      if (msgList.value.length > 0) {
        getMsgContent(msgList.value[0], 0);
      }
    }
  }
}

//初始化发送的用户
const initSendUser = async (fid: number) => {
  //遍历当前消息列表查找用户
  for (let i = 0; i < msgList.value.length; i++) {
    if (msgList.value[i].user.uid === fid) {
      getMsgContent(msgList.value[i], i);
      return;
    }
  }
  // 获取用户信息并添加到用户列表
  const res = await getUserBaseInfoAPI(fid);
  if (res.data.code === statusCode.OK) {
    const user = res.data.data.userInfo;
    msgList.value.unshift({
      user: {
        uid: user.uid,
        name: user.name,
        avatar: user.avatar
      },
      createdAt: new Date().toString(),
      status: true
    });
    targetUser.name = user.name;
    targetUser.avatar = user.avatar;
  }
}

//获取消息详情
const page = ref(1);
const loading = ref(false);
const noMore = ref(false);//没有更多
const allowToBottom = ref(true);//是否允许自动跳转到底部
const targetUser = reactive({
  name: "",
  avatar: ""
})

const getMsgContent = (item: WhisperListType, index: number) => {
  page.value = 1;
  noMore.value = false;
  msgDetails.value = [];
  msgForm.fid = item.user.uid;
  targetUser.name = item.user.name;
  targetUser.avatar = item.user.avatar;
  getMoreMsgContent(item.user.uid);
  //已读消息
  readWhisperAPI(item.user.uid);
  msgList.value[index].status = true;
}

//获取更多消息
const getMoreMsgContent = async (fid: number) => {
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
    getMoreMsgContent(msgForm.fid);
  }
}

//到达底部
const toBottom = () => {
  if (allowToBottom.value) {
    msgBox.value!.scrollTop = msgBox.value!.scrollHeight;
  } else {
    msgBox.value!.scrollTop = 100;
    allowToBottom.value = true;
  }
}

//发送消息
const sendLoading = ref(false);
const sendMsg = async () => {
  if (!msgForm.content) {
    ElMessage.error('不能发送空消息');
    return;
  }
  sendLoading.value = true;
  const res = await sendWhisperAPI(msgForm);
  if (res.data.code === statusCode.OK) {
    msgDetails.value.push({
      fid: 0,
      fromId: userInfo.value!.uid,
      content: msgForm.content,
      createdAt: ""
    });
    msgForm.content = "";
    sendLoading.value = false;
    nextTick(() => {
      toBottom()
    })
  } else {
    ElMessage.error(res.data.msg || '发送失败');

  }
  sendLoading.value = false;
}

//websocket
let SocketURL = "";
let websocket: WebSocket | null = null;
//初始化weosocket
const initWebSocket = async () => {
  const wsProtocol = globalConfig.https ? 'wss://' : 'ws://';
  SocketURL = `${wsProtocol}${globalConfig.domain}/api/v1/message/ws?token=${storageData.get("token")}`;
  websocket = new WebSocket(SocketURL);
  websocket.onmessage = websocketOnmessage;
  websocket.onerror = () => {};
  websocket.onclose = () => {};
}

//数据接收
const websocketOnmessage = (e: any) => {
  const res = JSON.parse(Base64.decode(e.data));
  // 收到后端 ping，立即回复 pong
  if (res.type === 'ping') {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
      websocket.send(JSON.stringify({ type: 'pong' }));
    }
    return;
  }
  if (msgForm.fid === res.fid) {
    msgDetails.value.push({
      fid: msgForm.fid,
      fromId: res.fid,
      content: res.content,
      createdAt: ""
    });
    nextTick(() => {
      toBottom();
    })
  } else {
    // 不再强制切换会话，改为标记未读并将该会话置顶
    const index = msgList.value.findIndex(item => item.user.uid === res.fid);
    if (index >= 0) {
      const chat = msgList.value[index];
      chat.status = false; // 未读
      chat.createdAt = new Date().toString();
      // 移动到顶部
      msgList.value.splice(index, 1);
      msgList.value.unshift(chat);
    } else {
      // 若列表中没有该会话，则新增到顶部，但不切换当前会话
      initSendUser(res.fid);
    }
  }
}

//加载和卸载页面
const route = useRoute();
onBeforeMount(async () => {
  try {
    if (route.query.target) {
      msgForm.fid = Number(route.query.target);
    }
    await getUserInfo();
    await initWebSocket();
    await getMsgList();
  } catch (e) {
    console.error('[whisper] onBeforeMount error:', e);
  } finally {
    initLoading.value = false;
  }
})

onBeforeUnmount(() => {
  if (websocket) {
    websocket.close();
  }
})
</script>

<style lang="scss" scoped>
.whisper-wrap {
  position: absolute;
  top: 0;
  left: 10px;
  right: 60px;
  bottom: 0;
  display: flex;
  background-color: var(--bg-elev-1);
}

.msg-left {
  display: flex;
  flex-direction: column;
  min-height: 0;
  width: 230px;
  border-right: 1px solid var(--border-color);

  .left-top {
    flex-shrink: 0;
    height: 40px;
    border-bottom: 1px solid var(--border-color);
  }

  .left-scroll {
    flex: 1;
    overflow-y: auto;
    min-height: 0;

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-thumb {
      background-color: var(--border-color);
      outline: none;
      border-radius: 2px;
    }

    &::-webkit-scrollbar-track {
      box-shadow: none;
      border-radius: 2px;
    }
  }

  .msg-user-item {
    width: 100%;
    height: 60px;
    display: flex;
    align-items: center;
    padding: 0 8px;
    position: relative;
    cursor: pointer;
    border-bottom: 1px solid var(--border-color);

    &:hover {
      background-color: var(--hover-bg);
    }

    // .msg-avatar {}

    .dot {
      top: 10px;
      left: 42px;
      width: 10px;
      height: 10px;
      position: absolute;
      border-radius: 50%;
      background-color: var(--font-danger);
    }

    .msg-name {
      position: absolute;
      top: 8px;
      left: 60px;
      color: var(--font-primary-1);
      font-size: 14px;
    }

    .msg-date {
      position: absolute;
      top: 32px;
      left: 60px;
      color: var(--font-primary-3);
      font-size: 12px;
    }
  }
}

.msg-right {
  display: flex;
  flex-direction: column;
  min-height: 0;
  width: calc(100% - 230px);

  .right-top {
    flex-shrink: 0;
    height: 40px;
    color: var(--font-primary-1);
    text-align: center;
    font-size: 16px;
    line-height: 40px;
    border-bottom: 1px solid var(--border-color);
  }

  .msg-main {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    background-color: var(--fill-1);
    border-bottom: 1px solid var(--border-color);

    /**修改滚动条样式 */
    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-thumb {
      background-color: var(--border-color);
      outline: none;
      border-radius: 2px;
    }

    &::-webkit-scrollbar-track {
      box-shadow: none;
      border-radius: 2px;
    }
  }
}

.right-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

/**聊天部分 */
.content-container {
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
    background-color: var(--primary-color);
    color: var(--primary-text-color);
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
    margin-top: 7px;
    max-width: 80%;
    background: var(--bg-elev-1);
    border: 1px solid var(--border-color);
    box-shadow: 0 1px 4px var(--shadow-weak);
    padding: 5px 10px 5px 10px;
    border-radius: 3px;
  }

  .content-left {
    font-size: 1rem;
    color: var(--font-primary-1);
  }
}

.msg-input {
  flex-shrink: 0;
  height: 150px;
  padding: 10px;
  position: relative;

  .btn-box {
    height: 50px;
    display: flex;
    margin-left: 10px;
    align-items: center;
    justify-content: space-between;

    button {
      width: 100px;
      height: 30px;
    }
  }
}
</style>