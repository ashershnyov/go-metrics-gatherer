package main

import (
	"flag"

	"github.com/ashershnyov/go-metrics-gatherer/internal/server"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/storage"
)

func main() {
	address := flag.String("a", "localhost:8080", "specifies the address for the server to start on")
	flag.Parse()
	storage := storage.NewMetricStorage()
	service := server.New(storage, server.SetAddress(address))
	err := service.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
