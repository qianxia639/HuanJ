package handler

import (
	"Dandelion/db/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
