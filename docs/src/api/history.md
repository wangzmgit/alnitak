# 历史记录相关接口

## 记录历史记录

#### 请求URL
- `/api/v1/history/video/addHistory `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型  | 说明        |
| :----- | :--- | :---- | ----------- |
| vid    | 是   | int   | 视频id      |
| part   | 否   | int   | 分P (默认1) |
| time   | 是   | float | 时间        |

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

## 获取历史记录

#### 请求URL
- `/api/v1/history/video/getHistory?page=页码&page_size=内容数量 `
  
#### 请求方式
- GET 

#### 返回示例 
``` json
{
  "code": 200,
    "data": {
    "videos": [
      {
        "vid": 2,
        "uid": 2,
        "title": "测试",
        "cover": "",
        "desc": "什么都没有",
        "time": 10.6,
        "updatedAt": "2022-06-06T08:42:13.525Z"
      }
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 

| 参数名    | 类型   | 说明         |
| :-------- | :----- | ------------ |
| vid       | int    | 视频ID       |
| uid       | int    | 用户ID       |
| title     | string | 标题         |
| cover     | string | 封面URL      |
| desc      | string | 视频简介     |
| time      | float  | 进度         |
| updatedAt | time   | 最后更新时间 |

#### 备注 
无

<!-- ************************ 分隔符 ************************ -->

## 获取播放进度

#### 请求URL
- `/api/v1/history/video/getProgress?vid=视频ID&part=分P `
  
#### 请求方式
- GET 

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "part": 1,
    "progress": 10.6
  },
  "msg": "ok"
}
```

#### 返回参数说明 

| 参数名   | 类型  | 说明              |
| :------- | :---- | ----------------- |
| part     | int   | 分P（v1.2.4新增） |
| progress | float | 进度              |

#### 备注 
无
