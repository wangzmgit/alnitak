package ffmpeg

import (
	"fmt"
	"math"
	"strconv"
)

// VideoEncodeArgs 构建视频编码的 ffmpeg 参数列表。
// 调用方需自行追加 -threads N -y outputFile 及 -progress pipe:1 -nostats。
func VideoEncodeArgs(inputFile, quality, rate, fps string, totalDuration float64, useGpu, useAv1, useHevc bool) []string {
	fpsFloat := ParseFPS(fps)
	gopSize := int(math.Round(fpsFloat * 2))
	if gopSize < 1 {
		gopSize = 60
	}
	gopSizeStr := strconv.Itoa(gopSize)
	targetRate, maxrate, bufsize := BuildRateControlParams(rate)
	// fps 过滤器强制精确 CFR
	scaleFilter := fmt.Sprintf("fps=%s,scale=%s:flags=lanczos", fps, quality)

	// CPU/GPU 线程策略在调用方决定，这里只构建编码参数
	args := []string{
		"-i", inputFile,
		"-filter_complex", fmt.Sprintf("[0:v]setpts=PTS-STARTPTS,%s", scaleFilter),
		"-an",
	}

	switch {
	case useGpu && useAv1:
		args = append(args,
			"-c:v", "av1_nvenc", "-cq", "30", "-preset", "p7", "-rc", "vbr",
			"-pix_fmt", "yuv420p", "-bf", "2",
			"-b_ref_mode", "middle", "-multipass", "fullres",
			"-aq-strength", "8",
			"-rc-lookahead", "32",
		)
	case useGpu && useHevc:
		args = append(args,
			"-c:v", "hevc_nvenc", "-cq", "24", "-preset", "p7", "-rc", "vbr",
			"-profile:v", "main10", "-pix_fmt", "yuv420p10le", "-bf", "2", "-b_ref_mode", "middle",
			"-forced-idr", "1", "-multipass", "fullres",
			"-spatial_aq", "1", "-temporal_aq", "1", "-aq-strength", "8",
			"-rc-lookahead", "32",
		)
	case useGpu && !useHevc:
		args = append(args,
			"-c:v", "h264_nvenc", "-cq", "23", "-preset", "p7", "-rc", "vbr",
			"-profile:v", "high", "-pix_fmt", "yuv420p", "-bf", "2", "-b_ref_mode", "middle",
			"-forced-idr", "1", "-multipass", "fullres",
			"-spatial_aq", "1", "-temporal_aq", "1", "-aq-strength", "8",
			"-rc-lookahead", "32",
		)
	case !useGpu && useAv1:
		args = append(args,
			"-c:v", "libsvtav1", "-preset", "6", "-crf", "30", "-tag:v", "av01",
			"-pix_fmt", "yuv420p",
		)
	case !useGpu && useHevc:
		args = append(args,
			"-c:v", "libx265", "-preset", "slow", "-tag:v", "hvc1",
			"-crf", "22",
			"-pix_fmt", "yuv420p10le",
			"-x265-params", "profile=main10:aq-mode=3:aq-strength=0.8:deblock=-1,-1:no-sao=1",
		)
	default:
		args = append(args,
			"-c:v", "libx264", "-preset", "medium", "-crf", "20", "-tune", "film",
			"-profile:v", "high", "-pix_fmt", "yuv420p", "-bf", "3", "-b_strategy", "2",
			"-flags", "+cgop", "-sc_threshold", "0",
			"-refs", "6", "-me_method", "hex", "-subq", "9",
		)
	}

	// SVT-AV1 CRF 模式不接受 -b:v
	if !useGpu && useAv1 {
		args = append(args,
			"-r", fps, "-g", gopSizeStr, "-keyint_min", gopSizeStr,
		)
	} else {
		args = append(args,
			"-b:v", targetRate, "-maxrate", maxrate, "-bufsize", bufsize,
			"-r", fps, "-g", gopSizeStr, "-keyint_min", gopSizeStr,
		)
	}

	if useGpu {
		args = append(args, "-strict_gop", "1")
	}

	args = append(args,
		"-fps_mode", "cfr",
		"-f", "mp4",
		"-frag_duration", FragDurationUs,
		"-movflags", "+frag_keyframe+empty_moov+default_base_moof+dash+global_sidx+negative_cts_offsets",
		"-avoid_negative_ts", "make_zero",
	)

	if dur := FfmpegOutputDurationArgs(totalDuration); dur != nil {
		args = append(args, dur...)
	}
	// 调用方自行追加 -threads N -y outputFile
	return args
}

// AudioEncodeArgs 构建单音轨 ffmpeg 参数。
// 调用方自行追加 -y outputFile。
func AudioEncodeArgs(inputFile string, bitRate, sampleRate, channels int, durationSec float64, presentationLeadMs int) []string {
	bitRateStr := fmt.Sprintf("%dk", bitRate/1000)
	sampleRateStr := strconv.Itoa(sampleRate)
	channelsStr := strconv.Itoa(channels)

	adelayArg := AdelayPerChannelArg(presentationLeadMs, channels)
	audioFilter := fmt.Sprintf("[0:a]asetpts=PTS-STARTPTS,aresample=osr=%s", sampleRateStr)
	if adelayArg != "" {
		audioFilter += ",adelay=" + adelayArg
	}
	audioFilter += "[aout]"
	audioOutputDur := durationSec
	if presentationLeadMs > 0 {
		audioOutputDur += float64(presentationLeadMs) / 1000.0
	}

	args := []string{
		"-i", inputFile,
		"-filter_complex", audioFilter,
		"-map", "[aout]", "-vn", "-c:a", "aac", "-b:a", bitRateStr,
		"-ar", sampleRateStr, "-ac", channelsStr,
		"-f", "mp4",
		"-frag_duration", FragDurationUs,
		"-movflags", "+frag_keyframe+empty_moov+default_base_moof+dash+global_sidx+negative_cts_offsets",
		"-avoid_negative_ts", "make_zero",
	}
	if dur := FfmpegOutputDurationArgs(audioOutputDur); dur != nil {
		args = append(args, dur...)
	}
	// 调用方自行追加 -y outputFile
	return args
}

// AudioTrackEncodeArgs 构建多音轨模式下指定音轨的 ffmpeg 参数。
// 调用方自行追加 -y outputFile。
func AudioTrackEncodeArgs(inputFile string, streamIndex, bitRate, sampleRate, channels int, durationSec float64, presentationLeadMs int) []string {
	bitRateStr := fmt.Sprintf("%dk", bitRate/1000)
	sampleRateStr := strconv.Itoa(sampleRate)
	channelsStr := strconv.Itoa(channels)

	adelayArg := AdelayPerChannelArg(presentationLeadMs, channels)
	audioFilter := fmt.Sprintf("[0:a:%d]asetpts=PTS-STARTPTS,aresample=osr=%s", streamIndex, sampleRateStr)
	if adelayArg != "" {
		audioFilter += ",adelay=" + adelayArg
	}
	audioFilter += "[aout]"
	audioOutputDur := durationSec
	if presentationLeadMs > 0 {
		audioOutputDur += float64(presentationLeadMs) / 1000.0
	}

	args := []string{
		"-i", inputFile,
		"-filter_complex", audioFilter,
		"-map", "[aout]", "-vn", "-c:a", "aac", "-b:a", bitRateStr,
		"-ar", sampleRateStr, "-ac", channelsStr,
		"-f", "mp4",
		"-frag_duration", FragDurationUs,
		"-movflags", "+frag_keyframe+empty_moov+default_base_moof+dash+global_sidx+negative_cts_offsets",
		"-avoid_negative_ts", "make_zero",
	}
	if dur := FfmpegOutputDurationArgs(audioOutputDur); dur != nil {
		args = append(args, dur...)
	}
	// 调用方自行追加 -y outputFile
	return args
}

// BuildRateControlParams 按目标码率字符串生成 -b:v / -maxrate / -bufsize 参数值。
func BuildRateControlParams(rate string) (targetRate, maxrate, bufsize string) {
	kbps := ParseBitrateKbps(rate)
	if kbps <= 0 {
		return rate, rate, rate
	}
	max := int(float64(kbps) * 1.2)
	buf := int(float64(kbps) * 2.0)
	return rate, FormatBitrateKbps(max), FormatBitrateKbps(buf)
}
