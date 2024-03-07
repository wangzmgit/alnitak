<template>
  <div class="player-container">

  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref, watch } from 'vue';

</script>

<style lang="scss" scoped>
.player-container {
  height: 0;
  width: 100%;
  padding-bottom: 56.25%;
  position: relative;
  // margin-bottom: 40px;

  .player {
    width: 100%;
    height: 100%;
    position: absolute;
    background-color: black;
  }
}
</style>