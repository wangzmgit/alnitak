# 用户相关接口

## 用户注册

#### 请求URL
- ` http://域名/api/v1/auth/register `
  
#### 请求方式
- POST 

####  请求头
- `"content-type": "application/json",`

#### 参数

| 参数名   | 必选 | 类型   | 说明   |
| :------- | :--- | :----- | ------ |
| email    | 是   | string | 邮箱   |
| password | 是   | string | 密码   |
| code     | 是   | string | 验证码 |

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