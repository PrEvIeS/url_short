package handler

import (
	"bytes"
	"github.com/PrEvIeS/url_short/internal/app/repository"
	"github.com/PrEvIeS/url_short/internal/app/service"
	"github.com/PrEvIeS/url_short/internal/pkg/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePost(t *testing.T) {
	urlStorage := storage.NewInMemoryStorage()
	urlRepo := repository.NewURLRepository(urlStorage)
	shortenerService := service.NewShortenerService(urlRepo)

	handler := NewShortenerHandler(shortenerService, nil)

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
	}

	expectedShortURL := "http://" + req.Host + "/"
	if !bytes.HasPrefix(rec.Body.Bytes(), []byte(expectedShortURL)) {
		t.Errorf("Expected body to start with %s; got %s", expectedShortURL, rec.Body.String())
	}

	shortID := string(bytes.TrimPrefix(rec.Body.Bytes(), []byte(expectedShortURL)))
	storedURL, exists := urlStorage.Get(shortID)
	if !exists {
		t.Errorf("Expected URL to be stored in storage, but it was not found")
	}
	if storedURL != originalURL {
		t.Errorf("Expected stored URL to be %s; got %s", originalURL, storedURL)
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
		t.Fatalf("Failed to set URL in storage: %v", err)
	}

	handler := NewShortenerHandler(shortenerService, nil)

	req := httptest.NewRequest(http.MethodGet, "/"+shortID, nil)
	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "shortID", Value: shortID}}

	handler.HandleGet(c)

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status %d; got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != originalURL {
		t.Errorf("Expected Location header to be %s; got %s", originalURL, location)
	}
}
