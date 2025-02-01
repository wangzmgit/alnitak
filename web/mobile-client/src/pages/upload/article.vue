<template>
  <div class="article-upload">
    <div class="title-box">
      <h2 class="title">投稿信息</h2>
    </div>
    <!-- 上传封面文件 -->
    <div class="cover">
      <div class="label">封面</div>
      <div class="cover-uploader-box">
        <cover-uploader class="uploader" v-if="!loadingForm" :cover="articleForm.cover" @finish="finishUpload" />
        <el-skeleton v-else class="uploader-skeleton" animated>
          <template #template>
            <el-skeleton-item style="width: 100%;height: 100%;" />
          </template>
        </el-skeleton>
      </div>
    </div>
    <!-- 视频信息表单 -->
    <el-form class="info-form" :model="articleForm" label-position="left" label-width="80px">
      <el-form-item label="标题" class="required">
        <el-input v-if="!loadingForm" v-model="articleForm.title" placeholder="请输入标题" maxlength="50" show-word-limit />
        <form-skeleton v-else></form-skeleton>
      </el-form-item>
      <el-form-item label="是否原创" class="required">
        <el-switch v-if="!loadingForm" :disabled="isEdit" v-model="articleForm.copyright" />
        <form-skeleton v-else></form-skeleton>
      </el-form-item>
      <el-form-item label="内容分区" class="required">
        <form-skeleton v-if="loadingForm"></form-skeleton>
        <el-input v-else-if="isEdit" disabled :value="partitionText"></el-input>
        <partition-selector v-else :partitions="partitionList" @selected="selectedPartition"></partition-selector>
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
      <div class="editor-box">
        <alnitak-editor v-if="!loadingForm" :content="articleForm.content"
          @update:content="contentChange"></alnitak-editor>
        <form-skeleton v-else style="height: 300px;"></form-skeleton>
      </div>
      <div class="upload-next-btn">
        <el-button type="primary" @click="submitArticleInfo">提交</el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import type { ElInput } from 'element-plus';
import { isLegalTag } from "@/utils/verify";
import PartitionSelector from "./components/PartitionSelector.vue";
import { getPartitionAPI } from '@/api/partition';
import CoverUploader from "./components/CoverUploader.vue";
import AlnitakEditor from '@/components/alnitak-editor/index.vue';
import { uploadArticleInfoAPI, editArticleAPI, getArticleStatusAPI } from "@/api/article";

definePageMeta({
  middleware: ['article']
})

const route = useRoute();

const isEdit = ref(false);
const waitingReview = ref(false);
const articleForm = reactive({
  aid: 0,
  title: "",
  cover: "",
  tags: "",
  content: "",
  copyright: true,
  partitionId: 0,
})

const loadingForm = ref(true);//加载表单
const partitionText = ref("");//分区名称

//封面上传完成
const finishUpload = (cover: string) => {
  articleForm.cover = cover;
}

const contentChange = (val: string) => {
  articleForm.content = val;
}

//选中分区
const selectedPartition = (value: number) => {
  articleForm.partitionId = value;
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

// 获取分区列表
const partitionList = ref<Array<PartitionType>>([]);//所有分区
const getPartition = async () => {
  const res = await getPartitionAPI("article");
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

const submitArticleInfo = async () => {
  if (!articleForm.title) {
    ElMessage.error("请填写标题");
    return;
  }
  if (!articleForm.partitionId) {
    ElMessage.error("请选择分区");
    return;
  }
  if (dynamicTags.value.length < 3) {
    ElMessage.error("标签不能低于3个");
    return;
  }
  if (!articleForm.content) {
    ElMessage.error("内容不能为空");
    return;
  }

  articleForm.tags = dynamicTags.value.join(',');
  const submitFunc = articleForm.aid ? editArticleAPI : uploadArticleInfoAPI;
  const res = await submitFunc(articleForm);
  if (res.data.code === statusCode.OK) {
    ElMessage.success("提交成功");
    navigateTo("/upload/article-manage");
  } else {
    ElMessage.error(res.data.msg || "提交失败");
  }
}

const getArticleStatus = async (aid: number) => {
  const res = await getArticleStatusAPI(aid);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.article) {
      // 判断是否在审核中
      if (res.data.data.article.status = reviewCode.WAITING_REVIEW) {
        waitingReview.value = true;
      }

      Object.assign(articleForm, res.data.data.article);
      if (articleForm.tags) {
        dynamicTags.value = articleForm.tags.split(',');
      }
      getPartitionName(articleForm.partitionId);
      isEdit.value = true;
    }
  }
}

onBeforeMount(async () => {
  await getPartition();
  const aid = Number(route.query.aid);
  if (aid) {
    await getArticleStatus(aid);
  }
  loadingForm.value = false;
})
</script>

<style lang="scss" scoped>
.article-upload {
  background-color: #fff;
  width: 100%;
  min-height: 100%;

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
}

.editor-box {
  padding: 0 0 20px;
  box-sizing: border-box;
}

.cover {
  display: flex;
  align-items: center;
  margin-bottom: 18px;

  .label {
    width: 80px;
    font-size: 14px;
    color: #606266;
    margin-left: 100px;
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
  width: calc(100% - 200px);
  margin-left: 100px;

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