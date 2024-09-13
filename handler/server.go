package handler

import "github.com/gin-gonic/gin"

type Server struct {
	Router *gin.Engine
}

func NewServer() *Server {
	var server Server

	server.setupRouter()

	return &server
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/login", Login)

	server.Router = router
}
