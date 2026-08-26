# PGC 演进：主任务提示词（供后续迭代 / AI 执行）

## 背景与既定约束

- **暂不接外部转码**：本地转码仍由 `internal/service/transcoding.go` 承担；后续外发转码需抽象 `Transcoder` 接口（`Submit/Status/Callback`），Episode 只依赖 `TranscodeJob` 与产物表。
- **需要分季（季度）**：当前工程已通过「同一 `media_id` 下多条 `pgc_content`（每条即一季）」实现播放面板中的 `seasons`；本阶段不强制改表，仅在文档与 API 上与该语义对齐说明。
- **流程偏好**：先建节目/季元数据，再逐集添加；剧集允许**先占位、后绑定视频**。

## 本阶段必达目标（已完成或待验收）

1. **剧集壳（无视频）**：`episodes[].vid` 允许为 `0`/缺省，创建占位剧集；不参与「可播过滤」SQL（`JOIN video` 已天然排除 `vid=0`）。
2. **绑定视频**：提供 `PUT /api/v1/pgc/:pgc_id/episodes/:id/bind`，将占位剧集绑定到已存在的 `vid`，并触发原有 `markVideosAsPGCAttached` 行为。
3. **`current_episodes` 语义**：统计「已绑定 `vid>0` 且状态为正常」的剧集数，与占位集区分。
4. **测试**：对纯函数与绑定语义补充单元测试；每次合并前执行 `go test ./...` 与 `go vet ./...`。

## 后续阶段（按优先级）

| 优先级 | 任务 | 说明 |
|--------|------|------|
| P1 | `Transcoder` 接口 + Job 表 | Episode 绑定 `asset_id`/`job_id`，本地实现与外发实现可切换 |
| P1 | 播放快照 / 发布单 | 后台编辑与 C 端播放解耦，避免「改到一半被读到」 |
| P2 | 分季管理 API | 显式 `POST /pgc/media/:media_id/season`（或文档化现有「新建 pgc_content 共 media_id」流程） |
| P2 | 审核状态机细化 | 与转码失败码、失败原因、重试策略对齐 |
| P3 | OpenAPI / 前端字段表 | 与 `PGC_API.md`、DTO 同步 |

## 执行自检清单（合并前）

- [ ] `go test ./...` 无失败
- [ ] `go vet ./...` 无报错
- [ ] 新接口已有最小文档片段（`PGC_API.md` 或 Quick Ref）
- [ ] 可播路径（推荐 / 播放面板）不展示未绑定 `vid` 的剧集

## 给 AI 的短提示词（复制使用）

> 在 `interastral-peace.com/alnitak` 的 PGC 模块中，保持「media 下多条 pgc_content = 多分季」现状。实现剧集占位：创建/添加剧集时 `vid` 可选；新增 `bind` 接口绑定 `vid`；`current_episodes` 只统计 `vid>0` 的正常集；播放面板列表不返回 `vid=0`。补充 `internal/service` 单元测试并确保 `go test ./...`、`go vet ./...` 通过。不改外部转码，Look for Transcoder 抽象仅留注释或后续 PR。
