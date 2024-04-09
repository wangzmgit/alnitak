package dto

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
