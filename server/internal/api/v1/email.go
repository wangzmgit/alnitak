package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/pkg/mail"
	"interastral-peace.com/alnitak/utils"
)

type EmailApi struct{}

func (e *EmailApi) SendRegisterEmailCode(ctx *gin.Context) {
	// 获取参数
	var sendEmailReq dto.SendEmailReq
	if err := ctx.Bind(&sendEmailReq); err != nil {
		service.AddFailOperation(ctx, "发送验证码", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if !utils.VerifyEmail(sendEmailReq.Email) {
		service.AddFailOperation(ctx, "发送验证码", "邮箱格式错误", gin.H{"邮箱": sendEmailReq.Email}, nil)
		resp.FailWithMessage(ctx, "邮箱格式错误")
		return
	}

	// 如果未进行人机验证
	if cache.GetCaptchaStatus(sendEmailReq.CaptchaId) == global.CAPTCHA_STATUS_ABSENT {
		captchaId := cache.CreateCaptchaStatus()

		service.AddFailOperation(ctx, "发送验证码", "需要人机验证", gin.H{"邮箱": sendEmailReq.Email, "滑块ID": captchaId}, nil)
		resp.Result(ctx, -1, gin.H{"captchaId": captchaId}, "需要人机验证")
		return
	}

	// 生成code
	code := utils.GenerateNumberCode(6)

	// 发送code
	if viper.GetBool("mail.debug") {
		// 验证码debug模式不发送邮件
		zap.L().Debug("邮箱:" + sendEmailReq.Email + ",验证码:" + code)
	} else {
		if err := mail.SendCaptcha(sendEmailReq.Email, code); err != nil {
			service.AddFailOperation(ctx, "发送验证码", "需要人机验证", gin.H{"邮箱": sendEmailReq.Email}, err)
			resp.FailWithMessage(ctx, "邮箱验证码发送失败")
			return
		}
	}

	// code放入缓存
	cache.SetEmailCode(sendEmailReq.Email, code)
	// 删除人机验证状态
	cache.DelCaptchaStatus(sendEmailReq.CaptchaId)

	service.AddOkOperation(ctx, "发送验证码", "发送验证码成功", gin.H{"邮箱": sendEmailReq.Email, "验证码": code})
	resp.Ok(ctx)
}
