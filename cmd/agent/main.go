package main

import (
	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/config"
	"github.com/ashershnyov/go-metrics-gatherer/internal/agent/service"
)

func main() {
	cfg := config.New()
	agent := service.New(cfg)
	agent.Run()
}
