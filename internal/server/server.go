package server

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashershnyov/go-metrics-gatherer/internal/buildinfo"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/audit"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/config"
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

// Server defines a server.
type Server struct {
	buildinfo buildinfo.BuildInfo
	router    chi.Router
	cfg       *config.Config
	db        *db.Postgres
}

// New creates a server using a provided cfg.
func New(bi buildinfo.BuildInfo, logger *zap.SugaredLogger, opts ...config.Option) (*Server, error) {
	cfg := config.NewConfig(opts...)

	router := chi.NewRouter()
	router.Use(
		middleware.Logging(logger),
		middleware.Gzip(),
		middleware.Hashing(cfg.Key),
	)

	var (
		pg  *db.Postgres
		err error
	)
	if cfg.DBAddress != "" {
		pg, err = db.NewPostgres(context.Background(), cfg.DBAddress, cfg.MaxRetries)
		if err != nil {
			return nil, fmt.Errorf("an error occured when opening DB: %w", err)
		}
	}

	return &Server{
		buildinfo: bi,
		router:    router,
		cfg:       cfg,
		db:        pg,
	}, nil
}

// ListenAndServe launches listening loop on the address provided in the config.
func (s *Server) ListenAndServe() {
	if err := http.ListenAndServe(s.cfg.Address, s.router); err != nil {
		log.Fatal(err)
	}
}

// Run executes server's loop.
func (s *Server) Run() error {
	var err error

	var stg service.MetricStorage
	if s.cfg.DBAddress != "" {
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
	if s.cfg.FilePath != "" {
		d, err = middleware.NewMetricDumper(service, s.cfg.StoreInterval, s.cfg.FilePath, s.cfg.RestoreMetrics)
		if err != nil {
			return fmt.Errorf("an error occurred when starting Server: %w", err)
		}
		defer d.Dump()
	}

	auditFileDst, err := audit.NewFileDst(s.cfg.Audit)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("an error occurred when starting Server: %w", err)
	}
	auditURLDist := audit.NewURLDst(s.cfg.Audit)
	auditLogger := audit.NewLogger(
		s.cfg.Audit,
		auditFileDst,
		auditURLDist,
	)
	defer auditLogger.CloseDestinations()

	h := handler.NewMetricsHandler(service, s.db, auditLogger)

	s.router.Get("/", h.ListMetrics())
	s.router.Post("/update/", d.Middleware(h.UpdateMetricJSON()))
	s.router.Post("/update/*", d.Middleware(h.UpdateMetric()))
	s.router.Post("/value/", h.GetMetricJSON())
	s.router.Get("/value/*", h.GetMetric())
	s.router.Get("/ping", h.PingDB())
	s.router.Post("/updates/", d.Middleware(h.UpdateMultipleJSON()))

	s.router.Handle("/debug/*", http.DefaultServeMux)

	fmt.Println(s.buildinfo.String())

	d.DumperLoop(service)
	go s.ListenAndServe()

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGTERM, syscall.SIGINT)
	<-term

	return nil
}
