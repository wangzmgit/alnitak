# 评论点赞 API 文档

## 概述

新增评论点赞功能，支持视频评论和文章评论的点赞/取消点赞，以及查询点赞状态。
点赞评论后会发送消息通知，被点赞的用户可以在"收到的赞"中查看。

## API 接口列表

### 视频评论点赞

| 方法 | 路径 | 描述 | 鉴权 |
|------|------|------|------|
| POST | `/api/v1/comment/video/like/:id` | 点赞视频评论 | 需要 |
| DELETE | `/api/v1/comment/video/like/:id` | 取消点赞视频评论 | 需要 |
| GET | `/api/v1/comment/video/getLikeStatus` | 获取视频评论点赞状态 | 需要 |

### 文章评论点赞

| 方法 | 路径 | 描述 | 鉴权 |
|------|------|------|------|
| POST | `/api/v1/comment/article/like/:id` | 点赞文章评论 | 需要 |
| DELETE | `/api/v1/comment/article/like/:id` | 取消点赞文章评论 | 需要 |
| GET | `/api/v1/comment/article/getLikeStatus` | 获取文章评论点赞状态 | 需要 |

---

## 接口详情

### 1. 点赞评论

**请求**
```
POST /api/v1/comment/video/like/:id
POST /api/v1/comment/article/like/:id
```

**参数**
- `:id` - 评论ID（路径参数）

**响应**
```json
{
  "code": 200,
  "msg": "success"
}
```

**错误码**
- `code: 400` - 评论不存在
- `code: 400` - 已点赞（重复点赞）

---

### 2. 取消点赞评论

**请求**
```
DELETE /api/v1/comment/video/like/:id
DELETE /api/v1/comment/article/like/:id
```

**参数**
- `:id` - 评论ID（路径参数）

**响应**
```json
{
  "code": 200,
  "msg": "success"
}
```

**错误码**
- `code: 400` - 评论不存在
- `code: 400` - 未点赞（重复取消）

---

### 3. 获取评论点赞状态

**请求**
```
GET /api/v1/comment/video/getLikeStatus?commentId=xxx
GET /api/v1/comment/article/getLikeStatus?commentId=xxx
```

**参数**
- `commentId` - 评论ID（查询参数，必须）

**响应**
```json
{
  "code": 200,
  "data": {
    "liked": true
  }
}
```

**响应字段说明**
| 字段 | 类型 | 描述 |
|------|------|------|
| liked | bool | 当前用户是否已点赞 |

---

## 数据库变更

### comment 表新增字段

```sql
ALTER TABLE comment ADD COLUMN likes INT DEFAULT 0;
```

### 新建 comment_like 表

自动通过 GORM AutoMigrate 创建：

| 字段 | 类型 | 描述 |
|------|------|------|
| id | uint | 主键 |
| comment_id | uint | 评论ID |
| uid | uint | 用户ID |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 删除时间 |

---

## 前端调用示例

### JavaScript (fetch)

```javascript
// 点赞评论
async function likeComment(commentId, type = 'video') {
  const response = await fetch(`/api/v1/comment/${type}/like/${commentId}`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
}

// 取消点赞
async function unlikeComment(commentId, type = 'video') {
  const response = await fetch(`/api/v1/comment/${type}/like/${commentId}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
}

// 查询点赞状态
async function getLikeStatus(commentId, type = 'video') {
  const response = await fetch(`/api/v1/comment/${type}/getLikeStatus?commentId=${commentId}`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
}
```

### 响应示例

```javascript
// 点赞成功
{ code: 200, msg: "success" }

// 取消点赞成功
{ code: 200, msg: "success" }

// 查询点赞状态
{ code: 200, data: { liked: true } }

// 错误示例
{ code: 400, msg: "已点赞" }
```

---

## 注意事项

1. 所有接口都需要登录认证
2. 评论的点赞数 `likes` 字段会实时更新
3. 同一用户对同一评论只能点赞一次
4. 取消点赞后可以再次点赞
5. **点赞评论会自动发送消息通知**，被点赞的用户可在"收到的赞"中查看

---

## 消息通知

点赞评论后，被点赞的用户会在"收到的赞"中看到通知。

### 获取点赞消息 API

**请求**
```
GET /api/v1/message/getLikeMsg?page=1&pageSize=10
```

**响应**
```json
{
  "code": 200,
  "data": {
    "total": 1,
    "messages": [
      {
        "id": 1,
        "cid": 123,
        "commentId": 456,
        "sid": 789,
        "type": 3,
        "parentType": 0,
        "createdAt": "2024-01-01T12:00:00Z",
        "user": {
          "uid": 789,
          "username": "点赞者用户名",
          "avatar": "头像URL"
        },
        "comment": {
          "id": 456,
          "content": "评论内容",
          "type": 0,
          "createdAt": "2024-01-01T11:00:00Z"
        }
      }
    ]
  }
}
```

**响应字段说明**
| 字段 | 类型 | 描述 |
|------|------|------|
| type | int | 类型：0视频点赞、1文章点赞、**3评论点赞** |
| parentType | int | 评论所属：0视频、1文章 |
| commentId | uint | 评论ID |
| comment | object | 评论信息 |

**type 字段说明**
| 值 | 描述 |
|----|------|
| 0 | 视频点赞 |
| 1 | 文章点赞 |
| 2 | 合集点赞 |
| 3 | 评论点赞 |