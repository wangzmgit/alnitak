<template>
  <div class="upload-video-info">
    <!-- 上传封面文件 -->
    <cover-uploader class="uploader" v-if="!loadingForm" :cover="videoForm.cover" @finish="finishUpload" />
    <el-skeleton v-else class="uploader-skeleton" animated>
      <template #template>
        <el-skeleton-item style="width: 100%;height: 100%;" />
      </template>
    </el-skeleton>
    <!-- 视频信息表单 -->
    <el-form class="info-form" :model="videoForm" label-position="left" label-width="80px">
      <el-form-item label="标题">
        <el-input v-if="!loadingForm" v-model="videoForm.title" placeholder="请输入标题" maxlength="50" show-word-limit />
        <form-skeleton v-else></form-skeleton>
      </el-form-item>
      <el-form-item label="视频简介">
        <el-input v-if="!loadingForm" v-model="videoForm.desc" placeholder="简单介绍一下视频~" maxlength="200" show-word-limit
          type="textarea" :rows="3" :autosize="descSize" resize="none" />
        <form-skeleton v-else style="height: 73px;"></form-skeleton>
      </el-form-item>
      <el-form-item label="是否原创">
        <el-switch v-if="!loadingForm" v-model="videoForm.copyright" />
        <form-skeleton v-else></form-skeleton>
      </el-form-item>
      <el-form-item label="视频分区">
        <form-skeleton v-if="loadingForm"></form-skeleton>
        <el-input v-else-if="partitionText" disabled :value="partitionText"></el-input>
        <partition-selector v-else :partitions="partitionList" @selected="selectedPartition"></partition-selector>
      </el-form-item>
      <div class="upload-next-btn">
        <el-button v-if="isModify" type="primary" @click="modifyVideoInfo">修改</el-button>
        <el-button v-else type="primary" @click="uploadInfo">下一步</el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { statusCode } from "@/utils/status-code";
import CoverUploader from "./CoverUploader.vue";
import PartitionSelector from "./PartitionSelector.vue";
import FormSkeleton from "@/components/form-skeleton/index.vue";
// import { getPartitionAPI, uploadVideoInfoAPI, modifyVideoInfoAPI } from "@leaf/apis";
import { uploadVideoInfoAPI } from "@/api/video";
import { getPartitionAPI } from '@/api/partition';
import { ElForm, ElFormItem, ElInput, ElSwitch, ElButton, ElSkeleton, ElSkeletonItem, ElMessage } from "element-plus";

const emits = defineEmits(["finish"]);
// const props = defineProps<{
//   info: VideoStatusType
// }>();

const descSize = { minRows: 3, maxRows: 3 };
const videoForm = reactive({
  vid: 0,
  title: "",
  cover: "",
  desc: "",
  copyright: true,
  partitionId: 0,
  created_at: ""
})

const isModify = ref(false);
const loadingForm = ref(true);//加载表单
const partitionText = ref("");//分区名称

//封面上传完成
const finishUpload = (cover: string) => {
  videoForm.cover = cover;
}

//选中分区
const selectedPartition = (value: number) => {
  videoForm.partitionId = value;
}

//上传视频信息
const uploadInfo = async () => {
  if (!videoForm.cover) {
    ElMessage.error("请上传视频封面");
    return;
  }
  if (!videoForm.title) {
    ElMessage.error("请填写视频标题");
    return;
  }
  if (!videoForm.partitionId) {
    ElMessage.error("请选择视频分区");
    return;
  }

  const res = await uploadVideoInfoAPI(videoForm);
  if (res.data.code === statusCode.OK) {
    emits('finish', res.data.data.videoId);
  } else {
    ElMessage.error(res.data.msg);
  }
}

//修改视频信息
const modifyVideoInfo = async () => {
  if (!videoForm.cover) {
    ElMessage.error("请上传视频封面");
    return;
  }
  if (!videoForm.title) {
    ElMessage.error("请填写视频标题");
    return;
  }

  // const res = await modifyVideoInfoAPI(videoForm);
  // if (res.data.code === statusCode.OK) {
  //   ElMessage.error("修改成功");
  // }
}

// 获取分区列表
const partitionList = ref<Array<PartitionType>>([]);//所有分区
const getPartition = async () => {
  const res = await getPartitionAPI();
  if (res.data.code === statusCode.OK) {
    partitionList.value = res.data.data.partitions;
  }
}

// 获取分区名
const getPartitionName = async (id: number) => {
  const subpartition = partitionList.value.find((item) => {
    return item.id === id;
  })

  const partition = partitionList.value.find((item) => {
    return item.id === subpartition?.parentId;
  })

  partitionText.value = `${partition?.name}/${subpartition?.name}`;
}

onMounted(async () => {
  await getPartition();
  // if (props.info.vid) {
  //   isModify.value = true;
  //   videoInfo.vid = props.info.vid;
  //   videoInfo.title = props.info.title;
  //   videoInfo.desc = props.info.desc;
  //   videoInfo.cover = props.info.cover;
  //   videoInfo.copyright = props.info.copyright;
  //   getPartitionName(props.info.partition);
  // }
  loadingForm.value = false;
})
</script>

<style lang="scss" scoped>
.uploader {
  margin: 50px auto;
}

.uploader-skeleton {
  width: 350px;
  margin: 50px auto;
  height: 200px;
  border-radius: 6px;
  background-color: #f0f2f5;
}

.info-form {
  width: calc(100% - 260px);
  margin-left: 130px;

  .upload-next-btn {
    float: right;
    margin-bottom: 20px;

    button {
      width: 160px;
      height: 40px;
    }
  }
}
</style>