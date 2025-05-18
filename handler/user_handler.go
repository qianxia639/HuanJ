package handler

import (
	db "Rejuv/db/sqlc"
	"Rejuv/logs"
	"Rejuv/token"
	"Rejuv/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username      string `json:"username" binding:"required"`       // 用户名
	Password      string `json:"password" binding:"required"`       // 用户密码
	CheckPassword string `json:"check_password" binding:"required"` // 确认密码
	Email         string `json:"email" binding:"required,email"`    // 用户邮箱
	Gender        int8   `json:"gender" binding:"required,gender"`  // 用户性别
}

func (h *Handler) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// ctx.JSON(http.StatusBadRequest, gin.H{
		// 	"message": "参数错误",
		// 	"error":   err.Error(),
		// })
		utils.ParamsError(ctx)
		return
	}

	if !utils.ValidatePassword(req.Password) {
		// ctx.JSON(http.StatusBadRequest, gin.H{"message": "密码格式不正确"})
		utils.ParamsError(ctx, "密码格式不正确")
		return
	}

	if !utils.ValidateUsername(req.Username) {
		// ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户名格式不正确"})
		utils.ParamsError(ctx, "用户名格式不正确")
		return
	}

	// 判断密码是否一致
	if req.Password != req.CheckPassword {
		// ctx.JSON(http.StatusBadRequest, gin.H{"message": "密码不一致"})
		utils.ParamsError(ctx, "确认密码不一致")
		return
	}

	// 判断用户名是否存在
	if i, _ := h.Store.ExistsUsername(ctx, req.Username); i > 0 {
		// ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户名已存在"})
		utils.ParamsError(ctx, "用户名已存在")
		return
	}

	// 判断邮箱是否存在
	if i, _ := h.Store.ExistsEmail(ctx, req.Email); i > 0 {
		// ctx.JSON(http.StatusBadRequest, gin.H{"message": "邮箱已存在"})
		utils.ParamsError(ctx, "邮箱已存在")
		return
	}

	// 密码加密
	hashPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		// ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		utils.ParamsError(ctx)
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
		// ctx.JSON(http.StatusBadRequest, gin.H{"message": "insert error", "error": err.Error()})
		logs.Errorf("Create User Fail,Err:[%v]", err)
		utils.ServerError(ctx)
		return
	}

	// ctx.JSON(http.StatusOK, gin.H{"message": "successfully"})
	utils.Success(ctx, "Create Success")
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
		// ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户名不存存", "error": err.Error()})
		// ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户名不存存"})
		Error(ctx, http.StatusUnauthorized, "用户名不存在")
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
	// ctx.JSON(http.StatusOK, gin.H{"message": "successfullt", "data": tokenStr, "user": user})
	Success(ctx, tokenStr)
}

func (h *Handler) getUser(ctx *gin.Context) {

	utils.Obj(ctx, h.CurrentUserInfo)

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
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if req.Nickname != nil && *req.Nickname != h.CurrentUserInfo.Nickname {
		// 判断用户昵称是否重复
		if i, _ := h.Store.ExistsNickname(ctx, *req.Nickname); i > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户昵称存在"})
			return
		}
		h.CurrentUserInfo.Nickname = *req.Nickname
	}

	if req.Gender != nil {
		h.CurrentUserInfo.Gender = *req.Gender
	}

	args := &db.UpdateUserParams{
		Gender:   h.CurrentUserInfo.Gender,
		ID:       h.CurrentUserInfo.ID,
		Nickname: h.CurrentUserInfo.Nickname,
	}
	err = h.Store.UpdateUser(ctx, args)
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
