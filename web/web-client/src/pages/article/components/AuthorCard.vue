<template>
  <div class="article-up-info">
    <div class="up-left">
      <div class="avatar-container">
        <common-avatar :size="44" :url="author.avatar"></common-avatar>
      </div>
      <div class="info-container">
        <span class="up-name">{{ author.name }}</span>
        <span class="up-sign">{{ author.sign }}</span>
      </div>
    </div>
    <el-button class="follow-btn" type="primary" @click="followBtnClick">{{ btnText }}</el-button>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, computed } from "vue";
import { ElMessage } from "element-plus";
import { statusCode } from "@/utils/status-code";
import { relationCode } from "@/utils/relation-code";
import CommonAvatar from "@/components/common-avatar/index.vue";
import { getUserRelationAPI, followAPI, unfollowAPI } from "@/api/relation";

const props = defineProps<{
  author: UserInfoType
}>()


const relation = ref(relationCode.NOT_FOLLOWING);
const getUserRelation = async () => {
  if (props.author.uid) {
    const res = await getUserRelationAPI(props.author.uid);
    if (res.data.code === statusCode.OK) {
      relation.value = res.data.data.relation;
    }
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
  if (props.author.uid) {
    const reqFunc = relation.value === relationCode.NOT_FOLLOWING ? followAPI : unfollowAPI;
    const res = await reqFunc(props.author.uid);
    if (res.data.code === statusCode.OK) {
      getUserRelation();
    } else {
      ElMessage.error(res.data.msg);
    }
  }
}

onBeforeMount(() => {
  getUserRelation();
})
</script>

<style lang="scss" scoped>
.article-up-info {
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 44px;
  justify-content: space-between;
  margin: 0 auto 20px;

  .up-left {
    display: flex;
    align-items: center;

    .avatar-container {
      flex-shrink: 0;
      display: block;
      position: relative;
      width: 44px;
      height: 44px;
      margin-right: 12px;
    }

    .info-container {
      .up-name {
        font-size: 14px;
        max-width: 150px;
        color: #212121;
        line-height: 23px;
        height: 23px;
        text-overflow: ellipsis;
        overflow: hidden;
        white-space: nowrap;
        display: inline-block;
        vertical-align: middle;
      }

      .up-sign {
        display: block;
        font-size: 13px;
        color: #999;
        line-height: 20px;
      }
    }
  }

  .follow-btn {
    width: 128px;
    height: 32px;
  }
}
</style>