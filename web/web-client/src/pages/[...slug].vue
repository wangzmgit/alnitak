<template>
  <div class="error-page">
    <header-bar class="header-bar"></header-bar>
    <div class="result-container">
      <h1 class="code">404</h1>
      <h2 class="desc">页面未找到</h2>
      <p class="detail">请求的路径不存在</p>
      <p class="back-btn" @click="handleError">返回首页</p>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: `404 - ${globalConfig.title}`
})

const handleError = () => {
  navigateTo('/')
}
</script>

<style lang="scss" scoped>
.error-page {
  height: 100vh;

  .header-bar {
    position: fixed;
    top: 0;
    width: 100%;
    z-index: 999;
  }

  .result-container {
    text-align: center;
    width: 800px;
    margin-left: -400px;
    position: absolute;
    top: 33%;
    left: 50%;

    .code {
      margin: 0;
      font-size: 150px;
      line-height: 150px;
      font-weight: bold;
      color: var(--primary-color);
    }

    .desc {
      margin-top: 20px;
      font-size: 30px;
      color: var(--font-primary-1);
    }

    .detail {
      margin-top: 10px;
      font-size: 16px;
      color: var(--font-primary-3);
    }

    .back-btn {
      cursor: pointer;
      color: var(--primary-color);
      margin-top: 20px;
      font-size: 18px;

      &:hover {
        color: var(--primary-hover-color);
      }
    }
  }
}
</style>
