package dto

type VideoListReq struct {
	Page     int
	PageSize int
}

type ReviewListReq struct {
	Page     int
	PageSize int
}

type UploadVideoReq struct {
	Title       string
	Cover       string
	Desc        string
	Copyright   bool
	Tags        string
	PartitionId uint //分区ID
}

type EditVideoReq struct {
	Vid       uint
	Title     string
	Cover     string
	Desc      string
	Tags      string
	Copyright bool
}
