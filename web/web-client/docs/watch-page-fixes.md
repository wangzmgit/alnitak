## watch 页面修复记录

### 问题清单

| 问题 | 根因 | 修复 |
|------|------|------|
| 构建元数据 500 | 旧 `.output/server/index.mjs` 进程残留（PID 23340 占用 3000 端口），build ID 与浏览器缓存不匹配 | `kill` 旧进程后重建 |
| Nuxt 水合不匹配 | 构建 500 的附带现象，非真实 SSR/CSR 内容差异 | 域名访问验证无错误 |
| 嵌套 `useAsyncData` | `asyncGetVideoInfoAPI` 内部已含 `useAsyncData`，外部再包一层导致数据结构异常和 500 | 恢复原始单层调用模式 |

### watch.vue 防御性修正

```diff
- 顶级 await 绑定 PGC
+ watch(videoInfo, () => { ... }, { immediate: true }) + process.client 守卫

- 无守卫直接 navigateTo('/404')
+ if (process.client) { navigateTo('/404') }

- 直接访问 videoApiData.code
+ if (videoApiData.value) { /* 使用 videoApiData.value.code */ }
```

#### 1. PGC 绑定移至 watch + process.client

避免 SSR 期间执行 PGC 客户端专属逻辑导致水合不匹配：

```typescript
// before: 顶级 await 绑定
const { data: pgcData } = await useFetch(...)

// after: 客户端安全绑定
watch(videoInfo, () => {
  if (process.client) {
    bindPGC(videoInfo.value?.pgcId)
  }
}, { immediate: true })
```

#### 2. 404 导航 SSR 保护

```typescript
// before: 无守卫
navigateTo('/404')

// after: 仅客户端执行
if (process.client) {
  navigateTo('/404')
}
```

#### 3. API 响应空指针保护

```typescript
// before: 可能访问 undefined.code
if (videoApiData.code === 0) { ... }

// after: 判空后访问
if (videoApiData.value && videoApiData.value.code === 0) { ... }
```

### 水合调试开关

`nuxt.config.ts` 中有条件开启（仅开发环境）：

```typescript
__VUE_PROD_HYDRATION_MISMATCH_DETAILS__: process.dev ? true : false
```

生产环境默认关闭。
