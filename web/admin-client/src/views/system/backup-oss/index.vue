<template>
  <div class="backup-oss-page">
    <n-card title="备用 OSS 上传失败记录" class="card">
      <template #header-extra>
        <n-space>
          <n-button size="small" @click="retryAll" :loading="retryingAll" type="warning">
            <template #icon><n-icon><RefreshOutline /></n-icon></template>
            重试全部
          </n-button>
          <n-button size="small" @click="fetchList" :loading="loading">
            <template #icon><n-icon><RefreshOutline /></n-icon></template>
            刷新
          </n-button>
        </n-space>
      </template>

      <n-empty v-if="!loading && records.length === 0" description="暂无失败记录" />

      <n-table v-else :bordered="true" :single-line="false" size="small">
        <thead>
          <tr>
            <th>ID</th>
            <th>来源</th>
            <th>OSS 路径</th>
            <th>本地文件</th>
            <th>错误信息</th>
            <th>重试次数</th>
            <th>失败时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in records" :key="r.ID">
            <td>{{ r.ID }}</td>
            <td>
              <n-tag :type="moduleTagType(r.Module)" size="small">
                {{ moduleLabel(r.Module) }}
              </n-tag>
            </td>
            <td class="path-cell">{{ r.ObjectKey }}</td>
            <td class="path-cell" :title="r.FilePath">{{ shortenPath(r.FilePath) }}</td>
            <td class="err-cell">
              <n-ellipsis :line-clamp="1">
                {{ r.ErrMsg }}
              </n-ellipsis>
            </td>
            <td>{{ r.RetryCount }}</td>
            <td>{{ formatTime(r.CreatedAt) }}</td>
            <td>
              <n-button size="tiny" @click="retryOne(r.ID)" :loading="retryingId === r.ID">
                重试
              </n-button>
            </td>
          </tr>
        </tbody>
      </n-table>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import {
  NCard, NButton, NIcon, NSpace, NTable, NTag, NEmpty, NEllipsis,
  useMessage, useDialog
} from 'naive-ui';
import { RefreshOutline } from '@vicons/ionicons5';
import request from '@/utils/request';

interface BackupFailureRecord {
  ID: number;
  ObjectKey: string;
  FilePath: string;
  Module: string;
  ErrMsg: string;
  RetryCount: number;
  CreatedAt: string;
}

const loading = ref(false);
const records = ref<BackupFailureRecord[]>([]);
const retryingId = ref<number | null>(null);
const retryingAll = ref(false);

const message = useMessage();
const dialog = useDialog();

const moduleLabel = (m: string) => {
  const map: Record<string, string> = { image: '图片', cover: '封面', subtitle: '字幕', video: '视频' };
  return map[m] || m;
};

const moduleTagType = (m: string) => {
  const map: Record<string, string> = { image: 'info', cover: 'warning', subtitle: 'success', video: 'error' };
  return map[m] || 'default';
};

const shortenPath = (p: string) => {
  if (p.length <= 50) return p;
  return '...' + p.slice(-47);
};

const formatTime = (t: string) => {
  if (!t) return '-';
  try {
    return t.replace('T', ' ').slice(0, 19);
  } catch {
    return t;
  }
};

const fetchList = async () => {
  loading.value = true;
  try {
    const res = await request.get('v1/backup/failures');
    if (res.data.code === 200) {
      records.value = res.data.data || [];
    }
  } catch {
    records.value = [];
  } finally {
    loading.value = false;
  }
};

const retryOne = async (id: number) => {
  retryingId.value = id;
  try {
    const res = await request.post(`v1/backup/retry/${id}`);
    if (res.data.code === 200) {
      message.success('重试成功');
      await fetchList();
    } else {
      message.warning(res.data.msg || '重试失败');
      await fetchList();
    }
  } catch (e: any) {
    message.error(e.message || '请求失败');
  } finally {
    retryingId.value = null;
  }
};

const retryAll = async () => {
  dialog.warning({
    title: '确认重试全部',
    content: '将重试所有失败记录，确定执行？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      retryingAll.value = true;
      try {
        const res = await request.post('v1/backup/retryAll');
        if (res.data.code === 200) {
          const d = res.data.data || {};
          message.success(`重试完成：成功 ${d.success}，失败 ${d.failed}`);
          await fetchList();
        } else {
          message.warning(res.data.msg || '重试失败');
        }
      } catch (e: any) {
        message.error(e.message || '请求失败');
      } finally {
        retryingAll.value = false;
      }
    },
  });
};

onMounted(() => {
  fetchList();
});
</script>

<style scoped lang="scss">
.backup-oss-page {
  height: 100%;

  .card {
    height: 100%;
  }
}

.path-cell {
  max-width: 280px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.err-cell {
  max-width: 200px;
}
</style>
