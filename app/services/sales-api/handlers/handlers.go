package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/labstack/echo/v4"
	"github.com/shalimski/ultimateservice/app/services/sales-api/handlers/debug/checkgrp"
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

func APIMux() *echo.Echo {
	e := echo.New()

	return e
}
