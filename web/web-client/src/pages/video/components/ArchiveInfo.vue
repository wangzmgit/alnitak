<template>
  <div class="archive-data">
    <!--点赞收藏-->
    <div class="archive-item">
      <el-icon :class="[likeAnimation, true ? 'active' : 'icon']" @click="likeClick">
        <like-icon></like-icon>
      </el-icon>
      <p>1</p>
    </div>
    <div class="archive-item">
      <el-icon :class="true ? 'active' : 'icon'" @click="showCollect = true">
        <collect-icon></collect-icon>
      </el-icon>
      <p>1</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ElIcon } from 'element-plus';
import LikeIcon from "@/components/icons/LikeIcon.vue";
import CollectIcon from "@/components/icons/CollectIcon.vue";



const showCollect = ref(false);
const likeAnimation = ref('');

//点赞点赞按钮
const likeClick = () => {
  likeAnimation.value = 'like-active';
}


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