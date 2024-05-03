<template>
  <div class="login-bg">
      <login-card :close="true" @close="closeCard" @success="closeCard"></login-card>
  </div>
</template>

<script setup lang="ts">
import LoginCard from "@/components/login-card/index.vue";
const emits = defineEmits(["close"]);

const closeCard = () => {
  emits("close");
}
</script>

<style lang="scss" scoped>
.login-bg {
  top: 0;
  z-index: 999;
  position: fixed;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.45);
}
</style>