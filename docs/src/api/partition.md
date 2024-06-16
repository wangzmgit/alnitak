# 分区相关接口

## 获取分区列表

#### 请求URL
- `/api/v1/partition/getPartitionList?type=类型(0：视频;1：文章) `
  
#### 请求方式
- GET 

#### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "partitions": [
      {
        "id": 1,
        "name": "动画",
        "type": 0,
        "parentId": 0,
        "createdAt": "2021-07-16T08:49:54Z",
      },
      {
        "id": 2,
        "name": "MAD·AMV",
        "type": 0,
        "parentId": 1,
        "createdAt": "2021-07-16T08:49:54Z",
      }
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名    | 类型   | 说明                  |
| :-------- | :----- | --------------------- |
| id        | int    | 分区ID                |
| name      | string | 分区名                |
| type      | int    | 类型，0：视频;1：文章 |
| parentId  | int    | 所属分区ID            |
| createdAt | string | 分区创建时间          |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 新增分区

#### 请求URL
- `/api/v1/partition/addPartition `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名   | 必选 | 类型   | 说明                  |
| :------- | :--- | :----- | --------------------- |
| name     | 是   | string | 分区名                |
| type     | 是   | int    | 类型，0：视频;1：文章 |
| parentId | 否   | int    | 所属分区id            |

#### 返回示例 

```json
{
  "code": 200,
  "data": null,
  "msg": "ok"
}
```

#### 备注 
无

<!-- ************************ 分隔符 ************************ -->

## 删除分区

#### 请求URL
- `/api/v1/partition/deletePartition/分区ID `
  
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