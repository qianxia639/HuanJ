package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	conf := LoadConfig(".")

	require.NotEmpty(t, conf)
}

func TestSingletonPattern(t *testing.T) {
	firstConf := LoadConfig(".")

	secondConf := LoadConfig(".")

	require.Equal(t, firstConf, secondConf)
}
