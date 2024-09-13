package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server Server
	DB     DB
}

type Server struct {
	ServerAddress string `env:"SERVER_ADDRESS" env-required:"true"`
}

type DB struct {
	PostgresConn string `env:"POSTGRES_CONN" env-required:"true"`
}

func LoadEnv() *Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(fmt.Errorf("read config: %v", err))
	}

	return &cfg
}
