package utils

import "testing"

func TestValidate(t *testing.T) {
	t.Log(ValidateUsername(RandomString(21))) // F
	t.Log(ValidateUsername("12345"))          // T
	t.Log(ValidateUsername("%^&*))"))         // F
	t.Log(ValidateUsername("_qaz123"))        // T
	t.Log(ValidateUsername("test@email.com")) // T
}
