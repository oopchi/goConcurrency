package app

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var tempPath string = filepath.Join(os.Getenv("TEMP"), "pipeline-temp")

func Start() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	dummyfilegenerator()
	durationWithoutConcurrency := renamefilegenerator()
	dummyfilegenerator()
	durationWithConcurrency := renamefilegeneratorConcurrent()

	log.Println("Duration without concurrency:", durationWithoutConcurrency, "seconds")
	log.Println("Duration with concurrency:", durationWithConcurrency, "seconds")
}
