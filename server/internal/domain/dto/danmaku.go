package dto

type DanmakuReq struct {
	Vid               string  `json:"vid" form:"vid"`
	Part              uint    `json:"part" form:"part"`
	Rid               string  `json:"rid" form:"rid"`
	Time              float32 `json:"time" form:"time"`
	Type              int     `json:"type" form:"type"`
	Color             string  `json:"color" form:"color"`
	Text              string  `json:"text" form:"text"`
}
