package dto

type RoleListReq struct {
	Page     int
	PageSize int
}

// 新增角色
type AddRoleReq struct {
	Name string
	Code string
	Desc string
}

// 编辑角色
type EditRoleReq struct {
	ID   uint
	Name string
	Desc string
}

// 编辑角色首页
type EditRoleHomeReq struct {
	ID   uint
	Home string
}

// 设置用户角色
type EditUserRoleReq struct {
	Uid  uint
	Code string
}
