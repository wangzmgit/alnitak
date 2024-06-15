# 点赞收藏相关接口

## 获取点赞收藏数据

#### 请求URL
- `/api/v1/archive/video/stat?vid=视频ID `
  
#### 请求方式
- GET 

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "stat": {
      "collect": 1,
      "like": 1
    }
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名  | 类型 | 说明     |
| :------ | :--- | -------- |
| collect | int  | 收藏数量 |
| like    | int  | 点赞数量 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 点赞

#### 请求URL
- `/api/v1/archive/video/like `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型 | 说明   |
| :----- | :--- | :--- | ------ |
| vid    | 是   | int  | 视频ID |

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

## 取消点赞

#### 请求URL
- `/api/v1/archive/video/cancelLike `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token `
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型 | 说明   |
| :----- | :--- | :--- | ------ |
| id     | 是   | int  | 视频ID |

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

## 获取点赞状态

#### 请求URL
- `/api/v1/archive/video/hasLike?vid=视频ID `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "like": true
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型 | 说明     |
| :----- | :--- | -------- |
| like   | bool | 是否点赞 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 添加收藏

#### 请求URL
- `/api/v1/archive/video/collect `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名     | 必选 | 类型  | 说明             |
| :--------- | :--- | :---- | ---------------- |
| vid        | 是   | int   | 视频ID           |
| addList    | 否   | []int | 添加的收藏夹数组 |
| cancelList | 否   | []int | 移除的收藏夹数组 |

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

## 获取已收藏的收藏夹

#### 请求URL
- `/api/v1/archive/video/getCollectInfo?vid=视频ID `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "collectionIds": [1, 2]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名        | 类型  | 说明               |
| :------------ | :---- | ------------------ |
| collectionIds | []int | 视频所在的收藏夹id |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取收藏状态

#### 请求URL
- `/api/v1/archive/video/hasCollect?vid=视频ID `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "collect": true
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名  | 类型 | 说明     |
| :------ | :--- | -------- |
| collect | bool | 是否收藏 |

#### 备注
无

