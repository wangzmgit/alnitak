# 评论回复相关接口

## 获取评论回复

#### 请求URL
- `/api/v1/comment/video/getComment?vid=视频ID&page=页码&pageSize=内容数量 `
  
#### 请求方式
- GET 

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "total": 1,
    "comments": [
      {
        "id": 1,
        "createdAt": "2022-06-20T13:42:40.625Z",
        "uid": 1,
        "content": "测试",
        "atUsernames": "",
        "atUserIds": "",
        "parentId": "",
        "replyCount": 2,
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
        "reply": [
          {
            "id": 1,
            "createdAt": "2022-06-20T13:42:40.625Z",
            "uid": 1,
            "content": "测试",
            "atUsernames": "",
            "atUserIds": "",
            "replyUserId": "",
            "replyUserName": "",
            "parentId": "",
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
          }
        ]
      }
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名   | 类型   | 说明         |
| :------- | :----- | ------------ |
| total    | int    | 数量         |
| comments | object | 评论回复数组 |

##### 评论comment
| 参数名      | 类型   | 说明                                                 |
| :---------- | :----- | ---------------------------------------------------- |
| id          | int    | 用户ID                                               |
| createdAt   | string | 发布时间                                             |
| uid         | int    | 用户ID                                               |
| content     | string | 评论内容                                             |
| atUsernames | string | @用户的用户名，使用`,`分隔                           |
| atUserIds   | string | @用户的ID，使用`,`分隔 ，与`atUsernames`中的顺序一致 |
| parentId    | int    | 所属评论ID                                           |
| replyCount  | int    | 回复数量                                             |
| author      | object | 作者信息                                             |
| reply       | object | 回复数组                                             |

##### 回复reply
| 参数名        | 类型   | 说明                                                 |
| :------------ | :----- | ---------------------------------------------------- |
| id            | int    | 用户ID                                               |
| createdAt     | string | 发布时间                                             |
| uid           | int    | 用户ID                                               |
| content       | string | 评论内容                                             |
| atUsernames   | string | @用户的用户名，使用`,`分隔                           |
| atUserIds     | string | @用户的ID，使用`,`分隔 ，与`atUsernames`中的顺序一致 |
| replyUserName | string | 回复用户的用户名                                     |
| replyUserId   | string | 回复用户的ID                                         |
| parentId      | int    | 所属评论ID                                           |
| author        | object | 作者信息                                             |

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
该接口返回的total为评论和回复总数，而评论和回复是树结构的，因此返回的总数无法用于列表的分页


## 获取回复详情

##### 请求URL
- `/api/v1/comment/video/getReply?commentId=评论ID&page=页码&pageSize=内容数量 `
  
##### 请求方式
- GET 

##### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "replies": [
      {
        "id": 1,
        "createdAt": "2022-06-20T13:42:40.625Z",
        "uid": 1,
        "content": "测试",
        "atUsernames": "",
        "atUserIds": "",
        "replyUserId": "",
        "replyUserName": "",
        "parentId": "",
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
      }
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名  | 类型   | 说明     |
| :------ | :----- | -------- |
| replies | object | 回复数组 |

##### 回复reply
| 参数名        | 类型   | 说明                                                 |
| :------------ | :----- | ---------------------------------------------------- |
| id            | int    | 用户ID                                               |
| createdAt     | string | 发布时间                                             |
| uid           | int    | 用户ID                                               |
| content       | string | 评论内容                                             |
| atUsernames   | string | @用户的用户名，使用`,`分隔                           |
| atUserIds     | string | @用户的ID，使用`,`分隔 ，与`atUsernames`中的顺序一致 |
| replyUserName | string | 回复用户的用户名                                     |
| replyUserId   | string | 回复用户的ID                                         |
| parentId      | int    | 所属评论ID                                           |
| author        | object | 作者信息                                             |

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
无

<!-- ************************ 分隔符 ************************ -->

## 发布评论回复

#### 请求URL
- `/api/v1/comment/video/addComment`
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数
| 参数名        | 必选 | 类型     | 说明                                 |
| :------------ | :--- | :------- | ------------------------------------ |
| cid           | 是   | int      | 视频ID                               |
| content       | 是   | string   | 评论内容                             |
| at            | 否   | []string | @的用户名数组                        |
| parentID      | 否   | int      | 所属评论ID                           |
| replyUserID   | 否   | int      | 回复用户的ID                         |
| replyUserName | 否   | string   | 回复用户的用户名                     |
| replyContent  | 否   | string   | 回复的评论或回复的内容，用于发送通知 |

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "comment": {
      "id": 1,
      "createdAt": "2022-06-20T13:42:40.625Z",
      "uid": 1,
      "content": "测试",
      "atUsernames": "",
      "atUserIds": "",
      "replyUserId": "",
      "replyUserName": "",
      "parentId": "",
    }
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名        | 类型   | 说明                                                 |
| :------------ | :----- | ---------------------------------------------------- |
| id            | int    | 用户ID                                               |
| createdAt     | string | 发布时间                                             |
| uid           | int    | 用户ID                                               |
| content       | string | 评论内容                                             |
| atUsernames   | string | @用户的用户名，使用`,`分隔                           |
| atUserIds     | string | @用户的ID，使用`,`分隔 ，与`atUsernames`中的顺序一致 |
| replyUserName | string | 回复用户的用户名                                     |
| replyUserId   | string | 回复用户的ID                                         |
| parentId      | int    | 所属评论ID                                           |
| author        | object | 作者信息                                             |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 删除评论

#### 请求URL
- `/api/v1/comment/video/deleteComment/评论或回复ID `
  
#### 请求方式
- DELETE 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

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

## 获取评论列表

#### 请求URL
- `/api/v1/comment/video/getCommentList?vid=视频ID&page=页码&pageSize=内容数量 `
  
#### 请求方式
- GET 

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
        "uid": 1,
        "createdAt": "2022-06-20T13:42:40.625Z",
        "content": "测试",
        "targetReplyContent": "",
        "rootContent": "",    
        "commentId": 1,
        "parentId": "",
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
        "target": {
          "uid": 2,
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
        }
      }
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名   | 类型   | 说明         |
| :------- | :----- | ------------ |
| total    | int    | 数量         |
| comments | object | 评论回复数组 |

##### 评论comment
| 参数名             | 类型   | 说明             |
| :----------------- | :----- | ---------------- |
| id                 | int    | id               |
| cid                | int    | 视频ID           |
| sid                | int    | 发送用户ID       |
| uid                | int    | 用户ID           |
| createdAt          | string | 发布时间         |
| content            | string | 评论或回复内容   |
| targetReplyContent | string | 回复目标内容     |
| rootContent        | string | 根评论内容       |
| commentId          | int    | 评论ID           |
| author             | object | 作者信息         |
| target             | object | 回复目标用户信息 |
| video              | object | 视频信息         |

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

##### 视频信息video
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

##### 视频资源resource
| 参数名    | 类型   | 说明     |
| :-------- | :----- | -------- |
| id        | int    | 分集ID   |
| createdAt | string | 上传时间 |
| vid       | int    | 视频ID   |
| title     | string | 标题     |
| duration  | float  | 视频时长 |
| status    | int    | 审核状态 |

#### 备注
无