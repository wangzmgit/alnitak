<template>
  <div class="archive-data">
    <!--点赞收藏-->
    <div class="archive-item">
      <el-icon :class="[likeAnimation, archive.hasLike ? 'active' : 'icon']" @click="likeClick">
        <like-icon></like-icon>
      </el-icon>
      <p>{{ stat.like }}</p>
    </div>
    <div class="archive-item">
      <el-icon :class="archive.hasCollect ? 'active' : 'icon'" @click="showCollect = true">
        <collect-icon></collect-icon>
      </el-icon>
      <p>{{ stat.collect }}</p>
    </div>
    <collection-list v-if="showCollect" :vid="vid" @close="closeCollectionCard"></collection-list>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount } from 'vue';
import { ElIcon } from 'element-plus';
import LikeIcon from "@/components/icons/LikeIcon.vue";
import CollectIcon from "@/components/icons/CollectIcon.vue";
import { getVideoArchiveStatAPI } from "@/api/archive";
import { getLikeVideoStatusAPI, likeVideoAPI, cancelLikeVideoAPI } from "@/api/like";
import { getCollectVideoStatusAPI } from '@/api/collect';
import CollectionList from './CollectionList.vue';

const props = defineProps<{
  vid: number;
}>();

// 点赞收藏数据
const stat = ref<{ like: number, collect: number }>({
  like: 0,
  collect: 0
});

const loading = ref(true);
const archive = reactive({ // 是否点赞收藏
  hasCollect: false,
  hasLike: false
})

//获取点赞收藏关注信息
const getArchiveStat = async () => {
  const res = await getVideoArchiveStatAPI(props.vid);
  if (res.data.code === statusCode.OK) {
    stat.value = res.data.data.stat;
  }
}

// 获取是否点赞
const getLikeStatus = async () => {
  const res = await getLikeVideoStatusAPI(props.vid);
  if (res.data.code === statusCode.OK) {
    archive.hasLike = res.data.data.like;
  }
}

// 获取是否收藏
const getCollectStatus = async () => {
  const res = await getCollectVideoStatusAPI(props.vid);
  if (res.data.code === statusCode.OK) {
    archive.hasCollect = res.data.data.collect;
  }
}

const likeAnimation = ref('');
const likeClick = async () => { // 点赞点赞按钮
  if (loading.value) return;
  if (!archive.hasLike) {
    //调用点赞接口
    await likeVideoAPI(props.vid);
    likeAnimation.value = 'like-active';
    stat.value.like++;
  } else {
    await cancelLikeVideoAPI(props.vid);
    stat.value.like--;
  }

  archive.hasLike = !archive.hasLike;
}


const showCollect = ref(false);
// 关闭收藏弹窗
const closeCollectionCard = (val: number) => {
  if (val === 1) {
    stat.value.collect++;
    archive.hasCollect = true;
  } else if (val === -1) {
    stat.value.collect--;
    archive.hasCollect = false;
  }
  
  showCollect.value = false;
}

onBeforeMount(async () => {
  await getArchiveStat();
  await getLikeStatus();
  await getCollectStatus();

  loading.value = false;
})
</script>

<style lang="scss" scoped>
.archive-data {
  height: 30px;

  .archive-item {
    float: left;
    user-select: none;
    margin-right: 20px;

    i {
      font-size: 26px;
      line-height: 30px;
      cursor: pointer;
    }

    p {
      font-size: 16px;
      float: right;
      margin: 0 6px;
      line-height: 30px;
    }

    .icon:hover {
      color: var(--primary-color);
    }

    .active {
      color: var(--primary-color);
    }

    .like-active {
      animation: scaleDraw .3s ease-in-out;
    }
  }
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
</style>