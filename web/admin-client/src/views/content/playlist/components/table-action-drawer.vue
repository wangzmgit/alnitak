<template>
  <n-drawer v-model:show="drawerVisible" :width="500">
    <n-drawer-content title="合集详情">
      <n-descriptions v-if="data" label-placement="top" :column="2">
        <n-descriptions-item label="合集ID">{{ data.id }}</n-descriptions-item>
        <n-descriptions-item label="作者">{{ data.author?.name }} (UID:{{ data.uid }})</n-descriptions-item>
        <n-descriptions-item label="标题" :span="2">{{ data.title }}</n-descriptions-item>
        <n-descriptions-item label="简介" :span="2">{{ data.desc || '暂无' }}</n-descriptions-item>
        <n-descriptions-item v-if="data.cover" label="封面" :span="2">
          <n-image :src="getResourceUrl(data.cover)" :width="200" />
        </n-descriptions-item>
        <n-descriptions-item label="视频数">{{ data.videoCount }}</n-descriptions-item>
        <n-descriptions-item label="浏览量">{{ data.views }}</n-descriptions-item>
        <n-descriptions-item label="公开">{{ data.isOpen ? '是' : '否' }}</n-descriptions-item>
        <n-descriptions-item label="创建时间">{{ formatTime(data.createdAt) }}</n-descriptions-item>
      </n-descriptions>

      <div class="video-box">
        <span>合集视频列表</span>
        <n-scrollbar style="max-height: 300px;">
          <div class="video-item" v-for="(item, index) in videoList" :key="item.vid || `${item.title}-${index}`">
            <div class="item-left">
              <span>P{{ index + 1 }} {{ item.title }}</span>
            </div>
          </div>
          <div v-if="videoList.length === 0" style="padding: 10px; color: #999;">暂无视频</div>
        </n-scrollbar>
      </div>
      <template #footer>
        <n-button class="btn" @click="closeDrawer">完成</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { formatTime } from '@/utils/format';
import { statusCode } from '@/utils/status-code';
import { getPlaylistVideoListAPI } from "@/api/playlist";
import { getResourceUrl } from '@/utils/resource';
import { NButton, NDrawer, NDrawerContent, NScrollbar, NDescriptions, NDescriptionsItem, NImage } from "naive-ui";

const emit = defineEmits(['update:visible']);
const props = withDefaults(defineProps<{
  visible: boolean;
  data: PlaylistType;
}>(), {
  visible: false,
})

const drawerVisible = computed({
  get() {
    return props.visible;
  },
  set(visible) {
    emit('update:visible', visible);
  }
});

const closeDrawer = () => {
  drawerVisible.value = false;
}

const videoList = ref<PlaylistVideoType[]>([]);

const getVideoList = async (id: number) => {
  const res = await getPlaylistVideoListAPI(id, 1, 200);
  if (res.data.code === statusCode.OK) {
    videoList.value = res.data.data.videos || [];
  } else {
    videoList.value = [];
  }
}

watch(() => props.visible, (newVal) => {
  if (newVal && props.data) {
    getVideoList(props.data.id);
  }
})
</script>

<style lang="scss" scoped>
.video-box {
  width: 100%;
  padding-bottom: 30px;

  .video-item {
    height: 36px;
    font-size: 14px;
    padding: 0 10px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-radius: 3px;
    border: 1px solid #efeff5;
    margin: 10px 0;

    .item-left {
      cursor: pointer;
    }
  }
}

.btn {
  width: 100px;
  margin-left: 10px;
}
</style>
