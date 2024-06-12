# 用户相关接口

## 用户注册

#### 请求URL
- `/api/v1/auth/register `
  
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

<!-- ************************ 分隔符 ************************ -->

## 账号密码登录

#### 请求URL
- `/api/v1/auth/login `
  
#### 请求方式
- POST 

####  请求头
- `"content-type": "application/json",`

#### 参数
| 参数名    | 必选 | 类型   | 说明       |
| :-------- | :--- | :----- | ---------- |
| email     | 是   | string | 邮箱       |
| password  | 是   | string | 密码       |
| captchaId | 否   | string | 人机验证ID |

#### 返回示例 
登录成功返回值
```json
{
  "code": 200,
  "data": {
    "token": "", 
    "refreshToken": "",
  },
  "msg":"ok"
}
```

连续登录失败三次会返回值
```json
{
  "code": -1,
  "data": {"captchaId": ""},
  "msg":"需要人机验证"
}
```

#### 返回参数说明 
| 参数名       | 类型   | 说明                          |
| :----------- | :----- | ----------------------------- |
| token        | string | 请求使用的token，1小时内有效  |
| refreshToken | string | 刷新token用的token，7天内有效 |

#### 备注
1. 同一邮箱连续登录失败三次会返回需要人机验证，此时需要等待30分钟或者调用人机验证接口并通过滑块验证才可以继续登录。
2. 登录成功会在Cookie中添加`user_id`和`user_id_md5`，其中Cookie中的内容暂时没有在后端使用，如需使用请自行修改相关代码。

<!-- ************************ 分隔符 ************************ -->

## 邮箱验证登录

#### 请求URL
- `/api/v1/auth/login/email `
  
#### 请求方式
- POST 

####  请求头
- `"content-type": "application/json",`

#### 参数

| 参数名    | 必选 | 类型   | 说明       |
| :-------- | :--- | :----- | ---------- |
| email     | 是   | string | 邮箱       |
| code      | 是   | string | 验证码     |
| captchaId | 否   | string | 人机验证ID |

#### 返回示例 
登录成功返回值
```json
{
  "code": 200,
  "data": {
    "token": "", 
    "refreshToken": "",
  },
  "msg":"ok"
}
```

连续登录失败三次会返回值
```json
{
  "code": -1,
  "data": {"captchaId": ""},
  "msg":"需要人机验证"
}
```

#### 返回参数说明 
| 参数名       | 类型   | 说明                          |
| :----------- | :----- | ----------------------------- |
| token        | string | 请求使用的token，1小时内有效  |
| refreshToken | string | 刷新token用的token，7天内有效 |

#### 备注
1. 同一邮箱连续登录失败三次会返回需要人机验证，此时需要等待30分钟或者调用人机验证接口并通过滑块验证才可以继续登录。
2. 登录成功会在Cookie中添加`user_id`和`user_id_md5`，其中Cookie中的内容暂时没有在后端使用，如需使用请自行修改相关代码。

<!-- ************************ 分隔符 ************************ -->

## 更新TOKEN

#### 请求URL
- `/api/v1/auth/updateToken `
  
#### 请求方式
- POST 

####  请求头
- `"content-type": "application/json",`

#### 参数
| 参数名       | 必选 | 类型   | 说明         |
| :----------- | :--- | :----- | ------------ |
| refreshToken | 是   | string | RefreshToken |

#### 返回示例 
成功返回值
```json
{
  "code": 200,
  "data": {
    "token": "", 
    "refreshToken": "",
  },
  "msg":"ok"
}
```

token失效
```json
{
  "code": 2000,
  "data": null,
  "msg":"error"
}
```

#### 返回参数说明 
| 参数名       | 类型   | 说明                          |
| :----------- | :----- | ----------------------------- |
| token        | string | 请求使用的token，1小时内有效  |
| refreshToken | string | 刷新token用的token，7天内有效 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 退出登录

#### 请求URL
- `Authorization': token `
- `/api/v1/auth/logout `
  
#### 请求方式
- POST 

####  请求头
- `"content-type": "application/json",`

#### 参数
| 参数名       | 必选 | 类型   | 说明         |
| :----------- | :--- | :----- | ------------ |
| refreshToken | 是   | string | RefreshToken |

#### 返回示例 
成功返回值
```json
{
  "code": 200,
  "data": null,
  "msg":"ok"
}
```

#### 备注
该接口会把RefreshToken失效并清除Cookie

<!-- ************************ 分隔符 ************************ -->

## 清除Cookie

#### 请求URL
- `/api/v1/auth/clearCookie `
  
#### 请求方式
- POST 

####  请求头
- `"content-type": "application/json",`

#### 参数
无
#### 返回示例 
无

#### 备注
该接只清除Cookie，对RefreshToken没有影响


<!-- ************************ 分隔符 ************************ -->

## 修改密码验证

#### 请求URL
- `/api/v1/auth/resetpwdCheck?eamil=用户邮箱&captchaId=人机验证滑块ID `
  
#### 请求方式
- GET 

#### 返回示例 

未进行人机验证返回值
```json
{
  "code": -1,
  "data": {"captchaId": ""},
  "msg":"需要人机验证"
}
```

成功返回值
``` json
{
  "code": 200,
  "data": null,
  "msg":"ok"
}
```

#### 备注
需要先进行人机验证，验证通过后可修改密码

<!-- ************************ 分隔符 ************************ -->

## 修改密码

#### 请求URL
- `/api/v1/auth/modifyPwd `
  
#### 请求方式
- POST 

####  请求头
- `"content-type": "application/json",`

#### 参数

| 参数名    | 必选 | 类型   | 说明           |
| :-------- | :--- | :----- | -------------- |
| email     | 是   | string | 邮箱           |
| password  | 是   | string | 新密码         |
| code      | 是   | string | 邮箱验证码     |
| captchaId | 否   | string | 人机验证滑块ID |


#### 返回示例 

```json
  {
    "code": 200,
    "data": null,
    "msg":"ok"
  }
```

#### 备注 
需要先调用[修改密码验证](#修改密码验证)接口 