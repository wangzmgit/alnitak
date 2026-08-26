# 用户认证模块开发文档

## 概述

本模块参考哔哩哔哩 UP主认证架构设计，实现了灵活的认证类型管理系统。

### 核心特性
- 认证类型可配置（数据库管理，无需改代码）
- 支持多种分类：创作者、身份、企业、其他
- 扩展字段 JSON 化，支持不同认证类型自定义数据
- 用户可拥有多种认证

---

## 数据库结构

### 表1: auth_type (认证类型配置表)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| code | varchar(32) | 认证类型唯一标识，如：up, musician |
| name | varchar(32) | 显示名称，如：UP主认证 |
| category | varchar(32) | 分类：creator/identity/enterprise/other |
| desc | varchar(255) | 描述 |
| icon | varchar(255) | 图标URL |
| color | varchar(16) | 颜色代码，如：#FB7299 |
| priority | int | 显示顺序，越大越靠前 |
| is_enable | bool | 是否启用 |
| required_fields | text | 必填字段JSON，如：["name","id_card"] |
| created_at, updated_at | timestamp | 时间戳 |

### 表2: user_auth (用户认证记录表)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| uid | uint | 用户ID |
| auth_type_id | uint | 关联认证类型ID |
| auth_type_code | varchar(32) | 认证类型code |
| title | varchar(64) | 认证头衔，如：知名UP主 |
| desc | varchar(255) | 认证描述 |
| extra_data | text | 扩展数据JSON |
| created_at, updated_at | timestamp | 时间戳 |

---

## API 接口

### 公开接口（无需登录）

#### 1. 获取认证类型列表
```
GET /api/v1/auth/type/list?category=creator
```
参数：
- category（可选）：creator/identity/enterprise/other

响应：
```json
{
  "code": 200,
  "data": {
    "list": [
      {"id":1,"code":"up","name":"UP主认证","category":"creator","color":"#FB7299",...}
    ]
  }
}
```

#### 2. 获取用户认证列表
```
GET /api/v1/auth/user/list?uid=1
```
参数：
- uid（必填）：用户ID

#### 3. 获取用户主要认证
```
GET /api/v1/auth/user/primary?uid=1
```
用于前端展示认证标识

#### 4. 获取指定用户的认证信息
```
GET /api/v1/auth/user/1/auth
```

### 管理接口（需登录）

#### 认证类型管理

```
POST   /api/v1/auth/type/add     添加认证类型
PUT    /api/v1/auth/type/edit    编辑认证类型
DELETE /api/v1/auth/type/:id     删除认证类型
GET    /api/v1/auth/type/all     获取所有认证类型（分页）
GET    /api/v1/auth/type/:id     获取认证类型详情
```

#### 用户认证管理

```
POST   /api/v1/auth/user/add     添加用户认证
PUT    /api/v1/auth/user/edit    编辑用户认证
DELETE /api/v1/auth/user         删除用户认证（需传 id 和 uid）
GET    /api/v1/auth/user/all     获取用户认证列表（带用户信息）
GET    /api/v1/auth/user/:id     获取用户认证详情
```

---

## 默认认证类型

系统初始化时会自动创建以下默认认证类型：

| code | name | category | color | priority |
|------|------|----------|-------|----------|
| up | UP主认证 | creator | #FB7299 | 100 |
| musician | 音乐人认证 | creator | #23ADE5 | 90 |
| artist | 画师认证 | creator | #F25F8D | 90 |
| vlog | Vlog认证 | creator | #FF6B6B | 90 |
| enterprise | 企业认证 | enterprise | #FFD700 | 80 |
| media | 媒体认证 | other | #4A90D9 | 70 |
| government | 政府认证 | other | #7CB342 | 70 |

---

## 前端使用示例

### 获取用户认证信息并显示

```javascript
// 获取用户认证
const res = await axios.get('/api/v1/auth/user/primary?uid=1')
const auth = res.data.data.auth

if (auth) {
  // 显示认证标识
  console.log(`
    <span style="color: ${auth.authTypeColor}">
      ${auth.authTypeName} · ${auth.title}
    </span>
  `)
}
```

### 显示用户所有认证

```javascript
const res = await axios.get('/api/v1/auth/user/list?uid=1')
const list = res.data.data.list

list.forEach(auth => {
  console.log(`
    ${auth.authTypeIcon ? `<img src="${auth.authTypeIcon}">` : ''}
    <span style="color: ${auth.authTypeColor}">${auth.authTypeName}</span>
    <span>${auth.title}</span>
  `)
})
```

---

## 扩展认证类型

如需添加新的认证类型，只需在数据库插入记录或调用管理接口：

```json
POST /api/v1/auth/type/add
{
  "code": "coder",
  "name": "程序员认证",
  "category": "creator",
  "desc": "编程领域创作者",
  "color": "#00FF00",
  "priority": 85,
  "isEnable": true,
  "requiredFields": "[\"github\"]"
}
```

---

## 目录结构

```
internal/
├── domain/
│   ├── model/
│   │   ├── auth_type.go      # 认证类型模型
│   │   └── user_auth.go      # 用户认证模型
│   ├── dto/
│   │   └── auth.go           # 请求参数
│   └── vo/
│       └── auth.go           # 响应结构
├── service/
│   └── auth.go               # 业务逻辑
├── api/v1/
│   └── auth.go               # API处理
├── routes/
│   └── auth_router.go        # 路由注册
└── initialize/
    └── tables.go             # 数据库初始化
```

---

## 后续开发建议

1. **前端申请入口**：用户申请认证的表单界面
2. **审核流程**：认证申请需后台审核
3. **认证到期**：支持认证有效期设置
4. **多认证展示**：用户可拥有多个认证，前端显示全部
5. **认证徽章**：前端展示认证图标组件

---

## 删除资源接口（含弹幕选项）

### 删除视频资源（支持选择是否同时删除弹幕）

```
DELETE /api/v1/resource/deleteResource/:id?deleteDanmaku=true
```

**参数：**
- `id` - 资源ID（路径参数）
- `deleteDanmaku` - 是否同时删除弹幕（可选，默认 false）
  - `true` - 删除资源同时删除该资源关联的所有弹幕
  - `false` - 仅删除资源，保留弹幕（弹幕绑定的是资源 shortID，不会丢失）

**示例：**

```bash
# 删除分P，保留弹幕（默认）
curl -X DELETE "https://domain/api/v1/resource/123"

# 删除分P，同时删除弹幕
curl -X DELETE "https://domain/api/v1/resource/123?deleteDanmaku=true"
```

**使用场景：**
- 用户删除分P后重新上传不同内容 → 建议勾选"同时删除弹幕"
- 用户删除分P后仅修改字幕重新上传 → 建议不删除弹幕