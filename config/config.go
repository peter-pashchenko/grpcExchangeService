package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Log
		PG
		GRPC
		HTTP
	}
	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL" `
	}
	PG struct {
		Host     string `env:"PG_HOST" env-required:"true" `
		Port     string `env:"PG_PORT" env-required:"true" `
		User     string `env:"PG_USER" env-required:"true" `
		Pass     string `env:"PG_PASSWORD" env-required:"true" `
		Database string `env:"PG_DATABASE" env-required:"true"`
	}
	GRPC struct {
		Port string `env:"GRPC_PORT" env-required:"true" `
	}
	HTTP struct {
		Host string `env:"HTTP_HOST" env-required:"true" `
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
