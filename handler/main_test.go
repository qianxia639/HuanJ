package handler

import (
	"Dandelion/config"
	db "Dandelion/db/service"
	"Dandelion/utils"
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

func newTestHandler(t *testing.T, queries *db.Queries) *Handler {
	conf := config.Config{
		Token: config.Token{
			TokenSymmetricKey:   utils.RandomString(32),
			AccessTokenDuration: time.Minute,
		},
	}

	h, err := NewHandler(conf, queries)
	require.NoError(t, err)

	return h
}
