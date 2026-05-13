package dto

type HistoryReq struct {
	Vid               interface{} `json:"vid" form:"vid"`
	Part              uint        `json:"part" form:"part"`
	Time              float64    `json:"time" form:"time"`
	Duration          int         `json:"duration" form:"duration"`
	Rid               string     `json:"rid" form:"rid"`
}
