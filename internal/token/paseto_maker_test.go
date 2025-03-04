package token

import (
	"Rejuv/internal/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker := NewPasetoMaker(utils.RandomString(32))

	username := utils.RandomString(6)
	duration := time.Minute

	issueAt := time.Now()
	expired := issueAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issueAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expired, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker := NewPasetoMaker(utils.RandomString(32))

	token, err := maker.CreateToken(utils.RandomString(6), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
