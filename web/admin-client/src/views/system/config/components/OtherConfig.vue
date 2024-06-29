<template>
  <div class="database">
    <n-form class="form" :model="otherForm" label-width="auto">
      <n-form-item label="跨域AllowOrigin">
        <n-input v-model:value="otherForm.allowOrigin"></n-input>
      </n-form-item>
      <n-form-item label="默认用户昵称前缀">
        <n-input v-model:value="otherForm.prefix"></n-input>
      </n-form-item>
      <div class="submit">
        <span></span>
        <n-button type="primary" @click="setConfig">保存</n-button>
      </div>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, reactive } from "vue";
import { statusCode } from "@/utils/status-code";
import { getOtherConfigAPI, setOtherConfigAPI } from "@/api/config";
import { NInput, NForm, NFormItem, NButton, useMessage } from "naive-ui";

const message = useMessage();

const otherForm = reactive({
  allowOrigin: "*",
  prefix: "",
});

const getConfig = async () => {
  const res = await getOtherConfigAPI();
  if (res.data.code === statusCode.OK) {
    const data = res.data.data.config;
    otherForm.allowOrigin = data.allowOrigin;
    otherForm.prefix = data.prefix;
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

:deep(.n-input .n-input__input-el) {
  height: auto;
}
</style>