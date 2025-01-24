package utils

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrSaltLenght = errors.New("salt lenght must gt 0 byte")

func GenerateSalt() string {
	u, err := uuid.NewRandom()
	if err != nil {
		log.Printf("generate salt failed: %v", err)
		return ""
	}

	return strings.ReplaceAll(u.String(), "-", "")
}

func passwordFmt(password, salt string) string {
	return password + "." + salt
}

func HashPassword(password, salt string) (string, error) {
	if len(salt) == 0 {
		return "", ErrSaltLenght
	}

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
