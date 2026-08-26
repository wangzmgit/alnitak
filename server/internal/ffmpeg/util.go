package ffmpeg

import (
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

// ============================================================
// 帧率 / 时长解析
// ============================================================

// ParseFPS 将帧率字符串（"30000/1001"、"29.97"）转为 float64。
func ParseFPS(fps string) float64 {
	if fps == "" {
		return 0
	}
	parts := strings.Split(fps, "/")
	if len(parts) == 2 {
		num, err1 := strconv.ParseFloat(parts[0], 64)
		den, err2 := strconv.ParseFloat(parts[1], 64)
		if err1 == nil && err2 == nil && den > 0 {
			return num / den
		}
	}
	v, err := strconv.ParseFloat(fps, 64)
	if err != nil {
		return 0
	}
	return v
}

// ParseFfmpegClockToSeconds 解析 ffmpeg clock 格式 "01:23:45.678" → 秒。
func ParseFfmpegClockToSeconds(clock string) float64 {
	parts := strings.Split(clock, ":")
	if len(parts) != 3 {
		return 0
	}
	h, _ := strconv.ParseFloat(parts[0], 64)
	m, _ := strconv.ParseFloat(parts[1], 64)
	s, _ := strconv.ParseFloat(parts[2], 64)
	return h*3600 + m*60 + s
}

// FfmpegOutputDurationArgs 生成 -t duration 参数。
func FfmpegOutputDurationArgs(seconds float64) []string {
	if seconds <= 0 || math.IsNaN(seconds) || math.IsInf(seconds, 0) {
		return nil
	}
	s := strings.TrimRight(strings.TrimRight(strconv.FormatFloat(seconds, 'f', 6, 64), "0"), ".")
	if s == "" {
		return nil
	}
	return []string{"-t", s}
}

// ============================================================
// 码率解析
// ============================================================

// ParseBitrateKbps 将 "8000k" "5000" 转为 int(kbps)。
func ParseBitrateKbps(rate string) int {
	rate = strings.TrimSpace(strings.TrimSuffix(rate, "k"))
	if v, err := strconv.Atoi(rate); err == nil && v > 0 {
		return v
	}
	return 0
}

// FormatBitrateKbps 将 int(kbps) 转为 "8000k"。
func FormatBitrateKbps(rateKbps int) string {
	if rateKbps < 1 {
		rateKbps = 1
	}
	return fmt.Sprintf("%dk", rateKbps)
}

// ============================================================
// B 帧对齐
// ============================================================

// BFramePresentationLeadMs 按帧率估计视频首帧展示延迟（ms）。
func BFramePresentationLeadMs(fps string) int {
	f := ParseFPS(fps)
	if f <= 0 {
		return 0
	}
	ms := int(math.Round(2000.0 / f))
	if ms < 1 {
		return 0
	}
	if ms > 200 {
		return 200
	}
	return ms
}

// AdelayPerChannelArg 生成 adelay 声道参数 "del0|del1|…"。
func AdelayPerChannelArg(delayMs, channels int) string {
	if delayMs <= 0 || channels < 1 {
		return ""
	}
	s := strconv.Itoa(delayMs)
	parts := make([]string, channels)
	for i := range parts {
		parts[i] = s
	}
	return strings.Join(parts, "|")
}

// ============================================================
// 分辨率 / 画质计算
// ============================================================

// ParseQualityInfo 从 quality name 解析宽、高、帧率、码率。
func ParseQualityInfo(quality string) (width, height int, bandwidth int, frameRate string) {
	// "1920x1080_8000k_30"
	parts := strings.Split(quality, "_")
	if len(parts) != 3 {
		return 0, 0, 0, ""
	}
	res := strings.Split(parts[0], "x")
	if len(res) == 2 {
		width, _ = strconv.Atoi(res[0])
		height, _ = strconv.Atoi(res[1])
	}
	bandwidth = ParseBitrateKbps(parts[1]) * 1000
	frameRate = parts[2]
	return
}

// ParseQualitySortKey 解析 quality name 用于排序（高度 > 宽度 > 帧率 > 码率）。
func parseQualitySortKey(quality string) (w, h, fps, br int) {
	parts := strings.Split(quality, "_")
	if len(parts) < 3 {
		return
	}
	res := strings.Split(parts[0], "x")
	if len(res) == 2 {
		w, _ = strconv.Atoi(res[0])
		h, _ = strconv.Atoi(res[1])
	}
	br = ParseBitrateKbps(parts[1])
	fpsVal := ParseFPS(parts[2])
	fps = int(math.Round(fpsVal))
	return
}

// ============================================================
// Quality Preset Helpers
// ============================================================

// GetMaxQualityLevel 根据原始分辨率返回最高支持的画质档位（短边高度）。
func GetMaxQualityLevel(width, height int) int {
	shortSide, longSide := width, height
	if height < width {
		shortSide, longSide = height, width
	}
	if longSide >= 1920 && shortSide >= 1000 {
		return 1080
	}
	if shortSide >= 1080 {
		return 1080
	}
	if shortSide >= 720 {
		return 720
	}
	if shortSide >= 480 {
		return 480
	}
	return 360
}

// CalcResolution 按目标短边计算等比例缩放后的实际分辨率（偶数对齐）。
func CalcResolution(srcWidth, srcHeight, targetShortSide int) (w, h int) {
	isPortrait := srcWidth < srcHeight
	srcAspect := float64(srcWidth) / float64(srcHeight)
	if isPortrait {
		w = targetShortSide
		h = int(float64(w) / srcAspect)
	} else {
		h = targetShortSide
		w = int(float64(h) * srcAspect)
	}
	return w &^ 1, h &^ 1
}

// GetPresetBitrateKbps 从 preset 中取码率（kbps），优先 60fps 值，fallback 到 30fps。
func GetPresetBitrateKbps(preset QualityPreset, isPortrait, fps60 bool) int {
	var bitrate string
	if fps60 {
		bitrate = preset.Bitrate60H
		if isPortrait {
			bitrate = preset.Bitrate60V
		}
		if kbps := ParseBitrateKbps(bitrate); kbps > 0 {
			return kbps
		}
	}
	bitrate = preset.BitrateH
	if isPortrait {
		bitrate = preset.BitrateV
	}
	return ParseBitrateKbps(bitrate)
}

// ScaleBitrateBySource 根据源视频码率缩放预设码率（保留动态范围）。
func ScaleBitrateBySource(sourceMaxKbps, maxPresetKbps, currentPresetKbps int) int {
	if sourceMaxKbps <= 0 {
		return currentPresetKbps
	}
	if maxPresetKbps <= 0 || currentPresetKbps <= 0 {
		return sourceMaxKbps
	}
	if currentPresetKbps >= maxPresetKbps {
		return sourceMaxKbps
	}
	scaled := int(math.Round(float64(sourceMaxKbps) * float64(currentPresetKbps) / float64(maxPresetKbps)))
	if scaled < 200 {
		scaled = 200
	}
	if scaled > sourceMaxKbps {
		scaled = sourceMaxKbps
	}
	return scaled
}

// GetDefaultVideoBitRateByLevel 按最高画质档位返回默认视频码率（bps）。
func GetDefaultVideoBitRateByLevel(maxLevel int, isPortrait bool) int {
	for _, preset := range QualityPresets {
		if preset.ShortSide == maxLevel {
			br := preset.BitrateH
			if isPortrait {
				br = preset.BitrateV
			}
			if kbps := ParseBitrateKbps(br); kbps > 0 {
				return kbps * 1000
			}
		}
	}
	return 1000000
}

// AudioFileNameForTrack 按语言代码生成音轨文件名。
// "und" 或空返回 "audio.m4s"，否则返回 "audio_{lang}.m4s"。
func AudioFileNameForTrack(language string) string {
	if language == "" || language == "und" {
		return "audio.m4s"
	}
	return "audio_" + language + ".m4s"
}

// GetFpsInfo 从平均帧率解析出 30fps/60fps 的 timebase。
func GetFpsInfo(avgFrameRate string, enable60 bool) (string, string) {
	parts := strings.Split(avgFrameRate, "/")
	if len(parts) == 2 {
		num, _ := strconv.Atoi(parts[0])
		den, _ := strconv.Atoi(parts[1])
		if den == 0 {
			return TimeBase30fps, ""
		}
		fps := float64(num) / float64(den)
		if fps < 30 {
			return avgFrameRate, ""
		}
		if fps >= 59 {
			if enable60 {
				return TimeBase30fps, TimeBase60fps
			}
			return TimeBase30fps, ""
		}
	}
	return TimeBase30fps, ""
}

// ============================================================
// codec string
// ============================================================

// TranscodingTarget 转码目标（清晰度档位）
type TranscodingTarget struct {
	Resolution  string
	BitrateRate string
	FPS         string
	FpsName     string
}

// GetTranscodingTargets 根据输入视频参数生成所有转码目标档位。
func GetTranscodingTargets(width, height, videoBitRate int, fps30, fps60 string, enable60 bool) []TranscodingTarget {
	targets := make([]TranscodingTarget, 0)
	maxLevel := GetMaxQualityLevel(width, height)
	isPortrait := width < height

	var maxPreset QualityPreset
	for _, p := range QualityPresets {
		if p.ShortSide == maxLevel {
			maxPreset = p
			break
		}
	}
	if maxPreset.ShortSide == 0 {
		return targets
	}

	maxPresetKbps30 := GetPresetBitrateKbps(maxPreset, isPortrait, false)
	maxPresetKbps60 := GetPresetBitrateKbps(maxPreset, isPortrait, true)
	if maxPresetKbps60 <= 0 {
		maxPresetKbps60 = maxPresetKbps30
	}

	sourceMaxKbps := videoBitRate / 1000
	if sourceMaxKbps <= 0 {
		sourceMaxKbps = maxPresetKbps30
	}

	gen60 := enable60 && fps60 != ""

	for _, preset := range QualityPresets {
		if preset.ShortSide > maxLevel {
			continue
		}
		w, h := CalcResolution(width, height, preset.ShortSide)
		resStr := fmt.Sprintf("%dx%d", w, h)

		currKbps30 := GetPresetBitrateKbps(preset, isPortrait, false)
		dynKbps30 := ScaleBitrateBySource(sourceMaxKbps, maxPresetKbps30, currKbps30)

		if preset.ShortSide == maxLevel && gen60 {
			currKbps60 := GetPresetBitrateKbps(preset, isPortrait, true)
			dynKbps60 := ScaleBitrateBySource(sourceMaxKbps, maxPresetKbps60, currKbps60)
			targets = append(targets, TranscodingTarget{Resolution: resStr, BitrateRate: FormatBitrateKbps(dynKbps60), FPS: fps60, FpsName: "60"})
		}
		targets = append(targets, TranscodingTarget{Resolution: resStr, BitrateRate: FormatBitrateKbps(dynKbps30), FPS: fps30, FpsName: "30"})
	}
	return targets
}

// Avc1CodecString 从 ffprobe profile+level 生成 avc1 编码字符串。
func Avc1CodecString(profile string, level int) string {
	// https://stackoverflow.com/questions/24834877/avc1-codec-string
	profileMap := map[string]int{
		"baseline": 66,
		"main":     77,
		"high":     100,
	}
	pi, ok := profileMap[strings.ToLower(profile)]
	if !ok {
		return ""
	}
	return fmt.Sprintf("avc1.%02X%02X%02X", pi, 0, level)
}

// ProbeVideoActualCodec 从已编码的视频文件中探测实际使用的编码器，
// 返回标准 codec 字符串（avc1.xxx / hvc1.xxx / av01.xxx）。
// 用于编码降级后正确记录 codec，而非依赖配置。
func ProbeVideoActualCodec(videoFile string) string {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0",
		"-show_entries", "stream=codec_name,profile,level",
		"-of", "default=nw=1:nk=1", videoFile)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return DefaultVideoCodec
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) < 3 {
		return DefaultVideoCodec
	}

	codecName := strings.TrimSpace(lines[0])
	profile := strings.TrimSpace(lines[1])
	level, _ := strconv.Atoi(strings.TrimSpace(lines[2]))

	switch strings.ToLower(codecName) {
	case "av1":
		return DefaultAV1Codec
	case "hevc", "h265":
		return DefaultHEVCCodec
	case "h264", "avc":
		if cs := Avc1CodecString(profile, level); cs != "" {
			return cs
		}
		return DefaultVideoCodec
	default:
		return DefaultVideoCodec
	}
}
