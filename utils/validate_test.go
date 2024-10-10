package utils

import "testing"

func TestValidate(t *testing.T) {
	t.Log(ValidateUsername(RandomString(21)))
	t.Log(ValidateUsername("12345"))
	t.Log(ValidateUsername("%^&*))"))
	t.Log(ValidateUsername("_qaz123"))
}
