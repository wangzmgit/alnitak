<template>
  <div class="login-bg">
    <login-card :close="true" @close="closeCard" @success="loginSuccess"></login-card>
  </div>
</template>

<script setup lang="ts">
import LoginCard from "@/components/login-card/index.vue";
const emit = defineEmits(["close", "success"]);

const closeCard = () => {
  emit("close");
}

const loginSuccess = () => {
  emit("success");
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