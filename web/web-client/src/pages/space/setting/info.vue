<template>
  <div class="modify-info-from">
    <el-form label-placement="left" label-width="80px">
      <div class="avatar-box">
        <p class="avatar-label">头像:</p>
        <common-avatar class="avatar" :url="userInfo.avatar" :size="50" @click="avatarClick"></common-avatar>
      </div>
      <el-form-item label="UID:">
        <p class="uid form-item">{{ userInfo.uid }}</p>
      </el-form-item>
      <el-form-item label="昵称:">
        <el-input class="form-input name" v-model="userInfo.name" placeholder="请输入昵称" maxlength="20" show-count />
      </el-form-item>
      <el-form-item label="性别:" label-for="gender-radio-group">
        <el-radio-group id="gender-radio-group" class="form-item" v-model="userInfo.gender">
          <el-radio-button label="0">保密</el-radio-button>
          <el-radio-button label="1">男</el-radio-button>
          <el-radio-button label="2">女</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="生日:">
        <el-date-picker class="form-item" date-format="YYYY-MM-DD" :clearable="false" placeholder="选择出生日期" type="date"
          v-model="birthdayData" @change="birthdayChange" />
      </el-form-item>
      <el-form-item label="个性签名:">
        <el-input class="form-input" v-model="userInfo.sign" placeholder="请输入个性签名" maxlength="50" show-word-limit
          type="textarea" :autosize="{ minRows: 3, maxRows: 3 }" resize="none" />
      </el-form-item>
      <div class="setting-save-btn">
        <el-button type="primary" @click="editUserInfo">修改</el-button>
      </div>
    </el-form>
    <image-cropper ref="cropperRef">
      <template #content="fileSlot">
        <avatar-cropper :file="fileSlot.file" @state-change="changeUpload"></avatar-cropper>
      </template>
    </image-cropper>
  </div>
</template>
<script setup lang="ts">
import moment from "moment";
import { ref, onBeforeMount } from "vue";
import { ElMessage } from 'element-plus';
import { getUserInfoAPI, editUserInfoAPI } from '@/api/user';
import CommonAvatar from "@/components/common-avatar/index.vue";
import ImageCropper from "@/components/image-cropper/index.vue";
import AvatarCropper from "@/components/image-cropper/components/AvatarCropper.vue";

const userInfo = ref<UserInfoType>({
  uid: 0,
  name: "",
  avatar: "",
  gender: 0,
  birthday: "",
});
const getUserInfo = async () => {
  const res = await getUserInfoAPI();
  if (res.data.code === statusCode.OK) {
    userInfo.value = res.data.data.userInfo;
    if (userInfo.value.birthday) {
      birthdayData.value = userInfo.value.birthday;
    }
  }
}

// 用户生日信息
const birthdayData = ref("1970-1-1");
const birthdayChange = (val: Date) => {
  userInfo.value.birthday = moment(val).format("YYYY-MM-DD")
}

const cropperRef = ref<InstanceType<typeof ImageCropper> | null>(null);
const avatarClick = () => {
  cropperRef.value?.open();
}

//上传变化的回调
const changeUpload = (status: string, data: any) => {
  switch (status) {
    case "success":
      userInfo.value.avatar = data.data.url;
      break;
    case "error":
      ElMessage.error('文件上传失败');
      break;
  }
  cropperRef.value?.closeCropper();
}

const editUserInfo = async () => {
  const data: EditUserInfoType = {
    avatar: userInfo.value.avatar,
    name: userInfo.value.name,
    sign: userInfo.value.sign || "",
    gender: Number(userInfo.value.gender) || 0,
    birthday: userInfo.value.birthday || '1970-1-1',
    spaceCover: userInfo.value.spaceCover,
  }

  if (!data.name) {
    ElMessage.error("请填写昵称");
    return;
  }

  const res = await editUserInfoAPI(data);
  if (res.data.code === statusCode.OK) {
    await getUserInfo();//获取用户信息
    ElMessage.success("修改成功");
  } else {
    ElMessage.error(res.data.msg || "修改失败");
  }
}

onBeforeMount(async () => {
  await getUserInfo();
})

</script>

<style lang="scss" scoped>
.modify-info-from {
  padding: 0 12px;
}

.avatar-box {
  position: relative;
  height: 50px;
  display: flex;
  align-items: center;
  margin-bottom: 18px;

  .avatar-label {
    box-sizing: border-box;
    color: #606266;
    width: 80px;
    margin: 0;
    padding: 0 12px 0 0;
    text-align: right;
    font-size: 14px;
    line-height: 60px;
  }

  .upload-avatar {
    display: none;
  }

  .avatar {
    position: relative;

    &:hover::after {
      display: block;
    }

    //头像遮罩
    &::after {
      color: #fff;
      font-size: 10px;
      display: none;
      left: 0;
      top: 0;
      position: absolute;
      content: "更换头像";
      line-height: 50px;
      text-align: center;
      vertical-align: middle;
      width: 100%;
      height: 100%;
      border-radius: 50%;
      cursor: pointer;
      background-color: rgba(0, 0, 0, 0.3);
    }
  }
}

.uid {
  margin: 0;
  color: #606266;
}

.setting-save-btn {
  float: right;
  margin-top: 16px;

  button {
    width: 120px;
    height: 40px;
  }
}
</style>
