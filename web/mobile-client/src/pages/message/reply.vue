<template>
  <div class="reply">
    <p class="reply-title">回复我的</p>
    <div class="reply-box">
      <el-scrollbar max-height="100%">
        <ul class="reply-list">
          <li class="reply-msg-item" v-for="(item, index) in replyMessageList" :key="index">
            <div class="item-left">
              <common-avatar class="avatar" :url="item.user.avatar" :size="45"></common-avatar>
            </div>
            <div class="item-center">
              <p class="title">
                <nuxt-link class="user-name" :to="`/user/${item.user.uid}`">{{ item.user.name }}</nuxt-link>
                <span v-if="item.rootContent"> 回复了你的评论</span>
                <span v-else> 对你的{{ item.type === 0 ? "视频" : "文章" }}发表评论</span>
              </p>
              <p class="content">{{ item.content }}</p>
              <p class="target-content" v-if="item.targetReplyContent">{{ item.targetReplyContent }}</p>
              <span class="msg-time"> {{ formatTime(item.createdAt) }}</span>
            </div>
            <div class="item-right">
              <nuxt-link v-if="item.type === 0" :to="`/video/${item.video.vid}`">
                <div class="root-content" v-if="item.rootContent">{{ item.rootContent }}</div>
                <el-image v-else class="img" :src="getResourceUrl(item.video.cover)" lazy alt="封面"></el-image>
              </nuxt-link>
              <nuxt-link v-else :to="`/article/${item.article.aid}`">
                <el-image v-if="item.article.cover" class="img" :src="getResourceUrl(item.article.cover)" lazy
                  alt="封面"></el-image>
                <div class="root-content" v-else>{{ item.article.title }}</div>
              </nuxt-link>
            </div>
          </li>
        </ul>
      </el-scrollbar>
    </div>
    <div v-show="showPagination" class="pagination">
      <el-pagination small background layout="prev, pager, next" :page="page" :page-size="pageSize" :total="total" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount } from "vue";
import { ElPagination } from 'element-plus';
import { formatTime } from "@/utils/format";
import { getReplyMsgAPI } from '@/api/msg-reply';
import CommonAvatar from '@/components/common-avatar/index.vue';

const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const showPagination = ref(true);
const replyMessageList = ref<ReplyMessageType[]>([]);
const getReplyMsgList = async () => {
  const res = await getReplyMsgAPI(page.value, pageSize.value);
  if (res.data.code === statusCode.OK) {
    total.value = res.data.data.total;
    replyMessageList.value = res.data.data.messages;
  }
}

onBeforeMount(() => {
  getReplyMsgList();
})
</script>

<style lang="scss" scoped>
.reply {
  padding: 0 18px 0;
  height: 100%;
  box-sizing: border-box;
  background-color: #fff;

  .reply-title {
    font-size: 18px;
    margin: 0;
    padding: 16px 0 10px;
  }

  .reply-box {
    width: 100%;
    height: calc(100% - 86px);
  }

  .pagination {
    height: 36px;
    display: flex;
    align-items: center;
  }
}

.reply-list {
  padding: 0;
  margin: 0;

  .reply-msg-item {
    display: flex;
    align-items: center;
    padding: 12px 0 16px;
    border-bottom: 1px solid #e1e2e3;

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
          font-weight: 600;

          &:visited {
            color: #222;
          }

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

      .target-content {
        color: #999;
        font-size: 12px;
        margin: 5px 0;
        padding-left: 6px;
        border-left: 2px solid #e7e7e7;
      }

      .msg-time {
        font-size: 12px;
        color: #999;
      }
    }

    .item-right {
      width: 90px;
      height: 55px;

      .root-content {
        width: 90px;
        font-size: 14px;
        color: #999;
        overflow: hidden;
        text-overflow: ellipsis;
        word-break: break-word;
        display: -webkit-box;
        -webkit-line-clamp: 3;
        -webkit-box-orient: vertical;
      }

      .img {
        cursor: pointer;
        width: 100%;
        height: 100%;
        border-radius: 2px;
      }
    }
  }
}
</style>