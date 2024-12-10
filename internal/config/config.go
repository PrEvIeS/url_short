package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/pflag"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func NewConfig() (*Config, error) {
	config := &Config{}

	pflag.StringVarP(
		&config.ServerAddress,
		"address",
		"a",
		"localhost:8080",
		"Адрес запуска HTTP-сервера",
	)
	pflag.StringVarP(
		&config.BaseURL,
		"base-url",
		"b",
		"http://localhost:8080",
		"Базовый адрес результирующего сокращённого URL",
	)

	pflag.Parse()

	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("error while parsing environment variables: %w", err)
	}

	log.Printf("server address: %s", config.ServerAddress)
	log.Printf("base url: %s", config.BaseURL)

	return config, nil
}
