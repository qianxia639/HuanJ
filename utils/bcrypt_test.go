package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestBcrypt(t *testing.T) {
	password := fmt.Sprintf("%06d", time.Now().Local().UnixNano())
	salt := GenerateSalt()

	hashPwd, err := HashPassword(password, salt)
	require.NoError(t, err)
	require.NotEmpty(t, hashPwd)

	err = ComparePassword(password, hashPwd, salt)
	require.NoError(t, err)

	wrongPwd := fmt.Sprintf("%06d", time.Now().Local().UnixNano())
	err = ComparePassword(wrongPwd, hashPwd, salt)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashPwd2, err := HashPassword(password, salt)
	require.NoError(t, err)
	require.NotEmpty(t, hashPwd2)
	require.NotEqual(t, hashPwd, hashPwd2)
}
