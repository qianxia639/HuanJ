package utils

import (
	"crypto/md5"
	"fmt"
)

func Md5(str string) string {
	hash := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", hash)
}
