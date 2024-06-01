package dto

type ArticleListReq struct {
	Page     int
	PageSize int
}

type UploadArticleReq struct {
	Title       string
	Cover       string
	Copyright   bool
	Tags        string
	Content     string
	PartitionId uint //分区ID
}

type EditArticleReq struct {
	Aid     uint
	Title   string
	Cover   string
	Tags    string
	Content string
}

type ReviewArticleListReq struct {
	Page     int
	PageSize int
}
