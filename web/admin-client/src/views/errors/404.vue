<template>
  <n-layout style="height: 100%">
    <n-result status="404" size="huge" title="404 资源不存在" description="生活总归带点荒谬">
      <template #footer>
        <n-button @click="handleReturn">回到过去</n-button>
      </template>
    </n-result>
  </n-layout>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router';
import { NLayout, NResult, NButton } from "naive-ui";

const router = useRouter();

const handleReturn = () => {
  router.back();
};
</script>
