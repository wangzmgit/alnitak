package dto

// 修改资源标题
type ModifyResourceTitleReq struct {
	ID    uint
	Title string
}

// 替换资源请求
type ReplaceResourceReq struct {
	ResourceID uint   `json:"resourceId"` // 要替换的资源ID
	Hash       string `json:"hash"`       // 新视频文件的hash
}

// 资源排序请求
type ReorderResourceReq struct {
	Vid         uint   `json:"vid"`         // 视频ID
	ResourceIds []uint `json:"resourceIds"` // 按顺序排列的资源ID列表
}

// 删除资源请求
type DeleteResourceReq struct {
	ResourceID     uint `json:"resourceId"`     // 资源ID
	DeleteDanmaku  bool `json:"deleteDanmaku"`   // 是否同时删除弹幕
}
