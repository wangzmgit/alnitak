<template>
  <ul class="follow-list">
    <li class="follow-card" v-for="(item, index) in followList" :key="index">
      <!--头像-->
      <div class="follow-avatar">
        <common-avatar :url="item.user.avatar" :size="50"></common-avatar>
      </div>
      <!--昵称和个签-->
      <div class="info-box">
        <nuxt-link class="follow-name" :to="`/user/${item.user.uid}`" target="_blank">{{ item.user.name }}</nuxt-link>
        <span class="follow-sign">{{ item.user.sign }}</span>
      </div>
      <div v-if="showBtn" class="follow-btn" @click="followBtnClick(item)">{{ getBtnText(item.relation) }}</div>
    </li>
  </ul>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed, onMounted, onUnmounted } from "vue";
import { getFollowingAPI, getFollowersAPI, followAPI, unfollowAPI, getUserRelationAPI } from '@/api/relation';
import { requireLogin } from "@/utils/require-login";
import { useAuthStore } from "@/stores/auth-store";

const props = withDefaults(defineProps<{
  userId: number,
  following?: boolean,
  showBtn?: boolean;
}>(), {
  userId: 0,
  following: false,
  showBtn: false,
})

const getBtnText = (relation: number) => {
  switch (relation) {
    case relationCode.FOLLOWED:
      return "已关注";
    case relationCode.NOT_FOLLOWING:
      return "关注";
    case relationCode.MUTUAL_FANS:
      return "已互粉";
  }
}

const pageInfo = reactive({
  current: 1,
  pageSize: 9
});

const noMore = ref(false);//没有更多
const loading = ref(false);//加载中
const followList = ref<RelationListType[]>([]);

const getRelationList = async () => {
  if (!props.userId) return;
  loading.value = true;
  const reqFunc = props.following ? getFollowingAPI : getFollowersAPI;
  const res = await reqFunc(props.userId, pageInfo.current, pageInfo.pageSize);
  if (res.data.data.users) {
    followList.value = followList.value.concat(res.data.data.users);
  } else {
    noMore.value = true;
  }
  loading.value = false;
}

const scrollLoad = () => {
  if (!noMore.value && !loading.value) {
    pageInfo.current++;
    getRelationList();
  }
}

const onWindowScroll = () => {
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop
  const scrollHeight = document.documentElement.scrollHeight || document.body.scrollHeight
  const clientHeight = document.documentElement.clientHeight
  if (scrollHeight - scrollTop - clientHeight < 100) scrollLoad()
}
onMounted(() => window.addEventListener('scroll', onWindowScroll))
onUnmounted(() => window.removeEventListener('scroll', onWindowScroll))

const followBtnClick = async (relation: RelationListType) => {
  if (!(await requireLogin('关注'))) return;
  const reqFunc = relation.relation === relationCode.NOT_FOLLOWING ? followAPI : unfollowAPI;
  const res = await reqFunc(relation.user.uid);
  if (res.data.code === statusCode.OK) {
    const resRelation = await getUserRelationAPI(relation.user.uid);
    if (resRelation.data.code === statusCode.OK) {
      relation.relation = resRelation.data.data.relation;
    }
  } else {
    ElMessage.error(res.data.msg);
  }
}

const auth = useAuthStore();
const viewerLoggedIn = computed(() => auth.isLoggedIn);

const resetAndReload = () => {
  pageInfo.current = 1;
  noMore.value = false;
  followList.value = [];
  getRelationList();
};

watch(
  () => [props.userId, props.following] as const,
  () => resetAndReload(),
  { immediate: true }
);

watch(
  () => viewerLoggedIn.value,
  () => resetAndReload()
);
</script>

<style lang="scss" scoped>
.follow-list {
  box-sizing: border-box;
  list-style: none;
  padding: 0;
}

.follow-card {
  padding: 0 16px;
  display: flex;
  align-items: center;
  height: 70px;
  position: relative;
  border-bottom: 1px solid var(--border-color);
}

.info-box {
  padding: 0 12px;
  flex: 1;

  .follow-name {
    display: block;
    color: var(--font-primary-1);
    font-size: 14px;
    margin-bottom: 6px;
    cursor: pointer;

    &:hover {
      color: var(--primary-hover-color);
    }
  }

  .follow-sign {
    font-size: 12px;
    color: var(--font-primary-3);
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 1;
    line-clamp: 1;
    -webkit-box-orient: vertical;
  }
}

.follow-btn {
  width: 70px;
  box-sizing: border-box;
  padding: 0;
  line-height: 26px;
  height: 26px;
  border-radius: 6px;
  margin-top: 5px;
  text-align: center;
  color: #fff;
  font-size: 12px;
  cursor: pointer;
  background: var(--primary-color);

  &:hover {
    background: var(--primary-hover-color);
  }
}
</style>