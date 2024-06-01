<template>
  <n-modal v-model:show="modalVisible" style="width: 700px;" preset="card" title="审核不通过">
    <n-form ref="formRef" :style="{ maxWidth: '640px' }">
      <n-grid :cols="24" :x-gap="18" v-for="(item, index) in dynamicForm">
        <n-form-item-gi :span="12" label="位置">
          <n-select v-model:value="item.position" :options="positionOptions" />
        </n-form-item-gi>
        <n-form-item-gi :span="12" label="问题">
          <n-select v-model:value="item.question" :options="questionOptions" />
        </n-form-item-gi>
      </n-grid>
    </n-form>
    <n-space :size="24" justify="space-between">
      <n-button class="form-btn" @click="addItem">追加</n-button>
      <n-button class="form-btn" type="primary" @click="handleSubmit">确定</n-button>
    </n-space>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { NModal, NForm, NGrid, NFormItemGi, NSelect, NSpace, NButton } from "naive-ui";
import { reviewArticleFailedAPI } from "@/api/review";
import { reviewCode } from "@/utils/review-code";
import { statusCode } from "@/utils/status-code";

const emit = defineEmits(['update:visible', 'finish']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  aid: number;
}>(), {
  visible: false,
})

const positionOptions = ref([
  { label: "封面", value: 0, },
  { label: "标题", value: 1, },
  { label: "内容", value: 2, },
])
const questionOptions = [
  { label: "其他", value: 0, },
  { label: "色情低俗", value: 1, },
  { label: "违法违禁", value: 2, },
  { label: "危险行为", value: 3, },
  { label: "血腥暴力", value: 4, },
  { label: "侵权著作权", value: 5, },
  { label: "传播不良价值观", value: 6, },
]

const modalVisible = computed({
  get() {
    return props.visible;
  },
  set(visible) {
    emit('update:visible', visible);
  }
});

const dynamicForm = ref([
  { position: 0, question: 0 },
])

const addItem = () => {
  dynamicForm.value.push({ position: 0, question: 0 });
}

const getPositionText = (position: number) => {
  return positionOptions.value.find(item => item.value === position)?.label;
}

const getQuestionText = (position: number) => {
  return questionOptions.find(item => item.value === position)?.label;
}

const handleSubmit = async () => {
  let remark = '';
  for (const item of dynamicForm.value) {
    remark += `${getPositionText(item.position)}存在${getQuestionText(item.question)}问题;`
  }

  const res = await reviewArticleFailedAPI({
    aid: props.aid,
    status: reviewCode.REVIEW_FAILED,
    remark: remark
  });
  if(res.data.code === statusCode.OK) {
    emit("finish");
  }
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    dynamicForm.value = [{ position: 0, question: 0 }];
  }
})
</script>

<style lang="scss" scoped>
.form-btn {
  width: 72px;
}
</style>

