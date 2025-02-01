package dto

type AddCollectionReq struct {
	Name string //收藏夹名称
}

type EditCollectionReq struct {
	ID    uint
	Cover string //封面图
	Name  string //收藏夹名称
	Desc  string //简介
	Open  bool   //是否公开
}
