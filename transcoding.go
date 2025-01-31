package main

import (
	"encoding/json"
	"os/exec"
	"io/ioutil"
	"log"
)

type Config struct {
	UseGPU bool `json:"use_gpu"`
}

type TranscodingTarget struct {
	Resolution  string
	BitrateRate string
	FPS         int
}

func getWidthRes(width int) int {
	if width >= 1920 {
		return 1080
	}
	if width >= 1280 {
		return 720
	}
	if width >= 854 {
		return 480
	}
	return 360
}

func getHeigthRes(height int) int {
	if height >= 1080 {
		return 1080
	}
	if height >= 720 {
		return 720
	}
	if height >= 480 {
		return 480
	}
	return 360
}

// 获取转码目标
func getTranscodingTarget(videoInfo *dto.TranscodingInfo) []TranscodingTarget {
	targets := make([]TranscodingTarget, 0)
	maxRresolution := utils.Max(getWidthRes(videoInfo.Width), getHeigthRes(videoInfo.Height))

	switch maxRresolution {
	case 1080:
		targets = append(targets, TranscodingTarget{Resolution: "1920x1080", BitrateRate: "3000k", FPS: videoInfo.FPS})
		fallthrough
	case 720:
		targets = append(targets, TranscodingTarget{Resolution: "1280x720", BitrateRate: "2000k", FPS: videoInfo.FPS})
		fallthrough
	case 480:
		targets = append(targets, TranscodingTarget{Resolution: "854x480", BitrateRate: "900k", FPS: videoInfo.FPS})
		fallthrough
	case 360:
		targets = append(targets, TranscodingTarget{Resolution: "640x360", BitrateRate: "500k", FPS: videoInfo.FPS})
	}

	return targets
}

// 获取视频信息
func getVideoInfo(input string) (info global.VideoInfo, err error) {
	cmd := exec.Command("ffprobe", "-i", input, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams")
	output, err := cmd.Output()
	if err != nil {
		return info, err
	}
	err = json.Unmarshal(output, &info)
	return info, err
}

// 读取配置文件
func readConfig() (Config, error) {
	var config Config
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}

// 执行转码
func transcode(input string, output string, target TranscodingTarget, useGPU bool) error {
	var cmd *exec.Cmd
	if useGPU {
		cmd = exec.Command("ffmpeg", "-i", input, "-c:v", "h264_nvenc", "-b:v", target.BitrateRate, "-s", target.Resolution, "-r", string(target.FPS), output)
	} else {
		cmd = exec.Command("ffmpeg", "-i", input, "-c:v", "libx264", "-b:v", target.BitrateRate, "-s", target.Resolution, "-r", string(target.FPS), output)
	}
	return cmd.Run()
}

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	videoInfo, err := getVideoInfo("input.mp4")
	if err != nil {
		log.Fatalf("获取视频信息失败: %v", err)
	}

	targets := getTranscodingTarget(&videoInfo)
	for _, target := range targets {
		err = transcode("input.mp4", "output_"+target.Resolution+".mp4", target, config.UseGPU)
		if err != nil {
			log.Fatalf("转码失败: %v", err)
		}
	}
}