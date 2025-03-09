package server

import (
	"fmt"

	"github.com/PrEvIeS/url_short/internal/config"
	"github.com/PrEvIeS/url_short/internal/handler"
	"github.com/PrEvIeS/url_short/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	handler *handler.ShortenerHandler
	config  *config.Config
	logger  *zap.Logger
}

func NewServer(hdl *handler.ShortenerHandler, cfg *config.Config, logger *zap.Logger) *Server {
	return &Server{
		handler: hdl,
		config:  cfg,
		logger:  logger,
	}
}

func (s *Server) Run(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Добавляем middleware Logger
	r.Use(middleware.Logger(s.logger))

	// Регистрируем обработчики
	r.POST("/", s.handler.HandlePost)
	r.GET("/:shortID", s.handler.HandleGet)
	r.POST("/api/shorten", s.handler.HandleJSONPost)
	// Запускаем сервер
	err := r.Run(addr)
	if err != nil {
		return fmt.Errorf("could not start server: %w", err)
	}
	return nil
}
