package config

import (
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQ
	Postgres
}

type RabbitMQ struct {
	Host     string `env:"RABBITMQ_HOST,required"`
	Port     string `env:"RABBITMQ_PORT"                  envDefault:"5672"`
	User     string `env:"RABBITMQ_DEFAULT_USER,required"`
	Password string `env:"RABBITMQ_DEFAULT_PASS,required"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Database string `env:"POSTGRES_DB,required"`
	SSLMode  string `env:"POSTGRES_SSL_MODE"          envDefault:"disable"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var cfg Config

	cfg, err = env.ParseAs[Config]()

	return &cfg, err
}
