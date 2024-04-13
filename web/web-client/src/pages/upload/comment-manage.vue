<template>
  <div class="comment-manage">
    <div class="header-box">
      <p class="title">评论管理</p>
      <client-only>
        <el-select v-model="videoId" style="width: 200px" @change="videoIdChange">
          <el-option label="全部视频" :value="0"></el-option>
          <el-option v-for="item in videoList" :key="item.vid" :label="item.title" :value="item.vid" />
        </el-select>
      </client-only>
    </div>
    <div class="comment-box">
      <el-scrollbar>
        <ul class="comment-list" v-infinite-scroll="scrollLoad">
          <li class="comment-msg-item" v-for="(item, index) in replyMessageList" :key="index">
            <div class="item-left">
              <common-avatar class="avatar" :url="item.author.avatar" :size="45"></common-avatar>
            </div>
            <div class="item-center">
              <p class="title">
                <nuxt-link class="user-name" :to="`/user/${item.author.uid}`">{{ item.author.name }}</nuxt-link>
                <span v-if="item.rootContent">
                  <span>回复</span>
                  <nuxt-link class="user-name" :to="`/user/${item.target.uid}`">@{{ item.target.name }}</nuxt-link>
                  <span>的{{ item.targetReplyContent ? '回复' : '评论' }}</span>
                  <el-tooltip :content="item.targetReplyContent || item.rootContent" effect="light" placement="bottom">
                    <span class="view-btn">查看{{ item.targetReplyContent ? '回复' : '评论' }}</span>
                  </el-tooltip>
                </span>
              </p>
              <p class="content">{{ item.content }}</p>
              <span class="msg-time"> {{ formatTime(item.createdAt) }}</span>
            </div>
            <div class="item-right">
              <nuxt-link :to="`/video/${item.video.vid}`">
                <el-image class="img" :src="getResourceUrl(item.video.cover)" lazy alt="封面"></el-image>
              </nuxt-link>
            </div>
          </li>
        </ul>
      </el-scrollbar>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount } from "vue";
import { formatTime } from "@/utils/format";
import { getCommentListAPI } from '@/api/comment';
import CommonAvatar from '@/components/common-avatar/index.vue';
import { getAllVideoAPI } from "@/api/video";

const videoId = ref<number>(0);
const videoIdChange = () => {
  page.value = 1;
  total.value = 0;
  noMore.value = false;
  loading.value = false;
  getReplyMsgList();
}

const videoList = ref<AllVideoType[]>([]);
const getAllVideo = async () => {
  const res = await getAllVideoAPI();
  if (res.data.code === statusCode.OK) {
    videoList.value = res.data.data.videos;
  }
}

const page = ref(1);
const total = ref(0);
const pageSize = ref(8);
const noMore = ref(false);
const loading = ref(false);
const replyMessageList = ref<CommentManageType[]>([]);
const getReplyMsgList = async () => {
  if (loading.value || noMore.value) return;
  loading.value = true;
  const res = await getCommentListAPI(videoId.value, page.value, pageSize.value);
  if (res.data.code === statusCode.OK) {
    total.value = res.data.data.total;
    if (res.data.data.comments) {
      replyMessageList.value = replyMessageList.value.concat(res.data.data.comments);
    } else {
      if (page.value === 1) { // 第一页且没有数据
        replyMessageList.value = [];
      }
      noMore.value = true;
    }
  }
  loading.value = false;
}

const scrollLoad = () => {
  if (!loading.value) {
    page.value++;
    getReplyMsgList();
  }
}

// TODO: 删除评论

onBeforeMount(() => {
  getAllVideo();
  getReplyMsgList();
})
</script>

<style lang="scss" scoped>
.comment-manage {
  padding: 0 18px 0;
  height: 100%;
  box-sizing: border-box;
  background-color: #fff;

  .header-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 0 10px;

    .title {
      font-size: 18px;
      margin: 0;
    }
  }

  .comment-box {
    width: 100%;
    height: calc(100% - 80px);
  }

  .comment-list {
    list-style: none;
    box-sizing: border-box;
    width: 100%;
    padding: 12px 0 0;
    margin: 0;

    .comment-msg-item {
      display: flex;
      width: 100%;
      height: 80px;
      margin-bottom: 12px;
      border-bottom: 1px solid #e6e6e6;
      padding-bottom: 12px;

      .item-left {
        width: 55px;
        height: 55px;
        align-self: flex-start;

        .avatar {
          margin: 5px;
        }
      }

      .item-center {
        flex: 1;
        min-height: 50px;
        padding-left: 8px;

        .title {
          font-size: 14px;
          color: #222;
          line-height: 18px;
          margin: 2px 0 8px;

          .user-name {
            cursor: pointer;
            color: #767676;
            margin-right: 3px;

            &:hover {
              color: var(--primary-hover-color);
            }
          }


          .view-btn {
            cursor: pointer;
            color: var(--primary-color);
            margin: 0 6px;

            &:hover {
              color: var(--primary-hover-color);
            }
          }
        }

        .content {
          color: #222;
          font-size: 14px;
          margin: 10px 0 6px;
        }

        .msg-time {
          font-size: 12px;
          color: #999;
        }
      }

      .item-right {
        width: 90px;
        height: 55px;
        padding-right: 10px;

        .img {
          cursor: pointer;
          width: 100%;
          height: 100%;
          border-radius: 2px;
        }
      }
    }
  }
}
</style>