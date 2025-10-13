package service

import (
	"net/http"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/config"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/handler"
)

// Server defines a server.
type Server struct {
	*http.ServeMux
	cfg *config.Config
}

// New creates a server using a provided cfg.
func New(cfg *config.Config) *Server {
	mux := http.NewServeMux()

	h := handler.NewMetricUpdateHandler()

	mux.Handle("POST /update/", h)

	return &Server{
		ServeMux: mux,
		cfg:      cfg,
	}
}

// ListenAndServe launches listening loop on the address provided in the config.
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.cfg.Address, s)
}
