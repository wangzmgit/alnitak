<template>
  <div class="follow">
    <follow-list v-if="!loading" :showBtn="true" :following="true" :userId="userId"></follow-list>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { getUserInfoAPI } from "@/api/user";
import FollowList from "@/components/follow-list/index.vue";

const loading = ref(true);
const userId = ref(0);
const getUserInfo = async () => {
  const res = await getUserInfoAPI();
  if (res.data.code === statusCode.OK) {
    userId.value = res.data.data.userInfo.uid;
    loading.value = false;
  }
}

onBeforeMount(() => {
  getUserInfo();
})
</script>