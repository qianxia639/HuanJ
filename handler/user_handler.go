package handler

import (
	"Ice/db/model"
	db "Ice/db/service"
	"Ice/internal/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	CheckPassword string `json:"check_password" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Answer        string `json:"answer" binding:"required"`
	Gender        int8   `json:"gender" binding:"required"`
}

func (h *Handler) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	if !utils.ValidatePassword(req.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "密码格式不正确"})
		return
	}

	if !utils.ValidateUsername(req.Username) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户名格式不正确"})

		return
	}

	// 判断密码是否一致
	if req.Password != req.CheckPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "密码不一致"})
		return
	}

	// 判断性别是否合法
	if _, exists := Gender[req.Gender]; !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "unknown gender"})
		return
	}
	// 判断用户名是否存在
	if i := h.Queries.ExistsUser(ctx, req.Username); i > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户名已存在"})
		return
	}

	if Answer != req.Answer {
		Error(ctx, http.StatusBadRequest, "验证失败")
		// ctx.JSON(http.StatusBadRequest, gin.H{"message": "验证失败"})
		return
	}

	// 判断邮箱是否存在
	if i := h.Queries.ExistsEmail(ctx, req.Email); i > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "邮箱已存在"})
		return
	}

	// 密码加密
	hashPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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

	err = h.Queries.CreateUser(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "insert error", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully"})
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) login(ctx *gin.Context) {

	ua := ctx.Request.Header.Get("User-Agent")
	if len(ua) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Can't find 'User-Agent' in header"})
		return
	}

	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 判断用户是否存在
	user, err := h.Queries.GetUser(ctx, req.Username)
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
	tokenStr, err := h.Token.CreateToken(user.Username, h.Conf.Token.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loginUserInfo := model.LoginUserInfo{
		User:      user,
		UserAgent: ua,
		LoginIp:   ctx.ClientIP(),
	}
	key := fmt.Sprintf("user:%s", user.Username)
	err = h.Redis.Set(ctx, key, &loginUserInfo, 24*time.Hour).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	// ctx.JSON(http.StatusOK, gin.H{"message": "successfullt", "data": tokenStr, "user": user})
	Success(ctx, tokenStr)
}

func (h *Handler) getUser(ctx *gin.Context) {

	h.CurrentUserInfo.Email = utils.MaskEmail(h.CurrentUserInfo.Email)

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully", "data": h.CurrentUserInfo})
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
	fields := strings.Fields(auth)
	payload, err := h.Token.VerifyToken(fields[1])
	if err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if req.Nickname != nil && len(*req.Nickname) > 3 && *req.Nickname != h.CurrentUserInfo.Nickname {
		// 判断用户昵称是否存在
		if i := h.Queries.ExistsNickname(ctx, *req.Nickname); i > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户昵称存在"})
			return
		}
		h.CurrentUserInfo.Nickname = *req.Nickname
	}

	if req.Gender != nil {
		h.CurrentUserInfo.Gender = *req.Gender
	}

	// h.CurrentUserInfo.UpdatedAt = time.Now()

	err = h.Queries.UpdateUser(ctx, h.CurrentUserInfo.User)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Update user failed", "error": err.Error()})
		return
	}

	_ = h.Redis.Set(ctx, fmt.Sprintf("user:%s", payload.Username), &h.CurrentUserInfo, 24*time.Hour)

	h.CurrentUserInfo.Email = utils.MaskEmail(h.CurrentUserInfo.Email)

	ctx.JSON(http.StatusOK, gin.H{"data": h.CurrentUserInfo})
}
