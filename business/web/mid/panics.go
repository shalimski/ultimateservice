package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/shalimski/ultimateservice/app/web"
	"github.com/shalimski/ultimateservice/business/sys/metrics"
)

// Panics recovers from panics and converts the panic to an error so it is handled in Errors.
func Panics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {

					// Stack trace will be provided.
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

					metrics.AddPanic(ctx)

				}
			}()

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
