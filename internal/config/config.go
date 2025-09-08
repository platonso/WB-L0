package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port    string `env:"PORT" env-default:":8080"`
	ConnStr string `env:"POSTGRES_URL" env-required:"true"`
	//	KafkaBrokers string
	//	KafkaTopic   string
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
