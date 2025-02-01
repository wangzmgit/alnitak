<template>
  <div class="history">
    <p class="history-title">历史记录</p>
    <div class="history-list">
      <history-video-list :video-list="historyList.today"></history-video-list>
      <span v-if="historyList.today.length && historyList.yesterday.length" class="date-title">昨天</span>
      <history-video-list :video-list="historyList.yesterday"></history-video-list>
      <span v-if="(historyList.today.length || historyList.yesterday.length) && historyList.earlier.length"
        class="date-title">更早</span>
      <history-video-list :video-list="historyList.earlier"></history-video-list>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import { getHistoryVideoAPI } from '@/api/history';
import HistoryVideoList from './components/HistoryVideoList.vue';


const historyList = reactive<{
  today: HistoryVideoType[];
  yesterday: HistoryVideoType[];
  earlier: HistoryVideoType[];
}>({
  today: [],
  yesterday: [],
  earlier: [],
})
const historyParams = reactive({
  page: 1,
  pageSize: 15,
  noMore: false,
  loading: false,
})
const getHistoryList = async () => {
  historyParams.loading = true;
  const res = await getHistoryVideoAPI(historyParams.page, historyParams.pageSize);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.videos.length < historyParams.pageSize) {
      historyParams.noMore = true;
    }

    res.data.data.videos.forEach((item: HistoryVideoType) => {
      const { group, time } = relativeDate(item.updatedAt);
      item.viewingDate = time;
      // @ts-ignore
      historyList[group].push(item);
    });
  }

  historyParams.loading = false;
}

onBeforeMount(() => {
  getHistoryList();
})
</script>

<style lang="scss" scoped>
.history {
  padding-left: 12px;

  .history-title {
    font-size: 18px;
    margin-top: 20px;
    padding-left: 10px;
  }

  .date-title {
    display: block;
    font-size: 16px;
    margin: 10px 10px;
  }
}
</style>