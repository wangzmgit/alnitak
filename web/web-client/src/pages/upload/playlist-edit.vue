<template>
  <div class="playlist-edit">
    <p class="title">{{ isEdit ? '编辑合集' : '创建合集' }}</p>
    <div class="form-container">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px" label-position="top">
        <el-form-item label="合集标题" prop="title">
          <el-input v-model="form.title" maxlength="100" show-word-limit placeholder="请输入合集标题"></el-input>
        </el-form-item>
        <el-form-item label="合集封面" prop="cover">
          <cover-uploader :cover="form.cover" @finish="coverFinish"></cover-uploader>
        </el-form-item>
        <el-form-item label="合集简介">
          <el-input v-model="form.desc" type="textarea" :rows="4" maxlength="2000" show-word-limit
            placeholder="请输入合集简介"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ isEdit ? '保存修改' : '创建合集' }}
          </el-button>
          <el-button @click="navigateTo('/upload/playlist-manage')">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onBeforeMount } from 'vue';
import { addPlaylistAPI, editPlaylistAPI } from '@/api/playlist';
import { statusCode } from '@/utils/status-code';
const route = useRoute();
const isEdit = computed(() => !!Number(route.query.id));
const submitting = ref(false);


const formRef = ref();

// 表单验证规则
const rules = ref({
  title: [
    { required: true, message: "合集标题不能为空", trigger: "blur" }
  ],
  cover: [
    { required: true, message: "合集封面不能为空", trigger: "change" }
  ]
});
const form = ref({
  id: 0,
  title: '',
  cover: '',
  desc: '',
});

const coverFinish = (url: string) => {
  form.value.cover = url;
}

const handleSubmit = async () => {
  // 表单验证
  if (!formRef.value) return;
  await formRef.value.validate((valid: boolean) => {
    if (!valid) {
      return;
    }
  });

  submitting.value = true;
  try {
    if (isEdit.value) {
      const res = await editPlaylistAPI({
        id: form.value.id,
        title: form.value.title,
        cover: form.value.cover,
        desc: form.value.desc,
      });
      if (res.data.code === statusCode.OK) {
        ElMessage.success('修改成功，等待审核');
        navigateTo('/upload/playlist-manage');
      } else {
        ElMessage.error(res.data.msg || '修改失败');
      }
    } else {
      const res = await addPlaylistAPI({
        title: form.value.title,
        cover: form.value.cover,
        desc: form.value.desc,
      });
      if (res.data.code === statusCode.OK) {
        const playlistId = res.data.data.id;
        try {
          await ElMessageBox.confirm('合集创建成功，等待审核。是否立即添加视频？', '提示', {
            confirmButtonText: '添加视频',
            cancelButtonText: '返回列表',
            type: 'success',
          });
          navigateTo(`/upload/playlist-video?id=${playlistId}`);
        } catch {
          navigateTo('/upload/playlist-manage');
        }
      } else {
        ElMessage.error(res.data.msg || '创建失败');
      }
    }
  } finally {
    submitting.value = false;
  }
}

onBeforeMount(async () => {
  const id = Number(route.query.id);
  if (id) {
    // 从API获取合集信息
    const { getMyPlaylistAPI } = await import('@/api/playlist');
    const res = await getMyPlaylistAPI();
    if (res.data.code === statusCode.OK) {
      const list = res.data.data.playlists || [];
      const target = list.find((p: any) => p.id === id);
      if (target) {
        form.value = {
          id: target.id,
          title: target.title,
          cover: target.cover,
          desc: target.desc,
        };
      }
    }
  }
})
</script>

<style lang="scss" scoped>
.playlist-edit {
  padding: 0 18px;
  height: 100%;
  box-sizing: border-box;
  background-color: var(--bg-elev-1);

  .title {
    font-size: 18px;
    margin: 0;
    padding: 16px 0 10px;
  }

  .form-container {
    max-width: 600px;
    padding: 20px 0;
  }
}
</style>
