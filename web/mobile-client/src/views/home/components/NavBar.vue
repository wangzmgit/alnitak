<template>
  <div class="switcher-header">
    <div class="switcher-header-after">
      <n-icon size="22" color="#aaa" class="icon">
        <down-icon></down-icon>
      </n-icon>
    </div>
    <div class="tabs-wrap">
      <ul ref="tabsRef" class="tabs-list header-start"
        :style="`transform: translateX(${translateX}px);${transitionDuration}`">
        <li class="tabs-item" :class="item.id === currentTab ? 'tabs-active' : ''" v-for="(item, index) in partitionList"
          :ref="el => setRefMap(index, el)" @click="partitionClick(index)">{{ item.name }}</li>
        <div class="switcher-header-anchor"
          :style="`width: ${tabAnchorWidth}px;transform: translateX(${tabAnchorLeft}px)`"></div>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeMount, nextTick } from "vue";
import { NIcon } from "naive-ui";
import { Down as DownIcon } from "@icon-park/vue-next";
import { getPartitionAPI } from "@/api/partition";
import { statusCode } from "@/utils/status-code";

const emit = defineEmits(["partitionChange"]);
const partitionList = ref<PartitionType[]>([
  { id: 0, name: "首页", parentId: 0 }
])

const currentTab = ref(-1);
const tabAnchorLeft = ref(0);
const tabAnchorWidth = ref(0);

const refMap = new Map<number, any>();
const tabsRef = ref<HTMLElement | null>(null);
const setRefMap = (key: number, el: any) => {
  refMap.set(key, el);
}

// 滑动滑动条
const slidingSlider = () => {
  if (tabsRef.value) {
    // PC端
    tabsRef.value.onmousedown = function (e) {
      let startX = e.clientX;
      document.onmousemove = function (e) {
        sliderValueChange(e, startX);
        startX = e.clientX;
      };
      document.onmouseup = function () {
        document.onmousemove = document.onmouseup = null;
      };
    };
    //移动端
    tabsRef.value.ontouchstart = function (e) {
      let startX = e.changedTouches[0].clientX;
      document.ontouchmove = function (e) {
        sliderValueChange(e, startX);
        startX = e.changedTouches[0].clientX;

      };
      document.ontouchend = function () {
        document.ontouchmove = document.ontouchend = null;
      };
    };
  }
}

//滑动条值改变
const translateX = ref(0);
const transitionDuration = ref("");
const sliderValueChange = (e: MouseEvent | TouchEvent, startX: number) => {
  let clientX: number;
  if (e instanceof MouseEvent) {
    clientX = e.clientX;
  } else {
    clientX = e.changedTouches[0].clientX;
  }

  translateX.value += (clientX - startX);
  startX = clientX;

  translateX.value = Math.min(translateX.value, 0);
  if (tabsRef.value) {
    transitionDuration.value = "";
    translateX.value = Math.max(translateX.value, tabsRef.value.clientWidth - tabsRef.value.scrollWidth);
  }
}

const partitionClick = (index: number) => {
  const tabDom = refMap.get(index);
  tabAnchorWidth.value = tabDom.clientWidth - 20;
  tabAnchorLeft.value = tabDom.offsetLeft + 10;

  let totalWidth = 0;
  if (index + 1 < partitionList.value.length) {
    for (let i = 0; i <= index + 1; i++) {
      totalWidth += refMap.get(i).clientWidth;
    }

    if (tabsRef.value?.clientWidth) {
      const offset = -(totalWidth - tabsRef.value.clientWidth);
      if (offset <= 0) {
        translateX.value = offset;
        transitionDuration.value = "transition-duration: 300ms;"
      }
    }
  }

  currentTab.value = partitionList.value[index].id;
  emit("partitionChange", currentTab.value);
}

onBeforeMount(async () => {
  const res = await getPartitionAPI("video");
  if (res.data.code === statusCode.OK) {
    const partitions = res.data.data.partitions as PartitionType[];
    partitionList.value.push(...partitions.filter(item => item.parentId === 0));
    nextTick(() => {
      partitionClick(0)
    })
  }
})

onMounted(() => {
  slidingSlider();
});
</script>

<style lang="scss" scoped>
.switcher-header {
  height: 46px;
  position: relative;
  box-sizing: border-box;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  overflow: hidden;
  background-color: #fff;

  &::before {
    content: "";
    position: absolute;
    box-sizing: border-box;
    left: 0;
    bottom: 0;
    width: 100%;
    height: 1px;
    transform: scaleY(.5);
    background-color: #e7e7e7;
  }

  .switcher-header-after {
    float: right;
    margin-top: 6px;
    order: 2;

    .icon {
      margin: 0 10px;
    }
  }

  .tabs-wrap {
    flex: 1;
    height: 100%;
    overflow: hidden;
  }

  .tabs-list {
    position: relative;
    font-size: 0;
    z-index: 1;
    padding: 0;
    margin: 0;
    flex: 1;
    white-space: nowrap;

    .tabs-item {
      font-size: 15px;
      flex-shrink: 0;
      height: 46px;
      line-height: 46px;
      display: inline-block;
      text-align: center;
      vertical-align: middle;
      white-space: nowrap;
      -webkit-user-select: none;
      -moz-user-select: none;
      user-select: none;
      padding: 0 16px;
    }

    .switcher-header-anchor {
      display: block;
      position: absolute;
      left: 0;
      bottom: 0;
      height: .53333vmin;
      border-radius: .53333vmin;
      transition-timing-function: cubic-bezier(.645, .045, .355, 1);
      transition-property: width, height, transform;
      pointer-events: none;
      background: var(--primary-color);
      transition-duration: 300ms;
    }
  }

  .header-start {
    text-align: left;
    display: flex;
    flex-direction: row;
  }
}

.tabs-active {
  color: var(--primary-color);
}
</style>