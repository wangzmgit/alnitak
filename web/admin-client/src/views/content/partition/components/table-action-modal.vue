<template>
  <n-modal v-model:show="modalVisible" style="width: 500px;" preset="card" title="新建分区">
    <n-form class="info-form">
      <n-form-item label="分区名">
        <n-input v-model:value="formModel.name" placeholder="请输入名称" maxlength="20" show-count />
      </n-form-item>
      <n-form-item label="所属分区">
        <n-select v-model:value="formModel.parentId" remote :options="selectOptions" />
      </n-form-item>
      <n-form-item label="类型">
        <n-select v-model:value="formModel.type" remote :options="typeOptions" />
      </n-form-item>
      <n-space :size="24" justify="end">
        <n-button type="primary" @click="submitForm">提交</n-button>
      </n-space>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import type { FormInst, FormRules } from 'naive-ui';
import { addPartitionAPI } from "@/api/partition";
import { NButton, NFormItem, NSpace, NModal, NForm, NInput, NSelect, useMessage } from 'naive-ui';
import { statusCode } from '@/utils/status-code';

const emit = defineEmits(['update:visible', 'refresh']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  partitions: PartitionType[]
}>(), {
  visible: false,
})

const message = useMessage();

const modalVisible = computed({
  get() {
    return props.visible;
  },
  set(visible) {
    emit('update:visible', visible);
  }
});

const typeOptions = [
  { label: "视频分区", value: 0 },
  { label: "文章分区", value: 1 },
]

const selectOptions = ref<Array<{ label: string, value: number }>>([]);
const formRef = ref<HTMLElement & FormInst>();
const formModel = reactive<AddPartitionType>({
  name: '',
  parentId: 0,
  type: 0,
});

const submitForm = async () => {
  if (!formModel.name) {
    message.error("请输入分区名称");
    return;
  }

  const res = await addPartitionAPI(formModel);
  if (res.data.code === statusCode.OK) {
    modalVisible.value = false;
    emit('refresh');
  } else {
    message.error(res.data.msg || "添加失败");
  }
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    // 初始化选择器选项
    selectOptions.value = props.partitions.map((item) => {
      return {
        label: item.name,
        value: item.id
      }
    });
    selectOptions.value.unshift({
      label: '一级分区',
      value: 0
    });
  }
})
</script>