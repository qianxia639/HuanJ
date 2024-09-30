package handler

import (
	db "Dandelion/db/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Router  *gin.Engine
	Queries *db.Queries
}

func NewHandler(queries *db.Queries) *Handler {
	handler := &Handler{
		Queries: queries,
	}

	handler.setupRouter()

	return handler
}

func (handler *Handler) setupRouter() {
	router := gin.Default()

	router.POST("/login", Login)
	router.POST("/user", handler.CreateUser)

	handler.Router = router
}
