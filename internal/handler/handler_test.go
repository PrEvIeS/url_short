package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PrEvIeS/url_short/internal/config"
	"github.com/PrEvIeS/url_short/internal/repository"
	"github.com/PrEvIeS/url_short/internal/service"
	"github.com/PrEvIeS/url_short/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestHandlePost(t *testing.T) {
	// Инициализация конфигурации
	cfg := &config.Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
	}

	// Инициализация хранилища и сервиса
	urlStorage := storage.NewInMemoryStorage()
	urlRepo := repository.NewURLRepository(urlStorage)
	shortenerService := service.NewShortenerService(urlRepo)

	// Создаем логгер
	logger := zap.NewNop()

	// Передаем логгер в NewShortenerHandler
	handler := NewShortenerHandler(shortenerService, cfg, logger)

	originalURL := "http://dehoy.ru/n1ldm7e8bh88/gxn0xloupjkjol/veghgaewpnuop"
	reqBody := bytes.NewBufferString(originalURL)
	req := httptest.NewRequest(http.MethodPost, "/", reqBody)
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	handler.HandlePost(c)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, rec.Code)
		return
	}

	expectedShortURL := cfg.BaseURL + "/"
	if !bytes.HasPrefix(rec.Body.Bytes(), []byte(expectedShortURL)) {
		t.Errorf("Expected body to start with %s; got %s", expectedShortURL, rec.Body.String())
		return
	}

	shortID := string(bytes.TrimPrefix(rec.Body.Bytes(), []byte(expectedShortURL)))
	storedURL, exists := urlStorage.Get(shortID)
	if !exists {
		t.Errorf("Expected URL to be stored in storage, but it was not found")
		return
	}
	if storedURL != originalURL {
		t.Errorf("Expected stored URL to be %s; got %s", originalURL, storedURL)
		return
	}
}

func TestHandleGet(t *testing.T) {
	cfg := &config.Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
	}

	urlStorage := storage.NewInMemoryStorage()
	urlRepo := repository.NewURLRepository(urlStorage)
	shortenerService := service.NewShortenerService(urlRepo)

	shortID := "pO92GVXi"
	originalURL := "https://practicum.yandex.ru/"
	err := urlStorage.Set(shortID, originalURL)
	if err != nil {
		t.Errorf("Failed to set URL in storage: %v", err)
		return
	}

	// Создаем логгер
	logger := zap.NewNop()

	// Передаем логгер в NewShortenerHandler
	handler := NewShortenerHandler(shortenerService, cfg, logger)

	req := httptest.NewRequest(http.MethodGet, "/"+shortID, http.NoBody)
	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "shortID", Value: shortID}}

	handler.HandleGet(c)

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status %d; got %d", http.StatusTemporaryRedirect, rec.Code)
		return
	}

	location := rec.Header().Get("Location")
	if location != originalURL {
		t.Errorf("Expected Location header to be %s; got %s", originalURL, location)
		return
	}
}
