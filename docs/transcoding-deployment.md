# 转码分离服务部署文档

## 概述

支持两种执行模式：

| 模式 | 说明 | 适用场景 |
|------|------|----------|
| `local`（默认） | 主进程内调用 ffmpeg 转码 | 单机部署、开发环境 |
| `remote` | 主进程推任务到 Redis Stream，`transcoder-worker` 消费执行 | GPU 分离、扩缩容 |

切换 mode 后对调用方透明——`VideoTransCoding()`、`GetVideoTranscodingProgress()`、`StopTranscodingAndCleanup()` 等接口行为一致。

## 架构

### Remote 模式

```
请求 → API → RemoteTranscoder → XADD Redis Stream
                                    │
                           transcoder-worker (N 个副本)
                                    │
                               ┌────┴────┐
                               │ encoding │ ← encoding_concurrency (cap=3)
                               │  signal  │     每视频同时跑几个 ffmpeg
                               └────┬────┘
                                    │
                          doEncodeWithFallback
                          ┌────┐  ┌──────┐
                          │GPU │  │ CPU  │
                          └─┬──┘  └──┬───┘
                       nvencSem  cpuSem
                        (cap=8)  (cap=2)
                      驱动上限   防线程打满
```

### 三层信号量

| 信号量 | 容量 | 作用 |
|--------|------|------|
| `sem` (worker_concurrency) | 可配置 | 整机同时处理几个视频 |
| `encodingSem` (encoding_concurrency) | 可配置 | 单个视频内同时跑几个 ffmpeg |
| `nvencSem` | 固定 8 | 全局 NVENC 硬件上限（GeForce 驱动限制），超了排队不降级 |
| `cpuSem` | 固定 2 | 全局 CPU 降级保护，最多 2 路 |

示例（`worker_concurrency=1, encoding_concurrency=3`，5 画质 + 1 音轨）：

- 3 个画质先拿到 `encodingSem` 跑 ffmpeg，2 个排队
- 3 个同时抢 `nvencSem`（有空位），全部 NVENC 硬编
- 完成一个释放 → 排队的拿到 `encodingSem` 接上
- GPU 炸了降级 CPU → 走 `cpuSem`，最多 2 路 CPU 同时跑
- 音频不受限，始终并行

## 前置条件

### Remote 模式
- Go 1.22+
- ffmpeg + ffprobe 在 PATH 中
- Redis 6.2+（Stream + Pub/Sub）
- Worker 节点可访问 OSS 和 Redis

## 主服务配置

```yaml
# server/conf/application.prod.yaml
transcoding:
  mode: remote
  generate_1080p60: true
  worker_concurrency: 3        # 主服务用不到，远程 Worker 读自己配置
  max_queue_depth: 10          # Redis Stream 积压上限

redis:
  host: 10.0.0.1               # 指向共享 Redis

storage:                       # 完整 OSS 配置
  oss_type: minio
  endpoint: oss.example.com:9002
  bucket: alnitak
  key_id: ""
  key_secret: ""
```

`worker_concurrency` 由 Worker 自己读自己的配置，主服务无需关心。

## Worker 部署

Worker 使用独立配置文件（从服务端复制后精简），只保留 `log`、`redis`、`storage`、`transcoding` 四段。

### 配置文件

```yaml
# conf/application.prod.yaml
log:
    level: info
    mode: prod
redis:
    host: 10.0.0.1
    port: 6379
storage:                       # 与主服务保持一致
    oss_type: minio
    endpoint: oss.example.com:9002
    bucket: alnitak
    key_id: ""
    key_secret: ""
    backup:                    # 可选备用 OSS
        oss_type: cloudflare
        endpoint: https://...
        key_id: ""
        key_secret: ""
transcoding:
    mode: remote
    generate_1080p60: true
    worker_concurrency: 1      # 同时处理几个视频
    encoding_concurrency: 3    # 每视频同时跑几个画质 ffmpeg
    work_dir: workdir          # 临时文件目录
    use_gpu: true
    use_h265: true
```

### 编译

```powershell
cd server
go build -o main.exe .
go build -o transcoder-worker.exe .\cmd\transcoder-worker
```

### Windows 服务注册（nssm）

```powershell
nssm install TranscoderWorker "E:\alnitak\worker\bin\transcoder-worker.exe"
nssm set TranscoderWorker AppDirectory "E:\alnitak\worker"
nssm set TranscoderWorker AppParameters "--env prod --health-port 9100"
nssm set TranscoderWorker Start SERVICE_AUTO_START
nssm start TranscoderWorker
```

参数说明：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--env` | prod | 加载 `conf/application.{env}.yaml` |
| `--health-port` | 9100 | 健康检查 HTTP 端口 |
| `--id` | hostname | Worker 标识 |

`--concurrency` 已废弃——并发数从配置文件的 `worker_concurrency` 读取。

### 启动验证

```powershell
curl http://localhost:9100/ready
curl http://localhost:9100/stats
```

```json
{
  "healthy": true,
  "concurrency": 1,
  "jobsActive": 0,
  "jobsTotal": 0,
  "jobsFailed": 0,
  "groupID": "transcoder-worker-PC-01"
}
```

## 健康检查

| 端点 | 说明 |
|------|------|
| `GET /health` | 存活检查，返回 `{"status":"ok"}` |
| `GET /ready` | 就绪检查，Redis 连通则 200 否则 503 |
| `GET /stats` | 运行时状态快照 |

## Redis Stream 监控

```bash
# 队列积压
redis-cli XLEN transcoding:queue

# 消费者组状态
redis-cli XINFO GROUP transcoding:queue transcoder

# 死信队列
redis-cli XLEN transcoding:deadletter

# 资源转码状态
redis-cli HGETALL transcoding:status:{resourceID}
```

## 多副本部署

```powershell
# 节点 A
transcoder-worker.exe --env prod --id "worker-a"
# 节点 B
transcoder-worker.exe --env prod --id "worker-b"
```

通过 Redis consumer group 自动分发任务。

## 故障恢复

- 消费使用 Redis Stream consumer group，`XAck` 确认
- Worker 崩溃后，未确认消息 5 分钟后被其他 Worker 认领
- 超过 3 次重试移入 `transcoding:deadletter`
- 启动时自动恢复 crash 遗留的 pending 消息

## 编码降级策略

```
GPU AV1 → GPU H.265 → GPU H.264 → CPU H.264
```

每步降级受信号量保护：
- GPU 阶段抢 `nvencSem`（8 路上限）
- CPU 阶段抢 `cpuSem`（2 路上限）

## 版本记录

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.0 | 2026-06-17 | 初始版，支持 local/remote 模式 |
| 1.1 | 2026-06-17 | 三层信号量架构，独立 Worker 配置 |
