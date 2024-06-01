<template>
  <div class="announce">
    <p class="announce-title">站内公告</p>
    <div class="announce-box">
      <el-scrollbar max-height="100%">
        <ul class="announce-list">
          <li class="announce-item" v-for="(item, index) in announceList" :key="index">
            <div class="title">
              <p class="item-title">{{ item.title }}</p>
              <span class="item-time"> {{ formatTime(item.createdAt) }}</span>
            </div>
            <span class="announce-content">{{ item.content }}</span>
            <span class="link" v-if="item.url">
              <nuxt-link :to="item.url">网页链接</nuxt-link>
            </span>
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
import { getAnnounceAPI } from '@/api/msg-announce';

const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const showPagination = ref(true);
const announceList = ref<AnnounceType[]>([]);
const getAnnounceList = async () => {
  const res = await getAnnounceAPI(page.value, pageSize.value);
  if (res.data.code === statusCode.OK) {
    announceList.value = res.data.data.announces;
  }
}

onBeforeMount(() => {
  getAnnounceList();
})
</script>

<style lang="scss" scoped>
.announce {
  padding: 0 18px 0;
  height: 100%;
  box-sizing: border-box;
  background-color: #fff;

  .announce-title {
    font-size: 18px;
    margin: 0;
    padding: 16px 0 10px;
  }

  .announce-box {
    width: 100%;
    height: calc(100% - 86px);
  }

  .pagination {
    height: 36px;
    display: flex;
    align-items: center;
  }
}

.announce-list {
  padding: 0;
  margin: 0;

  .announce-item {
    padding: 8px 0 16px;
    border-bottom: 1px solid #e6e6e6;

    .title {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin: 5px 0;

      .item-title {
        font-size: 15px;
        font-weight: 600;
        color: #333;
        width: calc(100% - 160px);
        margin: 0;
      }

      .item-time {
        color: #999;
        font-size: 12px;
      }
    }

    .announce-content {
      color: #666;
      font-size: 14px;
    }

    .link {
      font-size: 12px;
      margin-left: 10px;

      a {
        cursor: pointer;
        color: var(--primary-color);
        text-decoration: none;
      }
    }
  }
}
</style>