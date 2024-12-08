package config

import (
	"github.com/spf13/pflag"
	"log"
)

type Config struct {
	ServerAddress string
	BaseURL       string
}

func NewConfig() *Config {
	config := &Config{}

	pflag.StringVarP(&config.ServerAddress, "address", "a", "http://localhost:8080", "Адрес запуска HTTP-сервера")
	pflag.StringVarP(&config.BaseURL, "base-url", "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL")

	pflag.Parse()

	log.Printf("server address: %s", config.ServerAddress)
	log.Printf("base url: %s", config.BaseURL)

	return config
}
