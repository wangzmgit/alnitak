<template>
  <div>
    <header-bar></header-bar>
    <div class="find-password">
      <steps class="steps" :current="current" :data="['填写账号', '重置密码', '操作成功']"></steps>
      <div class="form" v-if="current === 1">
        <el-form :rules="rules" :model="modifyForm" label-width="auto">
          <el-form-item label="邮箱" prop="email">
            <el-input placeholder="请输入绑定的邮箱" v-model="modifyForm.email"></el-input>
          </el-form-item>
          <el-button class="submit" type="primary" @click="checkUser">验证</el-button>
        </el-form>
      </div>
      <div v-else-if="current === 2">
        <el-form class="form" ref="formRef" :rules="rules" :model="modifyForm" label-width="auto">
          <el-form-item label="邮箱">
            <span>{{ desensiEmail(modifyForm.email) }}</span>
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input placeholder="请输入密码" type="password" v-model="modifyForm.password" />
          </el-form-item>
          <el-form-item label="确认密码" prop="reenteredPassword">
            <el-input placeholder="请输入确认密码" type="password" v-model="modifyForm.reenteredPassword" />
          </el-form-item>
          <el-form-item class="form-code" label="验证码" prop="code">
            <el-input placeholder="请输入验证码" v-model="modifyForm.code">
              <template #append>
                <el-button :disabled="disabledSend" @click="beforeSendCode">{{ sendBtnText }}</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-button class="submit" type="primary" @click="submitForm">保存</el-button>
        </el-form>
      </div>
      <div v-else-if="current === 3">
        <el-result icon="success" title="重置成功" sub-title="已成功重置密码">
          <template #extra>
            <el-button @click="goHomePage">返回首页</el-button>
          </template>
        </el-result>
      </div>
    </div>
  </div>
  <slider-captcha v-model:show="showCaptcha" :captcha-id="captchaId"></slider-captcha>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import HeaderBar from '@/components/header-bar/index.vue';
import Steps from '@/components/alnitak-steps/index.vue';
import type { FormRules, FormInstance } from "element-plus";
import { ElMessage } from "element-plus";
import { resetpwdCheckAPI, mpdifyPwdAPI } from "@/api/auth";
import SliderCaptcha from "@/components/slider-captcha/index.vue";
import { sendEmailCodeAPI } from "@/api/code";

const current = ref(1);
const captchaId = ref("");
const showCaptcha = ref(false);
const modifyForm = reactive({
  email: "",
  password: "",
  reenteredPassword: "",
  code: "",
  captchaId: ""
})

const validatePasswordSame = (rule: any, value: string): boolean => {
  return value === modifyForm.password
}

const rules: FormRules = {
  email: [
    { required: true, message: "请输入邮箱", trigger: ['blur', 'input'] },
    { type: "email", message: "请输入正确的邮箱地址", trigger: ['blur', 'input'] },
  ],
  password: [{ required: true, message: '请输入密码', trigger: ['input', 'blur'] }],
  reenteredPassword: [
    { required: true, message: '请再次输入密码', trigger: ['input', 'blur'] },
    { validator: validatePasswordSame, message: '两次密码输入不一致', trigger: ['blur', 'password-input'] }
  ],
  code: { required: true, message: '请输入验证码', trigger: ['blur', 'input'] },
}

const checkUser = async () => {
  if (!modifyForm.email) {
    ElMessage.error("邮箱不能为空");
    return;
  };

  const res = await resetpwdCheckAPI({ email: modifyForm.email, captchaId: captchaId.value });
  if (res.data.code === statusCode.CAPTCHA_REQUIRED) {
    captchaId.value = res.data.data.captchaId;
    showCaptcha.value = true;
    return;
  }
  if (res.data.code === statusCode.OK) {
    current.value = 2;
  } else {
    ElMessage.error(res.data.msg);
  }
}

// 处理显示的邮箱
const desensiEmail = (email: string) => {
  return email.replace(/(.{0,3}).*@(.*)/, "$1***@$2")
}

const disabledSend = ref(false);//禁用发送按钮
const sendBtnText = ref("发送验证码");//发送按钮文字
const beforeSendCode = async () => {
  disabledSend.value = true;
  const res = await sendEmailCodeAPI({ email: modifyForm.email, captchaId: captchaId.value });
  if (res.data.code === statusCode.OK) {
    ElMessage.success("发送成功");
    //开启倒计时
    let count = 0;
    let tag = setInterval(() => {
      if (++count >= 60) {
        clearInterval(tag);
        disabledSend.value = false;
        sendBtnText.value = "发送验证码";
        return;
      }
      sendBtnText.value = `${60 - count}秒`;
    }, 1000);
  } else {
    disabledSend.value = false;
    if (res.data.code === statusCode.CAPTCHA_REQUIRED) {
      captchaId.value = res.data.data.captchaId;
      showCaptcha.value = true;
    }
    ElMessage.error(res.data.msg);
  }
}

// 提交表单
const formRef = ref<FormInstance | null>(null);
const submitForm = async () => {
  const valid = await formRef.value?.validate();
  if (valid) {
    modifyForm.captchaId = captchaId.value;
    const res = await mpdifyPwdAPI(modifyForm);
    if (res.data.code === statusCode.OK) {
      current.value = 3;
    } else {
      ElMessage.error(res.data.msg);
    }
  } else {
    ElMessage.error("请检查输入的数据");
  }
}

const goHomePage = () => {
  navigateTo("/");
}
</script>

<style lang="scss" scoped>
.find-password {
  width: 800px;
  margin: 60px auto;

  .steps {
    width: 100%;
    margin-left: 90px;
  }

  .form {
    width: 520px;
    margin: 30px auto;


    .submit {
      width: 100%;
      margin-top: 30px;
    }
  }
}
</style>