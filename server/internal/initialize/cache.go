package initialize

import (
	"time"

	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
)

type VideoData struct {
	ID          uint
	PartitionId uint
}

// 将所有视频ID根据分区写入redis，如果后期视频量很多再去使用其他方案
func InitCacheData() {
	offset := 0
	limit := 1000
	var videos []VideoData

	start := time.Now() // 获取当前时间
	zap.L().Info("开始同步视频ID", zap.String("module", "initialize"))

	cache.DelAllVideoId()

	for {
		if err := global.Mysql.Model(&model.Video{}).Where("status = ?", global.AUDIT_APPROVED).Offset(offset).
			Limit(limit).Scan(&videos).Error; err != nil {
			break
		}

		for _, video := range videos {
			cache.SetVideoId(global.VideoPartitionMap[video.PartitionId], video.ID)
		}

		if len(videos) < limit {
			break
		}

		offset += limit
	}

	zap.L().Info("视频ID同步完成，耗时:"+time.Since(start).String(), zap.String("initialize", "cron"))

}
