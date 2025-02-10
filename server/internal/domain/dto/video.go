package dto

type VideoListReq struct {
	Page     int
	PageSize int
}

type VideoFileReq struct {
	Hash string
}

type ReviewListReq struct {
	Page     int
	PageSize int
}

type UploadVideoReq struct {
	Vid         uint
	Title       string
	Cover       string
	Desc        string
	Copyright   bool
	Tags        string
	PartitionId uint //分区ID
}

type EditVideoReq struct {
	Vid   uint
	Title string
	Cover string
	Desc  string
	Tags  string
}

type SearchVideoReq struct {
	Page     int
	PageSize int
	KeyWords string
}
