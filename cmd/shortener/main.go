package main

import (
	"fmt"
	"github.com/PrEvIeS/url_short/internal/app/config"
	"github.com/PrEvIeS/url_short/internal/app/handler"
	"github.com/PrEvIeS/url_short/internal/app/repository"
	"github.com/PrEvIeS/url_short/internal/app/server"
	"github.com/PrEvIeS/url_short/internal/app/service"
	"github.com/PrEvIeS/url_short/internal/pkg/storage"
)

func main() {
	cfg := config.NewConfig()

	urlStorage := storage.NewInMemoryStorage()

	urlRepo := repository.NewURLRepository(urlStorage)

	shortenerService := service.NewShortenerService(urlRepo)

	shortenerHandler := handler.NewShortenerHandler(shortenerService, cfg)

	app := server.NewServer(shortenerHandler, cfg)

	fmt.Println("Starting server on :8080")
	app.Run(cfg.ServerAddress)
}
