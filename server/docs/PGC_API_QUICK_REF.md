# PGC API 快速参考

**基础路径**: `/api/v1/pgc`

## 认证说明
- 无需认证: PGC内容查询、剧集查询
- 需要认证: 创建、更新、删除PGC，添加/删除剧集
- 管理员权限: 用户组管理

---

## 📋 PGC内容 API

### 创建PGC
```
POST /api/v1/pgc/create
```
**认证**: 是

### 更新PGC
```
PUT /api/v1/pgc/update
```
**认证**: 是

### 删除PGC
```
DELETE /api/v1/pgc/:pgc_id
```
**认证**: 是

### 获取PGC列表
```
GET /api/v1/pgc/list?page=1&page_size=20&pgc_type=1&is_ongoing=true&keyword=xxx
```
**认证**: 否

### 获取PGC详情
```
GET /api/v1/pgc/detail?pgc_id=xxx
```
**认证**: 否

### 搜索PGC
```
GET /api/v1/pgc/search?keyword=xxx&pgc_type=1&page=1&page_size=20
```
**认证**: 否

### 按类型获取PGC
```
GET /api/v1/pgc/type/:type?page=1&page_size=20
```
**认证**: 否

### 获取连载中PGC
```
GET /api/v1/pgc/ongoing?page=1&page_size=20
```
**认证**: 否

### 获取推荐PGC
```
GET /api/v1/pgc/recommended?limit=10
```
**认证**: 否

### 获取PGC详情及剧集
```
GET /api/v1/pgc/detail-with-episodes?pgc_id=xxx
```
**认证**: 否

---

## 🎬 剧集管理 API

### 获取剧集列表
```
GET /api/v1/pgc/:pgc_id/episodes?page=1&page_size=20
```
**认证**: 否

### 添加剧集
```
POST /api/v1/pgc/:pgc_id/episodes/add
```
**认证**: 是  
**说明**: `vid` 可选；为 `0` 或不传时表示先创建占位剧集，后续调用「绑定视频」接口补充。

### 绑定剧集视频（占位 -> 已绑定）
```
PUT /api/v1/pgc/:pgc_id/episodes/:id/bind
```
**认证**: 是  
**Body**:
```json
{ "vid": 123456, "duration": 1440.0, "publish_time": "2024-01-01 10:00:00" }
```
`duration`、`publish_time` 可选；未传 `publish_time` 时默认取视频创建时间。

### 删除剧集
```
DELETE /api/v1/pgc/:pgc_id/episodes/:id
```
**认证**: 是

---

## 👥 用户组管理 API (管理员)

### 添加用户组
```
POST /api/v1/pgc/admin/user-group/add
```
**认证**: 管理员

### 移除用户组
```
DELETE /api/v1/pgc/admin/user-group/remove
```
**认证**: 管理员

### 获取用户组列表
```
GET /api/v1/pgc/admin/user-group/list?page=1&page_size=20&group_type=3
```
**认证**: 管理员

---

## 🔐 用户信息 API

### 获取当前用户组
```
GET /api/v1/pgc/user-group/mine
```
**认证**: 是

---

## 📦 常量说明

### PGC类型 (pgc_type)
- `1`: 国创(CN)
- `2`: 日创(JP)
- `3`: 纪录片
- `4`: 电影
- `5`: 电视剧

