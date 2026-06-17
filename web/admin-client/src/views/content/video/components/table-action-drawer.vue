<template>
  <n-drawer v-model:show="drawerVisible" :width="500">
    <n-drawer-content title="视频详情">
      <n-descriptions v-if="data" label-placement="top" :column="2">
        <n-descriptions-item label="作者ID">{{ data.author.uid }}</n-descriptions-item>
        <n-descriptions-item label="用户名">{{ data.author.name }}</n-descriptions-item>
        <n-descriptions-item label="上传时间">{{ formatTime(data.createdAt) }}</n-descriptions-item>
        <n-descriptions-item label="视频标签" :span="2">
          <n-tag class="tag" v-for="(item, index) in tagsList" :key="`${item}-${index}`">{{ item }}</n-tag>
        </n-descriptions-item>
      </n-descriptions>
      <div class="video-box">
        <span>视频列表</span>
        <n-scrollbar style="max-height: 300px;">
          <div class="video-item" v-for="(item, index) in resourceList" :key="item.id">
            <div class="item-left">
              <span>P{{ index + 1 }} {{ item.title }}</span>
            </div>
            <div class="item-right">
              <n-button text @click="playVideo(item)">查看</n-button>
            </div>
          </div>
        </n-scrollbar>
      </div>
      <video-modal v-model:visible="visibleVideoModal" :resource-id="currentResourceId"></video-modal>
      <template #footer>
        <n-button class="btn" @click="closeDrawer">完成</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { formatTime } from '@/utils/format';
import { getReviewResourceListAPI } from "@/api/video";
import { statusCode } from '@/utils/status-code';
import VideoModal from './video-modal.vue';
import { NButton, NTag, NDrawer, NDrawerContent, NScrollbar, NDescriptions, NDescriptionsItem } from "naive-ui";

const emit = defineEmits(['update:visible']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  data: VideoType;
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

const tagsList = computed<string[]>(() => {
  const tags = props.data?.tags;
  if (Array.isArray(tags)) return tags;
  if (typeof tags === 'string' && tags) return tags.split(',');
  return [];
});

const closeDrawer = () => {
  drawerVisible.value = false;
}

const resourceList = ref<ResourceType[]>([]);
const getReviewResourceList = async (vid: number) => {
  const res = await getReviewResourceListAPI(vid);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.resources) {
      resourceList.value = res.data.data.resources;
    } else {
      resourceList.value = [];
    }
  }
}

const currentResourceId = ref(0);
const visibleVideoModal = ref(false);
const playVideo = (r: ResourceType) => {
  currentResourceId.value = r.id;
  visibleVideoModal.value = true;
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    if (props.data) {
      getReviewResourceList(props.data.vid);
    }
  }
})
</script>

<style lang="scss" scoped>
.tag {
  margin-right: 10px;
}

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
