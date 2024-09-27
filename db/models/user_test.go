package models

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestUserSieOf(t *testing.T) {
	var user User
	user.ID = 1
	user.Username = "djkfhg"

	fmt.Printf("unsafe.Sizeof(&user): %v\n", unsafe.Sizeof(&user.ID))
}
