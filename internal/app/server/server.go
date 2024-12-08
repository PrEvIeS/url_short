package server

import (
	"github.com/PrEvIeS/url_short/internal/app/handler"
	"github.com/gin-gonic/gin"
)

type Server struct {
	handler *handler.ShortenerHandler
}

func NewServer(handler *handler.ShortenerHandler) *Server {
	return &Server{handler: handler}
}

func (s *Server) Run(addr string) {
	r := gin.Default()

	r.POST("/", s.handler.HandlePost)
	r.GET("/:shortID", s.handler.HandleGet)

	err := r.Run(addr)
	if err != nil {
		return
	}
}
