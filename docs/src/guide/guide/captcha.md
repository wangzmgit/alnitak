# 人机验证说明

:::tip 前言
以下部分简要说明了人机验证流程，使您更容易理解相关功能，或在二次开发中更轻松地引入该功能。

我们为该项目专门开发了拼图生成库，详细信息请参阅[jigsaw](github.com/wangzmgit/jigsaw)。
:::

## 人机验证状态码

```go
const CAPTCHA_STATUS_ABSENT = 0   // 人机验证状态不存在
const CAPTCHA_STATUS_PASS = 2     // 人机验证状态已通过
```

## 返回人机验证状态码
对于需要进行人机验证的相关接口，可以通过以下代码读取人机验证状态。
```go
cache.GetCaptchaStatus(captchaId)
```

如果人机验证状态不存在，可以通过以下代码，生成`captchaId`并返回`需要人机验证`状态
```go
captchaId := cache.CreateCaptchaStatus()

resp.Result(ctx, -1, gin.H{"captchaId": captchaId}, "需要人机验证")
```

验证通过后，请务必使用以下代码删除当前`captchaId`的验证状态
```go
cache.DelCaptchaStatus(captchaId)
```

## 前端处理验证码
前端通过以下代码引入和使用滑块验证组件，其中`v-model:show`控制组件的显示和隐藏，
`captcha-id`为后端返回的`captchaId`，`success`为验证成功的事件。
用户停止拖动滑块时会将`captchaId`和`滑块的X轴坐标`传回后端并进行验证。

```js
import SliderCaptcha from "@/components/slider-captcha/index.vue";
```
```vue
<slider-captcha 
  v-model:show="showCaptcha"
  :captcha-id="captchaId" 
  @success="captchaSuccess">
</slider-captcha>
```


