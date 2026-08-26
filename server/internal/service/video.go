package service

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

func UploadVideoInfo(ctx *gin.Context, uploadVideoReq dto.UploadVideoReq) error {
	userId := ctx.GetUint("userId")

	// 修复 1: 显式处理错误，防止 v 为 nil 导致的崩溃
	v, err := FindVideoById(uploadVideoReq.Vid)
	if err != nil {
		return errors.New("视频不存在")
	}

	if cache.GetUploadImage(uploadVideoReq.Cover) != userId {
		// 查询是否与旧封面图一致
		if v.Cover != uploadVideoReq.Cover {
			return errors.New("文件链接无效")
		}
	}

	if v.Uid != userId {
		return errors.New("无权修改该视频")
	}

	if v.PartitionId != 0 {
		return errors.New("视频信息已存在")
	}

	if !IsSubpartition(uploadVideoReq.PartitionId, global.CONTENT_TYPE_VIDEO) {
		return errors.New("分区不存在")
	}

	err = global.Mysql.Model(&model.Video{}).Where("id = ? and uid = ?", uploadVideoReq.Vid, userId).Updates(
		map[string]any{
			"title":        uploadVideoReq.Title,
			"cover":        uploadVideoReq.Cover,
			"desc":         uploadVideoReq.Desc,
			"tags":         uploadVideoReq.Tags,
			"copyright":    uploadVideoReq.Copyright,
			"partition_id": uploadVideoReq.PartitionId,
			"status":       getVideoStatus(uploadVideoReq.Vid),
		},
	).Error

	if err != nil {
		utils.ErrorLog("修改视频信息失败", "video", err.Error())
		return errors.New("修改失败")
	}

	// 清除缓存
	cache.DelVideoInfo(uploadVideoReq.Vid)
	SyncVideoTagsFromCSV(uploadVideoReq.Vid, uploadVideoReq.Tags)

	return nil
}

// ParseVideoID 兼容对外 shortId（随机不透明串）与数字自增 id。
func ParseVideoID(raw string) (uint, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, errors.New("视频ID不能为空")
	}

	if utils.IsAllASCIIDigits(raw) {
		id := utils.StringToUint(raw)
		if id == 0 {
			return 0, errors.New("视频不存在")
		}
		return id, nil
	}

	var video model.Video
	if err := global.Mysql.Where("short_id = ?", raw).First(&video).Error; err != nil || video.ID == 0 {
		return 0, errors.New("视频不存在")
	}
	return video.ID, nil
}

// ParseResourceID 兼容对外 shortId（随机不透明串）与数字自增 id。
func ParseResourceID(raw string) (uint, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, errors.New("资源ID不能为空")
	}

	if utils.IsAllASCIIDigits(raw) {
		id := utils.StringToUint(raw)
		if id == 0 {
			return 0, errors.New("资源不存在")
		}
		return id, nil
	}

	var resource model.Resource
	if err := global.Mysql.Where("short_id = ?", raw).First(&resource).Error; err != nil || resource.ID == 0 {
		return 0, errors.New("资源不存在")
	}
	return resource.ID, nil
}

func GetVideoStatus(ctx *gin.Context, vid uint) (video vo.VideoStatusResp, err error) {
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Video{}).Select(vo.VIDEO_STATUS_FIELD).Where("id = ? and uid = ?", vid, userId).Scan(&video)
	if video.ID == 0 {
		return video, errors.New("视频不存在")
	}

	video.Tags = LoadVideoTagNames(vid)
	//查询分区下的视频资源
	video.Resources = GetReviewResourceList(vid)

	return video, nil
}

// GetVideoFile 获取视频文件（返回 DASH MPD / HLS M3U8 / 播放信息 JSON）
func GetVideoFile(ctx *gin.Context, resourceId uint, quality, format string) (string, error) {
	// HLS 子清单无需校验审核状态（已有 key 鉴权）
	if format != "m3u8video" && format != "m3u8audio" {
		if !IsResourceExist(resourceId) {
			return "", errors.New("资源不存在")
		}
	}

	return getVideoFileInternal(ctx, resourceId, quality, format)
}

// =====================================================
// B站风格：SegmentBase 模式（音视频分离，字节范围请求）
// =====================================================

// getMediaFileURL 获取媒体文件的直链 URL（支持本地和OSS）
func getMediaFileURL(dirName, fileName, key string) string {
	if global.Config.Storage.OssType == "local" {
		domain := global.Config.Storage.Domain
		if domain == "" {
			return fmt.Sprintf("/api/v1/video/stream/%s?key=%s", fileName, key)
		}
		return fmt.Sprintf("%s/api/v1/video/stream/%s?key=%s", domain, fileName, key)
	}
	// OSS 存储：直接拼接公开URL
	return global.GetOssUrl("video/" + dirName + "/" + fileName)
}

// getMediaFileURLWithBackup 同 getMediaFileURL，但 useBackup=true 时优先使用备用 OSS URL。
// 备用 OSS 未配置或本地存储时回退到主 URL。
func getMediaFileURLWithBackup(dirName, fileName, key string, useBackup bool) string {
	if !useBackup {
		return getMediaFileURL(dirName, fileName, key)
	}
	if global.Config.Storage.OssType == "local" {
		return getMediaFileURL(dirName, fileName, key)
	}
	backupURL := global.GetBackupOssUrl("video/" + dirName + "/" + fileName)
	if backupURL == "" {
		return getMediaFileURL(dirName, fileName, key)
	}
	return backupURL
}

// buildMPDSegmentBase 生成显示单个画质的 DASH MPD，支持多音轨
func buildMPDSegmentBase(file *model.VideoIndexFile, audioTracks []model.AudioTrack, key string, useBackup ...bool) string {
	backup := len(useBackup) > 0 && useBackup[0]
	durationStr := formatDuration(file.TotalDuration)

	videoURL := getMediaFileURLWithBackup(file.DirName, file.VideoFile, key, backup)
	safeVideoURL := xmlEscape(videoURL)

	var safeBackupVideoURL string
	if !backup && global.StorageBackup != nil {
		if bv := global.GetBackupOssUrl("video/" + file.DirName + "/" + file.VideoFile); bv != "" {
			safeBackupVideoURL = xmlEscape(bv)
		}
	}

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString("\n")
	fmt.Fprintf(&sb, `<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" type="static" mediaPresentationDuration="%s" minBufferTime="PT1.5S" profiles="urn:mpeg:dash:profile:isoff-on-demand:2011">`, durationStr)
	sb.WriteString("\n  <Period>\n")

	// ========== 视频 AdaptationSet ==========
	sb.WriteString(`    <AdaptationSet mimeType="video/mp4" segmentAlignment="true" startWithSAP="1">`)
	sb.WriteString("\n")
	fmt.Fprintf(&sb, `      <Representation id="%s" bandwidth="%d" width="%d" height="%d" frameRate="%.3f" codecs="%s">`,
		file.Quality, file.VideoBandwidth, file.Width, file.Height, file.FrameRate, file.VideoCodec)
	sb.WriteString("\n")
	fmt.Fprintf(&sb, `        <BaseURL>%s</BaseURL>`, safeVideoURL)
	sb.WriteString("\n")
	if safeBackupVideoURL != "" {
		fmt.Fprintf(&sb, `        <BaseURL>%s</BaseURL>`, safeBackupVideoURL)
		sb.WriteString("\n")
	}
	fmt.Fprintf(&sb, `        <SegmentBase indexRange="%s">`, file.VideoIndexRange)
	sb.WriteString("\n")
	fmt.Fprintf(&sb, `          <Initialization range="%s"/>`, file.VideoInitRange)
	sb.WriteString("\n")
	sb.WriteString("        </SegmentBase>\n")
	sb.WriteString("      </Representation>\n")
	sb.WriteString("    </AdaptationSet>\n")

	// ========== 音频 AdaptationSet（多音轨支持） ==========
	writeAudioAdaptationSets(&sb, audioTracks, file, key, backup)

	sb.WriteString("  </Period>\n")
	sb.WriteString("</MPD>")

	return sb.String()
}

// writeAudioAdaptationSets 写入音频 AdaptationSet（多音轨或单音轨回退）
func writeAudioAdaptationSets(sb *strings.Builder, audioTracks []model.AudioTrack, fallbackFile *model.VideoIndexFile, key string, backup bool) {
	if len(audioTracks) > 0 {
		for _, at := range audioTracks {
			audioURL := getMediaFileURLWithBackup(at.DirName, at.AudioFile, key, backup)
			safeAudioURL := xmlEscape(audioURL)

			var safeBackupAudioURL string
			if !backup && global.StorageBackup != nil {
				if ba := global.GetBackupOssUrl("video/" + at.DirName + "/" + at.AudioFile); ba != "" {
					safeBackupAudioURL = xmlEscape(ba)
				}
			}

			label := at.Title
			if label == "" {
				label = "Audio"
			}
			defaultAttr := ""
			if at.IsDefault {
				defaultAttr = ` default="true"`
			}
			fmt.Fprintf(sb, `    <AdaptationSet mimeType="audio/mp4" lang="%s" label="%s" segmentAlignment="true" startWithSAP="1"%s>`,
				xmlEscape(at.Language), xmlEscape(label), defaultAttr)
			sb.WriteString("\n")
			fmt.Fprintf(sb, `      <Representation id="audio-%s" bandwidth="%d" codecs="%s">`,
				xmlEscape(at.Language), at.Bandwidth, at.Codec)
			sb.WriteString("\n")
			fmt.Fprintf(sb, `        <BaseURL>%s</BaseURL>`, safeAudioURL)
			sb.WriteString("\n")
			if safeBackupAudioURL != "" {
				fmt.Fprintf(sb, `        <BaseURL>%s</BaseURL>`, safeBackupAudioURL)
				sb.WriteString("\n")
			}
			if at.IndexRange != "" {
				fmt.Fprintf(sb, `        <SegmentBase indexRange="%s">`, at.IndexRange)
				sb.WriteString("\n")
				fmt.Fprintf(sb, `          <Initialization range="%s"/>`, at.InitRange)
				sb.WriteString("\n")
				sb.WriteString("        </SegmentBase>\n")
			}
			sb.WriteString("      </Representation>\n")
			sb.WriteString("    </AdaptationSet>\n")
		}
	} else {
		// 回退单音轨模式
		audioURL := getMediaFileURLWithBackup(fallbackFile.DirName, fallbackFile.AudioFile, key, backup)
		safeAudioURL := xmlEscape(audioURL)

		var safeBackupAudioURL string
		if !backup && global.StorageBackup != nil {
			if ba := global.GetBackupOssUrl("video/" + fallbackFile.DirName + "/" + fallbackFile.AudioFile); ba != "" {
				safeBackupAudioURL = xmlEscape(ba)
			}
		}

		sb.WriteString(`    <AdaptationSet mimeType="audio/mp4" segmentAlignment="true" startWithSAP="1">`)
		sb.WriteString("\n")
		fmt.Fprintf(sb, `      <Representation id="audio" bandwidth="%d" codecs="%s">`,
			fallbackFile.AudioBandwidth, fallbackFile.AudioCodec)
		sb.WriteString("\n")
		fmt.Fprintf(sb, `        <BaseURL>%s</BaseURL>`, safeAudioURL)
		sb.WriteString("\n")
		if safeBackupAudioURL != "" {
			fmt.Fprintf(sb, `        <BaseURL>%s</BaseURL>`, safeBackupAudioURL)
			sb.WriteString("\n")
		}
		if fallbackFile.AudioIndexRange != "" {
			fmt.Fprintf(sb, `        <SegmentBase indexRange="%s">`, fallbackFile.AudioIndexRange)
			sb.WriteString("\n")
			fmt.Fprintf(sb, `          <Initialization range="%s"/>`, fallbackFile.AudioInitRange)
			sb.WriteString("\n")
			sb.WriteString("        </SegmentBase>\n")
		}
		sb.WriteString("      </Representation>\n")
		sb.WriteString("    </AdaptationSet>\n")
	}
}

// buildMPDSegmentBaseUnified 生成包含所有清晰度的统一 DASH MPD（无缝切换）
// audioTracks 可选：传入则生成多音轨 AdaptationSet，否则回退 files[0] 单音轨
func buildMPDSegmentBaseUnified(files []model.VideoIndexFile, audioTracks []model.AudioTrack, key string, useBackup ...bool) string {
	backup := len(useBackup) > 0 && useBackup[0]
	// 按码率升序排列（dash.js Representation index 0 = 最低画质）
	sort.Slice(files, func(i, j int) bool {
		return dashRepresentationLess(files[i], files[j])
	})

	durationStr := formatDuration(files[0].TotalDuration)

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString("\n")
	fmt.Fprintf(&sb, `<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" type="static" mediaPresentationDuration="%s" minBufferTime="PT1.5S" profiles="urn:mpeg:dash:profile:isoff-on-demand:2011">`, durationStr)
	sb.WriteString("\n  <Period>\n")

	// ========== 视频 AdaptationSet（包含所有清晰度 Representation） ==========
	sb.WriteString(`    <AdaptationSet mimeType="video/mp4" segmentAlignment="true" startWithSAP="1">`)
	sb.WriteString("\n")
	for _, file := range files {
		videoURL := getMediaFileURLWithBackup(file.DirName, file.VideoFile, key, backup)
		safeVideoURL := xmlEscape(videoURL)

		// 备用 OSS 的第二 BaseURL（每个 Representation 独立）
		var safeBackupVideoURL string
		if !backup && global.StorageBackup != nil {
			if bv := global.GetBackupOssUrl("video/" + file.DirName + "/" + file.VideoFile); bv != "" {
				safeBackupVideoURL = xmlEscape(bv)
			}
		}

		fmt.Fprintf(&sb, `      <Representation id="%s" bandwidth="%d" width="%d" height="%d" frameRate="%.3f" codecs="%s">`,
			file.Quality, file.VideoBandwidth, file.Width, file.Height, file.FrameRate, file.VideoCodec)
		sb.WriteString("\n")
		fmt.Fprintf(&sb, `        <BaseURL>%s</BaseURL>`, safeVideoURL)
		sb.WriteString("\n")
		if safeBackupVideoURL != "" {
			fmt.Fprintf(&sb, `        <BaseURL>%s</BaseURL>`, safeBackupVideoURL)
			sb.WriteString("\n")
		}
		fmt.Fprintf(&sb, `        <SegmentBase indexRange="%s">`, file.VideoIndexRange)
		sb.WriteString("\n")
		fmt.Fprintf(&sb, `          <Initialization range="%s"/>`, file.VideoInitRange)
		sb.WriteString("\n")
		sb.WriteString("        </SegmentBase>\n")
		sb.WriteString("      </Representation>\n")
	}
	sb.WriteString("    </AdaptationSet>\n")

	// ========== 音频 AdaptationSet（多音轨支持） ==========
	writeAudioAdaptationSets(&sb, audioTracks, &files[0], key, backup)

	sb.WriteString("  </Period>\n")
	sb.WriteString("</MPD>")

	return sb.String()
}

// dashRepresentationLess 统一 DASH Representation 排序规则（稳定且可预期）
// 规则：码率升序 -> 像素面积升序 -> 帧率升序 -> quality 字符串升序
func dashRepresentationLess(a, b model.VideoIndexFile) bool {
	if a.VideoBandwidth != b.VideoBandwidth {
		return a.VideoBandwidth < b.VideoBandwidth
	}
	aPixels := a.Width * a.Height
	bPixels := b.Width * b.Height
	if aPixels != bPixels {
		return aPixels < bPixels
	}
	if a.FrameRate != b.FrameRate {
		return a.FrameRate < b.FrameRate
	}
	return a.Quality < b.Quality
}

// buildPlayURLJSON 生成类似B站的 playurl JSON 格式
// 当配置了备用 OSS 且不是 backup 模式本身时，同时返回 backupUrl（B站风格多源容灾）。
func buildPlayURLJSON(file *model.VideoIndexFile, key string, useBackup ...bool) string {
	backup := len(useBackup) > 0 && useBackup[0]
	// 【关键】使用直链而不是 API 端点，让 player 可以直接加载
	videoURL := getMediaFileURLWithBackup(file.DirName, file.VideoFile, key, backup)
	audioURL := getMediaFileURLWithBackup(file.DirName, file.AudioFile, key, backup)

	// 非 backup 模式且配置了备用 OSS 时，计算 backupUrl（B站风格的多源容灾）
	backupVideoURL := ""
	backupAudioURL := ""
	if !backup && global.StorageBackup != nil {
		backupVideoURL = getMediaFileURLWithBackup(file.DirName, file.VideoFile, key, true)
		backupAudioURL = getMediaFileURLWithBackup(file.DirName, file.AudioFile, key, true)
	}

	backupVideoJSON := ""
	backupAudioJSON := ""
	if backupVideoURL != "" {
		backupVideoJSON = fmt.Sprintf(`, "backupUrl": "%s"`, xmlEscape(backupVideoURL))
	}
	if backupAudioURL != "" {
		backupAudioJSON = fmt.Sprintf(`, "backupUrl": "%s"`, xmlEscape(backupAudioURL))
	}

	json := fmt.Sprintf(`{
  "code": 0,
  "message": "OK",
  "data": {
    "quality": "%s",
    "duration": %.3f,
    "dash": {
      "duration": %.0f,
      "minBufferTime": 1.5,
      "video": [{
        "id": "%s",
        "baseUrl": "%s",
        "bandwidth": %d,
        "mimeType": "video/mp4",
        "codecs": "%s",
        "width": %d,
        "height": %d,
        "frameRate": "%.3f",
        "SegmentBase": {
          "Initialization": "%s",
          "indexRange": "%s"
        }%s
      }],
      "audio": [{
        "id": "audio",
        "baseUrl": "%s",
        "bandwidth": %d,
        "mimeType": "audio/mp4",
        "codecs": "%s",
        "SegmentBase": {
          "Initialization": "%s",
          "indexRange": "%s"
        }%s
      }]
    }
  }
}`,
		xmlEscape(file.Quality), file.TotalDuration, file.TotalDuration,
		xmlEscape(file.Quality), xmlEscape(videoURL), file.VideoBandwidth, xmlEscape(file.VideoCodec), file.Width, file.Height, file.FrameRate,
		file.VideoInitRange, file.VideoIndexRange, backupVideoJSON,
		xmlEscape(audioURL), file.AudioBandwidth, xmlEscape(file.AudioCodec),
		file.AudioInitRange, file.AudioIndexRange, backupAudioJSON,
	)

	return json
}

// =====================================================
// HLS v7 over fMP4（SegmentBase 模式的 Safari/iOS 兼容）
// =====================================================

// sidxEntry 表示 sidx box 中的一个引用条目（一个 fragment）
type sidxEntry struct {
	Offset   int64   // fragment 在文件中的字节偏移
	Size     int64   // fragment 的字节大小
	Duration float64 // fragment 的时长（秒）
}

// parseSidxBox 从 fMP4 文件中解析 sidx box，提取所有 fragment 的字节范围
// sidx box 结构参考 ISO 14496-12
func parseSidxBox(filePath string) ([]sidxEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// 遍历顶层 box，找到 sidx
	offset := int64(0)

	for offset < fileSize {
		boxSize, boxType, err := readMP4BoxSizeAndType(file, offset, fileSize)
		if err != nil || boxSize <= 0 {
			break
		}

		if boxType == "sidx" {
			return decodeSidxBox(file, offset, boxSize)
		}

		// 遇到 moof 就停止搜索
		if boxType == "moof" {
			break
		}

		offset += boxSize
	}

	return nil, fmt.Errorf("sidx box not found in %s", filePath)
}

// decodeSidxBox 解码 sidx box 内容
func decodeSidxBox(reader io.ReaderAt, boxOffset, boxSize int64) ([]sidxEntry, error) {
	data := make([]byte, boxSize)
	if _, err := reader.ReadAt(data, boxOffset); err != nil {
		return nil, err
	}

	// box header: 8 bytes (size + type)
	pos := 8

	// FullBox header: version (1 byte) + flags (3 bytes)
	version := data[pos]
	pos += 4 // skip version + flags

	// reference_ID (4 bytes)
	pos += 4

	// timescale (4 bytes)
	var timescale uint32
	timescale = binary.BigEndian.Uint32(data[pos : pos+4])
	pos += 4

	if timescale == 0 {
		return nil, fmt.Errorf("sidx timescale is 0")
	}

	// earliest_presentation_time + first_offset
	var firstOffset int64
	if version == 0 {
		// 32-bit each
		pos += 4 // earliest_presentation_time
		firstOffset = int64(binary.BigEndian.Uint32(data[pos : pos+4]))
		pos += 4
	} else {
		// 64-bit each
		pos += 8 // earliest_presentation_time
		firstOffset = int64(binary.BigEndian.Uint64(data[pos : pos+8]))
		pos += 8
	}

	// reserved (2 bytes) + reference_count (2 bytes)
	pos += 2 // reserved
	refCount := int(binary.BigEndian.Uint16(data[pos : pos+2]))
	pos += 2

	// 第一个 fragment 的偏移 = sidx box 结束位置 + first_offset
	// ISO 14496-12: anchor_point = first byte after sidx box, first_offset 是额外偏移
	fragmentOffset := boxOffset + boxSize + firstOffset

	entries := make([]sidxEntry, 0, refCount)

	for range refCount {
		if pos+12 > len(data) {
			break
		}

		// reference_type (1 bit) + referenced_size (31 bits)
		refField := binary.BigEndian.Uint32(data[pos : pos+4])
		// referenceType := (refField >> 31) & 1  // 0 = media, 1 = index
		referencedSize := int64(refField & 0x7FFFFFFF)
		pos += 4

		// subsegment_duration (4 bytes)
		subsegmentDuration := binary.BigEndian.Uint32(data[pos : pos+4])
		pos += 4

		// SAP fields (4 bytes) - skip
		pos += 4

		duration := float64(subsegmentDuration) / float64(timescale)

		entries = append(entries, sidxEntry{
			Offset:   fragmentOffset,
			Size:     referencedSize,
			Duration: duration,
		})

		fragmentOffset += referencedSize
	}

	return entries, nil
}

// buildM3U8MasterSegmentBase 为 SegmentBase 模式生成 HLS v7 主播放列表（Master Playlist）
// Safari/iOS 通过此清单实现音视频分离播放
// format=m3u8 返回主清单，内部引用 m3u8video / m3u8audio 子清单
func buildM3U8MasterSegmentBase(file *model.VideoIndexFile, resourceId uint, key string, useBackup ...bool) string {
	backup := len(useBackup) > 0 && useBackup[0]
	var sb strings.Builder

	sb.WriteString("#EXTM3U\n")
	sb.WriteString("#EXT-X-VERSION:7\n")
	sb.WriteString("\n")

	// 【关键】构建子清单的完整 URL
	baseURL := getLocalBaseURL()

	// 备份后缀：通知子清单生成时也使用备用 OSS URL
	backupSuffix := ""
	if backup {
		backupSuffix = "&backup=true"
	}

	// 音频组：引用音频子清单
	audioURI := fmt.Sprintf("/api/v1/video/getVideoFile?resourceId=%d&quality=%s&format=m3u8audio&key=%s%s", resourceId, file.Quality, key, backupSuffix)
	if baseURL != "" {
		audioURI = baseURL + audioURI
	}
	fmt.Fprintf(&sb,
		"#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID=\"audio\",NAME=\"Audio\",DEFAULT=YES,AUTOSELECT=YES,URI=\"%s\"\n",
		audioURI)
	sb.WriteString("\n")

	// 视频流：引用视频子清单
	videoURI := fmt.Sprintf("/api/v1/video/getVideoFile?resourceId=%d&quality=%s&format=m3u8video&key=%s%s", resourceId, file.Quality, key, backupSuffix)
	if baseURL != "" {
		videoURI = baseURL + videoURI
	}
	fmt.Fprintf(&sb, "#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d,FRAME-RATE=%.3f,CODECS=\"%s,%s\",AUDIO=\"audio\"\n",
		file.VideoBandwidth+file.AudioBandwidth, file.Width, file.Height, file.FrameRate,
		file.VideoCodec, file.AudioCodec)
	sb.WriteString(videoURI)
	sb.WriteByte('\n')

	return sb.String()
}

// buildM3U8VideoSegmentBase 为视频流生成 HLS v7 媒体播放列表（从文件实时解析 sidx）
func buildM3U8VideoSegmentBase(file *model.VideoIndexFile, key string, useBackup ...bool) (string, error) {
	backup := len(useBackup) > 0 && useBackup[0]
	videoPath := "./upload/video/" + file.DirName + "/" + file.VideoFile
	entries, err := parseSidxBox(videoPath)
	if err != nil {
		return "", fmt.Errorf("parse video sidx failed: %w", err)
	}
	streamURL := getMediaFileURLWithBackup(file.DirName, file.VideoFile, key, backup)
	return buildByteRangeM3U8(entries, streamURL, file.VideoInitRange), nil
}

// buildM3U8AudioSegmentBase 为音频流生成 HLS v7 媒体播放列表（从文件实时解析 sidx）
func buildM3U8AudioSegmentBase(file *model.VideoIndexFile, key string, useBackup ...bool) (string, error) {
	backup := len(useBackup) > 0 && useBackup[0]
	audioPath := "./upload/video/" + file.DirName + "/" + file.AudioFile
	entries, err := parseSidxBox(audioPath)
	if err != nil {
		return "", fmt.Errorf("parse audio sidx failed: %w", err)
	}
	streamURL := getMediaFileURLWithBackup(file.DirName, file.AudioFile, key, backup)
	return buildByteRangeM3U8(entries, streamURL, file.AudioInitRange), nil
}

// buildByteRangeM3U8 从 sidx 条目生成带 EXT-X-BYTERANGE 的 HLS v7 媒体播放列表
func buildByteRangeM3U8(entries []sidxEntry, streamURL, initRange string) string {
	var sb strings.Builder

	sb.WriteString("#EXTM3U\n")
	sb.WriteString("#EXT-X-VERSION:7\n")

	// 计算最大 segment 时长
	maxDuration := 0.0
	for _, e := range entries {
		if e.Duration > maxDuration {
			maxDuration = e.Duration
		}
	}
	fmt.Fprintf(&sb, "#EXT-X-TARGETDURATION:%d\n", int(maxDuration)+1)
	sb.WriteString("#EXT-X-MEDIA-SEQUENCE:0\n")
	sb.WriteString("#EXT-X-PLAYLIST-TYPE:VOD\n")

	// EXT-X-MAP: 初始化段（ftyp + moov）
	// initRange 格式为 "0-927"，需要转换为 BYTERANGE 的 length@offset
	initParts := strings.Split(initRange, "-")
	if len(initParts) == 2 {
		start := parseRangeInt(initParts[0])
		end := parseRangeInt(initParts[1])
		length := end - start + 1
		fmt.Fprintf(&sb, "#EXT-X-MAP:URI=\"%s\",BYTERANGE=\"%d@%d\"\n", streamURL, length, start)
	}

	// 每个 fragment 作为一个 segment
	for _, entry := range entries {
		fmt.Fprintf(&sb, "#EXTINF:%.6f,\n", entry.Duration)
		fmt.Fprintf(&sb, "#EXT-X-BYTERANGE:%d@%d\n", entry.Size, entry.Offset)
		sb.WriteString(streamURL)
		sb.WriteByte('\n')
	}

	sb.WriteString("#EXT-X-ENDLIST\n")
	return sb.String()
}

// parseRangeInt 将范围字符串解析为 int64
func parseRangeInt(s string) int64 {
	val := int64(0)
	for _, c := range s {
		if c >= '0' && c <= '9' {
			val = val*10 + int64(c-'0')
		}
	}
	return val
}

// =====================================================
// 兼容模式：SegmentList 切片模式
// =====================================================

// buildM3U8SegmentList 从切片元数据生成 m3u8 (HLS)
func buildM3U8SegmentList(file *model.VideoIndexFile, key string) string {
	var sb strings.Builder

	isFMP4 := file.InitFile != ""

	if isFMP4 {
		sb.WriteString("#EXTM3U\n")
		sb.WriteString("#EXT-X-VERSION:7\n")
	} else {
		sb.WriteString("#EXTM3U\n")
		sb.WriteString("#EXT-X-VERSION:3\n")
	}

	fmt.Fprintf(&sb, "#EXT-X-TARGETDURATION:%d\n", int(file.SegmentDuration)+1)
	sb.WriteString("#EXT-X-MEDIA-SEQUENCE:0\n")
	sb.WriteString("#EXT-X-PLAYLIST-TYPE:VOD\n")

	// 【关键】构建基础 URL
	baseURL := getLocalBaseURL()

	if isFMP4 {
		initURI := fmt.Sprintf("/api/v1/video/slice/%s?key=%s", file.InitFile, key)
		if baseURL != "" {
			initURI = baseURL + initURI
		}
		fmt.Fprintf(&sb, "#EXT-X-MAP:URI=\"%s\"\n", initURI)
	}

	for i := range file.SegmentCount {
		duration := file.SegmentDuration
		if i == file.SegmentCount-1 {
			duration = file.LastSegmentDuration
		}

		ext := ".ts"
		if isFMP4 {
			ext = ".m4s"
		}
		fileName := fmt.Sprintf("%s_%05d%s", file.Quality, i, ext)

		fmt.Fprintf(&sb, "#EXTINF:%.3f,\n", duration)

		segmentURI := fmt.Sprintf("/api/v1/video/slice/%s?key=%s", fileName, key)
		if baseURL != "" {
			segmentURI = baseURL + segmentURI
		}
		sb.WriteString(segmentURI)
		sb.WriteByte('\n')
	}

	sb.WriteString("#EXT-X-ENDLIST\n")
	return sb.String()
}

// buildMPDSegmentList 从切片元数据生成 mpd (DASH SegmentList)
func buildMPDSegmentList(file *model.VideoIndexFile, key string) string {
	durationStr := formatDuration(file.TotalDuration)

	codec := file.Codec
	if codec == "" {
		codec = "avc1.640028"
	}

	// 【关键】构建基础 URL
	baseURL := getLocalBaseURL()

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString("\n")
	fmt.Fprintf(&sb, `<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" type="static" mediaPresentationDuration="%s" minBufferTime="PT2S" profiles="urn:mpeg:dash:profile:isoff-on-demand:2011">`, durationStr)
	sb.WriteString("\n  <Period>\n")
	sb.WriteString(`    <AdaptationSet mimeType="video/mp4" segmentAlignment="true" startWithSAP="1">`)
	sb.WriteString("\n")
	fmt.Fprintf(&sb, `      <Representation id="%s" bandwidth="%d" width="%d" height="%d" frameRate="%.0f" codecs="%s">`,
		file.Quality, file.Bandwidth, file.Width, file.Height, file.FrameRate, codec)
	sb.WriteString("\n")

	fmt.Fprintf(&sb, `        <SegmentList timescale="1000" duration="%d">`, int(file.SegmentDuration*1000))
	sb.WriteString("\n")
	if file.InitFile != "" {
		initURL := fmt.Sprintf("/api/v1/video/slice/%s?key=%s", file.InitFile, key)
		if baseURL != "" {
			initURL = baseURL + initURL
		}
		fmt.Fprintf(&sb, `          <Initialization sourceURL="%s"/>`, initURL)
		sb.WriteString("\n")
	}

	for i := range file.SegmentCount {
		ext := ".m4s"
		if file.InitFile == "" {
			ext = ".ts"
		}
		fileName := fmt.Sprintf("%s_%05d%s", file.Quality, i, ext)
		segmentURL := fmt.Sprintf("/api/v1/video/slice/%s?key=%s", fileName, key)
		if baseURL != "" {
			segmentURL = baseURL + segmentURL
		}
		fmt.Fprintf(&sb, "          <SegmentURL media=\"%s\"/>\n", segmentURL)
	}

	sb.WriteString("        </SegmentList>\n")
	sb.WriteString("      </Representation>\n")
	sb.WriteString("    </AdaptationSet>\n")
	sb.WriteString("  </Period>\n")
	sb.WriteString("</MPD>")

	return sb.String()
}

// =====================================================
// 兼容旧数据
// =====================================================

// buildFromLegacyContent 兼容旧数据（Content字段存储完整m3u8）
func buildFromLegacyContent(file *model.VideoIndexFile, key, format string) (string, error) {
	if file.Content == "" {
		return "", fmt.Errorf("no content found")
	}

	if format == "mpd" {
		return "", fmt.Errorf("legacy data only supports m3u8")
	}

	var res strings.Builder
	for _, line := range strings.Split(file.Content, "\n") {
		if strings.Contains(line, ".ts") || strings.Contains(line, ".m4s") {
			res.WriteString("/api/v1/video/slice/")
			res.WriteString(line)
			res.WriteString("?key=")
			res.WriteString(key)
			res.WriteByte('\n')
		} else {
			res.WriteString(line)
			res.WriteByte('\n')
		}
	}
	return res.String(), nil
}

// formatDuration 格式化时长为 ISO 8601 格式（四舍五入到整秒，避免截断导致 MPD 总时长短于实际、播放缺尾）
func formatDuration(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := seconds - float64(hours*3600+minutes*60)
	return fmt.Sprintf("PT%dH%dM%.3fS", hours, minutes, secs)
}

// 获取视频文件（后台管理，不校验审核状态）
func GetVideoFileManage(ctx *gin.Context, resourceId uint, quality, format string) (string, error) {
	return getVideoFileInternal(ctx, resourceId, quality, format)
}

// getVideoFileInternal 获取视频文件的公共逻辑（DASH MPD / HLS M3U8 / JSON）
// 支持 backup=true 查询参数：生成的 manifest 中所有媒体 URL 指向备用 OSS（播放容灾）。
func getVideoFileInternal(ctx *gin.Context, resourceId uint, quality, format string) (string, error) {
	useBackup := ctx.Query("backup") == "true"

	// ========== HLS 子清单请求（m3u8video / m3u8audio） ==========
	if format == "m3u8video" || format == "m3u8audio" {
		existingKey := ctx.Query("key")
		if existingKey == "" || resourceId == 0 {
			return "", errors.New("参数无效")
		}

		var file model.VideoIndexFile
		if err := global.Mysql.Where("resource_id = ? AND quality = ?", resourceId, quality).First(&file).Error; err != nil {
			return "", errors.New("视频索引不存在")
		}

		if cache.GetVideoSlice(existingKey) == "" {
			cache.SetVideoSlice(existingKey, file.DirName)
		}

		if format == "m3u8video" {
			return buildM3U8VideoSegmentBase(&file, existingKey, useBackup)
		}
		return buildM3U8AudioSegmentBase(&file, existingKey, useBackup)
	}

	// ========== 常规请求 ==========
	key := uuid.New().String()

	// 统一 DASH MPD（所有清晰度合并到一个 MPD）
	if format == "dash-unified" {
		var files []model.VideoIndexFile
		if err := global.Mysql.Where("resource_id = ?", resourceId).Find(&files).Error; err != nil || len(files) == 0 {
			return "", errors.New("视频索引不存在")
		}
		var sbFiles []model.VideoIndexFile
		for _, f := range files {
			if f.IsSegmentBase() {
				sbFiles = append(sbFiles, f)
			}
		}
		if len(sbFiles) == 0 {
			return "", errors.New("无SegmentBase资源")
		}
		cache.SetVideoSlice(key, sbFiles[0].DirName)
		var audioTracks []model.AudioTrack
		global.Mysql.Where("resource_id = ?", resourceId).Find(&audioTracks)
		return buildMPDSegmentBaseUnified(sbFiles, audioTracks, key, useBackup), nil
	}

	// 从 VideoIndexFile 元数据动态生成
	var file model.VideoIndexFile
	err := global.Mysql.Where("resource_id = ? AND quality = ?", resourceId, quality).First(&file).Error
	if err != nil {
		return "", errors.New("视频索引不存在")
	}

	cache.SetVideoSlice(key, file.DirName)

	// B站风格：SegmentBase 模式（音视频分离）
	if file.IsSegmentBase() {
		if format == "m3u8" {
			return buildM3U8MasterSegmentBase(&file, resourceId, key, useBackup), nil
		}
		if format == "dash" || format == "mpd" {
			var audioTracks []model.AudioTrack
			global.Mysql.Where("resource_id = ?", resourceId).Find(&audioTracks)
			return buildMPDSegmentBase(&file, audioTracks, key, useBackup), nil
		}
		return buildPlayURLJSON(&file, key, useBackup), nil
	}

	// 兼容模式：SegmentList 切片模式
	if file.IsSegmentList() {
		if format == "mpd" {
			return buildMPDSegmentList(&file, key), nil
		}
		return buildM3U8SegmentList(&file, key), nil
	}

	// 兼容旧数据（Content字段存储完整m3u8）
	return buildFromLegacyContent(&file, key, format)
}

// 获取视频切所在文件目录
func GetVideoSliceDir(key string) string {
	return cache.GetVideoSlice(key)
}

// 获取自己上传的视频
// category: all | published(0) | transcoding(100,200,300) | transcode_failed(3000) | pending(500) | rejected(2000)
func GetUploadVideoList(ctx *gin.Context, page, pageSize int, category string) (total int64, videos []vo.UploadVideoResp) {
	userId := ctx.GetUint("userId")

	// 个人投稿列表仅展示 UGC 视频，排除已绑定 PGC 剧集的视频
	db := global.Mysql.Model(&model.Video{}).Where("uid = ? AND pgc_attached = ?", userId, false)
	db = applyVideoCategoryFilter(db, category)
	db.Count(&total)
	db.Select(vo.UPLOAD_VIDEO_FIELD).
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&videos)

	// 确保 shortId 有值
	for i := range videos {
		if videos[i].ShortID == "" {
			videos[i].ShortID = utils.UintToString(videos[i].ID)
		}
	}

	// 更新播放量数据并收集需要返回转码进度的视频
	transcodingVideoIDs := make([]uint, 0)
	for i := 0; i < len(videos); i++ {
		videos[i].Clicks += GetVideoClicks(videos[i].ID)
		if isTranscodingStatus(videos[i].Status) {
			transcodingVideoIDs = append(transcodingVideoIDs, videos[i].ID)
		}
	}

	progressMap := collectVideoProgressSnapshots(transcodingVideoIDs)
	for i := 0; i < len(videos); i++ {
		if !isTranscodingStatus(videos[i].Status) {
			continue
		}
		if snapshot, ok := progressMap[videos[i].ID]; ok {
			videos[i].TranscodingProgress = snapshot.Overall
			videos[i].TranscodingDetails = snapshot.Details
			videos[i].UploadProgress = snapshot.UploadInfo
			if videos[i].UploadProgress != nil {
				videos[i].UploadProgress.OssType = "" // 用户端不暴露 OSS 类型
			}
		}
	}

	return
}

func applyVideoCategoryFilter(db *gorm.DB, category string) *gorm.DB {
	switch category {
	case "published":
		return db.Where("status = ?", global.AUDIT_APPROVED)
	case "transcoding":
		return db.Where("status IN ?", []int{global.CREATED_VIDEO, global.VIDEO_PROCESSING, global.SUBMIT_REVIEW})
	case "transcode_failed":
		return db.Where("status = ?", global.PROCESSING_FAIL)
	case "pending":
		return db.Where("status = ?", global.WAITING_REVIEW)
	case "rejected":
		return db.Where("status = ?", global.REVIEW_FAILED)
	default:
		return db
	}
}

type videoProgressSnapshot struct {
	Overall    float64
	Details    []vo.TranscodingProgressItem
	UploadInfo *vo.UploadProgressInfo
}

func isTranscodingStatus(status int) bool {
	return status == global.CREATED_VIDEO || status == global.VIDEO_PROCESSING || status == global.SUBMIT_REVIEW || status == global.WAITING_REVIEW
}

func collectVideoProgressSnapshots(videoIDs []uint) map[uint]videoProgressSnapshot {
	progressMap := make(map[uint]videoProgressSnapshot, len(videoIDs))
	if len(videoIDs) == 0 {
		return progressMap
	}

	resourceIDSet := make(map[uint]struct{})
	for _, videoID := range videoIDs {
		overall, details, uploadInfo := GetVideoTranscodingProgress(videoID)
		if len(details) == 0 {
			progressMap[videoID] = videoProgressSnapshot{Overall: overall, Details: details, UploadInfo: uploadInfo}
			continue
		}
		for _, item := range details {
			resourceIDSet[item.ResourceID] = struct{}{}
		}
		progressMap[videoID] = videoProgressSnapshot{Overall: overall, Details: details, UploadInfo: uploadInfo}
	}

	if len(resourceIDSet) == 0 {
		return progressMap
	}

	resourceIDs := make([]uint, 0, len(resourceIDSet))
	for resourceID := range resourceIDSet {
		resourceIDs = append(resourceIDs, resourceID)
	}

	var resources []model.Resource
	global.Mysql.Model(&model.Resource{}).Select("id,title").Where("id IN ?", resourceIDs).Find(&resources)

	titleMap := make(map[uint]string, len(resources))
	for _, resource := range resources {
		titleMap[resource.ID] = resource.Title
	}

	for videoID, snapshot := range progressMap {
		if len(snapshot.Details) == 0 {
			continue
		}
		for idx := range snapshot.Details {
			snapshot.Details[idx].ResourceTitle = titleMap[snapshot.Details[idx].ResourceID]
		}
		progressMap[videoID] = snapshot
	}

	return progressMap
}

func EditVideoInfo(ctx *gin.Context, editVideoReq dto.EditVideoReq) error {
	userId := ctx.GetUint("userId")
	oldVideo, err := FindVideoById(editVideoReq.Vid)
	if err != nil {
		return errors.New("视频不存在")
	}
	if oldVideo.Uid != userId {
		return errors.New("无权修改该视频")
	}
	if cache.GetUploadImage(editVideoReq.Cover) != userId {
		// 查询是否与旧封面图一致
		if oldVideo.Cover != editVideoReq.Cover {
			return errors.New("文件链接无效")
		}
	}

	// 准备更新的字段
	updateData := map[string]any{
		"title": editVideoReq.Title,
		"cover": editVideoReq.Cover,
		"desc":  editVideoReq.Desc,
		"tags":  editVideoReq.Tags,
	}

	// 重新计算视频状态（所有编辑都需要重新审核，防止替换违规内容）
	newStatus := getVideoStatus(editVideoReq.Vid)

	// 特殊处理：如果视频原本已审核通过，编辑后需要重新审核
	// 设置为WAITING_REVIEW(500)而不是根据资源状态计算，避免因为有转码中资源而变成SUBMIT_REVIEW(300)
	if oldVideo.Status == global.AUDIT_APPROVED {
		newStatus = global.WAITING_REVIEW
		utils.InfoLog(fmt.Sprintf("已发布视频被编辑，VideoID=%d，状态从AUDIT_APPROVED(0)改为WAITING_REVIEW(500)，需要重新审核", editVideoReq.Vid), "video")
	}

	updateData["status"] = newStatus

	if err := global.Mysql.Model(&model.Video{}).Where("id = ? and uid = ?", editVideoReq.Vid, userId).Updates(updateData).Error; err != nil {
		utils.ErrorLog("修改视频失败", "video", err.Error())
		return errors.New("修改失败")
	}

	// 如果是已发布视频被编辑，需要从分区列表中移除（因为状态变为待审核）
	if oldVideo.Status == global.AUDIT_APPROVED {
		cache.DelVideoId(oldVideo.PartitionId, oldVideo.ID)
	}

	// 删除视频信息缓存
	cache.DelVideoInfo(editVideoReq.Vid)
	SyncVideoTagsFromCSV(editVideoReq.Vid, editVideoReq.Tags)

	return nil
}

func DeleteVideo(ctx *gin.Context, id uint) error {
	var video model.Video
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Video{}).Where("id = ? and uid = ?", id, userId).First(&video)
	if video.ID == 0 {
		return errors.New("视频不存在")
	}

	return deleteVideoAndRelatedData(id, video.Uid, &video)
}

// deleteVideoAndRelatedData 删除视频及所有关联数据（公共逻辑）
// ownerUid: 视频所有者的UID（用于引用计数处理）
// video: 视频信息（可为nil，此时不清理分区/热门缓存）
func deleteVideoAndRelatedData(id, ownerUid uint, video *model.Video) error {
	// 停止正在进行的转码进程并清理文件
	StopTranscodingAndCleanup(id)

	// 先查询关联的资源记录，用于后续删除相关文件
	var resources []model.Resource
	global.Mysql.Where("vid = ?", id).Find(&resources)

	// 删除关联的m3u8索引文件记录和处理视频文件记录（通过resource_id关联）
	for _, resource := range resources {
		// 查找VideoIndexFile获取DirName
		var indexFile model.VideoIndexFile
		global.Mysql.Where("resource_id = ?", resource.ID).First(&indexFile)

		// 处理 VideoFile 引用计数（全局去重）
		if resource.FileID != 0 {
			decreaseVideoFileRefCount(resource.FileID, ownerUid, resource.ID, indexFile.DirName)
		} else if indexFile.DirName != "" {
			// 兼容旧数据：通过 DirName 查找 VideoFile
			if vf, err := findVideoFileByDirName(global.Mysql, indexFile.DirName); err == nil && vf != nil {
				decreaseVideoFileRefCount(vf.ID, ownerUid, resource.ID, indexFile.DirName)
			}
		}

		// 删除VideoIndexFile记录
		if err := global.Mysql.Where("resource_id = ?", resource.ID).Delete(&model.VideoIndexFile{}).Error; err != nil {
			utils.ErrorLog("删除视频关联m3u8索引文件失败", "video", err.Error())
		}
	}

	// 删除关联的资源记录
	if err := global.Mysql.Where("vid = ?", id).Delete(&model.Resource{}).Error; err != nil {
		utils.ErrorLog("删除视频关联资源失败", "video", err.Error())
	}

	// 删除关联的评论及级联点赞/点赞消息/回复消息
	CleanupContentComments(id, global.CONTENT_TYPE_VIDEO)

	// 删除关联的弹幕
	if video.ShortID != "" {
		if err := global.Mysql.Where("video_short_id = ?", video.ShortID).Delete(&model.Danmaku{}).Error; err != nil {
			utils.ErrorLog("删除视频关联弹幕失败", "video", err.Error())
		}
	}

	// 删除关联的收藏记录
	if err := global.Mysql.Where("vid = ?", id).Delete(&model.CollectVideo{}).Error; err != nil {
		utils.ErrorLog("删除视频关联收藏记录失败", "video", err.Error())
	}

	// 删除关联的点赞记录
	if err := global.Mysql.Where("vid = ?", id).Delete(&model.LikeVideo{}).Error; err != nil {
		utils.ErrorLog("删除视频关联点赞记录失败", "video", err.Error())
	}

	// 删除关联的历史记录
	if err := global.Mysql.Where("vid = ?", id).Delete(&model.History{}).Error; err != nil {
		utils.ErrorLog("删除视频关联历史记录失败", "video", err.Error())
	}

	// 删除关联的AT消息（cid是内容ID，type=0表示视频）
	if err := global.Mysql.Where("cid = ? and type = 0", id).Delete(&model.AtMessage{}).Error; err != nil {
		utils.ErrorLog("删除视频关联AT消息失败", "video", err.Error())
	}

	// 清理 PGC 关联：删除引用该视频的剧集记录，并同步回写对应 PGC 的当前集数
	// 防止后台/脚本直接删视频后，pgc_episode 残留脏数据。
	var pgcEpisodeRows []struct {
		PGCID uint64 `gorm:"column:pgc_id"`
	}
	if err := global.Mysql.Model(&model.PGCEpisode{}).
		Select("pgc_id").
		Where("vid = ?", id).
		Scan(&pgcEpisodeRows).Error; err != nil {
		utils.ErrorLog("查询视频关联PGC剧集失败", "video", err.Error())
	}
	if len(pgcEpisodeRows) > 0 {
		if err := global.Mysql.Where("vid = ?", id).Delete(&model.PGCEpisode{}).Error; err != nil {
			utils.ErrorLog("删除视频关联PGC剧集失败", "video", err.Error())
		} else {
			pgcSet := make(map[uint64]struct{}, len(pgcEpisodeRows))
			for _, row := range pgcEpisodeRows {
				if row.PGCID > 0 {
					pgcSet[row.PGCID] = struct{}{}
				}
			}
			for pgcID := range pgcSet {
				var currEp int64
				if err := global.Mysql.Model(&model.PGCEpisode{}).
					Where("pgc_id = ?", pgcID).
					Count(&currEp).Error; err != nil {
					utils.ErrorLog("统计PGC当前集数失败", "video", err.Error())
					continue
				}
				if err := global.Mysql.Model(&model.PGCContent{}).
					Where("pgc_id = ?", pgcID).
					Update("current_episodes", int(currEp)).Error; err != nil {
					utils.ErrorLog("回写PGC当前集数失败", "video", err.Error())
				}
			}
		}
	}

	if err := global.Mysql.Where("id = ?", id).Delete(&model.Video{}).Error; err != nil {
		utils.ErrorLog("删除视频失败", "video", err.Error())
		return errors.New("删除视频失败")
	}

	// 清理所有相关缓存
	if video != nil && video.ID != 0 {
		cache.DelVideoId(video.PartitionId, video.ID)
		cache.DelSingleHotVideoId(video.ID)
	}
	cache.DelVideoInfo(id)

	return nil
}

// 获取视频信息
func GetVideoById(ctx *gin.Context, videoId uint) (vo.VideoResp, error) {
	video := GetVideoInfo(videoId)
	if video.ID == 0 {
		return video, errors.New("视频信息不存在")
	}

	// 增加播放量(一个ip在同一个视频下，每30分钟可重新增加1播放量)
	AddVideoClicks(videoId, ctx.ClientIP())
	video.Clicks += GetVideoClicks(video.ID)

	var fans int64
	global.Mysql.Model(&model.Relation{}).
		Where("target_uid = ? and (relation = ? or relation = ?)", video.Uid, global.FOLLOWED, global.MUTUAL_FANS).
		Count(&fans)
	video.Fans = fans

	return video, nil
}

// 获取所有的视频列表
func GetAllVideoList(ctx *gin.Context, page, pageSize int) (total int64, videos []vo.AllVideoResp) {
	userId := ctx.GetUint("userId")
	global.Mysql.Model(&model.Video{}).Where("uid = ?", userId).Count(&total)
	global.Mysql.Model(&model.Video{}).Select("`id`,`title`,`cover`,`short_id`").
		Where("uid = ?", userId).
		Order("created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).Scan(&videos)

	for i := range videos {
		if videos[i].ShortID == "" {
			videos[i].ShortID = utils.UintToString(videos[i].ID)
		}
	}

	return
}

// 获取用户视频
func GetVideoByUser(ctx *gin.Context, userId uint, page, pageSize int) (total int64, videos []vo.UploadVideoResp) {
	global.Mysql.Model(&model.Video{}).
		Where("uid = ? and status = ? and pgc_attached = ?", userId, global.AUDIT_APPROVED, false).Count(&total)
	global.Mysql.Model(&model.Video{}).Select(vo.UPLOAD_VIDEO_FIELD).
		Where("uid = ? and status = ? and pgc_attached = ?", userId, global.AUDIT_APPROVED, false).
		Order("created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).Scan(&videos)

	// 更新播放量数据
	for i := 0; i < len(videos); i++ {
		videos[i].Clicks += GetVideoClicks(videos[i].ID)
	}

	return
}

// 获取视频列表(后台管理) - 仅展示 UGC 内容，PGC 内容由独立管理入口管理
func GetVideoListManage(videoListReq dto.VideoListReq) (total int64, videos []vo.VideoInfoManageResp) {
	baseWhere := "status = ? AND pgc_attached = ?"
	global.Mysql.Model(&model.Video{}).Where(baseWhere, global.AUDIT_APPROVED, false).Count(&total)
	global.Mysql.Model(&model.Video{}).Where(baseWhere, global.AUDIT_APPROVED, false).
		Order("created_at DESC").
		Limit(videoListReq.PageSize).Offset((videoListReq.Page - 1) * videoListReq.PageSize).Scan(&videos)

	// 更新播放量和作者数据
	for i := 0; i < len(videos); i++ {
		videos[i].Clicks += GetVideoClicks(videos[i].ID)
		videos[i].Author = GetUserBaseInfo(videos[i].Uid)
		videos[i].Tags = LoadVideoTagNames(videos[i].ID)
	}

	return
}

// 获取处理失败的视频列表（后台管理）
func GetFailedVideoList(videoListReq dto.VideoListReq) (total int64, videos []vo.VideoInfoManageResp) {
	baseWhere := "status = ? AND pgc_attached = ?"
	global.Mysql.Model(&model.Video{}).Where(baseWhere, global.PROCESSING_FAIL, false).Count(&total)
	global.Mysql.Model(&model.Video{}).Where(baseWhere, global.PROCESSING_FAIL, false).
		Order("created_at DESC").
		Limit(videoListReq.PageSize).Offset((videoListReq.Page - 1) * videoListReq.PageSize).Scan(&videos)

	for i := 0; i < len(videos); i++ {
		videos[i].Clicks += GetVideoClicks(videos[i].ID)
		videos[i].Author = GetUserBaseInfo(videos[i].Uid)
		videos[i].Tags = LoadVideoTagNames(videos[i].ID)

		// 加载每个失败视频的分P资源状态，用于展示具体失败的分P
		var resources []model.Resource
		global.Mysql.Where("vid = ?", videos[i].ID).Order("sort_order ASC").Find(&resources)
		if len(resources) > 0 {
			details := make([]vo.TranscodingProgressItem, 0, len(resources))
			for _, res := range resources {
				statusStr := global.ResourceStatusName(res.Status)
				progress := 0.0
				if res.Status == global.AUDIT_APPROVED {
					progress = 100.0
				}
				details = append(details, vo.TranscodingProgressItem{
					ResourceID:    res.ID,
					ResourceTitle: res.Title,
					Status:        statusStr,
					Progress:      progress,
				})
			}
			videos[i].TranscodingDetails = details
		}
	}

	return
}

// 获取处理中视频列表（后台管理）
func GetProcessingVideoList(videoListReq dto.VideoListReq) (total int64, videos []vo.VideoInfoManageResp) {
	processingStatuses := []int{global.CREATED_VIDEO, global.VIDEO_PROCESSING, global.SUBMIT_REVIEW}
	global.Mysql.Model(&model.Video{}).Where("status IN ? AND pgc_attached = ?", processingStatuses, false).Count(&total)
	global.Mysql.Model(&model.Video{}).Where("status IN ? AND pgc_attached = ?", processingStatuses, false).
		Order("created_at DESC").
		Limit(videoListReq.PageSize).Offset((videoListReq.Page - 1) * videoListReq.PageSize).Scan(&videos)

	transcodingVideoIDs := make([]uint, 0, len(videos))
	for i := 0; i < len(videos); i++ {
		videos[i].Clicks += GetVideoClicks(videos[i].ID)
		videos[i].Author = GetUserBaseInfo(videos[i].Uid)
		videos[i].Tags = LoadVideoTagNames(videos[i].ID)
		transcodingVideoIDs = append(transcodingVideoIDs, videos[i].ID)
	}

	progressMap := collectVideoProgressSnapshots(transcodingVideoIDs)
	for i := 0; i < len(videos); i++ {
		if snapshot, ok := progressMap[videos[i].ID]; ok {
			videos[i].TranscodingProgress = snapshot.Overall
			videos[i].TranscodingDetails = snapshot.Details
			videos[i].UploadProgress = snapshot.UploadInfo
		}
	}

	return
}

// 删除视频(后台管理)
func DeleteVideoManage(ctx *gin.Context, id uint) error {
	var video model.Video
	global.Mysql.Model(&model.Video{}).Where("id = ?", id).First(&video)
	if video.ID == 0 {
		return errors.New("视频不存在")
	}

	return deleteVideoAndRelatedData(id, video.Uid, &video)
}

// 获取待审核视频列表
func GetReviewList(reviewListReq dto.ReviewListReq) (total int64, videos []vo.ReviewListResp) {
	// 审核链路隔离：视频审核列表仅展示 UGC 待审视频，排除已绑定 PGC 的视频
	baseQuery := global.Mysql.Model(&model.Video{}).
		Where("status = ?", global.WAITING_REVIEW).
		Where("pgc_attached = ?", false)

	baseQuery.Count(&total)
	baseQuery.
		Order("created_at DESC").
		Limit(reviewListReq.PageSize).
		Offset((reviewListReq.Page - 1) * reviewListReq.PageSize).
		Scan(&videos)

	// 更新播放量和作者数据
	for i := 0; i < len(videos); i++ {
		videos[i].Author = GetUserBaseInfo(videos[i].Uid)
		videos[i].Tags = LoadVideoTagNames(videos[i].ID)
	}

	return
}

// 获取热门视频
func GetHotVideo(ctx *gin.Context, page, pageSize int) []vo.VideoResp {
	ids := cache.GetHotVideoId()
	videoIds := utils.SlicePagingStr(ids, page, pageSize)

	videos := make([]vo.VideoResp, 0, len(videoIds))
	for _, idStr := range videoIds {
		id := utils.StringToUint(idStr)
		if id == 0 {
			continue
		}
		v := GetVideoInfo(id)
		if v.ID == 0 {
			continue
		}
		v.Clicks += GetVideoClicks(id)
		v.DanmakuCount = GetDanmakuCount(id)
		videos = append(videos, v)
	}

	return videos
}

// 获取最近上传的视频
func GetLatestVideo(ctx *gin.Context, page, pageSize int) []vo.VideoResp {
	var videoIds []uint
	global.Mysql.Model(&model.Video{}).
		Where("status = ? AND pgc_attached = ?", global.AUDIT_APPROVED, false).
		Order("created_at DESC").
		Limit(pageSize).Offset((page-1)*pageSize).
		Pluck("id", &videoIds)

	videos := make([]vo.VideoResp, 0, len(videoIds))
	for _, id := range videoIds {
		v := GetVideoInfo(id)
		if v.ID == 0 {
			continue
		}
		v.Clicks += GetVideoClicks(id)
		v.DanmakuCount = GetDanmakuCount(id)
		videos = append(videos, v)
	}

	return videos
}

// 获取分区视频
func GetVideoListByPartition(ctx *gin.Context, size int, partitionId uint) []vo.VideoResp {
	videoIds := cache.GetVideoIdByPartition(partitionId, int64(size))

	videos := make([]vo.VideoResp, 0, len(videoIds))
	for _, idStr := range videoIds {
		id := utils.StringToUint(idStr)
		if id == 0 {
			continue
		}
		v := GetVideoInfo(id)
		if v.ID == 0 {
			continue
		}
		v.Clicks += GetVideoClicks(id)
		v.DanmakuCount = GetDanmakuCount(id)
		videos = append(videos, v)
	}

	return videos
}

// 获取相关推荐视频
func GetRelatedVideoList(ctx *gin.Context, videoId uint) []vo.VideoResp {
	video := GetVideoInfo(videoId)

	var videoIds []uint
	// 查询同作者的2个视频
	var authorVideoIds []uint
	global.Mysql.Model(&model.Video{}).
		Where("uid = ? and id != ? and `status` = ? and pgc_attached = ?", video.Uid, videoId, global.AUDIT_APPROVED, false).
		Limit(2).Pluck("id", &authorVideoIds)
	videoIds = append(videoIds, authorVideoIds...)

	// 查询同分区的7个视频（防止视频与同作者视频相同或者与当前视频相同）
	// 获取主分区ID用于推荐
	partitionId := video.PartitionId
	if parentId, exists := global.VideoPartitionMap[video.PartitionId]; exists && parentId != 0 {
		partitionId = parentId
	}

	for _, v := range cache.GetVideoIdByPartition(partitionId, 7) {
		id := utils.StringToUint(v)
		if id != videoId && !utils.IsUintInSlice(authorVideoIds, id) {
			videoIds = append(videoIds, id)
		}
	}

	videos := make([]vo.VideoResp, 0, len(videoIds))
	for _, id := range videoIds {
		v := GetVideoInfo(id)
		if v.ID == 0 {
			continue
		}
		v.Clicks += GetVideoClicks(id)
		v.DanmakuCount = GetDanmakuCount(id)
		videos = append(videos, v)
	}

	return videos
}

// 搜索视频
func SearchVideo(ctx *gin.Context, searchVideoReq dto.SearchVideoReq) []vo.VideoResp {
	var videoIds []uint

	page := searchVideoReq.Page
	if page < 1 {
		page = 1
	}
	ps := searchVideoReq.PageSize
	if ps < 1 {
		ps = 15
	}

	q := global.Mysql.Model(&model.Video{}).Where("`status` = ? AND pgc_attached = ?", global.AUDIT_APPROVED, false)

	tr := normalizeTimeRange(searchVideoReq.TimeRange)
	if start := parseTimeRangeStart(tr); start != nil {
		q = q.Where("created_at >= ?", *start)
	}

	if len(searchVideoReq.KeyWords) > 0 {
		// 直接用mysql模糊查询，之后可能会更换为es
		kw := "%" + escapeLikeKeyword(searchVideoReq.KeyWords) + "%"
		q = q.Where("(title LIKE ? ESCAPE '\\\\' OR tags LIKE ? ESCAPE '\\\\')", kw, kw)
	}

	// 排序：YouTube 风格（MVP）
	switch normalizeSort(searchVideoReq.Sort) {
	case "most_viewed":
		q = q.Order("clicks desc").Order("id desc")
	default: // newest / relevance
		q = q.Order("created_at desc").Order("id desc")
	}

	q.Limit(ps).Offset((page-1)*ps).Pluck("id", &videoIds)

	videos := make([]vo.VideoResp, 0, len(videoIds))
	for _, id := range videoIds {
		v := GetVideoInfo(id)
		v.Clicks += GetVideoClicks(id)
		v.DanmakuCount = GetDanmakuCount(id)
		videos = append(videos, v)
	}

	return videos
}

func CreateVideo(video *model.Video) (uint, error) {
	if video.ShortID == "" {
		sid, err := AllocateUniqueVideoShortID()
		if err != nil {
			return 0, err
		}
		video.ShortID = sid
	}

	if err := global.Mysql.Create(video).Error; err != nil {
		utils.ErrorLog("创建视频失败", "video", err.Error())
		return 0, errors.New("创建视频失败")
	}

	return video.ID, nil
}

// 重新转码视频（后台管理专用）
func ReTranscodeVideo(ctx *gin.Context, videoId uint) error {
	var video model.Video
	global.Mysql.Where("id = ?", videoId).First(&video)
	if video.ID == 0 {
		return errors.New("视频不存在")
	}

	// 记录重转码前的“原始审核状态”。
	// 兜底：若视频状态已被中间流程改成非通过，但仍存在已通过的有效分P，
	// 则按“原本已审核通过”处理，避免重转码结束后错误回落到待审核。
	originalVideoStatus := video.Status
	if originalVideoStatus != global.AUDIT_APPROVED {
		var approvedResourceCount int64
		global.Mysql.Model(&model.Resource{}).
			Where("vid = ? AND deleted_at IS NULL AND status = ?", videoId, global.AUDIT_APPROVED).
			Count(&approvedResourceCount)
		if approvedResourceCount > 0 {
			originalVideoStatus = global.AUDIT_APPROVED
			utils.InfoLog(
				fmt.Sprintf("【重转码状态兜底】VideoID=%d 当前status=%d，但存在%d个已通过分P，按AUDIT_APPROVED恢复",
					videoId, video.Status, approvedResourceCount),
				"video",
			)
		}
	}

	// 获取视频的所有资源（包含软删除的，因为之前失败的转码可能已经软删除了资源）
	var resources []model.Resource
	global.Mysql.Unscoped().Where("vid = ?", videoId).Find(&resources)
	if len(resources) == 0 {
		return errors.New("该视频没有可转码的资源")
	}
	utils.InfoLog(fmt.Sprintf("【重新转码】VideoID=%d, 找到%d个资源记录", videoId, len(resources)), "transcoding")

	// ========== 第1步：为每个分P查找原始文件 ==========
	// 注意：同一个 FileID 可能对应多个分P（秒传/复用场景）。
	// 重新转码必须按“资源(分P)”维度处理，否则会只转出一个分P。
	type resourceFileInfo struct {
		resource model.Resource
		vf       model.VideoFile
	}
	var rfInfos []resourceFileInfo

	for _, resource := range resources {
		// 只对未软删除的资源重新转码，避免“一个分P”被算成多个（历史软删除资源会被 Unscoped 查出）
		if resource.DeletedAt.Valid {
			continue
		}

		var vf model.VideoFile
		found := false

		// 方式1：通过 FileID 查找
		if resource.FileID != 0 {
			if err := global.Mysql.Where("id = ?", resource.FileID).First(&vf).Error; err == nil && vf.DirName != "" {
				found = true
				utils.InfoLog(fmt.Sprintf("【重新转码】通过FileID找到文件 resourceID=%d, FileID=%d, DirName=%s", resource.ID, resource.FileID, vf.DirName), "transcoding")
			} else {
				utils.ErrorLog("【重新转码】通过FileID查找失败", "transcoding", fmt.Sprintf("resourceID=%d, FileID=%d, err=%v", resource.ID, resource.FileID, err))
			}
		} else {
			utils.InfoLog(fmt.Sprintf("【重新转码】resourceID=%d 的 FileID=0，跳过方式1", resource.ID), "transcoding")
		}

		// 方式2：通过 VideoIndexFile.DirName 查找
		if !found {
			var indexFile model.VideoIndexFile
			if err := global.Mysql.Unscoped().Where("resource_id = ?", resource.ID).First(&indexFile).Error; err == nil && indexFile.DirName != "" {
				utils.InfoLog(fmt.Sprintf("【重新转码】找到VideoIndexFile resourceID=%d, DirName=%s", resource.ID, indexFile.DirName), "transcoding")
				if foundVF, err := findVideoFileByDirName(global.Mysql, indexFile.DirName); err == nil && foundVF != nil && foundVF.DirName != "" {
					vf = *foundVF
					found = true
					utils.InfoLog(fmt.Sprintf("【重新转码】通过DirName找到VideoFile resourceID=%d, DirName=%s", resource.ID, indexFile.DirName), "transcoding")
				} else {
					utils.ErrorLog("【重新转码】通过DirName查找VideoFile失败", "transcoding", fmt.Sprintf("resourceID=%d, DirName=%s, err=%v", resource.ID, indexFile.DirName, err))
				}
			} else {
				utils.ErrorLog("【重新转码】查找VideoIndexFile失败", "transcoding", fmt.Sprintf("resourceID=%d, err=%v", resource.ID, err))
			}
		}

		if !found {
			utils.ErrorLog("【重新转码】所有查找方式均失败，跳过此分P", "transcoding", fmt.Sprintf("resourceID=%d", resource.ID))
			continue
		}

		// 检查源文件是否存在（远程模式不需要，Worker 会从 OSS 下载）
		var inputPath string
		if global.Config.Transcoding.Mode != "remote" {
			suffix := utils.GetFileSuffix(vf.OriginalName)
			inputPath = "./upload/video/" + vf.DirName + "/upload" + suffix
			if _, err := os.Stat(inputPath); os.IsNotExist(err) {
				utils.ErrorLog("分P原始视频文件不存在，跳过", "transcoding", inputPath)
				continue
			}
		}

		rfInfos = append(rfInfos, resourceFileInfo{resource: resource, vf: vf})
	}

	if len(rfInfos) == 0 {
		return errors.New("所有分P的原始文件都不存在，无法重新转码")
	}

	// ========== 第2步：停止旧进程、清理旧数据 ==========
	if HasTranscodingProcess(videoId) {
		processCount := GetTranscodingProcessCount(videoId)
		utils.InfoLog(fmt.Sprintf("【重新转码】VideoID=%d 存在%d个转码进程，将停止后重新转码", videoId, processCount), "transcoding")
		StopTranscodingAndCleanup(videoId)
	}

	ResetGPUState()

	// 收集所有需要清理的目录（去重）
	dirsToClean := make(map[string]bool)
	for _, rf := range rfInfos {
		dirsToClean[rf.vf.DirName] = true
	}
	// 收集旧目录、删除旧 VideoIndexFile 与旧资源（Unscoped 确保能找到已软删除的记录）
	// 注意：这里必须硬删除旧 resource。用户侧数据（弹幕/历史/收藏/点赞）全部通过 short_id 或 vid 关联，
	// 新资源复用原 ShortID 即可保住绑定；按 resource.id 关联的 video_index_file / video_file_ref 本就由转码重建。
	// 如果软删除，旧行仍占用 UNIQUE(short_id) 索引，新资源 Create 时会报 Duplicate entry (Error 1062)。
	for _, resource := range resources {
		var indexFile model.VideoIndexFile
		if err := global.Mysql.Unscoped().Where("resource_id = ?", resource.ID).First(&indexFile).Error; err == nil && indexFile.DirName != "" {
			dirsToClean[indexFile.DirName] = true
		}
		global.Mysql.Unscoped().Where("resource_id = ?", resource.ID).Delete(&model.VideoIndexFile{})
		global.Mysql.Unscoped().Delete(&resource)
	}

	// 清理旧转码文件
	for d := range dirsToClean {
		deleteOldTranscodedFilesFromOSS(d)
		deleteOldTranscodedFilesLocal(d)
	}

	// 将稿件状态设为转码中，避免前端展示无资源的稿件
	global.Mysql.Model(&model.Video{}).Where("id = ?", videoId).Update("status", global.VIDEO_PROCESSING)
	cache.DelVideoInfo(videoId)

	// 保存原始审核状态到 Redis，供 cron 恢复时使用
	// 如果后续 Enqueue 失败，cron RecoverStuckTranscoding 的 reenqueueResource
	// 会读此值设置 OriginalVideoStatus，保证转码完成后恢复到正确状态。
	origStatusKey := fmt.Sprintf("transcoding:orig_status:%d", videoId)
	global.Redis.RawClient().Set(context.Background(), origStatusKey, originalVideoStatus, 72*time.Hour)

	// ========== 第3步：异步执行 ffprobe 探测、创建资源、启动转码 ==========
	// ffprobe 大文件耗时较长，放入 goroutine 避免 API 超时
	go func() {
		type transcodingTask struct {
			info *dto.TranscodingInfo
		}
		var tasks []transcodingTask

		for _, rf := range rfInfos {
			suffix := utils.GetFileSuffix(rf.vf.OriginalName)

			var info *dto.TranscodingInfo
			if global.Config.Transcoding.Mode == "remote" {
				// 远程模式：Worker 会从 OSS 下载并 ffprobe，此处用旧资源元数据
				info = &dto.TranscodingInfo{
					CodecName:           rf.resource.CodecName,
					Duration:            float64(rf.resource.Duration),
					OriginalVideoStatus: originalVideoStatus,
				}
			} else {
				// 本地模式：ffprobe 探测源文件获取元数据
				inputPath := "./upload/video/" + rf.vf.DirName + "/upload" + suffix
				var err error
				info, err = ProcessVideoInfo(inputPath)
				if err != nil {
					utils.ErrorLog("读取视频信息失败，跳过", "transcoding",
						fmt.Sprintf("resourceID=%d, file=%s, err=%v", rf.resource.ID, inputPath, err))
					continue
				}
			}

			// 复用原资源的 ShortID，确保历史记录/弹幕等绑定到 shortId 的数据不受影响
			shortID := rf.resource.ShortID
			if shortID == "" {
				shortID, _ = AllocateUniqueResourceShortID()
			}
			newResource := model.Resource{
				Vid:           videoId,
				Uid:           video.Uid,
				Title:         rf.resource.Title,
				CodecName:     info.CodecName,
				Status:        global.VIDEO_PROCESSING,
				Duration:      utils.SecFromFloat(info.Duration),
				FileID:        rf.vf.ID,
				ShortID:       shortID,
				SortOrder:     rf.resource.SortOrder,
				VisibleStatus: rf.resource.VisibleStatus, // 继承旧分P的可见性，避免重转码后变隐藏
			}
			if err := global.Mysql.Create(&newResource).Error; err != nil {
				utils.ErrorLog("创建新资源记录失败", "transcoding", err.Error())
				continue
			}

			// 重新转码直接输出到源目录，不创建新目录
			sourceDir := rf.vf.DirName

			transcodingInfo := &dto.TranscodingInfo{
				VideoID:             videoId,
				ResourceID:          newResource.ID,
				InputFile:           "./upload/video/" + sourceDir + "/upload" + suffix,
				OutputDir:           "./upload/video/" + sourceDir + "/",
				DirName:             sourceDir,
				Suffix:              suffix,
				Width:               info.Width,
				Height:              info.Height,
				Duration:            info.Duration,
				CodecName:           info.CodecName,
				FPS:                 info.FPS,
				FPS30:               info.FPS30,
				FPS60:               info.FPS60,
				VideoBitRate:        info.VideoBitRate,
				AudioBitRate:        info.AudioBitRate,
				AudioSampleRate:     info.AudioSampleRate,
				AudioChannels:       info.AudioChannels,
				OriginalVideoStatus: originalVideoStatus,
			}

			tasks = append(tasks, transcodingTask{info: transcodingInfo})
		}

		if len(tasks) == 0 {
			utils.ErrorLog("没有可转码的分P", "transcoding", fmt.Sprintf("videoId=%d", videoId))
			// 所有分P都探测失败，将稿件标记为处理失败
			global.Mysql.Model(&model.Video{}).Where("id = ?", videoId).Update("status", global.PROCESSING_FAIL)
			cache.DelVideoInfo(videoId)
			return
		}

		utils.InfoLog(fmt.Sprintf("【重新转码】VideoID=%d, 提交%d个分P转码任务", videoId, len(tasks)), "transcoding")

		// 通过 GetCurrentTranscoder 分发转码任务：
		//   mode=local  → LocalTranscoder.Enqueue (goroutine，进程内执行)
		//   mode=remote → RemoteTranscoder.Enqueue (推 Redis Stream，Worker 池消费)
		// 所有分P同时提交，local 模式下 semaphore 限制并发，remote 模式下 Worker 池调度。
		for _, task := range tasks {
			if err := GetCurrentTranscoder().Enqueue(context.Background(), task.info); err != nil {
				utils.ErrorLog("【重新转码】提交任务失败", "transcoding",
					fmt.Sprintf("ResourceID=%d, err=%v", task.info.ResourceID, err))
			}
		}
		utils.InfoLog(fmt.Sprintf("【重新转码提交完成】VideoID=%d, 共%d个分P", videoId, len(tasks)), "transcoding")
	}()

	return nil
}

// deleteOldTranscodedFilesFromOSS 删除OSS上指定目录的旧转码文件
func deleteOldTranscodedFilesFromOSS(dirName string) {
	if global.Config.Storage.OssType == "local" {
		return
	}

	if dirName == "" {
		return
	}

	localDir := "./upload/video/" + dirName + "/"
	files, err := os.ReadDir(localDir)
	if err != nil {
		utils.InfoLog(fmt.Sprintf("【OSS清理】读取目录失败: %s, 跳过OSS清理", localDir), "transcoding")
		return
	}

	objectKeys := make(map[string]bool)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		// 跳过原始上传文件（upload.mkv, upload.mp4 等）
		if strings.HasPrefix(fileName, "upload.") {
			continue
		}
		objectKeys[fileName] = true
	}

	if len(objectKeys) == 0 {
		utils.InfoLog(fmt.Sprintf("【OSS清理】目录 %s 没有需要清理的文件", dirName), "transcoding")
		return
	}

	deletedCount := 0
	for fileName := range objectKeys {
		objectKey := "video/" + dirName + "/" + fileName
		if err := global.Storage.DeleteObject(objectKey); err != nil {
			utils.ErrorLog(fmt.Sprintf("【OSS清理】删除文件失败: %s", objectKey), "transcoding", err.Error())
		} else {
			deletedCount++
		}
		// 同时删除备用OSS
		if global.StorageBackup != nil {
			if err := global.StorageBackup.DeleteObject(objectKey); err != nil {
				utils.ErrorLog(fmt.Sprintf("【备用OSS清理】删除文件失败: %s", objectKey), "transcoding", err.Error())
			}
		}
	}

	utils.InfoLog(fmt.Sprintf("【OSS清理】目录 %s 清理完成, 删除 %d/%d 个文件", dirName, deletedCount, len(objectKeys)), "transcoding")
}

// deleteOldTranscodedFilesLocal 删除本地旧转码文件
func deleteOldTranscodedFilesLocal(dirName string) {
	if dirName == "" {
		return
	}

	localDir := "./upload/video/" + dirName + "/"
	files, err := os.ReadDir(localDir)
	if err != nil {
		utils.InfoLog(fmt.Sprintf("【本地清理】读取目录失败: %s, 跳过本地清理", localDir), "transcoding")
		return
	}

	deletedCount := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		// 跳过原始上传文件和封面
		if strings.HasPrefix(fileName, "upload.") || fileName == "cover.jpg" {
			continue
		}
		filePath := localDir + fileName
		if err := os.Remove(filePath); err != nil {
			utils.ErrorLog(fmt.Sprintf("【本地清理】删除文件失败: %s", filePath), "transcoding", err.Error())
		} else {
			deletedCount++
		}
	}

	if deletedCount > 0 {
		utils.InfoLog(fmt.Sprintf("【本地清理】目录 %s 清理完成, 删除 %d 个文件", dirName, deletedCount), "transcoding")
	}
}

// ReUploadVideo 重新上传转码产物到 OSS。
// 当资源状态为 UPLOAD_FAILED（转码成功但上传失败）时调用，
// 重跑 uploadToOSS，成功后走 completeTransaction 流转状态。
func ReUploadVideo(ctx *gin.Context, videoId uint) error {
	// 查找视频
	var video model.Video
	global.Mysql.Where("id = ?", videoId).First(&video)
	if video.ID == 0 {
		return errors.New("视频不存在")
	}

	// 查找该视频下所有上传失败的资源
	var failedResources []model.Resource
	global.Mysql.Where("vid = ? AND status = ?", videoId, global.UPLOAD_FAILED).Find(&failedResources)
	if len(failedResources) == 0 {
		return errors.New("该视频没有上传失败的资源")
	}

	// 同步执行重传（调用方已开 goroutine 则无需再开）
	s := GetTranscoder()
	for _, resource := range failedResources {
		utils.InfoLog(fmt.Sprintf("【重新上传OSS】开始 ResourceID=%d, VideoID=%d", resource.ID, videoId), "video")

		// 通过 FileID 查找 VideoFile 获取 DirName 和 suffix
		var vf model.VideoFile
		if err := global.Mysql.Where("id = ?", resource.FileID).First(&vf).Error; err != nil {
			utils.ErrorLog("【重新上传OSS】查找视频文件失败", "video", fmt.Sprintf("resourceID=%d, err=%v", resource.ID, err))
			continue
		}
		if vf.DirName == "" {
			utils.ErrorLog("【重新上传OSS】VideoFile.DirName 为空", "video", fmt.Sprintf("resourceID=%d", resource.ID))
			continue
		}

		suffix := utils.GetFileSuffix(vf.OriginalName)
		sourceDir := vf.DirName
		outputDir := "./upload/video/" + sourceDir + "/"

		// 检查目录是否存在
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			utils.ErrorLog("【重新上传OSS】输出目录不存在", "video", outputDir)
			continue
		}

		// 将资源状态设为转码中（上传中）
		global.Mysql.Model(&model.Resource{}).Where("id = ?", resource.ID).Update("status", global.VIDEO_PROCESSING)

		info := &dto.TranscodingInfo{
			VideoID:    videoId,
			ResourceID: resource.ID,
			DirName:    sourceDir,
			OutputDir:  outputDir,
			Suffix:     suffix,
		}

		if err := s.uploadToOSS(ctx, info); err != nil {
			utils.ErrorLog("【重新上传OSS】上传失败", "video", fmt.Sprintf("resourceID=%d, err=%v", resource.ID, err))
			s.markUploadFailed(ctx, info)
			continue
		}

		// 上传成功 → 调用 completeTransaction 流转状态
		if err := s.completeTransaction(ctx, info, global.WAITING_REVIEW); err != nil {
			utils.ErrorLog("【重新上传OSS】状态流转失败", "video", fmt.Sprintf("resourceID=%d, err=%v", resource.ID, err))
		} else {
			utils.InfoLog(fmt.Sprintf("【重新上传OSS】成功 ResourceID=%d", resource.ID), "video")
			// 异步上传到备用 OSS（多源容灾，不影响主流程）
			if global.StorageBackup != nil {
				go s.uploadToBackup(info)
			}
		}
	}

	return nil
}

// 通过视频ID查询视频
func FindVideoById(id uint) (video model.Video, err error) {
	err = global.Mysql.Where("`id` = ?", id).First(&video).Error
	return
}

// 获取视频信息
func GetVideoInfo(videoId uint) (video vo.VideoResp) {
	video = cache.GetVideoInfo(videoId)
	if video.ID == 0 {
		// 缓存不存在，从数据库加载并写入缓存
		video = VideoWriteCache(videoId)
	}
	// 直接使用缓存中的完整数据，不再重新查询Resources和Author
	// 如果需要最新数据，应该先调用cache.DelVideoInfo删除缓存，再调用本函数

	return
}

// 视频信息写入缓存
func VideoWriteCache(videoId uint) (video vo.VideoResp) {
	global.Mysql.Model(&model.Video{}).Select(vo.VIDEO_FIELD).
		Where("id = ? and status = ?", videoId, global.AUDIT_APPROVED).Scan(&video)
	if video.ID == 0 {
		utils.ErrorLog("视频信息不存在", "video", fmt.Sprintf("videoId=%d", videoId))
		return
	}

	video.Author = GetAuthorPublic(video.Uid)
	video.Tags = LoadVideoTagNames(videoId)
	// 获取视频资源
	// 只返回对外可见的分P（VisibleStatus=1），处理中的新分P对观众隐藏
	video.Resources = GetVisibleResources(videoId)
	// 确保Resources不是nil，如果是nil则初始化为空切片，避免JSON序列化为null
	if video.Resources == nil {
		video.Resources = []vo.ResourceResp{}
	}

	// 存到redis
	cache.SetVideoInfo(video)

	return
}

// 获取视频状态
func getVideoStatus(videoId uint) int {
	var processingCount int64
	var failedCount int64
	var totalCount int64

	// 获取该视频下的所有资源总数
	global.Mysql.Model(&model.Resource{}).Where("vid = ?", videoId).Count(&totalCount)

	if totalCount == 0 {
		return global.VIDEO_PROCESSING
	}

	// 统计转码中和失败的数量
	global.Mysql.Model(&model.Resource{}).Where("vid = ? and status = ?", videoId, global.VIDEO_PROCESSING).Count(&processingCount)
	global.Mysql.Model(&model.Resource{}).Where("vid = ? and status = ?", videoId, global.PROCESSING_FAIL).Count(&failedCount)

	// 修复：逻辑优先级调整
	// 1. 只要有一个还在转码，整体就是转码中
	if processingCount > 0 {
		return global.VIDEO_PROCESSING
	}

	// 2. 如果没有转码中的，但有失败的，且失败+成功的总数等于总数（此处简化为失败>0）
	// 如果你希望只要有一个失败整体就显失败，用这个：
	if failedCount > 0 {
		return global.PROCESSING_FAIL
	}

	// 3. 全部资源都转码完成且无失败，进入待审核状态
	return global.WAITING_REVIEW
}

// ReTranscodeResource 重新入队单个分P（资源级别）的转码任务。
// quality 非空时只重试指定画质，为空则重试所有非 success 画质。
func ReTranscodeResource(ctx *gin.Context, resourceID uint, quality string) (*model.Resource, error) {
	var resource model.Resource
	if err := global.Mysql.Where("id = ?", resourceID).First(&resource).Error; err != nil {
		return nil, errors.New("分P不存在")
	}

	// 只允许重试失败或卡住的分P
	if resource.Status != global.PROCESSING_FAIL && resource.Status != global.VIDEO_PROCESSING {
		return nil, errors.New("该分P当前状态不允许重新转码，仅失败或卡住的分P可重试")
	}

	// 查找 VideoFile
	var vf model.VideoFile
	if resource.FileID != 0 {
		if err := global.Mysql.Where("id = ?", resource.FileID).First(&vf).Error; err != nil || vf.DirName == "" {
			// 通过 FileID 找不到，尝试从 VideoIndexFile 恢复
			var indexFile model.VideoIndexFile
			if err := global.Mysql.Unscoped().Where("resource_id = ?", resource.ID).First(&indexFile).Error; err == nil && indexFile.DirName != "" {
				if foundVF, err := findVideoFileByDirName(global.Mysql, indexFile.DirName); err == nil && foundVF != nil && foundVF.DirName != "" {
					vf = *foundVF
				}
			}
		}
	}
	if vf.DirName == "" {
		return nil, errors.New("找不到该分P的原始视频文件记录")
	}

	suffix := utils.GetFileSuffix(vf.OriginalName)

	// 远程模式不需要本地文件存在，Worker 会从 OSS 下载
	if global.Config.Transcoding.Mode != "remote" {
		inputPath := "./upload/video/" + vf.DirName + "/upload" + suffix
		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("源文件不存在: %s", inputPath)
		}
	}

	// 获取原始审核状态（优先从 Redis 读取，兜底用视频当前状态）
	origStatus := getOriginalVideoStatus(resource.Vid)
	if origStatus < 0 {
		var currentStatus int
		global.Mysql.Model(&model.Video{}).Select("status").Where("id = ?", resource.Vid).Scan(&currentStatus)
		origStatus = currentStatus
	}

	// 重置资源状态为转码中
	global.Mysql.Model(&resource).Update("status", global.VIDEO_PROCESSING)
	resource.Status = global.VIDEO_PROCESSING

	// 清理 Redis 中失败画质的进度，保留已完成的
	rdb := global.Redis.RawClient()
	statusKey := fmt.Sprintf("transcoding:status:%d", resource.ID)
	if quality != "" {
		// 只重置指定画质
		pipe := rdb.Pipeline()
		pipe.HSet(context.Background(), statusKey, fmt.Sprintf("status_%s", quality), "waiting")
		pipe.HSet(context.Background(), statusKey, fmt.Sprintf("progress_%s", quality), "0.00")
		pipe.HSet(context.Background(), statusKey, "status", "processing")
		pipe.HSet(context.Background(), statusKey, "updated", fmt.Sprintf("%d", time.Now().Unix()))
		pipe.Exec(context.Background())
	} else if hash, err := rdb.HGetAll(context.Background(), statusKey).Result(); err == nil && len(hash) > 0 {
		// 重置所有非 success 画质
		pipe := rdb.Pipeline()
		for field, val := range hash {
			if strings.HasPrefix(field, "status_") && val != "success" {
				pipe.HSet(context.Background(), statusKey, field, "waiting")
				progressField := "progress_" + field[len("status_"):]
				pipe.HSet(context.Background(), statusKey, progressField, "0.00")
			}
		}
		pipe.HSet(context.Background(), statusKey, "status", "processing")
		pipe.HSet(context.Background(), statusKey, "updated", fmt.Sprintf("%d", time.Now().Unix()))
		pipe.Exec(context.Background())
	}

	// 去重检查：指定画质已经是 waiting 说明已有重试任务在队列，跳过
	if quality != "" {
		curStatus, err := rdb.HGet(context.Background(), statusKey, fmt.Sprintf("status_%s", quality)).Result()
		if err == nil && curStatus == "waiting" {
			utils.InfoLog(fmt.Sprintf("【重试画质-跳过】ResourceID=%d, quality=%s 已在队列中", resource.ID, quality), "transcoding")
			return &resource, nil
		}
	}

	// 从 VideoIndexFile 读取分辨率（如果有的话），避免 Worker 用 0x0 计算 targets 导致 MinInt64
	var indexFile model.VideoIndexFile
	width, height := 0, 0
	if err := global.Mysql.Where("resource_id = ?", resource.ID).First(&indexFile).Error; err == nil {
		width, height = indexFile.Width, indexFile.Height
	}

	// 构建 TranscodingInfo 并入队
	info := &dto.TranscodingInfo{
		VideoID:             resource.Vid,
		ResourceID:          resource.ID,
		InputFile:           "./upload/video/" + vf.DirName + "/upload" + suffix,
		OutputDir:           "./upload/video/" + vf.DirName + "/",
		DirName:             vf.DirName,
		Suffix:              suffix,
		CodecName:           resource.CodecName,
		Duration:            float64(resource.Duration),
		Width:               width,
		Height:              height,
		OriginalVideoStatus: origStatus,
	}

	// 先发 cancel 信号杀掉旧的 pollProgress 协程，避免与新任务竞争
	if err := GetCurrentTranscoder().Cancel(context.Background(), resource.Vid); err != nil {
		utils.InfoLog(fmt.Sprintf("【重转码】Cancel 旧 pollProgress 失败（可忽略）: %v", err), "transcoding")
	}

	if err := GetCurrentTranscoder().Enqueue(context.Background(), info); err != nil {
		// 入队失败，回滚资源状态
		global.Mysql.Model(&resource).Update("status", global.PROCESSING_FAIL)
		return nil, fmt.Errorf("入队失败: %w", err)
	}

	utils.InfoLog(fmt.Sprintf("【重新转码-单分P】ResourceID=%d, VideoID=%d, DirName=%s",
		resource.ID, resource.Vid, vf.DirName), "transcoding")
	return &resource, nil
}

// getOriginalVideoStatus 从 Redis 读取原始审核状态。
// 如果键不存在或已过期，返回 -1。
func getOriginalVideoStatus(videoID uint) int {
	rdb := global.Redis.RawClient()
	ctx := context.Background()
	key := fmt.Sprintf("transcoding:orig_status:%d", videoID)
	val, err := rdb.Get(ctx, key).Int()
	if err != nil {
		return -1
	}
	return val
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
