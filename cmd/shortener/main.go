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

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Error syncing logger: %v", err)
		}
	}()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Инициализация хранилища
	var urlStorage storage.Storage
	if cfg.FileStoragePath != "" {
		urlStorage, err = storage.NewFileStorage(cfg.FileStoragePath)
		if err != nil {
			logger.Fatal("Failed to initialize file storage",
				zap.String("path", cfg.FileStoragePath),
				zap.Error(err),
			)
		}
	} else {
		urlStorage = storage.NewInMemoryStorage()
		logger.Info("Using in-memory storage")
	}

	// Создание цепочки зависимостей
	urlRepo := repository.NewURLRepository(urlStorage, logger)
	shortenerService := service.NewShortenerService(urlRepo, logger)
	shortenerHandler := handler.NewShortenerHandler(shortenerService, cfg, logger)

	// Запуск сервера
	app := server.NewServer(shortenerHandler, cfg, logger)
	if err := app.Run(cfg.ServerAddress); err != nil {
		logger.Fatal("Server failed to start",
			zap.String("address", cfg.ServerAddress),
			zap.Error(err),
		)
	}
}
