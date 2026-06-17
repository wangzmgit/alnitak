<template>
  <img :src="currentSrc" v-bind="$attrs" @error="onError" />
</template>

<script setup lang="ts">
/**
 * 多 OSS 容灾图片组件。
 * 1. 优先加载主线路 URL
 * 2. 加载失败时自动切换到备用线路
 * 3. 两条都失败后停止重试（避免循环）
 *
 * 用法：
 *   <oss-image :src="info.cover" alt="封面" class="cover" />
 *   <oss-image :src="getResourceUrl(info.cover)" alt="封面" />
 */
import { ref } from 'vue';
import { getResourceUrl } from '@/utils/resource';
import { getBackupUrl } from '@/utils/image-oss';

const props = defineProps<{
  src?: string;
}>();

const currentSrc = ref(props.src ? getResourceUrl(props.src) : '');
let hasTriedFallback = false;

function onError() {
  if (hasTriedFallback) return;
  hasTriedFallback = true;

  const cur = currentSrc.value;
  // 当前是主线路 → 切到备用
  if (!cur.includes('backup=true')) {
    const backup = cur ? getBackupUrl(cur) : null;
    if (backup) {
      currentSrc.value = backup;
      return;
    }
  }
  // 当前是备用（或备用 URL 不可用）→ 切回主线路
  const primary = props.src ? getResourceUrl(props.src) : '';
  if (primary && primary !== cur) {
    currentSrc.value = primary;
  }
}
</script>
