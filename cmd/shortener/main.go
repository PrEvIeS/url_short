package main

import (
	"log"

	"github.com/PrEvIeS/url_short/internal/storage"

	"github.com/PrEvIeS/url_short/internal/config"
	"github.com/PrEvIeS/url_short/internal/handler"
	"github.com/PrEvIeS/url_short/internal/repository"
	"github.com/PrEvIeS/url_short/internal/server"
	"github.com/PrEvIeS/url_short/internal/service"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	urlStorage := storage.NewInMemoryStorage()

	urlRepo := repository.NewURLRepository(urlStorage)

	shortenerService := service.NewShortenerService(urlRepo)

	shortenerHandler := handler.NewShortenerHandler(shortenerService, cfg)

	app := server.NewServer(shortenerHandler, cfg)

	err = app.Run(cfg.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
