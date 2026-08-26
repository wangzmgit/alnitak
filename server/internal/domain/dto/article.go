package dto

import "interastral-peace.com/alnitak/internal/domain/types"

type ArticleListReq struct {
	Page     int
	PageSize int
}

type UploadArticleReq struct {
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Copyright   types.CopyrightType `json:"copyright"`
	Tags        string `json:"tags"`
	Content     string `json:"content"`
	PartitionId uint   `json:"partitionId"`
}

type EditArticleReq struct {
	Aid     uint   `json:"aid"`
	Title   string `json:"title"`
	Cover   string `json:"cover"`
	Tags    string `json:"tags"`
	Content string `json:"content"`
}

type ReviewArticleListReq struct {
	Page     int
	PageSize int
}
