package handler

import (
	"HuanJ/config"
	db "HuanJ/db/sqlc"
	"HuanJ/logs"
	"HuanJ/mail"
	"HuanJ/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type resetPwdRequest struct {
	Email           string `json:"email" binding:"required,email"`       // 用户邮箱
	EmailCode       string `json:"email_code" binding:"required"`        // 邮件验证码
	ResetPwd        string `json:"reset_pwd" binding:"required"`         // 用户密码
	ConfirmResetPwd string `json:"confirm_reset_pwd" binding:"required"` // 确认密码
}

// 重置密码
func (h *Handler) resetPwd(ctx *gin.Context) {
	var req resetPwdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logs.Errorf("Bind params error: %v", err)
		h.ParamsError(ctx)
		return
	}

	if req.ResetPwd != req.ConfirmResetPwd {
		h.ParamsError(ctx, "两次输入密码不一致")
		return
	}

	if !utils.ValidatePassword(req.ConfirmResetPwd) {
		h.ParamsError(ctx, "密码格式不正确")
		return
	}

	// 校验邮箱是否存在
	user, _ := h.Store.GetUserByEmail(ctx, req.Email)
	if user.ID < 1 {
		h.Error(ctx, http.StatusBadRequest, "邮箱号不存在")
		return
	}

	// 校验邮件验证码是否正确
	ok, err := mail.VerifyEmailCode(h.RedisClient, req.Email, req.EmailCode, config.EmailCodeResetPwd)
	if err != nil {
		h.ServerError(ctx)
		return
	}

	if !ok {
		h.ParamsError(ctx, "邮箱验证码错误或过期")
		return
	}

	// 校验新旧密码是否相同
	if err := utils.ComparePassword(req.ResetPwd, user.Password); err != nil {
		h.ParamsError(ctx, "新密码不能与旧密码相同")
		return
	}

	hashPwd, err := utils.HashPassword(req.ConfirmResetPwd)
	if err != nil {
		h.ServerError(ctx)
		return
	}

	// 校验与上次修改密码的间隔时间
	if time.Since(user.PasswordChangedAt) < 7*24*time.Hour {
		h.ParamsError(ctx, "两次密码修改时间间隔不得低于7天")
		return
	}

	_ = h.Store.UpdatePwd(ctx, &db.UpdatePwdParams{
		Email:    req.Email,
		Password: hashPwd,
	})

	h.Success(ctx, "密码重置成功")
}
