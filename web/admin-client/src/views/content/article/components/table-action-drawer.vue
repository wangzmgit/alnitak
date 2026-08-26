<template>
  <n-drawer v-model:show="drawerVisible" :width="500">
    <n-drawer-content title="视频详情">
      <n-descriptions v-if="data" label-placement="top" :column="2">
        <n-descriptions-item label="作者ID">{{ data.author.uid }}</n-descriptions-item>
        <n-descriptions-item label="用户名">{{ data.author.name }}</n-descriptions-item>
        <n-descriptions-item label="内容标签" :span="2">
          <n-tag class="tag" v-for="(item, index) in tagsList" :key="`${item}-${index}`">{{ item }}</n-tag>
        </n-descriptions-item>
      </n-descriptions>
      <div class="video-box" v-if="data">
        <span>内容预览</span>
        <text-editor :content="data.content"></text-editor>
      </div>
      <template #footer>
        <n-button class="btn" @click="drawerVisible = false">完成</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import { statusCode } from '@/utils/status-code';
import TextEditor from './text-editor.vue';
import { NButton, NTag, NDrawer, NDrawerContent, NDescriptions, NDescriptionsItem } from "naive-ui";

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

const tagsList = computed<string[]>(() => {
  const tags = props.data?.tags;
  if (Array.isArray(tags)) return tags;
  if (typeof tags === 'string' && tags) return tags.split(',');
  return [];
});

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
