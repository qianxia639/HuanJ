package handler

import (
	"Ice/db/service"
	db "Ice/db/sqlc"
	"Ice/internal/config"
	"Ice/internal/token"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	Router          *gin.Engine
	Conf            config.Config
	Queries         *service.Queries
	Store           db.Store
	Token           token.Maker
	CurrentUserInfo db.LoginUserInfo
	Redis           *redis.Client
}

func NewHandler(conf config.Config, store db.Store, rdb *redis.Client) *Handler {

	maker := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)

	handler := &Handler{
		Conf:    conf,
		Queries: nil,
		Store:   store,
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

	// Friend Request Router
	authRouter.POST("/friend/request", handler.createFriendRequest)
	authRouter.POST("/friend/request/accept/:id", handler.acceptFriendRequest)
	authRouter.POST("/friend/request/reject/:id", handler.rejectFriendRequest)

	// Friendship Router
	authRouter.POST("/friendship", handler.createdFriend)
	authRouter.GET("/friendship", handler.getFriends)
	authRouter.DELETE("/friendship/:id", handler.deleteFriend)

	// Group Router
	authRouter.POST("/group", handler.createGroup)

	handler.Router = router
}
