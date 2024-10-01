package config

import (
	"fmt"
	"time"
)

type Config struct {
	Postgres Postgres
	Token    Token
}

type Postgres struct {
	Port         int    `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"db_name"`
	SslMode      string `mapstructure:"ssl_mode"`
	Driver       string `mapstructure:"driver"`
	MigrationUrl string `mapstructure:"migrate_url"`
}

func (p *Postgres) DatabaseSource() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.Username, p.Password, p.DbName, p.SslMode)
}

func (p *Postgres) DatabaseUrl() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		p.Driver, p.Username, p.Password, p.Host, p.Port, p.DbName, p.SslMode)
}

type Token struct {
	TokenSymmetricKey   string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration time.Duration `mapstructure:"access_token_duration"`
}
