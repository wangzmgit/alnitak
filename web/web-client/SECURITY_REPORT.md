# 安全审计报告

> 审计日期: 2026-06-19
> 审计范围: 前端 `web-client` (Nuxt 3 + Vue 3) / 后端 `server` (Go + Gin)

---

## 目录

1. [严重漏洞](#1-严重漏洞)
2. [高危漏洞](#2-高危漏洞)
3. [中危漏洞](#3-中危漏洞)
4. [低危/信息](#4-低危信息)
5. [依赖版本风险](#5-依赖版本风险)
6. [综合修复优先级](#6-综合修复优先级)

---

# 1. 严重漏洞

## C1. 配置文件明文存储全部基础设施凭据

| 属性 | 值 |
|------|-----|
| **文件** | `conf/application.prod.yaml`, `conf/application.dev.yaml` |
| **CWE** | CWE-312: 敏感信息明文存储 |
| **CVSS** | 9.1 (AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:N) |

### 泄露的凭据

| 凭据类型 | 值（部分脱敏） | 风险 |
|---------|---------------|------|
| JWT Access Secret | `3785657888448008` | 伪造任意用户 token |
| JWT Refresh Secret | `3785657888448008` | 与 access 相同，加倍风险 |
| User ID Salt | `3785657888448008` | 用户 ID 可预测 |
| MySQL 密码 | `12345678.aB` | 数据库完全访问 |
| Cloudflare R2 Key ID | `9fd91d5234991e653c7edc8487b988be` | 对象存储读写 |
| Cloudflare R2 Key Secret | `cfut_j0JKYDskiD...` | 对象存储读写 |
| Minio Key ID | `vEYLY7BwUVOHz3TVOwwY` | 本地对象存储访问 |
| Minio Key Secret | `DbVbRw11W83Hkx8nXzXBU1NcCCwzSOKSXpZYFwB9` | 本地对象存储访问 |
| Resend API Key（dev） | `re_NnWLaaTy_NBLqZWUj73f7DxWUibiLsvsd` | 邮件服务滥用 |
| SMTP 密码 | 空 / 明文 | 邮件服务滥用 |

### 攻击路径

```
攻击者拿到 Git 仓库（内部泄露/公开/CI泄露）
    ↓
读取 conf/application.prod.yaml
    ↓
   ├── 用 R2 Key 直接读取全部用户上传文件（视频、图片）
   ├── 用 MySQL 密码直连数据库，导出全部用户数据（邮箱、密码哈希）
   ├── 用 JWT Secret 伪造任意用户 token，登录任意账号
   └── 用 Resend API Key 滥用邮件服务发垃圾邮件
```

### 修复方案

1. 立即轮换所有凭据（生成新的随机值）
2. 将 `conf/*.yaml` 加入 `.gitignore`
3. 使用环境变量或密钥管理服务（如 Vault/ AWS Secrets Manager）注入敏感配置
4. 检查 Git 历史中是否已有凭据泄露，若已 push 到远程需立即轮换

---

## C2. Access Token 与 Refresh Token 使用完全相同的 JWT 密钥

| 属性 | 值 |
|------|-----|
| **文件** | `conf/application.prod.yaml:47-50`, `pkg/jwt/jwt.go:102-110` |
| **CWE** | CWE-321: 使用硬编码的密码学密钥 |

### 问题代码

```go
// pkg/jwt/jwt.go:102-110
switch claims.TokenType {
case 0: // accessToken
    secret = []byte(global.Config.Security.AccessJwtSecret)
case 1: // refreshToken
    secret = []byte(global.Config.Security.RefreshJwtSecret)  // 和 AccessJwtSecret 相同!
}
```

```yaml
# application.prod.yaml
access_jwt_secret: "3785657888448008"
refresh_jwt_secret: "3785657888448008"   # 完全一样
```

### 攻击路径

```
密码学分析/暴力破解 JWT Secret → 拿到 "3785657888448008"
    ↓
用同一密钥既可以伪造 accessToken（登录任意账号）
    ↓
也可以伪造 refreshToken（绕过 token 刷新机制，长期维持访问）
```

### 修复方案

1. 生成两个独立的 256-bit 随机密钥（`openssl rand -base64 32`）
2. 分别配置 `access_jwt_secret` 和 `refresh_jwt_secret`
3. 考虑使用不同的签名算法（如 access 用 RS256, refresh 用 HS256）

---

## C3. JWT 密钥强度不足（弱密钥）

| 属性 | 值 |
|------|-----|
| **文件** | `conf/application.prod.yaml:47` |
| **CWE** | CWE-1393: 使用过短的密码学密钥 |

### 问题

当前 JWT 密钥 `"3785657888448008"`：
- 仅 16 位纯数字
- 看起来像 Snowflake ID 或时间戳
- 暴力破解复杂度约 10^16 ≈ 2^53，远低于 HS256 推荐的 2^256

现代 GPU（如 RTX 4090）每秒可尝试约 10^9 次 HS256 签名验证：
- 理论穷举时间：10^16 / 10^9 = 10^7 秒 ≈ 115 天（但字典/模式攻击可大幅缩短）

### 修复方案

```bash
# 生成 256-bit (32字节) 随机密钥
openssl rand -base64 32
# 输出示例: f84K3hF8a9sL2xZn5Qr7Yu9W0bN1cM4vE6jH3gI5oP0=
```

将生成的密钥写入配置文件。

---

## C4. 生产环境 CORS 配置为通配符

| 属性 | 值 |
|------|-----|
| **文件** | `conf/application.prod.yaml:2`, `internal/middleware/cors.go:29-32` |
| **CWE** | CWE-942: 跨域权限管理不当 |

### 问题代码

```yaml
# application.prod.yaml
cors:
    allow_origin: '*'
```

```go
// internal/middleware/cors.go
func getAllowOrigin(origin string) string {
    if global.Config.Cors.AllowOrigin == "*" {
        return origin  // 直接将请求中的 Origin 返回，等于允许任意来源
    }
    // ...
}
```

同时前端构建时 `withCredentials: true`（`src/utils/request.ts:53`）。

### 攻击路径

```
攻击者在恶意站点 https://evil.com/ 上构造 JS 请求
    ↓
浏览器发现 CORS 头 Access-Control-Allow-Origin: * 且 withCredentials=true
    ↓
fetch('https://animes.ayypd.cn:9001/api/v1/user/info', { credentials: 'include' })
    ↓
用户浏览器自动带上 Cookie，API 返回用户敏感信息
    ↓
攻击者读取响应，获取用户邮箱、手机、登录态等
```

> ⚠️ **注意**: 前端配置了 `withCredentials: true`，浏览器对 `Access-Control-Allow-Origin: *` 的请求会**拒绝附带凭据**。但当前 CORS 逻辑并不是简单返回 `*`，而是**回传请求的 Origin**（`getAllowOrigin` 在 `*` 时返回请求的 Origin），浏览器将视为特定 Origin 从而允许附带凭据。

### 修复方案

```yaml
# application.prod.yaml
cors:
    allow_origin: "https://animes.ayypd.cn,https://www.ayypd.cn"
```

---

## C5. 前端 Token 存储在 localStorage（XSS 可窃取）

| 属性 | 值 |
|------|-----|
| **文件** | `src/stores/auth-store.ts:19-23`, `src/utils/storage-data.ts` |
| **CWE** | CWE-312: 敏感信息明文存储 |

### 问题代码

```typescript
// src/stores/auth-store.ts:19-23
const saveCredentials = (data) => {
  if (data.token) {
    storageData.set('token', data.token, 60);        // localStorage
    Cookies.set('token', data.token, { expires: 1 }); // Cookie（复制）
  }
  if (data.refreshToken) storageData.set('refreshToken', data.refreshToken, 7 * 24 * 60);
  // ...
};
```

Access token 和 refresh token 都写入 `localStorage`。同时，refresh token 也通过 HttpOnly Cookie 传递（后端 `setRefreshCookie`），但前端的 localStorage 副本完全不设 HttpOnly 保护。

### 攻击路径

```
页面存在任何 XSS 漏洞（评论、弹幕、个人信息等）
    ↓
攻击者执行: localStorage.getItem('token')
    ↓
窃取 accessToken → 立即发起 API 请求冒充用户
    ↓
窃取 refreshToken → 即使 accessToken 过期也能续签
```

### 修复方案

1. **仅使用 HttpOnly Cookie** 存储 refresh token（后端已实现 `setRefreshCookie`）
2. 前端不再将 token 写入 localStorage
3. 请求时依赖浏览器自动携带 Cookie，而非从 localStorage 读取手动设置 Authorization header
4. 或者使用 Web Worker 隔离的 Storage（需要额外架构改动）

---

# 2. 高危漏洞

## H1. 缺乏有效的 Content-Security-Policy

| 属性 | 值 |
|------|-----|
| **文件** | `server/middleware/security-headers.ts`（前端）|
| **CWE** | CWE-1021: 内容安全策略缺失/不当 |

### 问题代码

```typescript
// server/middleware/security-headers.ts:19-41
"script-src 'self' 'unsafe-inline' 'unsafe-eval'",
"style-src 'self' 'unsafe-inline'",
"img-src 'self' data: blob: *",
"connect-src 'self' *",
```

### 问题分析

| 指令 | 问题 |
|------|------|
| `'unsafe-inline'` | 允许任何内联脚本执行，CSP 的主要防护被绕过 |
| `'unsafe-eval'` | 允许 eval/setTimeout(string)/Function() 等执行字符串代码 |
| `img-src *` + `connect-src *` | 允许向任意外部地址发请求，可用于数据外传 |
| `frame-ancestors 'self'`（非 embed） | 这个做得好，可以防点击劫持 |

后端完全没有设置 CSP header。

### 攻击路径

```
攻击者存储 XSS payload（如评论、弹幕）
    ↓
CSP 无法阻止内联脚本执行（'unsafe-inline'）
    ↓
脚本读取 localStorage 中的 token 并外传到攻击者服务器
```

### 修复方案

逐步收紧 CSP：
```typescript
// 第一阶段：移除 'unsafe-eval'（多数现代框架不需要 eval）
"script-src 'self' 'unsafe-inline'",
// 第二阶段：使用 nonce 或 hash 替代 'unsafe-inline'
"script-src 'self' 'nonce-{random}'",
```

---

## H2. 操作日志明文记录请求体（含敏感字段）

| 属性 | 值 |
|------|-----|
| **文件** | `internal/middleware/operation.go:77` |
| **CWE** | CWE-532: 日志中包含敏感信息 |

### 问题代码

```go
// internal/middleware/operation.go:77
record := model.Operate{
    Ip:     c.ClientIP(),
    Method: c.Request.Method,
    Path:   c.Request.URL.Path,
    Agent:  c.Request.UserAgent(),
    Body:   string(body),   // ← 完整请求体原始内容
    UserID: userId,
}
```

### 风险

登录（`POST /api/v1/auth/login`）、注册、修改密码等请求的 body 中可能包含密码、验证码等敏感字段，这些会被完整记录到数据库的 `operate` 表中，且无过期清理机制。

### 修复方案

```go
// 对敏感字段进行脱敏
func sanitizeBody(body []byte, path string) string {
    if strings.Contains(path, "auth/login") || 
       strings.Contains(path, "auth/register") ||
       strings.Contains(path, "auth/modifyPwd") {
        return "[REDACTED]"
    }
    return string(body)
}
```

---

## H3. HTTP/HTTPS 双端口并行，无强制跳转

| 属性 | 值 |
|------|-----|
| **文件** | `internal/routes/router.go:61-73` |
| **CWE** | CWE-319: 明文传输敏感信息 |

### 问题代码

```go
// internal/routes/router.go
if sslConfig.Enabled {
    // 启动 HTTPS（goroutine）
    go func() {
        r.RunTLS(":"+httpsPort, sslConfig.CertFile, sslConfig.KeyFile)
    }()
    // 启动 HTTP（主协程）—— HTTP 端口仍然开放！
    r.Run(":" + httpPort)
}
```

### 攻击路径

```
用户在浏览器输入 example.com 或点击 http:// 链接
    ↓
DNS 解析到服务器 HTTP 端口
    ↓
攻击者在同一网络中进行 ARP 欺骗或 DNS 劫持
    ↓
用户请求被拦截，响应被篡改（SSL Stripping）
    ↓
用户凭据通过明文 HTTP 传输
```

### 修复方案

```go
// 方案1: 直接关闭 HTTP 端口，仅保留 HTTPS
if sslConfig.Enabled {
    r.RunTLS(":"+httpsPort, sslConfig.CertFile, sslConfig.KeyFile)
    return
}

// 方案2: HTTP 做 301 重定向到 HTTPS
go func() {
    httpSrv := &http.Server{
        Addr: ":" + httpPort,
        Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            http.Redirect(w, r, "https://"+r.Host+":"+httpsPort+r.RequestURI, http.StatusMovedPermanently)
        }),
    }
    httpSrv.ListenAndServe()
}()
r.RunTLS(":"+httpsPort, sslConfig.CertFile, sslConfig.KeyFile)
```

---

## H4. 重置密码仅凭邮箱验证码，无二次校验

| 属性 | 值 |
|------|-----|
| **文件** | `internal/api/v1/auth.go:252-280`, `internal/service/user.go:220-236` |
| **CWE** | CWE-640: 密码找回机制薄弱 |

### 问题代码

```go
// internal/service/user.go:220-236
func ModifyPwd(ctx *gin.Context, modifyPwdReq dto.ModifyPwdReq) error {
    if cache.GetEmailCode(modifyPwdReq.Email) != modifyPwdReq.Code {
        return errors.New("邮箱验证错误")
    }
    // 直接改密码，不需要旧密码验证或登录态验证
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(modifyPwdReq.Password), bcrypt.DefaultCost)
    global.Mysql.Model(&model.User{}).Where("email = ?", modifyPwdReq.Email).
        Update("password", hashedPassword).Error
}
```

### 攻击路径

```
攻击者获取目标的邮箱验证码（暴力枚举/social engineering/短信拦截）
    OR
邮箱验证码长度仅为 6 位数字（1,000,000 种组合），无限速
    ↓
攻击者直接调用 /api/v1/auth/resetpwdCheck + /api/v1/auth/modifyPwd
    ↓
账号密码被重置，攻击者登录账号
```

### 修复方案

```go
// 方案1: 要求已登录用户提供旧密码
func ModifyPwd(ctx *gin.Context, modifyPwdReq dto.ModifyPwdReq) error {
    // 从 token 获取用户 ID
    userId := ctx.GetUint("userId")
    if userId == 0 {
        return errors.New("需要登录")
    }
    // 验证旧密码
    user, _ := FindUserById(userId)
    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(modifyPwdReq.OldPassword)) != nil {
        return errors.New("旧密码错误")
    }
    // ... 其余逻辑
}

// 方案2: 重置密码后立即失效所有已有 token
```

---

## H5. 前端代理禁用 TLS 证书验证

| 属性 | 值 |
|------|-----|
| **文件** | `nuxt.config.ts:17` |
| **CWE** | CWE-295: 证书验证不当 |

### 问题代码

```typescript
// nuxt.config.ts:17
const proxyAgent = new https.Agent({ rejectUnauthorized: false });
```

### 风险

开发环境中代理到后端时跳过了 TLS 证书验证，中间人攻击可劫持 API 流量。如果该配置被带入生产构建，风险更大。

### 修复方案

```typescript
// 仅在开发环境跳过验证
const isDev = process.env.NODE_ENV === 'development';
const proxyAgent = isDev && useHttps
  ? new https.Agent({ rejectUnauthorized: false })
  : undefined;
```

---

## H6. CSRF 防护仅依赖 `X-Requested-With` 头

| 属性 | 值 |
|------|-----|
| **文件** | `src/utils/request.ts:56` |
| **CWE** | CWE-352: 跨站请求伪造 |

### 问题代码

```typescript
// src/utils/request.ts:56
const service: AxiosInstance = axios.create({
    withCredentials: true,
    headers: {
        'X-Requested-With': 'XMLHttpRequest',  // 唯一的 CSRF 防护
    },
});
```

### 风险

`X-Requested-With` 是一个非标准头，不能作为可靠的 CSRF 防护：
- `fetch` API 可以自由设置该头（`mode: 'cors'`）
- 攻击者可在自己的域下构造带有 `X-Requested-With` 的请求

### 修复方案

```go
// 后端: 生成 CSRF token（可基于 Redis + 用户 session）
func generateCSRFToken(userId uint) string {
    // 生成随机 token 并存入 Redis
}

// 前端: 每次请求从 meta tag 或 cookie 读取 CSRF token 添加到 header
service.interceptors.request.use((config) => {
    config.headers['X-CSRF-Token'] = getCSRFToken();
    return config;
});
```

---

# 3. 中危漏洞

## M1. 无全局 API 速率限制

| 属性 | 值 |
|------|-----|
| **文件** | 全局缺失 |
| **CWE** | CWE-799: 未限制资源交互次数 |

### 风险

| 接口 | 目前防护 | 风险 |
|------|---------|------|
| 登录 | 3 次尝试后触发滑块验证（Redis） | 基本防护，但 IP 级别未限制 |
| 注册 | 无限制 | 可批量注册账号 |
| 发送邮箱验证码 | 无限制 | 可消耗邮件配额、骚扰用户 |
| 重置密码检查 | 无限制 | 可枚举邮箱是否注册 |
| API 总请求 | 无限制 | 可 DoS 攻击 |

### 修复方案

在 Gin 中间件层面添加 IP + UserID 双维度的令牌桶限流：

```go
// 参考实现: internal/middleware/ratelimit.go
func RateLimit() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        ip := ctx.ClientIP()
        key := "ratelimit:" + ip
        // 从 Redis 获取当前计数，超过阈值则返回 429
        // ...
    }
}
```

---

## M2. WebSocket 认证 token 通过 URL Query 传递

| 属性 | 值 |
|------|-----|
| **文件** | `internal/middleware/auth.go:104` |
| **CWE** | CWE-598: 通过 GET 请求查询字符串传递敏感信息 |

### 问题代码

```go
func WsAuth() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        tokenString := ctx.Query("token")  // 从 URL 参数获取 token
        _, claims, err := jwt_parse.ParseToken(tokenString)
        // ...
    }
}
```

### 风险

- Query 参数出现在服务器访问日志中
- Query 参数出现在浏览器历史记录中
- 通过 Referer 头泄露给第三方
- WebSocket 握手阶段 URL 可能被中间设备记录

### 修复方案

```go
// 方案1: 在 WebSocket upgrade 请求中使用 Authorization header
tokenString := ctx.GetHeader("Authorization")

// 方案2: 使用单独的 short-lived ticket 机制：
// 先通过 REST API 获取一次性 ticket（有效期 30s），再用 ticket 建立 WebSocket
ticket := ctx.Query("ticket")
```

---

## M3. 无 HTTP 服务器超时设置

| 属性 | 值 |
|------|-----|
| **文件** | `internal/routes/router.go` |
| **CWE** | CWE-400: 未限制资源消耗 |

### 问题代码

```go
// internal/routes/router.go:71
r.Run(":" + httpPort)  // 使用默认 server，没有任何超时配置
```

### 风险

Gin/Go 默认 `http.Server` 不设超时可能导致：
- **Slowloris 攻击**: 攻击者慢慢发送请求头，长期占用连接
- **连接耗尽**: 慢速读取响应可耗尽服务器连接池

### 修复方案

```go
// internal/routes/router.go
srv := &http.Server{
    Addr:              ":" + httpPort,
    Handler:           r,
    ReadHeaderTimeout: 10 * time.Second,
    ReadTimeout:       30 * time.Second,
    WriteTimeout:      30 * time.Second,
    IdleTimeout:       120 * time.Second,
}
srv.ListenAndServe()
```

---

## M4. 注销不立即失效 Access Token

| 属性 | 值 |
|------|-----|
| **文件** | `internal/service/user.go:154-163` |
| **CWE** | CWE-613: 会话过期不充分 |

### 问题代码

```go
func Logout(_ *gin.Context, tokenReq dto.TokenReq) {
    if tokenReq.RefreshToken == "" {
        return
    }
    _, claims, err := jwt.ParseToken(tokenReq.RefreshToken)
    if err != nil {
        return
    }
    cache.DelRefreshToken(claims.UserId, tokenReq.RefreshToken)
    // 注意: 没有将 accessToken 加入黑名单
}
```

### 风险

用户注销后，已发出的 accessToken 仍有效直到过期（60 分钟）。如果用户因设备丢失或公共电脑注销，攻击者可能在这 60 分钟内使用已有的 accessToken。

### 修复方案

```go
// 将 accessToken 加入 Redis 黑名单（使用 accessToken 过期时间作为 TTL）
if tokenReq.AccessToken != "" {
    cache.SetBlacklistedAccessToken(tokenReq.AccessToken, claims.UserId, accessTokenExpiry)
}
```

或者在 Auth 中间件中检查 Redis 中该 token 是否已被吊销：

```go
if cache.IsAccessTokenRevoked(tokenString) {
    resp.Result(ctx, 3000, nil, "TOKEN已失效")
    ctx.Abort()
    return
}
```

---

## M5. 用户管理列表后端未限制最大 pageSize

| 属性 | 值 |
|------|-----|
| **文件** | `internal/service/user.go:239-244` |
| **CWE** | CWE-770: 未限制资源分配 |

### 问题代码

```go
func GetUserListManage(userListReq dto.UserListReq) (total int64, users []vo.UserInfoManageResp) {
    global.Mysql.Model(&model.User{}).Count(&total)
    global.Mysql.Model(&model.User{}).Limit(userListReq.PageSize).
        Offset((userListReq.Page - 1) * userListReq.PageSize).Scan(&users)
    return
}
```

### 风险

虽然前端有分页限制，但后端没有对 `pageSize` 做上限。攻击者可直接调用 API 传入大 pageSize（如 100000）一次性获取大量用户数据（包括邮箱等敏感信息）。

### 修复方案

```go
func GetUserListManage(userListReq dto.UserListReq) (total int64, users []vo.UserInfoManageResp) {
    if userListReq.PageSize <= 0 || userListReq.PageSize > 50 {
        userListReq.PageSize = 10  // 或 50 硬上限
    }
    // ...
}
```

---

## M6. 邮箱验证码可被暴力枚举

| 属性 | 值 |
|------|-----|
| **文件** | `internal/api/v1/auth.go`（注册、邮箱登录、修改密码） |
| **CWE** | CWE-307: 关键操作的暴力枚举防护不足 |

### 问题

邮箱验证码为 6 位数字（1,000,000 种组合），以下接口直接验证：
- 注册 (`POST /api/v1/auth/register`)
- 邮箱登录 (`POST /api/v1/auth/login/email`)
- 修改密码 (`POST /api/v1/auth/modifyPwd`)

没有验证码尝试次数限制。

### 攻击路径

```
攻击者知道目标的邮箱
    ↓
调用 POST /api/v1/auth/resetpwdCheck （触发发送验证码）
    ↓
每秒尝试 100 个验证码，3 小时内可穷举全部 100 万种组合
    ↓
验证码命中 → 立即调用 modifyPwd 修改密码
```

### 修复方案

```go
// 为每个邮箱限制验证码尝试次数
cache.IncrEmailCodeTryCount(email)
if cache.GetEmailCodeTryCount(email) > 5 {
    cache.DelEmailCode(email)  // 使当前验证码失效
    return errors.New("验证码尝试次数过多，请重新获取")
}
```

---

# 4. 低危/信息

## I1. 依赖版本偏旧

| 依赖 | 当前版本 | 最新版本 | 已知漏洞 |
|------|---------|---------|---------|
| `gin` | v1.9.1 | v1.10.x | 部分中间件安全改进 |
| `golang-jwt/jwt` | v4.5.0 | v5.x | v4 不再活跃维护 |
| `gorm` | v1.25.5 | v1.26+ | 安全修复 |
| `gorilla/websocket` | v1.5.1 | v1.5.3 | 有安全修复 |
| `go-playground/validator` | v10.14.0 | v10.22+ | 多项校验改进 |
| `axios`（前端） | — | 最新 | 需确认版本 |
| `element-plus`（前端） | — | 最新 | 需确认版本 |

## I2. 配置文件未加入 .gitignore

`conf/application.prod.yaml` 和 `conf/application.dev.yaml` 未在 `.gitignore` 中排除。如果仓库是公开的，凭据立即泄露。即使私有，所有有仓库访问权限的人都能看到生产密码。

## I3. 密码存储使用 bcrypt ✅（好评项）

```go
// internal/service/user.go:35
hashedPassword, _ := bcrypt.GenerateFromPassword(
    []byte(registerReq.Password), bcrypt.DefaultCost,
)
```

密码正确使用了 bcrypt 哈希存储，这是行业标准做法。

## I4. 使用 Casbin RBAC 权限控制 ✅

使用 Casbin 做 RBAC 权限控制（`internal/middleware/auth.go:66`），实现了 route-level 的细粒度授权。需要验证权限配置完整性。

## I5. 安全响应头部分设置 ✅

| 响应头 | 状态 |
|--------|------|
| `X-Content-Type-Options: nosniff` | ✅ 已设置 |
| `Referrer-Policy: strict-origin-when-cross-origin` | ✅ 已设置 |
| `Permissions-Policy` | ✅ 已设置 |
| `X-Frame-Options: DENY`（非 embed） | ✅ 已设置 |
| `Strict-Transport-Security` | ✅ 有条件设置（仅 HTTPS）|
| `Content-Security-Policy` | ❌ 未设置（后端）/ 设置但过松（前端）|

---

# 5. 依赖版本风险

## golang.org/x/crypto v0.33.0

当前使用的 `bcrypt` 来自该包，版本较新，无已知漏洞。

## github.com/golang-jwt/jwt/v4 v4.5.0

该包是 `github.com/dgrijalva/jwt-go` 的官方继承者，但：
- v4 版本已进入低维护模式
- 官方推荐升级到 v5（有 breaking changes）
- 注意 `ParseUnverified` 的使用（`pkg/jwt/jwt.go:82`）

## gin v1.9.1

- 没有启用 `TrustedPlatform` 相关配置
- `SetTrustedProxies(nil)` 等于不校验任何代理头（显式关闭了 `X-Forwarded-For` 信任）

---

# 6. 综合修复优先级

| 优先级 | 编号 | 问题 | 预估工作量 | 影响面 |
|--------|------|------|-----------|--------|
| 🔴 P0 | C1 | 配置文件明文凭据 | 1h（轮换+gitignore+环境变量）| 全部基础设施 |
| 🔴 P0 | C2 | Access/Refresh 同密钥 | 0.5h | 认证系统 |
| 🔴 P0 | C4 | CORS 通配符 | 0.2h | 全部 API |
| 🔴 P0 | C3 | JWT 弱密钥 | 0.3h | 认证系统 |
| 🟠 P1 | C5 | Token 存 localStorage | 2-4h | 全部用户会话 |
| 🟠 P1 | H1 | CSP 策略过松 | 1-2h | XSS 防护 |
| 🟠 P1 | H4 | 重置密码无二次校验 | 1h | 账号安全 |
| 🟠 P1 | H6 | CSRF 防护不足 | 2-3h | 全部写操作 |
| 🟡 P2 | H2 | 日志记录敏感数据 | 0.5h | 隐私合规 |
| 🟡 P2 | H3 | HTTP/HTTPS 双开 | 0.5h | 传输安全 |
| 🟡 P2 | H5 | 代理跳过 TLS 验证 | 0.3h | 开发环境安全 |
| 🟡 P2 | M1 | 缺少速率限制 | 2-4h | 整体抗攻击 |
| 🟡 P2 | M3 | HTTP 无超时 | 0.3h | DoS 防护 |
| 🔵 P3 | 其余 | 各项中低危问题 | 视情况 | 纵深防御 |

---

*报告结束。需要针对任何一项展开详细修复方案或直接动手修改吗？*
