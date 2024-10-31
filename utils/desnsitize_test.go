package utils

import (
	"testing"
)

func TestDesnsitize(t *testing.T) {
	testCases := []string{
		"we@example.com",
		"e@example.com",
		"tes@example.com",
		"test@example.com",
		"example@example.com",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			t.Log(DesnsitizeEmail(tc))
		})
	}
}
