<template>
  <div class="common-avatar" :style="{ width: `${size}px`, height: `${size}px` }">
    <el-image v-if="url" class="img" :src="getResourceUrl(url)" lazy alt="头像"></el-image>
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