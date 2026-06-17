<template>
  <div class="pgc-manage">
    <n-card class="pgc-card" :bordered="false">
      <div class="pgc-card-content">
        <n-tabs type="line" v-model:value="activeTab" @update:value="handleTabChange">
          <n-tab name="all">全部</n-tab>
          <n-tab name="approved">已通过</n-tab>
          <n-tab name="rejected">已驳回</n-tab>
          <n-tab name="offline">已下架</n-tab>
        </n-tabs>
        <n-space class="search-bar" justify="space-between">
          <n-space align="center" :size="18">
            <n-button :disabled="loading" size="small" type="primary" @click="getTableData">
              <n-icon>
                <refresh></refresh>
              </n-icon>
            </n-button>
            <n-input v-model:value="keyword" placeholder="搜索标题/简介" clearable size="small"
              style="width: 200px" @keyup.enter="handleSearch" @clear="handleSearch" />
            <n-select v-model:value="pgcTypeFilter" placeholder="类型" clearable size="small"
              style="width: 120px" :options="pgcTypeOptions" @update:value="handleSearch" />
          </n-space>
        </n-space>
        <n-data-table class="table" remote :columns="columns" :data="tableData" :loading="loading"
          :row-key="rowKey"
          :max-height="tableMaxHeight"
          scroll-x="1060"
          :pagination="pagination" />
        <detail-drawer v-model:visible="visibleDrawer" :data="detailsData!" />
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onBeforeMount, reactive, ref } from 'vue';
import useSystemStore from '@/stores/modules/system-store';
import { Refresh } from "@vicons/ionicons5";
import useLoading from '@/hooks/loading-hooks';
import { statusCode } from '@/utils/status-code';
import { getPGCManageListAPI, adminUpdatePGCStatusAPI, adminDeletePGCAPI } from '@/api/pgc';
import type { DataTableColumns, SelectOption } from 'naive-ui';
import { getResourceUrl } from '@/utils/resource';
import DetailDrawer from './components/detail-drawer.vue';
import { NCard, NImage, NIcon, NButton, NDataTable, NSpace, NTabs, NTab, NInput, NSelect, NTag, NPopconfirm, useMessage } from 'naive-ui';

const { loading, startLoading, endLoading } = useLoading(false);
const message = useMessage();
const systemStore = useSystemStore();

const tableMaxHeight = computed(() => {
  const h = systemStore.contentBorderBoxHeight;
  if (typeof h === 'number' && h > 160) {
    return Math.max(h - 160, 240);
  }
  return Math.min(window.innerHeight - 230, 720);
});

const rowKey = (row: PGCReviewRow) =>
  row.pgc_id != null && row.pgc_id !== '' ? String(row.pgc_id) : `id:${row.id ?? ''}`;

const activeTab = ref('all');
const keyword = ref('');
const pgcTypeFilter = ref<number | null>(null);

const pgcTypeOptions: SelectOption[] = [
  { label: '国创', value: 1 },
  { label: '日漫', value: 2 },
  { label: '纪录片', value: 3 },
  { label: '电影', value: 4 },
  { label: '电视剧', value: 5 },
];

const pgcTypeLabel = (t: number) => {
  const m: Record<number, string> = { 1: '国创', 2: '日漫', 3: '纪录片', 4: '电影', 5: '电视剧' };
  return m[t] ?? `类型${t}`;
};

const statusLabel = (s: number) => {
  if (s === -1) return '已下架';
  const m: Record<number, string> = { 0: '草稿', 100: '已提交', 200: '审核中', 300: '已通过', 400: '已驳回' };
  return m[s] ?? String(s);
};

const statusTagType = (s: number): 'default' | 'success' | 'warning' | 'error' | 'info' => {
  if (s === 300) return 'success';
  if (s === 400) return 'error';
  if (s === -1) return 'warning';
  if (s === 100 || s === 200) return 'info';
  return 'default';
};

const tabStatusMap: Record<string, number | undefined> = {
  all: undefined,
  approved: 300,
  rejected: 400,
  offline: -1,
};

// 详情抽屉
const visibleDrawer = ref(false);
const detailsData = ref<PGCReviewRow>();
const viewDetails = (row: PGCReviewRow) => {
  detailsData.value = row;
  visibleDrawer.value = true;
};

// 上架/下架
const toggleStatus = async (row: PGCReviewRow) => {
  const newStatus = row.status === 300 ? -1 : 300;
  const res = await adminUpdatePGCStatusAPI({ pgc_id: row.pgc_id, status: newStatus });
  const code = Number(res.data?.code);
  if (code === statusCode.OK) {
    message.success(newStatus === -1 ? '已下架' : '已上架');
    await getTableData();
  } else {
    message.error((res.data as { msg?: string })?.msg || '操作失败');
  }
};

// 删除
const deletePGC = async (row: PGCReviewRow) => {
  const res = await adminDeletePGCAPI(row.pgc_id);
  const code = Number(res.data?.code);
  if (code === statusCode.OK) {
    message.success('删除成功');
    await getTableData();
  } else {
    message.error((res.data as { msg?: string })?.msg || '删除失败');
  }
};

const columns: DataTableColumns<PGCReviewRow> = [
  {
    key: 'pgc_id',
    title: 'PGC ID',
    width: 160,
    align: 'center'
  },
  {
    key: 'cover',
    title: '封面',
    align: 'center',
    width: 80,
    render: row => {
      if (!row.cover) return '-';
      return h(NImage, { src: getResourceUrl(row.cover), width: 60, height: 32 });
    }
  },
  {
    key: 'title',
    title: '标题',
    align: 'center',
    ellipsis: { tooltip: true }
  },
  {
    key: 'pgc_type',
    title: '类型',
    width: 90,
    align: 'center',
    render: row => pgcTypeLabel(row.pgc_type),
  },
  {
    key: 'status',
    title: '状态',
    width: 90,
    align: 'center',
    render: row => h(NTag, { size: 'small', type: statusTagType(row.status) }, { default: () => statusLabel(row.status) })
  },
  {
    key: 'created_at',
    title: '创建时间',
    width: 170,
    align: 'center',
    render: row => row.created_at ? new Date(row.created_at).toLocaleString() : '-'
  },
  {
    key: 'actions',
    title: '操作',
    align: 'center',
    width: 220,
    render: row => {
      const btns = [
        h(NButton, { size: 'small', type: 'primary', onClick: () => viewDetails(row) }, { default: () => '详情' }),
      ];
      // 已通过 -> 可下架；已下架 -> 可上架
      if (row.status === 300) {
        btns.push(h(NPopconfirm, { onPositiveClick: () => toggleStatus(row) }, {
          default: () => '确定下架？',
          trigger: () => h(NButton, { size: 'small', type: 'warning' }, { default: () => '下架' })
        }));
      } else if (row.status === -1) {
        btns.push(h(NButton, { size: 'small', type: 'success', onClick: () => toggleStatus(row) }, { default: () => '上架' }));
      }
      btns.push(h(NPopconfirm, { onPositiveClick: () => deletePGC(row) }, {
        default: () => '确定删除？此操作不可恢复',
        trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' })
      }));
      return h(NSpace, { justify: 'center' }, { default: () => btns });
    }
  }
];

const tableData = ref<PGCReviewRow[]>([]);
const getTableData = async () => {
  startLoading();
  try {
    const params: PGCManageListParam = {
      page: pagination.page || 1,
      pageSize: pagination.pageSize || 10,
    };
    const tabStatus = tabStatusMap[activeTab.value];
    if (tabStatus !== undefined) {
      params.status = tabStatus;
    }
    if (keyword.value) {
      params.keyword = keyword.value;
    }
    if (pgcTypeFilter.value) {
      params.pgc_type = pgcTypeFilter.value;
    }

    const res = await getPGCManageListAPI(params);
    const code = Number(res.data?.code);
    const payload = res.data?.data as { list?: PGCReviewRow[]; total?: number } | undefined;
    if (code === statusCode.OK && payload) {
      tableData.value = Array.isArray(payload.list) ? payload.list : [];
      pagination.itemCount = payload.total ?? 0;
      return;
    }
    tableData.value = [];
    pagination.itemCount = 0;
    message.error((res.data as { msg?: string })?.msg || '获取 PGC 列表失败');
  } catch {
    tableData.value = [];
    pagination.itemCount = 0;
    message.error('获取 PGC 列表失败');
  } finally {
    endLoading();
  }
};

const handleSearch = () => {
  pagination.page = 1;
  getTableData();
};

const handleTabChange = () => {
  pagination.page = 1;
  getTableData();
};

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
  await getTableData();
});
</script>

<style lang="scss" scoped>
.pgc-manage {
  height: 100%;

  .pgc-card {
    border-radius: 0;

    :deep(.n-card__content) {
      height: calc(100vh - 130px);
      padding-right: 6px;

      .pgc-card-content {
        height: 100%;
        display: flex;
        flex-direction: column;

        .search-bar {
          margin: 12px 0;
        }

        .table {
          flex: 1;
          min-height: 0;
        }
      }
    }
  }
}
</style>
