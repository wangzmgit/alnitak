package types

import "encoding/json"

// CopyrightType 版权类型，统一处理前端 bool 和后端 int8 的 JSON 差异。
// 数据库/模型层用 int8 存储，通过自定义 MarshalJSON/UnmarshalJSON 在
// JSON 序列化时自动兼容 true/false（前端提交）和 0/1/2/3（后端存储/缓存）。
type CopyrightType int8

const (
	CopyrightUnknown CopyrightType = 0 // 未知版权（默认/历史遗留）
	CopyrightOriginal CopyrightType = 1 // 原创
	CopyrightReprint CopyrightType = 2 // 转载/搬运
	CopyrightPGC CopyrightType = 3 // PGC 授权内容
)

// MarshalJSON 统一输出为 int（0/1/2/3）
func (c CopyrightType) MarshalJSON() ([]byte, error) {
	return json.Marshal(int8(c))
}

// UnmarshalJSON 兼容 bool（true/false）和 int（0/1/2/3）
func (c *CopyrightType) UnmarshalJSON(data []byte) error {
	// 先试 int
	var i int8
	if err := json.Unmarshal(data, &i); err == nil {
		*c = CopyrightType(i)
		return nil
	}
	// 再试 bool（前端上传 true=原创/false=转载）
	var b bool
	if err := json.Unmarshal(data, &b); err != nil {
		*c = CopyrightUnknown
		return nil
	}
	if b {
		*c = CopyrightOriginal
	} else {
		*c = CopyrightReprint
	}
	return nil
}
