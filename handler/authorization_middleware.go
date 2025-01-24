package handler

import (
	"Dandelion/internal/db/model"
	"fmt"
	"net/http"
	"strings"

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

		if !strings.HasPrefix(authorization, authorizationPrefix) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header invalid"})
			return
		}

		fields := strings.Fields(authorization)
		payload, err := h.Token.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		var loginUserInfo model.LoginUserInfo
		err = h.Redis.Get(ctx, fmt.Sprintf("t_%s", payload.Username)).Scan(&loginUserInfo)
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
