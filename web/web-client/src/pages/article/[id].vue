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
              <span>{{ likeCount }}喜欢</span> ·
              <span>{{ commentCount }}评论</span>
            </div>
          </div>
        </div>
        <!-- 作者信息 -->
        <author-card v-if="articleInfo" :author="articleInfo.author"></author-card>
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
      <!-- 点赞收藏信息 -->
      <archive-info v-if="articleInfo" :aid="articleInfo.aid" @stat-change="statChange"
        @scroll-to-comment="scrollToComment"></archive-info>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";
import HeaderBar from "@/components/header-bar/index.vue";
import { asyncGetArticleInfoAPI } from "@/api/article";
import TextEditor from './components/TextEditor.vue';
import "@/assets/styles/article-html.scss";
import { scrollToViewCenter } from "@/utils/scroll";
import DOMPurify from 'isomorphic-dompurify';
import CommentList from "./components/CommentList.vue";
import AuthorCard from './components/AuthorCard.vue';
import ArchiveInfo from './components/ArchiveInfo.vue';
import { getArticleArchiveStatAPI } from "@/api/archive";

definePageMeta({
  middleware: ['article']
})

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
  // 处理文章信息不存在
  navigateTo('/404');
}

const getSanitizeHtml = (val: string = "") => {
  return DOMPurify.sanitize(val);
}

const likeCount = ref(0);
const commentCount = ref(0);
const updateCommentCount = (val: number) => {
  commentCount.value = val;
}

//获取点赞收藏关注信息
const getArchiveStat = async () => {
  const res = await getArticleArchiveStatAPI(articleId);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.stat) {
      likeCount.value = res.data.data.stat.like;
    }
  }
}

const statChange = (type: string, val: number) => {
  if (type === "like") {
    likeCount.value += val;
  }
}

const scrollToComment = () => {
  const el = document.getElementById("comment-title");
  if (el) scrollToViewCenter(el);
}

onBeforeMount(() => {
  getArchiveStat();
  nextTick(() => {
    loading.value = false;
  })
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
  position: relative;

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

@media (max-width: 1400px) {
  .article-container {
    width: 700px;
  }
}
</style>