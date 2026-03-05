package main

import (
	"log"

	"github.com/ashershnyov/go-metrics-gatherer/internal/buildinfo"
	"github.com/ashershnyov/go-metrics-gatherer/internal/server"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	bi := buildinfo.New(buildVersion, buildDate, buildCommit)

	srv, err := server.New(bi)
	if err != nil {
		log.Fatal(err)
	}

	if err = srv.Run(); err != nil {
		log.Fatal(err)
	}
}
