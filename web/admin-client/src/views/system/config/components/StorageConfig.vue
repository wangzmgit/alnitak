<template>
  <div class="database">
    <n-scrollbar :style="{ maxHeight: '550px' }">
      <n-form class="form" :model="storageForm" label-width="auto">
        <n-form-item label="最大图片文件大小(MB)">
          <n-input-number placeholder="文件域名" v-model:value="storageForm.maxImgSize" />
        </n-form-item>
        <n-form-item label="最大视频文件大小(MB)">
          <n-input-number placeholder="文件域名" v-model:value="storageForm.maxVideoSize" />
        </n-form-item>
        <n-form-item label="存储策略">
          <n-select v-model:value="storageForm.type" :options="ossOptions" />
        </n-form-item>
        <n-form-item v-show="storageForm.type !== 'local'" label="视频原始文件是否上传到OSS">
          <n-switch v-model:value="storageForm.uploadMp4File"></n-switch>
        </n-form-item>
        <n-form-item v-show="storageForm.type !== 'local'" label="OSS存储空间(Bucket)">
          <n-input placeholder="存储空间(Bucket)" v-model:value="storageForm.bucket" />
        </n-form-item>
        <n-form-item v-show="storageForm.type === 'aliyun'" label="OSS存储区域(Endpoint)">
          <n-input placeholder="存储区域(Endpoint)" v-model:value="storageForm.endpoint" />
        </n-form-item>
        <n-form-item v-show="storageForm.type === 'tencent'" label="存储桶AppID">
          <n-input placeholder="存储桶AppID" v-model:value="storageForm.appId" />
        </n-form-item>
        <n-form-item v-show="storageForm.type === 'tencent'" label="存储桶地域">
          <n-input placeholder="存储桶地域" v-model:value="storageForm.region" />
        </n-form-item>
        <n-form-item v-show="storageForm.type !== 'local'" label="域名">
          <n-input placeholder="域名" v-model:value="storageForm.domain" />
        </n-form-item>
        <n-form-item v-show="storageForm.type === 'qiniu'" label="是否私有">
          <n-switch v-model:value="storageForm.private"></n-switch>
        </n-form-item>
        <n-form-item v-show="storageForm.type !== 'local'" label="keyId">
          <n-input placeholder="keyId" v-model:value="storageForm.keyId" />
        </n-form-item>
        <n-form-item v-show="storageForm.type !== 'local'" label="keySecret">
          <n-input placeholder="keySecret" type="password" v-model:value="storageForm.keySecret" />
        </n-form-item>
        <div class="submit">
          <span>accesskeySecret如不需要修改请不要填写</span>
          <n-button type="primary" @click="setStorageConfig">保存</n-button>
        </div>
      </n-form>
    </n-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, reactive } from "vue";
import { statusCode } from "@/utils/status-code";
import { getStorageConfigAPI, setStorageConfigAPI } from '@/api/config';
import { NInput, NSwitch, NSelect, NInputNumber, NForm, NFormItem, NButton, NScrollbar, useMessage } from "naive-ui";

const message = useMessage();

const ossOptions = [
  {
    label: '本地',
    value: 'local'
  },
  {
    label: '阿里云',
    value: 'aliyun'
  },
  // {
  //   label: '腾讯云',
  //   value: 'tencent'
  // },
  // {
  //   label: '七牛云',
  //   value: 'qiniu'
  // },
  {
    label: 'MinIO',
    value: 'minio'
  },
];

const storageForm = reactive({
  maxImgSize: 5,
  maxVideoSize: 500,

  // 对象存储
  type: "local",
  keyId: "",
  keySecret: "",
  bucket: "",
  endpoint: "",
  appId: "",
  region: "",
  domain: "",
  private: false,
  uploadMp4File: false,
});

const getStorageConfig = async () => {
  const res = await getStorageConfigAPI();
  if (res.data.code === statusCode.OK) {
    const resData = res.data.data.config;
    storageForm.maxImgSize = resData.maxImgSize;
    storageForm.maxVideoSize = resData.maxVideoSize;
    storageForm.type = resData.type;
    storageForm.keyId = resData.keyId;
    storageForm.bucket = resData.bucket;
    storageForm.endpoint = resData.endpoint;
    storageForm.appId = resData.appId;
    storageForm.region = resData.region;
    storageForm.domain = resData.domain;
    storageForm.private = resData.private;
    storageForm.uploadMp4File = resData.uploadMp4File;
  } else {
    message.error("读取配置失败");
  }
};

const setStorageConfig = async () => {
  // 校验配置
  const msg = [];
  if (storageForm.type !== "type") {
    msg.push(...checkConfig({ 'bucket': storageForm.bucket, 'keyId': storageForm.keyId }));
  }

  // 使用不同OSS
  switch (storageForm.type) {
    case "aliyun":
      // 注释掉了与 endpoint 相关的校验
      msg.push(...checkConfig({ 'endpoint': storageForm.endpoint }));
      break;
    case "tencent":
      msg.push(...checkConfig({ 'appId': storageForm.appId, 'region': storageForm.region }));
      break;
    case "qiniu":
      msg.push(...checkConfig({ 'domain': storageForm.domain }));
      break;
    case "minio":
      // 注释掉了与 endpoint 相关的校验
       msg.push(...checkConfig({ 'endpoint': storageForm.endpoint, 'bucket': storageForm.bucket, 'keyId': storageForm.keyId }));
      break;
  }

  if (msg.length) {
    message.error(`使用${ossOptions.find(item => item.value === storageForm.type)?.label}存储${msg.join(',')}为必填项`);
    return;
  }

  const res = await setStorageConfigAPI(storageForm);
  if (res.data.code === statusCode.OK) {
    message.success("修改成功");
  } else {
    message.error(res.data.msg || "修改失败");
  }
};

const checkConfig = (config: { [key: string]: string }) => {
  const res = [];
  for (const key in config) {
    if (!config[key]) {
      res.push(key);
    }
  }
  return res;
};

onBeforeMount(() => {
  getStorageConfig();
});
</script>

<style lang="scss" scoped>
.form {
  padding: 20px;
}

.submit {
  display: flex;
  align-items: center;
  justify-content: space-between;

  span {
    color: #666;
  }
}
</style>
