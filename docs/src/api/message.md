# 消息相关接口

## 获取公告

#### 请求URL
- `/api/v1/message/getAnnounce?page=页码&pageSize=内容数量 `
  
#### 请求方式
- GET 

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "total": 1,
    "announces": [
      {
        "id": 1,
        "title": "测试",
        "content": "123",
        "url": "123",
        "createdAt": "2021-07-29T13:46:21Z",
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名    | 类型  | 说明     |
| :-------- | :---- | -------- |
| total     | int   | 总数量   |
| announces | array | 公告列表 |

##### 公告内容announces
| 参数名    | 类型   | 说明       |
| :-------- | :----- | ---------- |
| id        | int    | 公告ID     |
| title     | string | 标题       |
| content   | string | 内容       |
| url       | string | 指向的链接 |
| createdAt | time   | 发布时间   |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 添加公告

#### 请求URL
- `/api/v1/message/addAnnounce `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token `
- `"content-type": "application/json",`

#### 参数

| 参数名  | 必选 | 类型   | 说明       |
| :------ | :--- | :----- | ---------- |
| title   | 是   | string | 标题       |
| content | 是   | string | 内容       |
| url     | 否   | string | 指向的地址 |
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

## 删除公告

#### 请求URL
- `/api/v1/message/deleteAnnounce/公告ID `
  
#### 请求方式
- DELETE 

####  请求头
- `Authorization': token`

#### 返回示例 
``` json
{
  "code": 200,
  "data": null,
  "msg":"ok"
}
```

##### 备注 
无

<!-- ************************ 分隔符 ************************ -->

## 获取点赞消息

#### 请求URL
- `/api/v1/message/getLikeMsg?page=页码&pageSize=内容数量 `
  
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
    "messages": [
      {
        "id": 1,
        "cid": 1,
        "sid": 1,
        "createdAt": "2021-07-29T13:46:21Z",
        "type": 1,
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
        },
        "video": {
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
          "resources": [
            {
              "id": 1,
              "createdAt": "",
              "vid": 1,
              "title": "",
              "duration": 10,
              "status": 0
            }
          ]
        },
        "article": {
          "aid": 2,
          "uid": 2,
          "title": "测试1",
          "cover": "",
          "content": "",
          "status": 0,
          "copyright": true,
          "createdAt": "2021-07-16T08:49:54Z",
          "tags": "",
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
        },
      }
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名   | 类型  | 说明     |
| :------- | :---- | -------- |
| total    | int   | 总数量   |
| messages | array | 消息列表 |

##### 点赞消息
| 参数名    | 类型   | 说明                  |
| :-------- | :----- | --------------------- |
| id        | int    | 消息                  |
| cid       | int    | 内容ID（视频或文章）  |
| sid       | int    | 发送用户ID            |
| createdAt | string | 发布时间              |
| type      | int    | 类型，0：视频;1：文章 |
| user      | object | 点赞用户信息          |
| video     | object | 视频信息              |
| article   | object | 文章信息              |

##### 视频信息
| 参数名      | 类型   | 说明                      |
| :---------- | :----- | ------------------------- |
| vid         | int    | 视频ID                    |
| uid         | int    | 作者ID                    |
| title       | string | 标题                      |
| cover       | string | 封面URL                   |
| desc        | string | 视频简介                  |
| createdAt   | string | 上传时间                  |
| copyright   | bool   | 是否为原创视频            |
| tags        | string | 视频标签                  |
| duration    | float  | 视频时长                  |
| clicks      | int    | 视频点击量                |
| partitionId | int    | 分区ID                    |
| author      | object | 作者信息                  |
| resource    | array  | 视频资源,多个代表多个分集 |

##### 文章信息
| 参数名      | 类型   | 说明           |
| :---------- | :----- | -------------- |
| aid         | int    | 文章ID         |
| uid         | int    | 用户ID         |
| title       | string | 标题           |
| cover       | string | 封面图url      |
| content     | string | 内容           |
| status      | int    | 文章审核状态   |
| copyright   | bool   | 是否为原创文章 |
| clicks      | int    | 文章点击量     |
| createdAt   | string | 上传时间       |
| tags        | string | 文章标签       |
| partitionId | int    | 分区ID         |
| author      | object | 作者信息       |

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
实际使用中，视频信息和文章信息不会同时返回，具体返回哪一个以`type`为准

<!-- ************************ 分隔符 ************************ -->

## 获取@消息

#### 请求URL
- `/api/v1/message/getAtMsg?page=页码&pageSize=内容数量 `
  
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
    "messages": [
      {
        "id": 1,
        "cid": 1,
        "sid": 1,
        "createdAt": "2021-07-29T13:46:21Z",
        "type": 1,
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
        },
        "video": {
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
          "resources": [
            {
              "id": 1,
              "createdAt": "",
              "vid": 1,
              "title": "",
              "duration": 10,
              "status": 0
            }
          ]
        },
        "article": {
          "aid": 2,
          "uid": 2,
          "title": "测试1",
          "cover": "",
          "content": "",
          "status": 0,
          "copyright": true,
          "createdAt": "2021-07-16T08:49:54Z",
          "tags": "",
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
        },
      }
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名   | 类型  | 说明     |
| :------- | :---- | -------- |
| total    | int   | 总数量   |
| messages | array | 消息列表 |

##### 点赞消息
| 参数名    | 类型   | 说明                  |
| :-------- | :----- | --------------------- |
| id        | int    | 消息                  |
| cid       | int    | 内容ID（视频或文章）  |
| sid       | int    | 发送用户ID            |
| createdAt | string | 发布时间              |
| type      | int    | 类型，0：视频;1：文章 |
| user      | object | 点赞用户信息          |
| video     | object | 视频信息              |
| article   | object | 文章信息              |

##### 视频信息
| 参数名      | 类型   | 说明                      |
| :---------- | :----- | ------------------------- |
| vid         | int    | 视频ID                    |
| uid         | int    | 作者ID                    |
| title       | string | 标题                      |
| cover       | string | 封面URL                   |
| desc        | string | 视频简介                  |
| createdAt   | string | 上传时间                  |
| copyright   | bool   | 是否为原创视频            |
| tags        | string | 视频标签                  |
| duration    | float  | 视频时长                  |
| clicks      | int    | 视频点击量                |
| partitionId | int    | 分区ID                    |
| author      | object | 作者信息                  |
| resource    | array  | 视频资源,多个代表多个分集 |

##### 文章信息
| 参数名      | 类型   | 说明           |
| :---------- | :----- | -------------- |
| aid         | int    | 文章ID         |
| uid         | int    | 用户ID         |
| title       | string | 标题           |
| cover       | string | 封面图url      |
| content     | string | 内容           |
| status      | int    | 文章审核状态   |
| copyright   | bool   | 是否为原创文章 |
| clicks      | int    | 文章点击量     |
| createdAt   | string | 上传时间       |
| tags        | string | 文章标签       |
| partitionId | int    | 分区ID         |
| author      | object | 作者信息       |

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
实际使用中，视频信息和文章信息不会同时返回，具体返回哪一个以`type`为准

<!-- ************************ 分隔符 ************************ -->

## 获取回复消息

#### 请求URL
- `/api/v1/message/getReplyMsg?page=页码&pageSize=内容数量 `
  
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
    "comments": [
      {
        "id": 1,
        "cid": 1,
        "sid": 1,
        "createdAt": "2022-06-20T13:42:40.625Z",
        "content": "测试",
        "targetReplyContent": "",
        "rootContent": "",
        "commentId": 1,
        "parentId": "",
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
        },
        "video": {
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
          "resources": [
            {
              "id": 1,
              "createdAt": "",
              "vid": 1,
              "title": "",
              "duration": 10,
              "status": 0
            }
          ]
        },
        "article": {
          "aid": 2,
          "uid": 2,
          "title": "测试1",
          "cover": "",
          "content": "",
          "status": 0,
          "copyright": true,
          "createdAt": "2021-07-16T08:49:54Z",
          "tags": "",
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
        },
      }
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名   | 类型  | 说明     |
| :------- | :---- | -------- |
| total    | int   | 总数量   |
| messages | array | 消息列表 |

##### 评论回复消息
| 参数名             | 类型   | 说明                  |
| :----------------- | :----- | --------------------- |
| id                 | int    | id                    |
| cid                | int    | 视频ID                |
| sid                | int    | 发送用户ID            |
| uid                | int    | 用户ID                |
| createdAt          | string | 发布时间              |
| content            | string | 评论或回复内容        |
| targetReplyContent | string | 回复目标内容          |
| rootContent        | string | 根评论内容            |
| commentId          | int    | 评论ID                |
| type               | int    | 类型，0：视频;1：文章 |
| user               | object | 作者信息              |
| video              | object | 视频信息              |
| article            | object | 文章信息              |

##### 视频信息
| 参数名      | 类型   | 说明                      |
| :---------- | :----- | ------------------------- |
| vid         | int    | 视频ID                    |
| uid         | int    | 作者ID                    |
| title       | string | 标题                      |
| cover       | string | 封面URL                   |
| desc        | string | 视频简介                  |
| createdAt   | string | 上传时间                  |
| copyright   | bool   | 是否为原创视频            |
| tags        | string | 视频标签                  |
| duration    | float  | 视频时长                  |
| clicks      | int    | 视频点击量                |
| partitionId | int    | 分区ID                    |
| author      | object | 作者信息                  |
| resource    | array  | 视频资源,多个代表多个分集 |

##### 文章信息
| 参数名      | 类型   | 说明           |
| :---------- | :----- | -------------- |
| aid         | int    | 文章ID         |
| uid         | int    | 用户ID         |
| title       | string | 标题           |
| cover       | string | 封面图url      |
| content     | string | 内容           |
| status      | int    | 文章审核状态   |
| copyright   | bool   | 是否为原创文章 |
| clicks      | int    | 文章点击量     |
| createdAt   | string | 上传时间       |
| tags        | string | 文章标签       |
| partitionId | int    | 分区ID         |
| author      | object | 作者信息       |

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
实际使用中，视频信息和文章信息不会同时返回，具体返回哪一个以`type`为准

<!-- ************************ 分隔符 ************************ -->

## 私信websocket

#### 请求URL
- `/api/v1/message/ws?token=用户token `

#### 返回示例 
```json
{
  "fid": 1,//来自用户ID
  "content": "",
}
```

#### 备注
请求时使用ws或者wss协议，返回消息使用了base64编码

<!-- ************************ 分隔符 ************************ -->

## 获取私信列表

#### 请求URL
- `/api/v1/message/getWhisperList `
  
#### 请求方式
- GET 

#### 请求头
- `Authorization': token `

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "messages": [
      {
        "createdAt": "2021-07-29T13:46:21Z",
        "content": "",
        "fid": 1,
        "status": true,
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
        },
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名    | 类型   | 说明           |
| :-------- | :----- | -------------- |
| createdAt | string | 最新消息时间   |
| content   | string | 最新消息内容   |
| fid       | int    | 发送者ID       |
| status    | bool   | 是否已读       |
| user      | Object | 发送者用户信息 |

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

## 获取私信详情

#### 请求URL
- `/api/v1/message/getWhisperDetails?fid=目标用户ID&page=页码&pageSize=内容数量 `
  
#### 请求方式
- GET 

####  请求头
- `Authorization: token`

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "messages": [
      {
        "fid": 2,
        "fromId": 1,
        "content": "123",
        "createdAt": "2021-07-29T10:54:42Z"
      },
      {
        "fid": 2,
        "from_id": 1,
        "content": "456",
        "createdAt": "2021-07-29T10:55:31Z"
      }
    ],
  },
  "msg": "ok"
}
```

##### 返回参数说明 
| 参数名     | 类型   | 说明       |
| :--------- | :----- | ---------- |
| fid        | int    | 关联用户ID |
| fromId    | int    | 发送者id   |
| content    | string | 消息内容   |
| createdAt | string   | 发送时间   |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 发送私信

#### 请求URL
- `/api/v1/message/sendWhisper `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名  | 必选 | 类型   | 说明       |
| :------ | :--- | :----- | ---------- |
| fid     | 是   | int    | 目标用户id |
| content | 是   | string | 内容       |

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

## 私信已读

#### 请求URL
- `/api/v1/message/readWhisper `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`

#### 参数

| 参数名 | 必选 | 类型 | 说明       |
| :----- | :--- | :--- | ---------- |
| id     | 是   | int  | 目标用户id |

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