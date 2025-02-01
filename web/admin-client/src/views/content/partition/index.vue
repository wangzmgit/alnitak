<template>
  <div class="content-partition">
    <n-card class="partition-card" :bordered="false">
      <div class="partition-card-content">
        <n-space class="search-bar" justify="space-between">
          <n-space>
            <n-button type="primary" @click="openModal()">添加分区</n-button>
          </n-space>
          <n-space align="center" :size="18">
            <n-radio-group v-model:value="partitionType" class="radio-group" @update:value="partitionTypeChange">
              <n-radio class="radio" :key="0" value="video">视频</n-radio>
              <n-radio class="radio" :key="1" value="article">文章</n-radio>
            </n-radio-group>
            <n-button :disabled="loading" size="small" type="primary" @click="getPartition">
              <n-icon>
                <refresh></refresh>
              </n-icon>
            </n-button>
          </n-space>
        </n-space>
        <div class="table-container">
          <div class="table">
            <span class="title">一级分区</span>
            <n-data-table :row-key="(row: MenuItemType) => row.id" :columns="columns" :data="partitions"
              :loading="loading" remote />
          </div>
          <div class="table">
            <span class="title">二级分区</span>
            <n-data-table :row-key="(row: MenuItemType) => row.id" :columns="columns" :data="subpartition"
              :loading="loading" remote />
          </div>
        </div>
        <table-action-modal v-model:visible="visible" :partitions="partitions" @refresh="getPartition" />
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { h, onBeforeMount, ref } from 'vue';
import { Refresh } from "@vicons/ionicons5";
import { statusCode } from '@/utils/status-code';
import { getPartitionAPI, deletePartitionAPI } from "@/api/partition";
import useLoading from '@/hooks/loading-hooks';
import tableActionModal from './components/table-action-modal.vue';
import type { DataTableColumns } from 'naive-ui';
import {
  NCard, NIcon, NButton, NDataTable, NSpace,
  NRadioGroup, NRadio, NPopconfirm, useMessage
} from 'naive-ui';

const { loading, startLoading, endLoading } = useLoading(false);

const visible = ref(false);
const openModal = () => {
  visible.value = true;
}

const message = useMessage();

const columns: DataTableColumns<PartitionType> = [
  {
    key: 'id',
    title: 'ID',
    align: 'center',
    width: 120,
  },
  {
    key: 'name',
    title: '名称',
    align: 'center',
  },
  {
    key: 'actions',
    title: '操作',
    align: 'center',
    width: 220,
    render: row => {
      return h(NSpace, { justify: 'center' }, {
        default: () => {
          const btns = []
          if (row.parentId === 0) {
            btns.push(h(NButton, {
              size: 'small',
              onClick: () => showSubpartition(row.id)
            }, { default: () => '查看子分区' }))
          }
          // 删除
          btns.push(h(NPopconfirm, {
            onPositiveClick: () => deletePartition(row),
          }, {
            default: () => '是否删除分区?',
            trigger: () => h(NButton, {
              size: 'small',
            }, { default: () => '删除' })
          }))

          return btns;
        }
      })

    }
  }
]

//获取所有分区
const allPartition = ref<PartitionType[]>([]); // 所有分区
const partitions = ref<PartitionType[]>([]); // 一级分区
const subpartition = ref<PartitionType[]>([]);// 二级分区
const getPartition = async () => {
  startLoading();
  const res = await getPartitionAPI(partitionType.value);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.partitions) {
      allPartition.value = res.data.data.partitions;
      // 添加分区使用的一级分区
      partitions.value = allPartition.value.filter((item) => {
        return item.parentId === 0;
      });

      // 初始化子分区显示
      if (partitions.value && partitions.value[0]) {
        showSubpartition(partitions.value[0].id);
      } else {
        subpartition.value = [];
      }
    }
  }
  endLoading();
}

// 显示子分区
const showSubpartition = (id: number) => {
  subpartition.value = allPartition.value.filter((item) => {
    return item.parentId === id;
  });
}

const deletePartition = async (row: PartitionType) => {
  const res = await deletePartitionAPI(row.id);
  if (res.data.code === statusCode.OK) {
    message.success('删除成功');
    await getPartition();
  } else {
    message.error(res.data.msg || "删除失败");
  }
}

const partitionType = ref<"video" | "article">("video");
const partitionTypeChange = () => {
  getPartition();
}

onBeforeMount(async () => {
  await getPartition();
})
</script>


<style lang="scss" scoped>
.content-partition {
  height: 100%;

  .partition-card {
    height: 100%;

    .partition-card-content {
      height: 100%;
      display: flex;
      flex-direction: column;

      .search-bar {
        padding-bottom: 12px;

        .radio-group {
          height: 36px;
          display: flex;
          align-items: center;

          .radio {
            margin-left: 12px;
          }
        }
      }

      .table-container {
        display: flex;
        justify-content: space-between;

        .table {
          width: calc(50% - 30px);

          .title {
            font-size: 16px;
            font-weight: 500;
            line-height: 36px;
          }
        }
      }
    }
  }
}
</style>