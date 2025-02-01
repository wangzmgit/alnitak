<template>
  <n-modal style="width: 700px;" v-model:show="modalVisible" preset="card" :title="title">
    <n-form ref="formRef" label-placement="left" :label-width="96" :model="formModel" :rules="rules">
      <n-grid :cols="24" :x-gap="18">
        <n-form-item-grid-item :span="12" label="路由Name" path="name">
          <n-input v-model:value="formModel.name" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="12" label="路由Path" path="path">
          <n-input v-model:value="formModel.path" />
        </n-form-item-grid-item>
        <n-form-item-grid-item v-if="isRootMenu" :span="24" label="所属菜单" path="parentMenuName">
          <n-input :value="parentMenuName" disabled />
        </n-form-item-grid-item>
        <n-form-item-grid-item v-if="isRootMenu" :span="24" label="文件路径" path="component">
          <n-input v-model:value="formModel.component" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="12" label="标题" path="title">
          <n-input v-model:value="formModel.title" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="12" label="图标" path="icon">
          <n-popover trigger="hover" placement="bottom" :width="400">
            <template #trigger>
              <n-input readonly v-model:value="formModel.icon" />
            </template>
            <icon-selector @selected="iconSelected"></icon-selector>
          </n-popover>
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="12" label="是否隐藏" path="hidden">
          <n-select v-model:value="hidden" :options="boolOption" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="12" label="KeepAlive" path="keepAlive">
          <n-select v-model:value="keepAlive" :options="boolOption" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="24" label="简介" path="desc">
          <n-input v-model:value="formModel.desc" />
        </n-form-item-grid-item>
      </n-grid>
      <n-space :size="24" justify="end">
        <n-button class="form-btn" @click="closeModal">取消</n-button>
        <n-button class="form-btn" type="primary" @click="handleSubmit">确定</n-button>
      </n-space>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import type { FormInst, FormRules } from 'naive-ui';
import { statusCode } from '@/utils/status-code';
import { addMenuAPI, editMenuAPI } from '@/api/menu';
import IconSelector from "@/components/icon-selector/index.vue";
import {
  NGrid, NButton, NPopover, NSpace, NModal, NForm,
  NFormItemGridItem, NInput, NSelect, useMessage
} from 'naive-ui';

const emit = defineEmits(['update:visible', 'refresh']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  type?: 'add' | 'edit';
  editData?: MenuItemType;
  parentMenu?: MenuItemType;
}>(), {
  visible: false,
})

const message = useMessage();

const boolOption = [
  {
    label: "是",
    value: 1,
  },
  {
    label: "否",
    value: 0,
  },
]

const title = computed(() => {
  if (props.type === 'add') {
    return '新增菜单';
  }

  return '编辑菜单';
});

const modalVisible = computed({
  get() {
    return props.visible;
  },
  set(visible) {
    emit('update:visible', visible);
  }
});

const isRootMenu = ref(false);// 是否为根菜单
const parentMenuName = ref('');
const hidden = ref(0);
const keepAlive = ref(0);
const formRef = ref<HTMLElement & FormInst>();
const formModel = reactive<MenuFormType>({
  name: '',
  path: '',
  desc: '',
  sort: 1,
  component: '',
  title: '',
  icon: '',
  hidden: false,
  keepAlive: false,
});

const rules: FormRules = {
  name: { required: true, message: '请输入路由Name', trigger: ['blur', 'input'] },
  path: { required: true, message: '请输入路由Path', trigger: ['blur', 'input'] },
  component: { required: true, message: '请输入文件路径', trigger: ['blur', 'input'] },
  title: { required: true, message: '请输入标题', trigger: ['blur', 'input'] },
  icon: { required: true, message: '请选择图标', trigger: ['blur', 'change'] },
}

// 选择图标
const iconSelected = (val: string) => {
  formModel.icon = val;
}

const closeModal = () => {
  modalVisible.value = false;
}

const handleSubmit = async () => {
  await formRef.value?.validate();

  let success = false;
  formModel.hidden = hidden.value === 1 ? true : false;
  formModel.keepAlive = keepAlive.value === 1 ? true : false;
  const reqFunc = props.type === 'add' ? addMenuAPI : editMenuAPI;
  const res = await reqFunc(formModel);
  if (res.data.code === statusCode.OK) {
    success = true;
  } else {
    message.error(res.data.msg);
  }

  if (success) {
    emit('refresh');
    closeModal();
  }
}

// 初始化formModel
const initFormModel = () => {
  formModel.id = undefined;
  formModel.name = '';
  formModel.path = '';
  formModel.desc = '';
  formModel.sort = 1;
  formModel.component = '';
  formModel.title = '';
  formModel.icon = '';
  formModel.keepAlive = false;
  formModel.hidden = false;
  formModel.parentId = 0;
  parentMenuName.value = '';
  isRootMenu.value = false;
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    initFormModel();
    // 如果为编辑则赋值
    if (props.type === 'edit' && props.editData) {
      formModel.id = props.editData.id;
      formModel.name = props.editData.name;
      formModel.path = props.editData.path;
      formModel.desc = props.editData.desc;
      formModel.sort = props.editData.sort;
      formModel.component = props.editData.component;
      formModel.title = props.editData.meta.title;
      formModel.icon = props.editData.meta.icon;
      formModel.keepAlive = props.editData.meta.keepAlive;
      formModel.hidden = props.editData.meta.hidden;
      hidden.value = props.editData.meta.hidden ? 1 : 0;
      keepAlive.value = props.editData.meta.keepAlive ? 1 : 0;
    }

    if (props.parentMenu) {
      formModel.parentId = props.parentMenu.id;
      parentMenuName.value = props.parentMenu.meta.title;
      isRootMenu.value = true;
    }
  }
})
</script>

<style lang="scss" scoped>
.form-btn {
  width: 72px;
}
</style>
