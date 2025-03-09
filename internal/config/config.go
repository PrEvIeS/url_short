package config

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/pflag"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() (*Config, error) {
	config := &Config{}

	// Флаги
	pflag.StringVarP(&config.ServerAddress, "address", "a", "localhost:8080", "Server address")
	pflag.StringVarP(&config.BaseURL, "base-url", "b", "http://localhost:8080", "Base URL")
	pflag.StringVarP(&config.FileStoragePath, "file", "f", "urls.json", "File storage path")

	pflag.Parse()

	// Переменные окружения
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("error parsing env vars: %w", err)
	}

	// Приоритет переменных окружения
	if path := os.Getenv("FILE_STORAGE_PATH"); path != "" {
		config.FileStoragePath = path
	}

	log.Printf("Using storage file: %s", config.FileStoragePath)
	return config, nil
}
