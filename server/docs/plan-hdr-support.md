# HDR 视频转码支持 — 实施计划

## TL;DR

当前后端已支持 10-bit 编码（H.265 GPU/CPU），但**完全没有 HDR 处理逻辑**。本计划分 4 个阶段实现完整 HDR 支持，核心改动约 3 个文件。

---

## 阶段 0：前置知识 — HDR 视频的本质

| 概念 | 说明 |
|------|------|
| **色彩原色** | HDR 用 BT.2020（宽色域），SDR 用 BT.709 |
| **电光转移函数（EOTF）** | HDR 用 PQ（ST.2084）或 HLG，SDR 用 BT.709 |
| **色彩矩阵** | HDR 用 BT.2020nc（非恒定亮度），SDR 用 BT.709 |
| **位深** | HDR 最低 10-bit |
| **HDR10 静态元数据** | `MaxFALL`、`MaxCLL`、`master_display`（色度坐标 + 白点 + luminance min/max）|

ffprobe 输出示例（HDR 源）：
```json
{
  "streams": [{
    "codec_type": "video",
    "pix_fmt": "yuv420p10le",
    "color_primaries": "bt2020",
    "color_transfer": "smpte2084",
    "color_space": "bt2020nc",
    "color_range": "tv",
    "side_data_list": [{
      "side_data_type": "Mastering display metadata",
      "red_x": "0.708", "red_y": "0.292",
      "green_x": "0.170", "green_y": "0.797",
      "blue_x": "0.131", "blue_y": "0.046",
      "white_point_x": "0.3127", "white_point_y": "0.3290",
      "luminance_min": "0.0050", "luminance_max": "1000.0000"
    }, {
      "side_data_type": "Content light level metadata",
      "max_content": "848", "max_average": "236"
    }]
  }]
}
```

---

## 阶段 1：ffprobe 元数据采集扩展

### 改 `internal/global/transcoding.go`

**文件定位**：`internal/global/transcoding.go` — `VideoInfo.Streams` 结构体

**改动**：添加 HDR 相关字段

```go
type Streams struct {
    // ... 现有字段 ...
    PixFmt           string `json:"pix_fmt,omitempty"`
    // 新增 ↓
    ColorPrimaries   string        `json:"color_primaries,omitempty"`
    ColorTransfer    string        `json:"color_transfer,omitempty"`
    ColorSpace       string        `json:"color_space,omitempty"`
    ColorRange       string        `json:"color_range,omitempty"`
    BitsPerRawSample string        `json:"bits_per_raw_sample,omitempty"`
    SideDataList     []SideData    `json:"side_data_list,omitempty"`
    // ... 现有字段 ...
}

// 新增
type SideData struct {
    SideDataType string `json:"side_data_type,omitempty"`
    RedX         string `json:"red_x,omitempty"`
    RedY         string `json:"red_y,omitempty"`
    GreenX       string `json:"green_x,omitempty"`
    GreenY       string `json:"green_y,omitempty"`
    BlueX        string `json:"blue_x,omitempty"`
    BlueY        string `json:"blue_y,omitempty"`
    WhitePointX  string `json:"white_point_x,omitempty"`
    WhitePointY  string `json:"white_point_y,omitempty"`
    LuminanceMin string `json:"luminance_min,omitempty"`
    LuminanceMax string `json:"luminance_max,omitempty"`
    MaxContent   string `json:"max_content,omitempty"`
    MaxAverage   string `json:"max_average,omitempty"`
}
```

**注意**：ffprobe 的 side_data_list 字段名带空格（"Mastering display metadata"），用 JSON 的 `omitempty` 即可。现有 ffprobe 命令 `-show_streams` 已包含这些字段，不需要改 ffprobe 参数——Go 反序列化时忽略未知字段，新字段会被自动填充。

### 加判 HDR 工具函数（建议文件：`internal/service/hdr.go` 或加到 `transcoding.go`）

```go
// IsHDRSource 判断视频流是否为 HDR 内容
func IsHDRSource(stream *global.Streams) bool {
    if stream == nil {
        return false
    }
    // BT.2020 宽色域 + PQ/HLG 转移 = HDR
    return stream.ColorPrimaries == "bt2020" &&
        (stream.ColorTransfer == "smpte2084" || stream.ColorTransfer == "arib-std-b67" || stream.ColorTransfer == "smpte428")
}

// HasHDROMetadata 是否包含 HDR10 静态元数据
func HasHDROMetadata(stream *global.Streams) bool {
    for _, sd := range stream.SideDataList {
        if strings.Contains(sd.SideDataType, "Mastering display") {
            return true
        }
    }
    return false
}

// BuildColorParams 构造 FFmpeg 颜色参数（如果源是 HDR）
func BuildColorParams(stream *global.Streams) []string {
    if !IsHDRSource(stream) {
        return nil
    }
    // 从源继承色彩信息（如果源有明确值）
    prim := stream.ColorPrimaries   // "bt2020"
    trc  := stream.ColorTransfer    // "smpte2084" 或 "arib-std-b67"
    spc  := stream.ColorSpace       // "bt2020nc"
    if prim == "" { prim = "bt2020" }
    if trc  == "" { trc  = "smpte2084" }
    if spc  == "" { spc  = "bt2020nc" }
    return []string{
        "-color_primaries", prim,
        "-color_trc", trc,
        "-colorspace", spc,
    }
}

// BuildHDR10SideDataParams 构造 HDR10 静态元数据参数（仅 H.265/AV1）
func BuildHDR10SideDataParams(stream *global.Streams) []string {
    if !HasHDROMetadata(stream) {
        return nil
    }
    // 解析 mastering display 元数据 → x265-params 格式
    // 格式: "G(13250,34500)B(7500,3000)R(34000,16000)WP(15635,16450)L(10000000,1)"
    // 从 side_data 中提取、缩放到 50000 分度
    // ... 实现略，参考 x265 HDR10 文档
    // 暂时可以只传递 -x265-params "hdr10-opt=1:hdr10=1"
    // 完整 mastering display 需从 SideData 构造
    return nil // TODO
}
```

---

## 阶段 2：编码器链路 — HDR 参数注入

### 核心改动文件：`internal/service/transcoding.go`

#### 2a. `runVideoEncodeTask` 函数（~line 727）

**改动**：在函数开头获取 HDR 信息并传递给参数生成逻辑

```go
func (s *TranscodingService) runVideoEncodeTask(
    ctx context.Context, videoID, resourceID uint, inputFile, outputFile, resolution, rate, fps, qualityName string,
    totalDuration float64, useGpu bool, useAv1 bool, useHevc bool, cancelFunc context.CancelFunc,
) error {
    // 新增 ↓
    videoInfo, _ := getVideoInfo(inputFile)
    var videoStream *global.Streams
    for i := range videoInfo.Stream {
        if videoInfo.Stream[i].CodecType == "video" {
            videoStream = &videoInfo.Stream[i]
            break
        }
    }
    isHDR := IsHDRSource(videoStream)
    // ↑ 新增
```

> **注意**：`runVideoEncodeTask` 当前不走 `ProcessVideoInfo`，需要加一次轻量 ffprobe 获取颜色元数据。如果担心性能，可以在调用方提前探测并传入 isHDR 布尔值。

#### 2b. 编码器 switch 分支（~line 757-799）

每个分支根据 `isHDR` 追加颜色参数和像素格式：

```
switch {
case useGpu && useAv1:
    args = append(args,
        "-c:v", "av1_nvenc", ...
        // 如果 isHDR: "-pix_fmt", "yuv420p10le"  // 目前是 yuv420p 8-bit
        // 否则: "-pix_fmt", "yuv420p"
    )
    if isHDR {
        args = append(args, BuildColorParams(videoStream)...)
    }

case useGpu && useHevc:
    // 已是 yuv420p10le，无需改 pix_fmt
    if isHDR {
        args = append(args, BuildColorParams(videoStream)...)
        // 追加 HDR10 SEI："-sei", "+hdr10_plus+mastering_display"
        // 追加 x265-params 扩展
    }

case useGpu && !useHevc:
    // H.264 — 通常不用于 HDR，但可以保持颜色元数据
    // 如果 isHDR: 可以考虑 tone map 到 SDR 或保留元数据
    // 现阶段只追加颜色参数

case !useGpu && useAv1:
    // AV1 CPU — 目前 yuv420p 8-bit
    // 如果 isHDR: "-pix_fmt", "yuv420p10le"
    // 追加颜色参数

case !useGpu && useHevc:
    // 已是 yuv420p10le
    if isHDR {
        // 在 -x265-params 追加 "hdr10-opt=1:..."
        // 追加颜色参数
    }

default:  // H.264 CPU
    // 如果 isHDR: 追加颜色参数（保留元数据）
}
```

#### 2c. 各编码器分支具体改动表

| 分支 | 当前 pix_fmt | isHDR 时 | 颜色参数 | 额外 |
|------|-------------|----------|----------|------|
| **AV1 GPU** | yuv420p | yuv420p10le | 追加 | 加 `-svtav1-params` 传 color desc |
| **H.265 GPU** | yuv420p10le | 不变 | 追加 | 加 `-sei +hdr10_plus+mastering_display` |
| **H.264 GPU** | yuv420p | 不变 | 追加 | 不改 pix_fmt |
| **AV1 CPU** | yuv420p | yuv420p10le | 追加 | 加 `-svtav1-params` color desc |
| **H.265 CPU** | yuv420p10le | 不变 | 追加 | `-x265-params` 追加 `hdr10-opt=1` |
| **H.264 CPU** | yuv420p | 不变 | 追加 | 不改 pix_fmt |

#### 2d. 构建颜色参数的时机

颜色参数（`-color_primaries`、`-color_trc`、`-colorspace`、`-color_range`）**必须在编码器之后、输出文件之前**追加，顺序如下：

```
ffmpeg -i input
  -filter_complex ...scale...
  -c:v libx265 -pix_fmt yuv420p10le ...
  -color_primaries bt2020 -color_trc smpte2084 -colorspace bt2020nc
  -x265-params "hdr10-opt=1:..."
  output.mp4
```

---

## 阶段 3：HDR→SDR Tone Mapping（可选）

### 场景
- 源视频是 HDR（BT.2020/PQ），但用户选的是 H.264 或没开 HDR 模式
- 需要将 HDR 映射到 SDR，否则画面灰白（颜色发白）

### 实现方案

在 `-filter_complex` 的 scale 滤镜后追加 tone mapping 滤镜链：

```
-filter_complex "[0:v]setpts=PTS-STARTPTS,scale=...,zscale=t=linear:npl=100,format=gbrpf32le,zscale=p=bt709,t=bt709,format=yuv420p"
```

或者用 ffmpeg 内置的 `tonemap` 滤镜（更简单）：
```
-tonemap hable  // 或 tonemap=reinhard/mobius
```

### 具体改动

```go
// 在构建 -filter_complex 时
filter := fmt.Sprintf("[0:v]setpts=PTS-STARTPTS,%s", scaleFilter)
if isHDR && !wantHDR {  // wantHDR 可以由 config 或 codec 判定
    filter += ",tonemap=hable:desat=0"
}
```

### 判定 `wantHDR`
- 如果最终编码器是 H.265 或 AV1 的 10-bit → 可以 wantHDR=true（保留 HDR）
- 如果最终编码器是 H.264 8-bit → wantHDR=false（必须 tone map）
- 可以新增配置项 `hdr_mode: passthrough | tonemap | auto`

---

## 阶段 4：配置与 UI（可选）

### 后端配置（`internal/config/transcoding.go`）
```go
type Transcoding struct {
    // ...现有字段...
    HdrMode string `mapstructure:"hdr_mode" json:"hdr_mode" yaml:"hdr_mode"` // "passthrough" | "tonemap" | "auto"
}
```

### 默认值建议
- 默认 `auto`：源 HDR + 编码器支持 10-bit → passthrough，否则 tonemap
- `passthrough`：强制保留 HDR，编码器不支持时回退
- `tonemap`：强制降级到 SDR

### DTO/VO/API
- 与 `useAv1`/`useH265` 模式一致，加 `hdrMode` 字段
- 前端加选择器（下拉框）

---

## 文件清单 & 改动量估算

| 文件 | 改动类型 | 估算行数 |
|------|---------|---------|
| `internal/global/transcoding.go` | 结构体扩展 | ~30 行 |
| `internal/service/transcoding.go` | 核心逻辑：HDR 判定 + 参数注入 | ~100 行 |
| `internal/service/transcoding.go` | 新增 `BuildColorParams`/`IsHDRSource` | ~50 行 |
| `internal/config/transcoding.go` | 可选：配置字段 | ~3 行 |
| `internal/service/config.go` | 可选：get/set | ~20 行 |
| `internal/domain/dto/config.go` | 可选：DTO 字段 | ~3 行 |
| `internal/domain/vo/config.go` | 可选：VO 字段 | ~3 行 |
| 前端 `OtherConfig.vue` | 可选：UI 选择器 | ~10 行 |

**核心必改**：3 个文件，~180 行

---

## 实现顺序 & PR 策略

```
PR1 (阶段 1+2)  ─── 元数据 + 编码器 HDR 注入（核心功能）
     │
PR2 (可选)      ─── Tone mapping + 配置项 + UI
```

- **PR1 完成后**：HDR 源视频自动在 H.265/AV1 编码中保留 HDR 色彩信息，播放器可正确识别
- **PR2 完成后**：可配置 HDR→SDR 降级，H.264 场景下画面不发白

---

## 风险 & 注意事项

1. **AV1 10-bit 兼容性**：`yuv420p10le` + AV1 编码需要较新的编码器版本（SVT-AV1 ≥ 0.9，NVENC AV1 驱动 ≥ 525），需要在环境中验证
2. **颜色参数位置**：FFmpeg 中 `-color_primaries` 等必须在编码器之后，否则可能被编码器覆盖
3. **Side data 精度**：x265 的 master-display 参数需要 50000 分度整数（如 `G(13250,34500)`），ffprobe 输出是小数（如 `0.708`），需要做 `int(val * 50000)` 转换
4. **HDR 检测阈值**：部分 HDR 源可能只有 `pix_fmt=yuv420p10le` 但没有 `color_primaries=bt2020`，可以考虑加 fallback：pix_fmt 是 10-bit + 宽度 ≥ 1920 → 可能是 HDR
5. **回退安全**：`IsHDRSource` 和 `BuildColorParams` 在元数据缺失时返回 false/nil，不影响现有 SDR 源

---

## 参考命令

```bash
# 探测 HDR 元数据
ffprobe -i input.mp4 -v quiet -print_format json -show_streams -show_format | jq '.streams[0] | {pix_fmt, color_primaries, color_transfer, color_space, color_range, side_data_list}'

# 保留 HDR 的转码命令示例（H.265 CPU）
ffmpeg -i input.mkv -c:v libx265 -pix_fmt yuv420p10le -preset slow \
  -color_primaries bt2020 -color_trc smpte2084 -colorspace bt2020nc \
  -x265-params "hdr10-opt=1:master-display=G(13250,34500)B(7500,3000)R(34000,16000)WP(15635,16450)L(10000000,1):max-cll=848,236" \
  output.mp4

# HDR→SDR tone mapping（hable 算法）
ffmpeg -i input.mkv -c:v libx264 -preset medium -crf 20 \
  -vf "tonemap=hable:desat=0,zscale=p=bt709:t=bt709:m=bt709" \
  -pix_fmt yuv420p \
  output.mp4

# HDR→SDR（使用 zscale 更精确控制）
ffmpeg -i input.mkv -c:v libx264 \
  -vf "zscale=t=linear:npl=100,format=gbrpf32le,zscale=p=bt709:t=bt709:m=bt709,format=yuv420p" \
  output.mp4
```
