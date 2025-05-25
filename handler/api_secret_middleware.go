package handler

import (
	"HuanJ/logs"
	"HuanJ/utils"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html/atom"
)

func (h *Handler) secret() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		reqTime := time.Now().UnixMilli()
		timeStr := ctx.Request.Header.Get("Timestamp")
		if timeStr == "" {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		logs.Infof("timeStr: %s\n", timeStr)

		t, err := strconv.ParseInt(timeStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		timeSpan := reqTime - t
		// 左右不得超过5分钟
		if timeSpan > 300000 || timeSpan < -300000 {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		authorization := ctx.Request.Header.Get(authorizationHeader)

		logs.Infof("method: %s, path: %s, query: %s\n", method, path, query)
		logs.Infof("authorization: %s\n", authorization)

		var body string
		out, err := io.ReadAll(ctx.Request.Body)
		if err == nil {
			body = atom.String(out)
			// 重新生成流
			_ = ctx.Request.Body.Close() // 必须关闭
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(out))
		}

		logs.Infof("body: %v\n", body)

		// sign =Md5(m={Method}&p={Path}&q={Query}&t={Timestamp}&a={authorization}&b={Body})
		signStr := fmt.Sprintf("m=%s&p=%s&q=%s&t=%s&a=%s&b=%s",
			method, path, query, timeStr, authorization, body)

		signServer := utils.Md5(signStr)
		signClient := ctx.Request.Header.Get("Sign")
		logs.Info(signClient, signServer)
		logs.Info(signStr)
		// 校验接口签名是否一致
		if signServer != signClient {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// 判断签名是否已经被使用
		exists, err := h.RedisClient.SetNX(ctx, signClient, 1, 10*time.Minute).Result()
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// 如果exists为false，说明key存在，签名已被使用
		if !exists {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		ctx.Next()
	}
}
