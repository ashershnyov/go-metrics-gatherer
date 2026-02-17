package main

import (
	"github.com/ashershnyov/go-metrics-gatherer/internal/staticlint"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(staticlint.All()...)
}
