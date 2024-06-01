<template>
  <n-modal style="width: 600px;" v-model:show="modalVisible" preset="card" :title="title">
    <n-form ref="formRef" label-placement="left" :label-width="96" :model="formModel" :rules="rules">
      <n-grid :cols="24" :x-gap="18">
        <n-form-item-grid-item :span="24" label="路径" path="path">
          <n-input v-model:value="formModel.path" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="24" label="请求方法" path="method">
          <n-select v-model:value="formModel.method" :options="methodOption" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="24" label="分组" path="category">
          <n-input v-model:value="formModel.category" />
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
import { statusCode } from '@/utils/status-code';
import { addApiAPI, editApiAPI } from '@/api/api';
import type { FormInst, FormRules } from 'naive-ui';
import {
  NGrid, NButton, NSpace, NModal, NForm,
  NFormItemGridItem, NInput, NSelect, useMessage
} from 'naive-ui';

const emit = defineEmits(['update:visible', 'refresh']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  type?: 'add' | 'edit';
  editData?: ApiItemType;
}>(), {
  visible: false,
})

const message = useMessage();

const methodOption = [
  {
    label: "GET",
    value: "GET",
  },
  {
    label: "POST",
    value: "POST",
  },
  {
    label: "PUT",
    value: "PUT",
  },
  {
    label: "DELETE",
    value: "DELETE",
  },
]

const title = computed(() => {
  if (props.type === 'add') {
    return '新增API';
  }

  return '编辑API';
});

const modalVisible = computed({
  get() {
    return props.visible;
  },
  set(visible) {
    emit('update:visible', visible);
  }
});

const formRef = ref<HTMLElement & FormInst>();
const formModel = reactive<ApiFormType>({
  method: '',
  path: '',
  category: '',
  desc: '',
});

const rules: FormRules = {
  path: { required: true, message: '请输入API路径', trigger: ['blur', 'input'] },
  category: { required: true, message: '请输入API分组', trigger: ['blur', 'input'] },
  method: { required: true, message: '请选择图标', trigger: ['blur', 'change'] },
}

const closeModal = () => {
  modalVisible.value = false;
}

const handleSubmit = async () => {
  await formRef.value?.validate();

  let success = false;
  const reqFunc = props.type === 'add' ? addApiAPI : editApiAPI;
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
  formModel.method = 'GET';
  formModel.path = '';
  formModel.desc = '';
  formModel.category = '';
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    initFormModel();
    // 如果为编辑则赋值
    if (props.type === 'edit' && props.editData) {
      formModel.id = props.editData.id;
      formModel.method = props.editData.method;
      formModel.category = props.editData.category;
      formModel.path = props.editData.path;
      formModel.desc = props.editData.desc;
    }
  }
})
</script>

<style lang="scss" scoped>
.form-btn {
  width: 72px;
}
</style>
