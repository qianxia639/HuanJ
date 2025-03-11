package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	conf, err := LoadConfig(".")
	require.NoError(t, err)
	require.NotEmpty(t, conf)
}
