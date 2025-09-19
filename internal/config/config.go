package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port         string `env:"PORT" env-default:"8080"`
	PostgresUser string `env:"POSTGRES_USER" env-required:"true"`
	PostgresPass string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresDB   string `env:"POSTGRES_DB" env-required:"true"`
	PostgresHost string `env:"POSTGRES_HOST" env-required:"true"`
	PostgresPort string `env:"POSTGRES_PORT" env-required:"true"`
	KafkaBrokers string `env:"KAFKA_BROKERS" env-default:"localhost:9092"`
	KafkaTopic   string `env:"KAFKA_TOPIC" env-default:"orders"`
	KafkaGroupID string `env:"KAFKA_GROUP_ID" env-default:"order-viewer"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	_ = godotenv.Load()

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) GetConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		c.PostgresUser,
		c.PostgresPass,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDB,
	)
}
