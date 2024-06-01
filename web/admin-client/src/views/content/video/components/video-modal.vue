<template>
  <n-modal v-model:show="modalVisible" style="width: 700px;" preset="card" title="视频预览">
    <video-player ref="videoPlayerRef"></video-player>
    <n-space :size="24" justify="end">
      <n-button class="form-btn" type="primary" @click="handleClose">关闭</n-button>
    </n-space>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, ref, watch, nextTick } from "vue";
import { NModal, NSpace, NButton } from "naive-ui";
import VideoPlayer from "@/components/video-player/index.vue";

const emit = defineEmits(['update:visible']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  resourceId: number
}>(), {
  visible: false,
})

const modalVisible = computed({
  get() {
    return props.visible;
  },
  set(visible) {
    emit('update:visible', visible);
  }
});

const handleClose = async () => {
  modalVisible.value = false;
}

const videoPlayerRef = ref<InstanceType<typeof VideoPlayer> | null>(null);
watch(() => props.visible, (newVal) => {
  if (newVal) {
    nextTick(() => {
      videoPlayerRef.value?.loadVideo(props.resourceId);
    })
  }
})
</script>

<style lang="scss" scoped>
.form-btn {
  width: 72px;
}
</style>

