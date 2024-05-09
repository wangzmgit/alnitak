package dto

type AddCollectReq struct {
	Vid        uint   //视频ID
	AddList    []uint //添加的收藏夹数组
	CancelList []uint //取消的收藏夹数组
}

type CollectArticleReq struct {
	Aid uint
}
