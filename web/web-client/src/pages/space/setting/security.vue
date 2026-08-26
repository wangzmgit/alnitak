<template>
  <div class="security">
    <div class="security-section">
      <h3 class="section-title">修改密码</h3>
      <p class="section-desc">当前为已登录状态，修改后需重新登录</p>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        class="pwd-form"
      >
        <el-form-item label="旧密码" prop="oldPassword">
          <el-input
            v-model="form.oldPassword"
            type="password"
            show-password
            placeholder="请输入当前密码"
          />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="form.newPassword"
            type="password"
            show-password
            placeholder="请输入新密码（至少6位）"
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            show-password
            placeholder="请再次输入新密码"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="submitForm">
            保存
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="security-section">
      <h3 class="section-title">忘记密码？</h3>
      <p class="section-desc">如果忘记了当前密码，请通过邮箱验证重置</p>
      <nuxt-link to="/setpassword" class="reset-link">前往重置密码</nuxt-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import type { FormInstance, FormRules } from "element-plus";
import { changePasswordAPI } from "@/api/auth";
import { statusCode } from "@/utils/status-code";

const formRef = ref<FormInstance | null>(null);
const loading = ref(false);

const form = reactive({
  oldPassword: "",
  newPassword: "",
  confirmPassword: "",
});

const validateConfirm = (_rule: any, value: string, callback: any) => {
  if (value !== form.newPassword) {
    callback(new Error("两次输入的密码不一致"));
  } else {
    callback();
  }
};

const rules: FormRules = {
  oldPassword: [
    { required: true, message: "请输入当前密码", trigger: "blur" },
  ],
  newPassword: [
    { required: true, message: "请输入新密码", trigger: "blur" },
    { min: 6, message: "密码长度不能小于6位", trigger: "blur" },
  ],
  confirmPassword: [
    { required: true, message: "请再次输入新密码", trigger: "blur" },
    { validator: validateConfirm, trigger: "blur" },
  ],
};

const submitForm = async () => {
  const valid = await formRef.value?.validate();
  if (!valid) return;

  loading.value = true;
  try {
    const res = await changePasswordAPI({
      oldPassword: form.oldPassword,
      newPassword: form.newPassword,
    });
    if (res.data.code === statusCode.OK) {
      ElMessage.success(res.data.msg || "密码修改成功");
      form.oldPassword = "";
      form.newPassword = "";
      form.confirmPassword = "";
    } else {
      ElMessage.error(res.data.msg || "修改失败");
    }
  } catch {
    ElMessage.error("网络错误");
  } finally {
    loading.value = false;
  }
};
</script>

<style lang="scss" scoped>
.security {
  padding: 20px;
}

.security-section {
  background: var(--background-primary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 24px;
  margin-bottom: 20px;

  .section-title {
    margin: 0 0 4px 0;
    font-size: 16px;
    font-weight: 600;
    color: var(--font-primary-1);
  }

  .section-desc {
    margin: 0 0 20px 0;
    font-size: 13px;
    color: var(--font-primary-3);
  }
}

.pwd-form {
  max-width: 420px;
}

.reset-link {
  font-size: 14px;
  color: var(--primary-color);
  text-decoration: none;

  &:hover {
    color: var(--primary-hover-color);
  }
}
</style>
