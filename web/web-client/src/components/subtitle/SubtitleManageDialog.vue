<template>
  <el-dialog v-model="visible" :title="dialogTitle" width="640" :destroy-on-close="true" @closed="onClosed">
    <div v-if="loading" class="loading-text">加载中…</div>
    <div v-else-if="!videoData" class="loading-text">无法获取视频信息</div>
    <template v-else>
      <div v-if="videoData.resources.length === 0" class="loading-text">该视频暂无分P</div>
      <div v-for="r in videoData.resources" :key="r.id" class="resource-block">
        <div class="resource-title">{{ r.title || `分P ${r.sortOrder ?? r.id}` }}</div>
        <resource-subtitle-editor :resource-short-id="subtitleResourceKey(r)" :key="`sub-${r.id}`" />
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { getVideoStatusAPI } from '@/api/video';
import ResourceSubtitleEditor from '@/pages/upload/components/ResourceSubtitleEditor.vue';

const props = defineProps<{
  vid: number | string | null;
}>();

const emit = defineEmits<{
  close: [];
}>();

const visible = ref(false);
const loading = ref(false);
const videoData = ref<VideoStatusType | null>(null);
const dialogTitle = ref('字幕管理');

watch(
  () => props.vid,
  (val) => {
    if (val !== null && val > 0) {
      visible.value = true;
      loadVideo(val);
    } else {
      visible.value = false;
    }
  },
  { immediate: true },
);

const subtitleResourceKey = (item: ResourceType) => String(item.shortId || item.id);

const loadVideo = async (vid: number | string) => {
  loading.value = true;
  videoData.value = null;
  try {
    const res = await getVideoStatusAPI(vid);
    if (res.data.code === statusCode.OK) {
      videoData.value = res.data.data.video as VideoStatusType;
      dialogTitle.value = `字幕管理 - ${videoData.value.title}`;
    }
  } catch {
    videoData.value = null;
  } finally {
    loading.value = false;
  }
};

const onClosed = () => {
  videoData.value = null;
  emit('close');
};
</script>

<style scoped>
.loading-text {
  text-align: center;
  padding: 32px 0;
  font-size: 14px;
  color: var(--font-primary-3);
}

.resource-block {
  margin-bottom: 16px;
}

.resource-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--font-primary-2);
  margin-bottom: 4px;
}
</style>
