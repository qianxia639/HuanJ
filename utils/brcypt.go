package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateSalt() string {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		panic("generate salt failed: " + err.Error())
	}

	return hex.EncodeToString(salt)
}

func passwordFmt(password, salt string) string {
	return password + "." + salt
}

func HashPassword(password, salt string) (string, error) {
	password = passwordFmt(password, salt)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("password to hash password faild: %w", err)
	}

	return string(hashPassword), nil
}

func ComparePassword(password, hashedPassword, salt string) error {
	password = passwordFmt(password, salt)
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
