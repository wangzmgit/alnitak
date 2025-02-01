<template>
  <el-dialog v-model="dialogVisible" title="编辑收藏夹" width="600">
    <cover-uploader :cover="collectionForm.cover" @finish="finishUpload"></cover-uploader>
    <el-form class="info-form" :model="collectionForm" label-position="left" label-width="80px">
      <el-form-item label="标题">
        <el-input v-model="collectionForm.name" placeholder="请输入收藏夹名称" maxlength="20" show-word-limit />

      </el-form-item>
      <el-form-item label="简介">
        <el-input v-model="collectionForm.desc" placeholder="简单介绍一下视频~" maxlength="150" show-word-limit type="textarea"
          :rows="3" :autosize="descSize" resize="none" />
      </el-form-item>
      <el-form-item label="公开">
        <el-switch v-model="collectionForm.open" />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmClick">确认</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import CoverUploader from "@/components/cover-uploader/index.vue";
import { editCollectionAPI } from "@/api/collection";

const emit = defineEmits(["editFinish"]);

const dialogVisible = ref(false);
const descSize = { minRows: 3, maxRows: 3 };
const collectionForm = reactive<EditCollectionType>({
  id: 0,
  cover: "",
  name: "",
  desc: "",
  open: false,
});

//封面上传完成
const finishUpload = (cover: string) => {
  collectionForm.cover = cover;
}

const confirmClick = async () => {
  if (!collectionForm.name) {
    ElMessage.error("请输入收藏夹名称");
    return;
  }

  const res = await editCollectionAPI(collectionForm);
  if (res.data.code === statusCode.OK) {
    emit("editFinish", collectionForm);
  } else {
    ElMessage.error(res.data.msg || '编辑失败');
  }

  dialogVisible.value = false;
}

const open = (collection: CollectionType) => {
  Object.assign(collectionForm, collection);
  dialogVisible.value = true;
}

defineExpose({
  open
})
</script>

<style lang="scss" scoped></style> 