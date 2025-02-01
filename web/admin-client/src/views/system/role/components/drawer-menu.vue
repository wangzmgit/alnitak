<template>
  <n-drawer v-model:show="modalVisible" :width="500" placement="right">
    <n-drawer-content title="菜单管理" :loading="loading">
      <n-tree ref="treeRef" block-line :data="(treeData as any)" label-field="title" :default-checked-keys="selectedKeys"
        checkable cascade expand-on-click selectable />
      <template #footer>
        <n-button type="primary" @click="submit">保存</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { computed, h, ref, watch } from 'vue';
import useLoading from '@/hooks/loading-hooks';
import { statusCode } from '@/utils/status-code';
import { editRoleHomeAPI } from '@/api/role';
import { editRoleMenuAPI, getMenuAPI } from '@/api/menu';
import { NTree, NDrawer, NDrawerContent, NButton, useMessage } from 'naive-ui';

const emit = defineEmits(['update:visible', 'refresh']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  editData: {
    roleId: number;
    checkedIds: number[];
  };
}>(), {
  visible: false,
})

const message = useMessage();
const { loading, startLoading, endLoading } = useLoading(false);

const selectedKeys = ref<number[]>();
const modalVisible = computed({
  get() {
    return props.visible;
  },
  set(visible) {
    emit('update:visible', visible);
  }
});
const closeDrawer = () => {
  modalVisible.value = false;
};

// 处理菜单数据
const handelTreeData = (data: MenuType[]): MenuTreeType[] => {
  return data.map((item) => {
    return {
      key: item.id,
      id: item.id,
      title: item.meta.title,
      children: item.children ? handelTreeData(item.children) : null,
      suffix: () => {
        if (item.parentId === 0) return null;
        return h(NButton, {
          text: true,
          type: 'primary',
          disabled: !(selectedKeys.value?.includes(item.id)),
          onClick: () => setHomePage(item.name)
        }, {
          default: () => '设为首页'
        })
      }
    }
  })
}

const treeData = ref<MenuTreeType[]>([]);
const getTreeData = async () => {
  startLoading();
  const res = await getMenuAPI();
  if (res.data.code === statusCode.OK) {
    treeData.value = handelTreeData(res.data.data.menus)
    endLoading();
  }
}

// 修改用户首页
const setHomePage = async (name: string) => {
  const res = await editRoleHomeAPI(props.editData.roleId, name);
  if (res.data.code === statusCode.OK) {
    message.success('设置成功');
  }
}

const treeRef = ref<InstanceType<typeof NTree> | null>(null);
const submit = async () => {
  if (treeRef.value) {
    const { keys: checkedKeys } = treeRef.value.getCheckedData();
    const { keys: indeterminateKeys } = treeRef.value.getIndeterminateData();
    // 去重合并
    const keys = Array.from(new Set([...checkedKeys, ...indeterminateKeys])) as number[];

    const res = await editRoleMenuAPI({ id: props.editData.roleId, menuIds: keys });
    if (res.data.code === statusCode.OK) {
      message.success("编辑成功");
      emit('refresh');
      closeDrawer();
    } else {
      message.error(res.data.msg);
    }
  }
}

watch(() => props.visible, newValue => {
  if (newValue) {
    // 如果菜单存在子项则忽略
    selectedKeys.value = [];
    const ignoreKeys: number[] = [];
    treeData.value.forEach((item) => {
      if (item.children?.length) {
        ignoreKeys.push(item.id);
      }
    })

    props.editData.checkedIds.forEach((item) => {
      if (!ignoreKeys.includes(item)) {
        selectedKeys.value?.push(item);
      }
    })
  }
});

const init = () => {
  getTreeData();
}

// 初始化
init();
</script>