<template>
  <n-modal style="width: 600px;" v-model:show="modalVisible" preset="card" :title="title">
    <n-form ref="formRef" label-placement="left" :label-width="96" :model="formModel" :rules="rules">
      <n-grid :cols="24" :x-gap="18">
        <n-form-item-grid-item :span="24" label="标题" path="title">
          <n-input v-model:value="formModel.title" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="24" label="内容" path="content">
          <n-input v-model:value="formModel.content" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="24" label="URL" path="url">
          <n-input v-model:value="formModel.url" />
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
import { addAnnounceAPI } from '@/api/announce';
import type { FormInst, FormRules } from 'naive-ui';
import { NGrid, NButton, NSpace, NModal, NForm, NFormItemGridItem, NInput, useMessage } from 'naive-ui';

const emit = defineEmits(['update:visible', 'refresh']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  type?: 'add' | 'edit';
  editData?: AnnounceType;
}>(), {
  visible: false,
})

const message = useMessage();

const title = computed(() => {
  if (props.type === 'add') {
    return '新增公告';
  }

  return '编辑公告';
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
const formModel = reactive<AnnounceFormType>({
  id: 0,
  title: '',
  content: '',
  url: '',
});

const rules: FormRules = {
  title: { required: true, message: '请输入标题', trigger: ['blur', 'input'] },
  content: { required: true, message: '请输入内容', trigger: ['blur', 'input'] },
}

const closeModal = () => {
  modalVisible.value = false;
}

const handleSubmit = async () => {
  await formRef.value?.validate();

  let success = false;
  const res = await addAnnounceAPI(formModel);
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
  formModel.title = '';
  formModel.content = '';
  formModel.url = '';
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    initFormModel();
    // 如果为编辑则赋值
    if (props.type === 'edit' && props.editData) {
      formModel.id = props.editData.id;
      formModel.title = props.editData.title;
      formModel.content = props.editData.content;
      formModel.url = props.editData.url;
    }
  }
})
</script>

<style lang="scss" scoped>
.form-btn {
  width: 72px;
}
</style>
