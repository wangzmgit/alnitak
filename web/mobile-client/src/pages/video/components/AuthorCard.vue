<template>
  <div class="up-panel-container">
    <div class="up-info-container">
      <div class="up-info-left">
        <div class="up-avatar-wrap">
          <nuxt-link class="up-avatar" :to="`/user/${props.info.uid}`" target="_blank">
            <common-avatar :size="48" :url="props.info.avatar"></common-avatar>
          </nuxt-link>
        </div>
      </div>
      <div class="up-info-right">
        <div class="up-info-detail">
          <div class="up-detail">
            <div class="up-detail-top">
              <nuxt-link class="up-name" :to="`/user/${props.info.uid}`" target="_blank">{{ props.info.name }}</nuxt-link>
            </div>
            <div class="up-description up-detail-bottom">{{ props.info.sign }}</div>
          </div>
        </div>
        <!-- 关注按钮部分 -->
        <div class="up-info-btn-panel">
          <div class="up-info-btn" :class="disabledBtn ? 'btn-disabled' : ''" @click="followBtnClick">
            {{ btnText }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, computed } from "vue";
import { ElMessage } from "element-plus";
import { statusCode } from "@/utils/status-code";
import { relationCode } from "@/utils/relation-code";
import { getUserRelationAPI, followAPI, unfollowAPI } from "@/api/relation";

const props = defineProps<{
  info: UserInfoType
}>()

const disabledBtn = ref(false);
const relation = ref(relationCode.NOT_FOLLOWING);
const getUserRelation = async () => {
  const res = await getUserRelationAPI(props.info.uid);
  if (res.data.code === statusCode.OK) {
    relation.value = res.data.data.relation;
  } else {
    disabledBtn.value = true;
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
  if (disabledBtn.value) return;
  const reqFunc = relation.value === relationCode.NOT_FOLLOWING ? followAPI : unfollowAPI;
  const res = await reqFunc(props.info.uid);
  if (res.data.code === statusCode.OK) {
    getUserRelation();
  } else {
    ElMessage.error(res.data.msg);
  }
}

onBeforeMount(() => {
  getUserRelation();
})
</script>

<style lang="scss" scoped>
.up-info-container {
  box-sizing: border-box;
  height: 104px;
  display: flex;
  align-items: center;
}

.up-info-left {
  .up-avatar-wrap {
    width: 48px;
    height: 48px;
    flex-shrink: 0;
    display: flex;
    justify-content: center;
    align-items: center;

    .up-avatar {
      display: block;
      width: 100%;
      height: 100%;
      border-radius: 50%;
      background-color: #C9CCD0;
    }
  }
}

.up-info-right {
  margin-left: 12px;
  flex: 1;
  overflow: auto;

  .up-info-detail {
    margin-bottom: 5px;

    .up-detail {
      .up-detail-top {
        display: flex;
        align-items: center;

        .up-name {
          font-size: 15px;
          color: #18191C;
          font-weight: 500;
          position: relative;
          white-space: nowrap;
          text-overflow: ellipsis;
          overflow: hidden;
          margin-right: 12px;
          max-width: calc(100% - 12px - 56px);
        }
      }

      .up-description {
        margin-top: 2px;
        font-size: 13px;
        line-height: 16px;
        height: 16px;
        color: #9499A0;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}

.up-info-btn-panel {
  .up-info-btn {
    width: calc(100% - 22px);
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
    background: var(--primary-color);

    &:hover {
      background: var(--primary-hover-color);
    }
  }
}

.btn-disabled {
  background-color: var(--primary-hover-color) !important;
  cursor: not-allowed !important;

  &:hover {
    background-color: var(--primary-hover-color) !important;
  }
}
</style>