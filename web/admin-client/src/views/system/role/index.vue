<template>
  <div class="system-role">
    <n-card class="role-card" :bordered="false">
      <div class="role-card-content">
        <n-space class="search-bar" justify="space-between">
          <n-space>
            <n-button type="primary" @click="addRole()">添加角色</n-button>
          </n-space>
          <n-space align="center" :size="18">
            <n-button :disabled="loading" size="small" type="primary" @click="getTableData">
              <n-icon>
                <refresh></refresh>
              </n-icon>
            </n-button>
          </n-space>
        </n-space>
        <n-data-table class="table" remote :columns="columns" :data="tableData" :loading="loading" :pagination="pagination"
          flex-height />
        <table-action-modal v-model:visible="visible" :type="modalType" :edit-data="editData" @refresh="getTableData" />
        <drawer-api v-model:visible="visibleApiDrawer" :edit-data="editApiData"></drawer-api>
        <drawer-menu v-model:visible="visibleMenuDrawer" :edit-data="editMenuData"></drawer-menu>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { h, onBeforeMount, reactive, ref } from 'vue';
import { Refresh } from "@vicons/ionicons5";
import useLoading from '@/hooks/loading-hooks';
import { statusCode } from '@/utils/status-code';
import { getRoleApiAPI } from '@/api/api';
import { getRoleMenuAPI } from '@/api/menu';
import { getRoleListAPI, deleteRoleAPI } from '@/api/role';
import DrawerApi from './components/drawer-api.vue';
import DrawerMenu from './components/drawer-menu.vue';
import TableActionModal from './components/table-action-modal.vue';
import type { DataTableColumns } from 'naive-ui';
import { NCard, NIcon, NButton, NDataTable, NPopconfirm, NSpace, useMessage } from 'naive-ui';

const { loading, startLoading, endLoading } = useLoading(false);

const message = useMessage();

const modalType = ref<"add" | "edit">('add');
const visible = ref(false);
const openModal = () => {
  visible.value = true;
}

// 添加角色
const addRole = () => {
  modalType.value = 'add';
  openModal();
}

// 编辑角色
const editData = ref<RoleItemType>();
const editRole = (row: RoleItemType) => {
  editData.value = row;
  modalType.value = 'edit';
  openModal();
}

// 删除角色
const deleteRole = async (row: RoleItemType) => {
  const res = await deleteRoleAPI(row.id);
  if (res.data.code === statusCode.OK) {
    message.success('删除成功');
    await getTableData();
  } else {
    message.error(res.data.msg);
  }
}

// 角色API
const visibleApiDrawer = ref(false);
const openApiDrawer = () => {
  visibleApiDrawer.value = true;
}

const editApiData = ref<{ roleId: number, checkedIds: number[] }>({ roleId: 0, checkedIds: [] });
const editRoleApi = async (row: RoleItemType) => {
  editApiData.value = {
    roleId: row.id,
    checkedIds: []
  };
  const res = await getRoleApiAPI(row.code);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.list && res.data.data.list.length) {
      editApiData.value.checkedIds = res.data.data.list.map((item: ApiItemType) => {
        return item.id;
      })
    }
  }

  openApiDrawer();
}

// 角色菜单
const visibleMenuDrawer = ref(false);
const openMenuDrawer = () => {
  visibleMenuDrawer.value = true;
}

const editMenuData = ref<{ roleId: number, checkedIds: number[] }>({ roleId: 0, checkedIds: [] });
const editRoleMenu = async (row: RoleItemType) => {
  editMenuData.value = {
    roleId: row.id,
    checkedIds: []
  };

  const res = await getRoleMenuAPI(row.code);
  if (res.data.code === 200) {
    editMenuData.value.checkedIds = res.data.data.menus.map((item: MenuType) => {
      return item.id;
    })
  }

  openMenuDrawer();
}

const columns: DataTableColumns<RoleItemType> = [
  {
    key: 'id',
    title: '序号',
    align: 'center'
  },
  {
    key: 'name',
    title: '角色名',
    align: 'center'
  },
  {
    key: 'code',
    title: '角色代码',
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
            onClick: () => editRole(row)
          }, { default: () => '编辑' }),
          h(NButton, {
            size: 'small',
            onClick: () => editRoleApi(row)
          }, { default: () => '编辑API' }),
          h(NButton, {
            size: 'small',
            onClick: () => editRoleMenu(row)
          }, { default: () => '编辑菜单' }),
          h(NPopconfirm, {
            onPositiveClick: () => deleteRole(row),
          }, {
            default: () => '是否删除角色?',
            trigger: () => h(NButton, {
              size: 'small',
            }, { default: () => '删除' })
          })
        ]
      })

    }
  }
]

const tableData = ref<RoleItemType[]>([]);
const getTableData = async () => {
  startLoading();
  const page = pagination.page || 1;
  const pageSize = pagination.pageSize || 1;
  const res = await getRoleListAPI({ page, pageSize });
  if (res.data.code === statusCode.OK) {
    tableData.value = res.data.data.list;
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
.system-role {
  height: 100%;

  .role-card {
    height: 100%;

    .role-card-content {
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