package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func UploadImg(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	suffix := path.Ext(file.Filename)
	fileName := utils.GenerateImgFilename(suffix)
	// 参数校验
	if !utils.IsImgType(suffix) { // 文件后缀
		return "", errors.New("文件类型错误")
	}

	//文件大小限制
	if !utils.FileSize(ctx.GetHeader("Content-Length"), viper.GetInt64("file.max_img_size")) {
		return "", errors.New("文件大小超出限制")
	}

	//保存文件
	if err := ctx.SaveUploadedFile(file, "./upload/image/"+fileName); err != nil {
		return "", errors.New("文件上传失败")
	}

	url := generateFileUrl("image/" + fileName)
	if viper.GetString("oss.type") != "local" {
		// TODO: 上传到OSS
	}

	// 缓存url
	userId := ctx.GetUint("userId")
	cache.SetUploadImage(url, userId)

	// TODO:记录日志

	return url, nil
}

func UploadVideo(ctx *gin.Context, vid uint, file *multipart.FileHeader) error {
	userId := ctx.GetUint("userId")
	suffix := path.Ext(file.Filename)
	// 参数校验
	if video, _ := FindVideoById(vid); video.ID == 0 || video.Uid != userId {
		return errors.New("视频不存在")
	}

	if !utils.IsVideoType(suffix) { // 文件后缀
		return errors.New("视频上传失败")
	}

	//文件大小限制
	if !utils.FileSize(ctx.GetHeader("Content-Length"), viper.GetInt64("file.max_video_size")) {
		return errors.New("文件大小超出限制")
	}

	//保存文件
	newFileName := utils.GenerateVideoFilename()
	if err := os.Mkdir("./upload/video/"+newFileName, os.ModePerm); err != nil {
		return errors.New("视频上传失败")
	}

	uploadVideoPath := "./upload/video/" + newFileName + "/upload.mp4"
	if err := ctx.SaveUploadedFile(file, uploadVideoPath); err != nil {
		return errors.New("文件上传失败")
	}

	transcodingInfo, err := ProcessVideoInfo(uploadVideoPath)
	if err != nil {
		return errors.New("读取视频信息失败")
	}

	// 存入数据库
	resource := &model.Resource{
		Vid:      vid,
		Uid:      userId,
		Title:    "",
		Status:   global.VIDEO_PROCESSING,
		Duration: transcodingInfo.Duration,
	}
	if err := global.Mysql.Create(resource).Error; err != nil {
		return errors.New("保存视频失败")
	}

	// 启动转码服务
	transcodingInfo.VideoID = vid
	transcodingInfo.DirName = newFileName
	transcodingInfo.ResourceID = resource.ID
	transcodingInfo.OutputDir = "./upload/video/" + newFileName + "/"
	transcodingInfo.InputFile = transcodingInfo.OutputDir + "upload.mp4"
	go VideoTransCoding(transcodingInfo)

	return nil
}

// 生成文件url
func generateFileUrl(objectKey string) string {
	if viper.GetString("oss.type") != "local" {

		if viper.GetString("oss.domain") != "" {
			return fmt.Sprintf("https://%s/%s", viper.GetString("oss.domain"), objectKey)
		} else {
			// 暂时只支持阿里云
			return fmt.Sprintf("https://%s.%s/%s",
				viper.GetString("oss.bucket"),
				viper.GetString("oss.endpoint"),
				objectKey,
			)
		}
	}

	return "/api/" + objectKey
}
