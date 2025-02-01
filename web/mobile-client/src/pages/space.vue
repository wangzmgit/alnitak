<template>
  <div class="space">
    <header-bar class="header-bar "></header-bar>
    <div class="space-container">
      <div class="space-header">
        <button class="upload-btn" @click="uploadCoverClick">上传封面图片</button>
        <img class="cover" v-if="userInfo?.spaceCover" :src="getResourceUrl(userInfo.spaceCover)" alt="用户封面图" />
        <div class="header-inner">
          <common-avatar :url="userInfo?.avatar" :size="60" :iconSize="36"></common-avatar>
          <div class="header-info">
            <div class="name">
              <span>{{ userInfo?.name }}</span>
              <span class="gender-icon">
                <male v-if="userInfo?.gender === 1" fill="#1890ff"></male>
                <female v-else-if="userInfo?.gender === 2" fill="#eb2f96"></female>
              </span>
            </div>
            <div class="sign">{{ userInfo?.sign }}</div>
          </div>
          <!--关注粉丝信息-->
          <div class="user-data">
            <div>
              <p class="data-title">稿件</p>
              <p class="data-content">{{ videoCount }}</p>
            </div>
            <div>
              <p class="data-title">关注</p>
              <nuxt-link class="data-content" to="/space/following">{{ userData.followingCount }}</nuxt-link>
            </div>
            <div>
              <p class="data-title">粉丝</p>
              <nuxt-link class="data-content" to="/space/follower">{{ userData.followerCount }}</nuxt-link>
            </div>
          </div>
        </div>
      </div>
      <div class="space-content">
        <div class="space-menu">
          <nuxt-link v-for="item in menuList" :to="item.to" class="menu-item" :target="item.blank ? '_blank' : '_self'"
            :class="getMenuName(route.name) === item.key ? 'menu-item-active' : ''">
            <span class="menu-icon">
              <component :is="item.icon" size="18"></component>
            </span>
            <span class="menu-title">{{ item.name }}</span>
          </nuxt-link>
        </div>
        <div class="space-router">
          <NuxtPage />
        </div>
      </div>
    </div>
    <image-cropper ref="cropperRef">
      <template #content="fileSlot">
        <space-cover-cropper :file="fileSlot.file" @state-change="changeUpload"></space-cover-cropper>
      </template>
    </image-cropper>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia';
import HeaderBar from '@/components/header-bar/index.vue';
import CommonAvatar from '@/components/common-avatar/index.vue';
import { useVideoCountStore } from "@/composables/video-count-store";
import { getUserInfoAPI, asyncGetUserBaseInfoAPI, editUserInfoAPI } from '@/api/user';
import type { RouteRecordName } from 'vue-router';
import { getFollowDataAPI } from '@/api/relation';
import UploadIcon from "@/components/icons/UploadIcon.vue";
import ImageCropper from "@/components/image-cropper/index.vue";
import SpaceCoverCropper from "@/components/image-cropper/components/SpaceCoverCropper.vue";
import { Male, Female, VideoTwo, FolderFocus, History as HistoryIcon, Message as MessageIcon, Config } from '@icon-park/vue-next';

definePageMeta({
  middleware: ['auth', (to) => {
    if (to.name === 'space') {
      return navigateTo(`/space/manuscript`);
    }
  }]
})

const route = useRoute();
const menuList = [
  {
    key: 'space-manuscript',
    name: '稿件',
    to: '/space/manuscript',
    icon: VideoTwo,
  },
  {
    key: 'space-collection',
    name: '收藏夹',
    to: '/space/collection',
    icon: FolderFocus,
  },
  {
    key: 'space-history',
    name: '历史记录',
    to: '/space/history',
    icon: HistoryIcon,
  },
  {
    key: 'upload-video',
    name: '创作中心',
    to: '/upload/video',
    blank: true,
    icon: UploadIcon,
  },
  {
    key: 'space-message',
    name: '消息',
    to: '/message/announce',
    blank: true,
    icon: MessageIcon,
  },
  {
    key: 'space-setting',
    name: '设置',
    to: '/space/setting',
    icon: Config,
  },
]

const videoCountStore = useVideoCountStore();
const { videoCount } = storeToRefs(videoCountStore);

const userId = useCookie('user_id');
const userInfo = ref<UserInfoType>();
const { data } = await asyncGetUserBaseInfoAPI(userId.value!);
if ((data.value as any).code === statusCode.OK) {
  userInfo.value = (data.value as any).data.userInfo;
}
useHead({
  title: `${userInfo.value?.name}的个人中心`
})

const getMenuName = (name: RouteRecordName | null | undefined) => {
  if (name === "space-setting-info" || name === "space-setting-security") {
    return "space-setting"
  }

  return name;
}

// 上传封面图
const cropperRef = ref<InstanceType<typeof ImageCropper> | null>(null);
const uploadCoverClick = () => {
  cropperRef.value?.open();
}

//上传变化的回调
const changeUpload = async (status: string, data: any) => {
  switch (status) {
    case "success":
      //更新封面图
      if (!userInfo.value) return;
      userInfo.value.spaceCover = data.data.url;
      const res = await editUserInfoAPI(userInfo.value);
      if (res.data.code === statusCode.OK) {
        await getUserInfo();//获取用户信息
      } else {
        ElMessage.error(res.data.msg || "修改失败");
      }
      break;
    case "error":
      ElMessage.error('封面上传失败');
      break;
  }
  cropperRef.value?.closeCropper();
}

const getUserInfo = async () => {
  const res = await getUserInfoAPI();
  if (res.data.code === statusCode.OK) {
    userInfo.value = res.data.data.userInfo;
  }
}

onMounted(async () => {
  // 处理跳转页面回来用户信息为空的问题
  if (!userInfo.value || userInfo.value.uid === 0) {
    await getUserInfo();

    document.title = `${userInfo.value?.name}的个人中心`;
  }

  if (userInfo.value) {
    getFollowData(userInfo.value.uid);
  }
})

const userData = reactive({
  loading: true,
  followingCount: 0,
  followerCount: 0
})

//获取关注数和粉丝数
const getFollowData = async (id: number | string) => {
  const res = await getFollowDataAPI(id);
  if (res.data.code === statusCode.OK) {
    userData.followerCount = res.data.data.follower;
    userData.followingCount = res.data.data.following;
  }

  userData.loading = false;
}
</script>

<style lang="scss" scoped>
.space {
  width: 100%;

  .header-bar {
    position: fixed;
    top: 0;
  }

  .space-container {
    width: 1300px;
    margin: 70px auto 0;
  }
}

.space-header {
  position: relative;
  width: 100%;
  height: 220px;
  background-color: #9196a1;

  .upload-btn {
    position: absolute;
    top: 12px;
    right: 20px;
    width: 120px;
    height: 32px;
    color: #fff;
    background-color: rgba(0, 0, 0, .45);
    cursor: pointer;
    border-radius: 6px;
    border: none;

    &:hover {
      background-color: rgba(0, 0, 0, .5);
    }
  }

  .cover {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .header-inner {
    width: 100%;
    height: 100px;
    position: absolute;
    bottom: 0;
    display: flex;
    align-items: center;
    box-sizing: border-box;
    padding: 0 20px;

    .header-info {
      flex: 1;
      padding-left: 20px;

      .name {
        color: var(--font-primary-6);
        display: inline-block;
        margin-right: 5px;
        font-weight: 600;
        line-height: 18px;
        font-size: 18px;
        vertical-align: middle;
      }

      .gender-icon {
        display: inline-block;
        width: 18px;
        height: 18px;
        margin: 0 6px;
        vertical-align: middle;
      }

      .sign {
        font-size: 12px;
        color: #ccd0d7;
        margin: 6px 0 0 0;
      }
    }

    .user-data {
      width: 236px;
      display: flex;

      div {
        color: #fff;
        width: 78px;
        text-align: center;

        .data-title {
          margin-bottom: 6px;
        }

        .data-content {
          color: #fff;
          margin: 12px 0;
          display: block;
          cursor: pointer;
        }
      }
    }
  }
}

.space-content {
  display: flex;
  margin: 0 auto 30px;

  .space-menu {
    width: 200px;
    background-color: #f8f8f8;

    .menu-item {
      position: relative;
      height: 42px;
      padding-left: 24px;
      cursor: pointer;
      border-radius: 2px;
      margin: 6px 5px 0;
      box-sizing: border-box;
      display: flex;
      align-items: center;
      text-decoration: none;
      color: var(--font-primary-1);

      &:hover {
        background-color: #f3f3f5;
      }

      .menu-icon {
        width: 18px;
        height: 18px;
        margin-right: 8px;
      }
    }

    .menu-item-active {
      color: var(--primary-color) !important;
      background-color: var(--primary-color-active) !important;
    }
  }

  .space-router {
    margin-left: 10px;
    width: calc(100% - 210px);
    min-height: 630px;
    background-color: #fff;
  }
}

@media (max-width: 1400px) {
  .space {
    .space-container {
      width: 1100px;
    }
  }

  .space-header {
    height: 260px;
  }
}
</style>