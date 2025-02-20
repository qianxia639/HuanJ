package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	var cm ConfigManager
	conf := cm.LoadConfig(".")

	require.NotEmpty(t, conf)
}
