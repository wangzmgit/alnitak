<template>
  <n-modal style="width: 600px;" v-model:show="modalVisible" preset="card" title="编辑用户信息">
    <n-form ref="formRef" label-placement="left" :label-width="96" :model="formModel" :rules="rules">
      <n-grid :cols="24" :x-gap="18">
        <n-form-item-grid-item :span="24" label="头像" path="avatar">
          <common-avatar :url="formModel.avatar" :size="36"></common-avatar>
          <n-button class="remove-btn" @click="removeAvatar">移除头像</n-button>
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="24" label="背景图" path="spaceCover">
          <div class="space-cover">
            <n-image class="image" v-if="formModel.spaceCover" :src="getResourceUrl(formModel.spaceCover)" />
          </div>
          <n-button class="remove-btn" @click="removeSpaceCover">移除背景图</n-button>
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="24" label="用户名" path="name">
          <n-input v-model:value="formModel.name" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="24" label="邮箱" path="email">
          <n-input v-model:value="formModel.email" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="24" label="个性签名" path="sign">
          <n-input v-model:value="formModel.sign" />
        </n-form-item-grid-item>
      </n-grid>
      <n-space :size="24" justify="end">
        <n-button class="form-btn" @click="closeModal">取消</n-button>
        <n-button class="form-btn" type="primary" @click="handleSubmit">确定</n-button>
      </n-space>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import type { FormInst, FormRules } from 'naive-ui';
import { editUserInfoAPI } from '@/api/user';
import { statusCode } from '@/utils/status-code';
import { getResourceUrl } from "@/utils/resource";
import CommonAvatar from "@/components/common-avatar/index.vue";
import { NGrid, NButton, NSpace, NModal, NForm, NFormItemGridItem, NInput, NImage, useMessage } from 'naive-ui';

const emit = defineEmits(['update:visible', 'refresh']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  editData?: UserInfoType;
}>(), {
  visible: false,
})

const message = useMessage();

const modalVisible = computed({
  get() {
    return props.visible;
  },
  set(visible) {
    emit('update:visible', visible);
  }
});

const formRef = ref<HTMLElement & FormInst>();
const formModel = reactive<UserFormType>({
  uid: 0,
  name: "",
  avatar: "",
  spaceCover: "",
  email: "",
  sign: "",
  role: ""
});

const rules: FormRules = {
  name: { required: true, message: '请输入角色名', trigger: ['blur', 'input'] },
  email: [
    { required: true, message: "请输入邮箱", trigger: ['blur', 'input'] },
    { type: "email", message: "请输入正确的邮箱地址", trigger: ['blur', 'input'] },
  ],
}

const closeModal = () => {
  modalVisible.value = false;
}

const removeAvatar = () => {
  formModel.avatar = "";
}

const removeSpaceCover = () => {
  formModel.spaceCover = "";
}

const handleSubmit = async () => {
  await formRef.value?.validate();

  let success = false;
  const res = await editUserInfoAPI(formModel);
  if (res.data.code === statusCode.OK) {
    success = true;
  } else {
    message.error(res.data.msg);
  }

  if (success) {
    emit('refresh');
    closeModal();
  }
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    Object.assign(formModel, props.editData);
  }
})
</script>

<style lang="scss" scoped>
.form-btn {
  width: 72px;
}

.space-cover {
  width: 160px;
  height: 32px;
  background-color: #cccccc;

  .image {
    width: 160px;
    height: 32px;
  }
}

.remove-btn {
  margin-left: 16px;
}
</style>
