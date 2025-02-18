package handler

import (
	db "Ice/db/sqlc"
	"fmt"
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
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Not User-Agent"})
			return
		}

		payload, err := h.Token.VerifyToken(authorization)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		var loginUserInfo db.LoginUserInfo
		err = h.Redis.Get(ctx, fmt.Sprintf("user:%s", payload.Username)).Scan(&loginUserInfo)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// TODO: ua不一致标识异设备请求，可能需要做处理
		// if ua != loginUserInfo.UserAgent {
		// 	ctx.Abort()
		// 	return
		// }

		h.CurrentUserInfo = loginUserInfo

		ctx.Next()
	}
}
