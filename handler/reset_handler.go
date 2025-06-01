package handler

import (
	db "HuanJ/db/sqlc"
	"HuanJ/logs"
	"HuanJ/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type resetPwdRequest struct {
	Email     string `json:"email" binding:"required,email"` // 用户邮箱
	EmailCode string `json:"email_code" binding:"required"`  // 邮件验证码
	Password  string `json:"password" binding:"required"`    // 用户密码
	ResetPwd  string `json:"reset_pwd" binding:"required"`   // 确认密码
}

// 重置密码
func (h *Handler) resetPwd(ctx *gin.Context) {
	var req resetPwdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logs.Errorf("Bind params error: %v", err)
		h.ParamsError(ctx)
		return
	}

	if req.Password != req.ResetPwd {
		h.ParamsError(ctx, "两次输入密码不一致")
		return
	}

	if !utils.ValidatePassword(req.ResetPwd) {
		h.ParamsError(ctx, "密码格式不正确")
		return
	}

	// 校验邮箱是否存在
	if count, _ := h.Store.ExistsEmail(ctx, req.Email); count < 1 {
		h.Error(ctx, http.StatusInternalServerError, "邮箱号不存在")
		return
	}

	// 校验邮件验证码是否正确
	// ok, err := mail.VerifyEmailCode(h.RedisClient, req.Email, req.EmailCode, 2)
	// if err != nil {
	// 	h.ServerError(ctx)
	// 	return
	// }

	// if !ok {
	// 	h.ParamsError(ctx, "邮箱验证码错误或过期")
	// 	return
	// }

	hashPwd, err := utils.HashPassword(req.ResetPwd)
	if err != nil {
		h.ServerError(ctx)
		return
	}

	_ = h.Store.UpdatePwd(ctx, &db.UpdatePwdParams{
		ID:       h.CurrentUserInfo.ID,
		Email:    req.Email,
		Password: hashPwd,
	})

	h.Success(ctx, "密码重置成功")
}
