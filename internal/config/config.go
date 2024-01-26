package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env         string `env:"ENV"`
	ServerAddr  string `env:"SERVER_ADDR"`
	DatabaseUrl string `env:"DATABASE_URL"`
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
