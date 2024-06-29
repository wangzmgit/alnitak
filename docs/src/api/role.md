# 角色管理相关接口

## 获取个人角色信息

#### 请求URL
- `/api/v1/role/getRoleInfo `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "role": {
      "id": 2,
      "name": "",
      "code": "",
      "desc": "",
      "homePage": "",
      "createdAt": "2021-07-16T08:49:54Z",
    },
  },
  "msg": "ok"
}
```
#### 返回参数说明 
| 参数名    | 类型   | 说明     |
| :-------- | :----- | -------- |
| id        | int    | 角色ID   |
| name      | string | 角色名   |
| code      | string | 角色代码 |
| desc      | string | 简介     |
| homePage  | string | 角色首页 |
| createdAt | string | 创建时间 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取角色列表

#### 请求URL
- `/api/v1/role/getRoleList `
  
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
        "name": "",
        "code": "",
        "desc": "",
        "homePage": "",
        "createdAt": "2021-07-16T08:49:54Z",
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
| list   | object | 角色信息数组 |

##### 角色信息
| 参数名    | 类型   | 说明     |
| :-------- | :----- | -------- |
| id        | int    | 角色ID   |
| name      | string | 角色名   |
| code      | string | 角色代码 |
| desc      | string | 简介     |
| homePage  | string | 角色首页 |
| createdAt | string | 创建时间 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取全部角色列表

#### 请求URL
- `/api/v1/role/getAllRoleList`
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`
  
#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "roles": [
      {
        "name": "",
        "code": "",
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名 | 类型   | 说明     |
| :----- | :----- | -------- |
| name   | string | 角色名   |
| code   | string | 角色代码 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 新增角色

#### 请求URL
- `/api/v1/role/addRole`
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`

#### 参数

| 参数名 | 必选 | 类型   | 说明     |
| :----- | :--- | :----- | -------- |
| name   | 是   | string | 角色名   |
| code   | 是   | string | 角色代码 |
| desc   | 否   | string | 简介     |


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

## 编辑角色

#### 请求URL
- `/api/v1/role/editRole`
  
#### 请求方式
- PUT 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型   | 说明   |
| :----- | :--- | :----- | ------ |
| id     | 是   | int    | 角色ID |
| name   | 是   | string | 角色名 |
| desc   | 否   | string | 简介   |

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

## 删除角色

#### 请求URL
- `/api/v1/role/deleteRole/角色ID`
  
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

## 编辑角色首页

#### 请求URL
- `/api/v1/role/editRoleHome`
  
#### 请求方式
- PUT 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型   | 说明     |
| :----- | :--- | :----- | -------- |
| id     | 是   | int    | 角色ID   |
| home   | 是   | string | 角色首页 |

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
