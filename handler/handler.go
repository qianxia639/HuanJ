package handler

import (
	"Dandelion/internal/config"
	"Dandelion/internal/db/model"
	db "Dandelion/internal/db/service"
	"Dandelion/internal/token"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	Router          *gin.Engine
	Conf            config.Config
	Queries         *db.Queries
	Token           token.Maker
	CurrentUserInfo model.LoginUserInfo
	Redis           *redis.Client
}

func NewHandler(conf config.Config, queries *db.Queries, rdb *redis.Client) *Handler {

	maker := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)

	handler := &Handler{
		Conf:    conf,
		Queries: queries,
		Token:   maker,
		Redis:   rdb,
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
	authRouter.DELETE("/friend", handler.deleteFriend)

	handler.Router = router
}
