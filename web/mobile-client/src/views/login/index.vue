<template>
	<div class="bg">
		<div class="head-box">
			<div class="head-text">
				<span>您好，</span>
				<br>
				<span>欢迎回到 {{ globalConfig.title }}</span>
			</div>
		</div>
		<div class="login-div">
			<div class="login-box">
				<n-form ref="formRef" :rules="rules" :model="loginForm">
					<n-form-item label="邮箱" path="email">
						<n-input placeholder="请输入邮箱" v-model:value="loginForm.email" />
					</n-form-item>
					<n-form-item label="密码" path="password">
						<n-input placeholder="请输入密码" v-model:value="loginForm.password" type="password" />
					</n-form-item>
					<div class="card-btn">
						<n-button type="primary" @click="handelLogin()">登录</n-button>
						<n-button @click="goRegister">注册</n-button>
					</div>
				</n-form>
			</div>
		</div>
		<slider-captcha v-model:show="showCaptcha" :captcha-id="loginForm.captchaId"
			@success="captchaSuccess"></slider-captcha>
	</div>
</template>
<script setup lang="ts">
import router from "@/router";
import { reactive, ref } from 'vue';
import type { FormInst, FormRules } from 'naive-ui';
import { globalConfig } from "@/utils/global-config";
import { loginAPI, emailLoginAPI } from "@/api/auth";
import type { AxiosResponse } from "axios";
import { statusCode } from "@/utils/status-code";
import { storageData } from "@/utils/storage-data";
import SliderCaptcha from "@/components/slider-captcha/index.vue";
import { NForm, NFormItem, NButton, NInput, useMessage } from 'naive-ui';

const message = useMessage();//通知
const loginForm = reactive<UserLoginType>({
	email: '',
	password: '',
	code: '',
	captchaId: ''
})

const rules: FormRules = {
	email: [
		{ required: true, message: "请输入邮箱", trigger: ['blur', 'input'] },
		{ type: "email", message: "请输入正确的邮箱地址", trigger: ['blur', 'input'] },
	],
	password: { required: true, message: '请输入密码', trigger: ['blur', 'input'] },
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

const currentTab = ref('account');
const formRef = ref<FormInst | null>(null);
const handelLogin = async () => {
	await formRef.value?.validate()

	if (currentTab.value === 'account') {
		accountLogin();
	} else {
		codeLogin();
	}
}

// 账号登录
const accountLogin = async () => {
	const res = await loginAPI(loginForm);
	handelLoginRes(res);
}

const sendEmailCode = async () => {

}

const codeLogin = async () => {
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
			router.push({ name: 'Home' });
			break;
		default:
			message.error(res.data.msg);
	}
}

const goRegister = () => {
	router.push({ name: 'Register' });
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
