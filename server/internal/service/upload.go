package service

import (
	"errors"
	"mime/multipart"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func UploadImg(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	suffix := path.Ext(file.Filename)
	fileName := generateImgFilename(suffix)

	objectKey := "image/" + fileName
	filePath := "./upload/image/" + fileName
	// 参数校验
	if !utils.IsImgType(suffix) { // 文件后缀
		return "", errors.New("文件类型错误")
	}

	//文件大小限制
	if !utils.FileSize(ctx.GetHeader("Content-Length"), viper.GetInt64("file.max_img_size")) {
		return "", errors.New("文件大小超出限制")
	}

	//保存文件
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		return "", errors.New("文件上传失败")
	}

	url := generateFileUrl(objectKey)
	if viper.GetString("storage.oss_type") != "local" {
		// 上传到OSS
		global.Storage.PutObjectFromFile(objectKey, filePath)
	}

	// 缓存url
	userId := ctx.GetUint("userId")
	cache.SetUploadImage(url, userId)

	return url, nil
}

func UploadVideoCreate(ctx *gin.Context, file *multipart.FileHeader) (vo.ResourceResp, error) {
	userId := ctx.GetUint("userId")

	videoName, err := SaveUploadVideo(ctx, file)
	if err != nil {
		return vo.ResourceResp{}, err
	}

	uploadVideoPath := "./upload/video/" + videoName + "/upload.mp4"
	vid, _ := initVideo(userId, uploadVideoPath, file.Filename)
	if vid == 0 {
		return vo.ResourceResp{}, errors.New("创建失败")
	}

	resource, err := CompleteUploadVideo(vid, userId, videoName, file.Filename)
	if err != nil {
		return vo.ResourceResp{}, err
	}

	return resource, nil
}

func UploadVideoAdd(ctx *gin.Context, vid uint, file *multipart.FileHeader) (vo.ResourceResp, error) {
	userId := ctx.GetUint("userId")
	if video, _ := FindVideoById(vid); video.ID == 0 || video.Uid != userId {
		return vo.ResourceResp{}, errors.New("视频不存在")
	}

	videoName, err := SaveUploadVideo(ctx, file)
	if err != nil {
		return vo.ResourceResp{}, err
	}

	resource, err := CompleteUploadVideo(vid, userId, videoName, file.Filename)
	if err != nil {
		return vo.ResourceResp{}, err
	}

	return resource, nil
}

func SaveUploadVideo(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	suffix := path.Ext(file.Filename)
	if !utils.IsVideoType(suffix) { // 文件后缀
		return "", errors.New("视频上传失败")
	}

	//文件大小限制
	if !utils.FileSize(ctx.GetHeader("Content-Length"), viper.GetInt64("file.max_video_size")) {
		return "", errors.New("文件大小超出限制")
	}

	//保存文件
	videoName := generateVideoFilename()
	if err := os.Mkdir("./upload/video/"+videoName, os.ModePerm); err != nil {
		return "", errors.New("视频上传失败")
	}

	uploadVideoPath := "./upload/video/" + videoName + "/upload.mp4"
	if err := ctx.SaveUploadedFile(file, uploadVideoPath); err != nil {
		return "", errors.New("文件上传失败")
	}

	return videoName, nil
}

func CompleteUploadVideo(vid, userId uint, videoName, title string) (vo.ResourceResp, error) {
	uploadVideoPath := "./upload/video/" + videoName + "/upload.mp4"
	transcodingInfo, err := ProcessVideoInfo(uploadVideoPath)
	if err != nil {
		return vo.ResourceResp{}, errors.New("读取视频信息失败")
	}

	// 存入数据库
	resource := model.Resource{
		Vid:      vid,
		Uid:      userId,
		Title:    title,
		Status:   global.VIDEO_PROCESSING,
		Duration: transcodingInfo.Duration,
	}
	if err := global.Mysql.Create(&resource).Error; err != nil {
		return vo.ResourceResp{}, errors.New("保存视频失败")
	}

	// 启动转码服务
	transcodingInfo.VideoID = vid
	transcodingInfo.DirName = videoName
	transcodingInfo.ResourceID = resource.ID
	transcodingInfo.OutputDir = "./upload/video/" + videoName + "/"
	transcodingInfo.InputFile = transcodingInfo.OutputDir + "upload.mp4"
	go VideoTransCoding(transcodingInfo)

	return vo.ResourceToResourceResp(resource), nil
}

// 生成文件url
func generateFileUrl(objectKey string) string {
	if viper.GetString("storage.oss_type") != "local" {
		global.Storage.GetObjectUrl(objectKey)
	}

	return "/api/" + objectKey
}

// 初始化视频
func initVideo(userId uint, videoPath, title string) (uint, error) {
	// 生成封面
	coverName := generateImgFilename(".jpg")
	objectKey := "image/" + coverName
	filePath := "./upload/image/" + coverName

	GenerateCover(videoPath, filePath)
	if viper.GetString("storage.oss_type") != "local" {
		// 上传到OSS
		global.Storage.PutObjectFromFile(objectKey, filePath)
	}

	videoId, err := CreateVideo(&model.Video{
		Uid:    userId,
		Cover:  generateFileUrl(objectKey),
		Title:  title,
		Status: global.CREATED_VIDEO,
	})
	if err != nil {
		return 0, err
	}

	return videoId, nil
}

// 随机生成图片文件名
func generateImgFilename(suffix string) string {
	id := global.SnowflakeNode.Generate()
	return id.String() + suffix
}

// 随机视频文件名
func generateVideoFilename() string {
	id := global.SnowflakeNode.Generate()
	return id.String()
}
