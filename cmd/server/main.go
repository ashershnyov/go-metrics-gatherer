package main

import (
	"flag"
	"strconv"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
)

const (
	defaultFilename      string        = "metrics.json"
	defaultStoreInterval time.Duration = 300 * time.Second
)

type envs struct {
	Address       string `env:"ADDRESS" envDefault:""`
	FilePath      string `env:"FILE_STORAGE_PATH" envDefault:""`
	StoreInterval int    `env:"STORE_INTERVAL" envDefault:"-1"`
	Restore       string `env:"RESTORE" envDefault:""`
}

func main() {
	var err error

	address := flag.String("a", "localhost:8080", "specifies the address for the server to start on")
	storeInterval := flag.Int("i", int(defaultStoreInterval), "specifies the interval between writes to the specified file")
	filePath := flag.String("f", defaultFilename, "specifies the filepath to store metric values in")
	restore := flag.Bool("r", false, "indicates whether the stored metrics should be loaded from the specified file on server startup")
	flag.Parse()

	var envs envs
	err = env.Parse(&envs)
	if err != nil {
		panic(err)
	}

	if envs.Address != "" {
		address = &envs.Address
	}

	if envs.StoreInterval >= 0 {
		storeInterval = &envs.StoreInterval
	}

	if envs.FilePath != "" {
		filePath = &envs.FilePath
	}

	if envs.Restore != "" {
		v, err := strconv.ParseBool(envs.Restore)
		if err != nil {
			panic(err)
		}
		restore = &v
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugarLogger := logger.Sugar()

	srv, err := server.New(
		sugarLogger,
		server.SetAddress(address),
		server.SetFilePath(filePath),
		server.SetStoreInterval(storeInterval),
		server.SetRestoreMetrics(restore),
	)

	if err != nil {
		panic(err)
	}

	if err = srv.Run(); err != nil {
		panic(err)
	}
}
