<template>
  <n-scrollbar class="grid-wrapper">
    <n-grid :cols="4" style="height: 300px">
      <n-grid-item v-for="item of listData" :key="item">
        <div class="icon-wrapper" @click="onIconClick(item)">
          <n-icon size="20">
            <component :is="(icons as any)[item]"> </component>
          </n-icon>
          <n-ellipsis :line-clamp="1" style="font-size: 12px">{{ item }}</n-ellipsis>
        </div>
      </n-grid-item>
    </n-grid>
  </n-scrollbar>
  <div class="pagination">
    <n-pagination :page="page" :page-size="pageSize" :page-slot="5" :item-count="total" @update-page="onUpdatePage" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { NScrollbar, NEllipsis, NGrid, NGridItem, NIcon, NPagination } from 'naive-ui';
import * as icons from "@vicons/ionicons5";

const emit = defineEmits(["selected"]);

const page = ref(1);
const pageSize = 40;
const iconsKeys = Object.keys(icons);
const total = iconsKeys.length;
const listData = ref(iconsKeys.slice(0, pageSize))
const onUpdatePage = (val: number) => {
  page.value = val;
  const start = (page.value - 1) * pageSize;;
  listData.value = iconsKeys.slice(start, start + pageSize);
}

const onIconClick = (item: string) => {
  emit('selected', item);
}
</script>

<style lang="scss" scoped>
.grid-wrapper {
  .icon-wrapper {
    cursor: pointer;
    border: 1px solid #f5f5f5;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2px;
  }
}

.pagination {
  display: flex;
  justify-content: end;
  margin: 2px 0;
}
</style>