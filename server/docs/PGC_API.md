# PGC (专业生成内容) API 文档

## 概述

PGC (Professional Generated Content) 系统提供专业视频内容管理功能，类似于B站的番剧、纪录片、电影等内容管理。

**基础路径**: `/api/v1/pgc`

**认证方式**: 所有需要认证的接口都需要在请求头中携带 JWT Token
```
Authorization: Bearer <token>
```

---

## 一、PGC内容管理 API

### 1.1 创建PGC内容

创建新的PGC内容，包括标题、封面、简介等信息，以及关联的剧集。

**接口**: `POST /api/v1/pgc/create`

**认证**: 需要登录 + PGC权限

**请求参数**:
```json
{
  "pgc_type": 1,
  "title": "进击的巨人",
  "cover": "https://example.com/cover.jpg",
  "desc": "人类与巨人之间的战争",
  "year": 2013,
  "area": "日本",
  "rating": 9.5,
  "is_ongoing": true,
  "episodes": [
    {
      "episode_number": 1,
      "title": "致两千年后的你",
      "vid": 123456,
      "duration": 1440.5,
      "publish_time": "2024-01-01 10:00:00"
    },
    {
      "episode_number": 2,
      "title": "那一天",
      "vid": 123457,
      "duration": 1440.0,
      "publish_time": "2024-01-08 10:00:00"
    }
  ]
}
```

**参数说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| pgc_type | int | 是 | PGC类型 (1:国创(CN), 2:日创(JP), 3:纪录片, 4:电影, 5:电视剧) |
| title | string | 是 | 标题，最大255字符 |
| cover | string | 是 | 封面图片URL |
| desc | string | 否 | 简介，最大1000字符 |
| year | int | 否 | 年份 |
| area | string | 否 | 地区 (如：日本、美国、中国) |
| rating | float64 | 否 | 评分 (0-10) |
| is_ongoing | bool | 否 | 是否连载中 |
| episodes | array | 是 | 剧集列表，至少一集 |

**episodes 数组说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| episode_number | int | 是 | 剧集序号 (从1开始，不能重复) |
| title | string | 否 | 剧集标题 |
| vid | uint | 是 | 关联的视频ID |
| duration | float64 | 否 | 时长（秒） |
| publish_time | string | 否 | 发布时间 (格式: YYYY-MM-DD HH:mm:ss) |

**成功响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "pgc_id": 1234567890123456789
  }
}
```

**错误响应**:
```json
{
  "code": 4000,
  "msg": "标题不能为空",
  "data": null
}
```

---

### 1.2 更新PGC内容

更新已存在的PGC内容信息。

**接口**: `PUT /api/v1/pgc/update`

**认证**: 需要登录 + PGC权限

**请求参数**:
```json
{
  "pgc_id": 1234567890123456789,
  "title": "进击的巨人 最终季",
  "cover": "https://example.com/new-cover.jpg",
  "desc": "新的简介",
  "year": 2020,
  "area": "日本",
  "rating": 9.8,
  "is_ongoing": false
}
```

**参数说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| pgc_id | uint | 是 | PGC内容ID |
| title | string | 否 | 标题 (最大255字符) |
| cover | string | 否 | 封面URL |
| desc | string | 否 | 简介 |
| year | int | 否 | 年份 |
| area | string | 否 | 地区 |
| rating | float64 | 否 | 评分 (0-10) |
| is_ongoing | bool | 否 | 是否连载中 |

**成功响应**:
```json
{
  "code": 200,
  "msg": "更新成功",
  "data": null
}
```

---

### 1.3 删除PGC内容

删除PGC内容及其关联的所有剧集（物理删除）。

**接口**: `DELETE /api/v1/pgc/:pgc_id`

**认证**: 需要登录 + PGC权限

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| pgc_id | uint | 是 | PGC内容ID |

**成功响应**:
```json
{
  "code": 200,
  "msg": "删除成功",
  "data": null
}
```

---

### 1.4 获取PGC内容列表

分页获取PGC内容列表，支持多种筛选条件。

**接口**: `GET /api/v1/pgc/list`

**认证**: 无需认证

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码，从1开始 |
| page_size | int | 是 | 每页数量，最大100 |
| pgc_type | int | 否 | PGC类型筛选 |
| status | int | 否 | 审核状态筛选 |
| keyword | string | 否 | 关键词搜索 (标题或简介) |
| year | int | 否 | 年份筛选 |
| area | string | 否 | 地区筛选 |
| is_ongoing | bool | 否 | 是否连载中 |

**请求示例**:
```
GET /api/v1/pgc/list?page=1&page_size=20&pgc_type=1&is_ongoing=true
```

**成功响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 100,
    "list": [
      {
        "id": 1,
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z",
        "pgc_id": 1234567890123456789,
        "pgc_type": 1,
        "title": "进击的巨人",
        "cover": "https://example.com/cover.jpg",
        "desc": "人类与巨人之间的战争",
        "year": 2013,
        "area": "日本",
        "rating": 9.5,
        "is_ongoing": true,
        "total_episodes": 12,
        "current_episodes": 12,
        "status": 300,
        "operator_id": 0
      }
    ],
    "page": 1,
    "page_size": 20
  }
}
```

---

### 1.5 获取PGC内容详情

获取指定PGC内容的详细信息。

**接口**: `GET /api/v1/pgc/detail`

**认证**: 无需认证

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| pgc_id | uint | 是 | PGC内容ID |

**请求示例**:
```
GET /api/v1/pgc/detail?pgc_id=1234567890123456789
```

**成功响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "pgc": {
      "id": 1,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z",
      "pgc_id": 1234567890123456789,
      "pgc_type": 1,
      "title": "进击的巨人",
      "cover": "https://example.com/cover.jpg",
      "desc": "人类与巨人之间的战争",
      "year": 2013,
      "area": "日本",
      "rating": 9.5,
      "is_ongoing": true,
      "total_episodes": 12,
      "current_episodes": 12,
      "status": 300,
      "operator_id": 0
    }
  }
}
```

---

### 1.6 搜索PGC内容

根据关键词和类型搜索PGC内容。

**接口**: `GET /api/v1/pgc/search`

**认证**: 无需认证

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 否 | 搜索关键词 (标题或简介) |
| pgc_type | int | 否 | PGC类型筛选 |
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认20 |

**请求示例**:
```
GET /api/v1/pgc/search?keyword=巨人&pgc_type=1&page=1&page_size=20
```

**成功响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 5,
    "list": [...],
    "page": 1,
    "page_size": 20
  }
}
```

---

### 1.7 按类型获取PGC

获取指定类型的PGC内容列表。

**接口**: `GET /api/v1/pgc/type/:type`

**认证**: 无需认证

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | int | 是 | PGC类型 (1:国创(CN), 2:日创(JP), 3:纪录片, 4:电影, 5:电视剧) |

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| page_size | int | 是 | 每页数量 |

**请求示例**:
```
GET /api/v1/pgc/type/1?page=1&page_size=20
```

**成功响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 50,
    "list": [...],
    "page": 1,
    "page_size": 20
  }
}
```

---

### 1.8 获取连载中的PGC

获取所有正在连载的PGC内容。

**接口**: `GET /api/v1/pgc/ongoing`

**认证**: 无需认证

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| page_size | int | 是 | 每页数量 |

**请求示例**:
```
GET /api/v1/pgc/ongoing?page=1&page_size=20
```

**成功响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 25,
    "list": [...],
    "page": 1,
    "page_size": 20
  }
}
```

---

### 1.9 获取推荐PGC

获取推荐的PGC内容列表（按创建时间倒序）。

**接口**: `GET /api/v1/pgc/recommended`

**认证**: 无需认证

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| limit | int | 否 | 返回数量，默认10，最大50 |

**请求示例**:
```
GET /api/v1/pgc/recommended?limit=20
```

**成功响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [...]
  }
}
```

---

### 1.10 获取PGC详情及剧集

获取PGC内容详情及完整的剧集列表。

**接口**: `GET /api/v1/pgc/detail-with-episodes`

**认证**: 无需认证

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| pgc_id | uint | 是 | PGC内容ID |

**请求示例**:
```
GET /api/v1/pgc/detail-with-episodes?pgc_id=1234567890123456789
```

**成功响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "pgc": {
      "id": 1,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z",
      "pgc_id": 1234567890123456789,
      "pgc_type": 1,
      "title": "进击的巨人",
      "cover": "https://example.com/cover.jpg",
      "desc": "人类与巨人之间的战争",
      "year": 2013,
      "area": "日本",
      "rating": 9.5,
      "is_ongoing": true,
      "total_episodes": 12,
      "current_episodes": 12,
      "status": 300,
      "operator_id": 0
    },
    "episodes": [
      {
        "id": 1,
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z",
        "pgc_id": 1234567890123456789,
        "episode_number": 1,
        "title": "致两千年后的你",
        "vid": 123456,
        "duration": 1440.5,
        "status": 0,
        "publish_time": "2024-01-01 10:00:00"
      }
    ]
  }
}
```

---

## 二、PGC剧集管理 API

### 2.1 获取PGC剧集列表

获取指定PGC内容的剧集列表。

**接口**: `GET /api/v1/pgc/:pgc_id/episodes`

**认证**: 无需认证

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| pgc_id | uint | 是 | PGC内容ID |

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| page_size | int | 是 | 每页数量，最大100 |

**请求示例**:
```
GET /api/v1/pgc/1234567890123456789/episodes?page=1&page_size=20
```

**成功响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "episodes": [
      {
        "id": 1,
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z",
        "pgc_id": 1234567890123456789,
        "episode_number": 1,
        "title": "致两千年后的你",
        "vid": 123456,
        "duration": 1440.5,
        "status": 0,
        "publish_time": "2024-01-01 10:00:00"
      }
    ]
  }
}
```

---

### 2.2 添加剧集

向已存在的PGC内容添加新剧集。

**接口**: `POST /api/v1/pgc/:pgc_id/episodes/add`

**认证**: 需要登录 + PGC权限

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| pgc_id | uint | 是 | PGC内容ID |

**请求参数**:
```json
{
  "episode_number": 13,
  "title": "新的开始",
  "vid": 123458,
  "duration": 1440.0,
  "publish_time": "2024-01-15 10:00:00"
}
```

**参数说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| episode_number | int | 是 | 剧集序号 |
| title | string | 否 | 剧集标题 |
| vid | uint | 是 | 关联的视频ID |
| duration | float64 | 否 | 时长（秒） |
| publish_time | string | 否 | 发布时间 |

**成功响应**:
```json
{
  "code": 200,
  "msg": "添加成功",
  "data": null
}
```

---

### 2.3 删除剧集

删除PGC内容的指定剧集。

**接口**: `DELETE /api/v1/pgc/:pgc_id/episodes/:id`

**认证**: 需要登录 + PGC权限

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| pgc_id | uint | 是 | PGC内容ID |
| id | uint | 是 | 剧集ID |

**请求示例**:
```
DELETE /api/v1/pgc/1234567890123456789/episodes/1
```

**成功响应**:
```json
{
  "code": 200,
  "msg": "删除成功",
  "data": null
}
```

---

## 三、后台管理 API（Admin）

以下接口仅需登录（Auth 中间件），不需要 PGC 创作者权限，供后台管理面板使用。

### 3.1 获取PGC管理列表

分页获取所有 PGC 内容，支持按状态、类型、关键词筛选。用于后台内容管理页面。

**接口**: `POST /api/v1/pgc/getManageList`

**认证**: 需要登录

**请求参数**:
```json
{
  "page": 1,
  "pageSize": 10,
  "pgc_type": 1,
  "status": 300,
  "keyword": "巨人"
}
```

**参数说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码，从1开始 |
| pageSize | int | 是 | 每页数量，最大100 |
| pgc_type | int | 否 | PGC类型筛选 (1:国创, 2:日漫, 3:纪录片, 4:电影, 5:电视剧) |
| status | int | 否 | 状态筛选 (-1:已下架, 0:草稿, 100:已提交, 200:审核中, 300:已通过, 400:已驳回) |
| keyword | string | 否 | 关键词搜索（标题或简介） |

**成功响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 50,
    "list": [
      {
        "id": 1,
        "pgc_id": "1234567890123456789",
        "season_id": "1234567890123456789",
        "media_id": "1234567890123456789",
        "pgc_type": 1,
        "title": "进击的巨人",
        "cover": "cover.jpg",
        "desc": "简介",
        "year": 2013,
        "area": "日本",
        "rating": 9.5,
        "is_ongoing": false,
        "total_episodes": 12,
        "current_episodes": 12,
        "status": 300,
        "operator_id": 0,
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z"
      }
    ],
    "page": 1,
    "page_size": 10
  }
}
```

---

### 3.2 管理员修改PGC状态

管理员上架/下架 PGC 内容。仅支持在 **已通过(300)** 和 **已下架(-1)** 之间切换。

**接口**: `POST /api/v1/pgc/adminUpdateStatus`

**认证**: 需要登录

**请求参数**:
```json
{
  "pgc_id": "1234567890123456789",
  "status": -1
}
```

**参数说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| pgc_id | string | 是 | PGC内容ID（字符串形式，避免JS精度丢失） |
| status | int | 是 | 目标状态，仅允许 300(上架) 或 -1(下架) |

**成功响应**:
```json
{
  "code": 200,
  "msg": "操作成功",
  "data": null
}
```

**错误响应**:
```json
{
  "code": 4000,
  "msg": "无效的状态值",
  "data": null
}
```

---

### 3.3 管理员删除PGC

删除指定 PGC 内容及其关联的所有剧集。此操作不可恢复。

**接口**: `DELETE /api/v1/pgc/adminDelete/:pgc_id`

**认证**: 需要登录

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| pgc_id | uint | 是 | PGC内容ID |

**请求示例**:
```
DELETE /api/v1/pgc/adminDelete/1234567890123456789
```

**成功响应**:
```json
{
  "code": 200,
  "msg": "删除成功",
  "data": null
}
```

---

## 四、用户组管理

分页获取用户组列表。

**接口**: `GET /api/v1/pgc/admin/user-group/list`

**认证**: 需要登录 + 管理员权限

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| page_size | int | 是 | 每页数量，最大100 |
| group_type | int | 否 | 用户组类型筛选 |

**请求示例**:
```
GET /api/v1/pgc/admin/user-group/list?page=1&page_size=20&group_type=3
```

**成功响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 50,
    "list": [
      {
        "id": 1,
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z",
        "uid": 123456,
        "group_type": 3,
        "remark": "PGC内容创作者"
      }
    ],
    "page": 1,
    "page_size": 20
  }
}
```

---


获取当前登录用户所属的用户组列表。

**接口**: `GET /api/v1/pgc/user-group/mine`

**认证**: 需要登录

**成功响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "user_groups": [
      {
        "id": 1,
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z",
        "uid": 123456,
        "group_type": 3,
        "remark": "PGC内容创作者"
      }
    ]
  }
}
```

---

## 五、数据模型

### 5.1 PGC内容 (PGCContent)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键ID |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| pgc_id | uint | PGC内容ID (雪花ID) |
| pgc_type | int | PGC类型 |
| title | string | 标题 |
| cover | string | 封面URL |
| desc | string | 简介 |
| year | int | 年份 |
| area | string | 地区 |
| rating | float64 | 评分 (0-10) |
| is_ongoing | bool | 是否连载中 |
| total_episodes | int | 总集数 |
| current_episodes | int | 已更新集数 |
| status | int | 审核状态 |
| operator_id | uint | 运营人员ID |

---

### 5.2 PGC剧集 (PGCEpisode)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键ID |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| pgc_id | uint | PGC内容ID |
| episode_number | int | 剧集序号 |
| title | string | 剧集标题 |
| vid | uint | 关联视频ID |
| duration | float64 | 时长（秒） |
| status | int | 状态 |
| publish_time | string | 发布时间 |

---

