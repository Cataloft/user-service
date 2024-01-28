package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `env:"ENV"`
	ServerAddr string `env:"SERVER_ADDR"`
	Database   `env:"DATABASE"`
}

type Database struct {
	DatabaseURL      string        `env:"DATABASE_URL"`
	MaxAttempts      int           `env:"MAX_ATTEMPTS"`
	DurationAttempts time.Duration `env:"DURATION_ATTEMPTS"`
}

func MustLoad() *Config {
	configPath := os.Getenv("config_path")
	if configPath == "" {
		log.Fatal("config path is not set")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read cfg")
	}

	return &cfg
}
