<template>
  <n-drawer v-model:show="drawerVisible" :width="520">
    <n-drawer-content title="PGC 详情">
      <n-form v-if="data" label-placement="top">
        <n-grid :cols="24" :x-gap="18">
          <n-form-item-grid-item :span="12" label="PGC ID">{{ data.pgc_id }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="12" label="类型">{{ pgcTypeLabel(data.pgc_type) }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="12" label="状态">
            <n-tag size="small" :type="statusTagType(data.status)">{{ statusLabel(data.status) }}</n-tag>
          </n-form-item-grid-item>
          <n-form-item-grid-item :span="12" label="年份">{{ data.year || '-' }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="12" label="地区">{{ data.area || '-' }}</n-form-item-grid-item>
          <n-form-item-grid-item v-if="data.created_at" :span="12" label="创建时间">
            {{ new Date(data.created_at).toLocaleString() }}
          </n-form-item-grid-item>
          <n-form-item-grid-item :span="24" label="标题">{{ data.title }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="24" label="简介">{{ data.desc || '暂无' }}</n-form-item-grid-item>
          <n-form-item-grid-item v-if="data.cover" :span="24" label="封面">
            <n-image :src="getResourceUrl(data.cover)" :width="200" />
          </n-form-item-grid-item>
        </n-grid>
      </n-form>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { getResourceUrl } from '@/utils/resource';
import { NDrawer, NDrawerContent, NForm, NGrid, NFormItemGridItem, NImage, NTag } from "naive-ui";

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
