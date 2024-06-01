<template>
  <div class="side-bar">
    <div class="side-toolbar">
      <div class="toolbar-item">
        <el-icon :class="[likeAnimation, archive.hasLike ? 'active' : 'icon']" @click="likeClick">
          <like-icon></like-icon>
        </el-icon>
      </div>
      <div class="toolbar-item">
        <el-icon :class="[collectAnimation, archive.hasCollect ? 'active' : 'icon']" @click="collectClick">
          <collect-icon></collect-icon>
        </el-icon>
      </div>
      <div class="toolbar-item">
        <el-icon @click="scrollToComment">
          <comment-fill-icon></comment-fill-icon>
        </el-icon>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount } from 'vue';
import { ElIcon } from 'element-plus';
import LikeIcon from "@/components/icons/LikeIcon.vue";
import CollectIcon from "@/components/icons/CollectIcon.vue";
import CommentFillIcon from "@/components/icons/CommentFillIcon.vue";
import { likeArticleAPI, cancelLikeArticleAPI, getLikeArticleStatusAPI } from "@/api/like";
import { collectArticleAPI, cancelCollectArticleAPI, getCollectArticleStatusAPI } from '@/api/collect';

const emit = defineEmits(["statChange", "scrollToComment"]);
const props = defineProps<{
  aid: number;
}>();

const loading = ref(true);
const archive = reactive({ // 是否点赞收藏
  hasCollect: false,
  hasLike: false
})

// 获取是否点赞
const getLikeStatus = async () => {
  const res = await getLikeArticleStatusAPI(props.aid);
  if (res.data.code === statusCode.OK) {
    archive.hasLike = res.data.data.like;
  }
}

// 获取是否收藏
const getCollectStatus = async () => {
  const res = await getCollectArticleStatusAPI(props.aid);
  if (res.data.code === statusCode.OK) {
    archive.hasCollect = res.data.data.collect;
  }
}

const likeAnimation = ref('');
const likeClick = async () => { // 点赞点赞按钮
  if (loading.value) return;
  if (!archive.hasLike) {
    await likeArticleAPI(props.aid);
    likeAnimation.value = 'like-active';
    emit("statChange", "like", 1);
  } else {
    await cancelLikeArticleAPI(props.aid);
    emit("statChange", "like", -1);
  }

  archive.hasLike = !archive.hasLike;
}

const collectAnimation = ref('');
const collectClick = async () => { // 点赞点赞按钮
  if (loading.value) return;
  if (!archive.hasLike) {
    await collectArticleAPI(props.aid);
    collectAnimation.value = 'collect-active';
    emit("statChange", "collect", 1);

  } else {
    await cancelCollectArticleAPI(props.aid);
    emit("statChange", "collect", -1);
  }

  archive.hasCollect = !archive.hasCollect;
}

const scrollToComment = () => {
  emit("scrollToComment");
}

onBeforeMount(async () => {
  await getLikeStatus();
  await getCollectStatus();

  loading.value = false;
})
</script>

<style lang="scss" scoped>
.side-bar {
  position: fixed;
  margin-left: 930px;
  bottom: -60px;
  transition: .5s ease-in-out;
  margin-bottom: 80px;
  z-index: 10;

  .side-toolbar {
    .toolbar-item {
      position: relative;
      color: #8a919f;
      background: #fff;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 50px;
      width: 50px;
      margin-bottom: 20px;
      border-radius: 50%;
      cursor: pointer;
      font-size: 18px;
    }
  }
}

.icon:hover {
  color: var(--primary-hover-color);
}

.active {
  color: var(--primary-color);
}

.like-active {
  animation: scaleDraw .3s ease-in-out;
}

.collect-active {
  animation: scaleDraw .3s ease-in-out;
}

@keyframes scaleDraw {
  0% {
    transform: scale(1);
    /*开始为原始大小*/
  }

  25% {
    transform: scale(1.2);
    /*放大1.1倍*/
  }

  100% {
    transform: scale(1);
  }
}

@media (max-width: 1400px) {
  .side-bar {
    margin-left: 730px;
  }
}
</style>