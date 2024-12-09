package config

import (
	"log"

	"github.com/caarlos0/env/v11"

	"github.com/spf13/pflag"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func NewConfig() *Config {
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

	if config.ServerAddress == "" {
		config.ServerAddress = "localhost:8080"
	}
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:8080"
	}

	pflag.Parse()

	if err := env.Parse(config); err != nil {
		log.Println(err)
	}

	log.Printf("server address: %s", config.ServerAddress)
	log.Printf("base url: %s", config.BaseURL)

	return config
}
