# 收藏夹相关接口

## 获取收藏夹列表

#### 请求URL
- `/api/v1/collection/getCollectionList `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "total": 1,
    "collections": [
      {
        "id": 2,
        "name": "测试1",
        "cover": "",
        "desc": "",
        "open": true,
        "createdAt": "2021-07-16T08:49:54Z",
      }
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名      | 类型   | 说明           |
| :---------- | :----- | -------------- |
| total       | int    | 数量           |
| collections | object | 收藏夹信息数组 |

##### 收藏夹信息
| 参数名    | 类型   | 说明         |
| :-------- | :----- | ------------ |
| id        | int    | 收藏夹id     |
| name      | string | 收藏夹名称   |
| cover     | string | 收藏夹封面图 |
| desc      | string | 简介         |
| open      | bool   | 是否公开     |
| createdAt | string | 创建时间     |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取收藏夹信息

#### 请求URL
- `/api/v1/collection/getCollectionInfo?id=收藏夹id `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': access_token`

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "collection": {
      "id": 2,
      "uid": 1,
      "name": "测试1",
      "cover": "",
      "desc": "",
      "open": true,
      "createdAt": "2021-07-16T08:49:54Z",
      "author": {
        "uid": 1,
        "name": "user_1654250698886",
        "sign": "这个人很懒，什么都没有留下",
        "avatar": "",
        "spacecover": "",
        "gender": 0,
      }
    },
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名    | 类型   | 说明         |
| :-------- | :----- | ------------ |
| id        | int    | 收藏夹id     |
| uid       | int    | 用户id       |
| name      | string | 收藏夹名称   |
| cover     | string | 收藏夹封面图 |
| desc      | string | 简介         |
| open      | bool   | 是否公开     |
| createdAt | string | 创建时间     |
| author    | object | 作者信息     |

##### 作者信息author
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
对于公开收藏夹不需要登录携带请求头，对于非公开收藏夹需要登录，且仅限创建者可以访问。

<!-- ************************ 分隔符 ************************ -->

## 新建收藏夹

#### 请求URL
- `/api/v1/collection/addCollection `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型   | 说明       |
| :----- | :--- | :----- | ---------- |
| name   | 是   | string | 收藏夹名称 |

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

## 编辑收藏夹信息

#### 请求URL
- `/api/v1/collection/editCollection `
  
#### 请求方式
- PUT 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型   | 说明                |
| :----- | :--- | :----- | ------------------- |
| id     | 否   | int    | 收藏夹id            |
| name   | 是   | string | 收藏夹名称          |
| cover  | 否   | string | 收藏夹封面图        |
| desc   | 否   | string | 简介                |
| open   | 否   | bool   | 是否公开(默认false) |

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

## 删除收藏夹

#### 请求URL
- `/api/v1/collection/deleteCollection/收藏夹ID `
  
#### 请求方式
- DELETE 

####  请求头
- `"content-type": "application/json",`
- `Authorization': access_token`

#### 参数

| 参数名 | 必选 | 类型 | 说明     |
| :----- | :--- | :--- | -------- |
| id     | 是   | int  | 收藏夹id |

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

## 获取收藏夹的视频列表

#### 请求URL
- `/api/v1/video/collect?page=页码&pageSize=内容数量&cid=收藏夹id `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "videos": [
      {
        "vid": 1,
        "uid": 1,
        "title": "标题",
        "cover": "封面url",
        "desc": "视频简介",
        "createdAt": "",
        "copyright": true,
        "tags": "",
        "duration": 10,
        "clicks": 10,
        "partitionId": 1,
        "author": {
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
        },
        "resources": []
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型   | 说明         |
| :----- | :----- | ------------ |
| total  | int    | 数量         |
| videos | object | 视频信息数组 |

##### 视频信息video
| 参数名      | 类型   | 说明                        |
| :---------- | :----- | --------------------------- |
| vid         | int    | 视频ID                      |
| uid         | int    | 作者ID                      |
| title       | string | 标题                        |
| cover       | string | 封面URL                     |
| desc        | string | 视频简介                    |
| createdAt   | string | 上传时间                    |
| copyright   | bool   | 是否为原创视频              |
| tags        | string | 视频标签                    |
| duration    | float  | 视频时长                    |
| clicks      | int    | 视频点击量                  |
| partitionId | int    | 分区ID                      |
| author      | object | 作者信息                    |
| resource    | array  | 视频资源,该接口中不会有数据 |

##### 作者信息author
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
如果收藏夹是公开的则不需要携带请求头