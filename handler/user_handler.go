package handler

import (
	db "HuanJ/db/sqlc"
	"HuanJ/logs"
	"HuanJ/mail"
	"HuanJ/token"
	"HuanJ/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type createUserRequest struct {
	Username        string `json:"username" binding:"required"`         // 用户名
	Password        string `json:"password" binding:"required"`         // 用户密码
	ConfirmPassword string `json:"confirm_password" binding:"required"` // 确认密码
	Email           string `json:"email" binding:"required,email"`      // 用户邮箱
	Gender          int8   `json:"gender" binding:"required,gender"`    // 用户性别
}

func (h *Handler) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.ParamsError(ctx)
		return
	}

	if !utils.ValidatePassword(req.Password) {
		h.ParamsError(ctx, "密码格式不正确")
		return
	}

	if !utils.ValidateUsername(req.Username) {
		h.ParamsError(ctx, "用户名格式不正确")
		return
	}

	// 判断密码是否一致
	if req.Password != req.ConfirmPassword {
		h.ParamsError(ctx, "密码不一致")
		return
	}

	// 判断用户名是否存在
	if i, _ := h.Store.ExistsUsername(ctx, req.Username); i > 0 {
		h.ParamsError(ctx, "用户名已存在")
		return
	}

	// 判断邮箱是否存在
	if i, _ := h.Store.ExistsEmail(ctx, req.Email); i > 0 {
		h.ParamsError(ctx, "邮箱已存在")
		return
	}

	// 密码加密
	hashPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		h.ParamsError(ctx)
		return
	}

	// 创建用户
	args := &db.CreateUserParams{
		Username: req.Username,
		Nickname: req.Username,
		Password: hashPwd,
		Email:    req.Email,
		Gender:   req.Gender,
	}

	_, err = h.Store.CreateUser(ctx, args)
	if err != nil {
		logs.Errorf("Create User Fail,Err:[%v]", err)
		h.ServerError(ctx)
		return
	}

	h.Success(ctx, "Create Success")
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) login(ctx *gin.Context) {

	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 判断用户是否存在
	user, err := h.Store.GetUser(ctx, req.Username)
	if user.ID == 0 {
		h.Error(ctx, http.StatusUnauthorized, "用户名不存在")
		return
	}
	// 校验密码
	err = utils.ComparePassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// 生成Token
	args := token.Token{
		Username: user.Username,
		Duration: h.Conf.Token.AccessTokenDuration,
	}
	tokenStr, err := h.Token.CreateToken(args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 邮箱脱敏
	user.Email = utils.MaskEmail(user.Email)

	loginUserInfo := db.LoginUserInfo{
		User:      user,
		UserAgent: ctx.Request.Header.Get("User-Agent"),
		LoginIp:   ctx.ClientIP(),
	}
	key := fmt.Sprintf("user:%s", user.Username)
	err = h.RedisClient.Set(ctx, key, &loginUserInfo, 24*time.Hour).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	h.Success(ctx, tokenStr)
}

func (h *Handler) getUser(ctx *gin.Context) {

	h.Obj(ctx, h.CurrentUserInfo)

	// ctx.JSON(http.StatusOK, gin.H{"message": "successfully", "data": h.CurrentUserInfo})
}

type updateUserRequest struct {
	Gender   *int8   `json:"gender"`
	Nickname *string `json:"nickname"`
}

func (h *Handler) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	auth := ctx.Request.Header.Get(authorizationHeader)
	payload, err := h.Token.VerifyToken(auth)
	if err != nil {
		h.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if req.Nickname != nil && *req.Nickname != h.CurrentUserInfo.Nickname {
		// 判断用户昵称是否重复
		if i, _ := h.Store.ExistsNickname(ctx, *req.Nickname); i > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户昵称存在"})
			return
		}
	}

	args := &db.UpdateUserParams{
		ID: h.CurrentUserInfo.ID,
		Gender: pgtype.Int2{
			Int16: int16(*req.Gender),
			Valid: req.Gender != nil,
		},
		Nickname: pgtype.Text{
			String: *req.Nickname,
			Valid:  req.Nickname != nil,
		},
	}
	h.CurrentUserInfo.User, err = h.Store.UpdateUser(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Update user failed", "error": err.Error()})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Update user failed", "error": err.Error()})
		return
	}

	_ = h.RedisClient.Set(ctx, fmt.Sprintf("user:%s", payload.Username), &h.CurrentUserInfo, 24*time.Hour)

	h.CurrentUserInfo.Email = utils.MaskEmail(h.CurrentUserInfo.Email)

	ctx.JSON(http.StatusOK, gin.H{"data": h.CurrentUserInfo})
}

type updatePwd struct {
	OldPassword string `json:"old_pwd" binding:"required"`
	NewPassword string `json:"new_pwd" binding:"required"`
	EmailCode   string `json:"email_code" binding:"required"`
}

// 修改密码
func (h *Handler) updatePassword(ctx *gin.Context) {
	var req updatePwd
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.ParamsError(ctx)
		return
	}

	if req.OldPassword == req.NewPassword {
		h.ParamsError(ctx, "新旧密码不能相同")
		return
	}

	// 校验邮箱验证码
	ok, err := mail.VerifyEmailCode(h.RedisClient, h.CurrentUserInfo.Email, req.EmailCode, 1)
	if err != nil {
		h.ServerError(ctx)
		return
	}

	if !ok {
		h.Error(ctx, 1, "邮箱验证码错误")
		return
	}

	//
	if time.Since(h.CurrentUserInfo.PasswordChangedAt) < 7*24*time.Hour {
		h.ParamsError(ctx, "两次密码修改时间间隔不得低于7天")
		return
	}

	hashPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		h.ServerError(ctx)
		return
	}

	// 更新密码
	err = h.Store.UpdatePwd(ctx, &db.UpdatePwdParams{
		ID:       h.CurrentUserInfo.ID,
		Email:    h.CurrentUserInfo.Email,
		Password: hashPassword,
	})
	if err != nil {
		h.ServerError(ctx)
		return
	}

	h.Success(ctx, "Update success")
}

type sendEmail struct {
	Email         string `json:"email" binding:"required"`           // 邮箱
	EmailCodeType int8   `json:"email_code_type" binding:"required"` // 邮箱验证码类型
}

func (h *Handler) sendEmail(ctx *gin.Context) {
	var req sendEmail
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.ParamsError(ctx)
		return
	}

	// 发送邮箱验证码
	err := mail.SendEmailCode(h.RedisClient, req.Email, req.EmailCodeType)
	if err != nil {
		logs.Error(err)
		h.ServerError(ctx)
		return
	}

	h.Success(ctx, "Send email success")
}
