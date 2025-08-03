package handler

import (
	db "HuanJ/db/sqlc"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	authorizationHeader     = "authorization"
	authorizationPayloadKey = "authorization_payload"
	authorizationPrefix     = "Bearer "
)

func (h *Handler) authorizationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get(authorizationHeader)
		if len(authorization) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is not provided"})
			return
		}

		ua := ctx.Request.Header.Get("User-Agent")
		if len(ua) == 0 {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		payload, err := h.Token.VerifyToken(authorization)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		var loginUserInfo db.LoginUserInfo
		key := "user:" + payload.Username
		err = h.RedisClient.Get(ctx, key).Scan(&loginUserInfo)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			zap.L().Error("Redis error", zap.Error(err))
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// h.CurrentUserInfo = loginUserInfo
		ctx.Set("current_user_info", loginUserInfo)

		ctx.Next()
	}
}
