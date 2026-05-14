<template>
  <div class="home">
    <home-header class="home-header" @change-fold="changeMenuFold"></home-header>
    <div class="home-content">
      <div class="home-left" :class="menuFold ? 'home-left-fold' : ''">
        <home-sidebar class="home-sidebar" :fold="menuFold"></home-sidebar>
      </div>
      <div class="home-right">
        <div class="home-recommended">
          <!-- 空状态提示 -->
          <div v-if="!articleList || articleList.length === 0" class="empty-state">
            <el-empty description="暂无专栏文章" />
          </div>
          <!-- 文章列表 -->
          <ul v-else class="article-list">
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
                      <div v-if="item.tags" class="tag" v-for="tag in item.tags.split(',')" :key="tag">{{ tag }}</div>
                    </div>
                  </div>
                </div>
                <img v-if="item.cover" class="cover" :src="getResourceUrl(item.cover)" alt="封面">
              </nuxt-link>
            </li>
          </ul>
          <!-- 加载提示 -->
          <div v-if="loading" class="loading-tip">加载中...</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from 'vue';
import { PreviewOpen } from '@icon-park/vue-next';
import { asyncGetRandomArticleAPI, getRandomArticleAPI } from '@/api/article';
import { throttle } from "@/utils/debounce";

definePageMeta({
  middleware: ['article']
})

const menuFoldCookie = useCookie<boolean>('menu-fold-state', { default: () => false });
const menuFold = ref(menuFoldCookie.value);
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

const throttledLoading = throttle(lazyLoading, 150);

onMounted(() => {
  window.addEventListener('scroll', throttledLoading, true);
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', throttledLoading, true);
})
</script>

<style lang="scss" scoped>
.home {
  width: 100%;
  min-width: 1000px;
  overflow: hidden;
  min-height: 100vh;
  background-color: var(--bg-page);


  .home-header {
    position: fixed;
    top: 0;
    width: 100%;
    z-index: 999;
    background-color: var(--bg-elev-1);
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

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  background-color: var(--bg-elev-1);
  border-radius: 4px;
}

.loading-tip {
  text-align: center;
  padding: 20px 0;
  color: var(--font-primary-3);
  font-size: 14px;
}

.article-list {
  list-style: none;
  box-sizing: border-box;
  width: 100%;
  margin: 0;
  padding: 12px 0;
  background-color: var(--bg-elev-1);

  .article-item {
    height: 110px;
    padding: 0 20px 0;
    box-sizing: border-box;
  }
}

.content-wrapper {
  display: flex;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-color);
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
      color: var(--font-primary-1);
      width: 100%;
      display: -webkit-box;
      overflow: hidden;
      text-overflow: ellipsis;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 1;
      line-clamp: 1;
    }

    .abstract {
      min-height: 44px;
      margin-bottom: 4px;
      font-weight: 400;
      color: var(--font-primary-3);
      font-size: 13px;
      line-height: 22px;
      display: -webkit-box;
      overflow: hidden;
      text-overflow: ellipsis;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 2;
      line-clamp: 2;
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
          color: var(--font-primary-3);

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
          background-color: var(--fill-1);
          padding: 0 6px;
          border-radius: 2px;
          max-width: 76px;
          box-sizing: border-box;
          margin-left: 6px;
          color: var(--font-primary-3);
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
    background-color: var(--bg-elev-1);
    border-radius: 4px;
    border: 1px solid var(--border-color);
  }
}
</style>