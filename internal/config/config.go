package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `env:"ENV"`
	ServerAddr  string `env:"SERVER_ADDR"`
	DatabaseURL string `env:"DATABASE_URL"`
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
