package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Env string

const (
	Dev   Env = "dev"
	Prod  Env = "prod"
	Local Env = "local"
)

type AppConfig struct {
	Env      Env    `env:"ENV" envDefault:"local"`
	HttpPort string `env:"HTTP_PORT" envDefault:"8080"`
	DBConn   string `env:"DB_CONN" envDefault:""`
	Secret   string `env:"SECRET"`
}

func NewConfig() (*AppConfig, error) {
	_ = godotenv.Load()

	c := &AppConfig{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}
