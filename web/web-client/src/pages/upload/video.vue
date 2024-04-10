<template>
  <div class="video-upload">
    <steps class="step" :data="stepsData" :current="current" :status="currentStatus"></steps>
    <div class="upload-center">
      <upload-video-info v-if="current === 1" :info="videoInfo" @finish="infoFinish" />
      <upload-video-file v-else-if="current === 2" :vid="videoInfo.vid" :resources="videoInfo.resources"
        @review="videoFinish" />
      <el-result v-else class="result" :title="getStatusText(videoInfo.status)" :status="getStatusIcon(videoInfo.status)">
        <template #extra>
          <div v-if="[reviewCode.WRONG_VIDEO_CONTENT, reviewCode.WRONG_VIDEO_INFO].includes(videoInfo.status)">
            <el-button @click="modify">编辑</el-button>
          </div>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import Steps from '@/components/alnitak-steps/index.vue';
import UploadVideoInfo from "./components/UploadVideoInfo.vue";
import UploadVideoFile from './components/UploadVideoFile.vue';
import { reviewCode } from '@/utils/review-code';
import { getVideoStatusAPI } from "@/api/video";
import { ElResult, ElButton } from 'element-plus';

const route = useRoute();

// 当前步骤
const current = ref(2);
const currentStatus = ref("process");
const stepsData = ["投稿信息", "上传视频", "审核", "完成上传"];

// 获取状态文本
const getStatusText = (state: number) => {
  switch (state) {
    case reviewCode.VIDEO_PROCESSING:
      return '处理中';
    case reviewCode.SUBMIT_REVIEW:
    case reviewCode.WAITING_REVIEW:
      return '审核中';
    case reviewCode.AUDIT_APPROVED:
      return '审核通过';
    case reviewCode.WRONG_VIDEO_INFO:
      return '视频信息存在问题';
    case reviewCode.WRONG_VIDEO_CONTENT:
      return '视频内容存在问题';
    case reviewCode.PROCESSING_FAIL:
      return '处理失败';
    default:
      return '未知';
  }
}

// 状态图标
const getStatusIcon = (state: number) => {
  switch (state) {
    case reviewCode.VIDEO_PROCESSING:
    case reviewCode.WAITING_REVIEW:
      return "info";
    case reviewCode.AUDIT_APPROVED:
      return "success";
    case reviewCode.WRONG_VIDEO_INFO:
    case reviewCode.WRONG_VIDEO_CONTENT:
    case reviewCode.PROCESSING_FAIL:
      return "error";
    default:
      return "info";
  }
}

// 修改
const modify = () => {
  if (videoInfo.value.status === reviewCode.WRONG_VIDEO_CONTENT) {
    current.value = 2;
  } else if (videoInfo.value.status === reviewCode.WRONG_VIDEO_INFO) {
    current.value = 1;
  }
}

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

const infoFinish = (vid: number) => {
  videoInfo.value.vid = vid;
  current.value = 2;
}

const videoFinish = () => {
  current.value = 3;
}

onBeforeMount(async () => {
  const vid = Number(route.query.vid);
  const modify = (route.query.modify || "").toString();
  if (!vid) return;
  const res = await getVideoStatusAPI(vid);
  if (res.data.code === statusCode.OK) {
    videoInfo.value = res.data.data.video;
    switch (videoInfo.value.status) {
      case reviewCode.CREATED_VIDEO:
        current.value = 2;
        break;
      case reviewCode.AUDIT_APPROVED:
        videoInfo.value = res.data.data.video;
        if (modify === "info") {
          current.value = 1;
        } else if (modify === "video") {
          current.value = 2;
        }
        break;
      default:
        current.value = 4; //默认结果页
        currentStatus.value = "finish";
        break;
    }
  }
})
</script>

<style lang="scss" scoped>
.video-upload {
  background-color: #fff;
  width: 100%;
  height: 100%;

  .step {
    width: calc(100% - 100px);
    margin-left: 100px;
    padding-top: 30px;
  }

  .result {
    margin-top: 50px;
  }
}
</style>