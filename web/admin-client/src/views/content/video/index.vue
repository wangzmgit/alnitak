<template>
  <div class="user-manage">
    <n-card class="user-card" :bordered="false">
      <div class="user-card-content">
        <n-tabs type="line" v-model:value="activeTab" @update:value="handleTabChange">
          <n-tab name="published">已发布</n-tab>
          <n-tab name="processing">
            处理中
            <n-badge v-if="processingCount > 0" :value="processingCount" :max="99"
              style="margin-left: 6px;" />
          </n-tab>
          <n-tab name="failed">
            处理失败
            <n-badge v-if="failedCount > 0" :value="failedCount" :max="99"
              style="margin-left: 6px;" />
          </n-tab>
        </n-tabs>
        <n-space class="search-bar" justify="space-between">
          <n-space align="center" :size="18">
            <n-button :disabled="loading" size="small" type="primary" @click="handleRefreshClick">
              <n-icon>
                <refresh></refresh>
              </n-icon>
            </n-button>
            <n-button v-if="activeTab === 'failed'" :disabled="loading || tableData.length === 0"
              size="small" type="warning" @click="reTranscodeAll">
              全部重新转码
            </n-button>
          </n-space>
        </n-space>
        <n-data-table class="table" remote :columns="currentColumns" :data="tableData" :loading="loading"
          :pagination="pagination" :row-key="(row: VideoType) => row.vid" flex-height />
        <table-action-drawer v-model:visible="visibleDrawer" :data="editData!"></table-action-drawer>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { h, onBeforeMount, onBeforeUnmount, reactive, ref, computed } from 'vue';
import { Refresh } from "@vicons/ionicons5";
import useLoading from '@/hooks/loading-hooks';
import { statusCode } from '@/utils/status-code';
import { getVideoListAPI, getFailedVideoListAPI, getProcessingVideoListAPI, deleteVideoAPI, reTranscodeVideoAPI, reUploadVideoAPI } from '@/api/video';
import type { DataTableColumns } from 'naive-ui';
import { getResourceUrl } from '@/utils/resource';
import usePartition from '@/hooks/partition-hooks';
import TableActionDrawer from './components/table-action-drawer.vue';
import { NCard, NImage, NIcon, NButton, NDataTable, NPopconfirm, NSpace, NTabs, NTab, NBadge, NProgress, NTag, useMessage, useDialog } from 'naive-ui';

const { loading, startLoading, endLoading } = useLoading(false);
const { getPartition, getPartitionName } = usePartition("video");

const message = useMessage();
const dialog = useDialog();

const activeTab = ref('published');
const failedCount = ref(0);
const processingCount = ref(0);
let listRefreshTimer: number | null = null;

const visibleDrawer = ref(false);
const openDrawer = () => {
  visibleDrawer.value = true;
}

const refreshAfterReTranscode = async () => {
  await getTableData();
  await fetchFailedCount();
  await fetchProcessingCount();
}

// 编辑视频
const editData = ref<VideoType>();
const editVideo = (row: VideoType) => {
  editData.value = row;
  openDrawer();
}

// 删除视频
const deleteVideo = async (row: VideoType) => {
  const res = await deleteVideoAPI(row.vid);
  if (res.data.code === statusCode.OK) {
    message.success('删除成功');
    await getTableData();
    if (activeTab.value !== 'failed') await fetchFailedCount();
    if (activeTab.value !== 'processing') await fetchProcessingCount();
  } else {
    message.error(res.data.msg);
  }
}

// 重新转码视频（后端按 vid 维度处理所有分P，只需调用一次）
const reTranscodeVideo = async (row: VideoType) => {
  const res = await reTranscodeVideoAPI(row.vid);
  if (res.data.code === statusCode.OK) {
    message.success('重新转码任务已提交');
    await refreshAfterReTranscode();
  } else {
    message.error(res.data.msg);
  }
}

// 重新上传OSS（转码成功但上传失败时重置上传）
const reUploadVideo = async (row: VideoType) => {
  const res = await reUploadVideoAPI(row.vid);
  if (res.data.code === statusCode.OK) {
    message.success('重新上传任务已提交');
    await refreshAfterReTranscode();
  } else {
    message.error(res.data.msg);
  }
}

const fetchAllFailedVideos = async () => {
  const pageSize = 100;
  let page = 1;
  let total = 0;
  const result: VideoType[] = [];

  while (true) {
    const res = await getFailedVideoListAPI({ page, pageSize });
    if (res.data.code !== statusCode.OK) {
      throw new Error(res.data.msg || '获取失败视频列表失败');
    }

    const list: VideoType[] = res.data.data.list || [];
    total = res.data.data.total || 0;
    result.push(...list);

    if (result.length >= total || list.length === 0) {
      break;
    }
    page += 1;
  }

  return total > 0 ? result.slice(0, total) : result;
}

// 全部重新转码
const reTranscodeAll = () => {
  dialog.warning({
    title: '确认',
    content: `确定要对全部 ${pagination.itemCount} 个失败视频重新转码吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      startLoading();
      try {
        const allFailedRows = await fetchAllFailedVideos();
        if (allFailedRows.length === 0) {
          message.info('当前没有处理失败视频');
          return;
        }

        let successCount = 0;
        let failCount = 0;
        for (const row of allFailedRows) {
          try {
            const res = await reTranscodeVideoAPI(row.vid);
            if (res.data.code === statusCode.OK) successCount++;
            else failCount++;
          } catch {
            failCount++;
          }
        }
        message.success(`已提交 ${successCount} 个转码任务${failCount > 0 ? `，${failCount} 个失败` : ''}`);
        await refreshAfterReTranscode();
      } catch (e: any) {
        message.error(e?.message || '批量提交转码任务失败');
      } finally {
        endLoading();
      }
    }
  });
}

// 已发布列
const publishedColumns: DataTableColumns<VideoType> = [
  {
    key: 'vid',
    title: 'ID',
    width: 90,
    align: 'center'
  },
  {
    key: 'avatar',
    title: '封面',
    align: 'center',
    width: 80,
    render: row => {
      return h(NImage, {
        src: getResourceUrl(row.cover),
        width: 60,
        height: 32,
      })
    }
  },
  {
    key: 'title',
    title: '标题',
    align: 'center'
  },
  {
    key: 'desc',
    title: '简介',
    align: 'center',
  },
  {
    key: 'partition',
    title: '分区',
    align: 'center',
    render: row => {
      return getPartitionName(row.partitionId)
    }
  },
  {
    key: 'clicks',
    title: '播放量',
    align: 'center'
  },
  {
    key: 'actions',
    title: '操作',
    align: 'center',
    width: 240,
    render: row => {
      return h(NSpace, { justify: 'center' }, {
        default: () => [
          h(NButton, {
            size: 'small',
            onClick: () => editVideo(row)
          }, { default: () => '详情' }),
          h(NButton, {
            size: 'small',
            type: 'warning',
            onClick: () => reTranscodeVideo(row)
          }, { default: () => '转码' }),
          h(NPopconfirm, {
            onPositiveClick: () => deleteVideo(row),
          }, {
            default: () => '是否删除视频?',
            trigger: () => h(NButton, {
              size: 'small',
              type: 'error',
            }, { default: () => '删除' })
          })
        ]
      })
    }
  }
]

const ossTypeLabels: Record<string, string> = {
  aliyun: '阿里云OSS',
  minio: 'MinIO',
  cloudflare: 'Cloudflare R2',
  local: '本地存储',
}

const uploadProgressLabel = (info: UploadProgressInfo | undefined): string => {
  if (!info) return '';
  const ossName = ossTypeLabels[info.ossType] || info.ossType;
  if (info.status === 'local') return '本地存储，无需上传';
  if (info.status === 'uploading') return `上传至 ${ossName}`;
  if (info.status === 'success') return 'OSS 上传完成';
  if (info.status === 'fail') return 'OSS 上传失败';
  return '';
}

// 处理中列
const processingColumns: DataTableColumns<VideoType> = [
  {
    type: 'expand',
    renderExpand: row => {
      const transcodingParts: any[] = [];
      const details = row.transcodingDetails || [];
      if (details.length > 0) {
        transcodingParts.push(
          h('div', { style: 'font-size: 13px; font-weight: 600; margin-bottom: 8px;' }, '转码进度')
        );
        transcodingParts.push(...details.map(item => {
          const isWaiting = item.status === 'waiting';
          return h('div', { style: 'margin-bottom: 10px;' }, [
            h('div', { style: 'margin-bottom: 4px; font-size: 12px; display: flex; align-items: center; gap: 8px;' }, [
              `${item.resourceTitle || `资源#${item.resourceId}`} / ${item.quality}`,
              isWaiting ? h(NTag, { size: 'tiny', type: 'default' }, { default: () => '排队中' }) : null,
            ]),
            h(NProgress, {
              percentage: isWaiting ? 0 : Math.round(item.progress || 0),
              processing: item.status === 'processing',
              status: item.status === 'fail' ? 'error' : (item.status === 'success' ? 'success' : 'default'),
              showIndicator: true,
              indicatorPlacement: 'inside',
            })
          ]);
        }));
      }

      // 上传进度
      const up = row.uploadProgress;
      if (up && up.status && up.status !== 'local') {
        transcodingParts.push(
          h('div', { style: 'font-size: 13px; font-weight: 600; margin: 12px 0 8px 0;' }, 'OSS 上传进度')
        );
        transcodingParts.push(
          h('div', { style: 'margin-bottom: 4px; font-size: 12px;' }, uploadProgressLabel(up))
        );
        transcodingParts.push(
          h(NProgress, {
            percentage: Math.round(up.progress || 0),
            processing: up.status === 'uploading',
            status: up.status === 'fail' ? 'error' : (up.status === 'success' ? 'success' : 'default'),
            showIndicator: true,
            indicatorPlacement: 'inside',
          })
        );
        if (up.status === 'fail') {
          transcodingParts.push(
            h('div', { style: 'margin-top: 8px;' }, [
              h(NButton, {
                size: 'small',
                type: 'warning',
                onClick: () => reUploadVideo(row)
              }, { default: () => '重新上传' })
            ])
          );
        }
      }

      if (transcodingParts.length === 0) {
        return h('div', { style: 'padding: 8px 4px; color: var(--n-text-color-3);' }, '暂无进度明细');
      }
      return h('div', { style: 'padding: 6px 0;' }, transcodingParts);
    }
  },
  {
    key: 'vid',
    title: 'ID',
    width: 70,
    align: 'center'
  },
  {
    key: 'avatar',
    title: '封面',
    align: 'center',
    width: 80,
    render: row => {
      return h(NImage, {
        src: getResourceUrl(row.cover),
        width: 60,
        height: 32,
      })
    }
  },
  {
    key: 'title',
    title: '标题',
    align: 'center',
    ellipsis: { tooltip: true }
  },
  {
    key: 'author',
    title: '上传者',
    align: 'center',
    width: 120,
    render: row => row.author?.name || '-'
  },
  {
    key: 'progress',
    title: '总体转码进度',
    align: 'center',
    width: 200,
    render: row => h(NProgress, {
      percentage: Math.round(row.transcodingProgress || 0),
      processing: true,
      showIndicator: true
    })
  },
  {
    key: 'uploadProgress',
    title: '上传状态',
    align: 'center',
    width: 160,
    render: row => {
      const up = row.uploadProgress;
      if (!up || !up.status) return h('span', { style: 'color: var(--n-text-color-3); font-size: 12px;' }, '等待中');
      if (up.status === 'local') return h('span', { style: 'color: var(--n-text-color-3);' }, '本地存储');
      if (up.status === 'success') return h(NTag, { size: 'tiny', type: 'success' }, { default: () => '上传完成' });
      if (up.status === 'fail') return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NTag, { size: 'tiny', type: 'error' }, { default: () => '上传失败' }),
          h(NButton, { size: 'tiny', type: 'warning', onClick: () => reUploadVideo(row) }, { default: () => '重新上传' })
        ]
      });
      return h('div', { style: 'display: flex; align-items: center; gap: 6px;' }, [
        h(NProgress, {
          percentage: Math.round(up.progress || 0),
          processing: true,
          height: 16,
          showIndicator: true,
          style: { width: '100px' }
        }),
      ]);
    }
  },
  {
    key: 'createdAt',
    title: '创建时间',
    align: 'center',
    width: 170,
    render: row => new Date(row.createdAt).toLocaleString()
  }
]

// 失败列
const failedColumns: DataTableColumns<VideoType> = [
  {
    key: 'vid',
    title: 'ID',
    width: 70,
    align: 'center'
  },
  {
    key: 'avatar',
    title: '封面',
    align: 'center',
    width: 80,
    render: row => {
      return h(NImage, {
        src: getResourceUrl(row.cover),
        width: 60,
        height: 32,
      })
    }
  },
  {
    key: 'title',
    title: '标题',
    align: 'center',
    ellipsis: { tooltip: true }
  },
  {
    key: 'author',
    title: '上传者',
    align: 'center',
    width: 120,
    render: row => {
      return row.author?.name || '-';
    }
  },
  {
    key: 'createdAt',
    title: '创建时间',
    align: 'center',
    width: 170,
    render: row => {
      return new Date(row.createdAt).toLocaleString();
    }
  },
  {
    key: 'actions',
    title: '操作',
    align: 'center',
    width: 180,
    render: row => {
      return h(NSpace, { justify: 'center' }, {
        default: () => [
          h(NButton, {
            size: 'small',
            type: 'warning',
            onClick: () => reTranscodeVideo(row)
          }, { default: () => '重新转码' }),
          h(NPopconfirm, {
            onPositiveClick: () => deleteVideo(row),
          }, {
            default: () => '是否删除视频?',
            trigger: () => h(NButton, {
              size: 'small',
              type: 'error',
            }, { default: () => '删除' })
          })
        ]
      })
    }
  }
]

const currentColumns = computed(() => {
  if (activeTab.value === 'processing') return processingColumns;
  if (activeTab.value === 'failed') return failedColumns;
  return publishedColumns;
});

const tableData = ref<VideoType[]>([]);
const handleRefreshClick = () => getTableData();
const getTableData = async (options?: { silent?: boolean }) => {
  const silent = options?.silent === true;
  if (!silent) {
    startLoading();
  }
  try {
    const page = pagination.page || 1;
    const pageSize = pagination.pageSize || 1;

    const api = activeTab.value === 'published'
      ? getVideoListAPI
      : (activeTab.value === 'processing' ? getProcessingVideoListAPI : getFailedVideoListAPI);
    const res = await api({ page, pageSize });
    if (res.data.code === statusCode.OK) {
      tableData.value = res.data.data.list || [];
      pagination.itemCount = res.data.data.total;
      if (activeTab.value === 'failed') {
        failedCount.value = res.data.data.total;
      }
      if (activeTab.value === 'processing') {
        processingCount.value = res.data.data.total;
      }
      return;
    }
    tableData.value = [];
    pagination.itemCount = 0;
    if (!silent) {
      message.error(res.data.msg || '获取视频列表失败');
    }
  } catch {
    tableData.value = [];
    pagination.itemCount = 0;
    if (!silent) {
      message.error('获取视频列表失败');
    }
  } finally {
    if (!silent) {
      endLoading();
    }
  }
}

// 获取失败数量（用于badge显示）
const fetchFailedCount = async () => {
  try {
    const res = await getFailedVideoListAPI({ page: 1, pageSize: 1 });
    if (res.data.code === statusCode.OK) {
      failedCount.value = res.data.data.total;
    }
  } catch {
    // 忽略角标拉取失败，不影响主流程
  }
}

const fetchProcessingCount = async () => {
  try {
    const res = await getProcessingVideoListAPI({ page: 1, pageSize: 1 });
    if (res.data.code === statusCode.OK) {
      processingCount.value = res.data.data.total;
    }
  } catch {
    // 忽略角标拉取失败，不影响主流程
  }
}

const startListAutoRefresh = () => {
  if (listRefreshTimer) {
    window.clearInterval(listRefreshTimer);
  }
  const intervalMs = 3000;
  listRefreshTimer = window.setInterval(async () => {
    if (activeTab.value !== 'processing') {
      return;
    }
    await getTableData({ silent: true });
    await fetchFailedCount();
    await fetchProcessingCount();
  }, intervalMs);
}

const stopListAutoRefresh = () => {
  if (listRefreshTimer) {
    window.clearInterval(listRefreshTimer);
    listRefreshTimer = null;
  }
}

const handleTabChange = () => {
  pagination.page = 1;
  if (activeTab.value === 'processing') {
    startListAutoRefresh();
  } else {
    stopListAutoRefresh();
  }
  getTableData();
}

const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 15, 20, 25, 30],
  onChange: (page: number) => {
    pagination.page = page;
    getTableData();
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
    getTableData();
  }
});

onBeforeMount(async () => {
  await getPartition();
  await getTableData();
  await fetchFailedCount();
  await fetchProcessingCount();
  if (activeTab.value === 'processing') {
    startListAutoRefresh();
  }
})

onBeforeUnmount(() => {
  stopListAutoRefresh();
})
</script>

<style lang="scss" scoped>
.user-manage {
  height: 100%;

  .user-card {
    height: 100%;

    .user-card-content {
      height: 100%;
      display: flex;
      flex-direction: column;

      .search-bar {
        padding: 12px 0;
      }

      .table {
        flex: 1;
      }
    }
  }
}
</style>
