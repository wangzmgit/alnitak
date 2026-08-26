# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is "Alnitak", a Go-based video sharing platform server (similar to Bilibili/YouTube). It's a comprehensive backend service with video uploading, transcoding, user management, comments, and real-time features.

## Commands

### Development
```bash
# Run in development mode
go run cmd/main.go -env=dev

# Run in production mode (default)
go run cmd/main.go -env=prod

# Build the application
go build -o cmd.exe cmd/main.go

# Install dependencies
go mod tidy
```

### Testing
The project doesn't have a standard test command configured. Check for test files in individual packages before adding tests.

## Architecture

### Directory Structure
- `cmd/` - Application entry point (main.go)
- `internal/` - Private application code
  - `api/v1/` - HTTP API handlers
  - `routes/` - Route definitions and middleware setup
  - `service/` - Business logic layer
  - `global/` - Global variables and state
  - `config/` - Configuration management
  - `middleware/` - HTTP middleware
  - `domain/` - Domain models and DTOs
  - `initialize/` - Application initialization
  - `cron/` - Scheduled tasks
- `pkg/` - Reusable packages
  - `mysql/`, `redis/` - Database clients
  - `oss/` - Object storage interface
  - `logger/` - Logging utilities
  - `casbin/` - Authorization
  - `jwt/` - JWT authentication
- `conf/` - Configuration files (dev/prod YAML)
- `static/`, `upload/` - Static assets and uploads

### Key Components
1. **Web Framework**: Gin (HTTP router and middleware)
2. **Database**: GORM with MySQL, Redis for caching
3. **Authentication**: JWT tokens with Casbin RBAC
4. **Storage**: Multi-provider object storage (Minio, Aliyun OSS, Tencent COS)
5. **Real-time**: WebSocket support in `pkg/ws/`
6. **Video Processing**: Transcoding capabilities with GPU support
7. **ID Generation**: Snowflake algorithm for unique IDs
8. **Dynamic Manifest**: DASH MPD / HLS m3u8 generated on-the-fly

### Video Streaming Architecture (B站风格)

The system uses Bilibili-style audio/video separation with SegmentBase mode:

```
Transcoding flow:
  FFmpeg → Separate video.m4s + audio.m4s → Extract SegmentBase info → Save to DB

Playback flow:
  Request → Query video_index_file → Generate DASH MPD (SegmentBase) → Player uses byte-range requests

File structure per video:
  ./upload/video/{dir_name}/
    ├── upload.mp4              # Original file
    ├── 1920x1080_3000k_30_video.m4s   # 1080p video stream
    ├── 1280x720_2000k_30_video.m4s    # 720p video stream
    ├── 854x480_900k_30_video.m4s      # 480p video stream
    ├── 640x360_500k_30_video.m4s      # 360p video stream
    └── audio.m4s               # Shared audio stream (320kbps AAC)
```

**video_index_file table fields (SegmentBase mode):**
- Video stream: `video_file`, `video_bandwidth`, `video_codec`, `width`, `height`, `frame_rate`, `video_init_range`, `video_index_range`
- Audio stream: `audio_file`, `audio_bandwidth`, `audio_codec`, `audio_sample_rate`, `audio_init_range`, `audio_index_range`
- Common: `total_duration`, `quality`, `dir_name`

**API Endpoints:**
- `GET /api/v1/video/getVideoFile?resourceId=X&quality=Y&format=mpd` - Get DASH MPD manifest
- `GET /api/v1/video/stream/:file?key=K` - Stream file with byte-range support (SegmentBase)
- `GET /api/v1/video/slice/:file?key=K` - Legacy segment file access (SegmentList)

### Initialization Flow
The application initializes in this order (see `cmd/main.go:19-54`):
1. Config loading (dev/prod environment)
2. Logger setup
3. Captcha system
4. Object storage
5. Snowflake ID generator
6. MySQL database and tables
7. Video partition mapping
8. Redis cache
9. Casbin authorization
10. Background cron jobs
11. HTTP router (port 9000)

### Configuration
- Development: `conf/application.dev.yaml`
- Production: `conf/application.prod.yaml`
- Environment selected via `-env` flag (defaults to prod)

### API Structure
All endpoints are under `/api/v1/` with comprehensive route organization:
- Authentication, users, roles, permissions
- Video management, transcoding, partitions
- Comments, archives, collections
- Real-time messaging and online status
- File uploads and static serving

The server runs on port 9000 and serves static images at `/api/image/:file`.

## Development Rules (强制遵守)

### 1. 改动前必须检查原有逻辑
- 修改任何文件前，先读完相关文件的原有逻辑，确认新改动不会破坏已有功能
- 修改 API 前检查 git 历史确认原有函数保留
- 改路由前对比旧代码确认原有路由存在
- 改名/重命名前 grep 检查所有调用方

### 2. 改动后必须编译验证
- 后端：`go build ./...` + 打包 exe（`alnitak-server.exe` / `alnitak-worker.exe`）
- Flutter：`flutter analyze` + `flutter build apk --release --target-platform android-arm64 --split-per-abi`
- **Web 不编译**（`web-client`、`admin-client` 不运行 `yarn build`）

### 3. 前端三端同步（强制）
涉及后端接口变更（新增/修改/删除），**必须同步**以下所有前端：
- `E:\web\alnitak\web\web-client`（用户端 Web）
- `E:\web\alnitak\web\admin-client`（管理端 Web）
- `E:\alnitak_flutter`（Flutter App）
- 三端接口调用参数、响应解析、错误处理必须一致
