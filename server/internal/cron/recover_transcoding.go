package cron

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

const (
	transcodingQueueStream    = "transcoding:queue"
	transcodingStatusPrefix   = "transcoding:status:"
	retryCounterPrefix        = "transcoding:retry:"
	maxReenqueueRetries       = 3
)



// RecoverStuckTranscoding 定期恢复卡住的转码任务
//
// 转码任务入队（Enqueue）失败时错误被静默忽略（_ = Enqueue(...)），
// 导致资源状态停留在 VIDEO_PROCESSING 但实际没有 Worker 处理。
// 此函数检查所有处于 VIDEO_PROCESSING 的资源，如果没有对应的
// Redis 进度哈希（transcoding:status:{resourceID}），说明任务从未
// 成功入队，尝试重新入队。
func RecoverStuckTranscoding() {
	utils.InfoLog("开始检查卡住的转码任务", "cron")

	rdb := global.Redis.RawClient()
	ctx := context.Background()

	recovered := 0
	skipped := 0
	autoRetried := 0
	autoRetryFailed := 0
	errors := 0

	// ═══════════════════════════════════════════════════════════════
	// Part 1: 恢复 VIDEO_PROCESSING 状态的任务（原有逻辑）
	// 【注意】Local 模式使用进程内内存进度，不写 Redis status key，
	// 也不经过 Stream，因此 Redis 查询永远为空，会误判为卡住并重复入队。
	// Local 模式下跳过此部分。
	// ═══════════════════════════════════════════════════════════════
	if global.Config.Transcoding.Mode == "local" {
		utils.InfoLog("Local 模式跳过 Redis 恢复检查（使用内存进度）", "cron")
	} else {
	var processingResources []model.Resource
	global.Mysql.Model(&model.Resource{}).
		Where("status = ?", global.VIDEO_PROCESSING).
		Find(&processingResources)

	for _, res := range processingResources {
		redisKey := fmt.Sprintf("%s%d", transcodingStatusPrefix, res.ID)

		// 检查 Redis 中是否有该资源的进度哈希
		exists, err := rdb.Exists(ctx, redisKey).Result()
		if err != nil {
			utils.ErrorLog("检查Redis失败", "cron",
				fmt.Sprintf("ResourceID=%d, err=%v", res.ID, err))
			errors++
			continue
		}

		if exists > 0 {
			// 有进度哈希 → 检查是否过期（pollProgress 丢失后留下的残渣）
			status, err := rdb.HGet(ctx, redisKey, "status").Result()
			if err == redis.Nil {
				// 键存在但无 status 字段：Worker 进度 HSet 重建了 key 但未写入 status。
				// 检查是否有 progress_* 字段，有则任务仍在进行中。
				fieldCount, _ := rdb.HLen(ctx, redisKey).Result()
				if fieldCount > 0 {
					skipped++
					continue
				}
				// 空 key，删掉走下方 reenqueue 逻辑
				rdb.Del(ctx, redisKey)
			} else if err != nil {
				utils.ErrorLog("读取哈希状态失败", "cron",
					fmt.Sprintf("ResourceID=%d, err=%v", res.ID, err))
				errors++
				continue
			}
			switch status {
			case "completed", "failed", "partial_fail":
				// pollProgress 丢失（主服务重启）导致无人处理此完成状态，
				// 删掉哈希后重新入队，新的 pollProgress 会接管。
				utils.InfoLog(fmt.Sprintf("发现过期终态，重新入队: ResourceID=%d, VideoID=%d, status=%s",
					res.ID, res.Vid, status), "cron")
			case "processing":
				// "processing" → 检查 updated 字段判断是否过期
				updatedStr, err := rdb.HGet(ctx, redisKey, "updated").Result()
				if err != nil {
					utils.ErrorLog("读取 updated 字段失败", "cron",
						fmt.Sprintf("ResourceID=%d, err=%v", res.ID, err))
					errors++
					continue
				}
				updatedUnix, err := strconv.ParseInt(updatedStr, 10, 64)
				if err != nil {
					utils.ErrorLog("解析 updated 字段失败", "cron",
						fmt.Sprintf("ResourceID=%d, val=%s", res.ID, updatedStr))
					errors++
					continue
				}
				updatedAt := time.Unix(updatedUnix, 0)
				if time.Since(updatedAt) < 5*time.Minute {
					// 5分钟内更新过 → 任务正常进行中
					skipped++
					continue
				}
				// 超过5分钟未更新 → pollProgress 丢失（服务重启），重新入队
				utils.InfoLog(fmt.Sprintf("发现过期 processing 状态（updated=%v），重新入队: ResourceID=%d, VideoID=%d",
					time.Since(updatedAt).Round(time.Second), res.ID, res.Vid), "cron")
			default:
				// 未知 status，跳过
				skipped++
				continue
			}

			// 到达这里说明需要重新入队 → 删哈希，继续执行下面的 reenqueue 逻辑
			if err := rdb.Del(ctx, redisKey).Err(); err != nil {
				utils.ErrorLog("删除过期哈希失败", "cron",
					fmt.Sprintf("ResourceID=%d, err=%v", res.ID, err))
			}
		}

		// 再检查 Stream 中是否已有该 resourceID 的条目
		// （Enqueue 成功后 Worker 还没开始处理时，Hash 不存在但 Stream 里有消息）
		inStream, err := resourceInStream(ctx, rdb, res.ID)
		if err != nil {
			utils.ErrorLog("检查Stream失败", "cron",
				fmt.Sprintf("ResourceID=%d, err=%v", res.ID, err))
			errors++
			continue
		}
		if inStream {
			skipped++
			continue
		}

		// 没有 Redis 进度记录（或已被删除）+ Stream 中无条目
		// → Enqueue 从未成功或哈希残留已清理 → 重新入队
		utils.InfoLog(fmt.Sprintf("发现卡住的转码任务，准备重新入队: ResourceID=%d, VideoID=%d", res.ID, res.Vid), "cron")

		if err := reenqueueResource(res); err != nil {
			utils.ErrorLog("重新入队失败", "cron",
				fmt.Sprintf("ResourceID=%d, err=%v", res.ID, err))
			errors++
			continue
		}

		recovered++
	}
	} // end if mode != "local"

	// ═══════════════════════════════════════════════════════════════
	// Part 2: 自动重试 PROCESSING_FAIL 状态的任务（失败自动重试）
	// ═══════════════════════════════════════════════════════════════
	var failedResources []model.Resource
	global.Mysql.Model(&model.Resource{}).
		Where("status = ?", global.PROCESSING_FAIL).
		Find(&failedResources)

	for _, res := range failedResources {
		// 检查 Stream 中是否已有该 resourceID 的条目（避免重复入队）
		inStream, err := resourceInStream(ctx, rdb, res.ID)
		if err != nil {
			utils.ErrorLog("检查Stream失败(自动重试)", "cron",
				fmt.Sprintf("ResourceID=%d, err=%v", res.ID, err))
			errors++
			continue
		}
		if inStream {
			skipped++
			continue
		}

		// 检查重试上限
		retryKey := fmt.Sprintf("%s%d", retryCounterPrefix, res.ID)
		retryCount, err := rdb.Incr(ctx, retryKey).Result()
		if err != nil {
			utils.ErrorLog("重试计数器失败(自动重试)", "cron",
				fmt.Sprintf("ResourceID=%d, err=%v", res.ID, err))
			errors++
			continue
		}
		rdb.Expire(ctx, retryKey, 72*time.Hour)

		if retryCount > maxReenqueueRetries {
			utils.InfoLog(fmt.Sprintf("自动重试已达上限(%d)，放弃: ResourceID=%d, VideoID=%d",
				maxReenqueueRetries, res.ID, res.Vid), "cron")
			rdb.Del(ctx, retryKey)
			autoRetryFailed++
			continue
		}

		// 重置资源状态为转码中
		global.Mysql.Model(&model.Resource{}).Where("id = ?", res.ID).
			Update("status", global.VIDEO_PROCESSING)

		// 重置 Redis 中的画质状态
		statusKey := fmt.Sprintf("%s%d", transcodingStatusPrefix, res.ID)
		if hash, err := rdb.HGetAll(ctx, statusKey).Result(); err == nil && len(hash) > 0 {
			pipe := rdb.Pipeline()
			for field, val := range hash {
				if strings.HasPrefix(field, "status_") && val != "success" {
					pipe.HSet(ctx, statusKey, field, "waiting")
					progressField := "progress_" + field[len("status_"):]
					pipe.HSet(ctx, statusKey, progressField, "0.00")
				}
			}
			pipe.HSet(ctx, statusKey, "status", "processing")
			pipe.HSet(ctx, statusKey, "updated", fmt.Sprintf("%d", time.Now().Unix()))
			pipe.Exec(ctx)
		}

		utils.InfoLog(fmt.Sprintf("自动重试失败任务(第%d次): ResourceID=%d, VideoID=%d",
			retryCount, res.ID, res.Vid), "cron")

		if err := reenqueueResource(res); err != nil {
			utils.ErrorLog("自动重试入队失败", "cron",
				fmt.Sprintf("ResourceID=%d, err=%v", res.ID, err))
			// 恢复状态为失败
			global.Mysql.Model(&model.Resource{}).Where("id = ?", res.ID).
				Update("status", global.PROCESSING_FAIL)
			errors++
			continue
		}

		autoRetried++
	}

	utils.InfoLog(fmt.Sprintf("转码恢复任务完成: 恢复(卡住)=%d, 自动重试(失败)=%d, 自动重试放弃=%d, 跳过=%d, 错误=%d",
		recovered, autoRetried, autoRetryFailed, skipped, errors), "cron")
}

// getOriginalVideoStatus 读取 ReTranscodeVideo 写入 Redis 的原始审核状态。
// 如果键不存在或已过期，返回 -1（普通恢复，不改变审核状态）。
func getOriginalVideoStatus(videoID uint) int {
	rdb := global.Redis.RawClient()
	ctx := context.Background()
	key := fmt.Sprintf("transcoding:orig_status:%d", videoID)
	val, err := rdb.Get(ctx, key).Int()
	if err != nil {
		return -1
	}
	utils.InfoLog(fmt.Sprintf("读取到原始审核状态: VideoID=%d, status=%d", videoID, val), "cron")
	return val
}

// resourceInStream 检查指定 resourceID 是否已在 Redis Stream 中等待处理。
// 扫描 Stream 所有条目（maxDepth <= 10，最多 10 条），避免在 Worker 
// 还未开始处理（未写进度 Hash）时重复入队。
func resourceInStream(ctx context.Context, rdb *redis.Client, resourceID uint) (bool, error) {
	entries, err := rdb.XRange(ctx, transcodingQueueStream, "-", "+").Result()
	if err != nil {
		return false, err
	}
	target := strconv.FormatUint(uint64(resourceID), 10)
	for _, entry := range entries {
		if rid, ok := entry.Values["resourceID"].(string); ok && rid == target {
			return true, nil
		}
	}
	return false, nil
}

// reenqueueResource 为指定资源重建 TranscodingInfo 并重新入队
// 每次调用前会检查重试上限，超出则标记 PROCESSING_FAIL 不再重试。
func reenqueueResource(res model.Resource) error {
	rdb := global.Redis.RawClient()
	ctx := context.Background()

	// 0. 重试上限检查
	retryKey := fmt.Sprintf("%s%d", retryCounterPrefix, res.ID)
	retryCount, err := rdb.Incr(ctx, retryKey).Result()
	if err != nil {
		return fmt.Errorf("重试计数器失败: %w", err)
	}
	rdb.Expire(ctx, retryKey, 72*time.Hour) // 避免键堆积

	if retryCount > maxReenqueueRetries {
		utils.InfoLog(fmt.Sprintf("转码恢复已达上限(%d)，标记失败: ResourceID=%d",
			maxReenqueueRetries, res.ID), "cron")
		global.Mysql.Model(&model.Resource{}).Where("id = ?", res.ID).
			Update("status", global.PROCESSING_FAIL)
		rdb.Del(ctx, retryKey)
		// 同时清理可能残留的进度哈希
		statusKey := fmt.Sprintf("%s%d", transcodingStatusPrefix, res.ID)
		rdb.Del(ctx, statusKey)
		return fmt.Errorf("转码恢复已达上限(%d)，资源已标记为处理失败", maxReenqueueRetries)
	}

	// 1. 获取 VideoFile 信息（用于构建文件路径）
	var vf model.VideoFile
	if err := global.Mysql.Where("id = ?", res.FileID).First(&vf).Error; err != nil {
		return fmt.Errorf("查找VideoFile失败: %w", err)
	}
	if vf.DirName == "" {
		return fmt.Errorf("VideoFile.DirName 为空, FileID=%d", res.FileID)
	}

	suffix := utils.GetFileSuffix(vf.OriginalName)

	// 2. 获取原始审核状态（由 ReTranscodeVideo 写入 Redis）
	origStatus := getOriginalVideoStatus(res.Vid)

	// 3. 构建 TranscodingInfo
	//    远程模式：Worker 会从 OSS 下载文件并自行 ffprobe，只需基础字段
	//    本地模式：需要探测源文件
	inputPath := "./upload/video/" + vf.DirName + "/upload" + suffix

	var info *dto.TranscodingInfo

	if global.Config.Transcoding.Mode == "remote" {
		info = &dto.TranscodingInfo{
			VideoID:             res.Vid,
			ResourceID:          res.ID,
			InputFile:           inputPath,
			OutputDir:           "./upload/video/" + vf.DirName + "/",
			DirName:             vf.DirName,
			Suffix:              suffix,
			CodecName:           res.CodecName,
			Duration:            float64(res.Duration),
			OriginalVideoStatus: origStatus,
		}
	} else {
		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return fmt.Errorf("源文件不存在: %s", inputPath)
		}
		probed, err := service.ProcessVideoInfo(inputPath)
		if err != nil {
			return fmt.Errorf("探测视频信息失败: %w", err)
		}
		info = &dto.TranscodingInfo{
			VideoID:             res.Vid,
			ResourceID:          res.ID,
			InputFile:           inputPath,
			OutputDir:           "./upload/video/" + vf.DirName + "/",
			DirName:             vf.DirName,
			Suffix:              suffix,
			Width:               probed.Width,
			Height:              probed.Height,
			Duration:            probed.Duration,
			CodecName:           probed.CodecName,
			FPS:                 probed.FPS,
			FPS30:               probed.FPS30,
			FPS60:               probed.FPS60,
			VideoBitRate:        probed.VideoBitRate,
			AudioBitRate:        probed.AudioBitRate,
			AudioSampleRate:     probed.AudioSampleRate,
			AudioChannels:       probed.AudioChannels,
			AudioStreams:        probed.AudioStreams,
			OriginalVideoStatus: origStatus,
		}
	}

	// 4. 入队
	if err := service.GetCurrentTranscoder().Enqueue(context.Background(), info); err != nil {
		return fmt.Errorf("Enqueue失败: %w", err)
	}

	// 入队成功，清除重试计数器，下次再卡住从 1 开始计数
	rdb.Del(ctx, retryKey)

	utils.InfoLog(fmt.Sprintf("转码任务重新入队成功: ResourceID=%d, VideoID=%d, DirName=%s",
		res.ID, res.Vid, vf.DirName), "cron")

	return nil
}
