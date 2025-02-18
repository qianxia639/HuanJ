package config

import (
	"Ice/internal/logs"
	"sync"

	"github.com/spf13/viper"
)

type ConfigManager struct {
	instance *Config
	once     sync.Once
}

func (m *ConfigManager) LoadConfig(path string) Config {
	m.once.Do(func() {

		if err := setupViper(path); err != nil {
			logs.Fatalf("Failed to read config file: %v", err)
		}

		var conf Config
		if err := viper.Unmarshal(&conf); err != nil {
			logs.Fatalf("Failed to unmarshal config data: %v", err)
		}

		m.instance = &conf
	})

	return *m.instance
}

func setupViper(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	viper.AutomaticEnv()

	return viper.ReadInConfig()
}
