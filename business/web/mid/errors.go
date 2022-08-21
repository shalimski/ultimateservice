package mid

import (
	"context"
	"net/http"

	"github.com/shalimski/ultimateservice/app/web"
	"github.com/shalimski/ultimateservice/business/sys/validate"
	"go.uber.org/zap"
)

func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("valuse missing from contex")
			}

			if err := handler(ctx, w, r); err != nil {
				log.Errorw("ERROR", "traceID", v.TraceID, "ERROR", err)

				var er validate.ErrorResponse
				var status int
				switch {
				case validate.IsFieldErrors(err):
					fieldErrors := validate.GetFieldErrors(err)
					er = validate.ErrorResponse{
						Error:  "data validation error",
						Fields: fieldErrors.Error(),
					}
					status = http.StatusBadRequest
				case validate.IsRequestError(err):

					reqErr := validate.GetRequestError(err)

					er = validate.ErrorResponse{
						Error: er.Error,
					}

					status = reqErr.Status
				default:
					er = validate.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}

					status = http.StatusInternalServerError
				}
				// Respond with the error back to the client.
				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}

			}
			return nil
		}

		return h
	}

	return m
}
