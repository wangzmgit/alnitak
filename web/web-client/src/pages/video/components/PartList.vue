<template>
  <div class="part-list">
    <div class="part-head">
      <span class="title">分段列表</span>
      <span class="part">({{ current }}/{{ resources?.length }})</span>
    </div>
    <el-scrollbar max-height="340px">
      <ul class="list-box">
        <li :class="['list-item', current - 1 === index ? 'active-part' : '']" v-for="(item, index) in resources"
          @click="changePart(index)">
          <div class="item-content">
            <span class="part-num">P{{ index + 1 }}</span>
            <span class="part-title">{{ item.title }}</span>
          </div>
          <div class="duration">{{ toDuration(item.duration) }}</div>
        </li>
      </ul>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { toDuration } from "@/utils/format";

const emits = defineEmits(['change']);
const props = withDefaults(defineProps<{
  active?: number
  resources: Array<ResourceType>
}>(), {
  active: 1
})

const current = ref(props.active);

const changePart = (part: number) => {
  current.value = part + 1;
  emits('change', part + 1)
}
</script>

<style lang="scss" scoped>
.part-list {
  position: relative;
  border-radius: 6px;
  background-color: #F1F2F3;
  margin-bottom: 18px;


  .part-head {
    display: flex;
    padding: 14px 16px 0;
    align-items: center;
    justify-content: space-between;

    .part {
      color: #9499A0;
      font-size: 13px;
    }
  }
}

.list-box {
  padding: 0 6px;
  list-style: none;
  color: #18191c;
  margin: 0;
  padding: 0 10px 0 6px;

  .list-item {
    box-sizing: border-box;
    display: flex;
    justify-content: space-between;
    overflow: hidden;
    width: 100%;
    padding: 0 10px;
    height: 30px;
    line-height: 30px;
    color: #18191C;
    margin: 5px 0;
    transition: all 0.3s;
    cursor: pointer;
    white-space: nowrap;
    overflow: hidden;
    font-size: 13px;

    .item-content {
      display: flex;
      align-items: center;
      flex-shrink: 1;
      overflow: hidden;

      .part-num {
        margin-right: 10px;
      }

      .part-title {
        display: block;
        overflow: hidden;
        text-overflow: ellipsis;
        flex-shrink: 1;
      }
    }

    .duration {
      color: #9499A0;
    }
  }
}

.active-part {
  color: var(--primary-hover-color) !important;
  background-color: #fff;
  border-radius: 2px;
}
</style>