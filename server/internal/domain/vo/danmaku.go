package vo

import "time"

const (
    DANMAKU_FIELD = "`id`, `uid`, `time`, `type`, `color`, `text`, `created_at`"
)

type DanmakuResp struct {
    ID               uint      `json:"id"`
    Uid              uint      `json:"uid"`
    Time             float32   `json:"time"`
    Type             int       `json:"type"`
    Color            string    `json:"color"`
    Text             string    `json:"text"`
    Part             uint      `json:"part"`
    ResourceShortID  string    `json:"resourceShortId"`
    CreatedAt        time.Time `json:"createdAt"` 
}
