package service

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/spf13/viper"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

type TranscodingTarget struct {
	Resolution  string // 分辨率
	BitrateRate string // 码率
	FPS         string // 帧率
}

// 获取视频信息
func ProcessVideoInfo(input string) (*dto.TranscodingInfo, error) {
	var transcodingInfo dto.TranscodingInfo
	videoData, err := getVideoInfo(input)
	if err != nil {
		utils.ErrorLog("读取视频信息失败", "transcoding", err.Error())
		return &transcodingInfo, err
	}

	if videoData.Stream[0].CodecName != "h264" {
		utils.ErrorLog("视频编码不为h264", "transcoding", err.Error())
		return &transcodingInfo, errors.New("not h264")
	}

	//计算最大分辨率
	transcodingInfo.Width = videoData.Stream[0].Width
	transcodingInfo.Height = videoData.Stream[0].Height

	//获取视频时长
	transcodingInfo.Duration, _ = strconv.ParseFloat(videoData.Stream[0].Duration, 64)

	return &transcodingInfo, err
}

func VideoTransCoding(transcodingInfo *dto.TranscodingInfo) {
	var wg sync.WaitGroup
	targets := getTranscodingTarget(transcodingInfo)
	wg.Add(len(targets))
	for _, v := range targets {
		c := v // 处理协程引用循环变量问题
		go func() {
			fileName := c.Resolution + "_" + c.BitrateRate + "_" + c.FPS
			tsFileName := transcodingInfo.OutputDir + fileName + ".ts"
			// 压缩视频
			err := pressingVideo(transcodingInfo.InputFile, tsFileName, c.Resolution, c.BitrateRate)
			if err != nil {
				wg.Done()
				return
			}
			// 切片
			m3u8File, err := generateVideoSlices(tsFileName, transcodingInfo.OutputDir, fileName)
			if err != nil {
				wg.Done()
				return
			}
			// 读取m3u8写入数据库
			err = saveM3u8File(transcodingInfo, fileName, m3u8File)
			if err != nil {
				wg.Done()
				return
			}

			//删除临时文件
			os.Remove(tsFileName)
			os.Remove(m3u8File)

			wg.Done()
		}()
	}

	wg.Wait()

	// 上传oss
	if viper.GetString("storage.oss_type") != "local" {
		files, err := os.ReadDir(transcodingInfo.OutputDir)
		if err != nil {
			utils.ErrorLog("读取视频文件夹失败", "oss", err.Error())
			completeTransCoding(transcodingInfo.VideoID, transcodingInfo.ResourceID, global.PROCESSING_FAIL)
			return
		}

		for _, f := range files {
			if f.Name() == "upload.mp4" && !viper.GetBool("storage.upload_mp4_file") {
				continue
			}

			objectKey := "video/" + transcodingInfo.DirName + "/" + f.Name()
			filePath := "./upload/" + objectKey
			global.Storage.PutObjectFromFile(objectKey, filePath)
		}
	}

	// 更新状态
	completeTransCoding(transcodingInfo.VideoID, transcodingInfo.ResourceID, global.WAITING_REVIEW)
}

// 获取宽度支持的最大分辨率
func getWidthRes(width int) int {
	//1920*1080
	if width >= 1920 {
		return 1080
	}
	// 1280*720
	if width >= 1280 {
		return 720
	}
	// 720*480
	if width >= 720 {
		return 480
	}
	return 360
}

// 获取高度支持的最大分辨率
func getHeigthRes(height int) int {
	//1920*1080
	if height >= 1080 {
		return 1080
	}
	// 1280*720
	if height >= 720 {
		return 720
	}
	// 720*480
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
		targets = append(targets, TranscodingTarget{Resolution: "1920x1080", BitrateRate: "3000k", FPS: "30"})
		fallthrough
	case 720:
		targets = append(targets, TranscodingTarget{Resolution: "1080x720", BitrateRate: "2000k", FPS: "30"})
		fallthrough
	case 480:
		targets = append(targets, TranscodingTarget{Resolution: "854x480", BitrateRate: "900k", FPS: "30"})
		fallthrough
	case 360:
		targets = append(targets, TranscodingTarget{Resolution: "640x360", BitrateRate: "500k", FPS: "30"})
	}

	return targets
}

// 获取视频信息
func getVideoInfo(input string) (info global.VideoInfo, err error) {
	cmd := exec.Command("ffprobe", "-i", input, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams")
	out, err := utils.RunCmd(cmd)
	if err != nil {
		return info, err
	}

	// 反序列化
	err = json.Unmarshal(out.Bytes(), &info)
	if err != nil {
		return info, err
	}

	return info, nil
}

// 压缩视频
func pressingVideo(inputFile, outputFile, quality, rate string) error {
	command := []string{"-i", inputFile, "-crf", "20", "-s", quality,
		"-b:v", rate, "-c:v", "libx264", "-r", "30000/1001",
		"-c:a", "aac", "-f", "mpegts", outputFile,
	}

	_, err := utils.RunCmd(exec.Command("ffmpeg", command...))
	if err != nil {
		utils.ErrorLog("压缩视频失败", "transcoding", err.Error())
		return err
	}

	return nil
}

func generateVideoSlices(inputFile, outputDir, outputName string) (string, error) {
	outputM3U8 := outputDir + outputName + ".m3u8"
	outputTs := outputDir + outputName + "_%05d.ts"

	command := []string{"-i", inputFile, "-c", "copy",
		"-map", "0", "-f", "segment", "-segment_list",
		outputM3U8, "-segment_time", "10", outputTs,
	}

	_, err := utils.RunCmd(exec.Command("ffmpeg", command...))
	if err != nil {
		utils.ErrorLog("生成视频切片失败", "transcoding", err.Error())
		return outputM3U8, err
	}

	return outputM3U8, nil
}

// 保存m3u8文件
func saveM3u8File(transcodingInfo *dto.TranscodingInfo, fileName, m3u8File string) error {
	file, err := os.Open(m3u8File)
	if err != nil {
		utils.ErrorLog("打开m3u8文件失败", "transcoding", err.Error())
		return err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		utils.ErrorLog("读取m3u8文件失败", "transcoding", err.Error())
		return err
	}

	file.Close()

	global.Mysql.Create(&model.VideoIndexFile{
		ResourceID: transcodingInfo.ResourceID,
		Quality:    fileName,
		DirName:    transcodingInfo.DirName,
		Content:    string(bytes),
	})

	return nil
}

// 完成转码
func completeTransCoding(videoId, resourceId uint, status int) error {
	// 查询是否存在转码成功的视频文件
	var videoFileCount int64
	global.Mysql.Model(&model.VideoIndexFile{}).Where("resource_id = ?", resourceId).Count(&videoFileCount)
	if videoFileCount == 0 {
		status = global.PROCESSING_FAIL
	}

	// 更新资源状态
	if err := global.Mysql.Model(&model.Resource{}).Where("id = ?", resourceId).Updates(
		map[string]interface{}{
			"status": status,
		},
	).Error; err != nil {
		utils.ErrorLog("更新资源状态失败", "transcoding", err.Error())
		return err
	}

	// 获取转码中资源的数量
	var count int64
	global.Mysql.Model(&model.Resource{}).Where("vid = ? and status = ?", videoId, global.VIDEO_PROCESSING).Count(&count)
	// 如果没有转码中的视频，则更新视频为待审核
	if count == 0 {
		if err := global.Mysql.Model(&model.Video{}).Where("id = ? and status = ?", videoId, global.SUBMIT_REVIEW).Updates(
			map[string]interface{}{
				"status": global.WAITING_REVIEW,
			},
		).Error; err != nil {
			utils.ErrorLog("更新资源状态失败", "transcoding", err.Error())
			return err
		}
	}

	return nil
}
