<template>
  <div class="system-menu">
    <n-card class="menu-card" :bordered="false">
      <div class="menu-card-content">
        <n-space class="search-bar" justify="space-between">
          <n-space>
            <n-button type="primary" @click="addMenu()">添加根菜单</n-button>
          </n-space>
          <n-space align="center" :size="18">
            <n-button :disabled="loading" size="small" type="primary" @click="getTableData">
              <n-icon>
                <refresh></refresh>
              </n-icon>
            </n-button>
          </n-space>
        </n-space>
        <n-data-table class="table" :row-key="(row: MenuItemType) => row.id" :columns="columns" :data="tableData"
          :loading="loading" flex-height />
        <table-action-modal v-model:visible="visible" :parent-menu="parentMenu" :type="modalType" :edit-data="editData"
          @refresh="getTableData" />
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { h, onBeforeMount, ref } from 'vue';
import { Refresh } from "@vicons/ionicons5";
import * as icons from "@vicons/ionicons5";
import useLoading from '@/hooks/loading-hooks';
import { statusCode } from '@/utils/status-code';
import type { DataTableColumns } from 'naive-ui';
import { getMenuAPI, deleteMenuAPI } from '@/api/menu';
import TableActionModal from './components/table-action-modal.vue';
import { NCard, NIcon, NButton, NDataTable, NSpace, NPopconfirm, useMessage } from 'naive-ui';

const { loading, startLoading, endLoading } = useLoading(false);

const message = useMessage();

const modalType = ref<"edit" | "add">('add');
const visible = ref(false);
const openModal = () => {
  visible.value = true;
}

// 添加菜单
const addMenu = (row?: MenuItemType) => {
  parentMenu.value = row;
  modalType.value = 'add';
  openModal();
}

// 编辑数据
const parentMenu = ref<MenuItemType>();
const editData = ref<MenuItemType>();
const editMenu = (row: MenuItemType) => {
  editData.value = row;
  if (row.parentId !== 0) {
    parentMenu.value = tableData.value.find((item) => {
      return item.id === row.parentId;
    })
  } else {
    parentMenu.value = undefined;
  }
  modalType.value = 'edit';
  openModal();
}

// 删除数据
const deleteMenu = async (row: MenuItemType) => {
  const res = await deleteMenuAPI(row.id);
  if (res.data.code === statusCode.OK) {
    message.success('删除成功');
    await getTableData();
  } else {
    message.error(res.data.msg);
  }
}

const columns: DataTableColumns<MenuItemType> = [
  {
    key: 'id',
    title: '序号',
    align: 'center'
  },
  {
    key: 'title',
    title: '标题',
    align: 'center',
    render: row => {
      return row.meta.title;
    }
  },
  {
    key: 'icon',
    title: '图标',
    align: 'center',
    render(row) {
      const iconNode = (icons as any)[row.meta.icon];
      return h(NIcon, null, { default: () => h(iconNode) });
    }
  },
  {
    key: 'path',
    title: '路径',
    align: 'center',
  },
  {
    key: 'name',
    title: '路由名',
    align: 'center',
  },
  {
    key: 'desc',
    title: '简介',
    align: 'center',
  },
  {
    key: 'hidden',
    title: '隐藏',
    align: 'center',
    render: row => {
      return row.meta.hidden ? '是' : '否';
    }
  },
  {
    key: 'actions',
    title: '操作',
    align: 'center',
    width: 320,
    render: row => {
      return h(NSpace, { justify: 'center' }, {
        default: () => {
          const btns = []
          if (row.parentId === 0) {
            btns.push(h(NButton, {
              size: 'small',
              onClick: () => addMenu(row)
            }, { default: () => '添加子菜单' }))
          }
          // 编辑
          btns.push(h(NButton, {
            size: 'small',
            onClick: () => editMenu(row)
          }, { default: () => '编辑' }))
          // 删除
          btns.push(h(NPopconfirm, {
            onPositiveClick: () => deleteMenu(row),
          }, {
            default: () => '是否删除菜单?',
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

const tableData = ref<MenuItemType[]>([]);
const getTableData = async () => {
  startLoading();
  const res = await getMenuAPI();
  if (res.data.code === statusCode.OK) {
    tableData.value = res.data.data.menus;
    endLoading();
  }
}

onBeforeMount(() => {
  getTableData();
})
</script>

<style lang="scss" scoped>
.system-menu {
  height: 100%;

  .menu-card {
    height: 100%;

    .menu-card-content {
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