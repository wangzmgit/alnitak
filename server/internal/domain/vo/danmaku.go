package vo

import "time"

const (
    DANMAKU_FIELD = "`id`, `time`, `type`, `color`, `text`, `created_at`"
)

type DanmakuResp struct {
    ID        uint      `json:"id"`
    Time      float32   `json:"time"`
    Type      int       `json:"type"`
    Color     string    `json:"color"`
    Text      string    `json:"text"`
    CreatedAt time.Time `json:"createdAt"` 
}
