<template>
  <div class="collection-bg">
    <div class="collection-card">
      <div class="card-content">
        <n-input v-model:value="collectionName" resize="none" :rows="3" type="textarea" placeholder="请输入收藏夹名称" />
      </div>
      <div class="card-bottom">
        <div class="submit-btn" @click="addCollection">新建收藏夹</div>
        <div class="close-btn" @click="closeCard">返回</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { NInput, useMessage } from "naive-ui";
import { addCollectionAPI } from '@/api/collection';
import { statusCode } from '@/utils/status-code';

const emits = defineEmits(['close']);

const message = useMessage();

const collectionName = ref("");
const addCollection = async () => {
  if (!collectionName.value) {
    message.error('收藏夹名不能为空');
    return;
  }

  const res = await addCollectionAPI(collectionName.value);
  if (res.data.code === statusCode.OK) {
    collectionName.value = "";
    emits('close');
  } else {
    message.error('收藏夹创建失败');
  }
}

const closeCard = () => {
  emits('close');
}
</script>

<style lang="scss" scoped>
.collection-bg {
  top: 0;
  left: 0;
  z-index: 999;
  position: fixed;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.45);

  .collection-card {
    position: absolute;
    box-sizing: border-box;
    z-index: 2000;
    top: 50px;
    left: 0;
    width: 100%;
    height: calc(100vh - 50px);
    padding: 0 0;
    background-color: #fff;
    animation: fadein .3s ease-in;
    box-shadow: 16px 16px 50px -12px rgba(0, 0, 0, 0.8);
  }
}

.card-content {
  padding: 10px 10px;
}

.card-bottom {
  height: 76px;
  text-align: center;
  padding: 10px 10px;

  .submit-btn {
    font-size: 14px;
    width: 100%;
    height: 40px;
    line-height: 40px;
    margin-top: 18px;
    border: none;
    border-radius: 4px;
    background-color: var(--primary-color);
    color: #fff;
    cursor: pointer;

    &:hover {
      background-color: var(--primary-hover-color);
    }
  }

  .close-btn {
    font-size: 14px;
    width: 100%;
    height: 40px;
    line-height: 40px;
    margin-top: 18px;
    border: none;
    border-radius: 4px;
    background-color: #fff;
    color: #505050;
    cursor: pointer;
    border: 1px solid #eee;
  }
}
</style>