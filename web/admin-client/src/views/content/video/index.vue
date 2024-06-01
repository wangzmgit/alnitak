<template>
  <div class="user-manage">
    <n-card class="user-card" :bordered="false">
      <div class="user-card-content">
        <n-space class="search-bar" justify="space-between">
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
        <table-action-drawer v-model:visible="visibleDrawer" :data="editData!"></table-action-drawer>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { h, onBeforeMount, reactive, ref } from 'vue';
import { formatTime } from '@/utils/format';
import { Refresh } from "@vicons/ionicons5";
import useLoading from '@/hooks/loading-hooks';
import { statusCode } from '@/utils/status-code';
import { getVideoListAPI, deleteVideoAPI } from '@/api/video';
import type { DataTableColumns } from 'naive-ui';
import { getResourceUrl } from '@/utils/resource';
import usePartition from '@/hooks/partition-hooks';
import TableActionDrawer from './components/table-action-drawer.vue';
import { NCard, NImage, NIcon, NButton, NDataTable, NPopconfirm, NSpace, useMessage } from 'naive-ui';

const { loading, startLoading, endLoading } = useLoading(false);
const { getPartition, getPartitionName } = usePartition("video");

const message = useMessage();

const visibleDrawer = ref(false);
const openDrawer = () => {
  visibleDrawer.value = true;
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
  } else {
    message.error(res.data.msg);
  }
}

const columns: DataTableColumns<VideoType> = [
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
    width: 160,
    render: row => {
      return h(NSpace, { justify: 'center' }, {
        default: () => [
          h(NButton, {
            size: 'small',
            onClick: () => editVideo(row)
          }, { default: () => '详情' }),
          h(NPopconfirm, {
            onPositiveClick: () => deleteVideo(row),
          }, {
            default: () => '是否删除视频?',
            trigger: () => h(NButton, {
              size: 'small',
            }, { default: () => '删除' })
          })
        ]
      })

    }
  }
]

const tableData = ref<VideoType[]>([]);
const getTableData = async () => {
  startLoading();
  const page = pagination.page || 1;
  const pageSize = pagination.pageSize || 1;
  const res = await getVideoListAPI({ page, pageSize });
  if (res.data.code === statusCode.OK) {
    if (res.data.data.list) {
      tableData.value = res.data.data.list;
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

onBeforeMount(async () => {
  await getPartition();
  await getTableData();
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
        padding-bottom: 12px;
      }

      .table {
        flex: 1;
      }
    }
  }
}
</style>