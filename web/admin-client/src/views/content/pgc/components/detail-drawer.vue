<template>
  <n-drawer v-model:show="drawerVisible" :width="520">
    <n-drawer-content title="PGC 详情">
      <n-descriptions v-if="data" label-placement="top" :column="2">
        <n-descriptions-item label="PGC ID">{{ data.pgc_id }}</n-descriptions-item>
        <n-descriptions-item label="类型">{{ pgcTypeLabel(data.pgc_type) }}</n-descriptions-item>
        <n-descriptions-item label="状态">
          <n-tag size="small" :type="statusTagType(data.status)">{{ statusLabel(data.status) }}</n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="年份">{{ data.year || '-' }}</n-descriptions-item>
        <n-descriptions-item label="地区">{{ data.area || '-' }}</n-descriptions-item>
        <n-descriptions-item v-if="data.created_at" label="创建时间">
          {{ new Date(data.created_at).toLocaleString() }}
        </n-descriptions-item>
        <n-descriptions-item label="标题" :span="2">{{ data.title }}</n-descriptions-item>
        <n-descriptions-item label="简介" :span="2">{{ data.desc || '暂无' }}</n-descriptions-item>
        <n-descriptions-item v-if="data.cover" label="封面" :span="2">
          <n-image :src="getResourceUrl(data.cover)" :width="200" />
        </n-descriptions-item>
      </n-descriptions>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { getResourceUrl } from '@/utils/resource';
import { NDrawer, NDrawerContent, NDescriptions, NDescriptionsItem, NImage, NTag } from "naive-ui";

const emit = defineEmits(['update:visible']);
const props = withDefaults(defineProps<{
  visible: boolean;
  data: PGCReviewRow;
}>(), {
  visible: false,
});

const drawerVisible = computed({
  get() { return props.visible; },
  set(v) { emit('update:visible', v); }
});

const pgcTypeLabel = (t: number) => {
  const m: Record<number, string> = { 1: '国创', 2: '日漫', 3: '纪录片', 4: '电影', 5: '电视剧' };
  return m[t] ?? `类型${t}`;
};

const statusLabel = (s: number) => {
  if (s === -1) return '已下架';
  const m: Record<number, string> = { 0: '草稿', 100: '已提交', 200: '审核中', 300: '已通过', 400: '已驳回' };
  return m[s] ?? String(s);
};

const statusTagType = (s: number): 'default' | 'success' | 'warning' | 'error' | 'info' => {
  if (s === 300) return 'success';
  if (s === 400) return 'error';
  if (s === -1) return 'warning';
  if (s === 100 || s === 200) return 'info';
  return 'default';
};
</script>
