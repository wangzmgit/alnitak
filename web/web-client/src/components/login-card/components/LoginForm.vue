<template>
  <div class="login-form">
    <base-tabs :tabs="tabs" @tab-change="tabChange"></base-tabs>
    <div class="login-panel">
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
      <div class="button-group">
        <button class="btn-other" @click="emit('changeForm')">注册</button>
        <button class="btn-primary" @click="handelLogin">登录</button>
      </div>
    </div>
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
import BaseTabs from "@/components/base-tabs/index.vue";
import { sendEmailCodeAPI } from "@/api/code";

const emit = defineEmits(["success", "changeForm"]);

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
    handelLogin();
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
  captchaId: ''
})
const handelLogin = () => {
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

  const res = await loginAPI(loginForm);
  handelLoginRes(res);
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
  const res = await sendEmailCodeAPI(loginForm);
  switch (res.data.code) {
    case statusCode.OK:
      //开启倒计时
      startCountdown();
      ElMessage.success('发送成功');
      break;
    case statusCode.CAPTCHA_REQUIRED:
      captchaTrigger = "code";
      loginForm.captchaId = res.data.data.captchaId;
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

const codeLogin = async () => {
  if (!loginForm.code) {
    errorTips.codeError = '验证码不能为空';
    return;
  }

  if (loginForm.code.length !== 6) {
    errorTips.codeError = '验证码长度为6位';
    return;
  }

  const res = await emailLoginAPI(loginForm);
  handelLoginRes(res);
}

const handelLoginRes = async (res: AxiosResponse<any, any>) => {
  switch (res.data.code) {
    case statusCode.CAPTCHA_REQUIRED:
      captchaTrigger = "login";
      loginForm.captchaId = res.data.data.captchaId;
      showCaptcha.value = true;
      break;
    case statusCode.OK:
      storageData.set("token", res.data.data.token, 60);
      storageData.set("refreshToken", res.data.data.refreshToken, 7 * 24 * 60);

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