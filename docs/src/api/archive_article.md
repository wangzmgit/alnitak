# 点赞收藏相关接口

## 获取点赞收藏数据

#### 请求URL
- `/api/v1/archive/article/stat?aid=文章ID `
  
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
- `/api/v1/archive/article/like `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型 | 说明   |
| :----- | :--- | :--- | ------ |
| aid    | 是   | int  | 文章ID |

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
- `/api/v1/archive/article/cancelLike `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token `
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型 | 说明   |
| :----- | :--- | :--- | ------ |
| id     | 是   | int  | 文章ID |

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
- `/api/v1/archive/article/hasLike?aid=文章ID `
  
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
- `/api/v1/archive/article/collect `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名     | 必选 | 类型  | 说明             |
| :--------- | :--- | :---- | ---------------- |
| aid        | 是   | int   | 文章ID           |

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

## 取消收藏

#### 请求URL
- `/api/v1/archive/article/cancelCollect `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token `
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型 | 说明   |
| :----- | :--- | :--- | ------ |
| id     | 是   | int  | 文章ID |

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

## 获取收藏状态

#### 请求URL
- `/api/v1/archive/article/hasCollect?aid=文章ID `
  
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

