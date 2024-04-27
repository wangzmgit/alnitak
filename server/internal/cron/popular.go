package cron

import (
	"container/heap"
	"math"
	"time"

	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/service"
)

type VideoData struct {
	ID        uint
	Clicks    int64
	CreatedAt time.Time
	Score     float64 `gorm:"-"`
}

type IntHeap []VideoData

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i].Score > h[j].Score }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(VideoData))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func RefreshPopular() {
	offset := 0
	limit := 1000

	h := &IntHeap{}
	heap.Init(h)
	var videos []VideoData
	nowTime := time.Now() // 获取当前时间
	zap.L().Info("开始同步视频ID", zap.String("module", "initialize"))

	for {
		if err := global.Mysql.Model(&model.Video{}).Where("status = ?", global.AUDIT_APPROVED).
			Offset(offset).Limit(limit).Scan(&videos).Error; err != nil {
			break
		}

		for _, video := range videos {
			// 计算视频分值
			archive := service.FindArchiveData(video.ID)
			likeValue := 12.5 * float64(archive.Like)
			collectValue := 25 * float64(archive.Like)

			t := float64(nowTime.Sub(video.CreatedAt).Hours() / 24)
			video.Score = hot(collectValue+likeValue+float64(video.Clicks), t)

			heap.Push(h, video)
			if h.Len() > 100 {
				heap.Pop(h)
			}
		}

		if len(videos) < limit {
			break
		}

		offset += limit
	}

	// 写入缓存
	cache.DelHotVideoId()
	for i := 0; i < h.Len(); i++ {
		cache.SetHotVideoId(heap.Pop(h).(VideoData).ID)
	}

	zap.L().Info("视频ID同步完成，耗时:"+time.Since(nowTime).String(), zap.String("initialize", "cron"))
}

func hot(weight, t float64) float64 {
	return weight * math.Exp(-0.011*t)
}
