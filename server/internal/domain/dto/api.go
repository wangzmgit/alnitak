package dto

type ApiListReq struct {
	Page     int
	PageSize int
}

type AddApiReq struct {
	Path     string
	Category string
	Method   string
	Desc     string
}

type EditApiReq struct {
	ID       uint
	Path     string
	Category string
	Method   string
	Desc     string
}

// 编辑角色Api
type EditRoleApiReq struct {
	Id        uint   //角色ID
	AddIds    []uint //添加API ID数组
	RemoveIds []uint //移除API ID数组
}
