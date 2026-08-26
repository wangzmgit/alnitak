package service

import (
	"errors"
	"net"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/pkg/playtoken"
)

// AudioTrackInfo 音轨信息，用于返回给前端
type AudioTrackInfo struct {
	Language string `json:"language"`
	Title    string `json:"title"`
	IsDefault bool  `json:"isDefault"`
}

// PlayURLsResult 播放 URL 结果，包含主备音视频直链。
type PlayURLsResult struct {
	VideoURL    string `json:"video"`
	AudioURL    string `json:"audio"`
	Expires     int64  `json:"expires"`
	BackupVideo string `json:"backupVideo,omitempty"` // 备用 OSS 视频 URL（B站风格多源容灾）
	BackupAudio string `json:"backupAudio,omitempty"` // 备用 OSS 音频 URL
	AudioTracks []AudioTrackInfo `json:"audioTracks,omitempty"` // 可用音轨列表
}

const playGrantTTL = 8 * time.Minute
const streamSliceTTL = 2 * time.Hour

// IssuePlayGrantForResource 为分 P shortId 签发播放授权（公开稿件可无登录 uid=0）。
func IssuePlayGrantForResource(ctx *gin.Context, resourceShortID string) (token string, expiresUnix int64, err error) {
	var res model.Resource
	if err := global.Mysql.Where("short_id = ?", resourceShortID).First(&res).Error; err != nil || res.ID == 0 {
		return "", 0, errors.New("资源不存在")
	}
	var v model.Video
	if err := global.Mysql.Where("id = ?", res.Vid).First(&v).Error; err != nil || v.ID == 0 {
		return "", 0, errors.New("视频不存在")
	}
	if v.Status != global.AUDIT_APPROVED || res.Status != global.AUDIT_APPROVED {
		return "", 0, errors.New("内容不可播放")
	}
	uid := ctx.GetUint("userId")
	tok, err := playtoken.IssuePlayGrant(uid, v.ID, resourceShortID, playGrantTTL)
	if err != nil {
		return "", 0, err
	}
	return tok, time.Now().Add(playGrantTTL).Unix(), nil
}

// GetPlayURLs 校验 grant 后返回指定清晰度下音视频分离的直链（带 st）。
// 配置了备用 OSS 时同时返回 backupVideo/backupAudio（B站风格多源容灾）。
// audioLang 可选：指定音轨语言，为空时使用默认音轨。
func GetPlayURLs(ctx *gin.Context, resourceShortID, grantToken, quality, audioLang string) (result PlayURLsResult, err error) {
	if err := checkPlayAccess(ctx); err != nil {
		return result, err
	}
	claims, err := playtoken.ParsePlayGrant(grantToken)
	if err != nil {
		return result, errors.New("播放凭证无效")
	}
	if claims.ResourceShortID != resourceShortID {
		return result, errors.New("播放凭证与资源不匹配")
	}

	var res model.Resource
	if err := global.Mysql.Where("short_id = ?", resourceShortID).First(&res).Error; err != nil || res.ID == 0 {
		return result, errors.New("资源不存在")
	}
	if res.Vid != claims.VideoID {
		return result, errors.New("播放凭证已失效")
	}
	var v model.Video
	if err := global.Mysql.Where("id = ?", res.Vid).First(&v).Error; err != nil || v.ID == 0 {
		return result, errors.New("视频不存在")
	}
	if v.Status != global.AUDIT_APPROVED || res.Status != global.AUDIT_APPROVED {
		return result, errors.New("内容不可播放")
	}

	var files []model.VideoIndexFile
	if err := global.Mysql.Where("resource_id = ?", res.ID).Find(&files).Error; err != nil || len(files) == 0 {
		return result, errors.New("播放索引未就绪")
	}
	var candidates []model.VideoIndexFile
	for _, f := range files {
		if f.IsSegmentBase() {
			candidates = append(candidates, f)
		}
	}
	if len(candidates) == 0 {
		return result, errors.New("当前资源不支持该播放方式")
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].VideoBandwidth > candidates[j].VideoBandwidth
	})
	var file model.VideoIndexFile
	if quality != "" {
		found := false
		for _, f := range candidates {
			if f.Quality == quality {
				file = f
				found = true
				break
			}
		}
		if !found {
			return result, errors.New("清晰度不存在")
		}
	} else {
		file = candidates[0]
	}

	exp := time.Now().Add(streamSliceTTL).Unix()
	vTok, err := playtoken.IssueStreamToken(file.DirName, file.VideoFile, streamSliceTTL)
	if err != nil {
		return result, err
	}
	base := publicAPIBase(ctx)

	// ===== 多音轨支持 =====
	audioFile := file.AudioFile
	var audioTracks []AudioTrackInfo

	// 查询 AudioTrack 记录（如果存在）
	var trackRecords []model.AudioTrack
	global.Mysql.Where("resource_id = ?", res.ID).Find(&trackRecords)

	if len(trackRecords) > 0 {
		// 填充音轨信息到返回值
		for _, tr := range trackRecords {
			audioTracks = append(audioTracks, AudioTrackInfo{
				Language:  tr.Language,
				Title:     tr.Title,
				IsDefault: tr.IsDefault,
			})
		}
		result.AudioTracks = audioTracks

		// 根据 audioLang 参数选择音轨文件
		if audioLang != "" {
			found := false
			for _, tr := range trackRecords {
				if tr.Language == audioLang {
					audioFile = tr.AudioFile
					found = true
					break
				}
			}
			if !found {
				// audioLang 不存在时回退默认音轨
				for _, tr := range trackRecords {
					if tr.IsDefault {
						audioFile = tr.AudioFile
						break
					}
				}
			}
		} else {
			// 未指定时使用默认音轨
			for _, tr := range trackRecords {
				if tr.IsDefault {
					audioFile = tr.AudioFile
					break
				}
			}
		}
	}

	aTok, err := playtoken.IssueStreamToken(file.DirName, audioFile, streamSliceTTL)
	if err != nil {
		return result, err
	}

	dir := file.DirName
	result.VideoURL = base + "/api/v1/video/stream/" + url.PathEscape(file.VideoFile) + "?st=" + url.QueryEscape(vTok)
	result.AudioURL = base + "/api/v1/video/stream/" + url.PathEscape(audioFile) + "?st=" + url.QueryEscape(aTok)
	result.Expires = exp

	// B站风格：配置了备用 OSS 时附带 backup URL
	if global.StorageBackup != nil {
		if bv := global.GetBackupOssUrl("video/" + dir + "/" + file.VideoFile); bv != "" {
			result.BackupVideo = bv
		}
		if ba := global.GetBackupOssUrl("video/" + dir + "/" + audioFile); ba != "" {
			result.BackupAudio = ba
		}
	}

	return result, nil
}

// GetAudioTracks 查询指定资源的可用音轨列表
func GetAudioTracks(resourceID uint) ([]AudioTrackInfo, error) {
	var trackRecords []model.AudioTrack
	if err := global.Mysql.Where("resource_id = ?", resourceID).Find(&trackRecords).Error; err != nil {
		return nil, err
	}
	tracks := make([]AudioTrackInfo, 0, len(trackRecords))
	for _, tr := range trackRecords {
		tracks = append(tracks, AudioTrackInfo{
			Language:  tr.Language,
			Title:     tr.Title,
			IsDefault: tr.IsDefault,
		})
	}
	return tracks, nil
}

func publicAPIBase(ctx *gin.Context) string {
	if d := strings.TrimSpace(global.Config.Storage.Domain); d != "" {
		if strings.HasPrefix(d, "http://") || strings.HasPrefix(d, "https://") {
			return strings.TrimRight(d, "/")
		}
		scheme := "https://"
		if !global.Config.Storage.UseSSL {
			scheme = "http://"
		}
		return scheme + strings.TrimRight(d, "/")
	}
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + ctx.Request.Host
}

func checkPlayAccess(ctx *gin.Context) error {
	hosts := global.Config.Security.PlayAllowedRefererHosts
	if len(hosts) > 0 {
		ref := ctx.GetHeader("Referer")
		u, err := url.Parse(ref)
		if err != nil || u.Host == "" {
			return errors.New("拒绝访问")
		}
		h := strings.ToLower(u.Host)
		ok := false
		for _, allow := range hosts {
			allow = strings.ToLower(strings.TrimSpace(allow))
			if allow != "" && (h == allow || strings.HasSuffix(h, "."+allow)) {
				ok = true
				break
			}
		}
		if !ok {
			return errors.New("拒绝访问")
		}
	}
	cidrs := global.Config.Security.PlayAllowedCIDRs
	if len(cidrs) == 0 {
		return nil
	}
	ip := net.ParseIP(ctx.ClientIP())
	if ip == nil {
		return errors.New("拒绝访问")
	}
	for _, c := range cidrs {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		_, n, err := net.ParseCIDR(c)
		if err != nil {
			continue
		}
		if n.Contains(ip) {
			return nil
		}
	}
	return errors.New("拒绝访问")
}
