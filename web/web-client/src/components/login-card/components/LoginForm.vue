<template>
  <div class="login-form">
    <base-tabs :tabs="tabs" @tab-change="tabChange"></base-tabs>
    <form class="login-panel" @submit.prevent="handleLogin">
      <div v-if="currentTab === 'account'" class="input-group">
        <div class="input-box">
          <input v-model="loginForm.email" placeholder="请输入邮箱" class="input account-input" maxlength="64">
        </div>
        <div class="error-text">{{ errorTips.emailError }}</div>
        <div class="input-box">
          <input v-model="loginForm.password" placeholder="请输入密码" type="password" class="input input-password "
            maxlength="64">
          <nuxt-link class="find-password" to="/setpassword" target="blank">忘记密码</nuxt-link>
        </div>
        <div class="error-text">{{ errorTips.passwordError }}</div>
      </div>
      <div v-else class="input-group">
        <div class="input-box">
          <input v-model="loginForm.email" placeholder="请输入邮箱" class="input account-input" maxlength="64">
        </div>
        <div class="error-text">{{ errorTips.emailError }}</div>
        <div class="input-box">
          <input v-model="loginForm.code" placeholder="请输入验证码" class="input code-input" maxlength="6">
          <span class="send-code-btn" :class="disabledSend ? 'btn-disabled' : ''" @click="sendEmailCode">
            {{ sendBtnText }}
          </span>
        </div>
        <div class="error-text">{{ errorTips.codeError }}</div>
      </div>
      <div class="remember-row">
        <label class="remember-label">
          <input type="checkbox" v-model="loginForm.rememberMe" class="remember-checkbox" />
          <span class="remember-text">记住登录</span>
        </label>
      </div>
      <div class="button-group">
        <button class="btn-other" type="button" @click="emit('changeForm')">注册</button>
        <button class="btn-primary" type="submit" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </div>
    </form>
    <client-only>
      <slider-captcha v-model:show="showCaptcha" :captcha-id="loginForm.captchaId"
        @success="captchaSuccess"></slider-captcha>
    </client-only>
  </div>
</template>

<script setup lang="ts">
import { isEmail } from "@/utils/verify";
import { loginAPI, emailLoginAPI } from "@/api/auth";
import type { AxiosResponse } from "axios";
import { sendEmailCodeAPI } from "@/api/code";
import { saveCredentials, broadcastAuthChange } from "@/stores/auth-store";

const emit = defineEmits(["success", "changeForm"]);

const loading = ref(false);
const currentTab = ref('account');
const tabs = [{ key: 'account', label: '密码登录' }, { key: 'code', label: '邮箱登录' }];
const tabChange = (tab: string) => {
  currentTab.value = tab;
}

// 显示滑块验证
let captchaTrigger = "";
const showCaptcha = ref(false);
const captchaSuccess = () => {
  if (captchaTrigger === "login") {
    handleLogin();
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

const loginForm = reactive<UserLoginType>({
  email: '',
  password: '',
  code: '',
  captchaId: '',
  rememberMe: true,
})
const handleLogin = () => {
  initErrorTips();
  if (!loginForm.email) {
    errorTips.emailError = '邮箱不能为空';
    return;
  }
  if (!isEmail(loginForm.email)) {
    errorTips.emailError = '邮箱格式错误';
    return;
  }

  if (currentTab.value === 'account') {
    accountLogin();
  } else {
    codeLogin();
  }
}

// 账号登录
const accountLogin = async () => {
  if (!loginForm.password) {
    errorTips.passwordError = '密码不能为空';
    return;
  }

  if (loginForm.password.length < 6) {
    errorTips.passwordError = '密码长度不能小于6位';
    return;
  }

  loading.value = true;
  try {
    const res = await loginAPI(loginForm);
    handleLoginRes(res);
  } finally {
    loading.value = false;
  }
}

// 验证码登录
const disabledSend = ref(false);//禁用发送按钮
const sendBtnText = ref('发送验证码');//发送按钮文字
const startCountdown = (seconds: number) => {
  let remaining = seconds;
  sendBtnText.value = `${remaining}秒`;
  let tag = setInterval(() => {
    if (--remaining <= 0) {
      clearInterval(tag);
      disabledSend.value = false;
      sendBtnText.value = '发送验证码';
      return;
    }
    sendBtnText.value = `${remaining}秒`;
  }, 1000);
}
const sendEmailCode = async () => {
  if (disabledSend.value) return;
  //禁用发送按钮
  disabledSend.value = true;
  const res = await sendEmailCodeAPI({ email: loginForm.email, captchaId: loginForm.captchaId });
  switch (res.data.code) {
    case statusCode.OK:
      //开启倒计时，使用后端返回的冷却时间
      startCountdown(res.data.data.countdown || 60);
      ElMessage.success(res.data.msg || '发送成功');
      break;
    case statusCode.CAPTCHA_REQUIRED:
      captchaTrigger = "code";
      loginForm.captchaId = res.data.data.captchaId;
      showCaptcha.value = true;
      disabledSend.value = false;
      break;
    case statusCode.FAIL:
      //如果后端返回了冷却时间（发送过于频繁），开启倒计时
      if (res.data.data?.countdown) {
        startCountdown(res.data.data.countdown);
      } else {
        disabledSend.value = false;
        sendBtnText.value = '发送验证码';
      }
      ElMessage.error(res.data.msg);
      break;
    default:
      break;
  }
}

const codeLogin = async () => {
  if (!loginForm.code) {
    errorTips.codeError = '验证码不能为空';
    return;
  }

  if (loginForm.code.length !== 6) {
    errorTips.codeError = '验证码长度为6位';
    return;
  }

  loading.value = true;
  try {
    const res = await emailLoginAPI(loginForm);
    handleLoginRes(res);
  } finally {
    loading.value = false;
  }
}

const handleLoginRes = async (res: AxiosResponse<any, any>) => {
  switch (res.data.code) {
    case statusCode.CAPTCHA_REQUIRED:
      captchaTrigger = "login";
      loginForm.captchaId = res.data.data.captchaId;
      showCaptcha.value = true;
      break;
    case statusCode.OK:
      saveCredentials(res.data.data);
      broadcastAuthChange();
      emit("success");
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
  margin: 40px 0 0 50px;

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

    .find-password {
      cursor: pointer;
      font-size: 14px;
      position: absolute;
      right: 0;
      top: 0;
      font-weight: 500;
      padding: 10px;
      line-height: 17px;
      color: var(--font-primary-2);
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
  margin-top: 20px;
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
    color: var(--font-primary-1);
    background: var(--bg-elev-1);
    border: 1px solid var(--border-color);
    border-radius: 8px;
  }

  .btn-primary {
    box-sizing: border-box;
    width: calc(50% - 10px);
    height: 40px;
    cursor: pointer;
    color: var(--primary-text-color);
    background: var(--primary-color);
    border-radius: 8px;
    border: 1px solid var(--primary-color);

    &:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
  }
}

.btn-disabled {
  color: var(--font-primary-3) !important;
  cursor: not-allowed !important;

  &:hover {
    color: var(--font-primary-3) !important;
  }
}

.remember-row {
  display: flex;
  align-items: center;
  margin-top: 12px;
  padding: 0 10px;

  .remember-label {
    display: flex;
    align-items: center;
    cursor: pointer;
    user-select: none;
  }

  .remember-checkbox {
    width: 15px;
    height: 15px;
    accent-color: var(--primary-color);
    cursor: pointer;
  }

  .remember-text {
    margin-left: 6px;
    font-size: 13px;
    color: var(--font-primary-3);
  }
}
</style>