package server

import (
	"fmt"

	"github.com/PrEvIeS/url_short/internal/config"
	"github.com/PrEvIeS/url_short/internal/handler"
	"github.com/gin-gonic/gin"
)

type Server struct {
	handler *handler.ShortenerHandler
	config  *config.Config
}

func NewServer(hdl *handler.ShortenerHandler, cfg *config.Config) *Server {
	return &Server{handler: hdl, config: cfg}
}

func (s *Server) Run(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/", s.handler.HandlePost)
	r.GET("/:shortID", s.handler.HandleGet)

	err := r.Run(addr)
	if err != nil {
		return fmt.Errorf("could not start server: %w", err)
	}
	return nil
}
