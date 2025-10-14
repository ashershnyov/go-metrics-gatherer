package main

import "github.com/ashershnyov/go-metrics-gatherer/internal/agent"

func main() {
	agent := agent.New()
	agent.Run()
}
