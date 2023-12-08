package service

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/pkg/jwt"
	jwt_parse "interastral-peace.com/alnitak/pkg/jwt"
	"interastral-peace.com/alnitak/utils"
)

func UserRegister(ctx *gin.Context, registerReq dto.RegisterReq) error {
	if user, _ := FindUserByEmail(registerReq.Email); user.ID != 0 {
		AddFailOperation(ctx, "注册", "邮箱已存在", gin.H{"邮箱": registerReq.Email}, nil)
		return errors.New("邮箱已存在")
	}

	if cache.GetEmailCode(registerReq.Email) != registerReq.Code { // 验证邮箱验证码
		AddFailOperation(ctx, "注册", "邮箱验证码错误", gin.H{"邮箱": registerReq.Email, "验证码": registerReq.Code}, nil)
		return errors.New("邮箱验证码错误")
	}

	// 删除邮箱验证码
	cache.DelEmailCode(registerReq.Email)

	// 对密码加密
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	// 随机生成不重复的用户名
	username, err := generateUniqueUsername()
	if err != nil {
		AddFailOperation(ctx, "注册", "用户名生成失败", gin.H{"邮箱": registerReq.Email}, err)
		return errors.New("用户名生成失败")
	}

	// 保存到数据库
	return global.Mysql.Create(&model.User{
		Username: username,
		Email:    registerReq.Email,
		Password: string(hashedPassword),
	}).Error
}

func UserLogin(ctx *gin.Context, loginReq dto.LoginReq) (accessToken, refreshToken string, err error) {
	// 读取数据库
	user, err := FindUserByEmail(loginReq.Email)
	if err != nil {
		AddFailOperation(ctx, "登录", "获取用户信息失败", gin.H{"邮箱": loginReq.Email}, nil)
		return "", "", errors.New("获取用户信息失败")
	}

	// 验证账号密码
	passwordError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if passwordError != nil {
		// 记录登录尝试次数
		cache.IncrLoginTryCount(loginReq.Email)
		AddFailOperation(ctx, "登录", "用户名密码不匹配", gin.H{"邮箱": loginReq.Email}, nil)
		return "", "", errors.New("用户名密码不匹配")
	}

	// 生成验证token
	if accessToken, err = jwt.GenerateAccessToken(user.ID); err != nil {
		AddFailOperation(ctx, "登录", "AccessToken生成失败", gin.H{"邮箱": loginReq.Email}, err)
		return "", "", errors.New("验证token生成失败")
	}
	// 生成刷新token
	if refreshToken, err = jwt.GenerateRefreshToken(user.ID); err != nil {
		AddFailOperation(ctx, "登录", "RefreshToken生成失败", gin.H{"邮箱": loginReq.Email}, err)
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
	_, claims, err := jwt_parse.ParseToken(tokenReq.RefreshToken)
	if err != nil {
		AddFailOperation(ctx, "刷新TOKEN", "解析TOKEN出错", gin.H{"Token值": tokenReq.RefreshToken}, err)
		return "", "", errors.New("token验证失败")
	}

	// 读取缓存
	if !cache.IsRefreshTokenExist(claims.UserId, tokenReq.RefreshToken) { // refreshToken 存在
		AddFailOperation(ctx, "刷新TOKEN", "RefreshToken不在缓存中", gin.H{"Token值": tokenReq.RefreshToken}, nil)
		return "", "", errors.New("无效Token")

	}

	// 移除refreshToken
	cache.DelRefreshToken(claims.UserId, tokenReq.RefreshToken)

	// 重新生成refreshToken
	refreshToken, err = jwt.GenerateRefreshToken(claims.UserId)
	if err != nil {
		AddFailOperation(ctx, "刷新TOKEN", "RefreshToken生成失败", nil, err)
		return "", "", errors.New("Token生成失败")
	}

	// 刷新accessToken
	accessToken, err = jwt_parse.GenerateAccessToken(claims.UserId)
	if err != nil {
		AddFailOperation(ctx, "刷新TOKEN", "AccessToken生成失败", nil, err)
		return "", "", errors.New("Token生成失败")
	}

	// 用户ID写入Cookie
	SetUserIdCookie(ctx, claims.UserId)

	// 存入缓存
	cache.SetRefreshToken(claims.UserId, refreshToken)

	return accessToken, refreshToken, nil
}

// 生成用户Id和md5并写入Cookie
func SetUserIdCookie(ctx *gin.Context, userId uint) {
	salt := viper.GetString("user_id_salt")
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
			zap.L().Error(fmt.Sprintf("Service | 用户信息查询失败 | 错误信息: %s", err.Error()))
		}

		// 存到redis
		cache.SetUser(user)
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
			zap.L().Error(fmt.Sprintf("Service | 用户信息查询失败 | 错误信息: %s", err.Error()))
		}

		// 存到redis
		cache.SetUser(user)
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
