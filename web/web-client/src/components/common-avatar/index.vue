<template>
  <div class="common-avatar" :style="{ width: `${size}px`, height: `${size}px` }">
    <el-image v-if="url" class="img" :src="getResourceUrl(url)" alt="头像">
      <template #error>
        <div class="image-viewer-slot image-slot"></div>
      </template>
    </el-image>
  </div>
</template>

<script setup lang="ts">

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
    border-radius: 50%;
  }
}
</style>