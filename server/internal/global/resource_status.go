package global

// 资源（分 P）与视频稿件共用一套 status 数值（resource.status / video.status 语义一致处）。
// 完整枚举说明（上线文档与前端展示请与此保持一致）：
//
//   0    AUDIT_APPROVED   审核通过 / 已上线（可播放，具体是否可播仍受转码结果影响）
//   100  CREATED_VIDEO    已创建稿件，待完善信息或待上传
//   200  VIDEO_PROCESSING 转码处理中
//   300  SUBMIT_REVIEW    已提交审核（若业务未单独使用，可与 WAITING_REVIEW 择一）
//   500  WAITING_REVIEW   待人工审核
//   2000 REVIEW_FAILED    审核驳回
//   3000 PROCESSING_FAIL  转码/处理失败
//
// 资源特有：冲突、去重等业务可在扩展字段（如 conflict_resource_id）中表达，不单独占用 status。

// 分P 对外可见性（resource.visible_status）
const (
	VISIBLE_HIDDEN = 0 // 隐藏（处理中 / 排队中，不返回给前端）
	VISIBLE_SHOWN  = 1 // 可见（转码完成且已上线）
)

// ResourceStatusName 返回资源状态简短英文键，便于前端 i18n。
func ResourceStatusName(status int) string {
	switch status {
	case AUDIT_APPROVED:
		return "approved"
	case CREATED_VIDEO:
		return "created"
	case VIDEO_PROCESSING:
		return "processing"
	case SUBMIT_REVIEW:
		return "submitted"
	case WAITING_REVIEW:
		return "pending_review"
	case REVIEW_FAILED:
		return "rejected"
	case PROCESSING_FAIL:
		return "process_failed"
	case UPLOAD_FAILED:
		return "upload_failed"
	default:
		return "unknown"
	}
}
