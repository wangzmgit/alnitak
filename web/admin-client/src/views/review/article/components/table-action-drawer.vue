<template>
  <n-drawer v-model:show="drawerVisible" :width="500">
    <n-drawer-content title="视频详情">
      <n-form v-if="data" label-placement="top">
        <n-grid :cols="24" :x-gap="18">
          <n-form-item-grid-item :span="12" label="作者ID">{{ data.author.uid }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="12" label="用户名">{{ data.author.name }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="24" label="内容标签">
            <n-tag class="tag" v-for="item in  data.tags.split(',')">{{ item }}</n-tag>
          </n-form-item-grid-item>
        </n-grid>
      </n-form>
      <div class="video-box">
        <span>内容预览</span>
        <text-editor :content="data.content"></text-editor>
      </div>
      <template #footer>
        <n-button class="btn" @click="openModal">不通过</n-button>
        <n-button class="btn" type="primary" @click="reviewVideoApproved">通过</n-button>
      </template>
    </n-drawer-content>
    <review-modal v-model:visible="visibleModal" :aid="props.data.aid" @finish="reviewFinish"></review-modal>
  </n-drawer>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import { statusCode } from '@/utils/status-code';
import { reviewArticleApprovedAPI } from "@/api/review";
import ReviewModal from './review-modal.vue';
import TextEditor from './text-editor.vue';
import { NButton, NTag, NDrawer, NDrawerContent, NForm, NGrid, NFormItemGridItem } from "naive-ui";

const emit = defineEmits(['update:visible', 'finish']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  data: ArticleType;
}>(), {
  visible: false,
})

const drawerVisible = computed({
  get() {
    return props.visible;
  },
  set(visible) {
    emit('update:visible', visible);
  }
});

const visibleModal = ref(false);
const openModal = () => {
  visibleModal.value = true;
}

const reviewVideoApproved = async () => {
  const res = await reviewArticleApprovedAPI({ aid: props.data.aid });
  if (res.data.code === statusCode.OK) {
    reviewFinish();
  }
}

const reviewFinish = () => {
  visibleModal.value = false;
  drawerVisible.value = false;
  emit("finish");
}
</script>

<style lang="scss" scoped>
.tag {
  margin-right: 10px;
}

.video-box {
  width: 100%;
  padding-bottom: 30px;
}

.btn {
  width: 100px;
  margin-left: 10px;
}
</style>
