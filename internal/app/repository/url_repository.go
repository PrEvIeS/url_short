package repository

import (
	"errors"
	"log"

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
	log.Printf("Saving URL: %s with short ID: %s", originalURL, shortID)
	return r.storage.Set(shortID, originalURL)
}

func (r *URLRepositoryImpl) GetURL(shortID string) (string, error) {
	log.Printf("Fetching URL for short ID: %s", shortID)
	url, exists := r.storage.Get(shortID)
	if !exists {
		log.Printf("URL not found for short ID: %s", shortID)
		return "", errors.New("URL not found")
	}
	log.Printf("Fetched URL: %s for short ID: %s", url, shortID)
	return url, nil
}
