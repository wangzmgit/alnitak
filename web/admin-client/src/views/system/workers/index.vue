<template>
  <div class="workers-page">
    <n-card title="远程转码 Worker 状态" class="card">
      <template #header-extra>
        <n-button size="small" @click="refresh" :loading="loading">
          <template #icon><n-icon><RefreshOutline /></n-icon></template>
          刷新
        </n-button>
      </template>

      <n-empty v-if="!loading && workers.length === 0" description="暂无在线 Worker" />

      <n-grid v-else :cols="1" :x-gap="12" :y-gap="12">
        <n-grid-item v-for="w in workers" :key="w.groupID">
          <n-card size="small" :segmented="true" class="worker-card">
            <n-space align="center" justify="space-between">
              <n-space align="center">
                <n-icon size="24" :color="w.healthy ? '#18a058' : '#d03050'">
                  <CheckmarkCircle v-if="w.healthy" />
                  <CloseCircle v-else />
                </n-icon>
                <n-text strong>{{ w.groupID }}</n-text>
              </n-space>
              <n-tag :type="w.healthy ? 'success' : 'error'" size="small">
                {{ w.healthy ? '在线' : '离线' }}
              </n-tag>
            </n-space>

            <n-descriptions :column="3" bordered size="small" class="worker-desc">
              <n-descriptions-item label="运行时长">
                {{ w.uptime }}
              </n-descriptions-item>
              <n-descriptions-item label="并发数">
                {{ w.concurrency }}
              </n-descriptions-item>
              <n-descriptions-item label="活跃任务">
                {{ w.jobsActive }}
              </n-descriptions-item>
              <n-descriptions-item label="累计任务">
                {{ w.jobsTotal }}
              </n-descriptions-item>
              <n-descriptions-item label="失败任务">
                {{ w.jobsFailed }}
              </n-descriptions-item>
              <n-descriptions-item label="最后心跳">
                {{ formatTime(w.lastSeen) }}
              </n-descriptions-item>
            </n-descriptions>
          </n-card>
        </n-grid-item>
      </n-grid>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getWorkersAPI } from '@/api/workers';
import {
  NCard, NButton, NIcon, NGrid, NGridItem,
  NSpace, NText, NTag, NDescriptions, NDescriptionsItem, NEmpty
} from 'naive-ui';
import {
  RefreshOutline,
  CheckmarkCircle,
  CloseCircle
} from '@vicons/ionicons5';
import moment from 'moment';

interface WorkerHeartbeat {
  healthy: boolean;
  startedAt: string;
  uptime: string;
  concurrency: number;
  jobsActive: number;
  jobsTotal: number;
  jobsFailed: number;
  groupID: string;
  lastSeen: number;
}

const loading = ref(false);
const workers = ref<WorkerHeartbeat[]>([]);

const fetchWorkers = async () => {
  loading.value = true;
  try {
    const res = await getWorkersAPI();
    if (res.data.code === 200) {
      workers.value = res.data.data.workers || [];
    }
  } catch {
    workers.value = [];
  } finally {
    loading.value = false;
  }
};

const refresh = () => {
  fetchWorkers();
};

const formatTime = (ts: number) => {
  if (!ts) return '-';
  return moment.unix(ts).format('YYYY-MM-DD HH:mm:ss');
};

onMounted(() => {
  fetchWorkers();
});
</script>

<style scoped lang="scss">
.workers-page {
  height: 100%;

  .card {
    height: 100%;
  }
}

.worker-card {
  margin-bottom: 4px;
}

.worker-desc {
  margin-top: 12px;
}
</style>
