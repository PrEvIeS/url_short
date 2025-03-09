package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/PrEvIeS/url_short/internal/config"
	"github.com/PrEvIeS/url_short/internal/handler"
	"github.com/PrEvIeS/url_short/internal/repository"
	"github.com/PrEvIeS/url_short/internal/server"
	"github.com/PrEvIeS/url_short/internal/service"
	"github.com/PrEvIeS/url_short/internal/storage"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}
	defer func() {
		_ = logger.Sync() // Игнорируем ошибку Sync()
	}()
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	urlStorage := storage.NewInMemoryStorage()
	urlRepo := repository.NewURLRepository(urlStorage)
	shortenerService := service.NewShortenerService(urlRepo)

	// Передаем логгер в NewShortenerHandler
	shortenerHandler := handler.NewShortenerHandler(shortenerService, cfg, logger)

	app := server.NewServer(shortenerHandler, cfg, logger)

	err = app.Run(cfg.ServerAddress)
	if err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
	logger.Info("Application started successfully")
}
