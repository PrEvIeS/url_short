package server

import (
	"github.com/PrEvIeS/url_short/internal/app/handler"
	"net/http"
)

type Server struct {
	handler *handler.ShortenerHandler
}

func NewServer(handler *handler.ShortenerHandler) *Server {
	return &Server{handler: handler}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}
