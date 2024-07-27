<template>
  <div class="bg">
    <div class="head-box">
      <div class="head-text">
        <span>您好，</span>
        <br>
        <span>欢迎注册 {{ globalConfig.title }}</span>
      </div>
    </div>
    <div class="login-div">
      <div class="login-box">
        <n-form ref="formRef" :rules="rules" :model="registerForm" label-placement="left" label-width="auto">
          <n-form-item label="邮箱" path="email">
            <n-input placeholder="请输入邮箱" v-model:value="registerForm.email" />
          </n-form-item>
          <n-form-item label="密码" path="password">
            <n-input placeholder="请输入密码" v-model:value="registerForm.password" type="password" />
          </n-form-item>
          <n-form-item label="验证码" path="code">
            <n-input placeholder="请输入验证码" v-model:value="registerForm.code" />
            <n-button :disabled="disabledSend" @click="sendEmailCode">{{ sendBtnText }}</n-button>
          </n-form-item>
          <div class="card-btn">
            <n-button type="primary" @click="handelRegister()">注册</n-button>
            <n-button @click="goLogin">返回登录</n-button>
          </div>
        </n-form>
      </div>
    </div>
    <slider-captcha v-model:show="showCaptcha" :captcha-id="registerForm.captchaId"
      @success="captchaSuccess"></slider-captcha>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { reactive, ref } from 'vue';
import { registerAPI } from "@/api/auth";
import { sendEmailCodeAPI } from '@/api/code';
import { globalConfig } from "@/utils/global-config";
import { statusCode } from "@/utils/status-code";
import type { FormInst, FormRules } from 'naive-ui';
import SliderCaptcha from "@/components/slider-captcha/index.vue";
import { NForm, NFormItem, NButton, NInput, useMessage } from 'naive-ui';

const message = useMessage();

// 显示滑块验证
let captchaTrigger = "";
const showCaptcha = ref(false);
const captchaSuccess = () => {
  if (captchaTrigger === "register") {
    handelRegister();
  } else {
    sendEmailCode();
  }
}

const registerForm = reactive<UserRegisterType>({
  email: "",
  password: "",
  code: "",
  captchaId: ""
});

const rules: FormRules = {
  email: [
    { required: true, message: "请输入邮箱", trigger: ['blur', 'input'] },
    { type: "email", message: "请输入正确的邮箱地址", trigger: ['blur', 'input'] },
  ],
  password: { required: true, message: '请输入密码', trigger: ['blur', 'input'] },
  code: { required: true, message: '请输入验证码', trigger: ['blur', 'input'] },
}

const formRef = ref<FormInst | null>(null);
const handelRegister = async () => {
  await formRef.value?.validate();

  const res = await registerAPI(registerForm);
  switch (res.data.code) {
    case statusCode.CAPTCHA_REQUIRED:
      captchaTrigger = "register";
      registerForm.captchaId = res.data.data.captchaId;
      showCaptcha.value = true;
      break;
    case statusCode.OK:
      goLogin();
      message.success("注册成功");
      break;
    default:
      message.error(res.data.msg);
  }
}

//返回登录
const router = useRouter();
const goLogin = () => {
  router.push({ name: 'Login' });
}

// 发送验证码
const disabledSend = ref(false);//禁用发送按钮
const sendBtnText = ref('发送验证码');//发送按钮文字
const startCountdown = () => {
  let count = 0;
  let tag = setInterval(() => {
    if (++count >= 60) {
      clearInterval(tag);
      disabledSend.value = false;
      sendBtnText.value = '发送验证码';
      return;
    }
    sendBtnText.value = `${60 - count}秒`;
  }, 1000);
}
const sendEmailCode = async () => {
  if (disabledSend.value) return;
  //禁用发送按钮
  disabledSend.value = true;
  const res = await sendEmailCodeAPI(registerForm);
  switch (res.data.code) {
    case statusCode.OK:
      //开启倒计时
      startCountdown();
      message.success('发送成功');
      break;
    case statusCode.CAPTCHA_REQUIRED:
      captchaTrigger = "code";
      registerForm.captchaId = res.data.data.captchaId;
      showCaptcha.value = true;
      disabledSend.value = false;
      break;
    case statusCode.FAIL:
      disabledSend.value = false;
      sendBtnText.value = '发送验证码';
      message.error(res.data.msg);
      break;
    default:
      break;
  }
}

</script>

<style lang="scss" scoped>
.bg {
  position: fixed;
  top: 0;
  bottom: 0;
  width: 100%;
  background-color: var(--primary-color);
}

.head-box {
  width: 100%;
  height: 30rem;

  .head-text {
    text-align: left;
    font-size: 1.2rem;
    color: #fff;
    padding: 3rem 0 0 1rem;
    font-weight: bold;
    line-height: 2.2rem;
  }
}

.login-div {
  width: 100%;
  height: 100%;
  position: relative;
  margin-top: -20rem;
  border-radius: 2rem 2rem 0 0;
  background-color: #fff;

  .login-box {
    padding: 2rem 1rem;

    .card-btn {
      text-align: center;

      button {
        width: 100%;
        margin-top: 1rem;
      }
    }
  }
}
</style>
