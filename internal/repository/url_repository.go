package repository

import (
	"errors"
	"fmt"

	"github.com/PrEvIeS/url_short/internal/storage"
	"go.uber.org/zap"
)

type URLRepository interface {
	SaveURL(shortID, originalURL string) error
	GetURL(shortID string) (string, error)
}

type URLRepositoryImpl struct {
	storage storage.Storage
	logger  *zap.Logger
}

// NewURLRepository создает новый репозиторий с указанным хранилищем и логгером.
func NewURLRepository(stg storage.Storage, logger *zap.Logger) *URLRepositoryImpl {
	return &URLRepositoryImpl{
		storage: stg,
		logger:  logger,
	}
}

func (r *URLRepositoryImpl) SaveURL(shortID, originalURL string) error {
	if err := r.storage.Set(shortID, originalURL); err != nil {
		return fmt.Errorf("failed to save URL: %w", err)
	}
	return nil
}

func (r *URLRepositoryImpl) GetURL(shortID string) (string, error) {
	url, exists := r.storage.Get(shortID)
	if !exists {
		r.logger.Warn("URL not found", zap.String("shortID", shortID))
		return "", ErrURLNotFound
	}
	return url, nil
}

var ErrURLNotFound = errors.New("URL not found")
