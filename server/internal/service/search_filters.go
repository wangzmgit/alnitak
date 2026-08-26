package service

import (
	"strings"
	"time"
)

// parseTimeRangeStart 解析 timeRange，返回起始时间（nil 表示不限制）
func parseTimeRangeStart(timeRange string) *time.Time {
	now := time.Now()
	switch timeRange {
	case "24h":
		t := now.Add(-24 * time.Hour)
		return &t
	case "week":
		t := now.AddDate(0, 0, -7)
		return &t
	case "month":
		t := now.AddDate(0, -1, 0)
		return &t
	case "year":
		t := now.AddDate(-1, 0, 0)
		return &t
	default:
		return nil
	}
}

func normalizeSort(sort string) string {
	switch strings.TrimSpace(sort) {
	case "most_viewed":
		return "most_viewed"
	case "newest":
		return "newest"
	case "relevance":
		return "relevance"
	default:
		return "relevance"
	}
}

func normalizeTimeRange(tr string) string {
	switch strings.TrimSpace(tr) {
	case "24h", "week", "month", "year", "all":
		return strings.TrimSpace(tr)
	default:
		return "all"
	}
}

// escapeLikeKeyword 转义 LIKE 模式字符，避免用户输入 %/_ 造成误匹配或扫描
// 配合 SQL: `LIKE ?`（MySQL 默认转义符为 `\`）
func escapeLikeKeyword(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

