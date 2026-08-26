package dto

type ShareVideoReq struct {
	Vid string `json:"vid" binding:"required"`
}

type ShareArticleReq struct {
	Aid uint `json:"aid" binding:"required"`
}
