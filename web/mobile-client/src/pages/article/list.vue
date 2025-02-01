<template>
  <div class="home">
    <home-header class="home-header" @change-fold="changeMenuFold"></home-header>
    <div class="home-content">
      <div class="home-left" :class="menuFold ? 'home-left-fold' : ''">
        <home-sidebar class="home-sidebar" :fold="menuFold"></home-sidebar>
      </div>
      <div class="home-right">
        <div class="home-recommended">
          <ul class="article-list">
            <li class="article-item" v-for="(item, index) in articleList" :key="index">
              <nuxt-link class="content-wrapper" :to="`/article/${item.aid}`" target="_blank">
                <div class="content-main">
                  <div class="title-row" :to="`/article/${item.aid}`">{{ item.title }}</div>
                  <div class="abstract">{{ removeHtml(item.content) }}</div>
                  <div class="entry-footer">
                    <ul class="action-list">
                      <li class="item clicks">
                        <preview-open size="16" :strokeWidth="2" />
                        <span class="val">{{ item.clicks }}</span>
                      </li>
                    </ul>
                    <div class="entry-footer-tags">
                      <div v-if="item.tags" class="tag" v-for="tag in item.tags.split(',')">{{ tag }}</div>
                    </div>
                  </div>
                </div>
                <img v-if="item.cover" class="cover" :src="getResourceUrl(item.cover)" alt="封面">
              </nuxt-link>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from 'vue';
import HomeSidebar from '@/components/home-sidebar/index.vue';
import HomeHeader from "@/components/home-header/index.vue";
import { PreviewOpen } from '@icon-park/vue-next';
import { asyncGetRandomArticleAPI, getRandomArticleAPI } from '@/api/article';

definePageMeta({
  middleware: ['article']
})

const menuFold = ref(false);
const changeMenuFold = (val: boolean) => {
  menuFold.value = val;
}

// 获取分区
const size = 10;
const articleList = ref<ArticleType[]>([])
const { data } = await asyncGetRandomArticleAPI(size);
if ((data.value as any).code === statusCode.OK) {
  articleList.value = (data.value as any).data.articles;
}

const loading = ref(false);
const getArticleList = async () => {
  loading.value = true;
  const res = await getRandomArticleAPI(size);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.articles) {
      articleList.value = articleList.value.concat(res.data.data.articles);
    }
  }
  loading.value = false;
}

const lazyLoading = (e: Event) => {
  if((e.target as Document).location.pathname !== "/article/list") return;
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
  const clientHeight = document.documentElement.clientHeight;
  const scrollHeight = document.documentElement.scrollHeight;
  if (scrollTop + clientHeight >= scrollHeight) {
    if (!loading.value) {
      getArticleList();
    }
  }
}

onMounted(() => {
  window.addEventListener('scroll', lazyLoading, true);
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', lazyLoading, true);
})
</script>

<style lang="scss" scoped>
.home {
  width: 100%;
  min-width: 1000px;
  overflow: hidden;
  min-height: 100vh;
  background-color: #f2f3f5;


  .home-header {
    position: fixed;
    top: 0;
    width: 100%;
    z-index: 999;
    background-color: #fff;
  }
}

.home-content {
  display: flex;
  margin-top: 60px;
  z-index: 999;

  .home-left {
    width: 220px;
    transition: width .25s;

    .home-sidebar {
      position: fixed;
    }
  }

  .home-left-fold {
    width: 50px;
  }

  .home-right {
    flex: 1;
    margin-top: 12px;
  }
}

.home-recommended {
  margin-left: 16px;
  width: calc(100% - 32px);
}

.article-list {
  list-style: none;
  box-sizing: border-box;
  width: 100%;
  margin: 0;
  padding: 12px 0;
  background-color: #fff;

  .article-item {
    height: 110px;
  padding: 0 20px 0;
    box-sizing: border-box;
  }
}

.content-wrapper {
  display: flex;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(228, 230, 235, 0.5);
  width: 100%;
  margin-top: 2px;

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