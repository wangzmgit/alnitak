<template>
  <div class="like">
    <p class="like-title">收到的赞</p>
    <div class="like-box">
      <el-scrollbar max-height="100%">
        <ul class="like-list">
          <li class="like-msg-item" v-for="(item, index) in likeMessageList" :key="index">
            <div class="item-left">
              <common-avatar class="avatar" :url="item.user.avatar" :size="45"></common-avatar>
            </div>
            <div class="item-center">
              <p class="title">
                <nuxt-link class="user-name" :to="`/user/${item.user.uid}`">{{ item.user.name }}</nuxt-link>
                <span> 赞了我的{{ item.type === 0 ? "视频" : "文章" }}</span>
              </p>      
              <span class="msg-time"> {{ formatTime(item.createdAt) }}</span>
            </div>
            <div class="item-right">
              <nuxt-link class="user-name" :to="getContentUrl(item)">
                <el-image v-if="getCoverUrl(item)" class="img" :src="getResourceUrl(getCoverUrl(item))" lazy
                  alt="封面"></el-image>
                <div class="content-title" v-else>{{ getContentTitle(item) }}</div>
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
import { getLikeMsgAPI } from '@/api/msg-like';
import CommonAvatar from '@/components/common-avatar/index.vue';

const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const showPagination = ref(true);
const likeMessageList = ref<LikeMessageType[]>([]);
const getlikeMsgList = async () => {
  const res = await getLikeMsgAPI(page.value, pageSize.value);
  if (res.data.code === statusCode.OK) {
    total.value = res.data.data.total;
    likeMessageList.value = res.data.data.messages;
  }
}

const getContentUrl = (msg: LikeMessageType) => {
  if (msg.type === 0) {
    return `/video/${msg.video.vid}`;
  }

  return `/article/${msg.article.aid}`;
}

const getCoverUrl = (msg: LikeMessageType) => {
  if (msg.type === 0) {
    return msg.video.cover;
  }

  return msg.article.cover;
}

const getContentTitle = (msg: LikeMessageType) => {
  if (msg.type === 0) {
    return msg.video.title;
  }

  return msg.article.title;
}

onBeforeMount(() => {
  getlikeMsgList();
})
</script>

<style lang="scss" scoped>
.like {
  padding: 0 18px 0;
  height: 100%;
  box-sizing: border-box;
  background-color: #fff;

  .like-title {
    font-size: 18px;
    margin: 0;
    padding: 16px 0 10px;
  }

  .like-box {
    width: 100%;
    height: calc(100% - 86px);
  }

  .pagination {
    height: 36px;
    display: flex;
    align-items: center;
  }
}

.like-list {
  padding: 0;
  margin: 0;

  .like-msg-item {
    height: 60px;
    display: flex;
    align-items: center;
    padding: 8px 0 12px;
    border-bottom: 1px solid #e1e2e3;

    .item-left {
      width: 55px;
      height: 55px;

      .avatar {
        margin: 5px;

      }
    }

    .item-center {
      flex: 1;
      height: 50px;
      padding-left: 8px;

      .title {
        font-size: 14px;
        color: #222;
        line-height: 18px;
        margin: 2px 0 8px;

        .user-name {
          cursor: pointer;
          font-weight: 600;
          color: var(--primary-color);

          &:visited {
            color: #222;
          }

          &:hover {
            color: var(--primary-hover-color);
          }
        }
      }

      .msg-time {
        font-size: 12px;
        color: #999;
      }
    }

    .item-right {
      width: 90px;
      height: 55px;

      .img {
        cursor: pointer;
        width: 100%;
        height: 100%;
        border-radius: 2px;
      }

      .content-title {
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
    }
  }
}
</style>