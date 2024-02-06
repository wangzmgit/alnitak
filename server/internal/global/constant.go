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

	// 视频信息存在问题
	WRONG_VIDEO_INFO = 2100
	// 视频内容存在问题
	WRONG_VIDEO_CONTENT = 2200
	// 视频处理失败
	PROCESSING_FAIL = 2300
)
