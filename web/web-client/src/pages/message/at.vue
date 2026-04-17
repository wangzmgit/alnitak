<template>
  <div class="at">
    <p class="at-title">@我的</p>
    <div class="at-box">
      <el-scrollbar max-height="100%">
        <ul class="at-list">
          <li class="at-msg-item" v-for="(item, index) in atMessageList" :key="index">
            <div class="item-left">
              <common-avatar class="avatar" :url="item.user.avatar" :size="45"></common-avatar>
            </div>
            <div class="item-center">
              <p class="title">
                <nuxt-link class="user-name" :to="`/user/${item.user.uid}`">{{ item.user.name }}</nuxt-link>
                <span> 在评论中提到了你</span>
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
import { getAtMsgAPI } from '@/api/msg-at';

definePageMeta({
  middleware: ['auth']
})
import CommonAvatar from '@/components/common-avatar/index.vue';

const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const showPagination = ref(true);
const atMessageList = ref<AtMessageType[]>([]);
const getAtMsgList = async () => {
  const res = await getAtMsgAPI(page.value, pageSize.value);
  if (res.data.code === statusCode.OK) {
    total.value = res.data.data.total;
    atMessageList.value = res.data.data.messages;
  }
}

const getContentUrl = (msg: AtMessageType) => {
  if (msg.type === 0) {
    return `/watch?v=${msg.video.shortId || String(msg.video.vid)}`;
  }

  return `/article/${msg.article.aid}`;
}

const getCoverUrl = (msg: AtMessageType) => {
  if (msg.type === 0) {
    return msg.video.cover;
  }

  return msg.article.cover;
}

const getContentTitle = (msg: AtMessageType) => {
  if (msg.type === 0) {
    return msg.video.title;
  }

  return msg.article.title;
}

onBeforeMount(() => {
  getAtMsgList();
})
</script>

<style lang="scss" scoped>
.at {
  padding: 0 18px 0;
  height: 100%;
  box-sizing: border-box;
  background-color: var(--bg-elev-1);

  .at-title {
    font-size: 18px;
    margin: 0;
    padding: 16px 0 10px;
  }

  .at-box {
    width: 100%;
    height: calc(100% - 86px);
  }

  .pagination {
    height: 36px;
    display: flex;
    align-items: center;
  }
}

.at-list {
  padding: 0;
  margin: 0;

  .at-msg-item {
    height: 60px;
    display: flex;
    align-items: center;
    padding: 8px 0 12px;
      border-bottom: 1px solid var(--border-color);

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
        color: var(--font-primary-1);
        line-height: 18px;
        margin: 2px 0 8px;

        .user-name {
          cursor: pointer;
          font-weight: 600;

          &:visited {
            color: var(--font-primary-1);
          }

          &:hover {
            color: var(--primary-hover-color);
          }
        }
      }

      .msg-time {
        font-size: 12px;
        color: var(--font-primary-3);
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
        color: var(--font-primary-3);
        overflow: hidden;
        text-overflow: ellipsis;
        word-break: break-word;
        display: -webkit-box;
        -webkit-line-clamp: 3;
        line-clamp: 3;
        -webkit-box-orient: vertical;
      }
    }
  }
}
</style>