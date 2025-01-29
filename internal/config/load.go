package config

import (
	"Dandelion/internal/logs"
	"sync"

	"github.com/spf13/viper"
)

var instance *Config
var once sync.Once

func LoadConfig(path string) *Config {
	once.Do(func() {
		viper.AddConfigPath(path)
		viper.SetConfigName("config")
		viper.SetConfigType("toml")

		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			logs.Fatalf("Failed to read config file: %v", err)
		}

		var conf Config
		if err := viper.Unmarshal(&conf); err != nil {
			logs.Fatalf("Failed to unmarshal config data: %v", err)
		}

		instance = &conf
	})

	return instance
}
