package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

// ProgressFn 进度回调，pct 为 0-100。
type ProgressFn func(pct float64)

// EncodeVideo 执行 ffmpeg 视频编码，通过 progress 回调报告进度。
// 调用方需保证 totalDuration > 0 以启用进度报告。
// 返回 stderr 内容（用于错误诊断）和错误。
func EncodeVideo(ctx context.Context, args []string, totalDuration float64, progress ProgressFn) (stderrOut string, err error) {
	if len(args) == 0 {
		return "", fmt.Errorf("empty ffmpeg args")
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return stderr.String(), fmt.Errorf("ffmpeg start: %w", err)
	}

	// 异步读取 stdout 进度
	done := make(chan struct{})
	go func() {
		defer close(done)
		if totalDuration > 0 && progress != nil {
			watchProgress(stdout, totalDuration, progress)
		} else {
			// 不关心进度时消耗完 stdout 防止阻塞
			_, _ = io.Copy(io.Discard, stdout)
		}
	}()

	if err := cmd.Wait(); err != nil {
		<-done
		return stderr.String(), fmt.Errorf("ffmpeg: %w\n%s", err, stderr.String())
	}
	<-done
	return stderr.String(), nil
}

// EncodeAudio 执行 ffmpeg 音频编码（无进度报告，音频通常几秒完成）。
func EncodeAudio(ctx context.Context, args []string) error {
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("audio encode failed: %w, stderr: %s", err, stderr.String())
	}
	return nil
}

// watchProgress 解析 ffmpeg -progress pipe:1 输出并回调。
// 支持 out_time_ms= 和 out_time= 两种格式。
func watchProgress(r io.Reader, totalDuration float64, fn ProgressFn) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "out_time_ms="):
			raw := strings.TrimPrefix(line, "out_time_ms=")
			if ms, err := strconv.ParseFloat(raw, 64); err == nil {
				fn((ms / 1_000_000.0) / totalDuration * 100)
			}
		case strings.HasPrefix(line, "out_time="):
			raw := strings.TrimPrefix(line, "out_time=")
			if seconds := ParseFfmpegClockToSeconds(raw); seconds > 0 {
				fn(seconds / totalDuration * 100)
			}
		}
	}
}
