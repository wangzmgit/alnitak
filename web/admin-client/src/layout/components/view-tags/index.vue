<template>
  <div class="view-tags-container">
    <n-scrollbar id="view-tags-scroll" ref="viewTagsScrollRef" class="view-tags-scroll" x-scrollable
      @wheel="handleScrollWheel">
      <n-space align="center" :wrap="false">
        <div v-for="(item, index) in tagList" :id="`view-tag-${index}`" :key="item.path">
          <n-tag :type="currentRoute.path === item.path ? 'primary' : 'default'"
            :bordered="!(currentRoute.path === item.path)" :closable="!item.meta.isAffix" @click="handleTagClick(index)"
            @close="handleTagClose(index)">
            {{ item.meta.label }}
            <template #icon>
              <component :is="item.iconNode"></component>
            </template>
          </n-tag>
        </div>
      </n-space>
    </n-scrollbar>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch, nextTick } from "vue";
import { useRouter } from "vue-router";
import type { ScrollbarInst } from "naive-ui";
import { NScrollbar, NTag, NSpace } from "naive-ui";
import { globalConfig } from "@/utils/global-config";
import useHistoryStore from "@/stores/modules/history-store.js";

const router = useRouter();

const historyStore = useHistoryStore();

/** 标签列表  剔除空name字段页面路由 */
const tagList = computed(() => historyStore.routerHistory.filter(item => {
  return item.name && item.path !== `/${globalConfig.baseUrl}`
}));
const handleTagClose = (index: number) => {
  historyStore.removeRouterHistory(index);
};

/** 获取某个tag对于scroll的相对位置 */
const getTagMovePosition = (index: number) => {
  const tagEl = document.getElementById(`view-tag-${index}`);
  const scrollEl = document.getElementById("view-tags-scroll");

  if (scrollEl && tagEl) {
    const scrollWidth = scrollEl.offsetWidth;
    const tagWidth = tagEl.offsetWidth;
    const tagClientLeft = tagEl.getBoundingClientRect().left;
    const scrollClientLeft = scrollEl.getBoundingClientRect().left;

    // 在滚动条视野左侧外
    if (tagClientLeft - scrollClientLeft <= 0) {
      return tagClientLeft - scrollClientLeft;
    }
    // 在滚动条视野右侧外
    else if (tagClientLeft + tagWidth > scrollWidth + scrollClientLeft) {
      return tagWidth - (scrollWidth + scrollClientLeft - tagClientLeft);
    }
  }
  return 0;
};

const viewTagsScrollRef = ref<ScrollbarInst | null>(null);
const handleTagClick = (index: number) => {
  router.push(tagList.value[index].path);
};
// 历史记录变化自动滚动至最底部
watch(tagList, () => {
  nextTick(() => {
    viewTagsScrollRef.value?.scrollTo({ left: 99999 });
  });
});

const currentRoute = computed(() => {
  return router.currentRoute.value;
});
watch(
  () => currentRoute.value.path,
  val => {
    const index = tagList.value.findIndex(item => item.path === val);
    const movePosition = getTagMovePosition(index);
    viewTagsScrollRef.value?.scrollBy({ left: movePosition });
  }
);

const handleScrollWheel = (e: WheelEvent) => {
  viewTagsScrollRef.value?.scrollBy({ left: e.deltaY });
};
</script>

<style lang="scss" scoped>
.view-tags-container {
  height: 100%;
  display: flex;
  align-items: center;

  :deep(.n-scrollbar) {
    overflow: auto hidden;

    .n-scrollbar-content,
    .n-space {
      height: 100%;
    }

    >.n-scrollbar-rail.n-scrollbar-rail--horizontal {
      bottom: 1px;
      overflow: hidden;
    }
  }

  .n-tag {
    cursor: pointer;
  }
}
</style>
