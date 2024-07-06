<template>
  <div class="tabs-wrapper">
    <div class="tab-item" :ref="el => setRefMap(item.key, el)" v-for="item in props.tabs" v-show="!item.hidden"
      :class="[currentTab === item.key ? 'tab-item--active' : '', item.disabled ? 'tab-disabled' : '']"
      @click="tabChange(item.key, item.disabled)">
      {{ item.label }}
    </div>
    <div class="tabs-bar" :style="{ left: `${barLeft}px`, width: `${barWidth}px` }"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";

const emit = defineEmits(["tabChange"]);
const props = defineProps<{
  tabs: Array<{
    key: string;
    label: string;
    disabled?: boolean;
    hidden?: boolean;
  }>
  current?: string;
}>();

const refMap = new Map<string, any>();
const setRefMap = (key: string, el: any) => {
  refMap.set(key, el);
}

const barWidth = ref(0);
const barLeft = ref(0);
const currentTab = ref(props.current || "")
const tabChange = (tab: string, disabled: boolean = false) => {
  if (disabled) return;
  currentTab.value = tab;
  resize();
  emit("tabChange", currentTab.value);
}

const resize = () => {
  const tabDom = refMap.get(currentTab.value);
  barWidth.value = tabDom.clientWidth;
  barLeft.value = tabDom.offsetLeft;
}

onMounted(() => {
  if (props.tabs[0]) {
    tabChange(props.tabs[0].key)
  }
  window.addEventListener("resize", resize);
})

onBeforeUnmount(() => {
  window.removeEventListener("resize", resize);
})
</script>

<style lang="scss" scoped>
.tabs-wrapper {
  position: relative;
  width: 100%;
  display: flex;
  flex-wrap: nowrap;
  justify-content: space-evenly;

  .tab-item {
    font-weight: 400;
    box-sizing: border-box;
    vertical-align: bottom;
    cursor: pointer;
    white-space: nowrap;
    flex-wrap: nowrap;
    display: inline-flex;
    align-items: center;
    color: var(--font-primary-1);
    font-size: 16px;
    background-clip: padding-box;
    padding: 10px 0;
    transition: box-shadow .3s cubic-bezier(.4, 0, .2, 1),
      color .3s cubic-bezier(.4, 0, .2, 1),
      background-color .3s cubic-bezier(.4, 0, .2, 1),
      border-color .3s cubic-bezier(.4, 0, .2, 1);
  }

  .tab-item--active {
    color: var(--primary-color);
  }

  .tabs-bar {
    position: absolute;
    bottom: 0;
    // width: 64px;
    height: 2px;
    border-radius: 1px;
    background-color: var(--primary-color);
    transition: left 0.3s;
    -moz-transition: left 0.3s;
    -webkit-transition: left 0.3s;
    -o-transition: left 0.3s;
  }
}

.tab-disabled {
  color: #9499a0 !important;
  cursor: not-allowed !important;

  &:hover {
    color: #9499a0 !important;
  }
}
</style>