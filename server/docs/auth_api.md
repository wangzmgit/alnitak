# 用户认证凭证 API

## 公开接口（无需登录）

### 注册
```
POST /api/v1/auth/register
```

### 登录（密码）
```
POST /api/v1/auth/login
```

### 登录（邮箱验证码）
```
POST /api/v1/auth/login/email
```

### 获取当前会话用户
```
GET /api/v1/auth/me
```
支持 Authorization header 或 HttpOnly refresh_token Cookie

### 刷新 Token
```
POST /api/v1/auth/updateToken
```
通过 HttpOnly refresh_token Cookie 自动刷新

### 退出登录
```
POST /api/v1/auth/logout
```
失效 refresh_token + 将当前 accessToken 加入黑名单

### 重置密码检查
```
POST /api/v1/auth/resetpwdCheck
```
忘记密码场景：滑块验证 → 设置检查状态

### 重置密码
```
POST /api/v1/auth/modifyPwd
```
忘记密码场景：邮箱验证码 + 新密码

## 登录态接口（需 Authorization header）

### 修改密码
```
POST /api/v1/auth/changePassword
Authorization: <accessToken>
Content-Type: application/json

{
  "oldPassword": "当前密码",
  "newPassword": "新密码（至少6位）"
}
```
校验旧密码 → 更新为新密码 → 清除所有 refreshToken（其他会话需重新登录）
