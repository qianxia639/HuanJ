package handler

import (
	db "Rejuv/db/sqlc"
	"Rejuv/internal/config"
	"Rejuv/internal/utils"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func newTestHandler(t *testing.T, store db.Store) *Handler {
	conf := config.Config{
		Token: config.Token{
			TokenSymmetricKey:   utils.RandomString(32),
			AccessTokenDuration: time.Minute,
		},
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	h := NewHandler(conf, store, rdb)

	return h
}
