<template>
  <div class="upload-article">
    <p class="title">专栏管理</p>
    <div class="article-box">
      <el-scrollbar>
        <ul class="article-list" v-infinite-scroll="scrollLoad">
          <li class="article-item" v-for="(item, index) in articleList" :key="index">
            <div class="content-wrapper">
              <div class="content-main">
                <div class="title-row">{{ item.title }}</div>
                <div class="abstract">{{ removeHtml(item.content) }}</div>
                <div class="entry-footer">
                  <ul class="action-list">
                    <li v-if="item.status === reviewCode.AUDIT_APPROVED" class="item clicks">
                      <preview-open size="16" :strokeWidth="2" />
                      <span class="val">{{ item.clicks }}</span>
                    </li>
                    <li v-else class="item review-result">
                      <span class="status" :style="`color: ${getStatusTextColor(item.status)}`">
                        {{ getStatusText(item.status) }}
                      </span>
                      <span class="status status-btn" v-if="item.status === reviewCode.REVIEW_FAILED"
                        @click="showReason(item.aid)">查看原因</span>
                    </li>
                    <li class="item more">
                      <el-dropdown>
                        <el-icon size="16">
                          <more-icon></more-icon>
                        </el-icon>
                        <template #dropdown>
                          <el-dropdown-menu>
                            <el-dropdown-item @click="modifyArticle(item.aid)">编辑</el-dropdown-item>
                            <el-dropdown-item @click="deleteArticle(item.aid)">删除稿件</el-dropdown-item>
                          </el-dropdown-menu>
                        </template>
                      </el-dropdown>
                    </li>
                  </ul>
                  <div class="entry-footer-tags">
                    <div v-if="item.tags" class="tag" v-for="tag in item.tags.split(',')">{{ tag }}</div>
                  </div>
                </div>
              </div>
              <img v-if="item.cover" class="cover" :src="getResourceUrl(item.cover)" alt="封面">
            </div>
          </li>
        </ul>
      </el-scrollbar>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from 'vue';
import { getUploadArticleAPI } from '@/api/article';
import { More as MoreIcon, PreviewOpen } from '@icon-park/vue-next';
import { getArticleReviewRecordAPI } from '@/api/revies';
import { getResourceUrl } from '@/utils/resource';

const page = ref(1);
const total = ref(0);
const pageSize = 8;
const noMore = ref(false);
const loading = ref(false);
const articleList = ref<ArticleType[]>([]);
const getUploadArticle = async () => {
  if (loading.value || noMore.value) return;
  loading.value = true;
  const res = await getUploadArticleAPI(page.value, pageSize);
  if (res.data.code === statusCode.OK) {
    total.value = res.data.data.total;
    if (res.data.data.articles) {
      articleList.value = articleList.value.concat(res.data.data.articles);
    } else {
      noMore.value = true;
    }
  }
  loading.value = false;
}

const scrollLoad = () => {
  if (!loading.value) {
    page.value++;
    getUploadArticle();
  }
}

const deleteArticle = async (aid: number) => {
  // TODO: 删除提醒
  // const res = await deleteVideoAPI(vid);
  // if (res.data.code === statusCode.OK) {
  //   getUploadArticle();
  // }
}

const getStatusText = (status: number) => {
  switch (status) {
    case reviewCode.SUBMIT_REVIEW:
    case reviewCode.WAITING_REVIEW:
      return "审核中"
    case reviewCode.REVIEW_FAILED:
      return "审核不通过"
  }
}

const getStatusTextColor = (status: number) => {
  switch (status) {
    case reviewCode.CREATED_VIDEO:
      return "#999"
    case reviewCode.SUBMIT_REVIEW:
    case reviewCode.WAITING_REVIEW:
      return "var(--primary-hover-color)"
    case reviewCode.REVIEW_FAILED:
      return "#f56c6c"
    case reviewCode.PROCESSING_FAIL:
      return "#f56c6c"
  }
}

const showReason = async (vid: number) => {
  const res = await getArticleReviewRecordAPI(vid);
  if (res.data.code === statusCode.OK) {
    ElMessageBox.alert(res.data.data.review.remark, '', {
      confirmButtonText: '确认',
    })
  }
}

//前往修改视频
const modifyArticle = (aid: number) => {
  navigateTo({ name: "upload-article", query: { aid: aid } });
}

onBeforeMount(() => {
  getUploadArticle();
})
</script>

<style lang="scss" scoped>
.upload-article {
  padding: 0 18px 0;
  height: 100%;
  box-sizing: border-box;
  background-color: #fff;

  .title {
    font-size: 18px;
    margin: 0;
    padding: 16px 0 10px;
  }

  .article-box {
    height: calc(100% - 60px);
  }

  .article-list {
    list-style: none;
    box-sizing: border-box;
    width: 100%;
    margin: 0;
    padding: 0;

    .article-item {
      height: 110px;
      box-sizing: border-box;
    }
  }
}

.content-wrapper {
  display: flex;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(228, 230, 235, 0.5);
  width: 100%;

  .content-main {
    flex: 1 1 auto;

    .title-row {
      display: flex;
      margin-bottom: 2px;
      font-weight: 600;
      font-size: 16px;
      line-height: 24px;
      color: #252933;
      width: 100%;
      display: -webkit-box;
      overflow: hidden;
      text-overflow: ellipsis;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 1;
    }

    .abstract {
      min-height: 44px;
      margin-bottom: 4px;
      font-weight: 400;
      color: #8a919f;
      font-size: 13px;
      line-height: 22px;
      display: -webkit-box;
      overflow: hidden;
      text-overflow: ellipsis;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 2;
    }

    .entry-footer {
      display: flex;
      align-items: center;
      justify-content: space-between;
      flex-wrap: wrap;

      .action-list {
        display: flex;
        align-items: center;
        list-style: none;
        margin: 0;
        padding: 0;

        .item {
          display: flex;
          align-items: center;
          position: relative;
          margin-right: 24px;
          font-size: 13px;
          line-height: 20px;
          color: #8a919f;

          span {
            height: 16px;
            line-height: 16px;
          }

          .val {
            margin-left: 4px;
          }
        }

        .more {
          cursor: pointer;
        }
      }

      .entry-footer-tags {
        display: flex;
        align-items: center;

        .tag {
          background-color: #f2f3f5;
          padding: 0 6px;
          border-radius: 2px;
          max-width: 76px;
          box-sizing: border-box;
          margin-left: 6px;
          color: #8a919f;
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
          min-height: 18px;
          line-height: 18px;
          font-size: 12px;
        }
      }

      .review-result {
        font-size: 12px;

        .status {
          color: var(--primary-hover-color);
        }

        .status-btn {
          cursor: pointer;
          margin-left: 6px;
        }
      }
    }
  }

  .cover {
    flex: 0 0 auto;
    width: 108px;
    height: 72px;
    margin: 12px 0 0 24px;
    background-color: #fff;
    border-radius: 4px;
    border: 1px solid rgba(228, 230, 235, 0.5);
  }
}
</style>