package handler

import (
	"bytes"
	"github.com/PrEvIeS/url_short/internal/app/repository"
	"github.com/PrEvIeS/url_short/internal/app/service"
	"github.com/PrEvIeS/url_short/internal/pkg/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePost(t *testing.T) {
	urlStorage := storage.NewInMemoryStorage()
	urlRepo := repository.NewURLRepository(urlStorage)
	shortenerService := service.NewShortenerService(urlRepo)

	handler := NewShortenerHandler(shortenerService)

	originalURL := "https://practicum.yandex.ru/"
	reqBody := bytes.NewBufferString(originalURL)
	req := httptest.NewRequest(http.MethodPost, "/", reqBody)
	rec := httptest.NewRecorder()

	handler.handlePost(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, rec.Code)
	}

	expectedShortURL := "http://localhost:8080/"
	if !bytes.HasPrefix(rec.Body.Bytes(), []byte(expectedShortURL)) {
		t.Errorf("Expected body to start with %s; got %s", expectedShortURL, rec.Body.String())
	}
}
func TestHandleGet(t *testing.T) {
	urlStorage := storage.NewInMemoryStorage()
	urlRepo := repository.NewURLRepository(urlStorage)

	shortenerService := service.NewShortenerService(urlRepo)

	shortID := "pO92GVXi"
	originalURL := "https://practicum.yandex.ru/"
	err := urlStorage.Set(shortID, originalURL)
	if err != nil {
		return
	}

	handler := NewShortenerHandler(shortenerService)

	req := httptest.NewRequest(http.MethodGet, "/"+shortID, nil)
	rec := httptest.NewRecorder()

	handler.handleGet(rec, req)

	// Проверка статуса ответа
	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status %d; got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	// Проверка заголовка Location
	location := rec.Header().Get("Location")
	if location != originalURL {
		t.Errorf("Expected Location header to be %s; got %s", originalURL, location)
	}
}
