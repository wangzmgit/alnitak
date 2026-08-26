<template>
  <div class="follow">
    <follow-list :following="true" :userId="userId"></follow-list>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
const route = useRoute();
const userId = ref(parseInt(route.params.id.toString()));
</script>
