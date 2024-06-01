<template>
  <n-dropdown trigger="hover" :options="options" @select="handleDropdownClick">
    <n-avatar v-if="userInfo?.avatar" size="medium" :src="getResourceUrl(userInfo.avatar)"></n-avatar>
    <n-avatar v-else size="medium"></n-avatar>
  </n-dropdown>
</template>

<script lang="ts" setup>
import { h } from "vue";
import { useRouter } from "vue-router";
import { NText } from "naive-ui";
import { storeToRefs } from "pinia";
import { NDropdown, NAvatar } from "naive-ui";
import useUserStore from "@/stores/modules/user-store.js";
import { getResourceUrl } from "@/utils/resource";
import { renderIcon } from "@/utils/render";
import { LogInOutline } from "@vicons/ionicons5";


const router = useRouter();

const loginStore = useUserStore();
const { userInfo } = storeToRefs(loginStore);


const renderCustomDropdownHeader = () => {
  return h(
    "div",
    {
      style: "max-width: 200px; display: flex; align-items: center; padding: 8px 12px;",
    },
    [
      h("div", null, [
        h("div", null, [
          h(NText, { depth: 2 }, { default: () => userInfo.value?.name }),
        ]),
        h("div", { style: "font-size: 12px; margin-top: 8px;" }, [
          h(NText, { depth: 3 }, { default: () => userInfo.value?.sign }),
        ]),
      ]),
    ]
  );
};

const options = [
  {
    key: "header",
    type: "render",
    render: renderCustomDropdownHeader,
  },
  {
    key: "header-divider",
    type: "divider",
  },
  {
    label: "退出账号",
    key: "loginOut",
    icon: renderIcon(LogInOutline),
  },
];

const handleDropdownClick = (key: string) => {
  switch (key) {
    case "userInfo":
      router.push({ name: "userInfo" });
      break;
    case "loginOut":
      loginStore.logout();
      break;
    default:
      console.error(`u-dropdown has no key: ${key}`);
      break;
  }
};
</script>
