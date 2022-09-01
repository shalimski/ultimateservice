package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/shalimski/ultimateservice/app/services/sales-api/handlers"
	"github.com/shalimski/ultimateservice/business/sys/auth"
	"github.com/shalimski/ultimateservice/business/sys/database"
	"github.com/shalimski/ultimateservice/internal/config"
	"github.com/shalimski/ultimateservice/pkg/logger"
	"github.com/shalimski/ultimateservice/pkg/logger/keystore"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

var build = "develop"

func main() {
	log, err := logger.New("SALES-API")
	if err != nil {
		panic(err)
	}

	defer log.Sync() //nolint:errcheck

	cfg := config.New()
	log.Infof("configuration %+v", cfg)

	// auth support
	ks, err := keystore.NewFS(os.DirFS(cfg.AuthKeysFolder))
	if err != nil {
		log.Errorw("reading keys", err)
		return
	}

	auth, err := auth.NewAuth(cfg.AuthActiveKID, ks)
	if err != nil {
		log.Errorf("init auth: %w", err)
		return
	}

	// Database Support

	log.Infow("startup", "status", "initializing database support", "host", cfg.DB.Host)

	db, err := database.Open(database.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		log.Errorf("connecting to db: %w", err)
		return
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DB.Host)
		db.Close()
	}()

	// debug handlers
	debugMux := handlers.DebugMux(build, log, db)

	go func() {
		if err := http.ListenAndServe(cfg.DebugURI, debugMux); err != nil {
			log.Errorw("debug router is closer", err)
		}
	}()

	log.Infof("starting sales-api build[%s] CPU[%d]\n", build, runtime.NumCPU())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown: shutdown,
		Log:      log,
		Auth:     auth,
		DB:       db,
	})

	api := http.Server{
		Addr:         cfg.Host,
		Handler:      apiMux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "api handler started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		log.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			log.Errorw("could not stop server gracefully", "error", err)
		}
	}
}
