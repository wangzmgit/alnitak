<template>
  <div class="common-avatar" :style="{ width: `${size}px`, height: `${size}px` }"></div>
</template>

<script setup lang="ts">

const props = defineProps<{
  url?: string;
  size: number;
  iconSize: number;
}>();
</script>

<style lang="scss" scoped>
.common-avatar {
  background-color: #cccccc;
  border-radius: 50%;
}
</style>