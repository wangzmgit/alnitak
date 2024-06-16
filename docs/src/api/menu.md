# 菜单管理相关接口

## 获取菜单树

#### 请求URL
- `/api/v1/menu/getMenuTree `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "menus": [
      {
        "id": 2,
        "name": "",
        "path": "",
        "desc": "",
        "sort": 2,
        "component": "",
        "meta": {
          "title": "",
          "keepAlive": false,
          "icon": "",
          "hidden": false,
        },
        "children": [],
        "parentId": 0,
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名    | 类型   | 说明       |
| :-------- | :----- | ---------- |
| id        | int    | ID         |
| name      | string | 菜单名称   |
| path      | string | 菜单路径   |
| desc      | string | 简介       |
| sort      | int    | 排序       |
| component | string | 文件路径   |
| meta      | Meta   | -          |
| children  | []Menu | 子菜单     |
| parentId  | int    | 所属菜单ID |

##### Meta
| 参数名    | 类型   | 说明     |
| :-------- | :----- | -------- |
| title     | string | 标题     |
| keepAlive | bool   | 是否缓存 |
| icon      | string | 菜单图标 |
| hidden    | bool   | 是否隐藏 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 获取用户菜单树

#### 请求URL
- `/api/v1/menu/getUserMenu `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "menus": [
      {
        "id": 2,
        "name": "",
        "path": "",
        "desc": "",
        "sort": 2,
        "component": "",
        "meta": {
          "title": "",
          "keepAlive": false,
          "icon": "",
          "hidden": false,
        },
        "children": [],
        "parentId": 0,
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名    | 类型   | 说明       |
| :-------- | :----- | ---------- |
| id        | int    | ID         |
| name      | string | 菜单名称   |
| path      | string | 菜单路径   |
| desc      | string | 简介       |
| sort      | int    | 排序       |
| component | string | 文件路径   |
| meta      | Meta   | -          |
| children  | []Menu | 子菜单     |
| parentId  | int    | 所属菜单ID |

##### Meta
| 参数名    | 类型   | 说明     |
| :-------- | :----- | -------- |
| title     | string | 标题     |
| keepAlive | bool   | 是否缓存 |
| icon      | string | 菜单图标 |
| hidden    | bool   | 是否隐藏 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 新增菜单

#### 请求URL
- `/api/v1/menu/addMenu`
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`

#### 参数

| 参数名    | 必选 | 类型   | 说明       |
| :-------- | :--- | :----- | ---------- |
| name      | 是   | string | 菜单名称   |
| path      | 是   | string | 菜单路径   |
| component | 是   | string | 文件路径   |
| desc      | 否   | string | 简介       |
| sort      | 否   | int    | 排序       |
| parentId  | 否   | int    | 所属菜单ID |
| title     | 是   | string | 标题       |
| keepAlive | 否   | bool   | 是否缓存   |
| icon      | 是   | string | 菜单图标   |
| hidden    | 否   | bool   | 是否隐藏   |

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

## 编辑菜单

#### 请求URL
- `/api/v1/menu/editMenu`
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`

#### 参数

| 参数名    | 必选 | 类型   | 说明       |
| :-------- | :--- | :----- | ---------- |
| id        | 是   | int    | 菜单ID     |
| name      | 是   | string | 菜单名称   |
| path      | 是   | string | 菜单路径   |
| component | 是   | string | 文件路径   |
| desc      | 否   | string | 简介       |
| sort      | 否   | int    | 排序       |
| parentId  | 否   | int    | 所属菜单ID |
| title     | 是   | string | 标题       |
| keepAlive | 否   | bool   | 是否缓存   |
| icon      | 是   | string | 菜单图标   |
| hidden    | 否   | bool   | 是否隐藏   |

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

## 删除菜单

#### 请求URL
- `/api/v1/menu/deleteMenu/菜单ID`
  
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

## 获取角色菜单

#### 请求URL
- `/api/v1/menu/getRoleMenu?code=角色代码 `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

#### 返回示例 
``` json
{
  "code": 200,
  "data": {
    "menus": [
      {
        "id": 2,
        "name": "",
        "path": "",
        "desc": "",
        "sort": 2,
        "component": "",
        "meta": {
          "title": "",
          "keepAlive": false,
          "icon": "",
          "hidden": false,
        },
        "children": [],
        "parentId": 0,
      },
    ]
  },
  "msg": "ok"
}
```

#### 返回参数说明 
| 参数名    | 类型   | 说明       |
| :-------- | :----- | ---------- |
| id        | int    | ID         |
| name      | string | 菜单名称   |
| path      | string | 菜单路径   |
| desc      | string | 简介       |
| sort      | int    | 排序       |
| component | string | 文件路径   |
| meta      | Meta   | -          |
| children  | []Menu | 子菜单     |
| parentId  | int    | 所属菜单ID |

##### Meta
| 参数名    | 类型   | 说明     |
| :-------- | :----- | -------- |
| title     | string | 标题     |
| keepAlive | bool   | 是否缓存 |
| icon      | string | 菜单图标 |
| hidden    | bool   | 是否隐藏 |

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 编辑角色菜单

#### 请求URL
- `/api/v1/menu/editRoleMenu`
  
#### 请求方式
- PUT 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名    | 必选 | 类型  | 说明           |
| :-------- | :--- | :---- | -------------- |
| id        | 是   | int   | 角色ID         |
| menuIds   | 是   | []int | 菜单ID数组     |
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
