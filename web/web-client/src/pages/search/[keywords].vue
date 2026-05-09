<template>
  <div class="search-page">
    <header-bar></header-bar>
    <div class="search">
      <div class="search-form">
        <input
          class="input"
          v-model="searchKeywords"
          @keydown.enter="submitSearch"
        >
        <button class="btn" type="button" @click="submitSearch">
          <search-icon class="icon" size="16" />
        </button>
      </div>
    </div>

    <div class="filter-bar">
      <div class="chips">
        <el-tag v-if="timeRange !== 'all'" class="chip" type="info" effect="plain" closable @close="clearTimeRange">
          {{ timeRangeLabel }}
        </el-tag>
        <el-tag v-if="sort !== 'relevance'" class="chip" type="info" effect="plain" closable @close="clearSort">
          {{ sortLabel }}
        </el-tag>
      </div>
      <el-button class="filter-btn" text @click="openFilter">
        筛选
      </el-button>
    </div>

    <el-drawer v-model="filterOpen" title="搜索过滤条件" size="420px">
      <div class="filter-panel">
        <div class="filter-group">
          <div class="filter-title">排序依据</div>
          <el-radio-group v-model="sort" class="filter-radio">
            <el-radio-button label="relevance">相关性</el-radio-button>
            <el-radio-button label="newest">上传日期</el-radio-button>
            <el-radio-button label="most_viewed">观看次数</el-radio-button>
          </el-radio-group>
        </div>

        <div class="filter-group">
          <div class="filter-title">上传日期</div>
          <el-radio-group v-model="timeRange" class="filter-radio">
            <el-radio-button label="all">不限</el-radio-button>
            <el-radio-button label="24h">今天</el-radio-button>
            <el-radio-button label="week">本周</el-radio-button>
            <el-radio-button label="month">本月</el-radio-button>
            <el-radio-button label="year">今年</el-radio-button>
          </el-radio-group>
        </div>

        <div class="filter-actions">
          <el-button @click="resetFilters">重置</el-button>
          <el-button type="primary" @click="applyFilters">应用</el-button>
        </div>
      </div>
    </el-drawer>

    <el-tabs v-model="activeTab" class="search-tabs" @tab-change="onTabChange">
      <el-tab-pane label="视频" name="video">
        <div class="pgc-top-block" v-if="pgcList.length > 0">
          <div class="section-title">相关影视</div>
          <ul class="pgc-search-list">
            <li v-for="item in pgcList.slice(0, 3)" :key="item.pgc_id" class="pgc-search-item">
              <div class="pgc-link" @click="openPGC(item)">
                <img v-if="item.cover" class="cover" :src="getResourceUrl(item.cover)" alt="封面">
                <div class="meta">
                  <div class="title-row">{{ item.title }}</div>
                  <div class="desc">{{ item.desc || '暂无简介' }}</div>
                  <div class="entry-footer">
                    <span v-if="item.rating">评分 {{ item.rating }}</span>
                    <span v-if="item.current_episodes">更新至 {{ item.current_episodes }} 集</span>
                    <span v-if="item.area">{{ formatAreaName(item.area) }}</span>
                    <span v-if="item.year">{{ item.year }}</span>
                  </div>
                </div>
              </div>
            </li>
          </ul>
        </div>
        <div class="section-title">相关视频</div>
        <div class="card-list">
          <video-item
            v-for="item in videoList"
            :key="item.vid"
            :info="item"
            :keywords="routeKeywords"
          />
        </div>
        <p v-if="hydrated && !videoLoading && videoList.length === 0 && didSearch" class="empty-hint">暂无相关视频</p>
      </el-tab-pane>
      <el-tab-pane label="专栏" name="article">
        <ul class="article-search-list">
          <li v-for="item in articleList" :key="item.aid" class="article-search-item">
            <nuxt-link class="content-wrapper" :to="`/article/${item.aid}`" target="_blank">
              <div class="content-main">
                <div class="title-row">{{ item.title }}</div>
                <div class="abstract">{{ removeHtml(item.content) }}</div>
                <div class="entry-footer">
                  <span class="clicks">{{ item.clicks }} 阅读</span>
                  <span v-if="item.tags" class="tags">{{ item.tags }}</span>
                </div>
              </div>
              <img v-if="item.cover" class="cover" :src="getResourceUrl(item.cover)" alt="封面">
            </nuxt-link>
          </li>
        </ul>
        <p v-if="hydrated && !articleLoading && articleList.length === 0 && didSearch" class="empty-hint">暂无相关专栏</p>
      </el-tab-pane>
      <el-tab-pane label="UP主" name="user">
        <ul class="user-search-list">
          <li v-for="u in userList" :key="u.uid" class="user-search-item">
            <nuxt-link class="user-link" :to="`/user/${u.uid}`" target="_blank">
              <common-avatar class="avatar" :size="48" :url="u.avatar" />
              <div class="meta">
                <div class="name">{{ u.name }}</div>
                <div class="sign">{{ u.sign || '暂无签名' }}</div>
                <div class="fans" v-if="u.fans != null">粉丝 {{ u.fans }}</div>
              </div>
            </nuxt-link>
          </li>
        </ul>
        <p v-if="hydrated && !userLoading && userList.length === 0 && didSearch" class="empty-hint">暂无相关用户</p>
      </el-tab-pane>
      <el-tab-pane label="影视" name="pgc">
        <ul class="pgc-search-list">
          <li v-for="item in pgcList" :key="item.pgc_id" class="pgc-search-item">
            <div class="pgc-link" @click="openPGC(item)">
              <img v-if="item.cover" class="cover" :src="getResourceUrl(item.cover)" alt="封面">
              <div class="meta">
                <div class="title-row">{{ item.title }}</div>
                <div class="desc">{{ item.desc || '暂无简介' }}</div>
                <div class="entry-footer">
                  <span v-if="item.rating">评分 {{ item.rating }}</span>
                  <span v-if="item.current_episodes">更新至 {{ item.current_episodes }} 集</span>
                  <span v-if="item.area">{{ formatAreaName(item.area) }}</span>
                  <span v-if="item.year">{{ item.year }}</span>
                </div>
              </div>
            </div>
          </li>
        </ul>
        <p v-if="hydrated && !pgcLoading && pgcList.length === 0 && didSearch" class="empty-hint">暂无相关影视</p>
      </el-tab-pane>
    </el-tabs>

    <p v-if="hydrated && footerLoading" class="loading-more">加载中...</p>
    <p v-else-if="hydrated && currentNoMore && currentListNonEmpty" class="loading-more muted">没有更多了</p>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch, computed } from 'vue';
import { getVideoInfoAPI, searchVideoAPI } from '@/api/video';
import { searchArticleAPI } from '@/api/article';
import { searchUserAPI } from '@/api/user';
import { getPGCDetailWithEpisodesAPI, searchPGCAPI } from '@/api/pgc';
import VideoItem from './components/VideoItem.vue';
import { Search as SearchIcon } from '@icon-park/vue-next';
import { removeHtml } from '@/utils/format';
import { getResourceUrl } from '@/utils/resource';

const route = useRoute();
const router = useRouter();

const pageSize = 15;
const searchKeywords = ref('');
const routeKeywords = ref('');
const activeTab = ref<'video' | 'article' | 'user' | 'pgc'>('video');
const didSearch = ref(false);
const hydrated = ref(false);

const filterOpen = ref(false);
const sort = ref<'relevance' | 'newest' | 'most_viewed'>('relevance');
const timeRange = ref<'all' | '24h' | 'week' | 'month' | 'year'>('all');

const videoList = ref<VideoType[]>([]);
const videoPage = ref(1);
const videoNoMore = ref(false);
const videoLoading = ref(false);

const articleList = ref<ArticleType[]>([]);
const articlePage = ref(1);
const articleNoMore = ref(false);
const articleLoading = ref(false);

const userList = ref<UserInfoType[]>([]);
const userPage = ref(1);
const userNoMore = ref(false);
const userLoading = ref(false);

const pgcList = ref<PGCRecommendItem[]>([]);
const pgcPage = ref(1);
const pgcNoMore = ref(false);
const pgcLoading = ref(false);
const areaCodeMap: Record<string, string> = {
  CN: '中国大陆',
  JP: '日本',
  HK: '中国香港',
  TW: '中国台湾',
  KR: '韩国',
  US: '美国',
};
const formatAreaName = (raw: unknown) => {
  const code = String(raw ?? '').trim().toUpperCase();
  return areaCodeMap[code] || String(raw ?? '').trim();
};

const sortLabel = computed(() => {
  switch (sort.value) {
    case 'newest':
      return '上传日期'
    case 'most_viewed':
      return '观看次数'
    default:
      return '相关性'
  }
});

const timeRangeLabel = computed(() => {
  switch (timeRange.value) {
    case '24h':
      return '今天'
    case 'week':
      return '本周'
    case 'month':
      return '本月'
    case 'year':
      return '今年'
    default:
      return '不限'
  }
});

const readFiltersFromRoute = () => {
  const q: any = route.query || {};
  const s = String(q.sort || 'relevance');
  const t = String(q.timeRange || 'all');
  if (s === 'newest' || s === 'most_viewed' || s === 'relevance') {
    sort.value = s;
  } else {
    sort.value = 'relevance';
  }
  if (t === '24h' || t === 'week' || t === 'month' || t === 'year' || t === 'all') {
    timeRange.value = t;
  } else {
    timeRange.value = 'all';
  }
};

const syncKeywordsFromRoute = () => {
  const raw = route.params.keywords?.toString() ?? '';
  try {
    routeKeywords.value = decodeURIComponent(raw);
  } catch {
    routeKeywords.value = raw;
  }
  searchKeywords.value = routeKeywords.value;
};

const resetAllLists = () => {
  videoList.value = [];
  videoPage.value = 1;
  videoNoMore.value = false;
  articleList.value = [];
  articlePage.value = 1;
  articleNoMore.value = false;
  userList.value = [];
  userPage.value = 1;
  userNoMore.value = false;
  pgcList.value = [];
  pgcPage.value = 1;
  pgcNoMore.value = false;
  didSearch.value = false;
};

const openFilter = () => {
  filterOpen.value = true;
};

const resetFilters = () => {
  sort.value = 'relevance';
  timeRange.value = 'all';
};

const clearTimeRange = async () => {
  timeRange.value = 'all';
  await applyFilters();
};

const clearSort = async () => {
  sort.value = 'relevance';
  await applyFilters();
};

const loadVideos = async (init: boolean) => {
  if (videoLoading.value) return;
  videoLoading.value = true;
  if (init) {
    videoPage.value = 1;
    videoList.value = [];
    videoNoMore.value = false;
  }
  const res = await searchVideoAPI({
    page: videoPage.value,
    pageSize,
    keywords: routeKeywords.value,
    sort: sort.value,
    timeRange: timeRange.value,
  });
  if (res.data.code === statusCode.OK && res.data.data.videos) {
    const chunk = res.data.data.videos;
    videoList.value.push(...chunk);
    if (chunk.length < pageSize) videoNoMore.value = true;
  } else {
    videoNoMore.value = true;
    if (init) ElMessage.error(res.data.msg || '获取失败');
  }
  videoLoading.value = false;
  didSearch.value = true;
  if (init) {
    await loadPGC(true);
  }
};

const loadArticles = async (init: boolean) => {
  if (articleLoading.value) return;
  articleLoading.value = true;
  if (init) {
    articlePage.value = 1;
    articleList.value = [];
    articleNoMore.value = false;
  }
  const res = await searchArticleAPI({
    page: articlePage.value,
    pageSize,
    keywords: routeKeywords.value,
    sort: sort.value,
    timeRange: timeRange.value,
  });
  if (res.data.code === statusCode.OK && res.data.data.articles) {
    const chunk = res.data.data.articles;
    articleList.value.push(...chunk);
    if (chunk.length < pageSize) articleNoMore.value = true;
  } else {
    articleNoMore.value = true;
    if (init) ElMessage.error(res.data.msg || '获取失败');
  }
  articleLoading.value = false;
  didSearch.value = true;
};

const loadUsers = async (init: boolean) => {
  if (userLoading.value) return;
  userLoading.value = true;
  if (init) {
    userPage.value = 1;
    userList.value = [];
    userNoMore.value = false;
  }
  const res = await searchUserAPI({
    page: userPage.value,
    pageSize,
    keywords: routeKeywords.value,
    sort: sort.value,
    timeRange: timeRange.value,
  });
  if (res.data.code === statusCode.OK && res.data.data.users) {
    const chunk = res.data.data.users;
    userList.value.push(...chunk);
    if (chunk.length < pageSize) userNoMore.value = true;
  } else {
    userNoMore.value = true;
    if (init) ElMessage.error(res.data.msg || '获取失败');
  }
  userLoading.value = false;
  didSearch.value = true;
};

const loadPGC = async (init: boolean) => {
  if (pgcLoading.value) return;
  pgcLoading.value = true;
  if (init) {
    pgcPage.value = 1;
    pgcList.value = [];
    pgcNoMore.value = false;
  }
  const res = await searchPGCAPI({
    page: pgcPage.value,
    pageSize,
    keyword: routeKeywords.value,
    pgcType: 0,
  });
  if (res.data.code === statusCode.OK && res.data.data.list) {
    const chunk = res.data.data.list;
    pgcList.value.push(...chunk);
    if (chunk.length < pageSize) pgcNoMore.value = true;
  } else {
    pgcNoMore.value = true;
    if (init) ElMessage.error(res.data.msg || '获取失败');
  }
  pgcLoading.value = false;
  didSearch.value = true;
};

const loadCurrentTab = (init: boolean) => {
  if (activeTab.value === 'video') return loadVideos(init);
  if (activeTab.value === 'article') return loadArticles(init);
  if (activeTab.value === 'user') return loadUsers(init);
  return loadPGC(init);
};

const submitSearch = () => {
  const k = searchKeywords.value.trim();
  if (!k) {
    ElMessage.warning('请输入关键词');
    return;
  }
  router.push({
    path: `/search/${encodeURIComponent(k)}`,
    query: {
      sort: sort.value,
      timeRange: timeRange.value,
    },
  });
};

const onTabChange = (name: string) => {
  if (!routeKeywords.value.trim()) return;
  if (name === 'article' && articleList.value.length === 0) loadArticles(true);
  if (name === 'user' && userList.value.length === 0) loadUsers(true);
  if (name === 'pgc' && pgcList.value.length === 0) loadPGC(true);
};

const applyFilters = async () => {
  filterOpen.value = false;
  await router.replace({
    query: {
      ...route.query,
      sort: sort.value,
      timeRange: timeRange.value,
    },
  });
  resetAllLists();
  await loadCurrentTab(true);
};

const footerLoading = computed(() => {
  if (activeTab.value === 'video') return videoLoading.value;
  if (activeTab.value === 'article') return articleLoading.value;
  if (activeTab.value === 'user') return userLoading.value;
  return pgcLoading.value;
});

const currentNoMore = computed(() => {
  if (activeTab.value === 'video') return videoNoMore.value;
  if (activeTab.value === 'article') return articleNoMore.value;
  if (activeTab.value === 'user') return userNoMore.value;
  return pgcNoMore.value;
});

const currentListNonEmpty = computed(() => {
  if (activeTab.value === 'video') return videoList.value.length > 0;
  if (activeTab.value === 'article') return articleList.value.length > 0;
  if (activeTab.value === 'user') return userList.value.length > 0;
  return pgcList.value.length > 0;
});

const lazyLoading = () => {
  if (currentNoMore.value || footerLoading.value) return;
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
  const clientHeight = document.documentElement.clientHeight;
  const scrollHeight = document.documentElement.scrollHeight;
  if (scrollTop + clientHeight < scrollHeight - 10) return;

  if (activeTab.value === 'video') {
    videoPage.value++;
    loadVideos(false);
  } else if (activeTab.value === 'article') {
    articlePage.value++;
    loadArticles(false);
  } else {
    if (activeTab.value === 'user') {
      userPage.value++;
      loadUsers(false);
    } else {
      pgcPage.value++;
      loadPGC(false);
    }
  }
};

const openPGC = async (item: PGCRecommendItem) => {
  const toCleanVid = (raw: unknown): number => {
    const n = Number(String(raw ?? '').trim());
    return Number.isFinite(n) && n > 0 ? Math.floor(n) : 0;
  };

  const latestEpId = Number(item.latest_ep_id || 0);
  if (Number.isFinite(latestEpId) && latestEpId > 0) {
    await navigateTo(`/watch?ep=${Math.floor(latestEpId)}&mode=pgc`);
    return;
  }

  const candidates: number[] = [];
  const latestVid = toCleanVid(item.latest_vid);
  if (latestVid > 0) candidates.push(latestVid);

  // 兜底：无 latest_vid 或 latest_vid 不可播时，查剧集详情并依次尝试
  try {
    const res = await getPGCDetailWithEpisodesAPI(item.pgc_id);
    const eps = res?.data?.data?.episodes || [];
    for (const ep of eps) {
      const epId = Number(ep?.ep_id || ep?.id || 0);
      if (Number.isFinite(epId) && epId > 0) {
        await navigateTo(`/watch?ep=${Math.floor(epId)}&mode=pgc`);
        return;
      }
      const vid = toCleanVid(ep?.vid);
      if (vid > 0 && !candidates.includes(vid)) {
        candidates.push(vid);
      }
    }
  } catch {}

  for (const vid of candidates) {
    try {
      const infoRes = await getVideoInfoAPI(vid);
      if (infoRes?.data?.code === statusCode.OK) {
        await navigateTo(`/watch?v=${vid}&mode=pgc`);
        return;
      }
    } catch {}
  }

  ElMessage.warning('该影视暂不可播放');
};

watch(
  () => route.params.keywords,
  async () => {
    syncKeywordsFromRoute();
    readFiltersFromRoute();
    resetAllLists();
    activeTab.value = 'video';
    await loadVideos(true);
  },
);

onMounted(() => {
  hydrated.value = true;
  syncKeywordsFromRoute();
  readFiltersFromRoute();
  loadVideos(true);
  window.addEventListener('scroll', lazyLoading, true);
});

onBeforeUnmount(() => {
  window.removeEventListener('scroll', lazyLoading, true);
});
</script>

<style lang="scss" scoped>
.search-page {
  min-height: 100vh;
  padding-bottom: 40px;
  background: var(--bg-page);
}

.search {
  width: 100%;
  margin: 30px auto;
  max-width: 500px;

  .search-form {
    position: relative;

    .input {
      border: 1px solid #e5e5e5;
      outline: none;
      padding: 8px 30px 8px 11px;
      height: 36px;
      font-size: 12px;
      line-height: 14px;
      border-radius: 18px;
      width: 100%;
      max-width: 500px;
      vertical-align: top;
      color: var(--font-primary-1);
      box-sizing: border-box;
    }

    .btn {
      position: absolute;
      top: 0;
      right: 10px;
      border: none;
      width: 20px;
      height: 36px;
      line-height: 36px;
      font-size: 14px;
      vertical-align: top;
      background: transparent;
      padding: 0;
      cursor: pointer;

      .icon {
        display: block;
        margin-top: 3px;
        width: 20px;
        height: 36px;
        color: var(--font-primary-5);
      }
    }
  }
}

.search-tabs {
  width: 90%;
  margin: 0 auto;
}

.section-title {
  margin-top: 12px;
  margin-bottom: 10px;
  font-size: 16px;
  font-weight: 600;
  color: var(--font-primary-1);
}

.pgc-top-block {
  margin-top: 8px;
}

.filter-bar {
  width: 90%;
  margin: -10px auto 8px;
  display: flex;
  align-items: center;
  justify-content: space-between;

  .chips {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .chip {
    border-radius: 999px;
  }
}

.filter-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding-right: 6px;

  .filter-group {
    .filter-title {
      font-size: 13px;
      font-weight: 600;
      color: var(--font-primary-2);
      margin-bottom: 10px;
    }
    .filter-radio {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }
  }

  .filter-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 6px;
  }
}

.card-list {
  display: grid;
  column-gap: 16px;
  grid-template-columns: repeat(5, 1fr);
  margin-top: 16px;
}

.article-search-list {
  list-style: none;
  padding: 0;
  margin: 16px 0 0;
}

.article-search-item {
  margin-bottom: 16px;

  .content-wrapper {
    display: flex;
    gap: 16px;
    padding: 16px;
    background: var(--bg-elev-1);
    border-radius: 8px;
    text-decoration: none;
    color: inherit;
  }

  .content-main {
    flex: 1;
    min-width: 0;
  }

  .title-row {
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 8px;
  }

  .abstract {
    font-size: 13px;
    color: var(--font-primary-3);
    line-height: 1.5;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .entry-footer {
    margin-top: 8px;
    font-size: 12px;
    color: var(--font-primary-5);
    display: flex;
    gap: 12px;
  }

  .cover {
    width: 160px;
    height: 100px;
    object-fit: cover;
    border-radius: 6px;
    flex-shrink: 0;
  }
}

.user-search-list {
  list-style: none;
  padding: 0;
  margin: 16px 0 0;
}

.user-search-item {
  margin-bottom: 12px;

  .user-link {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 12px 16px;
    background: var(--bg-elev-1);
    border-radius: 8px;
    text-decoration: none;
    color: inherit;
  }

  .meta {
    flex: 1;
    min-width: 0;
  }

  .name {
    font-weight: 500;
    font-size: 15px;
  }

  .sign {
    font-size: 13px;
    color: var(--font-primary-3);
    margin-top: 4px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .fans {
    font-size: 12px;
    color: var(--font-primary-5);
    margin-top: 4px;
  }
}

.pgc-search-list {
  list-style: none;
  padding: 0;
  margin: 16px 0 0;
}

.pgc-search-item {
  margin-bottom: 12px;

  .pgc-link {
    display: flex;
    gap: 14px;
    padding: 12px 14px;
    background: var(--bg-elev-1);
    border-radius: 8px;
    text-decoration: none;
    color: inherit;
    cursor: pointer;
  }

  .cover {
    width: 180px;
    height: 102px;
    object-fit: cover;
    border-radius: 6px;
    flex-shrink: 0;
    background: rgba(0, 0, 0, .15);
  }

  .meta {
    flex: 1;
    min-width: 0;
  }

  .title-row {
    font-size: 16px;
    font-weight: 500;
    color: var(--font-primary-1);
  }

  .desc {
    margin-top: 8px;
    font-size: 13px;
    line-height: 1.5;
    color: var(--font-primary-3);
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .entry-footer {
    margin-top: 8px;
    font-size: 12px;
    color: var(--font-primary-5);
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
  }
}

.empty-hint {
  text-align: center;
  color: var(--font-primary-5);
  padding: 24px;
}

.loading-more {
  text-align: center;
  padding: 16px;
  &.muted {
    color: var(--font-primary-5);
    font-size: 13px;
  }
}
</style>
