<template>
  <div class="article">
    <header-bar class="header"></header-bar>
    <div class="article-container">
      <div class="article-content">
        <div class="title-container">
          <h1 class="title">{{ articleInfo?.title }}</h1>
          <div class="article-read-panel">
            <div class="article-read-info">
              <span class="publish-text">{{ articleInfo ? formatRelativeTime(articleInfo.createdAt) : '' }}</span>
              <span>{{ articleInfo?.clicks }}浏览</span> ·
              <!-- <span>{{ likeCount }}喜欢</span> · -->
              <span>{{ commentCount }}评论</span>
            </div>
          </div>
        </div>
        <div class="article-up-info">
          <div class="up-left">
            <div class="avatar-container">
              <common-avatar :size="44"></common-avatar>
            </div>
            <div class="info-container">
              <span class="up-name">{{ articleInfo?.author.name }}</span>
              <span class="up-sign">{{ articleInfo?.author.sign }}</span>
            </div>
          </div>
          <el-button class="follow-btn" type="primary" @click="followBtnClick">{{ btnText }}</el-button>
        </div>
        <div class="title-line"></div>
        <el-skeleton v-if="loading" :rows="10" />
        <client-only>
          <text-editor :content="articleInfo?.content"></text-editor>
        </client-only>
        <div class="content-html seowhy" v-html="getSanitizeHtml(articleInfo?.content)"></div>
      </div>
      <!-- 评论区 -->
      <div class="comment-list">
        <comment-list v-if="articleInfo" :aid="articleInfo.aid" @update-count="updateCommentCount"></comment-list>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";
import HeaderBar from "@/components/header-bar/index.vue";
import { asyncGetArticleInfoAPI } from "@/api/article";
import TextEditor from './components/TextEditor.vue';
import "@/assets/styles/article-html.scss";
import DOMPurify from 'isomorphic-dompurify';
import CommonAvatar from "@/components/common-avatar/index.vue";
import CommentList from "./components/CommentList.vue";
import { followAPI, getUserRelationAPI, unfollowAPI } from "@/api/relation";


const route = useRoute();
const router = useRouter();

const loading = ref(true);

// 获取文章信息
const articleInfo = ref<ArticleType>();
const articleId = route.params.id.toString();
const { data } = await asyncGetArticleInfoAPI(articleId);
if ((data.value as any).code === statusCode.OK) {
  articleInfo.value = (data.value as any).data.article as ArticleType;
} else {
  // TODO: 处理视频信息不存在
}

const getSanitizeHtml = (val: string = "") => {
  return DOMPurify.sanitize(val);
}

const likeCount = ref(0);
const commentCount = ref(0);
const updateCommentCount = (val: number) => {
  commentCount.value = val;
}

const relation = ref(relationCode.NOT_FOLLOWING);
const getUserRelation = async () => {
  if (articleInfo.value?.author.uid) {
    const res = await getUserRelationAPI(articleInfo.value.author.uid);
    if (res.data.code === statusCode.OK) {
      relation.value = res.data.data.relation;
    }
  }
}

const btnText = computed(() => {
  switch (relation.value) {
    case relationCode.FOLLOWED:
      return "已关注";
    case relationCode.NOT_FOLLOWING:
      return "关注";
    case relationCode.MUTUAL_FANS:
      return "已互粉";
  }
})

const followBtnClick = async () => {
  if (articleInfo.value?.author.uid) {
    const reqFunc = relation.value === relationCode.NOT_FOLLOWING ? followAPI : unfollowAPI;
    const res = await reqFunc(articleInfo.value.author.uid);
    if (res.data.code === statusCode.OK) {
      getUserRelation();
    } else {
      ElMessage.error(res.data.msg);
    }
  }
}

onBeforeMount(() => {
  getUserRelation();
  loading.value = false;
})
</script>

<style lang="scss" scoped>
.article {
  position: relative;
  min-height: 100vh;
  padding-top: 76px;
  background-color: #f6f7f9;

  .header {
    position: fixed;
    left: 0;
    top: 0;
  }
}

.article-container {
  width: 900px;
  margin: 0 auto;
}

.article-content {
  background-color: #fff;
  border-radius: 4px;
  padding: 20px 40px 40px;


  .title-container {
    padding: 0 20px;

    .title {
      min-height: 39px;
      font-size: 28px;
      color: #222;
      margin-bottom: 16px;
      font-weight: 700;
      line-height: 1.4;
    }

    .article-read-panel {
      margin-bottom: 20px;
      color: #999;

      .article-read-info {
        font-weight: 400;
        font-size: 13px;

        .publish-text {
          display: inline-block;
          margin-right: 10px;
        }

        span {
          margin: 0 3px;
        }
      }
    }
  }

  .article-up-info {
    display: flex;
    align-items: center;
    padding: 0 20px;
    height: 44px;
    justify-content: space-between;
    margin: 0 auto 20px;

    .up-left {
      display: flex;
      align-items: center;

      .avatar-container {
        flex-shrink: 0;
        display: block;
        position: relative;
        width: 44px;
        height: 44px;
        margin-right: 12px;
      }

      .info-container {
        .up-name {
          font-size: 14px;
          max-width: 150px;
          color: #212121;
          line-height: 23px;
          height: 23px;
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
          display: inline-block;
          vertical-align: middle;
        }

        .up-sign {
          display: block;
          font-size: 13px;
          color: #999;
          line-height: 20px;
        }
      }
    }

    .follow-btn {
      width: 128px;
      height: 32px;
    }
  }
}

.title-line {
  height: 1px;
  width: 100%;
  background: hsla(0, 0%, 59.2%, .21);
  margin-bottom: 20px;
}

.comment-list {
  margin-top: 20px;
  padding: 20px;
  background-color: #fff;
}
</style>