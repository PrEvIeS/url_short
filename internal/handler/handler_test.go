package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/PrEvIeS/url_short/internal/config"
	"github.com/PrEvIeS/url_short/internal/repository"
	"github.com/PrEvIeS/url_short/internal/service"
	"github.com/PrEvIeS/url_short/internal/storage"
)

func setupTestHandler() (*ShortenerHandler, *gin.Engine) {
	logger := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://test"}
	store := storage.NewInMemoryStorage()
	repo := repository.NewURLRepository(store, logger)
	svc := service.NewShortenerService(repo, logger)
	handler := NewShortenerHandler(svc, cfg, logger)

	router := gin.Default()
	router.POST("/", handler.HandlePost)
	router.GET("/:shortID", handler.HandleGet)
	router.POST("/api/shorten", handler.HandleJSONPost)

	return handler, router
}

func TestHandlePost(t *testing.T) {
	_, router := setupTestHandler()

	t.Run("Valid URL", func(t *testing.T) {
		body := bytes.NewBufferString("https://example.com")
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", "text/plain")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "http://test/")
	})

	t.Run("Empty Body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandleJSONPost(t *testing.T) {
	_, router := setupTestHandler()

	t.Run("Valid JSON", func(t *testing.T) {
		body := bytes.NewBufferString(`{"url":"https://json.example.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, `{"result":"http://test/edVPg3ks"}`, rec.Body.String())
	})
}

func TestHandleGet(t *testing.T) {
	handler, router := setupTestHandler()

	// Сначала создаем тестовую запись
	_, err := handler.service.CreateShortURL("https://redirect.example.com")
	if err != nil {
		return
	}

	t.Run("Valid Short URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/edVPg3ks", http.NoBody)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
		assert.Equal(t, "https://redirect.example.com", rec.Header().Get("Location"))
	})

	t.Run("Invalid Short URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/invalid", http.NoBody)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
