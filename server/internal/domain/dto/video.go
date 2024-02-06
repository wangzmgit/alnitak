package dto

type UploadVideoReq struct {
	Title       string
	Cover       string
	Desc        string
	Copyright   bool
	PartitionId uint //分区ID
}
