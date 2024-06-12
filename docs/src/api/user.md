# 用户相关接口

## 通过用户ID获取用户信息

#### 请求URL
- `/api/v1/user/getUserBaseInfo?userId=用户ID `
  
#### 请求方式
- GET 

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "userInfo": {
      "uid": 1,
      "name": "",
      "sign": "",
      "email": "",
      "phone": "",
      "avatar": "",
      "gender": 1,
      "spaceCover": "",
      "birthday": "",
      "createdAt": "",
    }
  },
  "msg": "ok"
}
```

#### 返回参数说明 

| 参数名     | 类型   | 说明                           |
| :--------- | :----- | ------------------------------ |
| uid        | int    | 用户ID                         |
| name       | string | 用户名                         |
| sign       | string | 个性签名                       |
| email      | string | 邮箱（内容为空）               |
| phone      | string | 手机号（内容为空）             |
| avatar     | string | 头像                           |
| gender     | int    | 用户性别，0：未知;1：男；2：女 |
| spacecover | string | 用户空间封面图                 |
| birthday   | time   | 生日                           |
| createdAt  | time   | 注册时间                       |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 用户获取个人信息

#### 请求URL
- `/api/v1/user/getUserInfo`
  
#### 请求方式
- GET 

#### 请求头
- `Authorization': token `
- `"content-type": "application/json"`

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "userInfo": {
      "uid": 1,
      "name": "",
      "sign": "",
      "email": "",
      "phone": "",
      "avatar": "",
      "gender": 1,
      "spaceCover": "",
      "birthday": "",
      "createdAt": "",
    }
  },
  "msg":"ok"
}
```

#### 返回参数说明 

| 参数名     | 类型   | 说明                           |
| :--------- | :----- | ------------------------------ |
| uid        | int    | 用户ID                         |
| name       | string | 用户名                         |
| sign       | string | 个性签名                       |
| email      | string | 邮箱                           |
| phone      | string | 手机号                         |
| avatar     | string | 头像                           |
| gender     | int    | 用户性别，0：未知;1：男；2：女 |
| spacecover | string | 用户空间封面图                 |
| birthday   | time   | 生日                           |
| createdAt  | time   | 注册时间                       |


#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 用户修改个人信息

#### 请求URL
- `/api/v1/user/editUserInfo`
  
#### 请求方式
- PUT 

#### 请求头
- `Authorization': token `
- `"content-type": "application/json"`

#### 参数
| 参数名     | 必选 | 类型               | 说明                 |
| :--------- | :--- | :----------------- | -------------------- |
| avatar     | 是   | string             | 用户头像             |
| name       | 是   | string             | 用户名               |
| gender     | 否   | int                | 性别（默认为0未知）  |
| birthday   | 是   | string（日期格式） | 生日（默认1970-1-1） |
| sign       | 否   | string             | 个性签名             |
| spaceCover | 是   | string             | 用户空间封面图       |

#### 返回示例 

``` json
{
  "code": 200,
  "data": null,
  "msg": "ok"
}
```

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 后台管理-获取用户列表

#### 请求URL
- `/api/v1/user/getUserListManage`
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`

#### 请求参数
| 参数名   | 必选 | 类型 | 说明              |
| :------- | :--- | :--- | ----------------- |
| page     | 是   | int  | 页码              |
| pageSize | 是   | int  | 页面数量，最大100 |

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "total": 1,
    "list": [
      {
        "uid": 1,
        "name": "昵称",
        "sign": "个性签名",
        "email": "1****1@qq.com",
        "phone": "",
        "avatar": "",
        "gender": 1,
        "spaceCover": "",
        "birthday":"",
        "createdAt": "",
        "role": "",
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 

| 参数名     | 类型   | 说明                           |
| :--------- | :----- | ------------------------------ |
| uid        | int    | 用户ID                         |
| name       | string | 用户名                         |
| sign       | string | 个性签名                       |
| email      | string | 用户邮箱                       |
| phone      | string | 手机号                         |
| avatar     | string | 头像                           |
| spaceCover | string | 用户空间封面图                 |
| gender     | int    | 用户性别，0：未知;1：男；2：女 |
| birthday   | time   | 生日                           |
| createdAt  | time   | 注册时间                       |
| role       | string | 用户角色代码                   |


#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 后台管理-修改用户信息

#### 请求URL
- `/api/v1/user/editUserInfoManage`
  
#### 请求方式
- PUT 

#### 请求头
- `Authorization': token `
- `"content-type": "application/json"`

#### 参数
| 参数名     | 必选 | 类型   | 说明           |
| :--------- | :--- | :----- | -------------- |
| uid        | 是   | int    | 用户id         |
| avatar     | 是   | string | 用户头像       |
| spaceCover | 是   | string | 用户空间封面图 |
| name       | 是   | string | 用户名         |
| email      | 是   | string | 邮箱           |
| sign       | 否   | string | 个性签名       |

#### 返回示例 

``` json
{
  "code": 200,
  "data": null,
  "msg": "ok"
}
```

#### 备注 
无

<!-- ************************ 分隔符 ************************ -->

## 后台管理-设置用户角色

#### 请求URL
- `/api/v1/user/editUserRole`
  
#### 请求方式
- PUT 

#### 请求头
- `Authorization': token `
- `"content-type": "application/json"`

#### 参数
| 参数名 | 必选 | 类型   | 说明     |
| :----- | :--- | :----- | -------- |
| uid    | 是   | int    | 用户id   |
| code   | 是   | string | 角色代码 |

#### 返回示例 

``` json
{
  "code": 200,
  "data": null,
  "msg": "ok"
}
```

#### 备注 
无

<!-- ************************ 分隔符 ************************ -->

## 后台管理-删除用户

#### 请求URL
- `/api/v1/user/deleteUser/用户ID`
  
#### 请求方式
- DELETE 

#### 请求头
- `Authorization': token `
- `"content-type": "application/json"`

#### 返回示例 

``` json
{
  "code": 200,
  "data": null,
  "msg": "ok"
}
```

#### 备注 
无