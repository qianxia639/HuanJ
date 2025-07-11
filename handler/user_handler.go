package handler

import (
	"HuanJ/config"
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

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) login(ctx *gin.Context) {

	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.ParamsError(ctx)
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
		h.Error(ctx, http.StatusUnauthorized, "密码错误")
		return
	}
	// 生成Token
	resp, err := h.generateTokens(user.Username)
	if err != nil {
		logs.Errorf("Token生成失败 [%s]: %v", user.Username, err)
		h.ServerError(ctx)
		return
	}

	// 邮箱脱敏
	// user.Email = utils.MaskEmail(user.Email)

	loginUserInfo := db.LoginUserInfo{
		User:      user,
		UserAgent: ctx.Request.Header.Get("User-Agent"),
		LoginIp:   ctx.ClientIP(),
	}
	key := "user:" + user.Username
	err = h.RedisClient.Set(ctx, key, &loginUserInfo, 24*time.Hour).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	h.Obj(ctx, resp)
}

func (h *Handler) generateTokens(username string) (loginResponse, error) {
	// 生成access token
	accessToken, err := h.Token.CreateToken(token.Token{
		Username: username,
		Duration: h.Conf.Token.AccessTokenDuration,
	})
	if err != nil {
		return loginResponse{}, err
	}

	// 生成refresh token
	refreshToken, err := h.Token.CreateToken(token.Token{
		Username: username,
		Duration: h.Conf.Token.RefreshTokenDuration,
	})

	return loginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, err
}

type refreshTokenRequest struct {
	Token string `json:"token"`
}

func (h *Handler) refreshToken(ctx *gin.Context) {
	// 获取token
	var req refreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.ParamsError(ctx)
		return
	}
	// 校验token
	payload, err := h.Token.VerifyToken(req.Token)
	if err != nil {
		h.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	// 生成新的access token
	resp, err := h.generateTokens(payload.Username)
	if err != nil {
		logs.Errorf("Token生成失败 [%s]: %v", payload.Username, err)
		h.ServerError(ctx)
		return
	}
	// 返回新的access token
	h.Obj(ctx, resp)
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
	OldPassword string `json:"old_pwd" binding:"required"`    // 旧密码
	NewPassword string `json:"new_pwd" binding:"required"`    // 新密码
	EmailCode   string `json:"email_code" binding:"required"` // 邮箱验证码
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
	ok, err := mail.VerifyEmailCode(h.RedisClient, h.CurrentUserInfo.Email, req.EmailCode, config.EmailCodeResetPwd)
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
		Email:    h.CurrentUserInfo.Email,
		Password: hashPassword,
	})
	if err != nil {
		h.ServerError(ctx)
		return
	}
	// TODO 是否有必要删除缓存
	// key := "user:" + h.CurrentUserInfo.Username
	// _ = h.RedisClient.Del(ctx, key).Err()

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

// 用户登出
func (h *Handler) logout(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get(authorizationHeader)
	if len(authorization) == 0 {
		h.ParamsError(ctx, "缺少Token")
		return
	}

	payload, err := h.Token.VerifyToken(authorization)
	if err != nil {
		logs.Errorf("解析token失败: %v\n", err)
		h.ServerError(ctx)
		return
	}

	// 计算剩余时间(秒)
	ttl := time.Until(payload.ExpiredAt)

	if ttl < time.Second {
		ttl = time.Second // 确保最小生存时间
	}

	// 将令牌加入黑名单
	err = h.RedisClient.Set(ctx, "token_blacklist:"+authorization, 1, ttl).Err()
	if err != nil {
		logs.Errorf("token黑名单设置失败 %v\n", err)
		h.ServerError(ctx)
		return
	}

	h.Success(ctx, "登出成功")
}
