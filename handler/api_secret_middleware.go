package handler

import (
	"HuanJ/utils"
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) secret() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		reqTime := time.Now().UnixMilli()
		timeStr := ctx.Request.Header.Get("Timestamp")
		if timeStr == "" {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		timestamp, err := strconv.ParseInt(timeStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		timeDiff := reqTime - timestamp
		// 左右不得超过5分钟
		if timeDiff > 300000 || timeDiff < -300000 {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// 读取并恢复请求体
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer ctx.Request.Body.Close()
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // 恢复原始body

		// 签名要素
		// sign =Md5(m={Method}&p={Path}&q={Query}&t={Timestamp}&a={authorization}&b={Body})
		params := map[string]string{
			"m": ctx.Request.Method,
			"p": ctx.Request.URL.Path,
			"q": ctx.Request.URL.RawQuery,
			"t": timeStr,
			"a": ctx.Request.Header.Get(authorizationHeader),
			"b": string(body),
		}

		signStr := buildSignString(params)
		signServer := utils.Md5(signStr)
		signClient := ctx.Request.Header.Get("Sign")

		// 校验接口签名是否一致
		if signServer != signClient {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// 重放攻击防护
		// 判断签名是否已经被使用
		exists, err := h.RedisClient.SetNX(ctx, signClient, 1, 5*time.Minute).Result()
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// 如果exists为false，说明key存在，签名已被使用
		if !exists {
			zap.L().Warn("检测到重放攻击", zap.String("sign_client", signClient))
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		ctx.Next()
	}
}

// 构建签名字符串
func buildSignString(params map[string]string) string {
	var builder strings.Builder
	keys := []string{"m", "p", "q", "t", "a", "b"}
	for i, k := range keys {
		if i > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(k)
		builder.WriteByte('=')
		builder.WriteString(params[k])
	}
	return builder.String()
}
