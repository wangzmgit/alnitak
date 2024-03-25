package dto

type DanmakuReq struct {
	Vid   uint
	Part  uint
	Time  float32 //时间
	Type  int     //类型0滚动;1顶部;2底部
	Color string
	Text  string
}
