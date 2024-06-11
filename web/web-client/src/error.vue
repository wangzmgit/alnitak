<template>
  <div class="error-page" >
    <header-bar class="header-bar "></header-bar>
    <div class="result-container">
      <h1 class="code">{{ error.statusCode }}</h1>
      <h2 class="desc">{{ getStatusDesc(error.statusCode) }}</h2>
      <p class="back-btn" @click="handleError">返回首页</p>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  error: any
}>()

useHead({
  title: `${props.error.statusCode} - ${globalConfig.title}`
})

const getStatusDesc = (code: number) => {
  switch (code) {
    case 404:
      return 'Not Found';
    case 500:
      return '服务器错误'
  }
}

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