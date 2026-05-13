package api

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

// GetSubtitleList GET /api/v1/video/subtitle/list?resourceShortId=
func GetSubtitleList(ctx *gin.Context) {
	rsid := strings.TrimSpace(ctx.Query("resourceShortId"))
	if rsid == "" {
		resp.FailWithMessage(ctx, "resourceShortId 不能为空")
		return
	}
	list, err := service.ListSubtitleTracks(ctx, rsid)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.OkWithData(ctx, gin.H{"tracks": list})
}

// PostSubtitleUpload POST /api/v1/video/subtitle/upload （multipart）
func PostSubtitleUpload(ctx *gin.Context) {
	rsid := strings.TrimSpace(ctx.PostForm("resourceShortId"))
	lang := strings.TrimSpace(ctx.PostForm("lang"))
	label := strings.TrimSpace(ctx.PostForm("label"))
	isDefault := ctx.PostForm("isDefault") == "1" || strings.EqualFold(ctx.PostForm("isDefault"), "true")

	if rsid == "" {
		resp.FailWithMessage(ctx, "resourceShortId 不能为空")
		return
	}
	file, err := ctx.FormFile("file")
	if err != nil || file == nil {
		resp.FailWithMessage(ctx, "请上传字幕文件")
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".srt" && ext != ".vtt" {
		resp.FailWithMessage(ctx, "仅支持 .srt 或 .vtt")
		return
	}
	src, err := file.Open()
	if err != nil {
		resp.FailWithMessage(ctx, "读取文件失败")
		return
	}
	defer src.Close()

	err = service.CreateSubtitleTrack(ctx, rsid, lang, label, isDefault, ext, src)
	if err != nil {
		if errors.Is(err, service.ErrSubtitleDuplicateLang) {
			resp.FailWithMessage(ctx, err.Error())
			return
		}
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.Ok(ctx)
}

type putSubtitleBody struct {
	Label     string `json:"label"`
	IsDefault *bool  `json:"isDefault"`
}

// PutSubtitle PUT /api/v1/video/subtitle/:id
// 可选 multipart：file + 与上传相同扩展名；仅 JSON 则只更新元数据
func PutSubtitle(ctx *gin.Context) {
	id := utils.StringToUint(ctx.Param("id"))
	if id == 0 {
		resp.FailWithMessage(ctx, "字幕不存在")
		return
	}
	req := service.UpdateSubtitleTrackReq{}
	if ctx.ContentType() != "" && strings.Contains(ctx.ContentType(), "application/json") {
		var body putSubtitleBody
		if err := ctx.ShouldBindJSON(&body); err != nil {
			resp.FailWithMessage(ctx, "参数有误")
			return
		}
		req.Label = body.Label
		req.IsDefault = body.IsDefault
	} else {
		req.Label = strings.TrimSpace(ctx.PostForm("label"))
		if ctx.PostForm("isDefault") != "" {
			v := ctx.PostForm("isDefault") == "1" || strings.EqualFold(ctx.PostForm("isDefault"), "true")
			req.IsDefault = &v
		}
		file, ferr := ctx.FormFile("file")
		if ferr == nil && file != nil {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext != ".srt" && ext != ".vtt" {
				resp.FailWithMessage(ctx, "仅支持 .srt 或 .vtt")
				return
			}
			src, err := file.Open()
			if err != nil {
				resp.FailWithMessage(ctx, "读取文件失败")
				return
			}
			defer src.Close()
			req.ReplaceVTT = src
			req.SrcExt = ext
		}
	}

	if err := service.UpdateSubtitleTrack(ctx, id, req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.Ok(ctx)
}

// DeleteSubtitle DELETE /api/v1/video/subtitle/:id
func DeleteSubtitle(ctx *gin.Context) {
	id := utils.StringToUint(ctx.Param("id"))
	if id == 0 {
		resp.FailWithMessage(ctx, "字幕不存在")
		return
	}
	if err := service.DeleteSubtitleTrack(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.Ok(ctx)
}
