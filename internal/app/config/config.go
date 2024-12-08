package config

import (
	"github.com/spf13/pflag"
	"log"
)

type Config struct {
	ServerAddress string
	BaseUrl       string
}

func NewConfig() *Config {
	config := &Config{}

	pflag.StringVarP(&config.ServerAddress, "address", "a", "localhost:8080", "Адрес запуска HTTP-сервера")
	pflag.StringVarP(&config.BaseUrl, "base-url", "b", "localhost:8080", "Адрес запуска HTTP-сервера")

	pflag.Parse()

	log.Printf("server address: %s", config.ServerAddress)
	log.Printf("base url: %s", config.BaseUrl)

	return config
}
