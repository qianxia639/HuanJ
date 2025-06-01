package handler

import (
	db "HuanJ/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
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
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		h.CurrentUserInfo = loginUserInfo

		ctx.Next()
	}
}
