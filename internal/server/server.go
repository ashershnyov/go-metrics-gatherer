package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server/db"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/handler"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/middleware"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/service"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/storage"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"go.uber.org/zap"
)

const (
	// defaultAddress specifies the address used unless overridden by starting params.
	defaultAddress = ":8080"
	// defaultStoreInterval specifies the default interval to dump metrics to file.
	defaultStoreInterval = 300 * time.Second
	// defaultFilePath specifies the default path to file to dump metrics to.
	defaultFilePath = "metrics.json"
	// defaultRestoreMetrics specifies the default value of metrics restoration flag.
	defaultRestoreMetrics = true
	// defaultMaxRetries sets the default amount of retries upon send errors.
	defaultMaxRetries = 3
)

// config stores the server's configuration.
type config struct {
	address        string
	storeInterval  time.Duration
	filePath       string
	restoreMetrics bool
	dbAddress      string
	maxRetries     int
}

// newConfig constructs a config with default values, overrides with opts if passed.
func newConfig(opts ...option) *config {
	c := &config{
		address:        defaultAddress,
		storeInterval:  defaultStoreInterval,
		filePath:       defaultFilePath,
		restoreMetrics: defaultRestoreMetrics,
		maxRetries:     defaultMaxRetries,
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

// SetStoreInterval sets custom interval between metric dumps.
func SetStoreInterval(interval *int) option {
	return func(c *config) {
		if interval != nil {
			c.storeInterval = time.Duration(*interval) * time.Second
		}
	}
}

// SetFilePath sets custom path to the file to dump metrics to.
func SetFilePath(path *string) option {
	return func(c *config) {
		if path != nil {
			c.filePath = *path
		}
	}
}

// SetRestoreMetrics sets metric restoration from file flag.
func SetRestoreMetrics(flag *bool) option {
	return func(c *config) {
		if flag != nil {
			c.restoreMetrics = *flag
		}
	}
}

func SetDBAddress(address *string) option {
	return func(c *config) {
		if address != nil {
			c.dbAddress = *address
		}
	}
}

// Server defines a server.
type Server struct {
	router chi.Router
	cfg    *config
	db     *db.Postgres
}

// New creates a server using a provided cfg.
func New(logger *zap.SugaredLogger, opts ...option) (*Server, error) {
	cfg := newConfig(opts...)

	router := chi.NewRouter()
	router.Use(middleware.Logging(logger))
	router.Use(middleware.Gzip())

	var (
		pg  *db.Postgres
		err error
	)
	if cfg.dbAddress != "" {
		pg, err = db.NewPostgres(context.Background(), cfg.dbAddress, cfg.maxRetries)
		if err != nil {
			return nil, fmt.Errorf("an error occured when opening DB: %w", err)
		}
	}

	return &Server{
		router: router,
		cfg:    cfg,
		db:     pg,
	}, nil
}

// ListenAndServe launches listening loop on the address provided in the config.
func (s *Server) ListenAndServe() {
	if err := http.ListenAndServe(s.cfg.address, s.router); err != nil {
		log.Fatal(err)
	}
}

// Run executes server's loop.
func (s *Server) Run() error {
	var err error

	var stg service.MetricStorage
	if s.cfg.dbAddress != "" {
		stg = storage.NewDB(s.db)
		err = goose.Up(s.db.SQLDB(), "./migrations")
		if err != nil {
			return fmt.Errorf("an error occurred when starting Server: %w", err)
		}
		defer goose.Down(s.db.SQLDB(), "./migrations")
		defer s.db.Close()
	} else {
		stg = storage.NewInMemory()
	}

	service := service.NewService(stg)

	var d *middleware.MetricDumper
	if s.cfg.filePath != "" {
		d, err = middleware.NewMetricDumper(service, s.cfg.storeInterval, s.cfg.filePath, s.cfg.restoreMetrics)
		if err != nil {
			return fmt.Errorf("an error occurred when starting Server: %w", err)
		}
		defer d.Dump()
	}

	h := handler.NewMetricsHandler(service, s.db)

	s.router.Get("/", h.ListMetrics())
	s.router.Post("/update/", d.Middleware(h.UpdateMetricJSON()))
	s.router.Post("/update/*", d.Middleware(h.UpdateMetric()))
	s.router.Post("/value/", h.GetMetricJSON())
	s.router.Get("/value/*", h.GetMetric())
	s.router.Get("/ping", h.PingDB())
	s.router.Post("/updates/", d.Middleware(h.UpdateMultipleJSON()))

	d.DumperLoop(service)
	go s.ListenAndServe()

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGTERM, syscall.SIGINT)
	<-term

	return nil
}
