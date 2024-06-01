<template>
  <div class="system-announce">
    <n-card class="announce-card" :bordered="false">
      <div class="announce-card-content">
        <n-space class="search-bar" justify="space-between">
          <n-space>
            <n-button type="primary" @click="addAnnounce()">添加公告</n-button>
          </n-space>
          <n-space align="center" :size="18">
            <n-button :disabled="loading" size="small" type="primary" @click="getTableData">
              <n-icon>
                <refresh></refresh>
              </n-icon>
            </n-button>
          </n-space>
        </n-space>
        <n-data-table class="table" remote :columns="columns" :data="tableData" :loading="loading"
          :pagination="pagination" flex-height />
        <table-action-modal v-model:visible="visible" :type="modalType" :edit-data="editData" @refresh="getTableData" />
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, reactive, ref, h } from 'vue';
import { Refresh } from "@vicons/ionicons5";
import useLoading from '@/hooks/loading-hooks';
import { statusCode } from '@/utils/status-code';
import type { DataTableColumns } from 'naive-ui';
import TableActionModal from './components/table-action-modal.vue';
import { NCard, NPopconfirm, NIcon, NButton, NDataTable, NSpace, useMessage } from 'naive-ui';
import { formatTime } from '@/utils/format';
import { deleteAnnounceAPI, getAnnounceListAPI } from '@/api/announce';

const { loading, startLoading, endLoading } = useLoading(false);

const message = useMessage();

const modalType = ref<"add" | "edit">('add');
const visible = ref(false);
const openModal = () => {
  visible.value = true;
}

// 添加API
const addAnnounce = () => {
  modalType.value = 'add';
  openModal();
}

const editData = ref<AnnounceType>();
const editAnnounce = (row: AnnounceType) => {
  editData.value = row;
  modalType.value = 'edit';
  openModal();
}

// 删除数据
const deleteAnnounce = async (row: AnnounceType) => {
  const res = await deleteAnnounceAPI(row.id);
  if (res.data.code === statusCode.OK) {
    message.success('删除成功');
    await getTableData();
  } else {
    message.error(res.data.msg);
  }
}

const columns: DataTableColumns<AnnounceType> = [
  {
    key: 'id',
    title: 'ID',
    align: 'center'
  },
  {
    key: 'title',
    title: '标题',
    align: 'center'
  },
  {
    key: 'content',
    title: '内容',
    align: 'center',

  },
  {
    key: 'url',
    title: 'URL',
    align: 'center',
  },
  {
    key: 'createdAt',
    title: '发布时间',
    align: 'center',
    render: row => {
      return row.createdAt ? formatTime(row.createdAt) : ""
    }
  },
  {
    key: 'actions',
    title: '操作',
    align: 'center',
    width: 320,
    render: row => {
      return h(NSpace, { justify: 'center' }, {
        default: () => [
          // h(NButton, {
          //   size: 'small',
          //   onClick: () => editAnnounce(row)
          // }, { default: () => '编辑' }),
          h(NPopconfirm, {
            onPositiveClick: () => deleteAnnounce(row),
          }, {
            default: () => '是否删除公告?',
            trigger: () => h(NButton, {
              size: 'small',
            }, { default: () => '删除' })
          })
        ]
      })
    }
  }
]

const tableData = ref<AnnounceType[]>([]);
async function getTableData() {
  startLoading();
  const page = pagination.page || 1;
  const pageSize = pagination.pageSize || 1;
  const res = await getAnnounceListAPI(page, pageSize);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.announces) {
      tableData.value = res.data.data.announces;
    } else {
      tableData.value = [];
    }

    pagination.itemCount = res.data.data.total;
    endLoading();
  }
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

onBeforeMount(() => {
  getTableData();
})
</script>

<style lang="scss" scoped>
.system-announce {
  height: 100%;

  .announce-card {
    height: 100%;

    .announce-card-content {
      height: 100%;
      display: flex;
      flex-direction: column;

      .search-bar {
        padding-bottom: 12px;
      }

      .table {
        flex: 1;
      }
    }
  }
}
</style>
