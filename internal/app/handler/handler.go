package handler

import (
	"github.com/PrEvIeS/url_short/internal/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ShortenerHandler struct {
	service *service.ShortenerService
}

func NewShortenerHandler(service *service.ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{service: service}
}

func (h *ShortenerHandler) HandlePost(c *gin.Context) {
	var requestBody string

	if err := c.ShouldBind(&requestBody); err != nil {
		c.String(http.StatusBadRequest, "URL is required")
		return
	}

	shortID, err := h.service.CreateShortURL(requestBody)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create short URL")
		return
	}

	shortURL := "http://" + c.Request.Host + "/" + shortID
	c.String(http.StatusCreated, shortURL)
}

func (h *ShortenerHandler) HandleGet(c *gin.Context) {
	shortID := c.Param("shortID")

	originalURL, err := h.service.GetOriginalURL(shortID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short URL not found"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, originalURL)
}
