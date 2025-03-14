package handler

import (
	"Rejuv/config"
	db "Rejuv/db/sqlc"
	"Rejuv/logs"
	"Rejuv/token"
	"crypto/ed25519"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

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

	// maker := token.NewPasetoMaker(conf.Token.TokenSymmetricKey)

	privateKey, publicKey, err := parseKeypair("../token/private_key.pem", "../token/public_key.pem")
	if err != nil {
		logs.Error(err)
		return nil
	}

	maker := token.NewPasetoMakerV2(privateKey, publicKey)

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
		authRouter.POST("/friend/request/process", handler.processFriendRequest)
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
