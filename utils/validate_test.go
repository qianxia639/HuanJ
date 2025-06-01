package utils

import (
	"testing"
)

func TestValidate(t *testing.T) {
	t.Log(ValidateUsername(RandomString(31))) // F
	t.Log(ValidateUsername("12345"))          // T
	t.Log(ValidateUsername("%^&*))"))         // F
	t.Log(ValidateUsername("_qaz123"))        // T
	t.Log(ValidateUsername("test@email.com")) // T
}

func TestValidatePassword(t *testing.T) {
	t.Log(ValidatePassword("12345"))            // F
	t.Log(ValidatePassword("123456"))           // F
	t.Log(ValidatePassword("$123456"))          // F
	t.Log(ValidatePassword("_qaz123"))          // T
	t.Log(ValidatePassword("example@mail.com")) // T
}
