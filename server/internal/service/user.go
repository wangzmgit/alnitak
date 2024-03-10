package service

import (
	"errors"
	"math"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
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

	if cache.GetEmailCode(registerReq.Email) != registerReq.Code { // 验证邮箱验证码
		return errors.New("邮箱验证码错误")
	}

	// 删除邮箱验证码
	cache.DelEmailCode(registerReq.Email)

	// 对密码加密
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	// 随机生成不重复的用户名
	username, err := generateUniqueUsername()
	if err != nil {
		return errors.New("用户名生成失败")
	}

	// 保存到数据库
	if err := global.Mysql.Create(&model.User{
		Username: username,
		Email:    registerReq.Email,
		Password: string(hashedPassword),
	}).Error; err != nil {
		utils.ErrorLog("创建用户失败", "user", err.Error())
		return errors.New("创建用户失败")
	}

	return nil
}

func UserLogin(ctx *gin.Context, loginReq dto.LoginReq) (accessToken, refreshToken string, err error) {
	// 读取数据库
	user, err := FindUserByEmail(loginReq.Email)
	if err != nil {
		return "", "", errors.New("获取用户信息失败")
	}

	// 验证账号密码
	passwordError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if passwordError != nil {
		// 记录登录尝试次数
		cache.IncrLoginTryCount(loginReq.Email)
		return "", "", errors.New("用户名密码不匹配")
	}

	// 生成验证token
	if accessToken, err = jwt.GenerateAccessToken(user.ID); err != nil {
		return "", "", errors.New("验证token生成失败")
	}
	// 生成刷新token
	if refreshToken, err = jwt.GenerateRefreshToken(user.ID); err != nil {
		return "", "", errors.New("刷新token生成失败")
	}

	// 用户ID写入Cookie
	SetUserIdCookie(ctx, user.ID)

	// 存入缓存
	cache.SetRefreshToken(user.ID, refreshToken)

	return accessToken, refreshToken, nil
}

func UpdateToken(ctx *gin.Context, tokenReq dto.TokenReq) (accessToken, refreshToken string, err error) {
	// 验证并解析token
	_, claims, err := jwt.ParseToken(tokenReq.RefreshToken)
	if err != nil {
		utils.ErrorLog("token验证失败", "user", err.Error())
		return "", "", errors.New("token验证失败")
	}

	// 读取缓存
	if !cache.IsRefreshTokenExist(claims.UserId, tokenReq.RefreshToken) { // refreshToken 存在
		utils.ErrorLog("无效Token", "user", "token不存在")
		return "", "", errors.New("无效Token")
	}

	// 移除refreshToken
	cache.DelRefreshToken(claims.UserId, tokenReq.RefreshToken)

	// 重新生成refreshToken
	refreshToken, err = jwt.GenerateRefreshToken(claims.UserId)
	if err != nil {
		utils.ErrorLog("Token生成失败", "user", err.Error())
		return "", "", errors.New("Token生成失败")
	}

	// 刷新accessToken
	accessToken, err = jwt.GenerateAccessToken(claims.UserId)
	if err != nil {
		utils.ErrorLog("AccessToken生成失败", "user", err.Error())
		return "", "", errors.New("Token生成失败")
	}

	// 用户ID写入Cookie
	SetUserIdCookie(ctx, claims.UserId)

	// 存入缓存
	cache.SetRefreshToken(claims.UserId, refreshToken)

	return accessToken, refreshToken, nil
}

func Logout() {
	// 删除token
	// 移除用户cookie

}

// 生成用户Id和md5并写入Cookie
func SetUserIdCookie(ctx *gin.Context, userId uint) {
	salt := viper.GetString("security.user_id_salt")
	ckMd5 := utils.GenerateSaltedMD5(strconv.Itoa(int(userId)), salt)

	ctx.SetCookie("user_id", strconv.Itoa(int(userId)), math.MaxInt32, "/", "", false, true)
	ctx.SetCookie("user_id_md5", ckMd5, math.MaxInt32, "/", "", false, true)
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

	// 过滤掉敏感信息
	user.Email = ""
	user.Phone = ""

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

// 随机生成一个不重复的用户名
func generateUniqueUsername() (string, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return "", err
	}

	id := node.Generate()

	// 前缀 + snowflake ID(36进制)
	return viper.GetString("user.prefix") + strconv.FormatInt(id.Int64(), 36), nil
}
