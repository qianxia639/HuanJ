package handler

import (
	"Dandelion/config"
	db "Dandelion/db/service"
	"Dandelion/token"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Router  *gin.Engine
	Queries *db.Queries
	Token   token.Maker
	Conf    config.Config
}

func NewHandler(queries *db.Queries, conf config.Config) (*Handler, error) {

	maker, err := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	handler := &Handler{
		Queries: queries,
		Token:   maker,
		Conf:    conf,
	}

	handler.setupRouter()

	return handler, nil
}

func (handler *Handler) setupRouter() {
	router := gin.Default()

	router.POST("/user/login", handler.login)
	router.POST("/user", handler.createUser)

	handler.Router = router
}
