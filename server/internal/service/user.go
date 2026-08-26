package service

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/pkg/jwt"
	"interastral-peace.com/alnitak/utils"
)

func UserRegister(ctx *gin.Context, registerReq dto.RegisterReq) error {
	if user, _ := FindUserByEmail(registerReq.Email); user.ID != 0 {
		return errors.New("邮箱已存在")
	}

	// 验证码暴力枚举防护
	if cache.GetEmailCodeTryCount(registerReq.Email) > 5 {
		return errors.New("验证码错误次数过多，请重新获取")
	}

	if cache.GetEmailCode(registerReq.Email) != registerReq.Code { // 验证邮箱验证码
		cache.IncrEmailCodeTryCount(registerReq.Email)
		return errors.New("邮箱验证码错误")
	}

	// 删除邮箱验证码及计数
	cache.DelEmailCode(registerReq.Email)
	cache.DelEmailCodeTryCount(registerReq.Email)

	// 对密码加密
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)

	// 保存到数据库
	if err := global.Mysql.Create(&model.User{
		Username: generateUniqueUsername(), // 随机生成不重复的用户名
		Email:    registerReq.Email,
		Password: string(hashedPassword),
	}).Error; err != nil {
		utils.ErrorLog("创建用户失败", "user", err.Error())
		return errors.New("创建用户失败")
	}

	return nil
}

// generateTokenPair 生成 accessToken + refreshToken 并存入缓存
func generateTokenPair(userId uint, rememberMe ...bool) (accessToken, refreshToken string, err error) {
	if accessToken, err = jwt.GenerateAccessToken(userId); err != nil {
		return "", "", errors.New("验证token生成失败")
	}

	useLong := len(rememberMe) > 0 && rememberMe[0]
	var refreshTTL time.Duration
	if useLong {
		refreshTTL = cache.REFRESH_TOKEN_LONG_EXPIRATION_TIME
	} else {
		refreshTTL = cache.REFRESH_TOKEN_EXPIRATION_TIME
	}
	if refreshToken, err = jwt.GenerateRefreshToken(userId, refreshTTL); err != nil {
		return "", "", errors.New("刷新token生成失败")
	}
	cache.SetRefreshToken(userId, refreshToken, refreshTTL)
	return accessToken, refreshToken, nil
}

func UserLogin(ctx *gin.Context, loginReq dto.LoginReq) (accessToken, refreshToken string, userId uint, err error) {
	user, err := FindUserByEmail(loginReq.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cache.IncrLoginTryCount(loginReq.Email)
			return "", "", 0, errors.New("用户名密码不匹配")
		}
		return "", "", 0, errors.New("获取用户信息失败")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)) != nil {
		cache.IncrLoginTryCount(loginReq.Email)
		return "", "", 0, errors.New("用户名密码不匹配")
	}

	accessToken, refreshToken, err = generateTokenPair(user.ID, loginReq.RememberMe)
	if err != nil {
		return "", "", 0, err
	}
	return accessToken, refreshToken, user.ID, nil
}

func EmailLogin(ctx *gin.Context, loginReq dto.EmailLoginReq) (accessToken, refreshToken string, userId uint, err error) {
	user, err := FindUserByEmail(loginReq.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cache.IncrLoginTryCount(loginReq.Email)
			return "", "", 0, errors.New("用户名密码不匹配")
		}
		return "", "", 0, errors.New("获取用户信息失败")
	}

	if cache.GetEmailCode(loginReq.Email) != loginReq.Code {
		cache.IncrLoginTryCount(loginReq.Email)
		return "", "", 0, errors.New("邮箱验证错误")
	}

	accessToken, refreshToken, err = generateTokenPair(user.ID, loginReq.RememberMe)
	if err != nil {
		return "", "", 0, err
	}
	return accessToken, refreshToken, user.ID, nil
}

func UpdateToken(ctx *gin.Context, tokenReq dto.TokenReq) (accessToken, refreshToken string, userId uint, err error) {
	// 验证并解析token
	_, claims, err := jwt.ParseToken(tokenReq.RefreshToken)
	if err != nil {
		utils.ErrorLog("token验证失败", "user", err.Error())
		return "", "", 0, errors.New("token验证失败")
	}

	// 必须为 refreshToken（TokenType=1），拒绝 accessToken 被用于刷新
	if claims.TokenType != 1 {
		return "", "", 0, errors.New("token类型错误")
	}

	// 读取缓存
	if !cache.IsRefreshTokenExist(claims.UserId, tokenReq.RefreshToken) { // refreshToken 存在
		utils.ErrorLog("无效Token", "user", "token不存在")
		return "", "", 0, errors.New("无效Token")
	}

	// refreshToken过期时间小于缓冲时间
	if claims.ExpiresAt.Before(time.Now().Add(cache.REFRESH_TOKEN_BUFFER_TIME)) {
		// 移除refreshToken
		cache.DelRefreshToken(claims.UserId, tokenReq.RefreshToken)
		// 重新生成refreshToken（活跃用户用长过期时间）
		refreshToken, err = jwt.GenerateRefreshToken(claims.UserId, cache.REFRESH_TOKEN_LONG_EXPIRATION_TIME)
		if err != nil {
			utils.ErrorLog("Token生成失败", "user", err.Error())
			return "", "", 0, errors.New("Token生成失败")
		}
	}

	// 刷新accessToken
	accessToken, err = jwt.GenerateAccessToken(claims.UserId)
	if err != nil {
		utils.ErrorLog("AccessToken生成失败", "user", err.Error())
		return "", "", 0, errors.New("Token生成失败")
	}

	// 存入缓存（新 token 用长过期时间）
	if refreshToken != "" {
		cache.SetRefreshToken(claims.UserId, refreshToken, cache.REFRESH_TOKEN_LONG_EXPIRATION_TIME)
	}

	return accessToken, refreshToken, claims.UserId, nil
}

// Logout 吊销 refresh：阶段 B 下登出接口不再依赖 Auth 中间件注入 userId，
// 须从 refresh JWT 解析出用户 id，与行业标准「仅凭 Cookie 即可退出」一致。
func Logout(_ *gin.Context, tokenReq dto.TokenReq) {
	if tokenReq.RefreshToken == "" {
		return
	}
	_, claims, err := jwt.ParseToken(tokenReq.RefreshToken)
	if err != nil {
		return
	}
	cache.DelRefreshToken(claims.UserId, tokenReq.RefreshToken)
}

// 编辑用户信息
func EditUserInfo(ctx *gin.Context, editUserInfoReq dto.EditUserInfoReq) error {
	userId := ctx.GetUint("userId")

	// 检查用户名是否重复
	if u, _ := FindUserByName(editUserInfoReq.Name); u.ID != 0 && u.ID != userId {
		return errors.New("用户名已存在")
	}

	// 检测文件是否有效
	current, _ := FindUserById(userId)
	if editUserInfoReq.Avatar != current.Avatar && cache.GetUploadImage(editUserInfoReq.Avatar) != userId {
		editUserInfoReq.Avatar = current.Avatar
	}

	if editUserInfoReq.SpaceCover != current.SpaceCover && cache.GetUploadImage(editUserInfoReq.SpaceCover) != userId {
		editUserInfoReq.SpaceCover = current.SpaceCover
	}

	birthday, err := time.Parse("2006-01-02", editUserInfoReq.Birthday)
	if err != nil {
		utils.ErrorLog("日期转换失败", "user", err.Error())
		birthday = time.Unix(0, 0)
	}
	if err := global.Mysql.Model(&model.User{}).Where("id = ?", userId).Updates(
		map[string]interface{}{
			"avatar":      editUserInfoReq.Avatar,
			"username":    editUserInfoReq.Name,
			"gender":      editUserInfoReq.Gender,
			"birthday":    birthday,
			"space_cover": editUserInfoReq.SpaceCover,
			"sign":        editUserInfoReq.Sign,
		},
	).Error; err != nil {
		return err
	}

	//移除缓存
	cache.DelUserInfo(userId)

	return nil
}

func ResetPwdCheck(ctx *gin.Context, email string) error {
	user, _ := FindUserByEmail(email)
	if user.ID == 0 {
		return errors.New("用户不存在")
	}

	// 更新验证状态
	cache.SetResetPwdCheckStatus(email, 1)

	return nil
}

func ModifyPwd(ctx *gin.Context, modifyPwdReq dto.ModifyPwdReq) error {
	// 校验：必须先通过 ResetPwdCheck（人机验证）才能改密码
	if cache.GetResetPwdCheckStatus(modifyPwdReq.Email) != 1 {
		return errors.New("请先完成安全验证")
	}

	// 校验：邮箱验证码
	if cache.GetEmailCode(modifyPwdReq.Email) != modifyPwdReq.Code {
		// 验证码错误计数
		cache.IncrEmailCodeTryCount(modifyPwdReq.Email)
		if cache.GetEmailCodeTryCount(modifyPwdReq.Email) > 5 {
			cache.DelEmailCode(modifyPwdReq.Email)            // 使当前验证码失效
			cache.DelResetPwdCheckStatus(modifyPwdReq.Email) // 需要重新走安全验证
			return errors.New("验证码错误次数过多，请重新获取")
		}
		return errors.New("邮箱验证错误")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(modifyPwdReq.Password), bcrypt.DefaultCost)
	user, err := FindUserByEmail(modifyPwdReq.Email)
	if err != nil {
		return errors.New("用户不存在")
	}
	if err := global.Mysql.Model(&model.User{}).Where("id = ?", user.ID).
		Update("password", string(hashedPassword)).Error; err != nil {
		utils.ErrorLog("修改密码失败", "user", err.Error())
		return errors.New("修改失败")
	}

	// 清理验证状态
	cache.DelEmailCode(modifyPwdReq.Email)
	cache.DelEmailCodeTryCount(modifyPwdReq.Email)
	cache.DelResetPwdCheckStatus(modifyPwdReq.Email)
	// 改密码后失效该用户的所有已有 refreshToken，防止已被窃取的会话续签
	cache.DelAllRefreshToken(user.ID)

	return nil
}

// ChangePassword 修改密码（已登录用户，知道旧密码）
func ChangePassword(ctx *gin.Context, userId uint, req dto.ChangePasswordReq) error {
	if len(req.NewPassword) < 6 {
		return errors.New("新密码长度不能小于6位")
	}

	user, err := FindUserById(userId)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 校验旧密码
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)) != nil {
		return errors.New("旧密码错误")
	}

	// 新密码不能和旧密码相同
	if req.OldPassword == req.NewPassword {
		return errors.New("新密码不能与旧密码相同")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err := global.Mysql.Model(&model.User{}).Where("id = ?", user.ID).
		Update("password", string(hashedPassword)).Error; err != nil {
		utils.ErrorLog("修改密码失败", "user", err.Error())
		return errors.New("修改失败")
	}

	// 改密码后失效所有 refreshToken，防止已被窃取的会话续签
	cache.DelAllRefreshToken(user.ID)

	return nil
}

// 获取用户列表
func GetUserListManage(userListReq dto.UserListReq) (total int64, users []vo.UserInfoManageResp) {
	global.Mysql.Model(&model.User{}).Count(&total)
	global.Mysql.Model(&model.User{}).Limit(userListReq.PageSize).
		Offset((userListReq.Page - 1) * userListReq.PageSize).Scan(&users)
	return
}

// 编辑用户信息(后台管理)
func EditUserInfoManage(ctx *gin.Context, editUserInfoReq dto.EditUserInfoManageReq) error {
	// 检查用户名是否重复
	if u, _ := FindUserByName(editUserInfoReq.Name); u.ID != 0 && u.ID != editUserInfoReq.Uid {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否重复
	if u, _ := FindUserByEmail(editUserInfoReq.Email); u.ID != 0 && u.ID != editUserInfoReq.Uid {
		return errors.New("邮箱已存在")
	}

	if err := global.Mysql.Model(&model.User{}).Where("id = ?", editUserInfoReq.Uid).Updates(
		map[string]interface{}{
			"avatar":      editUserInfoReq.Avatar,
			"space_cover": editUserInfoReq.SpaceCover,
			"username":    editUserInfoReq.Name,
			"email":       editUserInfoReq.Email,
			"sign":        editUserInfoReq.Sign,
		},
	).Error; err != nil {
		return err
	}

	//移除缓存
	cache.DelUserInfo(editUserInfoReq.Uid)

	return nil
}

// 编辑角色首页
func EditUserRole(ctx *gin.Context, editUserRoleReq dto.EditUserRoleReq) error {
	if err := global.Mysql.Model(&model.User{}).Where("id = ?", editUserRoleReq.Uid).Updates(
		map[string]interface{}{
			"role": editUserRoleReq.Code,
		},
	).Error; err != nil {
		utils.ErrorLog("修改用户角色失败", "user", err.Error())
		return errors.New("修改失败")
	}
	return nil
}

// 封禁用户
func BanUser(ctx *gin.Context, banUserReq dto.BanUserReq) error {
	var err error
	var endTime time.Time
	if banUserReq.EndTime != "" {
		endTime, err = time.Parse("2006-01-02", banUserReq.EndTime)
		if err != nil {
			utils.ErrorLog("时间转换失败", "user", err.Error())
			return errors.New("时间转换失败")
		}
	}

	// 判断封禁目标是否是自己
	if banUserReq.Uid == ctx.GetUint("userId") {
		return errors.New("不能封禁自己")
	}

	tx := global.Mysql.Begin()

	// 将用户进行中的封禁记录更新为已取消
	if err := tx.Model(&model.UserBan{}).Where("uid = ? and status = 0", banUserReq.Uid).Updates(
		map[string]interface{}{
			"status": 4,
		},
	).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("封禁用户失败", "user", err.Error())
		return errors.New("封禁失败")
	}

	// 更新用户状态
	if err := tx.Model(&model.User{}).Where("id = ?", banUserReq.Uid).Updates(
		map[string]interface{}{
			"status": 1,
		},
	).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("封禁用户失败", "user", err.Error())
		return errors.New("封禁失败")
	}

	// 添加封禁记录
	if err := tx.Create(&model.UserBan{
		EndTime:  endTime,
		Reason:   banUserReq.Reason,
		Uid:      banUserReq.Uid,
		Operator: ctx.GetUint("userId"),
	}).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("封禁用户失败", "user", err.Error())
		return errors.New("封禁失败")
	}

	tx.Commit()

	//移除缓存
	cache.DelUserInfo(banUserReq.Uid)

	return nil
}

// 解封用户
func UnBanUser(ctx *gin.Context, id uint) error {
	// 查询封禁记录
	var banUser model.UserBan
	if err := global.Mysql.Where("id = ?", id).First(&banUser).Error; err != nil {
		utils.ErrorLog("解封用户失败", "user", err.Error())
		return errors.New("解封失败")
	}

	tx := global.Mysql.Begin()

	// 更新用户状态
	if err := tx.Model(&model.User{}).Where("id = ?", banUser.Uid).Updates(
		map[string]interface{}{
			"status": 0,
		},
	).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("解封用户失败", "user", err.Error())
		return errors.New("解封失败")
	}

	// 更新封禁记录
	if err := tx.Model(&model.UserBan{}).Where("id = ?", id).Updates(
		map[string]interface{}{
			"status": 1,
		},
	).Error; err != nil {
		tx.Rollback()
		utils.ErrorLog("解封用户失败", "user", err.Error())
		return errors.New("解封失败")
	}

	tx.Commit()

	//移除缓存
	cache.DelUserInfo(banUser.Uid)

	return nil
}

// 获取封禁记录
func GetUserBanRecord(userId uint) (bans []vo.UserBanRecordResp) {
	global.Mysql.Model(&model.UserBan{}).Where("uid = ?", userId).Order("id desc").Limit(3).Scan(&bans)
	return
}

// 删除用户
func DeleteUser(ctx *gin.Context, id uint) error {
	// 删除用户的所有视频（包括视频的所有关联数据）
	var videos []model.Video
	global.Mysql.Where("uid = ?", id).Find(&videos)
	for _, video := range videos {
		// 删除视频的所有关联资源
		global.Mysql.Where("vid = ?", video.ID).Delete(&model.Resource{})
		global.Mysql.Where("vid = ?", video.ID).Delete(&model.VideoIndexFile{})
		// 评论级联（含 comment_like / msg_like type=3 / msg_like type=0 / msg_reply）
		CleanupContentComments(video.ID, global.CONTENT_TYPE_VIDEO)
		if video.ShortID != "" {
			global.Mysql.Where("video_short_id = ?", video.ShortID).Delete(&model.Danmaku{})
		}
		global.Mysql.Where("vid = ?", video.ID).Delete(&model.CollectVideo{})
		global.Mysql.Where("vid = ?", video.ID).Delete(&model.LikeVideo{})
		global.Mysql.Where("vid = ?", video.ID).Delete(&model.History{})
		global.Mysql.Where("cid = ? and `type` = ?", video.ID, global.CONTENT_TYPE_VIDEO).Delete(&model.AtMessage{})
		// 删除视频本身
		global.Mysql.Where("id = ?", video.ID).Delete(&model.Video{})
	}

	// 删除用户的所有文章（包括文章的所有关联数据）
	var articles []model.Article
	global.Mysql.Where("uid = ?", id).Find(&articles)
	for _, article := range articles {
		// 评论级联清理
		CleanupContentComments(article.ID, global.CONTENT_TYPE_ARTICLE)
		global.Mysql.Where("cid = ? and `type` = ?", article.ID, global.CONTENT_TYPE_ARTICLE).Delete(&model.AtMessage{})
		global.Mysql.Where("aid = ?", article.ID).Delete(&model.CollectArticle{})
		global.Mysql.Where("aid = ?", article.ID).Delete(&model.LikeArticle{})
		// 删除文章本身
		global.Mysql.Where("id = ?", article.ID).Delete(&model.Article{})
	}

	// 删除用户的所有评论
	global.Mysql.Where("uid = ?", id).Delete(&model.Comment{})

	// 删除用户的所有收藏夹及其内容
	var collections []model.Collection
	global.Mysql.Where("uid = ?", id).Find(&collections)
	for _, collection := range collections {
		global.Mysql.Where("collection_id = ?", collection.ID).Delete(&model.CollectVideo{})
		global.Mysql.Where("id = ?", collection.ID).Delete(&model.Collection{})
	}

	// 删除用户的点赞记录
	global.Mysql.Where("uid = ?", id).Delete(&model.LikeVideo{})
	global.Mysql.Where("uid = ?", id).Delete(&model.LikeArticle{})

	// 删除用户的收藏记录
	global.Mysql.Where("uid = ?", id).Delete(&model.CollectVideo{})
	global.Mysql.Where("uid = ?", id).Delete(&model.CollectArticle{})

	// 删除用户的历史记录
	global.Mysql.Where("uid = ?", id).Delete(&model.History{})

	// 删除用户的关注关系
	global.Mysql.Where("uid = ? or fid = ?", id, id).Delete(&model.Relation{})

	// 删除用户的消息记录
	global.Mysql.Where("uid = ? or sid = ?", id, id).Delete(&model.AtMessage{})
	global.Mysql.Where("uid = ? or sid = ?", id, id).Delete(&model.ReplyMessage{})
	global.Mysql.Where("uid = ? or sid = ?", id, id).Delete(&model.LikeMessage{})
	global.Mysql.Where("uid = ? or sid = ?", id, id).Delete(&model.Whisper{})

	// 删除用户的封禁记录
	global.Mysql.Where("uid = ?", id).Delete(&model.UserBan{})

	// 删除用户的视频文件记录
	global.Mysql.Where("uid = ?", id).Delete(&model.VideoFile{})

	if err := global.Mysql.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		utils.ErrorLog("删除用户失败", "user", err.Error())
		return errors.New("删除失败")
	}

	//移除缓存
	cache.DelUserInfo(id)

	return nil
}

// 获取用户信息
func GetUserInfo(userId uint) (user vo.UserInfoResp) {
	var err error

	user = cache.GetUserInfo(userId)
	if user.ID == 0 {
		user, err = FindUserInfoById(userId)
		if err != nil {
			utils.ErrorLog("用户信息查询失败", "user", err.Error())
			return
		}

		// 存到redis
		cache.SetUserInfo(user)
	}

	return
}

// 获取用户基本信息
func GetUserBaseInfo(userId uint) (user vo.UserInfoResp) {
	var err error

	user = cache.GetUserInfo(userId)
	if user.ID == 0 {
		user, err = FindUserInfoById(userId)
		if err != nil {
			utils.ErrorLog("用户信息查询失败", "user", err.Error())
			return
		}

		// 存到redis
		cache.SetUserInfo(user)
	}

	// 查询粉丝数量（实时查询，不缓存）
	global.Mysql.Model(&model.Relation{}).
		Where("target_uid = ? and (relation = ? or relation = ?)", userId, global.FOLLOWED, global.MUTUAL_FANS).
		Count(&user.Fans)

	return
}

// GetAuthorPublic 视频区等场景使用的作者公开信息（仅 uid、昵称、头像）。
func GetAuthorPublic(userId uint) (a vo.AuthorPublicResp) {
	u := GetUserBaseInfo(userId)
	if u.ID == 0 {
		return
	}
	a.UID = u.ID
	a.Name = u.Username
	a.Avatar = u.Avatar
	return
}

// 通过用户ID查询用户
func FindUserById(id uint) (user model.User, err error) {
	err = global.Mysql.Where("`id` = ?", id).First(&user).Error
	return
}

// 通过用户ID查询用户
func FindUserInfoById(id uint) (user vo.UserInfoResp, err error) {
	err = global.Mysql.Model(&model.User{}).Where("`id` = ?", id).Scan(&user).Error

	return
}

// 通过用户邮箱查询用户
func FindUserByEmail(email string) (user model.User, err error) {
	err = global.Mysql.Where("`email` = ?", email).First(&user).Error
	return
}

// 通过用户名查询用户
func FindUserByName(name string) (user model.User, err error) {
	err = global.Mysql.Where("`username` = ?", name).First(&user).Error
	return
}

// 通过用户邮箱或用户名查询用户
func FindUserByEmailOrUsername(username string) (user model.User, err error) {
	err = global.Mysql.Where("`email` = ? or `username` = ?", username, username).First(&user).Error
	return
}

// 通过用户名查询多个用户
func FindUserIdsByName(names []string) (ids []uint) {
	global.Mysql.Model(model.User{}).Where("username in (?)", names).Pluck("id", &ids)
	return
}

// 获取用户的最新封禁信息
func FindUserLastBan(userId uint) (ban vo.BanResp) {
	global.Mysql.Model(&model.UserBan{}).Where("uid = ?", userId).Order("id desc").First(&ban)
	return
}

// SearchUser 按用户名、签名模糊搜索正常状态用户（status=0），返回公开字段（经 GetUserBaseInfo）
func SearchUser(ctx *gin.Context, req dto.SearchKeywordPageReq) []vo.UserInfoResp {
	var userIds []uint
	page := req.Page
	if page < 1 {
		page = 1
	}
	ps := req.PageSize
	if ps < 1 {
		ps = 15
	}
	kw := strings.TrimSpace(req.KeyWords)
	if utf8.RuneCountInString(kw) > 100 {
		kw = string([]rune(kw)[:100])
	}

	q := global.Mysql.Model(&model.User{}).Where("status = ?", 0)

	tr := normalizeTimeRange(req.TimeRange)
	if start := parseTimeRangeStart(tr); start != nil {
		q = q.Where("created_at >= ?", *start)
	}

	if len(kw) > 0 {
		pattern := "%" + escapeLikeKeyword(kw) + "%"
		q = q.Where("(username LIKE ? ESCAPE '\\\\' OR sign LIKE ? ESCAPE '\\\\')", pattern, pattern)
	}

	// 用户排序：most_fans 后续需要稳定 fans_count；MVP 默认 newest
	switch normalizeSort(req.Sort) {
	default:
		q = q.Order("created_at desc").Order("id desc")
	}

	q.Limit(ps).Offset((page - 1) * ps).Pluck("id", &userIds)

	users := make([]vo.UserInfoResp, 0, len(userIds))
	for _, id := range userIds {
		u := GetUserBaseInfo(id)
		if u.ID == 0 {
			continue
		}
		users = append(users, u)
	}
	return users
}

// 随机生成一个不重复的用户名
func generateUniqueUsername() string {
	id := global.SnowflakeNode.Generate()

	// 前缀 + snowflake ID(36进制)
	return global.Config.User.Prefix + strconv.FormatInt(id.Int64(), 36)
}
