package handler

import (
	"HuanJ/config"
	db "HuanJ/db/sqlc"
	"HuanJ/logs"
	"HuanJ/token"
	"HuanJ/ws"
	"crypto/ed25519"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

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
	RedisClient     *redis.Client
}

func NewHandler(conf config.Config, store db.Store, rdb *redis.Client) *Handler {

	maker := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)

	// privateKey, publicKey, err := parseKeypair("../token/private_key.pem", "../token/public_key.pem")
	// if err != nil {
	// 	logs.Error(err)
	// 	return nil
	// }

	// maker := token.NewPasetoMakerV2(privateKey, publicKey)

	handler := &Handler{
		Conf:        conf,
		Store:       store,
		Token:       maker,
		RedisClient: rdb,
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

	connManager := ws.NewConnectionManager()
	go connManager.Run()

	if handler.Conf.Secret.Enable {
		router.Use(handler.secret())
	}

	router.GET("/ws", handler.wsHandler)
	router.GET("/wss", func(ctx *gin.Context) {
		handler.wssHandler(connManager, ctx.Writer, ctx.Request)
	})

	authRouter := router.Group("")
	authRouter.Use(handler.authorizationMiddleware())

	// Router
	{
		router.POST("/refresh")
		router.PUT("/reset/pwd", handler.resetPwd)
		router.POST("/email/code", handler.sendEmail)
	}

	// User Router
	{
		router.POST("/login", handler.login)
		authRouter.POST("/logout", handler.logout)
		router.POST("/user", handler.createUser)

		authRouter.GET("/user", handler.getUser)
		authRouter.PUT("/user", handler.updateUser)
		authRouter.PUT("/user/pwd", handler.updatePassword)
	}

	// Friend Request Router
	{
		authRouter.POST("/friend/request", handler.sendFriendRequest)
		authRouter.POST("/friend/request/process", handler.processFriendRequest)
		authRouter.POST("/friend/request/list", handler.listFriendRequest)
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

func LogFuncExecTime() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request

		logs.Infof("HttpUrl: %s|%s Started,IP:%s", req.Method, req.URL, ctx.RemoteIP())
		req.URL.Path = strings.ToLower(req.URL.Path)
		t := time.Now()
		ctx.Next()
		logs.Infof("HttpUrl: %s|%s Finish,Execute Time %5dms", req.Method, req.URL, time.Now().Sub(t).Milliseconds())
	}
}

func parseKeypair(privateKeyPath, publicKeyPath string) (ed25519.PrivateKey, ed25519.PublicKey, error) {

	// 读取私钥
	privateKeyPem, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, err
	}

	// 解析PEM格式的私钥
	privateKeyBlock, _ := pem.Decode(privateKeyPem)
	if privateKeyBlock == nil || privateKeyBlock.Type != "PRIVATE KEY" {
		return nil, nil, fmt.Errorf("无效的PEM文件")
	}

	privateKey := ed25519.PrivateKey(privateKeyBlock.Bytes)

	// 读取私钥
	publicKeyPem, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, nil, err
	}

	// 解析PEM格式的私钥
	publicKeyBlock, _ := pem.Decode(publicKeyPem)
	if publicKeyBlock == nil || publicKeyBlock.Type != "PUBLIC KEY" {
		return nil, nil, fmt.Errorf("无效的PEM文件")
	}

	publicKey := ed25519.PublicKey(publicKeyBlock.Bytes)

	return privateKey, publicKey, nil
}
