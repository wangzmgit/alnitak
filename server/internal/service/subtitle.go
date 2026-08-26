package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/pkg/subtitle"
	"interastral-peace.com/alnitak/utils"
)

// subtitlePublicURL 与图片等媒体一致：local 返回同源路由；OSS 返回可直访的存储 URL（由浏览器请求，需在 OSS 配 CORS，且播放器 video 启用 crossOrigin）
func subtitlePublicURL(objectKey string) string {
	base := filepath.Base(objectKey)
	if base == "" || base == "." || base == "/" {
		return ""
	}
	if global.Config.Storage.OssType == "local" {
		return "/api/subtitle/" + base
	}
	return global.GetOssUrl(objectKey)
}

const subtitleMaxBytes = 2 * 1024 * 1024 // 2MB

// ErrSubtitleDuplicateLang 同分 P 同语言唯一约束冲突
var ErrSubtitleDuplicateLang = errors.New("该分 P 已存在相同语言的字幕")

var (
	subtitleServeFileRegexp = regexp.MustCompile(`^[0-9]+\.vtt$`)
	subtitleLangRegexp      = regexp.MustCompile(`^[a-zA-Z]{2,3}(-[a-zA-Z0-9]+)*$`)
)

func normalizeSubtitleLang(s string) string {
	return strings.TrimSpace(s)
}

func validateSubtitleLang(lang string) error {
	if lang == "" || len(lang) > 20 {
		return errors.New("语言代码无效")
	}
	if !subtitleLangRegexp.MatchString(lang) {
		return errors.New("语言代码格式无效（请使用如 zh-Hans、en）")
	}
	return nil
}

func loadResourceWithVideoByRaw(resourceRaw string) (res model.Resource, v model.Video, err error) {
	rid, perr := ParseResourceID(resourceRaw)
	if perr != nil || rid == 0 {
		return res, v, errors.New("资源不存在")
	}
	if err := global.Mysql.Where("id = ?", rid).First(&res).Error; err != nil || res.ID == 0 {
		return res, v, errors.New("资源不存在")
	}
	if err := global.Mysql.Where("id = ?", res.Vid).First(&v).Error; err != nil || v.ID == 0 {
		return res, v, errors.New("视频不存在")
	}
	return res, v, nil
}

func subtitleVisibleToViewer(ctx *gin.Context, v model.Video, res model.Resource) bool {
	uid := ctx.GetUint("userId")
	if uid != 0 && v.Uid == uid {
		return true
	}
	return v.Status == global.AUDIT_APPROVED && res.Status == global.AUDIT_APPROVED
}

func assertSubtitleManager(ctx *gin.Context, v model.Video) error {
	uid := ctx.GetUint("userId")
	if uid == 0 {
		return errors.New("请先登录")
	}
	if v.Uid != uid {
		return errors.New("无权限操作该稿件字幕")
	}
	return nil
}

func deleteSubtitleObject(objectKey, localPath string) error {
	if localPath != "" {
		_ = os.Remove(localPath)
	}
	if global.Config.Storage.OssType != "local" && objectKey != "" {
		if err := global.Storage.DeleteObject(objectKey); err != nil {
			utils.ErrorLog("删除字幕对象失败", "subtitle", err.Error())
			return err
		}
		// 同时删除备用OSS
		if global.StorageBackup != nil {
			if err := global.StorageBackup.DeleteObject(objectKey); err != nil {
				utils.ErrorLog("删除备用OSS字幕对象失败", "subtitle", err.Error())
			}
		}
	}
	return nil
}

func unsetOtherSubtitleDefaults(tx *gorm.DB, resourceShortID string, keepID uint) error {
	return tx.Model(&model.SubtitleTrack{}).
		Where("resource_short_id = ? AND id != ? AND is_default = ?", resourceShortID, keepID, true).
		Update("is_default", false).Error
}

func parseBodyToVTT(body []byte, ext string) ([]byte, int, error) {
	var vtt []byte
	var sourceFmt int
	var err error
	switch strings.ToLower(ext) {
	case ".srt":
		vtt, err = subtitle.SRTToWebVTT(body)
		sourceFmt = model.SubtitleFormatSRT
	case ".vtt":
		vtt, err = subtitle.ValidateOrNormalizeWebVTT(body)
		sourceFmt = model.SubtitleFormatVTT
	default:
		return nil, 0, errors.New("仅支持 .srt 或 .vtt")
	}
	if err != nil {
		return nil, 0, err
	}
	return vtt, sourceFmt, nil
}

func writeSubtitleVTTFile(objectKey string, vtt []byte) (localPath string, err error) {
	localPath = filepath.Join(".", "upload", filepath.FromSlash(objectKey))
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		utils.ErrorLog("创建字幕目录失败", "subtitle", err.Error())
		return "", errors.New("保存字幕失败")
	}
	if err := os.WriteFile(localPath, vtt, 0644); err != nil {
		utils.ErrorLog("写入字幕失败", "subtitle", err.Error())
		return "", errors.New("保存字幕失败")
	}
	if global.Config.Storage.OssType != "local" {
		if err := global.Storage.PutObjectFromFile(objectKey, localPath); err != nil {
			_ = os.Remove(localPath)
			utils.ErrorLog("字幕上传到对象存储失败", "subtitle", err.Error())
			return "", errors.New("上传字幕失败")
		}
		// 上传到备用 OSS（带重试 + 失败持久化）
		go UploadToBackupWithRetry(objectKey, localPath, "subtitle")
	}
	return localPath, nil
}

// ListSubtitleTracks 字幕列表（公开：仅已发布；UP 主：含未发布分 P）
func ListSubtitleTracks(ctx *gin.Context, resourceShortID string) ([]vo.SubtitleTrackItem, error) {
	res, v, err := loadResourceWithVideoByRaw(resourceShortID)
	if err != nil {
		return nil, err
	}
	if !subtitleVisibleToViewer(ctx, v, res) {
		return nil, errors.New("内容不可访问")
	}
	var tracks []model.SubtitleTrack
	if err := global.Mysql.Where("resource_short_id = ? AND status = ?", res.ShortID, model.SubtitleStatusActive).
		Order("is_default DESC, id ASC").Find(&tracks).Error; err != nil {
		utils.ErrorLog("查询字幕失败", "subtitle", err.Error())
		return nil, errors.New("获取字幕失败")
	}
	out := make([]vo.SubtitleTrackItem, 0, len(tracks))
	for _, t := range tracks {
		out = append(out, vo.SubtitleTrackItem{
			ID:        t.ID,
			ShortID:   t.ShortID,
			Lang:      t.Lang,
			Label:     t.Label,
			URL:       subtitlePublicURL(t.ObjectKey),
			IsDefault: t.IsDefault,
		})
	}
	return out, nil
}

// ParseSubtitleTrackID 兼容对外 shortId 与数字自增 id。
func ParseSubtitleTrackID(raw string) (uint, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, errors.New("字幕ID不能为空")
	}
	if utils.IsAllASCIIDigits(raw) {
		id := utils.StringToUint(raw)
		if id == 0 {
			return 0, errors.New("字幕不存在")
		}
		return id, nil
	}
	var track model.SubtitleTrack
	if err := global.Mysql.Where("short_id = ?", raw).First(&track).Error; err != nil || track.ID == 0 {
		return 0, errors.New("字幕不存在")
	}
	return track.ID, nil
}

// CreateSubtitleTrack 上传新建字幕轨
func CreateSubtitleTrack(ctx *gin.Context, resourceShortID, lang, label string, isDefault bool, fileExt string, r io.Reader) error {
	res, v, err := loadResourceWithVideoByRaw(resourceShortID)
	if err != nil {
		return err
	}
	if err := assertSubtitleManager(ctx, v); err != nil {
		return err
	}
	lang = normalizeSubtitleLang(lang)
	if err := validateSubtitleLang(lang); err != nil {
		return err
	}
	label = strings.TrimSpace(label)
	if len(label) > 64 {
		return errors.New("字幕显示名过长")
	}

	body, err := io.ReadAll(io.LimitReader(r, subtitleMaxBytes+1))
	if err != nil {
		return errors.New("读取文件失败")
	}
	if len(body) > subtitleMaxBytes {
		return errors.New("字幕文件过大")
	}

	vtt, sourceFmt, err := parseBodyToVTT(body, fileExt)
	if err != nil {
		return err
	}

	fileName := global.SnowflakeNode.Generate().String() + ".vtt"
	objectKey := "subtitle/" + fileName
	localPath, err := writeSubtitleVTTFile(objectKey, vtt)
	if err != nil {
		return err
	}

	uid := ctx.GetUint("userId")
	
	sid, err := AllocateUniqueSubtitleShortID()
	if err != nil {
		return errors.New("分配字幕ID失败")
	}
	
	track := model.SubtitleTrack{
		ShortID:         sid,
		ResourceShortID: res.ShortID,
		Vid:             v.ID,
		Lang:            lang,
		Label:           label,
		SourceFmt:       sourceFmt,
		ObjectKey:       objectKey,
		IsDefault:       isDefault,
		Origin:          model.SubtitleOriginUser,
		Status:          model.SubtitleStatusActive,
		CreatedBy:       uid,
	}

	err = global.Mysql.Transaction(func(tx *gorm.DB) error {
		if isDefault {
			if err := tx.Model(&model.SubtitleTrack{}).
				Where("resource_short_id = ?", res.ShortID).
				Update("is_default", false).Error; err != nil {
				return err
			}
		}
		if err := tx.Create(&track).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate") {
				utils.ErrorLog("字幕 Duplicate 键冲突", "subtitle", fmt.Sprintf("short_id=%s lang=%s err=%v", res.ShortID, lang, err))
				return ErrSubtitleDuplicateLang
			}
			return err
		}
		return nil
	})
	if err != nil {
		_ = deleteSubtitleObject(objectKey, localPath)
		if errors.Is(err, ErrSubtitleDuplicateLang) {
			return err
		}
		utils.ErrorLog("写入字幕记录失败", "subtitle", fmt.Sprintf("%v", err))
		return errors.New("保存字幕失败")
	}
	return nil
}

// UpdateSubtitleTrackReq 更新参数
type UpdateSubtitleTrackReq struct {
	Label      string
	IsDefault  *bool
	ReplaceVTT io.Reader
	SrcExt     string
}

// UpdateSubtitleTrack 更新元数据或替换文件
func UpdateSubtitleTrack(ctx *gin.Context, trackID uint, req UpdateSubtitleTrackReq) error {
	if trackID == 0 {
		return errors.New("字幕不存在")
	}
	var track model.SubtitleTrack
	if err := global.Mysql.Where("id = ?", trackID).First(&track).Error; err != nil || track.ID == 0 {
		return errors.New("字幕不存在")
	}
	var v model.Video
	if err := global.Mysql.Where("id = ?", track.Vid).First(&v).Error; err != nil || v.ID == 0 {
		return errors.New("视频不存在")
	}
	if err := assertSubtitleManager(ctx, v); err != nil {
		return err
	}

	var err error
	updates := map[string]interface{}{}
	trimmedLabel := strings.TrimSpace(req.Label)
	if trimmedLabel != track.Label {
		if len(trimmedLabel) > 64 {
			return errors.New("字幕显示名过长")
		}
		updates["label"] = trimmedLabel
	}
	if req.IsDefault != nil {
		updates["is_default"] = *req.IsDefault
	}

	oldKey := track.ObjectKey
	oldLocal := filepath.Join(".", "upload", filepath.FromSlash(oldKey))

	var newKey string
	var newLocal string

	if req.ReplaceVTT != nil {
		body, err := io.ReadAll(io.LimitReader(req.ReplaceVTT, subtitleMaxBytes+1))
		if err != nil {
			return errors.New("读取文件失败")
		}
		if len(body) > subtitleMaxBytes {
			return errors.New("字幕文件过大")
		}
		vtt, sourceFmt, err := parseBodyToVTT(body, req.SrcExt)
		if err != nil {
			return err
		}
		newName := global.SnowflakeNode.Generate().String() + ".vtt"
		newKey = "subtitle/" + newName
		newLocal, err = writeSubtitleVTTFile(newKey, vtt)
		if err != nil {
			return err
		}
		updates["object_key"] = newKey
		updates["source_fmt"] = sourceFmt
	}

	if len(updates) == 0 {
		if newKey != "" {
			_ = deleteSubtitleObject(newKey, newLocal)
		}
		return nil
	}

	err = global.Mysql.Transaction(func(tx *gorm.DB) error {
		if req.IsDefault != nil && *req.IsDefault {
			if err := unsetOtherSubtitleDefaults(tx, track.ResourceShortID, track.ID); err != nil {
				return err
			}
		}
		return tx.Model(&model.SubtitleTrack{}).Where("id = ?", track.ID).Updates(updates).Error
	})
	if err != nil {
		if newKey != "" {
			_ = deleteSubtitleObject(newKey, newLocal)
		}
		utils.ErrorLog("更新字幕失败", "subtitle", err.Error())
		return errors.New("更新字幕失败")
	}

	if newKey != "" {
		_ = deleteSubtitleObject(oldKey, oldLocal)
	}
	return nil
}

// DeleteSubtitleTrack 删除字幕轨
func DeleteSubtitleTrack(ctx *gin.Context, trackID uint) error {
	if trackID == 0 {
		return errors.New("字幕不存在")
	}
	var track model.SubtitleTrack
	if err := global.Mysql.Where("id = ?", trackID).First(&track).Error; err != nil || track.ID == 0 {
		return errors.New("字幕不存在")
	}
	var v model.Video
	if err := global.Mysql.Where("id = ?", track.Vid).First(&v).Error; err != nil || v.ID == 0 {
		return errors.New("视频不存在")
	}
	if err := assertSubtitleManager(ctx, v); err != nil {
		return err
	}

	localPath := filepath.Join(".", "upload", filepath.FromSlash(track.ObjectKey))
	// 须物理删除：软删后唯一索引 uk_subtitle_resource_lang 仍会占用，同语言无法再上传
	if err := global.Mysql.Unscoped().Delete(&model.SubtitleTrack{}, trackID).Error; err != nil {
		utils.ErrorLog("删除字幕记录失败", "subtitle", err.Error())
		return errors.New("删除字幕失败")
	}
	_ = deleteSubtitleObject(track.ObjectKey, localPath)
	return nil
}

// GetSubtitleTrackForFileServe 校验静态字幕访问并返回本地路径与 objectKey
func GetSubtitleTrackForFileServe(ctx *gin.Context, fileName string) (localPath, objectKey string, ok bool) {
	if !subtitleServeFileRegexp.MatchString(fileName) {
		return "", "", false
	}
	objectKey = "subtitle/" + fileName
	var tr model.SubtitleTrack
	if err := global.Mysql.Where("object_key = ? AND status = ?", objectKey, model.SubtitleStatusActive).First(&tr).Error; err != nil || tr.ID == 0 {
		return "", "", false
	}
	var res model.Resource
	if err := global.Mysql.Where("short_id = ?", tr.ResourceShortID).First(&res).Error; err != nil || res.ID == 0 {
		return "", "", false
	}
	var v model.Video
	if err := global.Mysql.Where("id = ?", res.Vid).First(&v).Error; err != nil || v.ID == 0 {
		return "", "", false
	}
	if !subtitleVisibleToViewer(ctx, v, res) {
		return "", "", false
	}
	localPath = filepath.Join(".", "upload", filepath.FromSlash(objectKey))
	return localPath, objectKey, true
}
