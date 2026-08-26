# 转码系统

## 架构概览

```
上传入口                       转码核心                          存储出口
─────────                   ──────────                        ─────────
upload.go                   ProcessVideo()                    uploadToOSS()
resource.go                   ├── processSingleQuality() 1..N    ├── 主 OSS (aliyun/minio/tencent/cloudflare)
  └── Enqueue()               │     ├── encodeVideoWithFallback()  └── 备用 OSS (backup, 异步容灾)
video.go                      │     │     ├── GPU: nvenc (AV1/H.265/H.264)
  └── VideoTransCoding()      │     │     └── CPU: libsvtav1/libx265/libx264
                              │     ├── encodeAudioTrack() / encodeAudioOnly()
                              │     └── saveIndexRecord()
                              ├── saveAudioTrackRecords()      DB 收尾
                              └── completeTransaction()        └── resource.status / video.status
```

### 组件角色

| 组件 | 职责 |
|------|------|
| `upload.go` | 视频上传完成后调用 `Transcoder.Enqueue()` 提交转码 |
| `resource.go` | 资源替换（换源）后触发转码 |
| `video.go` | 重新转码（重试/修复）时串行调用 `VideoTransCoding()` |
| `TranscodeService` | 本地转码核心：并发 quality + ffmpeg 进程管理 |
| `LocalTranscoder` | `Transcoder` 接口的本地实现，包装 `TranscodeService` |
| `RemoteTranscoder` | `Transcoder` 接口的远程实现（Phase 4，当前为桩） |

---

## Transcoder 接口

所有转码调用通过 `Transcoder` 接口解耦，支持 `local` / `remote` 两种后端。

```go
// internal/service/transcoder.go
type Transcoder interface {
    Enqueue(ctx context.Context, info *dto.TranscodingInfo) error
    GetProgress(ctx context.Context, resourceID uint) (*TranscodingProgress, error)
    Cancel(ctx context.Context, videoID uint) error
}
```

### 接口实现

#### LocalTranscoder（mode=local）

当前默认实现，行为与重构前完全一致：

| 方法 | 行为 |
|------|------|
| `Enqueue` | 启动 goroutine 执行 `ProcessVideo()`，立即返回 |
| `GetProgress` | 从 `TranscodeService` 内存进度表读取 |
| `Cancel` | 调用 `StopTranscodingAndCleanup()` 终止 ffmpeg 进程 + 清理产物 |

#### RemoteTranscoder（mode=remote）

Phase 4 实现，当前为安全桩（`ErrRemoteNotImplemented`）：

```go
// 计划流程：
// Enqueue  → XADD transcoding:queue (Redis Streams)
// GetProgress → HGETALL transcoding:status:{resourceId}
// Cancel  → PUBLISH transcoding:cancel {videoID}
```

### 调用链重构

| 调用方 | 原写法 | 当前写法 |
|--------|--------|---------|
| `upload.go` | `go VideoTransCoding(info)` | `GetCurrentTranscoder().Enqueue(ctx, info)` |
| `resource.go` | `go VideoTransCoding(info)` | `GetCurrentTranscoder().Enqueue(ctx, info)` |
| `video.go` | `VideoTransCoding(info)` | 保留（串行重转码需要阻塞语义） |

`VideoTransCoding()` 保留为阻塞函数——内部直接调用 `ProcessVideo()`，供 video.go 的串行 re-transcode 使用。

---

## 模式切换

```yaml
# config.yaml
transcoding:
  mode: local          # local（默认） | remote（需部署 transcoder-worker）

  # local 模式配置（现有参数不变）
  use_gpu: false
  use_h265: false
  use_av1: false
  generate_1080p60: false
  max_cpu_concurrency: 2
  max_gpu_concurrency: 2
```

`mode=local` 零行为变更。`mode=remote` 需先部署 Worker（Phase 4），否则所有 `Enqueue` 返回 `ErrRemoteNotImplemented`。

---

## 视频转码流程

### 入参（dto.TranscodingInfo）

```go
type TranscodingInfo struct {
    VideoID      uint    // 视频稿件 ID
    ResourceID   uint    // 分 P 资源 ID
    InputFile    string  // 源文件路径
    OutputDir    string  // 输出目录
    DirName      string  // OSS 目录名
    Duration     float64 // 视频时长（秒）
    Width, Height int    // 原始分辨率
    CodecName    string  // 源编码
    FPS, FPS30, FPS60 string // 帧率
    VideoBitRate int     // 视频码率
    AudioBitRate int     // 音频码率
    AudioSampleRate int  // 音频采样率
    AudioChannels   int  // 音频声道数
    AudioStreams  []AudioStreamProbe // 全部音轨探测结果（多音轨支持）
    OriginalVideoStatus int  // 重新转码时记录原始状态
}
```

### 核心流水线 ProcessVideo

```
ProcessVideo(ctx, info)
│
├── getTranscodingTarget(info) → 确定目标画质列表
│   └── 按原始分辨率 + 码率确定最高档位，向下生成 360p/480p/720p/1080p/4K
│
├── 并发转码（每个 quality 一个 goroutine）
│   ├── processSingleQuality(ctx, info, target)
│   │   ├── encodeVideoWithFallback()  ← GPU/CPU + AV1/H.265/H.264 逐级降级
│   │   │   └── runVideoEncodeTask()   ← 拼接 ffmpeg 参数
│   │   ├── encodeAudioTrack() / encodeAudioOnly()
│   │   │   └── SharedTask 机制：跨 quality 共享音频编码结果（只编一次）
│   │   └── saveIndexRecord()          ← 解析 m4s init/index range
│   └── wg.Wait() 等待全部 quality 完成
│
├── saveAudioTrackRecords()  ← 多音轨记录落库
│
├── uploadToOSS(ctx, info)   ← 并发上传产物到 OSS（10 worker, 3 次重试）
│   └── PutObjectFromFile(objectKey, filePath)
│
└── completeTransaction(ctx, info, status) ← DB 状态流转
    ├── resource.status = WAITING_REVIEW | PROCESSING_FAIL
    ├── video.status   = WAITING_REVIEW | AUDIT_APPROVED
    └── cache.DelVideoInfo(videoID)
```

---

## 多音轨支持

### ffprobe 探测

`ProcessVideoInfo()` 自动探测源视频的全部音频流：

```go
type AudioStreamProbe struct {
    StreamIndex int    // ffprobe stream index
    Language    string // 语言代码 (ISO 639-2, 如 jpn/eng)
    SampleRate  int    // 采样率
    Channels    int    // 声道数
    BitRate     int    // 码率 (bps)
}
```

### 编码

- 每个音轨独立编码，通过 `SharedTask` 保证所有 quality goroutine 共享结果（只编一次）
- 默认音轨（第一条）编码失败 → 对应 quality 标记失败
- 附加音轨编码失败 → 仅记录日志，不影响 quality 状态

### 文件名约定

| 音轨 | 文件名 |
|------|--------|
| 默认 / 未识别语言 | `audio.m4s` |
| 日语 | `audio_jpn.m4s` |
| 英语 | `audio_eng.m4s` |
| 中文 | `audio_chi.m4s` |

### 音轨表

`audio_track` 表记录每个音轨的 Language、Codec、Init/Index Range、Bandwidth，播放器据此生成音轨选择 UI。

---

## GPU 降级策略

`encodeVideoWithFallback` 实现 4 级降级：

```
GPU AV1  → GPU H.265 → GPU H.264 → CPU H.264
(rtx40+)    (hevc)     (h264)      (兜底)
```

- 连续失败 3 次 → 禁用 GPU，后续 quality 及视频全部走 CPU
- GPU 成功一次 → 失败计数器归零
- 可通过 `max_gpu_fail_threshold` 常量调整
- `ResetGPUState()` 手动恢复（管理后台）

---

## HDR / 10-bit 检测

`ProcessVideoInfo` 中检测源视频位深：

- 10-bit 源 + 8-bit 目标（H.264）→ 日志告警色彩信息丢失
- 10-bit 源 + 10-bit 目标（H.265/AV1）→ 日志确认保留 HDR

---

## 编解码参数

| 编码器 | GPU/CPU | profile | 适用场景 |
|--------|---------|---------|----------|
| `av1_nvenc` | GPU (RTX 40+) | main | 最高压缩率，硬件要求高 |
| `hevc_nvenc` | GPU | main10 | H.265 10-bit，NVENC |
| `h264_nvenc` | GPU | high | H.264 硬件编码，兼容最好 |
| `libsvtav1` | CPU | - | AV1 软件编码，preset 6 |
| `libx265` | CPU | main10 | H.265 10-bit，slow preset |
| `libx264` | CPU | high | H.264 软件编码，medium preset |

## 并发控制

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `max_cpu_concurrency` | `maxCPUConcurrentTranscoding = 2` | CPU 模式最大并行转码视频数 |
| `max_gpu_concurrency` | `maxGPUConcurrentTranscoding = 2` | GPU 模式最大并行数 |
| `ossUploadMaxConcurrency` | 10 | OSS 上传并发数 |
| `ossUploadMaxRetries` | 3 | OSS 上传重试次数 |
| `ossUploadBackoff` | 1s, 2s, 4s | 上传重试指数退避 |

线程分配：GPU 编码固定 4 线程；CPU 编码按 `(runtime.NumCPU() - 2) / maxConcur` 动态分配，上限 8。

---

## 备用 OSS（多源容灾）

主 OSS 上传完成后，异步上传到备用 OSS：

```yaml
storage:
  backup:
    oss_type: "cloudflare"
    endpoint: "https://xxx.r2.cloudflarestorage.com"
    bucket: "backup-bucket"
    key_id: "..."
    key_secret: "..."
```

- 失败记录到 `backup_failure` 表，后续可重试
- 不影响主流程（goroutine 异步、错误仅日志）

---

## 转码相关 API

| 方法 | 路径 | 描述 | 鉴权 |
|------|------|------|------|
| POST | `/api/v1/video/reTranscodeVideo?vid={id}` | 重新转码整个视频（所有分P入队） | 后台管理 |
| POST | `/api/v1/video/reTranscodeResource?resourceID={id}` | 重新转码单个分P（仅失败/卡住的分P） | 后台管理 |

### reTranscodeResource 说明

用于在源文件问题修复后单独重试失败分P，不影响其他分P的状态。

- **请求参数**: `resourceID` — Resource 表主键
- **允许状态**: `PROCESSING_FAIL` 或 `VIDEO_PROCESSING`
- **响应**:
```json
{
  "code": 200,
  "data": { "resourceID": 123, "videoID": 456, "status": 1 },
  "msg": "已触发分P重新转码"
}
```

## 近期修复

| 问题 | 原因 | 修复 |
|------|------|------|
| 构建元数据 500 | 旧版 `.output/server/index.mjs` 进程残留，build ID 不匹配 | 杀掉 3000 端口旧进程 |
| Nuxt 水合不匹配 | 构建 500 的附带现象，非真实 SSR/CSR 差异 | 已验证域名访问无错误 |
| watch.vue PGC 水合 | 顶级 `await` PGC 绑定导致 SSR 不匹配 | 移到 `watch(videoInfo) + process.client` |
| watch.vue API 空指针 | 无 `videoApiData.value` 守卫 | 加 `if (videoApiData.value)` 判断 |
| 404 导航 SSR 异常 | `navigateTo('/404')` 在服务端执行 | 加 `process.client` 守卫 |
| 嵌套 `useAsyncData` | `asyncGetVideoInfoAPI` 已含 `useAsyncData`，外部再包一层 | 恢复原始调用模式 |

---

## 配置参考

详见 [`CONFIG_CHANGE.md`](../CONFIG_CHANGE.md) 的 storage 配置说明。

```yaml
transcoding:
  mode: local              # local | remote
  use_gpu: false
  use_h265: false
  use_av1: false
  generate_1080p60: false
  max_cpu_concurrency: 2
  max_gpu_concurrency: 2
```
