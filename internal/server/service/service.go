package service

import (
	"net/http"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/config"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/handler"
)

// Service defines a server.
type Service struct {
	*http.ServeMux
	cfg *config.Config
}

// New creates a server using a provided cfg.
func New(cfg *config.Config) *Service {
	mux := http.NewServeMux()

	h := handler.NewMetricUpdateHandler()

	mux.Handle("POST /update/", h)

	return &Service{
		ServeMux: mux,
		cfg:      cfg,
	}
}

// ListenAndServe launches listening loop on the address provided in the config.
func (s *Service) ListenAndServe() error {
	return http.ListenAndServe(s.cfg.Address, s)
}
