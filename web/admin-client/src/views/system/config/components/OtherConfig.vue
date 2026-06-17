<template>
  <div class="config-container">
    <n-scrollbar class="config-scrollbar">
      <n-form class="form" :model="otherForm" label-width="auto">
        <n-divider title-placement="left">基础配置</n-divider>
        <n-form-item label="跨域AllowOrigin">
          <n-input v-model:value="otherForm.allowOrigin"></n-input>
        </n-form-item>
        <n-form-item label="默认用户昵称前缀">
          <n-input v-model:value="otherForm.prefix"></n-input>
        </n-form-item>
        <n-form-item label="生成1080p60帧视频">
          <n-switch v-model:value="otherForm.generate1080p60"></n-switch>
        </n-form-item>
        <n-form-item label="转码开启gpu加速">
          <n-switch v-model:value="otherForm.useGpu"></n-switch>
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
          不满足时自动降级到 CPU 编码（libsvtav1），速度显著变慢，建议仅在 GPU H.265 无法满足需求时选用。
        </n-alert>

        <n-divider title-placement="left">服务器配置</n-divider>
        <n-form-item label="HTTP端口">
          <n-input v-model:value="otherForm.serverPort" placeholder="默认9000"></n-input>
        </n-form-item>
        <n-form-item label="启用HTTPS">
          <n-switch v-model:value="otherForm.sslEnabled"></n-switch>
        </n-form-item>
        <n-form-item v-show="otherForm.sslEnabled" label="HTTPS端口">
          <n-input v-model:value="otherForm.sslPort" placeholder="默认443"></n-input>
        </n-form-item>
        <n-form-item v-show="otherForm.sslEnabled" label="证书文件路径">
          <n-input v-model:value="otherForm.sslCertFile" placeholder="如: ./cert/server.crt"></n-input>
        </n-form-item>
        <n-form-item v-show="otherForm.sslEnabled" label="私钥文件路径">
          <n-input v-model:value="otherForm.sslKeyFile" placeholder="如: ./cert/server.key"></n-input>
        </n-form-item>

        <n-alert v-if="otherForm.sslEnabled" type="warning" style="margin-bottom: 16px;">
          修改服务器配置后需要重启服务才能生效
        </n-alert>

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
import { getOtherConfigAPI, setOtherConfigAPI } from "@/api/config";
import { NInput, NSwitch, NForm, NFormItem, NButton, NDivider, NScrollbar, NAlert, NRadioGroup, NRadioButton, useMessage } from "naive-ui";

const message = useMessage();

const otherForm = reactive({
  allowOrigin: "*",
  prefix: "",
  generate1080p60: false,
  useGpu: false,
  useH265: false,
  useAv1: false,
  // 服务器配置
  serverPort: "9000",
  sslEnabled: false,
  sslPort: "443",
  sslCertFile: "",
  sslKeyFile: "",
});

const getConfig = async () => {
  const res = await getOtherConfigAPI();
  if (res.data.code === statusCode.OK) {
    const data = res.data.data.config;
    otherForm.allowOrigin = data.allowOrigin;
    otherForm.prefix = data.prefix;
    otherForm.generate1080p60 = data.generate1080p60;
    otherForm.useGpu = data.useGpu;
    otherForm.useH265 = data.useH265;
    otherForm.useAv1 = data.useAv1;
    codecMode.value = data.useAv1 ? 'av1' : data.useH265 ? 'h265' : 'h264';
    // 服务器配置
    otherForm.serverPort = data.serverPort || "9000";
    otherForm.sslEnabled = data.sslEnabled || false;
    otherForm.sslPort = data.sslPort || "443";
    otherForm.sslCertFile = data.sslCertFile || "";
    otherForm.sslKeyFile = data.sslKeyFile || "";
  } else {
    message.error("配置加载失败");
  }
}

const setConfig = async () => {
  const res = await setOtherConfigAPI(otherForm);
  if (res.data.code === statusCode.OK) {
    message.success("修改成功");
  } else {
    message.error(res.data.msg || "修改失败");
  }
}

const codecMode = ref<'h264' | 'h265' | 'av1'>('h264');
watch(codecMode, (mode) => {
  otherForm.useH265 = mode === 'h265';
  otherForm.useAv1 = mode === 'av1';
});

onBeforeMount(() => {
  getConfig();
})

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