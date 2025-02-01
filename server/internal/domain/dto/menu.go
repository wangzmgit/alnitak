package dto

type AddMenuReq struct {
	Name      string
	Path      string
	Component string
	Desc      string
	Sort      uint
	ParentId  uint
	// Meta内容
	Title     string
	Icon      string
	Hidden    bool
	KeepAlive bool
}

type EditMenuReq struct {
	ID        uint
	Name      string
	Path      string
	Component string
	Desc      string
	Sort      uint
	ParentId  uint
	// Meta内容
	Title     string
	Icon      string
	Hidden    bool
	KeepAlive bool
}

// 修改角色菜单
type EditRoleMenuReq struct {
	Id      uint   //角色ID
	MenuIds []uint //菜单ID数组
}
