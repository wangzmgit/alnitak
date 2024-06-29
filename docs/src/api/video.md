# 视频相关接口

## 上传视频信息

#### 请求URL
- `/api/v1/video/uploadVideoInfo`
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`

#### 参数

| 参数名      | 必选 | 类型   | 说明                       |
| :---------- | :--- | :----- | -------------------------- |
| vid         | 是   | int    | 视频ID                     |
| title       | 是   | string | 视频标题                   |
| cover       | 是   | string | 封面图url                  |
| desc        | 否   | string | 视频简介 默认:"什么都没有" |
| copyright   | 是   | bool   | 是否为原创视频             |
| tags        | 否   | string | 视频标签，使用`,`分隔      |
| partitionId | 是   | int    | 视频分区                   |

#### 返回示例 
``` json
{
  "code": 200,
  "data": null,
  "msg":"ok"
}
```

#### 备注
需要先上传视频拿到视频ID

<!-- ************************ 分隔符 ************************ -->

## 获取视频状态信息

#### 请求URL
- `/api/v1/video/getVideoStatus?vid=视频ID `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "video": {
      "vid": 1,
      "title": "标题",
      "cover": "封面url",
      "desc": "视频简介",
      "copyright": true,
      "status": 0,
      "partitionId": 1,
      "tags": "",
      "duration": 10,
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
  },
  "msg": "ok"
}
```
#### 返回参数说明 

##### 视频信息video
| 参数名      | 类型   | 说明                      |
| :---------- | :----- | ------------------------- |
| vid         | int    | 视频ID                    |
| title       | string | 标题                      |
| cover       | string | 封面URL                   |
| desc        | string | 视频简介                  |
| copyright   | bool   | 是否为原创视频            |
| status      | int    | 视频审核状态              |
| partitionId | int    | 分区ID                    |
| tags        | string | 视频标签                  |
| duration    | float  | 视频时长                  |
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

<!-- ************************ 分隔符 ************************ -->

## 获取上传视频列表

#### 请求URL
- `/api/v1/video/getUploadVideo?page=页码&pageSize=内容数量 `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token `

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "total": 1,
    "videos": [
      {
        "vid": 2,
        "title": "测试1",
        "cover": "",
        "desc": "",
        "status": 0,
        "copyright": true,
        "createdAt": "2021-07-16T08:49:54Z",
        "clicks": 10,
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


##### 视频信息
| 参数名    | 类型   | 说明           |
| :-------- | :----- | -------------- |
| vid       | int    | 视频ID         |
| title     | string | 标题           |
| cover     | string | 封面图url      |
| desc      | string | 简介           |
| status    | int    | 视频审核状态   |
| copyright | bool   | 是否为原创视频 |
| clicks    | int    | 视频点击量     |
| createdAt | string | 上传时间       |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 编辑视频信息

#### 请求URL
- `/api/v1/video/editVideoInfo`
  
#### 请求方式
- PUT 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型   | 说明                  |
| :----- | :--- | :----- | --------------------- |
| vid    | 是   | int    | 视频ID                |
| title  | 是   | string | 视频标题              |
| cover  | 是   | string | 封面图url             |
| desc   | 否   | string | 视频简介              |
| tags   | 否   | string | 视频标签，使用`,`分隔 |

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

## 删除视频

#### 请求URL
- `/api/v1/video/deleteVideo/视频ID`
  
#### 请求方式
- DELETE 

#### 请求头
- `Authorization': token `

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

## 获取所有的视频列表

#### 请求URL
- `/api/v1/video/getAllVideoList `
  
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
        "title": ""
      }
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


##### 视频信息
| 参数名 | 类型   | 说明   |
| :----- | :----- | ------ |
| vid    | int    | 视频ID |
| title  | string | 标题   |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取视频信息

#### 请求URL
- `/api/v1/video/getVideoById?vid=视频ID `
  
#### 请求方式
- GET 

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
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
  },
  "msg": "ok"
}
```
#### 返回参数说明 

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

<!-- ************************ 分隔符 ************************ -->

## 获取视频分辨率信息

#### 请求URL
- `/api/v1/video/getResourceQuality?resourceId=资源ID`
  
#### 请求方式
- GET 

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "quality": [
      "640x360_500k",
      "854x480_900k",
      "1080x720_2000k",
      "1920x1080_3000k",
    ]
  },
  "msg": "ok"
}
```

#### 备注 
无

<!-- ************************ 分隔符 ************************ -->

## 获取视频文件

#### 请求URL
- `/api/v1/video/getVideoFile?resourceId=资源ID&quality=资源分辨率信息`
  
#### 请求方式
- GET 

#### 返回示例 

``` text
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-MEDIA-SEQUENCE:0
#EXT-X-ALLOW-CACHE:YES
#EXT-X-TARGETDURATION:15
#EXTINF:14.647967,
854x480_900k_300000.ts
...
#EXTINF:9.509500,
854x480_900k_300014.ts
#EXT-X-ENDLIST
```

#### 备注 
无

<!-- ************************ 分隔符 ************************ -->

## 获取视频切片

#### 请求URL
- `/api/v1/video/slice/文件?key=文件key`
  
#### 请求方式
- GET 

#### 返回示例 

如果存储为本地存储会返回文件

如果存储为OSS则会重定向到OSS上的文件

#### 备注 
该接口的请求参数为获取视频文件接口返回的信息

<!-- ************************ 分隔符 ************************ -->

## 通过用户ID获取视频列表

#### 请求URL
- `/api/v1/video/getVideoByUser?userId=用户ID&page=页码&pageSize=内容数量 `
  
#### 请求方式
- GET 

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "total": 1,
    "videos": [
      {
        "vid": 2,
        "title": "测试1",
        "cover": "",
        "desc": "",
        "status": 0,
        "copyright": true,
        "createdAt": "2021-07-16T08:49:54Z",
        "clicks": 10,
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

##### 视频信息
| 参数名    | 类型   | 说明           |
| :-------- | :----- | -------------- |
| vid       | int    | 视频ID         |
| title     | string | 标题           |
| cover     | string | 封面图url      |
| desc      | string | 简介           |
| status    | int    | 视频审核状态   |
| copyright | bool   | 是否为原创视频 |
| clicks    | int    | 视频点击量     |
| createdAt | string | 上传时间       |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取热门视频

#### 请求URL
- `/api/v1/video/getHotVideo?page=页码&pageSize=内容数量 `
  
#### 请求方式
- GET 

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
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型   | 说明         |
| :----- | :----- | ------------ |
| videos | object | 视频信息数组 |

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
内容数量最大为30

<!-- ************************ 分隔符 ************************ -->

## 随机获取指定分区的视频

#### 请求URL
- `/api/v1/video/getVideoListByPartition?size=内容数量&partitionId=分区ID `
  
#### 请求方式
- GET 

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
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型   | 说明         |
| :----- | :----- | ------------ |
| videos | object | 视频信息数组 |

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
内容数量最大为30

<!-- ************************ 分隔符 ************************ -->

## 获取相关推荐视频

#### 请求URL
- `/api/v1/video/getRelatedVideoList?vid=视频ID `
  
#### 请求方式
- GET 

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
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型   | 说明         |
| :----- | :----- | ------------ |
| videos | object | 视频信息数组 |

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

<!-- ************************ 分隔符 ************************ -->

## 搜索视频

#### 请求URL
- `/api/v1/video/searchVideo `
  
#### 请求方式
- POST 

####  请求头
- `"content-type": "application/json"`

#### 请求参数
| 参数名   | 必选 | 类型   | 说明             |
| :------- | :--- | :----- | ---------------- |
| page     | 是   | int    | 页码             |
| pageSize | 是   | int    | 页面数量，最大30 |
| keyWords | 是   | string | 关键词           |

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
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型   | 说明         |
| :----- | :----- | ------------ |
| videos | object | 视频信息数组 |

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
内容数量最大为30

<!-- ************************ 分隔符 ************************ -->

## 后台管理-获取审核列表

#### 请求URL
- `/api/v1/video/getReviewList?page=页码&page_size=内容数量 `
  
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
          "vid": 1,
          "uid": 1,
          "title": "测试1",
          "cover": "",
          "desc": "",
          "createdAt": "2022-06-06T08:42:13.525Z",
          "copyright": true,
          "tags": "",
          "duration": 100,
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
      ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型  | 说明     |
| :----- | :---- | -------- |
| total  | int   | 数量     |
| list   | array | 视频信息 |


##### 视频信息
| 参数名      | 类型   | 说明           |
| :---------- | :----- | -------------- |
| vid         | int    | 视频ID         |
| uid         | int    | 用户ID         |
| title       | string | 标题           |
| cover       | string | 封面URL        |
| desc        | string | 视频简介       |
| createdAt   | time   | 发布时间       |
| copyright   | bool   | 是否为原创视频 |
| tags        | string | 视频标签       |
| duration    | float  | 视频时长       |
| partitionId | int    | 分区ID         |
| author      | object | 作者信息       |

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

##### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 后台管理-获取待审核资源列表

#### 请求URL
- `/api/v1/video/getReviewResourceList?vid=视频ID `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "resources": [
      {
        "id": 1,
        "createdAt": "",
        "vid": 1,
        "title": "",
        "duration": 10,
        "status": 0,
      }
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名    | 类型   | 说明     |
| :-------- | :----- | -------- |
| id        | int    | 分集ID   |
| createdAt | time   | 上传时间 |
| vid       | int    | 视频ID   |
| title     | string | 标题     |
| duration  | float  | 视频时长 |
| status    | int    | 审核状态 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 后台管理-获取视频列表

#### 请求URL
- `/api/v1/video/getVideoListManage`
  
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
      "videos": [
        {
          "vid": 1,
          "title": "测试1",
          "cover": "",
          "desc": "",
          "created_at": "2022-06-06T08:42:13.525Z",
          "copyright": true,
          "author": {
            "uid": 1,
            "name": "user_1654250698886",
            "sign": "这个人很懒，什么都没有留下1",
            "avatar": "",
            "spacecover": "",
            "gender": 0,
          },
          "clicks": 86,
          "partition": 1
        },
      ]
  },
  "msg": "ok"
}
```

#### 返回示例 

``` json
{
  "code": 200,
    "data": {
      "total": 1,
      "list": [
        {
          "vid": 1,
          "uid": 1,
          "title": "测试1",
          "cover": "",
          "desc": "",
          "createdAt": "2022-06-06T08:42:13.525Z",
          "copyright": true,
          "tags": "",
          "clicks": 1,
          "duration": 100,
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
      ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型  | 说明     |
| :----- | :---- | -------- |
| total  | int   | 数量     |
| list   | array | 视频信息 |


##### 视频信息videos
| 参数名      | 类型   | 说明           |
| :---------- | :----- | -------------- |
| vid         | int    | 视频ID         |
| uid         | int    | 用户ID         |
| title       | string | 标题           |
| cover       | string | 封面URL        |
| desc        | string | 视频简介       |
| createdAt   | time   | 发布时间       |
| copyright   | bool   | 是否为原创视频 |
| tags        | string | 视频标签       |
| clicks      | string | 播放量         |
| duration    | float  | 视频时长       |
| partitionId | int    | 分区ID         |
| author      | object | 作者信息       |

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

##### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 后台管理-删除视频

#### 请求URL
- `/api/v1/video/deleteVideoManage/视频ID`
  
#### 请求方式
- DELETE 

#### 请求头
- `Authorization': token `

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

## 后台管理-获取视频分辨率信息

#### 请求URL
- `/api/v1/video/getResourceQualityManage?resourceId=资源ID`
  
#### 请求方式
- GET 

#### 请求头
- `Authorization': token `

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "quality": [
      "640x360_500k",
      "854x480_900k",
      "1080x720_2000k",
      "1920x1080_3000k",
    ]
  },
  "msg": "ok"
}
```

#### 备注 
无

<!-- ************************ 分隔符 ************************ -->

## 后台管理-获取视频文件

#### 请求URL
- `/api/v1/video/getVideoFileManage?resourceId=资源ID&quality=资源分辨率信息`
  
#### 请求方式
- GET 

#### 请求头
- `Authorization': token `

#### 返回示例 

``` text
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-MEDIA-SEQUENCE:0
#EXT-X-ALLOW-CACHE:YES
#EXT-X-TARGETDURATION:15
#EXTINF:14.647967,
854x480_900k_300000.ts
...
#EXTINF:9.509500,
854x480_900k_300014.ts
#EXT-X-ENDLIST
```

#### 备注 
无

