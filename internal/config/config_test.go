package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	var cm ConfigManager
	conf := cm.LoadConfig(".")

	t.Logf("config: %+v\n", conf)

	require.NotEmpty(t, conf)
}
