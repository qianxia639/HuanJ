package handler

import (
	"Dandelion/db/models"
	"Dandelion/db/service"
	"Dandelion/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	CheckPassword string `json:"check_password" binding:"required"`
	Nickname      string `json:"nickname" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Gender        int8   `json:"gender"`
}

func (h *Handler) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
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
	if i := h.Store.ExistsUsername(ctx, req.Username); i > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "用户名已存在"})
		return
	}
	// 判断昵称是否存在
	if i := h.Store.ExistsNickname(ctx, req.Nickname); i > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "昵称已存在"})
		return
	}

	// 密码加密
	salt := fmt.Sprintf("%d", time.Now().Local().UnixNano())
	hashPwd, err := utils.HashPassword(req.Password, salt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Encoding password faild", "error": err.Error()})
		return
	}
	// 创建用户

	now := time.Now()

	args := &service.CreateUserParams{
		Username:  req.Username,
		Nickname:  req.Nickname,
		Password:  hashPwd,
		Salt:      salt,
		Email:     req.Email,
		Gender:    req.Gender,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = h.Store.CreateUser(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "insert error", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully"})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context) {

	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := "admin"

	if req.Username != arg || req.Password != arg {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "权限受限"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: req.Password,
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "login success...", "user": user})
}
