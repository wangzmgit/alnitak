<template>
  <n-carousel autoplay show-arrow>
    <div class="carousel" v-for="(item, index) in carousels" :key="index">
      <img class="carousel-img" :src="getResourceUrl(item.img)" alt="轮播图" />
    </div>
  </n-carousel>
</template>

<script setup lang="ts">
import { NCarousel } from 'naive-ui';
import { ref, onBeforeMount } from 'vue';
import { getCarouselAPI } from '@/api/carousel';
import { statusCode } from '@/utils/status-code';
import { getResourceUrl } from '@/utils/resource';

const carousels = ref(Array<CarouselType>());
const getCarousel = async () => {
  const res = await getCarouselAPI("0");
  if (res.data.code === statusCode.OK) {
    carousels.value = res.data.data.carousels;
  }
}

onBeforeMount(() => {
  getCarousel();
})
</script>

<style lang="scss" scoped>
.carousel {
  width: 100%;
  height: 100%;

  .carousel-img {
    width: 100%;
    height: 100%;
    object-fit: fill;
  }
}
</style>