package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

// LocalTranscoder 包装现有的 TranscodeService，在本地进程内执行转码。
// 这是 "mode=local" 的默认实现，行为与重构前完全一致。
type LocalTranscoder struct {
	inner *TranscodeService
}

func NewLocalTranscoder() *LocalTranscoder {
	return &LocalTranscoder{inner: GetTranscoder()}
}

// Enqueue 启动 goroutine 执行转码，立即返回。
// 与原来调用方写 `go VideoTransCoding(info)` 等效。
func (t *LocalTranscoder) Enqueue(ctx context.Context, info *dto.TranscodingInfo) error {
	if info == nil {
		return fmt.Errorf("transcoding info is nil")
	}
	go func() {
		// 确保本地源文件可用：如果本地不存在则从 OSS 下载（兼容远程→本地切换场景）
		if err := ensureLocalSourceFile(info); err != nil {
			utils.ErrorLog("【本地转码】无法获取源文件，取消转码", "transcoding", err.Error())
			return
		}
		t.inner.ProcessVideo(ctx, info)
	}()
	return nil
}

// ensureLocalSourceFile 确保本地源文件存在。如果本地没有且 OSS 上有备份，自动下载。
// 解决场景：远程模式运行一段时间后切回本地模式，源文件已被清理，需从 OSS 恢复。
func ensureLocalSourceFile(info *dto.TranscodingInfo) error {
	if utils.IsFileExists(info.InputFile) {
		return nil // 本地已有，直接使用
	}

	// local 存储模式没有 OSS 备份，无法恢复，只能抛错
	if global.Config.Storage.OssType == "local" || global.Storage == nil {
		return fmt.Errorf("源文件不存在且无法从OSS下载（local存储）: %s", info.InputFile)
	}

	objectKey := fmt.Sprintf("video/%s/upload%s", info.DirName, info.Suffix)
	utils.InfoLog(fmt.Sprintf("【本地转码】源文件不存在，从OSS下载 key=%s -> %s", objectKey, info.InputFile), "transcoding")

	// 确保父目录存在
	parentDir := filepath.Dir(info.InputFile)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败 %s: %w", parentDir, err)
	}

	if err := global.Storage.GetObjectToFile(objectKey, info.InputFile); err != nil {
		return fmt.Errorf("从OSS下载源文件失败 key=%s: %w", objectKey, err)
	}

	utils.InfoLog(fmt.Sprintf("【本地转码】OSS源文件下载完成: %s", info.InputFile), "transcoding")
	return nil
}

// GetProgress 从 TranscodeService 的内存进度表查询指定 resource 的转码进度。
func (t *LocalTranscoder) GetProgress(ctx context.Context, resourceID uint) (*TranscodingProgress, error) {
	state := t.inner.getResourceProgress(resourceID)
	if state == nil {
		return nil, nil
	}

	qualities := make([]QualityProgress, 0, len(state.Details))
	var total float64
	for _, item := range state.Details {
		qualities = append(qualities, QualityProgress{
			Quality:  item.Quality,
			Progress: item.Progress,
			Status:   item.Status,
		})
		total += item.Progress
	}

	overall := float64(0)
	if len(qualities) > 0 {
		overall = total / float64(len(qualities))
	}

	status := "processing"
	uploadStatus := state.UploadStatus
	switch uploadStatus {
	case "success", "local":
		status = "success"
	case "fail":
		status = "fail"
	}

	var upload *UploadProgress
	if state.UploadStatus != "" {
		upload = &UploadProgress{
			OssType:  state.UploadOSS,
			Progress: state.UploadProgress,
			Status:   state.UploadStatus,
		}
	}

	return &TranscodingProgress{
		ResourceID:      resourceID,
		OverallProgress: overall,
		Status:          status,
		Qualities:       qualities,
		Upload:          upload,
	}, nil
}

// Cancel 通过 TranscodeService 中止指定 video 的所有转码进程并清理产物。
func (t *LocalTranscoder) Cancel(ctx context.Context, videoID uint) error {
	return t.inner.StopTranscodingAndCleanup(videoID)
}

// DispatchPending Local 模式无 pending 队列，空操作。
func (t *LocalTranscoder) DispatchPending(ctx context.Context) {}

// ListenDispatch Local 模式无 pending 队列，空操作。
func (t *LocalTranscoder) ListenDispatch(ctx context.Context) {}
