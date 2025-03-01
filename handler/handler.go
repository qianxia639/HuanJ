package handler

import (
	db "Ice/db/sqlc"
	"Ice/internal/config"
	"Ice/internal/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	Router          *gin.Engine
	Conf            config.Config
	Store           db.Store
	Token           token.Maker
	CurrentUserInfo db.LoginUserInfo
	Redis           *redis.Client
}

func NewHandler(conf config.Config, store db.Store, rdb *redis.Client) *Handler {

	maker := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)

	handler := &Handler{
		Conf:  conf,
		Store: store,
		Token: maker,
		Redis: rdb,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("gender", validGender)
	}

	handler.setupRouter()

	return handler
}

func (handler *Handler) setupRouter() {
	router := gin.Default()

	router.Use(handler.CORS())

	router.GET("/secret", handler.secret(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully..."})
	})

	router.GET("/ws", handler.wsHandler)

	authRouter := router.Group("")
	authRouter.Use(handler.authorizationMiddleware())

	// User Router
	{
		router.POST("/login", handler.login)
		router.POST("/user", handler.createUser)

		authRouter.GET("/user", handler.getUser)
		authRouter.PUT("/user", handler.updateUser)
	}

	// Friend Request Router
	{
		authRouter.POST("/friend/request", handler.createFriendRequest)
		authRouter.POST("/friend/request/accept/:id", handler.acceptFriendRequest)
		authRouter.POST("/friend/request/reject/:id", handler.rejectFriendRequest)
	}

	// Friendship Router
	// authRouter.POST("/friendship", handler.createdFriend)
	{
		authRouter.GET("/friendship", handler.getFriends)
		authRouter.DELETE("/friendship/:id", handler.deleteFriend)
	}

	// Group Router
	{
		authRouter.POST("/group", handler.createGroup)
	}

	handler.Router = router
}
