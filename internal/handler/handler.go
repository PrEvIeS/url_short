package handler

import (
	"bytes"
	"encoding/json"
	"github.com/PrEvIeS/url_short/internal/model"
	"net/http"

	"github.com/PrEvIeS/url_short/internal/config"
	"github.com/PrEvIeS/url_short/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	shortIDKey = "shortID" // Константа для ключа "shortID"
)

type ShortenerHandler struct {
	service *service.ShortenerService
	config  *config.Config
	logger  *zap.Logger
}

func NewShortenerHandler(
	shortenerService *service.ShortenerService,
	cfg *config.Config, logger *zap.Logger) *ShortenerHandler {
	return &ShortenerHandler{
		service: shortenerService,
		config:  cfg,
		logger:  logger,
	}
}

func (h *ShortenerHandler) HandlePost(c *gin.Context) {
	requestBody := c.Request.Body

	if requestBody == nil {
		h.logger.Error("Request body is nil")
		c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(requestBody)
	if err != nil {
		h.logger.Error("Failed to read request body", zap.Error(err))
		c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	originalURL := buf.String()

	shortID, err := h.service.CreateShortURL(originalURL)
	if err != nil {
		h.logger.Error("Failed to create short URL", zap.Error(err))
		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	shortURL := h.config.BaseURL + "/" + shortID
	c.String(http.StatusCreated, shortURL)

	h.logger.Info("Created short URL", zap.String(shortIDKey, shortID)) // Используем константу
}

func (h *ShortenerHandler) HandleGet(c *gin.Context) {
	shortID := c.Param(shortIDKey) // Используем константу

	originalURL, err := h.service.GetOriginalURL(shortID)
	if err != nil {
		h.logger.Error("Failed to get original URL", zap.String(shortIDKey, shortID), zap.Error(err)) // Используем константу
		c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, originalURL)

	h.logger.Info("Expanded short URL", zap.String(shortIDKey, shortID)) // Используем константу
}
func (h ShortenerHandler) HandleJSONPost(c *gin.Context) {
	if c.Request.Body == nil {
		h.logger.Error("Request body is nil")
		c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	var request model.CreateShortUrl
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to decode JSON", zap.Error(err))
		c.String(http.StatusBadRequest, "Invalid JSON format")
		return
	}

	originalURL := request.Url

	shortID, err := h.service.CreateShortURL(originalURL)
	if err != nil {
		h.logger.Error("Failed to create short URL", zap.Error(err))
		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	shortURL := h.config.BaseURL + "/" + shortID
	result := model.UrlResponse{Result: shortURL}

	c.JSON(http.StatusCreated, result)

	h.logger.Info("Created short URL", zap.String("short_id", shortID))
}
