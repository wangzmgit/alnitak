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
import { onBeforeMount, reactive } from "vue";
import { statusCode } from "@/utils/status-code";
import { getOtherConfigAPI, setOtherConfigAPI } from "@/api/config";
import { NInput, NSwitch, NForm, NFormItem, NButton, NDivider, NScrollbar, NAlert, useMessage } from "naive-ui";

const message = useMessage();

const otherForm = reactive({
  allowOrigin: "*",
  prefix: "",
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