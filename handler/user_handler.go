package handler

import (
	"Dandelion/db/models"
	"fmt"
	"log"
	"net/http"
	"time"

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

func (h *Handler) Get(ctx *gin.Context) {
	var u models.User
	rows, err := h.DB.QueryContext(ctx, "SELECT *FROM users")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "query error", "error": err.Error()})
		return
	}

	for rows.Next() {
		err = rows.Scan(&u)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "scan error", "error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
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

	// 判断用户是否存在

	// 密码加密

	// 创建用户
	sql := `
		INSERT INTO users (username, nickname, password, salt, email, gender, created_at, updated_at) 
		VALUES ($1, $2, $3, $4,$5, $6, $7, $8)
	`
	salt := fmt.Sprintf("%d", time.Now().Local().UnixNano())
	log.Printf("salt: %s\n", string(salt))
	t := time.Now()
	rows, err := h.DB.QueryContext(ctx, sql, req.Username, req.Nickname, req.Password, string(salt), req.Email, req.Gender, t, t)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "insert error", "error": err.Error()})
		return
	}
	var u models.User
	for rows.Next() {
		_ = rows.Scan(&u)
	}

	ctx.JSON(http.StatusOK, gin.H{"data": u, "error": err.Error()})
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
