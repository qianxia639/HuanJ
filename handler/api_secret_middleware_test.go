package handler

import (
	"HuanJ/utils"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSecretMiddleware(t *testing.T) {

	t.Run("OK", func(t *testing.T) {
		h := newTestHandler(t, nil)
		secretPath := "/test-secret"
		h.Router.POST(secretPath,
			h.secret(),
			h.authorizationMiddleware(),
			func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			},
		)

		body := `{"key": "value"}`

		req := httptest.NewRequest(http.MethodPost, "/test-secret", bytes.NewBufferString(body))
		timeStr := strconv.FormatInt(time.Now().UnixMilli(), 10)
		req.Header.Set("Timestamp", timeStr)

		signStr := fmt.Sprintf("m=%s&p=%s&q=%s&t=%s&a=%s&b=%s",
			req.Method, req.URL.Path, req.URL.RawQuery, timeStr, "", body)

		sign := utils.Md5(signStr)

		req.Header.Set("Sign", sign)

		err := h.RedisClient.Set(context.Background(), sign, 1, 0).Err()
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		h.Router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
	})
}
