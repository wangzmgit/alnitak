# 轮播图相关接口

## 获取轮播图

#### 请求URL
- `/api/v1/carousel/getCarousel?partitionId=分区ID `
  
#### 请求方式
- GET 

#### 返回示例 

``` json
{
  "code": 200,
    "data": {
    "carousels": [
      {
        "id": 1,
        "img": "",
        "url": "",
        "title": "",
        "color": "",
        "use": true,
        "partitionId": 1,
        "createdAt": "2021-08-01T10:39:25Z"
      },
    ]
  },
  "msg": "ok"
}
```
#### 返回参数说明
| 参数名      | 类型   | 说明               |
| :---------- | :----- | ------------------ |
| id          | int    | 轮播图id           |
| img         | string | 图片地址           |
| url         | string | 单击图片前往的地址 |
| title       | string | 标题               |
| color       | string | 主题色             |
| use         | bool   | 是否启用           |
| partitionId | int    | 分区ID             |
| createdAt   | time   | 发布时间           |

#### 备注 
无

<!-- ************************ 分隔符 ************************ -->

## 上传轮播图信息

#### 请求URL
- `/api/v1/carousel/addCarousel `
  
#### 请求方式
- POST 

####  请求头
- `"content-type": "application/json",`
- `Authorization': access_token`

#### 参数

| 参数名      | 必选 | 类型   | 说明                |
| :---------- | :--- | :----- | ------------------- |
| img         | 是   | string | 轮播图url           |
| url         | 否   | string | 单击轮播图打开的url |
| title       | 否   | string | 标题                |
| color       | 否   | string | 主题色              |
| use         | 是   | bool   | 是否启用            |
| partitionId | 是   | int    | 分区ID              |

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

## 获取轮播图列表

#### 请求URL
- `/api/v1/carousel/getCarouselList `
  
#### 请求方式
- POST 

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
        "id": 1,
        "img": "",
        "url": "",
        "title": "",
        "color": "",
        "use": true,
        "partitionId": 1,
        "createdAt": "2021-08-01T10:39:25Z"
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型   | 说明           |
| :----- | :----- | -------------- |
| total  | int    | 数量           |
| list   | object | 轮播图信息数组 |


##### 轮播图信息
| 参数名      | 类型   | 说明               |
| :---------- | :----- | ------------------ |
| id          | int    | 轮播图id           |
| img         | string | 图片地址           |
| url         | string | 单击图片前往的地址 |
| title       | string | 标题               |
| color       | string | 主题色             |
| use         | bool   | 是否启用           |
| partitionId | int    | 分区ID             |
| createdAt   | time   | 发布时间           |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 编辑API

#### 请求URL
- `/api/v1/carousel/editCarousel`
  
#### 请求方式
- PUT 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名      | 必选 | 类型   | 说明                |
| :---------- | :--- | :----- | ------------------- |
| id          | 是   | int    | 轮播图ID            |
| img         | 是   | string | 轮播图url           |
| url         | 否   | string | 单击轮播图打开的url |
| title       | 否   | string | 标题                |
| color       | 否   | string | 主题色              |
| use         | 是   | bool   | 是否启用            |
| partitionId | 是   | int    | 分区ID              |

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
    
## 删除轮播图

#### 请求URL
- `/api/v1/carousel/deleteCarousel/轮播图ID `
  
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