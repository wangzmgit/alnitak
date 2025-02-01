<template>
  <div class="video-upload">
    <div v-if="videoInfo.vid === 0" class="video-uploader">
      <video-uploader @finish="finishUploadVideo"></video-uploader>
    </div>
    <div v-else>
      <upload-video-file :vid="videoInfo.vid" :resources="videoInfo.resources" />
      <div class="divider"></div>
      <div class="title-box">
        <h2 class="title">基本信息</h2>
      </div>
      <upload-video-info :info="videoInfo" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { getVideoStatusAPI } from '@/api/video';
import VideoUploader from "./components/VideoUploader.vue";
import UploadVideoInfo from "./components/UploadVideoInfo.vue";
import UploadVideoFile from './components/UploadVideoFile.vue';

const route = useRoute();

const videoInfo = ref<VideoStatusType>({
  vid: 0,
  status: 0,
  title: "",
  cover: "",
  desc: "",
  tags: "",
  copyright: false,
  partitionId: 0,
  resources: [],
  createdAt: ""
});

const finishUploadVideo = async (data: any) => {
  if (data && data.data && data.data.resource) {
    videoInfo.value.vid = data.data.resource.vid;
    getVideoStatus(videoInfo.value.vid);
  }
}

const getVideoStatus = async (vid: number) => {
  const res = await getVideoStatusAPI(vid);
  if (res.data.code === statusCode.OK) {
    videoInfo.value = res.data.data.video;
  }
}

onBeforeMount(async () => {
  const vid = Number(route.query.vid);
  if (vid) {
    getVideoStatus(vid);
  }
})
</script>


<style lang="scss" scoped>
.video-upload {
  background-color: #fff;
  width: 100%;
  min-height: 100%;
}

.video-uploader {
  padding-top: 20px;
}

.title-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  height: 50px;
  padding: 10px 20px;

  .title {
    margin: 0;
    font-size: 16px;
    color: #212121;
    font-weight: 600;
    line-height: 50px;
  }
}

.divider {
  width: 100%;
  background: #fafafa;
  height: 24px;
}
</style>