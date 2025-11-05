package main

import (
	"flag"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/storage"
	"github.com/caarlos0/env"
)

type envs struct {
	Address string `env:"ADDRESS" envDefault:""`
}

func main() {
	var err error

	address := flag.String("a", "localhost:8080", "specifies the address for the server to start on")
	flag.Parse()

	var envs envs
	err = env.Parse(&envs)
	if err != nil {
		panic(err)
	}

	if envs.Address != "" {
		address = &envs.Address
	}

	storage := storage.NewMetricStorage()
	service := server.New(storage, server.SetAddress(address))
	err = service.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
