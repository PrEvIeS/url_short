package service

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/PrEvIeS/url_short/internal/app/repository"
)

type ShortenerService struct {
	repo repository.URLRepository
}

func NewShortenerService(repo repository.URLRepository) *ShortenerService {
	return &ShortenerService{repo: repo}
}

func (s *ShortenerService) CreateShortURL(originalURL string) (string, error) {
	shortID := generateShortID(8)
	err := s.repo.SaveURL(shortID, originalURL)
	if err != nil {
		return "", err
	}
	return shortID, nil
}

func (s *ShortenerService) GetOriginalURL(shortID string) (string, error) {
	return s.repo.GetURL(shortID)
}

func generateShortID(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}
