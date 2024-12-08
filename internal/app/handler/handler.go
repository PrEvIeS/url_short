package handler

import (
	"fmt"
	"github.com/PrEvIeS/url_short/internal/app/service"
	"io"
	"net/http"
	"strings"
)

type ShortenerHandler struct {
	service *service.ShortenerService
}

func NewShortenerHandler(service *service.ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{service: service}
}

func (h *ShortenerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodGet:
		h.handleGet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ShortenerHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	originalURL := strings.TrimSpace(string(body))
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortID, err := h.service.CreateShortURL(originalURL)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortID)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		return
	}
}

func (h *ShortenerHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	shortID := strings.TrimPrefix(r.URL.Path, "/")

	originalURL, err := h.service.GetOriginalURL(shortID)
	if err != nil {
		http.Error(w, "Short URL not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
