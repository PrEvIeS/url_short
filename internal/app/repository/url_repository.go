package repository

import (
	"errors"
	"github.com/PrEvIeS/url_short/internal/pkg/storage"
)

type URLRepository interface {
	SaveURL(shortID, originalURL string) error
	GetURL(shortID string) (string, error)
}

type URLRepositoryImpl struct {
	storage storage.Storage
}

func NewURLRepository(storage storage.Storage) *URLRepositoryImpl {
	return &URLRepositoryImpl{storage: storage}
}

func (r *URLRepositoryImpl) SaveURL(shortID, originalURL string) error {
	return r.storage.Set(shortID, originalURL)
}

func (r *URLRepositoryImpl) GetURL(shortID string) (string, error) {
	url, exists := r.storage.Get(shortID)
	if !exists {
		return "", errors.New("URL not found")
	}
	return url, nil
}
