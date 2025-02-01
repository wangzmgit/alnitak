<template>
  <div class="database">
    <n-form class="form" :model="emailForm" label-width="auto">
      <n-form-item label="发件人名称">
        <n-input placeholder="发件人名称" v-model:value="emailForm.addresser" />
      </n-form-item>
      <n-form-item label="主机">
        <n-input placeholder="主机" v-model:value="emailForm.host" />
      </n-form-item>
      <n-form-item label="端口">
        <n-input-number placeholder="端口" v-model:value="emailForm.port" />
      </n-form-item>
      <n-form-item label="邮箱地址">
        <n-input placeholder="邮箱地址" v-model:value="emailForm.user" />
      </n-form-item>
      <n-form-item label="授权码">
        <n-input placeholder="授权码" type="password" v-model:value="emailForm.pass" />
      </n-form-item>
      <div class="submit">
        <span>如不修改密码请留空</span>
        <n-button type="primary" @click="setEmailConfig">保存</n-button>
      </div>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, reactive } from "vue";
import { statusCode } from "@/utils/status-code"
import { getEmailConfigAPI, setEmailConfigAPI } from '@/api/config';
import { NInput, NInputNumber, NForm, NFormItem, NButton, useMessage } from "naive-ui";

const emailForm = reactive({
  addresser: '',
  host: '',
  port: 465,
  user: '',
  pass: ''
});

const message = useMessage();

const getEmailConfig = async () => {
  const res = await getEmailConfigAPI();
  if (res.data.code === statusCode.OK) {
    const data = res.data.data.config;
    emailForm.addresser = data.addresser;
    emailForm.host = data.host;
    emailForm.port = data.port;
    emailForm.user = data.user;
  } else {
    message.error("加载邮箱配置失败");
  }
}

const setEmailConfig = async () => {
  const res = await setEmailConfigAPI(emailForm);
  if (res.data.code === statusCode.OK) {
    message.success('修改成功');
  } else {
    message.error(res.data.msg || '修改失败');
  }
}

onBeforeMount(() => {
  getEmailConfig();
})
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