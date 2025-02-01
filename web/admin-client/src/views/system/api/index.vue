<template>
  <div class="system-api">
    <n-card class="api-card" :bordered="false">
      <div class="api-card-content">
        <n-space class="search-bar" justify="space-between">
          <n-space>
            <n-button type="primary" @click="addMenu()">添加API</n-button>
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
import { getApiListAPI, deleteApiAPI } from '@/api/api';
import { statusCode } from '@/utils/status-code';
import type { DataTableColumns } from 'naive-ui';
import TableActionModal from './components/table-action-modal.vue';
import { NCard, NPopconfirm, NIcon, NButton, NDataTable, NSpace, useMessage } from 'naive-ui';

const { loading, startLoading, endLoading } = useLoading(false);

const message = useMessage();

const modalType = ref<"add" | "edit">('add');
const visible = ref(false);
const openModal = () => {
  visible.value = true;
}

// 添加API
const addMenu = () => {
  modalType.value = 'add';
  openModal();
}

const editData = ref<ApiItemType>();
const editApi = (row: ApiItemType) => {
  editData.value = row;
  modalType.value = 'edit';
  openModal();
}

// 删除数据
const deleteApi = async (row: ApiItemType) => {
  const res = await deleteApiAPI(row.id);
  if (res.data.code === statusCode.OK) {
    message.success('删除成功');
    await getTableData();
  } else {
    message.error(res.data.msg);
  }
}

const columns: DataTableColumns<ApiItemType> = [
  {
    key: 'id',
    title: '序号',
    align: 'center'
  },
  {
    key: 'method',
    title: '请求方法',
    align: 'center'
  },
  {
    key: 'path',
    title: '请求路径',
    align: 'center',

  },
  {
    key: 'category',
    title: '分类',
    align: 'center',
  },
  {
    key: 'desc',
    title: '简介',
    align: 'center',
  },
  {
    key: 'actions',
    title: '操作',
    align: 'center',
    width: 320,
    render: row => {
      return h(NSpace, { justify: 'center' }, {
        default: () => [
          h(NButton, {
            size: 'small',
            onClick: () => editApi(row)
          }, { default: () => '编辑' }),
          h(NPopconfirm, {
            onPositiveClick: () => deleteApi(row),
          }, {
            default: () => '是否删除API?',
            trigger: () => h(NButton, {
              size: 'small',
            }, { default: () => '删除' })
          })
        ]
      })
    }
  }
]

const tableData = ref<ApiItemType[]>([]);
async function getTableData() {
  startLoading();
  const page = pagination.page || 1;
  const pageSize = pagination.pageSize || 1;
  const res = await getApiListAPI({ page, pageSize });
  if (res.data.code === statusCode.OK) {
    if ( res.data.data.list) {
      tableData.value = res.data.data.list;
      pagination.itemCount = res.data.data.total;
    }
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
.system-api {
  height: 100%;

  .api-card {
    height: 100%;

    .api-card-content {
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
