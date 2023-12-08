package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wangzmgit/jigsaw"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
)

type CaptchaApi struct{}

// 获取滑块验证码
func (c *CaptchaApi) GetSliderCaptcha(ctx *gin.Context) {
	//获取参数
	captchaId := ctx.DefaultQuery("captchaId", "")

	// 参数校验
	if cache.GetCaptchaStatus(captchaId) != global.CAPTCHA_STATUS_NOT_USED {
		service.AddFailOperation(ctx, "获取滑块", "滑块ID不为未验证状态", gin.H{"滑块ID": captchaId}, nil)
		resp.FailWithMessage(ctx, "获取失败")
		return
	}

	// 生成滑块验证
	slider, bg, x, y, err := jigsaw.Create()
	if err != nil {
		service.AddFailOperation(ctx, "获取滑块", "滑块生成失败", gin.H{"滑块ID": captchaId}, err)
		resp.FailWithMessage(ctx, "滑块生成失败")
		return
	}

	// 保存x坐标到缓存
	cache.SetSliderX(captchaId, x)

	// 返回给前端
	service.AddOkOperation(ctx, "获取滑块", "滑块生成失败", gin.H{"滑块ID": captchaId, "位置X": x, "位置Y": y})
	resp.OkWithData(ctx, gin.H{"slider_captcha": vo.CaptchaResp{
		BgImg:     bg,
		SliderImg: slider,
		Y:         y,
	}})
}

// 验证滑块
func (c *CaptchaApi) ValidateSlider(ctx *gin.Context) {
	// 获取参数
	var sliderReq dto.SliderReq
	if err := ctx.Bind(&sliderReq); err != nil {
		service.AddFailOperation(ctx, "验证滑块", "请求参数有误", nil, nil)
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	//获取缓存中的x并删除缓存
	x := cache.GetSliderX(sliderReq.CaptchaId)
	cache.DelSlider(sliderReq.CaptchaId)
	//与用户拖动位置对比
	if sliderReq.X < x-3 || sliderReq.X > x+3 {
		service.AddFailOperation(ctx, "验证滑块", "滑块验证失败", gin.H{"滑块ID": sliderReq.CaptchaId, "位置X": sliderReq.X}, nil)
		resp.FailWithMessage(ctx, "滑块验证失败")
		return
	}

	cache.SetCaptchaStatus(sliderReq.CaptchaId, global.CAPTCHA_STATUS_PASS)
	// 返回给前端
	service.AddOkOperation(ctx, "验证滑块", "请求参数有误", gin.H{"滑块ID": sliderReq.CaptchaId})
	resp.Ok(ctx)
}
