package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedisClient 模拟Redis客户端
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	args := m.Called(ctx, key, value, expiration)
	return args.Bool(0), args.Error(1)
}

// 测试签名构建函数
func TestBuildSignString(t *testing.T) {
	params := map[string]string{
		"m": "POST",
		"p": "/api/v1/users",
		"q": "page=1&limit=10",
		"t": "123456789",
		"a": "Bearer token123",
		"b": `{"name":"John"}`,
	}

	expected := "m=POST&p=/api/v1/users&q=page=1&limit=10&t=123456789&a=Bearer token123&b={\"name\":\"John\"}"
	result := buildSignString(params)

	assert.Equal(t, expected, result)
}

// 测试中间件
func TestSecretMiddlewareV2(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	// 创建模拟Redis客户端
	mockRedis := new(MockRedisClient)
	h := &Handler{
		RedisClient: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}

	// 当前时间戳（毫秒）
	now := time.Now().UnixMilli()
	currentTimestamp := strconv.FormatInt(now, 10)

	// 测试路由
	router := gin.New()
	router.Use(h.secret())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	// 生成有效签名
	generateSignature := func(timestamp, method, path, query, auth, body string) string {
		signStr := fmt.Sprintf("m=%s&p=%s&q=%s&t=%s&a=%s&b=%s",
			method, path, query, timestamp, auth, body)
		hash := md5.Sum([]byte(signStr))
		return hex.EncodeToString(hash[:])
	}

	t.Run("有效请求", func(t *testing.T) {
		// 设置Redis模拟 - SetNX返回true
		mockRedis.On("SetNX", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).Once()

		// 生成签名
		signature := generateSignature(currentTimestamp, "GET", "/test", "", "", "")

		// 创建请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Timestamp", currentTimestamp)
		req.Header.Set("Sign", signature)

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证
		assert.Equal(t, http.StatusOK, w.Code)
		mockRedis.AssertExpectations(t)
	})

	t.Run("缺少时间戳头", func(t *testing.T) {
		// 创建请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Sign", "dummy_signature")

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("无效时间戳格式", func(t *testing.T) {
		// 创建请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Timestamp", "invalid_timestamp")
		req.Header.Set("Sign", "dummy_signature")

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("时间戳过期", func(t *testing.T) {
		// 创建过期时间戳（6分钟前）
		expiredTime := now - 6*60*1000
		expiredTimestamp := strconv.FormatInt(expiredTime, 10)

		// 创建请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Timestamp", expiredTimestamp)
		req.Header.Set("Sign", "dummy_signature")

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("时间戳超前", func(t *testing.T) {
		// 创建未来时间戳（6分钟后）
		futureTime := now + 6*60*1000
		futureTimestamp := strconv.FormatInt(futureTime, 10)

		// 创建请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Timestamp", futureTimestamp)
		req.Header.Set("Sign", "dummy_signature")

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("签名不匹配", func(t *testing.T) {
		// 生成正确签名
		correctSignature := generateSignature(currentTimestamp, "GET", "/test", "", "", "")
		wrongSignature := correctSignature + "tampered"

		// 创建请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Timestamp", currentTimestamp)
		req.Header.Set("Sign", wrongSignature)

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("重放攻击", func(t *testing.T) {
		// 生成签名
		signature := generateSignature(currentTimestamp, "GET", "/test", "", "", "")

		// 设置Redis模拟 - SetNX返回false（签名已存在）
		mockRedis.On("SetNX", mock.Anything, signature, mock.Anything, mock.Anything).
			Return(false, nil).Once()

		// 创建请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Timestamp", currentTimestamp)
		req.Header.Set("Sign", signature)

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockRedis.AssertExpectations(t)
	})

	t.Run("Redis错误", func(t *testing.T) {
		// 生成签名
		signature := generateSignature(currentTimestamp, "GET", "/test", "", "", "")

		// 设置Redis模拟 - 返回错误
		mockRedis.On("SetNX", mock.Anything, signature, mock.Anything, mock.Anything).
			Return(false, fmt.Errorf("redis error")).Once()

		// 创建请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Timestamp", currentTimestamp)
		req.Header.Set("Sign", signature)

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockRedis.AssertExpectations(t)
	})

	t.Run("带请求体的POST请求", func(t *testing.T) {
		// 设置Redis模拟 - SetNX返回true
		mockRedis.On("SetNX", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).Once()

		// 请求体
		body := `{"name":"John","email":"john@example.com"}`
		bodyReader := strings.NewReader(body)

		// 生成签名
		signature := generateSignature(
			currentTimestamp,
			"POST",
			"/test",
			"",
			"Bearer token123",
			body,
		)

		// 创建请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", bodyReader)
		req.Header.Set("Timestamp", currentTimestamp)
		req.Header.Set("Sign", signature)
		req.Header.Set("Authorization", "Bearer token123")

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证
		assert.Equal(t, http.StatusOK, w.Code)
		mockRedis.AssertExpectations(t)
	})

	t.Run("带查询参数的请求", func(t *testing.T) {
		// 设置Redis模拟 - SetNX返回true
		mockRedis.On("SetNX", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).Once()

		// 生成签名
		signature := generateSignature(
			currentTimestamp,
			"GET",
			"/test",
			"page=2&limit=20",
			"",
			"",
		)

		// 创建请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test?page=2&limit=20", nil)
		req.Header.Set("Timestamp", currentTimestamp)
		req.Header.Set("Sign", signature)

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证
		assert.Equal(t, http.StatusOK, w.Code)
		mockRedis.AssertExpectations(t)
	})

	t.Run("请求体读取错误", func(t *testing.T) {
		// 创建一个会返回错误的特殊读取器
		errorReader := &errorReader{}

		// 创建请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", errorReader)
		req.Header.Set("Timestamp", currentTimestamp)
		req.Header.Set("Sign", "dummy_signature")

		// 执行请求
		router.ServeHTTP(w, req)

		// 验证
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// errorReader 用于模拟读取错误
type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated read error")
}
