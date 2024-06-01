<template>
  <n-drawer v-model:show="drawerVisible" :width="500">
    <n-drawer-content :title="isEdit ? '编辑轮播图' : '新增轮播图'">
      <carousel-upload :cover="carouselForm.img" @finish="finishUpload"></carousel-upload>
      <n-form class="info-form">
        <n-form-item label="URL">
          <n-input v-model:value="carouselForm.url" placeholder="请输入URL" maxlength="100" show-count />
        </n-form-item>
        <n-form-item label="标题">
          <n-input v-model:value="carouselForm.title" placeholder="请输入标题" maxlength="20" show-count />
        </n-form-item>
        <n-form-item label="主题色">
          <n-color-picker v-model:value="carouselForm.color" :show-alpha="false" />
        </n-form-item>
        <n-form-item label="所属分区">
          <n-select v-model:value="carouselForm.partitionId" label-field="name" value-field="id" remote
            :options="(partitions as any)" />
        </n-form-item>
        <n-form-item label="是否启用">
          <n-switch v-model:value="carouselForm.use" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="btn" type="primary" @click="submitForm">保存</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import { getMainColor } from '@/utils/color';
import { statusCode } from '@/utils/status-code';
import { addCarouselAPI, editCarouselAPI } from "@/api/carousel";
import CarouselUpload from './carousel-uploader.vue';
import {
  NButton, NDrawer, NDrawerContent, NForm, NFormItem,
  NInput, NSelect, NSwitch, NColorPicker, useMessage
} from "naive-ui";

const emit = defineEmits(['update:visible', 'finish']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  data: CarouselType;
  partitions: PartitionType[];
}>(), {
  visible: false,
})

const message = useMessage();

const isEdit = ref(false);
const drawerVisible = computed({
  get() {
    return props.visible;
  },
  set(visible) {
    emit('update:visible', visible);
  }
});

const carouselForm = reactive<AddOrEditCarouselType>({
  img: "",
  url: "",
  color: "",
  title: "",
  use: true,
  partitionId: 0,
})

//封面上传完成
const finishUpload = async (cover: string) => {
  carouselForm.img = cover;
  const res = await getMainColor(cover);
  if (res) {
    carouselForm.color = String(res);
  }
}

const submitForm = async () => {
  const submitFunc = isEdit.value ? editCarouselAPI : addCarouselAPI;
  const res = await submitFunc(carouselForm);
  if (res.data.code === statusCode.OK) {
    emit("finish");
    drawerVisible.value = false;
  } else {
    message.error(res.data.msg || `${isEdit.value ? '编辑' : '添加'}失败`);
  }
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    if (props.data) {
      isEdit.value = true;
      Object.assign(carouselForm, props.data);
    } else {
      isEdit.value = false;
    }
  }
})
</script>

<style lang="scss" scoped>
.btn {
  width: 100px;
  margin-left: 10px;
}
</style>
