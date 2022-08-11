package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	_ "go.uber.org/automaxprocs"
)

var build = "develop"

func main() {
	log.Printf("starting sales-api build[%s] CPU[%d]\n", build, runtime.NumCPU())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("stopping sales-api")
}
