package handler

import (
	"Ice/internal/logs"
	"Ice/internal/utils"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html/atom"
)

func (h *Handler) secret() gin.HandlerFunc {
	return func(ctx *gin.Context) {

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

		// sign =Md5(m={Method}&p={Path}&q={Query}&a={authorization}&b={Body})
		sign := fmt.Sprintf("m=%s&p=%s&q=%s&a=%s&b=%s", method, path, query, authorization, body)

		signServer := utils.Md5(sign)
		signClient := ctx.Request.Header.Get("Sign")
		// 校验接口签名是否一致
		if signServer != signClient {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		ctx.Next()
	}
}
