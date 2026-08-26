<template>
  <div class="config-container">
    <n-scrollbar class="config-scrollbar">
      <n-form class="form" :model="form" label-width="auto">
        <n-divider title-placement="left">转码模式</n-divider>
        <n-form-item label="执行模式">
          <n-radio-group v-model:value="form.mode">
            <n-radio-button value="local">本地进程 (local)</n-radio-button>
            <n-radio-button value="remote">远程 Worker 池 (remote)</n-radio-button>
          </n-radio-group>
        </n-form-item>
        <n-alert type="warning" style="margin-bottom: 16px;">
          <template #header>切换模式需重启服务</template>
          remote 模式需要部署并运行 transcoder-worker 二进制，
          切换前请确保 Worker 已就绪，否则转码任务不会被执行。
        </n-alert>

        <n-divider title-placement="left">编码参数</n-divider>
        <n-form-item label="GPU 加速">
          <n-switch v-model:value="form.useGpu"></n-switch>
        </n-form-item>
        <n-form-item label="视频编码">
          <n-radio-group v-model:value="codecMode">
            <n-radio-button value="h264">H.264</n-radio-button>
            <n-radio-button value="h265">H.265 (10-bit)</n-radio-button>
            <n-radio-button value="av1">AV1</n-radio-button>
          </n-radio-group>
        </n-form-item>
        <n-alert v-if="codecMode === 'av1'" type="warning" style="margin-bottom: 16px;">
          <template #header>AV1 硬件要求</template>
          AV1 硬件编码需要 <strong>NVIDIA RTX 40 系列及以上</strong> 显卡。
          不满足时自动降级到 CPU 编码。
        </n-alert>
        <n-form-item label="生成 1080p60 帧视频">
          <n-switch v-model:value="form.generate1080p60"></n-switch>
        </n-form-item>

        <n-divider title-placement="left">本地模式参数</n-divider>
        <n-form-item label="CPU 最大并发数">
          <n-input-number v-model:value="form.maxCpuConcurrency" :min="0" :max="64" style="width: 120px" />
        </n-form-item>
        <n-form-item label="GPU 最大并发数">
          <n-input-number v-model:value="form.maxGpuConcurrency" :min="0" :max="16" style="width: 120px" />
        </n-form-item>

        <n-divider title-placement="left">远程 Worker 参数</n-divider>
        <n-form-item label="Worker 并发数">
          <n-input-number v-model:value="form.workerConcurrency" :min="1" :max="32" style="width: 120px" />
        </n-form-item>
        <n-form-item label="编码并发数">
          <n-input-number v-model:value="form.encodingConcurrency" :min="0" :max="16" style="width: 120px" />
          <n-text style="margin-left: 8px; color: #888;">0=不限制</n-text>
        </n-form-item>
        <n-form-item label="队列深度上限">
          <n-input-number v-model:value="form.maxQueueDepth" :min="1" :max="200" style="width: 120px" />
          <n-text style="margin-left: 8px; color: #888;">超过此数量后新任务会被拒绝</n-text>
        </n-form-item>
        <n-form-item label="工作目录">
          <n-input v-model:value="form.workDir" placeholder="默认为系统临时目录"></n-input>
        </n-form-item>

        <div class="submit">
          <span></span>
          <n-button type="primary" @click="setConfig">保存</n-button>
        </div>
      </n-form>
    </n-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, reactive, ref, watch } from "vue";
import { statusCode } from "@/utils/status-code";
import { getTranscodingConfigAPI, setTranscodingConfigAPI } from "@/api/config";
import {
  NInput, NInputNumber, NSwitch, NForm, NFormItem, NButton,
  NDivider, NScrollbar, NAlert, NRadioGroup, NRadioButton,
  NText, useMessage
} from "naive-ui";

const message = useMessage();

const form = reactive<TranscodingConfigType>({
  mode: "local",
  useGpu: false,
  useH265: false,
  useAv1: false,
  generate1080p60: false,
  maxCpuConcurrency: 0,
  maxGpuConcurrency: 0,
  workerConcurrency: 1,
  encodingConcurrency: 0,
  maxQueueDepth: 10,
  workDir: "",
});

const getConfig = async () => {
  const res = await getTranscodingConfigAPI();
  if (res.data.code === statusCode.OK) {
    const data = res.data.data.config;
    form.mode = data.mode || "local";
    form.useGpu = data.useGpu;
    form.useH265 = data.useH265;
    form.useAv1 = data.useAv1;
    form.generate1080p60 = data.generate1080p60;
    form.maxCpuConcurrency = data.maxCpuConcurrency || 0;
    form.maxGpuConcurrency = data.maxGpuConcurrency || 0;
    form.workerConcurrency = data.workerConcurrency || 1;
    form.encodingConcurrency = data.encodingConcurrency || 0;
    form.maxQueueDepth = data.maxQueueDepth || 10;
    form.workDir = data.workDir || "";
    codecMode.value = data.useAv1 ? 'av1' : data.useH265 ? 'h265' : 'h264';
  } else {
    message.error("配置加载失败");
  }
};

const setConfig = async () => {
  const res = await setTranscodingConfigAPI(form);
  if (res.data.code === statusCode.OK) {
    message.success("修改成功，部分参数需要重启服务生效");
  } else {
    message.error(res.data.msg || "修改失败");
  }
};

const codecMode = ref<'h264' | 'h265' | 'av1'>('h264');
watch(codecMode, (mode) => {
  form.useH265 = mode === 'h265';
  form.useAv1 = mode === 'av1';
});

onBeforeMount(() => {
  getConfig();
});
</script>

<style lang="scss" scoped>
.config-container {
  height: calc(100vh - 280px);
  min-height: 400px;
}
.config-scrollbar {
  height: 100%;
}
.form {
  padding: 20px;
  padding-bottom: 40px;
}
.submit {
  display: flex;
  align-items: center;
  justify-content: space-between;
  span {
    color: #666;
  }
}
:deep(.n-input .n-input__input-el) {
  height: auto;
}
</style>
