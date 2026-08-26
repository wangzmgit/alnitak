package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
)

type playGrantReq struct {
	ResourceShortID string `json:"resourceShortId" form:"resourceShortId"`
}

// PostPlayGrant 获取播放授权 token（可登录也可匿名；前端在过期前换发以实现无感续期）。
// 可选：在路由上使用 middleware.AuthOptional 若你后续需要区分会员权益。
func PostPlayGrant(ctx *gin.Context) {
	var req playGrantReq
	_ = ctx.ShouldBind(&req)
	if req.ResourceShortID == "" {
		req.ResourceShortID = ctx.Query("resourceShortId")
	}
	if req.ResourceShortID == "" {
		resp.FailWithMessage(ctx, "resourceShortId 不能为空")
		return
	}
	tok, exp, err := service.IssuePlayGrantForResource(ctx, req.ResourceShortID)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.OkWithData(ctx, gin.H{"token": tok, "expires": exp})
}

// GetPlayByResource 核心播放接口：grant token + 资源 shortId → 带签名的音视频 URL。
// 配置了备用 OSS 时同时返回 backupVideo/backupAudio（B站风格多源容灾）。
// GET /api/v1/play/:resourceShortId?token=JWT&quality=720p&audioLang=jpn（quality/audioLang 可选）
func GetPlayByResource(ctx *gin.Context) {
	rsid := ctx.Param("resourceShortId")
	tok := ctx.Query("token")
	if rsid == "" || tok == "" {
		resp.FailWithMessage(ctx, "缺少资源或 token")
		return
	}
	quality := ctx.Query("quality")
	audioLang := ctx.Query("audioLang")
	result, err := service.GetPlayURLs(ctx, rsid, tok, quality, audioLang)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.OkWithData(ctx, result)
}
