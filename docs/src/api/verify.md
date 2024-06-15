# 验证相关接口

## 获取滑块验证

#### 请求URL
- `/api/v1/verify/captcha/get?captchaId=人机验证ID `
  
#### 请求方式
- GET 

#### 返回示例 

``` json
{
  "code": 200,
    "data": {
    "slider_captcha": {
      "slider_img": "",
      "bg_img": "",
      "y": 1,
    }
  },
  "msg": "ok"
}
```

#### 返回参数说明 

| 参数名     | 类型   | 说明            |
| :--------- | :----- | --------------- |
| slider_img | string | base64滑块图    |
| bg_img     | string | base64背景图    |
| y          | int    | 滑块左上角y坐标 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 验证滑块

#### 请求URL
- `/api/v1/verify/captcha/validate `
  
#### 请求方式
- POST 

####  请求头
- `"content-type": "application/json",`

#### 参数

| 参数名    | 必选 | 类型   | 说明            |
| :-------- | :--- | :----- | --------------- |
| captchaId | 是   | string | 人机验证ID      |
| x         | 是   | int    | 滑块左上角x坐标 |

#### 返回示例 

``` json
{
  "code": 200,
  "data": null,
  "msg":"ok"
}
```

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 发送邮箱验证码

#### 请求URL
- `/api/v1/verify/getEmailCode`
  
#### 请求方式
- POST 

####  请求头
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型   | 说明     |
| :----- | :--- | :----- | -------- |
| email  | 是   | string | 目标邮箱 |
| captchaId | 是   | string | 人机验证ID      |

#### 返回示例 

```json
  {
    "code": 200,
    "data": null,
    "msg":"ok"
  }
```

#### 备注 
无