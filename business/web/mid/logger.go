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
			v, err := web.GetValues(ctx)
			if err != nil {
				return err
			}

			log.Infow("request start", "traceID", v.TraceID, "method", r.Method, "path", r.URL.Path)

			err = handler(ctx, w, r)

			log.Infow("request completed", "traceID", v.TraceID, "method", r.Method, "path", r.URL.Path, "duration", time.Since(v.Now), "status", v.StatusCode)

			return err
		}

		return h
	}

	return m
}
