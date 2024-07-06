<template>
  <div class="comment-manage">
    <div class="header-box">
      <base-tabs class="tabs" :tabs="tabsOptions" @tab-change="tabChange" />
      <client-only>
        <el-select v-if="commentType === '0'" v-model="videoId" style="width: 200px" @change="videoIdChange">
          <el-option label="全部视频" :value="0"></el-option>
          <el-option v-for="item in videoList" :key="item.vid" :label="item.title" :value="item.vid" />
        </el-select>
        <el-select v-else v-model="articleId" style="width: 200px" @change="articleIdChange">
          <el-option label="全部专栏" :value="0"></el-option>
          <el-option v-for="item in articleList" :key="item.aid" :label="item.title" :value="item.aid" />
        </el-select>
      </client-only>
    </div>
    <div class="comment-box">
      <el-scrollbar>
        <ul class="comment-list" v-infinite-scroll="scrollLoad">
          <li class="comment-msg-item" v-for="(item, index) in commentList" :key="index">
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
import { getVideoCommentListAPI, getArticleCommentListAPI } from '@/api/comment';

import CommonAvatar from '@/components/common-avatar/index.vue';
import { getAllVideoAPI } from "@/api/video";
import { getAllArticleAPI } from "~/api/article";

const tabsOptions = [
  { key: '0', label: '视频评论' },
  { key: '1', label: '专栏评论', hidden: !globalConfig.article }
];

const commentType = ref('0');
const tabChange = (val: string) => {
  if (commentType.value !== val) { // 类型改变初始化数据
    page.value = 1;
    total.value = 0;
    noMore.value = false;
    loading.value = false;
    commentList.value = [];
    commentType.value = val;
    getReplyMsgList();
  }
}

const videoId = ref<number>(0);
const videoIdChange = () => {
  page.value = 1;
  total.value = 0;
  noMore.value = false;
  loading.value = false;
  commentList.value = [];
  getReplyMsgList();
}

const articleId = ref<number>(0);
const articleIdChange = () => {
  page.value = 1;
  total.value = 0;
  noMore.value = false;
  loading.value = false;
  commentList.value = [];
  getReplyMsgList();
}

const videoList = ref<AllVideoType[]>([]);
const getAllVideo = async () => {
  const res = await getAllVideoAPI();
  if (res.data.code === statusCode.OK) {
    videoList.value = res.data.data.videos;
  }
}

const articleList = ref<AllArticleType[]>([]);
const getAllArticle = async () => {
  const res = await getAllArticleAPI();
  if (res.data.code === statusCode.OK) {
    articleList.value = res.data.data.articles;
  }
}

const page = ref(1);
const total = ref(0);
const pageSize = ref(8);
const noMore = ref(false);
const loading = ref(false);
const commentList = ref<CommentManageType[]>([]);
const getReplyMsgList = async () => {
  if (loading.value || noMore.value) return;
  loading.value = true;
  const { reqFunc, param } = getFuncAndParam();
  const res = await reqFunc(param, page.value, pageSize.value);
  if (res.data.code === statusCode.OK) {
    total.value = res.data.data.total;
    if (res.data.data.comments) {
      commentList.value = commentList.value.concat(res.data.data.comments);
    } else {
      if (page.value === 1) { // 第一页且没有数据
        commentList.value = [];
      }
      noMore.value = true;
    }
  }
  loading.value = false;
}

const getFuncAndParam = () => {
  if (commentType.value === '0') {
    return {
      reqFunc: getVideoCommentListAPI,
      param: videoId.value,
    }
  }

  return {
    reqFunc: getArticleCommentListAPI,
    param: articleId.value,
  }
}

const scrollLoad = () => {
  if (!loading.value) {
    page.value++;
    getReplyMsgList();
  }
}

// TODO: 增加删除评论

onBeforeMount(() => {
  getAllVideo();
  getAllArticle();
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

    // .title {
    //   font-size: 18px;
    //   margin: 0;
    // }

    .tabs {
      width: 150px;
      justify-content: space-between;
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