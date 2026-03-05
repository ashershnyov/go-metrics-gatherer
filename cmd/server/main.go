package main

import (
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/ashershnyov/go-metrics-gatherer/internal/buildinfo"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/config"
	"github.com/caarlos0/env"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

type envVars struct {
	Address       string `env:"ADDRESS" envDefault:""`
	FilePath      string `env:"FILE_STORAGE_PATH" envDefault:""`
	Restore       string `env:"RESTORE" envDefault:""`
	DBAddress     string `env:"DATABASE_DSN" envDefault:""`
	Key           string `env:"KEY" envDefault:""`
	AuditURL      string `env:"AUDIT_URL" envDefault:""`
	AuditFile     string `env:"AUDIT_FILE" envDefault:""`
	StoreInterval int    `env:"STORE_INTERVAL" envDefault:"-1"`
	CryptoKey     string `env:"CRYPTO_KEY" envDefault:""`
}

func main() {
	var err error

	address := flag.String("a", "localhost:8080", "specifies the address for the server to start on")
	storeInterval := flag.Int("i", int(300*time.Second), "specifies the interval between writes to the specified file")
	filePath := flag.String("f", "metrics.json", "specifies the filepath to store metric values in")
	restore := flag.Bool("r", true, "indicates whether the stored metrics should be loaded from the specified file on server startup")
	dbAddress := flag.String("d", "", "specifies the address of the DB")
	key := flag.String("k", "", "specifies the key to use to hash the response body")
	auditURL := flag.String("audit-url", "", "specifies the URL to send audit logs to")
	auditFile := flag.String("audit-file", "", "specifies the filepath to write audit logs to")
	cryptoKey := flag.String("crypto-key", "", "specifies the filepath to server's private key")
	flag.Parse()

	var envs envVars
	err = env.Parse(&envs)
	if err != nil {
		log.Fatal(err)
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

	if envs.DBAddress != "" {
		dbAddress = &envs.DBAddress
	}

	if envs.Key != "" {
		key = &envs.Key
	}

	if envs.AuditFile != "" {
		auditFile = &envs.AuditFile
	}

	if envs.AuditURL != "" {
		auditURL = &envs.AuditURL
	}

	if envs.CryptoKey != "" {
		cryptoKey = &envs.CryptoKey
	}

	if envs.Restore != "" {
		v, err := strconv.ParseBool(envs.Restore)
		if err != nil {
			log.Fatal(err)
		}
		restore = &v
	}

	bi := buildinfo.New(buildVersion, buildDate, buildCommit)

	srv, err := server.New(
		bi,
		config.SetAddress(address),
		config.SetFilePath(filePath),
		config.SetStoreInterval(storeInterval),
		config.SetRestoreMetrics(restore),
		config.SetDBAddress(dbAddress),
		config.SetKey(key),
		config.SetAuditFilePath(auditFile),
		config.SetAuditURL(auditURL),
		config.SetCryptoKeyPath(cryptoKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = srv.Run(); err != nil {
		log.Fatal(err)
	}
}
