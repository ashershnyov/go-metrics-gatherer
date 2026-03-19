package server

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/buildinfo"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/audit"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/config"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/db"
	grpchandler "github.com/ashershnyov/go-metrics-gatherer/internal/server/handler/grpc"
	httphandler "github.com/ashershnyov/go-metrics-gatherer/internal/server/handler/http"
	grpcmiddleware "github.com/ashershnyov/go-metrics-gatherer/internal/server/middleware/grpc"
	httpmiddleware "github.com/ashershnyov/go-metrics-gatherer/internal/server/middleware/http"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/service"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/storage"
	"github.com/ashershnyov/go-metrics-gatherer/pkg/api/proto"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Server defines a server.
type Server struct {
	http      http.Server
	grpc      *grpc.Server
	buildinfo buildinfo.BuildInfo
	cfg       *config.Config
	db        *db.Postgres
}

// New creates a server using a provided cfg.
func New(bi buildinfo.BuildInfo) (*Server, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating server: %w", err)
	}

	var pg *db.Postgres

	if cfg.DBAddress != "" {
		pg, err = db.NewPostgres(context.Background(), cfg.DBAddress, cfg.MaxRetries)
		if err != nil {
			return nil, fmt.Errorf("an error occured when opening DB: %w", err)
		}
	}

	return &Server{
		http: http.Server{
			Addr: cfg.Address,
		},
		cfg: cfg,
		db:  pg,
	}, nil
}

// ListenAndServeHTTP launches listening loop on the address provided in the config.
func (s *Server) ListenAndServeHTTP() {
	if err := s.http.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// ListenAndServeGrpc launches listening loop on the GRPC address provided in the config.
func (s *Server) ListenAndServeGrpc(listener net.Listener) {
	if err := s.grpc.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) loadCryptoKey() (*rsa.PrivateKey, error) {
	bytes, err := os.ReadFile(s.cfg.CryptoKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading reading crypto key file: %w", err)
	}

	pemBlock, _ := pem.Decode(bytes)
	if pemBlock == nil {
		return nil, fmt.Errorf("error parsing crypto key: no PEM block found")
	}

	key, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing crypto key: %w", err)
	}

	return key, nil
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

	var d *httpmiddleware.MetricDumper
	if s.cfg.FilePath != "" {
		d, err = httpmiddleware.NewMetricDumper(
			service, s.cfg.StoreInterval,
			s.cfg.FilePath, s.cfg.RestoreMetrics,
		)
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

	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("an error occurred when starting Server: %w", err)
	}
	defer logger.Sync()
	sugarLogger := logger.Sugar()

	var key *rsa.PrivateKey
	if s.cfg.CryptoKeyPath != "" {
		key, err = s.loadCryptoKey()
		if err != nil {
			return fmt.Errorf("an error occurred when starting Server: %w", err)
		}
	}

	router := chi.NewRouter()
	router.Use(
		httpmiddleware.Logging(sugarLogger),
		httpmiddleware.Gzip(),
		httpmiddleware.Hashing(s.cfg.Key),
		httpmiddleware.Decrypt(key),
		httpmiddleware.CheckIP(s.cfg.TrustedSubnet.IPNet),
	)

	h := httphandler.NewMetricsHandler(service, s.db, auditLogger)

	router.Get("/", h.ListMetrics())
	router.Post("/update/", d.Middleware(h.UpdateMetricJSON()))
	router.Post("/update/*", d.Middleware(h.UpdateMetric()))
	router.Post("/value/", h.GetMetricJSON())
	router.Get("/value/*", h.GetMetric())
	router.Get("/ping", h.PingDB())
	router.Post("/updates/", d.Middleware(h.UpdateMultipleJSON()))
	router.Handle("/debug/*", http.DefaultServeMux)

	s.http.Handler = router

	log.Println(s.buildinfo.String())

	d.DumperLoop(service)
	go s.ListenAndServeHTTP()

	if s.cfg.GrpcAddress != "" {
		srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
			grpcmiddleware.CheckIP(s.cfg.TrustedSubnet.IPNet)))

		handler := grpchandler.NewMetricsHandler(service, s.db)
		proto.RegisterMetricsServer(srv, handler)

		s.grpc = srv

		listen, err := net.Listen("tcp", s.cfg.GrpcAddress)
		if err != nil {
			return fmt.Errorf("could not setup gRPC listener: %w", err)
		}

		go s.ListenAndServeGrpc(listen)
	}

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-term

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down server: %w", err)
	}

	return nil
}
