<template>
  <div class="login-form">
    <base-tabs :tabs="tabs" @tab-change="tabChange"></base-tabs>
    <div class="login-panel">
      <div v-if="currentTab === 'email'" class="input-group">
        <div class="input-box">
          <input v-model="registerForm.email" placeholder="请输入邮箱" class="input account-input" maxlength="64">
        </div>
        <div class="error-text">{{ errorTips.emailError }}</div>
        <div class="input-box">
          <input v-model="registerForm.password" placeholder="请输入密码" type="password" class="input input-password "
            maxlength="64">
        </div>
        <div class="error-text">{{ errorTips.passwordError }}</div>
        <div class="input-box">
          <input v-model="registerForm.code" placeholder="请输入验证码" class="input code-input" maxlength="6">
          <span class="send-code-btn" :class="disabledSend ? 'btn-disabled' : ''" @click="sendEmailCode">
            {{ sendBtnText }}
          </span>
        </div>
        <div class="error-text">{{ errorTips.codeError }}</div>
      </div>
      <div class="button-group">
        <button class="btn-other" @click="emit('changeForm')">登录</button>
        <button class="btn-primary" @click="handelRegister">注册</button>
      </div>
    </div>
    <client-only>
      <slider-captcha v-model:show="showCaptcha" :captcha-id="registerForm.captchaId"
        @success="captchaSuccess"></slider-captcha>
    </client-only>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { isEmail } from "@/utils/verify";
import { statusCode } from '@/utils/status-code';
import { registerAPI } from "@/api/auth";
import BaseTabs from "@/components/base-tabs/index.vue";
import SliderCaptcha from "@/components/slider-captcha/index.vue";
import { sendEmailCodeAPI } from '~/api/code';

const emit = defineEmits(["changeForm"]);

const currentTab = ref('email');
const tabs = [{ key: 'email', label: '邮箱注册' }, { key: 'phone', label: '手机注册', disabled: true }];
const tabChange = (tab: string) => {
  currentTab.value = tab;
}

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

const errorTips = reactive<Record<string, string>>({
  emailError: '',
  passwordError: '',
  codeError: '',
});

const initErrorTips = () => {
  for (const key in errorTips) {
    errorTips[key] = '';
  }
}

// 注册表单
const registerForm = reactive<UserRegisterType>({
  email: "",
  password: "",
  code: "",
  captchaId: ""
});

const handelRegister = () => {
  initErrorTips();
  if (!registerForm.email) {
    errorTips.emailError = '邮箱不能为空';
    return;
  }
  if (!isEmail(registerForm.email)) {
    errorTips.emailError = '邮箱格式错误';
    return;
  }

  if (!registerForm.code) {
    errorTips.codeError = '验证码不能为空';
    return;
  }

  if (registerForm.code.length !== 6) {
    errorTips.codeError = '验证码长度为6位';
    return;
  }

  if (currentTab.value === 'email') {
    emailRegister();
  }
}

// 验证码登录
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
      ElMessage.success('发送成功');
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
      ElMessage.error(res.data.msg);
      break;
    default:
      break;
  }
}


const emailRegister = async () => {
  const res = await registerAPI(registerForm);
  switch (res.data.code) {
    case statusCode.CAPTCHA_REQUIRED:
      captchaTrigger = "register";
      registerForm.captchaId = res.data.data.captchaId;
      showCaptcha.value = true;
      break;
    case statusCode.OK:
      ElMessage.success("注册成功");
      emit("changeForm");
      break;
    default:
      errorTips.emailError = res.data.msg;
  }
}
</script>

<style lang="scss" scoped>
.login-form {
  width: 100%;
  height: 100%;
}

.login-panel {
  width: calc(100% - 70px);
  height: calc(100% - 80px);
  margin: 22px 0 0 50px;

  .input-box {
    margin-top: 3px;
    position: relative;
    display: flex;
    background-color: var(--input-bg-color);
    border-radius: 4px;

    .input {
      color: var(--font-primary-1);
      padding: 10px;
      width: 100%;
      border: 1px solid transparent;
      border-radius: 4px;
      outline: none;
      box-sizing: border-box;
      background-color: transparent;

      &:focus {
        border-color: var(--primary-color);
        background-color: var(--input-focus-bg-color);
      }
    }

    .input-password {
      padding-right: 6em;
    }

    .send-code-btn {
      width: 76px;
      text-align: center;
      cursor: pointer;
      font-size: 14px;
      position: absolute;
      right: 0;
      top: 0;
      font-weight: 500;
      padding: 10px;
      line-height: 17px;
      color: var(--primary-color);
    }
  }

  .error-text {
    color: var(--font-danger);
    font-size: 12px;
    line-height: 20px;
    min-height: 20px;
    text-align: left;
    padding: 0 10px;
  }
}

.button-group {
  display: flex;
  justify-content: space-between;
  margin-top: 6px;
  width: 100%;
  font-style: normal;
  font-weight: 400;
  font-size: 14px;
  line-height: 40px;
  text-align: center;

  .btn-other {
    box-sizing: border-box;
    width: calc(50% - 10px);
    height: 40px;
    cursor: pointer;
    color: #18191c;
    background: #fff;
    border: 1px solid #e3e5e7;
    border-radius: 8px;
  }

  .btn-primary {
    box-sizing: border-box;
    width: calc(50% - 10px);
    height: 40px;
    cursor: pointer;
    color: #fff;
    background: var(--primary-color);
    border-radius: 8px;
    border: 1px solid #e3e5e7;
  }
}

.btn-disabled {
  color: #9499a0 !important;
  cursor: not-allowed !important;

  &:hover {
    color: #9499a0 !important;
  }
}
</style>