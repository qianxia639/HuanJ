package handler

import (
	"Dandelion/db/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Router *gin.Engine
	Store  service.Store
}

func NewHandler(store service.Store) *Handler {
	handler := &Handler{
		Store: store,
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
