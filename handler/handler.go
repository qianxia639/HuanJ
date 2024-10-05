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

func NewHandler(conf config.Config, queries *db.Queries) (*Handler, error) {

	maker, err := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	handler := &Handler{
		Conf:    conf,
		Queries: queries,
		Token:   maker,
	}

	handler.setupRouter()

	return handler, nil
}

func (handler *Handler) setupRouter() {
	router := gin.Default()

	userRouter := router.Group("/user")
	{
		userRouter.POST("/login", handler.login)
		userRouter.POST("/", handler.createUser)
	}
	userRouterAuth := userRouter.Use(handler.authorizationMiddleware())
	{

		userRouterAuth.GET("/", handler.getUser)
	}

	handler.Router = router
}
