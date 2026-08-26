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
import { h, defineComponent, type PropType, onBeforeMount, onBeforeUnmount, reactive, ref, computed } from 'vue';
import { Refresh } from "@vicons/ionicons5";
import useLoading from '@/hooks/loading-hooks';
import { statusCode } from '@/utils/status-code';
import { getVideoListAPI, getFailedVideoListAPI, getProcessingVideoListAPI, deleteVideoAPI, reTranscodeVideoAPI, reTranscodeResourceAPI, reUploadVideoAPI } from '@/api/video';
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

// 重新转码单个分P（仅重试失败的资源，不影响其他分P）
const reTranscodeResource = async (resourceID: number) => {
  const res = await reTranscodeResourceAPI(resourceID);
  if (res.data.code === statusCode.OK) {
    message.success('单分P转码任务已提交');
    await refreshAfterReTranscode();
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

// 分P转码进度组（可展开查看各画质明细）
const TranscodingResourceGroup = defineComponent({
  props: {
    title: String,
    items: Array as PropType<TranscodingProgressItem[]>,
  },
  setup(props) {
    const expanded = ref(false)
    return () => {
      const items = props.items || []
      const n = items.length
      if (n === 0) return null
      const totalPct = items.reduce((s, i) => s + (i.progress || 0), 0)
      const avgPct = Math.round(totalPct / n)
      const anyFail = items.some(i => i.status === 'fail')
      const allSuccess = items.every(i => i.status === 'success')
      const allWaiting = items.every(i => i.status === 'waiting')
      const done = !allWaiting && !anyFail && allSuccess
      const upload = items.find(i => i.upload)

      return h('div', {
        style: 'margin-bottom: 8px; border: 1px solid var(--n-border-color); border-radius: 6px; overflow: hidden;'
      }, [
        // 头部行 — 点击展开/折叠
        h('div', {
          style: 'display: flex; align-items: center; gap: 10px; padding: 8px 10px; cursor: pointer; user-select: none; background: var(--n-action-color);',
          onClick: () => { expanded.value = !expanded.value }
        }, [
          h('span', {
            style: 'font-size: 10px; color: var(--n-text-color-3); width: 14px; text-align: center; transition: transform .2s; flex-shrink: 0;' + (expanded.value ? ' transform: rotate(90deg);' : '')
          }, '▶'),
          h('span', { style: 'font-size: 12px; font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; flex-shrink: 1; min-width: 0;' }, props.title || ''),
          h('span', { style: 'font-size: 11px; color: var(--n-text-color-3); white-space: nowrap; flex-shrink: 0; margin-right: auto;' }, `${n} 个画质`),
          h('div', { style: 'width: 52px; flex-shrink: 0; text-align: center;' }, [
            allWaiting ? h(NTag, { size: 'tiny', type: 'default' }, { default: () => '排队中' }) : null,
            anyFail    ? h(NTag, { size: 'tiny', type: 'error' }, { default: () => '失败' }) : null,
            done       ? h(NTag, { size: 'tiny', type: 'success' }, { default: () => '完成' }) : null,
          ]),
          h('div', { style: 'min-width: 160px; max-width: 320px; flex-shrink: 0;' }, [
            h(NProgress, {
              class: 'transcoding-progress',
              percentage: allWaiting ? 0 : avgPct,
              processing: !anyFail && !allSuccess && !allWaiting,
              status: anyFail ? 'error' : (allSuccess ? 'success' : 'default'),
              height: 18,
              showIndicator: true,
              indicatorPlacement: 'inside',
            })
          ]),
        ]),
        // 折叠面板 — 各画质明细 + 上传进度
        expanded.value ? h('div', { style: 'padding: 8px 12px 4px 26px; border-top: 1px solid var(--n-border-color);' }, [
          ...items.map(item => {
            const waiting = item.status === 'waiting'
            return h('div', { style: 'margin-bottom: 8px;' }, [
              h('div', { style: 'margin-bottom: 2px; font-size: 12px; display: flex; align-items: center; gap: 6px;' }, [
                h('span', { style: 'color: var(--n-text-color-3);' }, item.quality),
                waiting ? h(NTag, { size: 'tiny', type: 'default' }, { default: () => '排队中' }) : null,
              ]),
              h(NProgress, {
                class: 'transcoding-progress',
                percentage: waiting ? 0 : Math.round(item.progress || 0),
                processing: item.status === 'processing',
                status: item.status === 'fail' ? 'error' : (item.status === 'success' ? 'success' : 'default'),
                showIndicator: true,
                indicatorPlacement: 'inside',
                height: 14,
              }),
              item.status === 'fail' ? h('div', { style: 'margin-top: 4px;' }, [
                h(NButton, { size: 'tiny', type: 'warning', onClick: () => reTranscodeResource(item.resourceId) }, { default: () => '单P重试' }),
              ]) : null,
            ])
          }),
          // 上传进度
          upload && upload.status && upload.status !== 'local'
            ? h('div', { style: 'margin-top: 4px; padding-top: 8px; border-top: 1px dashed var(--n-border-color);' }, [
                h('div', { style: 'margin-bottom: 2px; font-size: 12px; display: flex; align-items: center; gap: 6px;' }, [
                  h('span', { style: 'color: var(--n-text-color-3);' }, 'OSS 上传'),
                ]),
                h(NProgress, {
                  percentage: Math.round(upload.progress || 0),
                  processing: upload.status === 'uploading',
                  status: upload.status === 'fail' ? 'error' : (upload.status === 'success' ? 'success' : 'default'),
                  showIndicator: true,
                  indicatorPlacement: 'inside',
                  height: 14,
                }),
                upload.status === 'fail' ? h('div', { style: 'margin-top: 4px;' }, [
                  h(NButton, { size: 'tiny', type: 'warning', onClick: () => reTranscodeResource(items[0].resourceId) }, { default: () => '重新上传' }),
                ]) : null,
              ])
            : null,
        ]) : null
      ])
    }
  }
})

// 处理中列
const processingColumns: DataTableColumns<VideoType> = [
  {
    type: 'expand',
    renderExpand: row => {
      const transcodingParts: any[] = [];
      const details = row.transcodingDetails || [];
      if (details.length > 0) {
        // 按 resourceId 分组
        const groups = new Map<number, { title: string; items: TranscodingProgressItem[] }>()
        for (const item of details) {
          const g = groups.get(item.resourceId)
          if (g) {
            g.items.push(item)
          } else {
            groups.set(item.resourceId, { title: item.resourceTitle || `资源#${item.resourceId}`, items: [item] })
          }
        }
        transcodingParts.push(
          h('div', { style: 'font-size: 13px; font-weight: 600; margin-bottom: 8px;' }, '转码进度')
        );
        for (const [, group] of groups) {
          transcodingParts.push(h(TranscodingResourceGroup, { title: group.title, items: group.items }));
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
    type: 'expand',
    renderExpand: row => {
      const parts: any[] = [];
      const details = row.transcodingDetails || [];
      if (details.length > 0) {
        parts.push(
          h('div', { style: 'font-size: 13px; font-weight: 600; margin-bottom: 8px;' }, '分P状态')
        );
        parts.push(h('div', { style: 'display: flex; flex-direction: column; gap: 8px;' }, details.map(item => {
          const isFail = item.status === 'process_failed' || item.status === 'upload_failed';
          const isSuccess = item.status === 'approved';
          const statusLabel = isFail ? '失败' : (isSuccess ? '成功' : item.status);
          return h('div', { style: 'display: flex; align-items: center; gap: 10px; padding: 4px 0;' }, [
            h('span', { style: 'font-size: 12px; min-width: 120px;' }, item.resourceTitle || `资源#${item.resourceId}`),
            isFail
              ? h(NTag, { size: 'tiny', type: 'error' }, { default: () => statusLabel })
              : isSuccess
                ? h(NTag, { size: 'tiny', type: 'success' }, { default: () => statusLabel })
                : h(NTag, { size: 'tiny', type: 'default' }, { default: () => statusLabel }),
            isFail ? h(NButton, {
              size: 'tiny',
              type: 'warning',
              onClick: () => reTranscodeResource(item.resourceId)
            }, { default: () => '单P重试' }) : null
          ]);
        })));
      }
      if (parts.length === 0) {
        return h('div', { style: 'padding: 8px 4px; color: var(--n-text-color-3);' }, '暂无分P信息');
      }
      return h('div', { style: 'padding: 6px 0;' }, parts);
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

