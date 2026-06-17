## 多音轨切换

### 概述

播放器支持多音轨视频（如日语原声 + 英语配音），通过 dash.js 原生 API 实现音轨切换。仅 DASH 模式支持无缝切换；HLS 模式回退默认音轨并日志告警。

### 触发条件

后端返回的 dash MPD `<AdaptationSet mimeType="audio">` 包含多个 `<Representation>`（对应不同语言）时，播放器自动检测并显示音轨选择按钮。

### 音轨选择按钮

- 仅当 `audioTrackCount > 1` 时显示
- 单音轨视频不显示 UI 变化
- 按钮默认展示当前音轨标题（如"日本語"、"English"）
- 点击弹出语言列表

### 实现

**文件**: `src/components/video-player/index.vue`

```typescript
// dash.js 音轨 API
const player = dashjs.MediaPlayer().create()
player.initialize(videoElement, manifestUrl, autoPlay)

// 获取可用音轨
const tracks = player.getTracksFor('audio')

// 切换音轨
player.setCurrentTrack(selectedTrack)
```

### 后端数据

MPD 中每个音轨以独立的 `<Representation>` 出现：

```xml
<AdaptationSet mimeType="audio" contentType="audio">
  <Role schemeIdUri="urn:mpeg:dash:role:2011" value="main" />
  <Representation id="audio-jpn" bandwidth="192000">
    <AudioChannelConfiguration schemeIdUri="urn:mpeg:dash:23003:3:audio_channel_configuration" value="2" />
    <BaseURL>audio_jpn.m4s</BaseURL>
    <SegmentBase indexRange="1234-5678">
      <Initialization range="0-1233" />
    </SegmentBase>
  </Representation>
  <Representation id="audio-eng" bandwidth="192000">
    ...
    <BaseURL>audio_eng.m4s</BaseURL>
  </Representation>
</AdaptationSet>
```

### 音轨表

DB `audio_track` 表记录：

| 字段 | 说明 |
|------|------|
| `resource_id` | 所属资源 |
| `dir_name` | OSS 目录 |
| `language` | ISO 639-2 代码（jpn/eng/chi） |
| `title` | 可读标题（日语/英语/中文） |
| `is_default` | 是否默认音轨 |
| `audio_file` | 文件名（audio_jpn.m4s） |
| `init_range` / `index_range` | m4s 索引范围 |
| `bandwidth` | 码率 |
