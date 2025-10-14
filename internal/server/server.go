package server

import (
	"net/http"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/handler"
)

const (
	// defaultAddress specifies the address used unless overridden by starting params.
	defaultAddress = ":8080"
)

// config stores the server's configuration.
type config struct {
	address string
}

// New constructs a config with default values, overrides with opts if passed.
func newConfig() *config {
	return &config{
		address: defaultAddress,
	}
}

// Server defines a server.
type Server struct {
	*http.ServeMux
	cfg *config
}

// New creates a server using a provided cfg.
func New() *Server {
	cfg := newConfig()
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
	return http.ListenAndServe(s.cfg.address, s)
}
