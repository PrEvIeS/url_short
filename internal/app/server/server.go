package server

import (
	"github.com/PrEvIeS/url_short/internal/app/config"
	"github.com/PrEvIeS/url_short/internal/app/handler"
	"github.com/gin-gonic/gin"
)

type Server struct {
	handler *handler.ShortenerHandler
	config  *config.Config
}

func NewServer(handler *handler.ShortenerHandler, cfg *config.Config) *Server {
	return &Server{handler: handler, config: cfg}
}

func (s *Server) Run(addr string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/", s.handler.HandlePost)
	r.GET("/:shortID", s.handler.HandleGet)

	err := r.Run(addr)
	if err != nil {
		panic(err)
	}
}
