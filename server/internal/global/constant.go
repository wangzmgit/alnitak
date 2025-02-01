package global

const CAPTCHA_STATUS_ABSENT = 0   // 人机验证状态不存在
const CAPTCHA_STATUS_NOT_USED = 1 // 人机验证状态未使用
const CAPTCHA_STATUS_PASS = 2     // 人机验证状态已通过

// 视频审核状态
const (
	// 审核通过
	AUDIT_APPROVED = 0
	// 成功创建视频
	CREATED_VIDEO = 100
	// 视频转码中
	VIDEO_PROCESSING = 200
	// 提交审核
	SUBMIT_REVIEW = 300
	// 等待审核
	WAITING_REVIEW = 500
	// 审核不通过
	REVIEW_FAILED = 2000
	// 视频处理失败
	PROCESSING_FAIL = 3000
)

// 用户关系
const (
	// 未关注
	NOT_FOLLOWING = 0
	// 已关注
	FOLLOWED = 1
	// 互粉
	MUTUAL_FANS = 2
)

const (
	CONTENT_TYPE_VIDEO   = 0
	CONTENT_TYPE_ARTICLE = 1
)
