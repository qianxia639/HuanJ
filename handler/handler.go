package handler

import (
	"Dandelion/config"
	db "Dandelion/db/service"
	"Dandelion/token"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Router  *gin.Engine
	Conf    config.Config
	Queries *db.Queries
	Token   token.Maker
}

func NewHandler(conf config.Config, queries *db.Queries) *Handler {

	maker := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)

	handler := &Handler{
		Conf:    conf,
		Queries: queries,
		Token:   maker,
	}

	handler.setupRouter()

	return handler
}

func (handler *Handler) setupRouter() {
	router := gin.Default()

	router.Use(handler.CORS())

	authRouter := router.Group("")
	authRouter.Use(handler.authorizationMiddleware())

	// User Router
	router.POST("/login", handler.login)
	router.POST("/user", handler.createUser)

	authRouter.GET("/user", handler.getUser)
	authRouter.PUT("/user", handler.updateUser)

	// Friend Router
	authRouter.POST("/friend", handler.createdFriend)
	authRouter.GET("/friend", handler.getFriends)

	handler.Router = router
}
