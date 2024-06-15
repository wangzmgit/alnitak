<template>
  <div class="user-manage">
    <n-card class="user-card" :bordered="false">
      <div class="user-card-content">
        <n-space class="search-bar" justify="space-between">
          <n-space>
            <n-button type="primary" @click="addCarousel()">添加轮播图</n-button>
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
        <table-action-drawer v-model:visible="visibleDrawer" :data="detailsData!" :partitions="partitions"
          @finish="reviewFinish"></table-action-drawer>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { h, onBeforeMount, reactive, ref } from 'vue';
import { Refresh } from "@vicons/ionicons5";
import useLoading from '@/hooks/loading-hooks';
import {constant} from "@/utils/constant";
import { statusCode } from '@/utils/status-code';
import type { DataTableColumns } from 'naive-ui';
import { getResourceUrl } from '@/utils/resource';
import TableActionDrawer from './components/table-action-drawer.vue';
import { getPartitionAPI } from "@/api/partition";
import { getCarouselListAPI, deleteCarouselAPI } from '@/api/carousel';
import { NCard, NImage, NIcon, NButton, NDataTable, NSpace, NPopconfirm, useMessage } from 'naive-ui';

const { loading, startLoading, endLoading } = useLoading(false);

const message = useMessage();

const visibleDrawer = ref(false);
const openDrawer = () => {
  visibleDrawer.value = true;
}

// 查看详情
const detailsData = ref<CarouselType>();
const viewDetails = (row: CarouselType) => {
  detailsData.value = row;
  openDrawer();
}

// 添加轮播图
const addCarousel = () => {
  openDrawer();
}

const columns: DataTableColumns<CarouselType> = [
  {
    key: 'id',
    title: 'ID',
    width: 90,
    align: 'center'
  },
  {
    key: 'img',
    title: '封面',
    align: 'center',
    width: 80,
    render: row => {
      return h(NImage, {
        src: getResourceUrl(row.img),
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
    key: 'partition',
    title: '分区',
    align: 'center',
    render: row => {
      return partitions.value.find(item => item.id === row.partitionId)?.name || '未知'
    }
  },
  {
    key: 'color',
    title: '主题色',
    align: 'center',
    render: row => {
      return h("div", { style: { display: "flex", alignItems: "center", justifyContent: "center" } }, [
        h("div", { style: { width: "16px", height: "16px", backgroundColor: row.color } })
      ])
    }
  },
  {
    key: 'url',
    title: 'URL',
    align: 'center',
  },
  {
    key: 'use',
    title: '是否启用',
    align: 'center',
    render: row => {
      return row.use ? '是' : '否'
    }
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
            onClick: () => viewDetails(row)
          }, { default: () => '编辑' }),
          h(NPopconfirm, {
            onPositiveClick: () => deleteCarousel(row),
          }, {
            default: () => '是否删除轮播图?',
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
  const res = await getCarouselListAPI({ page, pageSize });
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

const reviewFinish = () => {
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

// 分区
const partitions = ref<PartitionType[]>([]);
const getPartition = async () => {
  const res = await getPartitionAPI("video");
  if (res.data.code === statusCode.OK) {
    if (res.data.data.partitions) {
      partitions.value = res.data.data.partitions.filter((item: PartitionType) => {
        return item.parentId === 0;
      })
      partitions.value.unshift({
        id: 0,
        name: "首页",
        parentId: 0,
        type: constant.CONTNET_TYPE_VIDEO
      })
    }
  }
}

const deleteCarousel = async (row:CarouselType) => {
  const res = await deleteCarouselAPI(row.id);
  if (res.data.code === statusCode.OK) {
    message.success('删除成功');
    await getTableData();
  } else {
    message.error(res.data.msg);
  }
}

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