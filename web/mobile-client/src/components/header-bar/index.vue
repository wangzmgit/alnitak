<template>
  <div class="header-box">
    <common-avatar :url="userInfo?.avatar" :size="36" @click="headerRouter('Space')" />
    <div class="search-box">
      <n-input v-model:value="keywords" round placeholder="搜索关键词~" @keydown.enter="search">
        <template #suffix>
          <n-icon @click="search">
            <search-icon></search-icon>
          </n-icon>
        </template>
      </n-input>
    </div>
    <n-icon class="msg" @click="headerRouter('Message')">
      <message-icon></message-icon>
    </n-icon>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount } from "vue";
import { useRouter, useRoute } from "vue-router";
import { NIcon, NInput } from "naive-ui";
import { statusCode } from "@/utils/status-code";
import SearchIcon from "@/components/icons/SearchIcon.vue";
import MessageIcon from "@/components/icons/MessageIcon.vue";
import CommonAvatar from "@/components/common-avatar/index.vue";
import { getUserInfoAPI } from "@/api/user";

const route = useRoute();
const router = useRouter();

const userInfo = ref<UserInfoType>();
const getUserBaseInfo = async () => {
  const res = await getUserInfoAPI();
  if (res.data.code === statusCode.OK) {
    userInfo.value = res.data.data.userInfo;
  }
}

const headerRouter = (page: string) => {
  const path = route.name;
  switch (page) {
    case 'Space':
      if (path == "SpaceInfo")
        page = "Home";
      break;
    case 'Message':
      if (path == "Message")
        page = "Home";
      break;
  }
  router.push({ name: page });
}

// 搜索
const keywords = ref("");
const search = () => {
  router.push({ name: "Search", query: { keywords: keywords.value } });
}

onBeforeMount(() => {
  getUserBaseInfo();
})
</script>

<style lang="scss" scoped>
.header-box {
  width: calc(100% - 20px);
  height: 50px;
  display: flex;
  align-items: center;
  padding: 0 10px;
  background-color: var(--primary-color);
  justify-content: space-between;
  -webkit-box-shadow: 0px 0px 3px #c8c8c8;
  -moz-box-shadow: 0px 0px 3px #c8c8c8;
  box-shadow: 0px 0px 3px #c8c8c8;

  .search-box {
    width: calc(100% - 130px);
  }

  .msg {
    color: #fff;
    font-size: 26px;
  }
}
</style>