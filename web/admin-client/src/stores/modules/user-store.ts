import { ref } from "vue";
import router from "@/router";
import { defineStore } from "pinia";
import { statusCode } from "@/utils/status-code";
import { getUserInfoAPI  } from "@/api/user";
import { getRoleInfoAPI } from "@/api/role";
import { logoutAPI } from "@/api/auth";
import { storageData } from "@/utils/storage-data";

const useUserStore = defineStore("login", () => {
  const token = ref("");
  const userInfo = ref<UserInfoType>();
  const homePage = ref('');

  token.value = localStorage.getItem('refreshToken') || ""

  /** 登录 */
  const login = async () => {
    const { redirect = "/" } = router.currentRoute.value.query;
    // 获取用户信息
    await getUserInfo();
    window.$message.success("登陆成功，正在跳转...");
    const redirectUrl = decodeURIComponent(redirect as string);
    router.push(redirectUrl);
  }

  /** 登出 */
  const logout = () => {
    window.$dialog.error({
      title: "警告",
      content: "确认退出账号吗？",
      negativeText: "取消",
      positiveText: "确认",
      onPositiveClick: async () => {
        await logoutAPI(storageData.get('refreshToken'));
        storageData.remove("token");
        storageData.remove('refreshToken');
        
        router.push({ name: "login" });
        window.$message.success("退出账号成功");
      },
    });
  }

  /** 获取用户信息 */
  const getUserInfo = async () => {
    const res = await getUserInfoAPI();
    if (res.data.code === statusCode.OK) {
      userInfo.value = res.data.data.userInfo;
    }

    const roleRes = await getRoleInfoAPI();
    if (roleRes.data.code === statusCode.OK) {
      homePage.value = roleRes.data.data.role.homePage;
    }
  }

  return {
    token,
    homePage,
    userInfo,
    login,
    logout,
    getUserInfo
  }
})

export default useUserStore;
