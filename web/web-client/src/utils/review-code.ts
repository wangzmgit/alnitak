export const reviewCode = {
    AUDIT_APPROVED: 0, //审核通过

    CREATED_VIDEO: 100, //成功创建视频
    VIDEO_PROCESSING: 200, //视频转码中
    SUBMIT_REVIEW: 300, // 提交审核
    WAITING_REVIEW: 500, //等待审核

    //审核不通过
    WRONG_VIDEO_INFO: 2100, //视频信息存在问题
    WRONG_VIDEO_CONTENT: 2200, //视频内容存在问题
    PROCESSING_FAIL: 2300, //视频处理失败
}