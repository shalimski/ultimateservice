package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/shalimski/ultimateservice/pkg/logger"
	_ "go.uber.org/automaxprocs"
)

var build = "develop"

func main() {
	log, err := logger.New("SALES-API")
	if err != nil {
		panic(err)
	}

	defer log.Sync()

	log.Infof("starting sales-api build[%s] CPU[%d]\n", build, runtime.NumCPU())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Info("stopping sales-api")
}
