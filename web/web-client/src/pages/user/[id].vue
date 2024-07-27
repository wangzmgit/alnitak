<template>
  <div class="space">
    <header-bar class="header-bar"></header-bar>
    <div class="space-container">
      <div class="space-header">
        <img class="cover" v-if="userInfo?.spaceCover" :src="getResourceUrl(userInfo.spaceCover)" alt="用户封面图" />
        <div class="header-inner">
          <common-avatar :url="userInfo?.avatar" :size="60" :iconsize="36"></common-avatar>
          <div class="header-info">
            <div class="name">
              <span>{{ userInfo?.name }}</span>
              <span class="gender-icon">
                <male v-if="userInfo?.gender === 1" fill="#1890ff"></male>
                <female v-else-if="userInfo?.gender === 2" fill="#eb2f96"></female>
              </span>
            </div>
            <p class="sign">{{ userInfo?.sign }}</p>
          </div>
          <!--关注粉丝信息-->
          <div class="card-btn">
            <button class="btn" @click="goWhisperPage">私信</button>
            <button class="btn" @click="followBtnClick">{{ btnText }}</button>
          </div>
        </div>
      </div>
      <div class="user-menu">
        <ul class="menu-list">
          <nuxt-link class="menu-item" :class="route.name === 'user-id-manuscript' ? 'menu-active' : ''"
            :to="`/user/${route.params.id}/manuscript`">投稿</nuxt-link>
          <nuxt-link class="menu-item" :class="route.name === 'user-id-following' ? 'menu-active' : ''"
            :to="`/user/${route.params.id}/following`">关注</nuxt-link>
          <nuxt-link class="menu-item" :class="route.name === 'user-id-follower' ? 'menu-active' : ''"
            :to="`/user/${route.params.id}/follower`">粉丝</nuxt-link>
        </ul>
      </div>
      <div class="user-content">
        <NuxtPage />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { asyncGetUserBaseInfoAPI } from '@/api/user';
import { followAPI, getFollowDataAPI, getUserRelationAPI, unfollowAPI } from '@/api/relation';
import { Male, Female } from '@icon-park/vue-next';

definePageMeta({
  middleware: [(to) => {
    if (to.name === 'user-id') {
      return navigateTo({ name: "user-id-manuscript", params: { id: to.params.id } });
    }
  }]
})

const route = useRoute();

const userId = parseInt(route.params.id.toString())
const userInfo = ref<UserInfoType>();
const { data } = await asyncGetUserBaseInfoAPI(userId);
if ((data.value as any).code === statusCode.OK) {
  userInfo.value = (data.value as any).data.userInfo;
}

useHead({
  title: `${userInfo.value?.name}的个人中心`
})

const userData = reactive({
  loading: true,
  followingCount: 0,
  followerCount: 0
})

const relation = ref(relationCode.NOT_FOLLOWING);
const getUserRelation = async () => {
  const res = await getUserRelationAPI(userId);
  if (res.data.code === statusCode.OK) {
    relation.value = res.data.data.relation;
  }
}

const btnText = computed(() => {
  switch (relation.value) {
    case relationCode.FOLLOWED:
      return "已关注";
    case relationCode.NOT_FOLLOWING:
      return "关注";
    case relationCode.MUTUAL_FANS:
      return "已互粉";
  }
})

const followBtnClick = async () => {
  if (!storageData.get("refreshToken")) {
    return navigateTo(`/login?redirect=${location.pathname}`);
  }
  const reqFunc = relation.value === relationCode.NOT_FOLLOWING ? followAPI : unfollowAPI;
  const res = await reqFunc(userId);
  if (res.data.code === statusCode.OK) {
    getUserRelation();
  } else {
    ElMessage.error(res.data.msg);
  }
}

//获取关注数和粉丝数
const getFollowData = async (id: number | string) => {
  const res = await getFollowDataAPI(id);
  if (res.data.code === statusCode.OK) {
    userData.followerCount = res.data.data.follower;
    userData.followingCount = res.data.data.following;
  }

  userData.loading = false;
}

const goWhisperPage = () => {
  navigateTo({ name: 'message-whisper', query: { target: userId } });
}


onBeforeMount(() => {
  getUserRelation()
  // getFollowData(route.params.id.toString());
})
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

    .card-btn {
      width: 160px;
      display: flex;
      margin-right: 20px;
      justify-content: space-between;

      .btn {
        width: 70px;
        box-sizing: border-box;
        padding: 0;
        line-height: 30px;
        height: 30px;
        border-radius: 6px;
        margin-top: 5px;
        text-align: center;
        color: #fff;
        font-size: 14px;
        cursor: pointer;
        border: 1px solid var(--primary-color);
        background: var(--primary-color);

        &:hover {
          background: var(--primary-hover-color);
        }
      }
    }
  }
}

.user-menu {
  height: 46px;
  padding: 0 20px;
  box-shadow: 0 0 0 1px #eee;

  .menu-list {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    align-items: center;
    height: 100%;

    .menu-item {
      color: #222;
      box-sizing: border-box;
      height: 44px;
      line-height: 44px;
      padding: 0 4px;
      margin-right: 20px;
      font-size: 14px;
      cursor: pointer;
    }
  }
}

.menu-active {
  color: var(--primary-color) !important;
  border-bottom: 3px solid;
}

.user-content {
  margin-top: 10px;
  min-height: 300px;
  background-color: #fff;
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