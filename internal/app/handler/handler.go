package handler

import (
	"bytes"
	"github.com/PrEvIeS/url_short/internal/app/config"
	"github.com/PrEvIeS/url_short/internal/app/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ShortenerHandler struct {
	service *service.ShortenerService
	config  *config.Config
}

func NewShortenerHandler(service *service.ShortenerService, config *config.Config) *ShortenerHandler {
	return &ShortenerHandler{service: service, config: config}
}

func (h *ShortenerHandler) HandlePost(c *gin.Context) {
	requestBody := c.Request.Body

	if requestBody == nil {
		c.String(http.StatusBadRequest, "URL is required")
		return
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(requestBody)
	if err != nil {
		return
	}
	originalURL := buf.String()

	shortID, err := h.service.CreateShortURL(originalURL)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create short URL")
		return
	}

	shortURL := h.config.BaseURL + "/" + shortID
	c.String(http.StatusCreated, shortURL)

	log.Printf("Created short URL: %s", shortID)
}

func (h *ShortenerHandler) HandleGet(c *gin.Context) {
	shortID := c.Param("shortID")

	originalURL, err := h.service.GetOriginalURL(shortID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short URL not found"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, originalURL)

	log.Printf("Expanded short URL: %s", shortID)
}
