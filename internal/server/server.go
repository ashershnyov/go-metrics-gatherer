package server

import (
	"net/http"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/handler"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/middleware"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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
func newConfig(opts ...option) *config {
	c := &config{
		address: defaultAddress,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type option func(*config)

// SetAddress sets custom address for an agent to send metrics to.
func SetAddress(addr *string) option {
	return func(c *config) {
		if addr != nil {
			c.address = *addr
		}
	}
}

// Server defines a server.
type Server struct {
	chi.Router
	cfg *config
}

// New creates a server using a provided cfg.
func New(s service.MetricStorage, opts ...option) *Server {
	cfg := newConfig(opts...)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugarLogger := logger.Sugar()

	h := handler.NewMetricsHandler(s)

	router := chi.NewRouter()
	router.Post("/update/*", middleware.Logging(sugarLogger, h.UpdateMetric()))
	router.Get("/value/*", middleware.Logging(sugarLogger, h.GetMetric()))
	return &Server{
		Router: router,
		cfg:    cfg,
	}
}

// ListenAndServe launches listening loop on the address provided in the config.
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.cfg.address, s)
}
