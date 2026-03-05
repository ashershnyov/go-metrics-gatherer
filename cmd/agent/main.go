package main

import (
	"log"

	"github.com/ashershnyov/go-metrics-gatherer/internal/agent"
	"github.com/ashershnyov/go-metrics-gatherer/internal/buildinfo"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	var err error

	bi := buildinfo.New(buildVersion, buildDate, buildCommit)

	agent, err := agent.New(bi)
	if err != nil {
		log.Fatal(err)
	}

	agent.Run()
}
