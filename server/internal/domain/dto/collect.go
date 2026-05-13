package dto

type AddCollectReq struct {
	Vid        string  `json:"vid" form:"vid"`
	AddList    []uint  `json:"addList" form:"addList"`
	CancelList []uint  `json:"cancelList" form:"cancelList"`
}

type CollectArticleReq struct {
	Aid uint `json:"aid" form:"aid"`
}
