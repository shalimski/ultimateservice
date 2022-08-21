package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/shalimski/ultimateservice/app/web"
)

// Panics recovers from panics and converts the panic to an error so it is handled in Errors.
func Panics() web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {

					// Stack trace will be provided.
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

				}
			}()

			// Call the next handler and set its return value in the err variable.
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
