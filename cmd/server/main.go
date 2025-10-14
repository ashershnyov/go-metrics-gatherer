package main

import (
	"github.com/ashershnyov/go-metrics-gatherer/internal/server"
)

func main() {
	service := server.New()
	err := service.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
