package service

import "interastral-peace.com/alnitak/internal/domain/dto"

// countEpisodeReqWithBoundVideo 统计创建请求中已绑定视频（vid>0）的剧集数量。
func countEpisodeReqWithBoundVideo(episodes []dto.EpisodeReq) int {
	n := 0
	for _, e := range episodes {
		if e.VID > 0 {
			n++
		}
	}
	return n
}

// uniquePositiveVidsFromEpisodeReq 从剧集请求中提取去重后的正数 vid（用于存在性校验与 markVideosAsPGCAttached）。
func uniquePositiveVidsFromEpisodeReq(episodes []dto.EpisodeReq) []uint {
	if len(episodes) == 0 {
		return nil
	}
	set := make(map[uint]struct{}, len(episodes))
	out := make([]uint, 0, len(episodes))
	for _, e := range episodes {
		if e.VID == 0 {
			continue
		}
		if _, ok := set[e.VID]; ok {
			continue
		}
		set[e.VID] = struct{}{}
		out = append(out, e.VID)
	}
	return out
}
