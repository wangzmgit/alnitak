# 粉丝关注相关接口

## 通过用户ID获取关注列表

#### 请求URL
- `/api/v1/relation/getFollowings?userId=用户ID&page=页码&pageSize=内容数量 `
  
#### 请求方式
- GET 

#### 返回示例 

``` json
{
  "code": 200,
    "data": {
    "users": [
      {
        "relation": 1,
        "user": {
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
      }
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名   | 类型   | 说明                                |
| :------- | :----- | ----------------------------------- |
| relation | int    | 关注状态，0未关注、1关注、2互相关注 |
| user     | string | 用户信息                            |

##### 用户信息
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

## 通过用户ID获取粉丝列表

#### 请求URL
- `/api/v1/relation/getFollowers?userId=用户ID&page=页码&pageSize=内容数量 `
  
#### 请求方式
- GET 

#### 返回示例 
``` json
{
  "code": 200,
    "data": {
    "users": [
      {
        "relation": 1,
        "user": {
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
      }
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名   | 类型   | 说明                                |
| :------- | :----- | ----------------------------------- |
| relation | int    | 关注状态，0未关注、1关注、2互相关注 |
| user     | string | 用户信息                            |

##### 用户信息
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

## 获取关注和粉丝数

#### 请求URL
- `/api/v1/relation/getFollowCount?userId=用户ID `
  
#### 请求方式
- GET 

#### 返回示例 

``` json

{
  "code": 200,
  "data": {
    "followers": 1,
    "following": 0
  },
  "msg": "ok"
}
```

#### 返回参数说明 

| 参数名    | 类型 | 说明   |
| :-------- | :--- | ------ |
| followers | int  | 粉丝数 |
| following | int  | 关注数 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取用户关系

#### 请求URL
- `/api/v1/relation/getUserRelation?userId=目标用户ID `
  
#### 请求方式
- GET 

#### 请求头
- `Authorization': token `

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "relation": 1,
  },
  "msg": "ok"
}
```

#### 返回参数说明 

| 参数名   | 类型 | 说明                                |
| :------- | :--- | ----------------------------------- |
| relation | int  | 关注状态，0未关注、1关注、2互相关注 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

##  关注用户

#### 请求URL
- `/api/v1/relation/follow `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`

#### 参数

| 参数名 | 必选 | 类型 | 说明       |
| :----- | :--- | :--- | ---------- |
| id     | 是   | int  | 目标用户id |

### 返回示例 

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

##  取消关注用户

#### 请求URL
- `/api/v1/relation/unfollow `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`
  
#### 参数

| 参数名 | 必选 | 类型 | 说明       |
| :----- | :--- | :--- | ---------- |
| id     | 是   | int  | 目标用户id |

### 返回示例 

``` json
{
  "code": 200,
  "data": null,
  "msg":"ok"
}
```

#### 备注
无