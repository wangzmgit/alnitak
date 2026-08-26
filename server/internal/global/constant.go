package global

import "interastral-peace.com/alnitak/internal/domain/types"

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
	// 转码成功，上传 OSS 失败（可重试上传，无需重新转码）
	UPLOAD_FAILED = 3001
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
	CONTENT_TYPE_VIDEO    = 0
	CONTENT_TYPE_ARTICLE  = 1
	CONTENT_TYPE_PLAYLIST = 2
	CONTENT_TYPE_COMMENT  = 3 // 评论
)

// ========== PGC内容类型 ==========
const (
	// PGCTypeNone 非PGC
	PGCTypeNone = 0
	// PGCTypeCN 国创（番剧 CN/国漫）
	PGCTypeCN = 1
	// PGCTypeJP 日创（番剧 JP/日漫）
	PGCTypeJP = 2
	// PGCTypeDocumentary 纪录片
	PGCTypeDocumentary = 3
	// PGCTypeMovie 电影
	PGCTypeMovie = 4
	// PGCTypeTVSeries 电视剧
	PGCTypeTVSeries = 5
)

// ========== PGC审核状态 ==========
const (
	// PGCAuditDraft 草稿
	PGCAuditDraft = 0
	// PGCAuditSubmitted 已提交
	PGCAuditSubmitted = 100
	// PGCAuditProcessing 审核中
	PGCAuditProcessing = 200
	// PGCAuditApproved 审核通过
	PGCAuditApproved = 300
	// PGCAuditRejected 审核驳回
	PGCAuditRejected = 400
	// PGCAuditOffline 下架（管理侧操作）
	PGCAuditOffline = -1
)

// ========== 版权类型 ==========
// 类型定义和常量值在 types.CopyrightType 中统一管理
// 此处为便利别名，让现有代码无需修改 import
const (
	CopyrightUnknown types.CopyrightType = 0 // 未知版权（默认/历史遗留）
	CopyrightOriginal types.CopyrightType = 1 // 原创
	CopyrightReprint types.CopyrightType = 2 // 转载/搬运
	CopyrightPGC types.CopyrightType = 3 // PGC 授权内容
)

// ========== PGC剧集状态 ==========
const (
	// PGCEpisodeNormal 正常
	PGCEpisodeNormal = 0
	// PGCEpisodeOffline 下架
	PGCEpisodeOffline = -1
)
