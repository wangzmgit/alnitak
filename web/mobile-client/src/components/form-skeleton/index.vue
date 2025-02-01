<template>
  <el-skeleton animated style="width: 100%;height: 100%;">
    <template #template>
      <el-skeleton-item style="width: 100%;height: 100%;" />
    </template>
  </el-skeleton>
</template>

<script setup lang="ts">
import { ElSkeleton, ElSkeletonItem } from "element-plus";
</script>