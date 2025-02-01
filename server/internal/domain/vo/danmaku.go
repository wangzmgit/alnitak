package vo

const (
	DANMAKU_FIELD = "`id`,`time`,`type`,`color`,`text`"
)

type DanmakuResp struct {
	ID    uint    `json:"id"`
	Time  float32 `json:"time"`
	Type  int     `json:"type"`
	Color string  `json:"color"`
	Text  string  `json:"text"`
}
