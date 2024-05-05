package dto

type ReviewVideoReq struct {
	Vid    uint
	Status int
	Remark string
}

type ReviewArticleReq struct {
	Aid    uint
	Status int
	Remark string
}
