# 数据库迁移 + 多 OSS 文件同步 方案设计

## 一、数据库迁移系统

### 现状问题

- 只用 GORM AutoMigrate（只能加字段，不能删/改）
- 无版本管理，无回滚
- 数据补齐（backfill）散落在各处

### 方案

```
server/cmd/migrate/main.go  →  统一入口，子命令模式
  go run cmd/migrate/main.go -env=dev up          # 执行所有待迁移
  go run cmd/migrate/main.go -env=dev down         # 回滚最后一批
  go run cmd/migrate/main.go -env=dev status       # 查看状态
  go run cmd/migrate/main.go -env=dev create xxx   # 生成迁移模板
```

### 迁移文件结构

```
server/internal/migrate/
├── migrator.go          # 核心引擎（执行/回滚/记录）
├── files/               # 版本迁移文件
│   ├── 20260514_001_backfill_pgc_media.go
│   ├── 20260514_002_add_video_duration_column.go
│   └── ...
├── migrate.go           # 注册所有迁移文件
└── schema_migration.go  # schema_migrations 表模型
```

### 每个迁移文件格式

```go
// 20260514_001_xxx.go
func init() {
    Register(&Migration{
        Version: "20260514_001",
        Name:    "backfill_pgc_media",
        Up: func(tx *gorm.DB) error {
            // 迁移逻辑
        },
        Down: func(tx *gorm.DB) error {
            // 回滚逻辑
        },
    })
}
```

### 与启动流程结合

- 启动时检查 `schema_migrations` 表，有未执行迁移则打印警告（或自动执行）
- 不再依赖 `InitTables()` 里的 backfill 代码，统一走迁移

---

## 二、多 OSS 文件同步

### 现状问题

- `config.Storage` 是单例，只能配一个 OSS
- 想要多 OSS（如 MinIO 做主存储 + 阿里云做灾备）需要改架构

### 配置变更

```yaml
# application.yaml
storage:
  primary:
    oss_type: minio
    bucket: alnitak
    endpoint: oss.ayypd.cn:9002
    key_id: xxx
    key_secret: xxx
    domain: oss.ayypd.cn:9002
    private: true
    use_ssl: true
    sync_db: true   # 本 OSS 的数据库表是否同步

  replicas:
    - name: aliyun-backup
      oss_type: aliyun
      bucket: alnitak-backup
      endpoint: oss-cn-hangzhou.aliyuncs.com
      key_id: xxx
      key_secret: xxx
      region: cn-hangzhou
      domain: alnitak-backup.oss-cn-hangzhou.aliyuncs.com
      private: false
      use_ssl: true
      sync_db: false    # 只同步文件，不同步数据库
```

### 全局变量变更

```go
// global.go
var (
    Storage          oss.Storage           // 主存储（读写）
    StorageReplicas  []oss.ReplicaStorage  // 副本存储（只写）
)
```

### 同步方式（三种策略）

| 策略 | 说明 | 实时性 | 复杂度 |
|------|------|--------|--------|
| **同步写入** | 上传时同时写入所有 replica | 实时 | 低（可加异步队列） |
| **异步队列** | 上传主存储后，发消息到 channel，worker 消费后同步到 replica | 近实时 | 中 |
| **定时同步** | cron 定期扫描主存储，对比 replica 差异并补齐 | 有延迟 | 高（需要文件清单） |

### 推荐：异步队列模式

```
上传请求 → 主 OSS（同步） → 写入 sync_queue 表 → worker goroutine 消费 → 同步到 replica OSS
```

### sync_queue 表结构

```sql
CREATE TABLE file_sync_tasks (
    id            BIGINT PRIMARY KEY,
    object_key    VARCHAR(512) NOT NULL,     -- 文件路径
    source_oss    VARCHAR(64) NOT NULL,      -- 源 OSS 名称
    target_oss    VARCHAR(64) NOT NULL,      -- 目标 OSS 名称
    action        VARCHAR(16) NOT NULL,      -- sync / delete
    status        TINYINT DEFAULT 0,         -- 0=待处理 1=处理中 2=成功 3=失败
    retry_count   INT DEFAULT 0,
    error_msg     TEXT,
    created_at    DATETIME,
    updated_at    DATETIME
);
```

### 后台管理界面接入点

- **迁移管理**：查看迁移状态、手动执行/回滚
- **同步管理**：查看同步队列状态、手动触发同步、查看同步日志

---

## 三、实施路线

| 阶段 | 内容 |
|------|------|
| **Phase 1** | 实现迁移引擎 + 将所有现存 backfill 改为迁移文件 |
| **Phase 2** | 多 OSS 配置 + `oss.ReplicaStorage` 接口 + 同步队列 |
| **Phase 3** | 后台管理 UI（迁移管理 + 同步监控页面） |
