# 文章相关接口

## 上传文章信息

#### 请求URL
- `/api/v1/article/uploadArticleInfo`
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`

#### 参数

| 参数名      | 必选 | 类型   | 说明                  |
| :---------- | :--- | :----- | --------------------- |
| title       | 是   | string | 文章标题              |
| cover       | 否   | string | 封面图url             |
| copyright   | 是   | bool   | 是否为原创文章        |
| tags        | 否   | string | 文章标签，使用`,`分隔 |
| content     | 是   | string | 文章内容              |
| partitionId | 是   | int    | 文章分区              |

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

## 获取文章状态信息

#### 请求URL
- `/api/v1/article/getArticleStatus?aid=文章ID `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "article": {
      "aid": 1,
      "title": "标题",
      "cover": "封面url",
      "content": "文章内容",
      "copyright": true,
      "status": 0,
      "partitionId": 1,
      "tags": "",
    }
  },
  "msg": "ok"
}
```
#### 返回参数说明 

##### 文章信息article
| 参数名      | 类型   | 说明           |
| :---------- | :----- | -------------- |
| aid         | int    | 文章ID         |
| title       | string | 标题           |
| cover       | string | 封面URL        |
| content     | string | 文章内容       |
| copyright   | bool   | 是否为原创文章 |
| status      | int    | 审核状态       |
| partitionId | int    | 分区ID         |
| tags        | string | 文章标签       |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取上传的文章列表

#### 请求URL
- `/api/v1/article/getUploadArticle?page=页码&pageSize=内容数量 `
  
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
    "articles": [
      {
        "aid": 2,
        "title": "测试1",
        "cover": "",
        "content": "",
        "status": 0,
        "copyright": true,
        "createdAt": "2021-07-16T08:49:54Z",
        "tags": "",
        "clicks": 10,
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名   | 类型   | 说明         |
| :------- | :----- | ------------ |
| total    | int    | 数量         |
| articles | object | 文章信息数组 |


##### 文章信息
| 参数名    | 类型   | 说明           |
| :-------- | :----- | -------------- |
| aid       | int    | 文章ID         |
| title     | string | 标题           |
| cover     | string | 封面图url      |
| content   | string | 内容           |
| status    | int    | 文章审核状态   |
| copyright | bool   | 是否为原创文章 |
| clicks    | int    | 文章点击量     |
| createdAt | string | 上传时间       |
| tags      | string | 文章标签       |

#### 备注
返回的文章内容仅为文章的前200字

<!-- ************************ 分隔符 ************************ -->

## 编辑文章信息

#### 请求URL
- `/api/v1/article/editArticleInfo`
  
#### 请求方式
- PUT 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名  | 必选 | 类型   | 说明                  |
| :------ | :--- | :----- | --------------------- |
| aid     | 是   | int    | 文章ID                |
| title   | 是   | string | 文章标题              |
| cover   | 是   | string | 封面图url             |
| tags    | 否   | string | 文章标签，使用`,`分隔 |
| content | 否   | string | 文章内容              |

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

## 删除文章

#### 请求URL
- `/api/v1/article/deleteArticle/文章ID`
  
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

## 获取所有的文章列表

#### 请求URL
- `/api/v1/article/getAllArticleList `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "articles": [
      {
        "aid": 1,
        "title": ""
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
| articles | object | 视频信息数组 |


##### 视频信息
| 参数名 | 类型   | 说明   |
| :----- | :----- | ------ |
| aid    | int    | 文章ID |
| title  | string | 标题   |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取文章信息

#### 请求URL
- `/api/v1/article/getArticleById?aid=文章ID`
  
#### 请求方式
- GET 

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
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
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型   | 说明         |
| :----- | :----- | ------------ |
| total  | int    | 数量         |
| list   | object | 文章信息数组 |


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

## 通过用户ID获取文章列表

#### 请求URL
- `/api/v1/article/getArticleByUser?userId=用户ID&page=页码&pageSize=内容数量`
  
#### 请求方式
- GET 

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "total": 1,
    "articles": [
      {
        "aid": 2,
        "title": "测试1",
        "cover": "",
        "content": "",
        "status": 0,
        "copyright": true,
        "createdAt": "2021-07-16T08:49:54Z",
        "tags": "",
        "clicks": 10,
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名   | 类型   | 说明         |
| :------- | :----- | ------------ |
| total    | int    | 数量         |
| articles | object | 文章信息数组 |


##### 文章信息
| 参数名    | 类型   | 说明           |
| :-------- | :----- | -------------- |
| aid       | int    | 文章ID         |
| title     | string | 标题           |
| cover     | string | 封面图url      |
| content   | string | 内容           |
| status    | int    | 文章审核状态   |
| copyright | bool   | 是否为原创文章 |
| clicks    | int    | 文章点击量     |
| createdAt | string | 上传时间       |
| tags      | string | 文章标签       |

#### 备注
返回的文章内容仅为文章的前200字

<!-- ************************ 分隔符 ************************ -->

## 获取随机文章列表

#### 请求URL
- `/api/v1/article/getRandomArticleList?size=内容数量`
  
#### 请求方式
- GET 

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "articles": [
      {
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
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名   | 类型   | 说明         |
| :------- | :----- | ------------ |
| articles | object | 文章信息数组 |

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

## 后台管理-获取审核列表

#### 请求URL
- `/api/v1/article/getReviewArticleList?page=页码&pageSize=内容数量 `
  
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
        "aid": 1,
        "uid": 1,
        "title": "测试1",
        "cover": "",
        "content": "",
        "createdAt": "2022-06-06T08:42:13.525Z",
        "copyright": true,
        "tags": "",
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
| list   | array | 文章信息 |


##### 文章信息
| 参数名      | 类型   | 说明           |
| :---------- | :----- | -------------- |
| aid         | int    | 文章ID         |
| uid         | int    | 用户ID         |
| title       | string | 标题           |
| cover       | string | 封面URL        |
| content     | string | 文章内容       |
| createdAt   | time   | 发布时间       |
| copyright   | bool   | 是否为原创视频 |
| tags        | string | 视频标签       |
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

## 后台管理-获取文章列表

#### 请求URL
- `/api/v1/article/getArticleListManage`
  
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
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型   | 说明         |
| :----- | :----- | ------------ |
| total  | int    | 数量         |
| list   | object | 文章信息数组 |


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

## 后台管理-删除文章

#### 请求URL
- `/api/v1/article/deleteArticleManage/文章ID`
  
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
