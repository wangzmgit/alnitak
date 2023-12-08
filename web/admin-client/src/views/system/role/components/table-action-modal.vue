<template>
  <n-modal style="width: 600px;" v-model:show="modalVisible" preset="card" :title="title">
    <n-form ref="formRef" label-placement="left" :label-width="96" :model="formModel" :rules="rules">
      <n-grid :cols="24" :x-gap="18">
        <n-form-item-grid-item :span="24" label="角色名" path="path">
          <n-input v-model:value="formModel.name" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="24" label="角色代码" path="category">
          <n-input v-model:value="formModel.code" :disabled="$props.type === 'edit'" />
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
import { addRoleAPI, editRoleAPI } from '@/api/role';
import {
  NGrid, NButton, NSpace, NModal, NForm,
  NFormItemGridItem, NInput, useMessage
} from 'naive-ui';

const emit = defineEmits(['update:visible', 'refresh']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  type?: 'add' | 'edit';
  editData?: RoleItemType;
}>(), {
  visible: false,
})

const message = useMessage();

const title = computed(() => {
  if (props.type === 'add') {
    return '新增角色';
  }

  return '编辑角色';
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
const formModel = reactive<RoleFormType>({
  name: '',
  code: '',
  desc: '',
});

const rules: FormRules = {
  name: { required: true, message: '请输入角色名', trigger: ['blur', 'input'] },
  code: { required: true, message: '请输入角色代码', trigger: ['blur', 'change'] },
}

const closeModal = () => {
  modalVisible.value = false;
}

const handleSubmit = async () => {
  await formRef.value?.validate();

  let success = false;
  const reqFunc = props.type === 'add' ? addRoleAPI : editRoleAPI;
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
  formModel.code = '';
  formModel.desc = '';
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    initFormModel();
    // 如果为编辑则赋值
    if (props.type === 'edit' && props.editData) {
      formModel.id = props.editData.id;
      formModel.name = props.editData.name;
      formModel.code = props.editData.code;
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
