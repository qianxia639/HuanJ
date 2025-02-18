package handler

import (
	db "Ice/db/sqlc"
	"Ice/internal/config"
	"Ice/internal/utils"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
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

	h := NewHandler(conf, store, nil)
	require.Equal(t, 1, 1)

	return h
}
