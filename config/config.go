package config

import (
	"fmt"
	"time"
)

type Config struct {
	Http     Http     `mapstructure:"http"`
	Postgres Postgres `mapstructure:"postgres"`
	Token    Token    `mapstructure:"token"`
	Redis    Redis    `mapstructure:"redis"`
	Secret   Secret   `mapstructure:"secret"`
}

type Http struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func (h *Http) Address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
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
	TokenSymmetricKey    string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

type Redis struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func (r *Redis) Address() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type Secret struct {
	Enable bool `mapstructure:"enable"`
}
