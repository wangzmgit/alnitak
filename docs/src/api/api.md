# API管理相关接口

<!-- ************************ 分隔符 ************************ -->

## 获取API列表

#### 请求URL
- `/api/v1/api/getApiList `
  
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
        "id": 2,
        "method": "GET",
        "path": "",
        "category": "",
        "desc": "",
        "createdAt": "2021-07-16T08:49:54Z",
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型   | 说明        |
| :----- | :----- | ----------- |
| total  | int    | 数量        |
| list   | object | API信息数组 |


##### API信息
| 参数名    | 类型   | 说明     |
| :-------- | :----- | -------- |
| id        | int    | ID       |
| method    | string | 请求方法 |
| path      | string | 请求路径 |
| category  | string | 分组     |
| desc      | string | 简介     |
| createdAt | string | 上传时间 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取全部API列表

#### 请求URL
- `/api/v1/api/getAllApiList`
  
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
    "list": [
      {
        "id": 2,
        "method": "GET",
        "path": "",
        "category": "",
        "desc": "",
        "createdAt": "2021-07-16T08:49:54Z",
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型   | 说明        |
| :----- | :----- | ----------- |
| total  | int    | 数量        |
| list   | object | API信息数组 |


##### API信息
| 参数名    | 类型   | 说明     |
| :-------- | :----- | -------- |
| id        | int    | ID       |
| method    | string | 请求方法 |
| path      | string | 请求路径 |
| category  | string | 分组     |
| desc      | string | 简介     |
| createdAt | string | 上传时间 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 新增API

#### 请求URL
- `/api/v1/api/addAPI`
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`

#### 参数

| 参数名   | 必选 | 类型   | 说明     |
| :------- | :--- | :----- | -------- |
| method   | 是   | string | 请求方法 |
| path     | 是   | string | 请求路径 |
| category | 是   | string | 分组     |
| desc     | 是   | string | 简介     |


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

## 编辑API

#### 请求URL
- `/api/v1/api/editApi`
  
#### 请求方式
- PUT 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名   | 必选 | 类型   | 说明     |
| :------- | :--- | :----- | -------- |
| id       | 是   | int    | ID       |
| method   | 是   | string | 请求方法 |
| path     | 是   | string | 请求路径 |
| category | 是   | string | 分组     |
| desc     | 是   | string | 简介     |

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

## 删除API

#### 请求URL
- `/api/v1/api/deleteApi/apiID`
  
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

## 获取角色API

#### 请求URL
- `/api/v1/api/getRoleApi?code=角色代码 `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 2,
        "method": "GET",
        "path": "",
        "category": "",
        "desc": "",
        "createdAt": "2021-07-16T08:49:54Z",
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型   | 说明        |
| :----- | :----- | ----------- |
| list   | object | API信息数组 |

##### API信息
| 参数名    | 类型   | 说明     |
| :-------- | :----- | -------- |
| id        | int    | ID       |
| method    | string | 请求方法 |
| path      | string | 请求路径 |
| category  | string | 分组     |
| desc      | string | 简介     |
| createdAt | string | 上传时间 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 编辑角色API

#### 请求URL
- `/api/v1/api/editRoleApi`
  
#### 请求方式
- PUT 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名    | 必选 | 类型  | 说明           |
| :-------- | :--- | :---- | -------------- |
| id        | 是   | int   | 角色ID         |
| addIds    | 是   | []int | 添加API ID数组 |
| removeIds | 是   | []int | 移除API ID数组 |

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
