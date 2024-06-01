package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/cache"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
	"interastral-peace.com/alnitak/utils"
)

func GetUserInfo(ctx *gin.Context) {
	userId := ctx.GetUint("userId")
	user := service.GetUserInfo(userId)

	resp.OkWithData(ctx, gin.H{"userInfo": user})
}

func EditUserInfo(ctx *gin.Context) {
	var editUserInfoReq dto.EditUserInfoReq
	if err := ctx.Bind(&editUserInfoReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(editUserInfoReq.Name, "=", 0) {
		resp.FailWithMessage(ctx, "用户名不能为空")
		return
	}

	if err := service.EditUserInfo(ctx, editUserInfoReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 获取用户基本信息
func GetUserBaseInfo(ctx *gin.Context) {
	userId := utils.StringToUint(ctx.Query("userId"))

	user := service.GetUserBaseInfo(userId)
	if user.ID == 0 {
		resp.FailWithMessage(ctx, "用户不存在")
		return
	}

	resp.OkWithData(ctx, gin.H{"userInfo": user})
}

// 重置密码检查
func ResetPwdCheck(ctx *gin.Context) {
	// 获取参数
	email := ctx.Query("email")
	captchaId := ctx.Query("captchaId")

	// 参数校验
	if !utils.VerifyEmail(email) {
		resp.FailWithMessage(ctx, "邮箱格式错误")
		return
	}

	if utils.VerifyStringLength(captchaId, "=", 0) {
		captchaId := cache.CreateCaptchaStatus()
		resp.Result(ctx, -1, gin.H{"captchaId": captchaId}, "需要人机验证")
		return
	}

	switch cache.GetCaptchaStatus(captchaId) {
	case global.CAPTCHA_STATUS_ABSENT: // 如果未进行人机验证
		captchaId := cache.CreateCaptchaStatus()
		resp.Result(ctx, -1, gin.H{"captchaId": captchaId}, "需要人机验证")
		return
	case global.CAPTCHA_STATUS_PASS:
		cache.DelCaptchaStatus(captchaId) // 删除人机验证状态
	}

	if err := service.ResetPwdCheck(ctx, email); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 修改密码
func ModifyPwd(ctx *gin.Context) {
	var modifyPwdReq dto.ModifyPwdReq
	if err := ctx.Bind(&modifyPwdReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if !utils.VerifyEmail(modifyPwdReq.Email) {
		resp.FailWithMessage(ctx, "邮箱格式错误")
		return
	}

	if utils.VerifyStringLength(modifyPwdReq.Password, "<", 6) {
		resp.FailWithMessage(ctx, "密码长度不能小于6位")
		return
	}

	if !utils.VerifyStringLength(modifyPwdReq.Code, "=", 6) {
		resp.FailWithMessage(ctx, "验证码长度为6位")
		return
	}

	if err := service.ModifyPwd(ctx, modifyPwdReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 获取用户列表(后台管理)
func GetUserListManage(ctx *gin.Context) {
	// 获取参数
	var userListReq dto.UserListReq
	if err := ctx.Bind(&userListReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if userListReq.PageSize > 100 {
		resp.FailWithMessage(ctx, "请求数量过多")
		return
	}

	total, roles := service.GetUserListManage(userListReq)

	// 返回给前端
	resp.OkWithData(ctx, gin.H{"list": roles, "total": total})
}

// 编辑用户信息(后台管理)
func EditUserInfoManage(ctx *gin.Context) {
	var editUserInfoManageReq dto.EditUserInfoManageReq
	if err := ctx.Bind(&editUserInfoManageReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 参数校验
	if utils.VerifyStringLength(editUserInfoManageReq.Name, "=", 0) {
		resp.FailWithMessage(ctx, "用户名不能为空")
		return
	}

	if err := service.EditUserInfoManage(ctx, editUserInfoManageReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回给前端
	resp.Ok(ctx)
}

// 设置用户角色
func EditUserRole(ctx *gin.Context) {
	// 获取参数
	var editUserRoleReq dto.EditUserRoleReq
	if err := ctx.Bind(&editUserRoleReq); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.EditUserRole(ctx, editUserRoleReq); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}

// 删除用户
func DeleteUser(ctx *gin.Context) {
	// 获取参数
	id := utils.StringToUint(ctx.Param("id"))

	if err := service.DeleteUser(ctx, id); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	// 返回
	resp.Ok(ctx)
}
