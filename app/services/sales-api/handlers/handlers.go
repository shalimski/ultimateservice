package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/shalimski/ultimateservice/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/shalimski/ultimateservice/app/web"
	"github.com/shalimski/ultimateservice/business/sys/auth"
	"github.com/shalimski/ultimateservice/business/web/mid"
	"go.uber.org/zap"
)

// DebugMux pprof handlers for profiling
func DebugMux(build string, log *zap.SugaredLogger) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	chg := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}

	mux.HandleFunc("/debug/liveness", chg.Liveness)
	mux.HandleFunc("/debug/readyness", chg.Readiness)

	return mux
}

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
}

func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)
	v1(app, cfg)

	return app
}

// v1 binds all the version 1 routes
func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"

	// create handlers

	app.Handle(http.MethodGet, version, "/test", nil)
}
