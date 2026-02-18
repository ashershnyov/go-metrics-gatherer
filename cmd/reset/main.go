package main

import (
	"path/filepath"
	"runtime"

	"github.com/ashershnyov/go-metrics-gatherer/internal/resetgen"
)

var ignoreSuffixes = []string{
	"_test.go",
	"_mock.go",
	".gen.go",
}

func main() {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get current file path")
	}
	root := filepath.Dir(filepath.Dir(currentFile))
	err := resetgen.Do(root, ignoreSuffixes, nil)
	if err != nil {
		panic(err)
	}
}
