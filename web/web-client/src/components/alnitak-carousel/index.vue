<template>
  <div class="carousel-area">
    <div ref="carouselContainer" class="carousel-container"
      :style="`width: ${data.containerWidth}%;transition: ${data.transition};transform: ${data.transform}`"
      @mouseover="manualSwitching(null, false)" @mouseleave="manualSwitching(null, true)">
      <div class="carousel-item" :style="`width: ${data.itemWidth}% `" v-for="item in data.playList">
        <template v-if="item.url">
          <a :href="item.url" target="_blank" class="carousel-link">
        <img class="carousel-img" :src="getResourceUrl(item.img)" :alt="item.title" />
        <div class="carousel-mask" :style="`background-color: ${item.color}`"></div>
          </a>
        </template>
        <template v-else>
          <img class="carousel-img" :src="getResourceUrl(item.img)" :alt="item.title" />
          <div class="carousel-mask" :style="`background-color: ${item.color}`"></div>
        </template>
      </div>
    </div>
    <div class="carousel-footer">
      <div class="tool">
        <p class="title" v-if="carouselList[data.currentIndex]">
          <template v-if="carouselList[data.currentIndex].url">
            <a :href="carouselList[data.currentIndex].url" target="_blank" class="carousel-title-link">
              {{ carouselList[data.currentIndex].title || "" }}
            </a>
          </template>
          <template v-else>
            {{ carouselList[data.currentIndex].title || "" }}
          </template>
        </p>
        <div class="dots" v-if="data.carouselCount > 1">
          <div class="dot" v-for="index of data.carouselCount" :key="index"
            :class="data.currentIndex === index - 1 ? 'dot-active' : ''" @click="changeCurrentIndex(index - 1)">
          </div>
        </div>
      </div>
      <div class="arrow" v-if="data.carouselCount > 1">
        <button @click="arrowClick(false)">
          <arrow-left></arrow-left>
        </button>
        <button @click="arrowClick(true)">
          <arrow-right></arrow-right>
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { nextTick, onBeforeMount, reactive, ref } from 'vue';
import { getCarouselAPI } from "@/api/carousel";
import ArrowLeft from "@/components/icons/ArrowLeftIcon.vue";
import ArrowRight from "@/components/icons/ArrowRightIcon.vue";

const props = withDefaults(defineProps<{
  partitionId?: number | string
}>(), {
  partitionId: 0
})

const data = reactive<{
  carouselCount: number
  itemWidth: number
  containerWidth: number
  currentIndex: number
  carouselTimer: null | number
  playList: Array<CarouselType>
  transition: string
  transform: string
}>({
  carouselCount: 0,//轮播图数量
  itemWidth: 100,
  containerWidth: 100,
  currentIndex: 0,
  carouselTimer: null,//自动切换计时器
  playList: [],//播放列表
  transition: "",
  transform: "",
});

const carouselList = ref<Array<CarouselType>>([]);

const carouselContainer = ref<HTMLElement | null>(null);

//初始化轮播图
const initCarousel = () => {
  data.carouselCount = carouselList.value.length;
  if (data.carouselCount <= 1) {
    data.itemWidth = 100;
    data.containerWidth = 100;
    data.playList = carouselList.value.slice();
    data.transition = "";
    data.transform = "none";
    return;
  }
  data.itemWidth = 100 / data.carouselCount;
  data.containerWidth = 100 * data.carouselCount;

  // 将最后一张图挪到最前面
  data.playList = handelSequence(carouselList.value, data.carouselCount - 1);

  //设置初始状态
  nextTick(() => {
    //设置初始偏移量
    if (carouselContainer.value) {
      carouselContainer.value.style.transform = `translateX(-${data.itemWidth}%)`;
    }
  })
}

// 开启定期切换
const startInterval = () => {
  if (data.carouselCount <= 1) return; // 只有多于一张图才自动切换
  data.carouselTimer = window.setInterval(() => {
    changeCurrentImg(true);
  }, 3000)
}

//切换当前图片
const changeCurrentImg = (isNext: boolean) => {
  let translateX = 0;
  if (isNext) {
    translateX = -data.itemWidth * 2;
    data.currentIndex = data.currentIndex + 1 >= data.carouselCount ? 0 : data.currentIndex + 1;
  } else {
    data.currentIndex = data.currentIndex - 1 < 0 ? data.carouselCount - 1 : data.currentIndex - 1;
  }
  data.transition = "transform 300ms ease 0s";
  data.transform = `translateX(${translateX}%)`;

  setTimeout(() => {
    data.transition = "transform 0ms ease 0s";
    data.transform = `translateX(-${data.itemWidth}%)`;
    if (isNext) {
      //将第一张房到最后
      data.playList = handelSequence(data.playList, 1);
    } else {
      //将最后一张放最前面
      data.playList = handelSequence(data.playList, data.carouselCount - 1);
    }

  }, 300)
}

//根据index切换图片
const changeCurrentIndex = (i: number) => {
  if (i === data.currentIndex) return;
  manualSwitching(() => {
    if (i === 0) {
      // 将最后一张图挪到最前面
      data.playList = handelSequence(carouselList.value, data.carouselCount - 1);
    } else {
      // 当前图片前一张在最前面，其他图片依次排列
      data.playList = handelSequence(carouselList.value, i - 1);
    }
    data.currentIndex = i;
  });
}

//点击切换按钮
const arrowClick = (() => {
  //节流
  var timer: number | null = null;
  return (isNext: boolean) => {
    if (timer) return;

    manualSwitching(() => {
      changeCurrentImg(isNext);
    });

    timer = window.setTimeout(() => {
      timer = null;
    }, 300);
  }
})();

//手动切换
const manualSwitching = (fn: null | (() => void), recoverable = true) => {
  if (data.carouselTimer) {
    clearInterval(data.carouselTimer);
    data.carouselTimer = null;
  }
  if (typeof fn === "function") fn();
  //判断是否需要恢复计时
  if (!recoverable) return;
  startInterval();
}

//处理图片顺序
const handelSequence = (data: any[], start: number) => {
  const before = data.slice(0, start);
  return [...data.slice(start, data.length), ...before];
}

onBeforeMount(async () => {
  const res = await getCarouselAPI(props.partitionId);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.carousels) {
      carouselList.value = res.data.data.carousels;
      initCarousel();
      startInterval();
    }
  }
})
</script>

<style lang="scss" scoped>
.carousel-area {
  width: 100%;
  height: 100%;
  overflow: hidden;
  position: relative;
}

.carousel-container {
  height: 100%;
}

.carousel-item {
  float: left;
  position: relative;
  height: 100%;

  .carousel-link {
    display: block;
    width: 100%;
    height: 100%;
  }
  
  .carousel-img {
    width: 100%;
    height: 320px;
    object-fit: cover;
  }

  .carousel-mask {
    height: 100%;
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    user-select: none;
    pointer-events: none;
    mask-image: linear-gradient(0, #2f3238 120px, transparent 220px);
    -webkit-mask-image: linear-gradient(0, #2f3238 120px, transparent 220px);
  }
}


.carousel-footer {
  box-sizing: border-box;
  position: absolute;
  bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;

  width: 100%;
  height: 70px;
  padding: 0 16px;

  .tool {

    //标题
    .title {
      color: #fff;
      font-size: 18px;
      margin-bottom: 6px;
    }

    //指示点
    .dots {
      .dot {
        display: inline-block;
        height: 4px;
        width: 20px;
        margin: 0 3px;
        border-radius: 2px;
        background-color: rgba(255, 255, 255, .3);
        cursor: pointer;
        transition: width .5s;
      }

      .dot-active {
        width: 36px;
        background-color: #fff;
      }
    }

  }

  .arrow {
    display: flex;
    flex-wrap: nowrap;

    button {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 28px;
      height: 28px;
      border-radius: 8px;

      color: #fff;
      font-size: 18px;
      background-color: rgba(255, 255, 255, .2);

      border: none;
      outline: none;
      cursor: pointer;
      user-select: none;
      -webkit-user-select: none;

      &:first-child {
        margin-right: 12px;
      }

      &:hover {
        background-color: rgba(255, 255, 255, .3);
      }
    }
  }
}

//标题高亮交互！
.carousel-title-link {
  color: inherit;
  text-decoration: none;
  transition: color 0.2s, text-shadow 0.2s;
}
.carousel-title-link:hover {
  color: #fff;
  text-shadow: 0 0 8px #fff, 0 0 16px #fff;
}
  
</style>
