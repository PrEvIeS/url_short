package repository

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type Storage interface {
	Set(key, value string) error
	Get(key string) (string, bool)
}
type URLRepository interface {
	SaveURL(shortID, originalURL string) error
	GetURL(shortID string) (string, error)
}

type URLRepositoryImpl struct {
	storage Storage
}

func NewURLRepository(srg Storage) *URLRepositoryImpl {
	return &URLRepositoryImpl{storage: srg}
}

func (r *URLRepositoryImpl) SaveURL(shortID, originalURL string) error {
	var mutex sync.Mutex
	mutex.Lock()
	log.Printf("Saving URL: %s with short ID: %s", originalURL, shortID)
	err := r.storage.Set(shortID, originalURL)
	if err != nil {
		return fmt.Errorf("saving URL: %w", err)
	}
	mutex.Unlock()
	return nil
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
