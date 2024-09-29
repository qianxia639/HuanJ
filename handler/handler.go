package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	Router *gin.Engine
	DB     *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	handler := &Handler{
		DB: db,
	}

	handler.setupRouter()

	return handler
}

func (handler *Handler) setupRouter() {
	router := gin.Default()

	router.POST("/login", Login)
	router.POST("/user", handler.CreateUser)
	router.GET("/get", handler.Get)

	handler.Router = router
}
