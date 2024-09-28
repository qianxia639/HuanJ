package handler

import (
	"Dandelion/db/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	CheckPassword string `json:"check_password" binding:"required"`
	Nickname      string `json:"nickname" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Gender        int8   `json:"gender"`
}

func CreateUser(ctx *gin.Context) {
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

	// 判断用户是否存在

	// 密码加密

	// 创建用户

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
