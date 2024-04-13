<template>
  <ul class="follow-list" v-infinite-scroll="scrollLoad">
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
import { onBeforeMount, ref, reactive } from "vue";
import { getFollowingAPI, getFollowersAPI, followAPI, unfollowAPI, getUserRelationAPI } from '@/api/relation';

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

const followBtnClick = async (relation: RelationListType) => {
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

onBeforeMount(() => {
  getRelationList();
})
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
  border-bottom: 1px solid #eee;
}

.info-box {
  padding: 0 12px;
  flex: 1;

  .follow-name {
    display: block;
    color: #333;
    font-size: 14px;
    margin-bottom: 6px;
    cursor: pointer;

    &:hover {
      color: var(--primary-hover-color);
    }
  }

  .follow-sign {
    font-size: 12px;
    color: #666;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 1;
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