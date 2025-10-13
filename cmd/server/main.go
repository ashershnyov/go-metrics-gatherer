package main

import (
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/config"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server/service"
)

func main() {
	cfg := config.New()
	service := service.New(cfg)
	err := service.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
