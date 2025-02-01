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
        <table-action-modal v-model:visible="visible" :edit-data="editData" @refresh="getTableData" />
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
import { deleteUserAPI, editUserRoleAPI, getUserListAPI } from '@/api/user';
import { getAllRoleListAPI } from '@/api/role';
import CommonAvatar from "@/components/common-avatar/index.vue";
import TableActionModal from './components/table-action-modal.vue';
import type { DataTableColumns } from 'naive-ui';
import { NCard, NIcon, NSelect, NButton, NDataTable, NPopconfirm, NSpace, useMessage } from 'naive-ui';

const { loading, startLoading, endLoading } = useLoading(false);

const message = useMessage();

const visible = ref(false);
const openModal = () => {
  visible.value = true;
}

const roleList = ref<Array<{ code: string, name: string }>>([])
const getAllRoleList = async () => {
  const res = await getAllRoleListAPI();
  if (res.data.code === statusCode.OK) {
    roleList.value = res.data.data.roles;
  }
}

// 编辑用户
const editData = ref<UserInfoType>();
const editUser = (row: UserInfoType) => {
  editData.value = row;
  openModal();
}

// 删除用户
const deleteUser = async (row: UserInfoType) => {
  const res = await deleteUserAPI(row.uid);
  if (res.data.code === statusCode.OK) {
    message.success('删除成功');
    await getTableData();
  } else {
    message.error(res.data.msg);
  }
}

const columns: DataTableColumns<UserInfoType> = [
  {
    key: 'uid',
    title: '序号',
    width: 90,
    align: 'center'
  },
  {
    key: 'avatar',
    title: '头像',
    align: 'center',
    width: 60,
    render: row => {
      return h(CommonAvatar, {
        url: row.avatar,
        size: 30,
        style: {
          marginLeft: '3px'
        }
      })
    }
  },
  {
    key: 'name',
    title: '用户名',
    align: 'center'
  },
  {
    key: 'email',
    title: '邮箱',
    align: 'center',
  },
  {
    key: 'sign',
    title: '个性签名',
    align: 'center',
  },
  {
    key: 'role',
    title: '角色',
    align: 'center',
    render: row => {
      return h(NSelect, {
        size: "small",
        defaultValue: row.role,
        labelField: "name",
        valueField: "code",
        options: roleList.value,
        onUpdateValue: (val: string) => roleChange(row, val)
      })
    }
  },
  {
    key: 'createdAt',
    title: '注册时间',
    align: 'center',
    render: row => {
      return row.createdAt ? formatTime(row.createdAt) : ""
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
            onClick: () => editUser(row)
          }, { default: () => '编辑' }),
          h(NPopconfirm, {
            onPositiveClick: () => deleteUser(row),
          }, {
            default: () => '是否删除用户?',
            trigger: () => h(NButton, {
              size: 'small',
            }, { default: () => '删除' })
          })
        ]
      })

    }
  }
]

const roleChange = async (row: UserInfoType, val: string) => {
  const res = await editUserRoleAPI({ uid: row.uid, code: val });
  if (res.data.code === statusCode.OK) {
    message.success("修改成功");
  } else {
    message.error(res.data.msg || "修改失败");
  }
}

const tableData = ref<UserInfoType[]>([]);
const getTableData = async () => {
  startLoading();
  const page = pagination.page || 1;
  const pageSize = pagination.pageSize || 1;
  const res = await getUserListAPI({ page, pageSize });
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


onBeforeMount(async () => {
  await getAllRoleList();
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