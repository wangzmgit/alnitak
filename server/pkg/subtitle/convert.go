package subtitle

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	// ErrInvalidSRT 无法解析的 SRT
	ErrInvalidSRT = errors.New("无效的 SRT 字幕")
	// ErrInvalidVTT 无法作为字幕使用的 WebVTT
	ErrInvalidVTT = errors.New("无效的 WebVTT 字幕")
)

// SRTToWebVTT 将 SRT 正文转为 WebVTT（UTF-8）
func SRTToWebVTT(srt []byte) ([]byte, error) {
	text := string(bytes.TrimPrefix(srt, []byte{0xEF, 0xBB, 0xBF}))
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, ErrInvalidSRT
	}

	blocks := strings.Split(text, "\n\n")
	var out strings.Builder
	out.WriteString("WEBVTT\n\n")
	cues := 0
	for _, block := range blocks {
		lines := strings.Split(strings.TrimSpace(block), "\n")
		if len(lines) < 2 {
			continue
		}
		idx := 0
		if _, err := strconv.Atoi(strings.TrimSpace(lines[0])); err == nil && len(lines) >= 2 {
			idx = 1
		}
		if idx >= len(lines) {
			continue
		}
		timeLine := strings.TrimSpace(lines[idx])
		start, end, ok := parseSRTTimeLine(timeLine)
		if !ok {
			continue
		}
		textLines := lines[idx+1:]
		if len(textLines) == 0 {
			continue
		}
		fmt.Fprintf(&out, "%s --> %s\n", start, end)
		for _, tl := range textLines {
			out.WriteString(strings.TrimRight(tl, "\r"))
			out.WriteByte('\n')
		}
		out.WriteByte('\n')
		cues++
	}
	if cues == 0 {
		return nil, ErrInvalidSRT
	}
	return []byte(out.String()), nil
}

func parseSRTTimeLine(line string) (start, end string, ok bool) {
	parts := strings.Split(line, "-->")
	if len(parts) != 2 {
		return "", "", false
	}
	s := strings.TrimSpace(strings.ReplaceAll(parts[0], ",", "."))
	e := strings.TrimSpace(strings.ReplaceAll(parts[1], ",", "."))
	// 去掉可能附带的 V4 样式坐标
	if i := strings.IndexByte(e, ' '); i > 0 {
		// WebVTT 可保留部分 cue settings；SRT 行尾 X1 等需剥离
		head := e[:i]
		if isTimeToken(head) {
			e = head
		}
	}
	if !isTimeToken(s) || !isTimeToken(e) {
		return "", "", false
	}
	return s, e, true
}

func isTimeToken(t string) bool {
	// 00:00:00.000 或 00:00:00,000（调用前已统一为点）
	if len(t) < 12 {
		return false
	}
	if t[2] != ':' || t[5] != ':' || (t[8] != '.' && t[8] != ',') {
		return false
	}
	return true
}

// ValidateOrNormalizeWebVTT 校验上传的 WebVTT，必要时补充 WEBVTT 头
func ValidateOrNormalizeWebVTT(raw []byte) ([]byte, error) {
	text := string(bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF}))
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if utf8.ValidString(text) == false {
		return nil, ErrInvalidVTT
	}
	trim := strings.TrimSpace(text)
	if trim == "" {
		return nil, ErrInvalidVTT
	}
	if !strings.HasPrefix(strings.TrimLeft(text, "\n"), "WEBVTT") {
		text = "WEBVTT\n\n" + text
	}
	// 至少应有一条时间轴行包含 -->
	if !strings.Contains(text, "-->") {
		return nil, ErrInvalidVTT
	}
	return []byte(text), nil
}
