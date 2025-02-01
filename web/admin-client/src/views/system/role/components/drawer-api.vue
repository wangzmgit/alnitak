<template>
  <n-drawer v-model:show="modalVisible" :width="500" placement="right">
    <n-drawer-content title="API管理" :loading="loading">
      <n-scrollbar>
        <n-tree ref="treeRef" cascade block-line :data="treeData" :default-checked-keys="selectedKeys" checkable
          expand-on-click selectable />
      </n-scrollbar>
      <template #footer>
        <n-button type="primary" @click="submit">保存</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { computed, onBeforeMount, ref, watch } from 'vue';
import { v1 as uuid } from "uuid";
import { statusCode } from '@/utils/status-code';
import useLoading from '@/hooks/loading-hooks';
import { getAllApiListAPI, editRoleApiAPI } from '@/api/api';
import { NTree, NDrawer, NDrawerContent, NButton, NScrollbar, useMessage } from 'naive-ui';

const emit = defineEmits(['update:visible', 'refresh']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  type?: 'add' | 'edit';
  editData: {
    roleId: number;
    checkedIds: number[];
  };
}>(), {
  visible: false,
})

const message = useMessage();
const { loading, startLoading, endLoading } = useLoading(false);

const selectedKeys = ref<string[]>([]);
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

const treeData = ref<ApiTreeType[]>([]);
const getTreeData = async () => {
  startLoading();
  const res = await getAllApiListAPI();
  if (res.data.code === statusCode.OK) {
    const tree: Record<string, ApiTreeType> = {};
    res.data.data.list.forEach((item: ApiItemType) => {
      const category = item.category;
      if (!tree[category]) {
        tree[category] = {
          key: uuid(),
          isLeaf: false,
          label: category,
          children: []
        };
      }
      tree[category].children?.push({
        id: item.id,
        key: `${item.id}`,
        isLeaf: true,
        label: item.desc
      });
    });
    treeData.value = Object.values(tree);
    endLoading();
  }
}

const treeRef = ref<InstanceType<typeof NTree> | null>(null);
const submit = async () => {
  if (treeRef.value) {
    let keys: number[] = [];
    const { options } = treeRef.value.getCheckedData();
    options.forEach((item) => {
      if (item && item.isLeaf && item.id) {
        keys.push(item.id as number);
      }
    })

    //原数组不存在新数组存在表示添加
    const addIds = keys.filter((v) => {
      return selectedKeys.value.indexOf(v.toString()) == -1
    })

    //原数组存在新数组不存在表示删除
    const removeIds = selectedKeys.value.filter((v) => {
      return keys.indexOf(parseInt(v)) == -1
    }).map((item) => {
      return parseInt(item);
    })

    const res = await editRoleApiAPI({ id: props.editData.roleId, addIds, removeIds });
    if (res.data.code === statusCode.OK) {
      message.success("编辑成功");
      emit('refresh');
      closeDrawer();
    } else {
      message.error(res.data.msg);
    }
  }
}

watch(() => props.visible, (newValue) => {
  if (newValue) {
    selectedKeys.value = props.editData.checkedIds.map((item) => {
      return item.toString();
    });
  }
});

onBeforeMount(() => {
  getTreeData();
})
</script>