<template>
  <n-layout style="height: 100%">
    <n-result status="403" size="huge" title="403 禁止访问" description="总有些门是对你关闭的">
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
