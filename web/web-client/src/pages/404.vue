<template>
  <div class="error-page" >
    <header-bar class="header-bar "></header-bar>
    <div class="result-container">
      <h1 class="code">404</h1>
      <h2 class="desc">Not Found</h2>
      <p class="back-btn" @click="handleError">返回首页</p>
    </div>
  </div>
</template>

<script setup lang="ts">

useHead({
  title: `404 - ${globalConfig.title}`
})

const handleError = () => clearError({ redirect: '/' })
</script>

<style lang="scss" scoped>
.error-page {
  height: 100vh;

  .header-bar {
    position: fixed;
    top: 0;
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
      // color: var(--primary-color);
      margin-top: 20px;
      font-size: 30px;
    }

    .back-btn {
      cursor: pointer;
      color: var(--primary-color);

      &:hover {
        color: var(--primary-hover-color);
      }
    }
  }
}
</style>