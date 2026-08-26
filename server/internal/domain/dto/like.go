package dto

type LikeVideoReq struct {
	Vid string `json:"vid" form:"vid"`
}

type LikeArticleReq struct {
	Aid uint `json:"aid" form:"aid"`
}
