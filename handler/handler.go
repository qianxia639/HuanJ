package handler

import (
	"Dandelion/config"
	db "Dandelion/db/service"
	"Dandelion/token"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	handlerOptions

	Router *gin.Engine
}

type handlerOptions struct {
	queries    *db.Queries
	tokenMaker token.Maker
	conf       config.Config
}

func NewHandler(opts ...HandlerOption) (*Handler, error) {
	var handlerOptions handlerOptions
	for _, o := range opts {
		o.apply(&handlerOptions)
	}

	maker, err := token.NewPasetoMaker(handlerOptions.conf.Token.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	handlerOptions.tokenMaker = maker

	handler := &Handler{
		handlerOptions: handlerOptions,
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

type HandlerOption interface {
	apply(*handlerOptions)
}

type funchandlerOption struct {
	f func(*handlerOptions)
}

func (fdo *funchandlerOption) apply(do *handlerOptions) {
	fdo.f(do)
}

func newFuncHandlerOption(f func(*handlerOptions)) *funchandlerOption {
	return &funchandlerOption{
		f: f,
	}
}

func InitiaQueries(queries *db.Queries) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.queries = queries
	})
}

func InitiaToken(tokenMaker token.Maker) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.tokenMaker = tokenMaker
	})
}

func InitiaConfig(conf config.Config) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.conf = conf
	})
}
