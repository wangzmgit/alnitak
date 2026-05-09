<template>
  <div class="subtitle-panel">
    <div class="subtitle-title">字幕</div>
    <div v-if="loading" class="subtitle-muted">加载中…</div>
    <ul v-else class="track-list">
      <li v-for="t in tracks" :key="t.id" class="track-row">
        <span class="track-name">{{ t.label || t.lang }}</span>
        <span class="track-lang">{{ t.lang }}</span>
        <el-tag v-if="t.isDefault" size="small" type="success" class="tag-def">默认</el-tag>
        <el-popconfirm title="确定删除该字幕？" width="220" @confirm="onDelete(t.id)">
          <template #reference>
            <el-button link type="danger" size="small">删除</el-button>
          </template>
        </el-popconfirm>
      </li>
      <li v-if="tracks.length === 0" class="subtitle-muted empty">暂无字幕，可上传 .srt / .vtt</li>
    </ul>
    <div class="subtitle-form">
      <el-input
        v-model="form.lang"
        placeholder="语言 zh-Hans"
        size="small"
        maxlength="20"
        show-word-limit
        class="inp-lang"
      />
      <el-input
        v-model="form.label"
        placeholder="显示名（可选）"
        size="small"
        maxlength="64"
        show-word-limit
        class="inp-label"
      />
      <el-checkbox v-model="form.isDefault" size="small">默认</el-checkbox>
      <el-upload
        :show-file-list="false"
        accept=".srt,.vtt"
        :before-upload="beforeSubtitleFile"
        :http-request="handleSubtitleUpload"
      >
        <el-button type="primary" size="small" :loading="uploading">上传字幕</el-button>
      </el-upload>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import type { UploadRequestOptions } from 'element-plus';
import { ElMessage } from 'element-plus';
import { getSubtitleListAPI, uploadSubtitleAPI, deleteSubtitleAPI } from '@/api/subtitle';

const props = defineProps<{
  /** 分 P shortId，无则传数字 id 字符串 */
  resourceShortId: string;
}>();

const loading = ref(false);
const uploading = ref(false);
const tracks = ref<SubtitleTrackItemType[]>([]);

const form = ref({
  lang: 'zh-Hans',
  label: '',
  isDefault: false,
});

const loadTracks = async () => {
  if (!props.resourceShortId) return;
  loading.value = true;
  try {
    const res = await getSubtitleListAPI(props.resourceShortId);
    if (res.data.code === statusCode.OK) {
      tracks.value = (res.data.data?.tracks as SubtitleTrackItemType[]) ?? [];
    } else {
      tracks.value = [];
    }
  } catch {
    tracks.value = [];
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  void loadTracks();
});

watch(
  () => props.resourceShortId,
  () => {
    void loadTracks();
  },
);

const beforeSubtitleFile = (file: File) => {
  const name = file.name.toLowerCase();
  const ok = name.endsWith('.srt') || name.endsWith('.vtt');
  if (!ok) {
    ElMessage.error('仅支持 .srt 或 .vtt');
  }
  return ok;
};

const handleSubtitleUpload = async (opt: UploadRequestOptions) => {
  const lang = form.value.lang.trim();
  if (!lang) {
    ElMessage.warning('请填写语言代码');
    return;
  }
  const fd = new FormData();
  fd.append('resourceShortId', props.resourceShortId);
  fd.append('lang', lang);
  const label = form.value.label.trim();
  if (label) {
    fd.append('label', label);
  }
  if (form.value.isDefault) {
    fd.append('isDefault', 'true');
  }
  fd.append('file', opt.file as File);

  uploading.value = true;
  try {
    const res = await uploadSubtitleAPI(fd);
    if (res.data.code === statusCode.OK) {
      ElMessage.success('字幕已上传');
      form.value.label = '';
      form.value.isDefault = false;
      await loadTracks();
    } else {
      ElMessage.error(res.data.msg || '上传失败');
    }
  } catch (e: unknown) {
    ElMessage.error('上传失败');
  } finally {
    uploading.value = false;
  }
};

const onDelete = async (id: number) => {
  try {
    const res = await deleteSubtitleAPI(id);
    if (res.data.code === statusCode.OK) {
      ElMessage.success('已删除');
      await loadTracks();
    } else {
      ElMessage.error(res.data.msg || '删除失败');
    }
  } catch {
    ElMessage.error('删除失败');
  }
};
</script>

<style lang="scss" scoped>
.subtitle-panel {
  margin-top: 12px;
  padding-top: 10px;
  border-top: 1px dashed var(--fill-1, #e8e8e8);
}

.subtitle-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--font-primary-2);
  margin-bottom: 8px;
}

.subtitle-muted {
  font-size: 12px;
  color: var(--font-primary-3);
}

.track-list {
  list-style: none;
  margin: 0 0 10px;
  padding: 0;
}

.track-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  padding: 4px 0;
  flex-wrap: wrap;
}

.track-name {
  color: var(--font-primary-1);
}

.track-lang {
  color: var(--font-primary-3);
  font-family: monospace;
}

.tag-def {
  margin-right: 4px;
}

.empty {
  padding: 4px 0;
}

.subtitle-form {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.inp-lang {
  width: 140px;
  max-width: 40%;
}

.inp-label {
  width: 160px;
  max-width: 45%;
}
</style>
