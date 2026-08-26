# 开发规范与注意事项

## 代码修改规则

### 1. 修改已有功能前必须先对比旧代码
- 修改 API 文件前，先查看 git 历史或备份，确认原有函数保留
- 修改路由文件前，对比旧代码确认原有路由函数存在
- 禁止直接覆盖/删除已有函数，如需重构先备份

### 2. 重命名需谨慎
- 改函数名/变量名前检查是否有其他地方调用
- 改路由函数名需同步更新 router.go 中的调用
- 改 API 名需同步更新 initialize/data.go 中的权限配置

### 3. 新增功能检查清单
- [ ] Model 层：新增表结构
- [ ] DTO/VO 层：请求/响应结构
- [ ] Service 层：业务逻辑
- [ ] API 层：处理函数
- [ ] Routes 层：路由注册
- [ ] initialize/data.go：API 权限配置
- [ ] initialize/tables.go：AutoMigrate（如需）
- [ ] 文档：更新 API 文档

### 4. 改动前检查 + 改动后编译（强制）
- **改动前**：必须先读完相关文件的原有逻辑，确认新改动不会破坏已有功能
- **改动后**：必须编译验证：
  - 后端：`go build ./...` + 打包 exe
  - Flutter：`flutter analyze` + `flutter build apk --release --target-platform android-arm64 --split-per-abi`
  - Web：**不编译**

### 5. 前端同步规则（强制）
- 涉及后端接口变更（新增/修改/删除），**必须同步**以下所有前端项目：
  - `E:\web\alnitak\web\web-client`（用户端 Web）
  - `E:\web\alnitak\web\admin-client`（管理端 Web）
  - `E:\alnitak_flutter`（Flutter App）
- 三端接口调用参数、响应解析、错误处理必须一致
- 修改后逐个检查三端是否编译通过：
  - Web：`yarn build`（如未禁止）
  - Flutter：`flutter analyze` + `flutter build apk --release --target-platform android-arm64 --split-per-abi`

### 5. 测试验证
- 修改后用 curl 测试关键 API 是否正常
- 启动服务检查路由注册日志
- 确认原有的公开/管理接口未受影响