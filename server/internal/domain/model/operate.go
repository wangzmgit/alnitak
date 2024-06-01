package model

import (
	"time"

	"gorm.io/gorm"
)

type Operate struct {
	gorm.Model
	Ip           string        `gorm:"column:ip;comment:请求ip"`                         // 请求ip
	Method       string        `gorm:"column:method;comment:请求方法"`                     // 请求方法
	Path         string        `gorm:"column:path;comment:请求路径"`                       // 请求路径
	Status       int           `gorm:"column:status;comment:请求状态"`                     // 请求状态
	Latency      time.Duration `gorm:"column:latency;comment:延迟" swaggertype:"string"` // 延迟
	Agent        string        `gorm:"column:agent;comment:代理"`                        // 代理
	ErrorMessage string        `gorm:"column:error_message;comment:错误信息"`              // 错误信息
	Body         string        `gorm:"type:text;column:body;comment:请求Body"`           // 请求Body
	Msg          string        `gorm:"type:text;column:msg;comment:响应Msg"`             // 响应Msg
	UserID       int           `gorm:"column:user_id;comment:用户id"`                    // 用户id
}

func (table *Operate) TableName() string {
	return "operate"
}
