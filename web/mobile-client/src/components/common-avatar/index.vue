<template>
  <div class="common-avatar" :style="{ width: `${size}px`, height: `${size}px` }">
    <img v-if="url" class="img" :src="getResourceUrl(url)" lazy alt="头像" />
  </div>
</template>

<script setup lang="ts">
import { getResourceUrl } from "@/utils/resource";

const props = defineProps<{
  url?: string;
  size: number;
}>();
</script>

<style lang="scss" scoped>
.common-avatar {
  background-color: #cccccc;
  border-radius: 50%;

  .img {
    width: 100%;
    height: 100%;
    border-radius: 50%;
  }
}
</style>