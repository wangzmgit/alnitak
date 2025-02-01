<template>
  <div class="upload-video-info">
    <!-- 上传封面文件 -->
    <div class="cover">
      <div class="label">封面</div>
      <div class="cover-uploader-box">
        <cover-uploader class="uploader" v-if="!loadingForm" :cover="videoForm.cover" @finish="finishUpload" />
        <el-skeleton v-else class="uploader-skeleton" animated>
          <template #template>
            <el-skeleton-item style="width: 100%;height: 100%;" />
          </template>
        </el-skeleton>
      </div>
    </div>
    <!-- 视频信息表单 -->
    <el-form class="info-form" :model="videoForm" label-position="left" label-width="80px">
      <el-form-item label="标题" class="required">
        <el-input v-if="!loadingForm" v-model="videoForm.title" placeholder="请输入标题" maxlength="50" show-word-limit />
        <form-skeleton v-else></form-skeleton>
      </el-form-item>
      <el-form-item label="视频简介">
        <el-input v-if="!loadingForm" v-model="videoForm.desc" placeholder="简单介绍一下视频~" maxlength="200" show-word-limit
          type="textarea" :rows="3" :autosize="descSize" resize="none" />
        <form-skeleton v-else style="height: 73px;"></form-skeleton>
      </el-form-item>
      <el-form-item label="是否原创" class="required">
        <el-switch v-if="!loadingForm" :disabled="isEdit" v-model="videoForm.copyright" />
        <form-skeleton v-else></form-skeleton>
      </el-form-item>
      <el-form-item label="标签" class="required">
        <div v-if="!loadingForm" class="tags-box">
          <el-tag class="tag" v-for="tag in dynamicTags" :key="tag" closable @close="closeTag(tag)">{{ tag }}</el-tag>
          <el-input v-if="inputVisible" class="tag-input" ref="inputRef" v-model="inputValue" size="small"
            @keyup.enter="handleInputConfirm" @blur="handleInputConfirm" />
          <el-button v-else class="btn-new-tag" size="small" @click="showInput">+ 添加</el-button>
        </div>
        <form-skeleton v-else></form-skeleton>
      </el-form-item>
      <el-form-item label="视频分区" class="required">
        <form-skeleton v-if="loadingForm"></form-skeleton>
        <el-input v-else-if="props.info.partitionId !== 0" disabled :value="partitionText"></el-input>
        <partition-selector v-else :partitions="partitionList" @selected="selectedPartition"></partition-selector>
      </el-form-item>
      <div class="upload-next-btn">
        <el-button type="primary" @click="submitVideoInfo">提交</el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { statusCode } from "@/utils/status-code";
import { isLegalTag } from "@/utils/verify";
import CoverUploader from "./CoverUploader.vue";
import PartitionSelector from "./PartitionSelector.vue";
import FormSkeleton from "@/components/form-skeleton/index.vue";
import { uploadVideoInfoAPI, editVideoAPI } from "@/api/video";
import { getPartitionAPI } from '@/api/partition';
import { ElForm, ElFormItem, ElInput, ElSwitch, ElButton, ElSkeleton, ElSkeletonItem, ElMessage } from "element-plus";

const emits = defineEmits(["finish"]);
const props = defineProps<{
  info: VideoStatusType
}>();

const isEdit = ref(false);
const descSize = { minRows: 3, maxRows: 3 };
const videoForm = reactive({
  vid: 0,
  title: "",
  cover: "",
  desc: "",
  tags: "",
  copyright: true,
  partitionId: 0,
})

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

// 标签
const inputValue = ref('');
const inputVisible = ref(false);
const dynamicTags = ref<string[]>([]);
const inputRef = ref<InstanceType<typeof ElInput>>();
const showInput = () => {
  inputVisible.value = true
  nextTick(() => {
    inputRef.value!.input!.focus()
  })
}

const handleInputConfirm = () => {
  if (inputValue.value) {
    if (!dynamicTags.value.includes(inputValue.value)) {
      if (isLegalTag(inputValue.value)) {
        dynamicTags.value.push(inputValue.value);
      } else {
        ElMessage.error("标签不可包含特殊字符");
      }
    } else {
      ElMessage.error("不能重复添加标签");
    }
  }
  inputVisible.value = false;
  inputValue.value = '';
}

const closeTag = (tag: string) => {
  dynamicTags.value.splice(dynamicTags.value.indexOf(tag), 1)
}

const submitVideoInfo = async () => {
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
  if (dynamicTags.value.length < 3) {
    ElMessage.error("标签不能低于3个");
    return;
  }

  videoForm.tags = dynamicTags.value.join(',');
  const reqFunc = isEdit.value ? editVideoAPI : uploadVideoInfoAPI;
  const res = await reqFunc(videoForm);
  if (res.data.code === statusCode.OK) {
    ElMessage.success("提交成功");
    navigateTo("/upload/video-manage");
  } else {
    ElMessage.error(res.data.msg || "提交失败");
  }
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

const loadVideoInfo = () => {
  if (props.info.vid) {
    if (props.info.tags) {
      dynamicTags.value = props.info.tags.split(',');
    }
    Object.assign(videoForm, props.info);
    getPartitionName(props.info.partitionId);

    isEdit.value = props.info.partitionId !== 0;
  }
}

watch(() => props.info, () => {
  loadingForm.value = true;
  loadVideoInfo();
  nextTick(() => {
    loadingForm.value = false;
  })
})

onMounted(async () => {
  await getPartition();
  loadVideoInfo();
  loadingForm.value = false;
})
</script>

<style lang="scss" scoped>
.cover {
  display: flex;
  align-items: center;
  margin-bottom: 18px;

  .label {
    position: relative;
    width: 80px;
    font-size: 14px;
    color: #606266;
    margin-left: 130px;

    &::before {
      position: absolute;
      font-size: 12px;
      left: -8px;
      top: 0px;
      content: "*";
      color: #ff3b30;
    }
  }

  .cover-uploader-box {
    width: 169px;
    height: 127px;

    .uploader {
      width: 169px;
      height: 127px;
    }
  }



  .uploader-skeleton {
    width: 169px;
    height: 127px;
    border-radius: 3px;
    background-color: #f0f2f5;
  }
}

.tags-box {
  width: 100%;
  height: 32px;
  box-sizing: border-box;
  padding: 0 10px;
  border-radius: 4px;
  background-color: #f5f7fa;
  box-shadow: 0 0 0 1px #e4e7ed inset;

  .tag {
    margin-right: 6px;
  }

  .tag-input {
    width: 80px;
  }
}

.info-form {
  width: calc(100% - 260px);
  margin-left: 130px;

  .upload-next-btn {
    display: flex;
    justify-content: flex-end;
    padding-bottom: 16px;

    button {
      width: 160px;
      height: 40px;
    }
  }
}

.required {
  position: relative;

  &::before {
    position: absolute;
    font-size: 12px;
    left: -8px;
    top: 6px;
    content: "*";
    color: #ff3b30;
  }
}
</style>