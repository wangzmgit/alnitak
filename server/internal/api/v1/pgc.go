package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
)

func CreatePGC(ctx *gin.Context) {
	var req dto.CreatePGCReq
	if err := ctx.ShouldBind(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}

	pgcID, err := service.CreatePGCContent(req)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 必须用字符串返回：JS JSON.parse 对超过 Number.MAX_SAFE_INTEGER 的整数会丢精度，导致跳转编辑页 ID 错误
	resp.OkWithData(ctx, gin.H{"pgc_id": strconv.FormatUint(pgcID, 10)})
}

func UpdatePGC(ctx *gin.Context) {
	var req dto.UpdatePGCReq
	if err := ctx.ShouldBind(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}

	if err := service.UpdatePGCContent(req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithMessage(ctx, "更新成功")
}

func DeletePGC(ctx *gin.Context) {
	pgcID := ctx.Param("pgc_id")

	pgcIDUint, err := convertToUint(pgcID)
	if err != nil {
		resp.FailWithMessage(ctx, "无效的PGC ID")
		return
	}

	if err := service.DeletePGCContent(pgcIDUint); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithMessage(ctx, "删除成功")
}

func GetPGCList(ctx *gin.Context) {
	var req dto.PGCListReq
	if err := ctx.ShouldBind(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}

	total, list, err := service.GetPGCContentList(req)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	formatted := make([]gin.H, 0, len(list))
	for _, item := range list {
		formatted = append(formatted, formatPGCContent(item))
	}

	resp.OkWithData(ctx, gin.H{
		"total":     total,
		"list":      formatted,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetPGCManageList 后台：PGC 内容管理列表（所有状态）
func GetPGCManageList(ctx *gin.Context) {
	var req dto.PGCManageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}

	// 复用已有 service 函数
	listReq := dto.PGCListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		PGCType:  req.PGCType,
		Status:   req.Status,
		Keyword:  req.Keyword,
	}

	total, list, err := service.GetPGCContentList(listReq)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	formatted := make([]gin.H, 0, len(list))
	for _, item := range list {
		formatted = append(formatted, formatPGCContent(item))
	}

	resp.OkWithData(ctx, gin.H{
		"total":     total,
		"list":      formatted,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetPGCReviewList 后台：待审 PGC 列表
func GetPGCReviewList(ctx *gin.Context) {
	var req dto.PGCReviewListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}

	total, list, err := service.GetPGCReviewList(req.Page, req.PageSize)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	formatted := make([]gin.H, 0, len(list))
	for _, item := range list {
		formatted = append(formatted, formatPGCContent(item))
	}

	resp.OkWithData(ctx, gin.H{
		"total": total,
		"list":  formatted,
	})
}

// ReviewPGCApproved 后台：审核通过
func ReviewPGCApproved(ctx *gin.Context) {
	var req dto.PGCReviewActionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}
	if err := service.AdminReviewPGC(req.PGCID, global.PGCAuditApproved); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.OkWithMessage(ctx, "操作成功")
}

// ReviewPGCFailed 后台：审核驳回
func ReviewPGCFailed(ctx *gin.Context) {
	var req dto.PGCReviewActionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}
	if err := service.AdminReviewPGC(req.PGCID, global.PGCAuditRejected); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.OkWithMessage(ctx, "操作成功")
}

// AdminUpdatePGCStatus 后台：管理员修改 PGC 状态（上架/下架）
func AdminUpdatePGCStatus(ctx *gin.Context) {
	var req dto.UpdatePGCStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdatePGCStatus(req.PGCID, req.Status); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.OkWithMessage(ctx, "操作成功")
}

// AdminDeletePGC 后台：管理员删除 PGC
func AdminDeletePGC(ctx *gin.Context) {
	pgcID := ctx.Param("pgc_id")
	pgcIDUint, err := convertToUint(pgcID)
	if err != nil {
		resp.FailWithMessage(ctx, "无效的PGC ID")
		return
	}
	if err := service.DeletePGCContent(pgcIDUint); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.OkWithMessage(ctx, "删除成功")
}

func GetPGCDetail(ctx *gin.Context) {
	pgcID := ctx.Query("pgc_id")

	pgcIDUint, err := convertToUint(pgcID)
	if err != nil {
		resp.FailWithMessage(ctx, "无效的PGC ID")
		return
	}

	pgc, err := service.GetPGCContentDetail(pgcIDUint)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, gin.H{"pgc": formatPGCContent(*pgc)})
}

func GetPGCEpisodes(ctx *gin.Context) {
	pgcIDStr := ctx.Param("pgc_id")
	pgcIDUint, err := convertToUint(pgcIDStr)
	if err != nil {
		resp.FailWithMessage(ctx, "无效的PGC ID")
		return
	}
	var req dto.EpisodeListReq
	if err := ctx.ShouldBind(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}

	episodes, err := service.GetPGCEpisodeList(pgcIDUint, req)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	formatted := make([]gin.H, 0, len(episodes))
	for _, item := range episodes {
		formatted = append(formatted, formatPGCEpisode(item))
	}

	resp.OkWithData(ctx, gin.H{"episodes": formatted})
}

func AddPGCEpisode(ctx *gin.Context) {
	pgcID := ctx.Param("pgc_id")

	pgcIDUint, err := convertToUint(pgcID)
	if err != nil {
		resp.FailWithMessage(ctx, "无效的PGC ID")
		return
	}

	var req dto.EpisodeReq
	if err := ctx.ShouldBind(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}

	if err := service.AddPGCEpisode(pgcIDUint, req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithMessage(ctx, "添加成功")
}

// BindPGCEpisodeVideo 占位剧集绑定已有视频
func BindPGCEpisodeVideo(ctx *gin.Context) {
	pgcID := ctx.Param("pgc_id")
	episodeID := ctx.Param("id")
	pgcIDUint, err1 := convertToUint(pgcID)
	episodeIDUint, err2 := convertToUint(episodeID)
	if err1 != nil || err2 != nil {
		resp.FailWithMessage(ctx, "无效的ID")
		return
	}
	var req dto.BindPGCEpisodeVideoReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}
	if err := service.BindPGCEpisodeVideo(pgcIDUint, episodeIDUint, req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.OkWithMessage(ctx, "绑定成功")
}

func DeletePGCEpisode(ctx *gin.Context) {
	pgcID := ctx.Param("pgc_id")
	episodeID := ctx.Param("id")

	pgcIDUint, err1 := convertToUint(pgcID)
	episodeIDUint, err2 := convertToUint(episodeID)
	if err1 != nil || err2 != nil {
		resp.FailWithMessage(ctx, "无效的ID")
		return
	}

	if err := service.DeletePGCEpisode(pgcIDUint, episodeIDUint); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithMessage(ctx, "删除成功")
}

func UpdatePGCEpisode(ctx *gin.Context) {
	pgcID := ctx.Param("pgc_id")
	episodeID := ctx.Param("id")
	pgcIDUint, err1 := convertToUint(pgcID)
	episodeIDUint, err2 := convertToUint(episodeID)
	if err1 != nil || err2 != nil {
		resp.FailWithMessage(ctx, "无效的ID")
		return
	}
	var req dto.UpdatePGCEpisodeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdatePGCEpisode(pgcIDUint, episodeIDUint, req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.OkWithMessage(ctx, "更新成功")
}

func UpdatePGCStatus(ctx *gin.Context) {
	pgcID := ctx.Param("pgc_id")
	pgcIDUint, err := convertToUint(pgcID)
	if err != nil {
		resp.FailWithMessage(ctx, "无效的PGC ID")
		return
	}
	var req dto.UpdatePGCStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}
	// 以 path 为准
	req.PGCID = pgcIDUint
	if err := service.UpdatePGCStatus(req.PGCID, req.Status); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.OkWithMessage(ctx, "更新成功")
}

func UpdatePGCEpisodeStatus(ctx *gin.Context) {
	pgcID := ctx.Param("pgc_id")
	episodeID := ctx.Param("id")
	pgcIDUint, err1 := convertToUint(pgcID)
	episodeIDUint, err2 := convertToUint(episodeID)
	if err1 != nil || err2 != nil {
		resp.FailWithMessage(ctx, "无效的ID")
		return
	}
	var req dto.UpdatePGCEpisodeStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdatePGCEpisodeStatus(pgcIDUint, episodeIDUint, req.Status); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.OkWithMessage(ctx, "更新成功")
}

func SearchPGC(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	pgcType := ctx.Query("pgc_type")
	page := ctx.Query("page")
	pageSize := ctx.Query("page_size")

	// pgc_type 可选：为空时默认 0（全部类型）
	pgcTypeInt := 0
	if pgcType != "" {
		v, err := convertToInt(pgcType)
		if err != nil {
			resp.FailWithMessage(ctx, "无效的参数")
			return
		}
		pgcTypeInt = v
	}
	pageInt, err2 := convertToInt(page)
	pageSizeInt, err3 := convertToInt(pageSize)
	if err2 != nil || err3 != nil {
		resp.FailWithMessage(ctx, "无效的参数")
		return
	}

	total, list, err := service.SearchPGC(keyword, pgcTypeInt, pageInt, pageSizeInt)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	formatted := make([]gin.H, 0, len(list))
	for _, item := range list {
		formatted = append(formatted, formatPGCContent(item))
	}

	resp.OkWithData(ctx, gin.H{
		"total":     total,
		"list":      formatted,
		"page":      pageInt,
		"page_size": pageSizeInt,
	})
}

func GetPGCByType(ctx *gin.Context) {
	pgcType := ctx.Query("pgc_type")
	page := ctx.Query("page")
	pageSize := ctx.Query("page_size")

	pgcTypeInt, err1 := convertToInt(pgcType)
	pageInt, err2 := convertToInt(page)
	pageSizeInt, err3 := convertToInt(pageSize)
	if err1 != nil || err2 != nil || err3 != nil {
		resp.FailWithMessage(ctx, "无效的参数")
		return
	}

	total, list, err := service.GetPGCByType(pgcTypeInt, pageInt, pageSizeInt)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, gin.H{
		"total":     total,
		"list":      list,
		"page":      pageInt,
		"page_size": pageSizeInt,
	})
}

func GetOngoingPGC(ctx *gin.Context) {
	page := ctx.Query("page")
	pageSize := ctx.Query("page_size")

	pageInt, err1 := convertToInt(page)
	pageSizeInt, err2 := convertToInt(pageSize)
	if err1 != nil || err2 != nil {
		resp.FailWithMessage(ctx, "无效的分页参数")
		return
	}

	total, list, err := service.GetOngoingPGC(pageInt, pageSizeInt)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, gin.H{
		"total":     total,
		"list":      list,
		"page":      pageInt,
		"page_size": pageSizeInt,
	})
}

func GetRecommendedPGC(ctx *gin.Context) {
	limit := ctx.Query("limit")
	limitInt, err := convertToInt(limit)
	if err != nil {
		resp.FailWithMessage(ctx, "无效的limit参数")
		return
	}

	list, err := service.GetRecommendedPGC(limitInt)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, gin.H{"list": list})
}

// RecommendPGC PGC 推荐（参考 B 站：按 seed/type 召回，并过滤不可播放内容）
//
// Query:
// - page, page_size: 分页（page_size <= 50）
// - pgc_type: 可选，强制指定类型
// - seed_pgc_id: 可选，种子 season_id（未传 pgc_type 时会用 seed 推断类型）
// - scene: 预留（home/detail）
func RecommendPGC(ctx *gin.Context) {
	var req dto.PGCRecommendReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		resp.FailWithMessage(ctx, "参数错误: "+err.Error())
		return
	}

	total, list, latest, err := service.RecommendPGC(req)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	formatted := make([]gin.H, 0, len(list))
	for _, item := range list {
		card := formatPGCContent(item)
		if ep, ok := latest[item.PGCID]; ok {
			card["latest_ep_id"] = ep.ID
			card["latest_ep_number"] = ep.EpisodeNumber
			card["latest_ep_title"] = ep.Title
			card["latest_vid"] = ep.VID
			card["latest_publish_time"] = ep.PublishTime
			// 对齐 B 站 new_ep 语义（简化版）
			if ep.EpisodeNumber > 0 {
				card["new_ep"] = gin.H{
					"index_show": fmt.Sprintf("第%d话", ep.EpisodeNumber),
					"title":      ep.Title,
				}
			}
		}
		// 简单角标：连载中 / 评分
		if item.IsOngoing {
			card["badge"] = "连载中"
		} else if item.Rating > 0 {
			card["badge"] = fmt.Sprintf("%.1f分", item.Rating)
		}
		formatted = append(formatted, card)
	}

	resp.OkWithData(ctx, gin.H{
		"total":     total,
		"list":      formatted,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// RecommendPGCByVideo 播放页推荐：给定 vid，返回同类 PGC 推荐列表。
func RecommendPGCByVideo(ctx *gin.Context) {
	vidStr := ctx.Query("vid")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	vidInt, err1 := convertToInt(vidStr)
	pageInt, err2 := convertToInt(pageStr)
	pageSizeInt, err3 := convertToInt(pageSizeStr)
	if err1 != nil || err2 != nil || err3 != nil || vidInt <= 0 {
		resp.FailWithMessage(ctx, "无效的参数")
		return
	}

	total, list, latest, err := service.RecommendPGCByVideo(uint(vidInt), pageInt, pageSizeInt)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	formatted := make([]gin.H, 0, len(list))
	for _, item := range list {
		card := formatPGCContent(item)
		if ep, ok := latest[item.PGCID]; ok {
			card["latest_ep_id"] = ep.ID
			card["latest_ep_number"] = ep.EpisodeNumber
			card["latest_ep_title"] = ep.Title
			card["latest_vid"] = ep.VID
			card["latest_publish_time"] = ep.PublishTime
			if ep.EpisodeNumber > 0 {
				card["new_ep"] = gin.H{
					"index_show": fmt.Sprintf("第%d话", ep.EpisodeNumber),
					"title":      ep.Title,
				}
			}
		}
		if item.IsOngoing {
			card["badge"] = "连载中"
		} else if item.Rating > 0 {
			card["badge"] = fmt.Sprintf("%.1f分", item.Rating)
		}
		formatted = append(formatted, card)
	}

	resp.OkWithData(ctx, gin.H{
		"total":     total,
		"list":      formatted,
		"page":      pageInt,
		"page_size": pageSizeInt,
	})
}

func GetPGCDetailWithEpisodes(ctx *gin.Context) {
	pgcID := ctx.Query("pgc_id")

	pgcIDUint, err := convertToUint(pgcID)
	if err != nil {
		resp.FailWithMessage(ctx, "无效的PGC ID")
		return
	}

	pgc, episodes, err := service.GetPGCWithEpisodes(pgcIDUint)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	formattedEpisodes := make([]gin.H, 0, len(episodes))
	for _, item := range episodes {
		formattedEpisodes = append(formattedEpisodes, formatPGCEpisode(item))
	}

	resp.OkWithData(ctx, gin.H{
		"pgc":      formatPGCContent(*pgc),
		"episodes": formattedEpisodes,
	})
}

// GetPGCPlayPanelByVideo 播放页：按 vid 获取“季度 + 剧集”面板数据
func GetPGCPlayPanelByVideo(ctx *gin.Context) {
	vidStr := ctx.Query("vid")
	seasonStr := ctx.DefaultQuery("season_id", "0")

	vidInt, err1 := convertToInt(vidStr)
	seasonUint, err2 := convertToUint(seasonStr)
	if err1 != nil || err2 != nil || vidInt <= 0 {
		resp.FailWithMessage(ctx, "无效的参数")
		return
	}

	current, seasons, episodes, activeSeasonID, err := service.GetPGCPlayPanelByVideo(uint(vidInt), seasonUint)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	formattedSeasons := make([]gin.H, 0, len(seasons))
	for _, item := range seasons {
		formattedSeasons = append(formattedSeasons, formatPGCContent(item))
	}
	formattedEpisodes := make([]gin.H, 0, len(episodes))
	for _, item := range episodes {
		formattedEpisodes = append(formattedEpisodes, formatPGCEpisode(item))
	}

	resp.OkWithData(ctx, gin.H{
		"current":          formatPGCContent(*current),
		"seasons":          formattedSeasons,
		"episodes":         formattedEpisodes,
		"active_season_id": strconv.FormatUint(activeSeasonID, 10),
	})
}

func GetPGCEpisodeDetail(ctx *gin.Context) {
	epID := ctx.Query("ep_id")
	epIDUint, err := convertToUint(epID)
	if err != nil {
		resp.FailWithMessage(ctx, "无效的剧集ID")
		return
	}
	ep, err := service.GetPGCEpisodeDetail(epIDUint)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}
	resp.OkWithData(ctx, gin.H{"episode": formatPGCEpisode(*ep)})
}

func convertToUint(s string) (uint64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty")
	}
	return strconv.ParseUint(s, 10, 64)
}

func convertToInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// 向 B 站 season/media 语义对齐，同时保留当前项目字段兼容
func formatPGCContent(item model.PGCContent) gin.H {
	id := strconv.FormatUint(item.PGCID, 10)
	mediaID := item.MediaID
	if mediaID == 0 {
		// 兼容历史数据：早期未落 media_id 时回退 season_id
		mediaID = item.PGCID
	}
	mediaIDStr := strconv.FormatUint(mediaID, 10)
	return gin.H{
		"id":               item.ID,
		"pgc_id":           id,
		"season_id":        id, // season_id 对应当前 pgc_id
		"media_id":         mediaIDStr,
		"pgc_type":         item.PGCType,
		"title":            item.Title,
		"cover":            item.Cover,
		"desc":             item.Desc,
		"year":             item.Year,
		"area":             item.Area,
		"rating":           item.Rating,
		"is_ongoing":       item.IsOngoing,
		"total_episodes":   item.TotalEpisodes,
		"current_episodes": item.CurrentEpisodes,
		"status":           item.Status,
		"operator_id":      item.OperatorID,
		"created_at":       item.CreatedAt,
		"updated_at":       item.UpdatedAt,
	}
}

// 向 B 站 ep 语义对齐，同时保留当前项目字段兼容
func formatPGCEpisode(item model.PGCEpisode) gin.H {
	seasonID := strconv.FormatUint(item.PGCID, 10)
	return gin.H{
		"id":             item.ID,
		"ep_id":          item.ID,
		"pgc_id":         seasonID,
		"season_id":      seasonID,
		"episode_number": item.EpisodeNumber,
		"title":          item.Title,
		"vid":            item.VID,
		"duration":       item.Duration,
		"status":         item.Status,
		"publish_time":   item.PublishTime,
		"created_at":     item.CreatedAt,
		"updated_at":     item.UpdatedAt,
	}
}
