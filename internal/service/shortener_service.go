package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/PrEvIeS/url_short/internal/repository"
)

type ShortenerService struct {
	repo repository.URLRepository
}

func NewShortenerService(repo repository.URLRepository) *ShortenerService {
	return &ShortenerService{repo: repo}
}

func (s *ShortenerService) CreateShortURL(originalURL string) (string, error) {
	shortID := generateShortID(8)
	log.Printf("Generated short ID: %s", shortID)

	err := s.repo.SaveURL(shortID, originalURL)
	if err != nil {
		log.Printf("Failed to save URL: %v", err)
		return "", fmt.Errorf("failed to save URL: %w", err)
	}

	log.Printf("Saved URL: %s with short ID: %s", originalURL, shortID)
	return shortID, nil
}

func (s *ShortenerService) GetOriginalURL(shortID string) (string, error) {
	log.Printf("Fetching original URL for short ID: %s", shortID)
	url, err := s.repo.GetURL(shortID)
	if err != nil {
		return "", fmt.Errorf("failed to get URL: %w", err)
	}
	return url, nil
}

func generateShortID(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}
