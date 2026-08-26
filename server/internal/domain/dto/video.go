package dto

import "interastral-peace.com/alnitak/internal/domain/types"

type VideoListReq struct {
	Page     int
	PageSize int
}

type VideoFileReq struct {
	Hash              string  `json:"hash"`
	FileID            string  `json:"fileID"`
	Size              int64   `json:"size"`
	Cover             string  `json:"cover"`              // 前端截取的封面 objectKey，为空则 ffmpeg 截取
	Duration          float64 `json:"duration"`           // 视频时长（秒），前端传入
	Width             int     `json:"width"`              // 视频宽度
	Height            int     `json:"height"`             // 视频高度
	CodecName         string  `json:"codecName"`          // 视频编码名称
	ReplaceResourceID uint    `json:"replaceResourceID"` // >0 表示替换该旧资源
}

type VideoCheckResp struct {
	Chunks []int  `json:"chunks"`
	FileID string `json:"fileID"`
}

type ReviewListReq struct {
	Page     int
	PageSize int
}

type UploadVideoReq struct {
	Vid         uint   `json:"vid"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Desc        string `json:"desc"`
	Copyright   types.CopyrightType `json:"copyright"`
	Tags        string `json:"tags"`
	PartitionId uint   `json:"partitionId"`
}

type EditVideoReq struct {
	Vid   uint   `json:"vid"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Desc  string `json:"desc"`
	Tags  string `json:"tags"`
}

type SearchVideoReq struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
	KeyWords  string `json:"keywords"`
	Sort      string `json:"sort"`
	TimeRange string `json:"timeRange"`
}

// ── 直传 OSS 相关 DTO ──

// PresignImageReq 请求预签名图片上传 URL。
type PresignImageReq struct {
	FileName string `json:"fileName" binding:"required"`
	FileSize int64  `json:"fileSize" binding:"required"`
}

// PresignImageResp 返回预签名 URL 和 objectKey。
type PresignImageResp struct {
	PresignURL string `json:"presignURL"`
	ObjectKey  string `json:"objectKey"`
}

// ConfirmImageReq 确认图片已直传到 OSS。
type ConfirmImageReq struct {
	ObjectKey string `json:"objectKey" binding:"required"`
	Hash      string `json:"hash" binding:"required"`
}

// InitVideoUploadReq 发起视频分片直传。
type InitVideoUploadReq struct {
	Hash         string `json:"hash" binding:"required"`
	Size         int64  `json:"size" binding:"required"`
	FileName     string `json:"fileName" binding:"required"`
	TotalChunks  int    `json:"totalChunks" binding:"required"`
}

// InitVideoUploadResp 返回 uploadID、objectKey 和第一批分片的预签名 URL。
type InitVideoUploadResp struct {
	FileID         string              `json:"fileID"`
	UploadID       string              `json:"uploadID"`
	ObjectKey      string              `json:"objectKey"`
	TotalChunks    int                 `json:"totalChunks"`
	Chunks         []PresignChunkResp  `json:"chunks"`
	NextBatchStart int                 `json:"nextBatchStart"` // 下一批起始 index，-1 表示全部签完
}

// PresignUploadChunksReq 前端上传完一批分片后续签下一批。
type PresignUploadChunksReq struct {
	FileID  string `json:"fileID" binding:"required"`
	Start   int    `json:"start" binding:"required"`   // 起始分片 index
	Count   int    `json:"count" binding:"required"`    // 本次请求数量
}

// PresignUploadChunksResp 返回续签的分片预签名 URL。
type PresignUploadChunksResp struct {
	Chunks         []PresignChunkResp `json:"chunks"`
	NextBatchStart int                `json:"nextBatchStart"` // -1 表示全部签完
}

// PresignChunkResp 单个分片的预签名信息。
type PresignChunkResp struct {
	Index      int    `json:"index"`
	PartNumber int    `json:"partNumber"`
	PresignURL string `json:"presignURL"`
}

// CompleteVideoUploadReq 完成视频分片直传。
type CompleteVideoUploadReq struct {
	FileID   string         `json:"fileID" binding:"required"`
	UploadID string         `json:"uploadID" binding:"required"`
	Parts    []CompletePart `json:"parts" binding:"required"`
}

// CompletePart 表示一个已上传分片的 ETag。
type CompletePart struct {
	PartNumber int    `json:"partNumber"`
	ETag       string `json:"etag"`
}
