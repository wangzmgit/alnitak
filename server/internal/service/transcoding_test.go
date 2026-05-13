package service

import (
	"testing"

	"interastral-peace.com/alnitak/internal/global"
)

func TestAvc1CodecStringFromH264ProfileLevel(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		level   int
		want    string
	}{
		{name: "High-4.2", profile: "High", level: 42, want: "avc1.64002A"},
		{name: "High-4.0", profile: "High", level: 40, want: "avc1.640028"},
		{name: "Main-4.2", profile: "Main", level: 42, want: "avc1.4D002A"},
		{name: "Baseline-3.1", profile: "baseline", level: 31, want: "avc1.42001F"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := avc1CodecStringFromH264ProfileLevel(tt.profile, tt.level)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("codec string = %q, want %q", got, tt.want)
			}
		})
	}

	t.Run("unsupported-profile", func(t *testing.T) {
		_, err := avc1CodecStringFromH264ProfileLevel("Unknown", 42)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}

func TestGetMaxQualityLevel(t *testing.T) {
	tests := []struct {
		name          string
		width, height int
		expected      int
	}{
		{"1080p标准", 1920, 1080, 1080},
		{"1080p非标（1920x1038）", 1920, 1038, 1080},
		{"1080p非标短边刚好1000", 1920, 1000, 1080},
		{"720p标准", 1280, 720, 720},
		{"480p标准", 854, 480, 480},
		{"360p标准", 640, 360, 360},
		{"低于360p", 320, 240, 360},
		{"竖屏1080p", 1080, 1920, 1080},
		{"竖屏720p", 720, 1280, 720},
		{"4K", 3840, 2160, 1080},
		{"接近1080p但短边不足1000", 1920, 960, 720},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMaxQualityLevel(tt.width, tt.height)
			if result != tt.expected {
				t.Errorf("getMaxQualityLevel(%d, %d) = %d, want %d",
					tt.width, tt.height, result, tt.expected)
			}
		})
	}
}

func TestCalcResolution(t *testing.T) {
	tests := []struct {
		name                          string
		srcWidth, srcHeight, targetSS int
		expectedW, expectedH         int
	}{
		{"16:9横屏→1080", 1920, 1080, 1080, 1920, 1080},
		{"16:9横屏→720", 1920, 1080, 720, 1280, 720},
		{"16:9横屏→360", 1920, 1080, 360, 640, 360},
		{"9:16竖屏→1080", 1080, 1920, 1080, 1080, 1920},
		{"4:3横屏→720", 1440, 1080, 720, 960, 720},
		{"偶数对齐", 1920, 1080, 480, 852, 480},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h := calcResolution(tt.srcWidth, tt.srcHeight, tt.targetSS)
			if w != tt.expectedW || h != tt.expectedH {
				t.Errorf("calcResolution(%d, %d, %d) = (%d, %d), want (%d, %d)",
					tt.srcWidth, tt.srcHeight, tt.targetSS, w, h, tt.expectedW, tt.expectedH)
			}
		})
	}
}

func TestScaleBitrateBySource(t *testing.T) {
	tests := []struct {
		name                                            string
		sourceMaxKbps, maxPresetKbps, currentPresetKbps int
		expected                                        int
	}{
		{"源码率充足", 10000, 8000, 5000, 6250},
		{"源码率不足", 3000, 8000, 5000, 1875},
		{"当前档等于最高档", 5000, 5000, 5000, 5000},
		{"源码率为0用默认", 0, 8000, 5000, 5000},
		{"源码率极低时受sourceMax限制", 100, 8000, 5000, 100},
		{"结果不超过源码率", 2000, 8000, 8000, 2000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scaleBitrateBySource(tt.sourceMaxKbps, tt.maxPresetKbps, tt.currentPresetKbps)
			if result != tt.expected {
				t.Errorf("scaleBitrateBySource(%d, %d, %d) = %d, want %d",
					tt.sourceMaxKbps, tt.maxPresetKbps, tt.currentPresetKbps, result, tt.expected)
			}
		})
	}
}

func TestParseQualityInfo(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantW  int
		wantH  int
		wantBW int
		wantFR float64
	}{
		{"标准1080p30", "1920x1080_8000k_30", 1920, 1080, 8000000, 30},
		{"标准720p60", "1280x720_5000k_60", 1280, 720, 5000000, 60},
		{"竖屏480p", "480x854_1500k_30", 480, 854, 1500000, 30},
		{"360p", "640x360_1000k_30", 640, 360, 1000000, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h, bw, fr := parseQualityInfo(tt.input)
			if w != tt.wantW || h != tt.wantH || bw != tt.wantBW || fr != tt.wantFR {
				t.Errorf("parseQualityInfo(%q) = (%d, %d, %d, %f), want (%d, %d, %d, %f)",
					tt.input, w, h, bw, fr, tt.wantW, tt.wantH, tt.wantBW, tt.wantFR)
			}
		})
	}
}

func TestParseFPS(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"30fps分数", "30000/1000", 30},
		{"59.94fps", "60000/1001", 60000.0 / 1001.0},
		{"23.976fps", "24000/1001", 24000.0 / 1001.0},
		{"整数", "30", 30},
		{"零除零", "0/0", 0},
		{"空字符串", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseFPS(tt.input)
			if result != tt.expected {
				t.Errorf("parseFPS(%q) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseFfmpegClockToSeconds(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"零", "00:00:00.000", 0},
		{"1分30秒", "00:01:30.500", 90.5},
		{"1小时", "01:00:00.000", 3600},
		{"无效格式", "invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseFfmpegClockToSeconds(tt.input)
			if result != tt.expected {
				t.Errorf("parseFfmpegClockToSeconds(%q) = %f, want %f",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestBuildRateControlParams(t *testing.T) {
	tests := []struct {
		rate       string
		wantTarget string
		wantMax    string
		wantBuf    string
	}{
		{"8000k", "8000k", "11200k", "22400k"},
		{"1000k", "1000k", "1400k", "2800k"},
		{"200k", "200k", "280k", "560k"},
	}

	for _, tt := range tests {
		t.Run(tt.rate, func(t *testing.T) {
			target, maxRate, bufSize := buildRateControlParams(tt.rate)
			if target != tt.wantTarget || maxRate != tt.wantMax || bufSize != tt.wantBuf {
				t.Errorf("buildRateControlParams(%q) = (%q, %q, %q), want (%q, %q, %q)",
					tt.rate, target, maxRate, bufSize, tt.wantTarget, tt.wantMax, tt.wantBuf)
			}
		})
	}
}

func TestParseBitrateKbps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"正常值", "8000k", 8000},
		{"小值", "500k", 500},
		{"零值", "0k", 0},
		{"空字符串", "", 0},
		{"非数字", "abck", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseBitrateKbps(tt.input)
			if result != tt.expected {
				t.Errorf("parseBitrateKbps(%q) = %d, want %d",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestFfmpegOutputDurationArgs(t *testing.T) {
	if got := ffmpegOutputDurationArgs(0); got != nil {
		t.Fatalf("expected nil for 0, got %#v", got)
	}
	got := ffmpegOutputDurationArgs(120.5)
	if len(got) != 2 || got[0] != "-t" || got[1] != "120.5" {
		t.Fatalf("120.5: got %#v", got)
	}
	got = ffmpegOutputDurationArgs(210.043121)
	if len(got) != 2 || got[0] != "-t" || got[1] != "210.043121" {
		t.Fatalf("210.043121: got %#v", got)
	}
}

func TestMinEncodeDurationSeconds(t *testing.T) {
	v := &global.Streams{Duration: "210.05"}
	a := &global.Streams{Duration: "209.90"}
	if m := minEncodeDurationSeconds("210.10", v, a); m != 209.90 {
		t.Fatalf("want 209.90 min, got %v", m)
	}
	if m := minEncodeDurationSeconds("", v, nil); m != 210.05 {
		t.Fatalf("want 210.05, got %v", m)
	}
}

func TestBFramePresentationLeadMs(t *testing.T) {
	if bFramePresentationLeadMs("30") != 67 {
		t.Fatalf("30fps: want 67ms, got %d", bFramePresentationLeadMs("30"))
	}
	if bFramePresentationLeadMs("60000/1001") != 33 {
		t.Fatalf("59.94: want 33ms, got %d", bFramePresentationLeadMs("60000/1001"))
	}
}

func TestAdelayPerChannelArg(t *testing.T) {
	if adelayPerChannelArg(0, 2) != "" {
		t.Fatal()
	}
	if adelayPerChannelArg(67, 2) != "67|67" {
		t.Fatalf("got %q", adelayPerChannelArg(67, 2))
	}
	if adelayPerChannelArg(50, 1) != "50" {
		t.Fatalf("got %q", adelayPerChannelArg(50, 1))
	}
}

