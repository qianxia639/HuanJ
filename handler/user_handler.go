package handler

import (
	"Dandelion/db/models"
	db "Dandelion/db/service"
	"Dandelion/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	CheckPassword string `json:"check_password" binding:"required"`
	Nickname      string `json:"nickname" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Gender        int8   `json:"gender"`
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
	if i := h.Queries.ExistsUser(ctx, req.Username, req.Email); i > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户名已存在"})
		return
	}
	// 判断昵称是否存在
	if i := h.Queries.ExistsNickname(ctx, req.Nickname); i > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "昵称已存在"})
		return
	}

	// 判断邮箱是否存在

	// 密码加密
	salt := utils.GenerateSalt()
	hashPwd, err := utils.HashPassword(req.Password, salt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Encoding password faild", "error": err.Error()})
		return
	}
	// 创建用户

	now := time.Now()

	args := &db.CreateUserParams{
		Username:  req.Username,
		Nickname:  req.Nickname,
		Password:  hashPwd,
		Salt:      salt,
		Email:     req.Email,
		Gender:    req.Gender,
		CreatedAt: now,
		UpdatedAt: now,
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
	err = utils.ComparePassword(req.Password, user.Password, user.Salt)
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

	// 返回结果
	// ctx.JSON(http.StatusOK, gin.H{"message": "successfullt", "data": tokenStr, "user": user})
	Success(ctx, tokenStr)
}

func (h *Handler) getUser(ctx *gin.Context) {

	// k, exists := ctx.Get(authorizationPayloadKey)
	// if !exists {
	// 	ctx.JSON(http.StatusNotFound, gin.H{"error": "Not user"})
	// 	return
	// }
	// payload, ok := k.(*token.Payload)
	// if !ok {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "type assertion failed"})
	// 	return
	// }

	// user, err := h.Queries.GetUser(ctx, payload.Username)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully", "data": h.CurrentUser})
}

type updateUserRequest struct {
	Gender   *int8   `json:"gender"`
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
}

func (h *Handler) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {

	}
	user, ok := h.CurrentUser.(models.User)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "断言失败"})
		return
	}

	if req.Nickname != nil && *req.Nickname != user.Nickname {
		// 判断用户昵称是否重复
		if i := h.Queries.ExistsNickname(ctx, *req.Nickname); i > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户昵称重复"})
			return
		}
		user.Nickname = *req.Nickname
	}

	if req.Gender != nil {
		user.Gender = *req.Gender
	}

	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	h.CurrentUser = user

	err := h.Queries.UpdateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Update user failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": h.CurrentUser})
}
