<template>
  <div class="msg" v-title data-title="消息">
    <header-bar></header-bar>
    <div>
      <div class="announce-box" @click="goAnnounce()">
        <img src="@/assets/images/common/announce.png" />
        <span class="announce-name">系统公告</span>
      </div>
    </div>
    <!-- 消息列表 -->
    <div v-for="(item, index) in msgList" :key="index">
      <div class="msg-user-item" @click="goMsgDetails(item)">
        <common-avatar class="msg-avatar" :url="item.user.avatar" :size="45"></common-avatar>
        <span class="msg-name">{{ item.user.name }}</span>
        <span class="msg-date">
          <n-time :time="new Date(item.createdAt)"></n-time>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import { NTime } from "naive-ui";
import { useRouter } from "vue-router";
import { getWhisperListAPI } from "@/api/msg-whisper";
import { statusCode } from "@/utils/status-code";
import HeaderBar from "@/components/header-bar/index.vue";
import CommonAvatar from "@/components/common-avatar/index.vue";

const router = useRouter();

const msgList = ref<WhisperListType[]>([]);
const getMsgList = async () => {
  const res = await getWhisperListAPI();
  if (res.data.code === statusCode.OK) {
    msgList.value = res.data.data.messages;
  }
}

const goAnnounce = () => {
  router.push({ name: "Announce" });
}

//前往私信
const goMsgDetails = (msg: WhisperListType) => {
  router.push({
    name: "Whisper", query: {
      fid: msg.user.uid,
      name: msg.user.name,
      avatar: msg.user.avatar
    }
  });
}

onBeforeMount(() => {
  getMsgList();
})
</script>

<style lang="scss" scoped>
.msg {
  height: 100%;
  width: 100%;
  top: 0;
  bottom: 0;
  position: fixed;
  background: #f4f4f4;
  overflow-y: auto;
}

.announce-box {
  width: 100%;
  height: 60px;
  display: flex;
  align-items: center;
  background-color: #fff;
  border-bottom: 1px solid #e0e0e0;
}

.announce-box>img {
  width: 45px;
  height: 45px;
  margin: 8px 0 0 6px;
}

.announce-name {
  margin-left: 10px;
  font-size: 16px;
}

.msg-user-item {
  width: 100%;
  height: 60px;
  background-color: #fff;
  position: relative;
  border-bottom: 1px solid #e0e0e0;
}

.msg-avatar {
  margin: 8px 0 0 6px;
}

.msg-name {
  position: absolute;
  top: 8px;
  left: 60px;
  font-size: 16px;
}

.msg-date {
  position: absolute;
  top: 28px;
  left: 60px;
  color: #808080;
}
</style>