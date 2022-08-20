package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/shalimski/ultimateservice/app/web"
	"go.uber.org/zap"
)

// Logger middleware for request
func Logger(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			log.Infow("request start", "method", r.Method, "path", r.URL.Path)
			start := time.Now()

			err := handler(ctx, w, r)

			log.Infow("request completed", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))

			return err
		}

		return h
	}

	return m
}
