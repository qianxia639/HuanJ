package utils

import (
	"testing"
)

func TestMaskEmail(t *testing.T) {
	testCases := []string{
		"we@example.com",
		"e@example.com",
		"tes@example.com",
		"test@example.com",
		"example@example.com",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			t.Log(MaskEmail(tc))
		})
	}
}

func TestMaskPhone(t *testing.T) {
	testCases := []string{
		"15500001111",
		"8734261",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			t.Log(MaskPhone(tc))
		})
	}
}
