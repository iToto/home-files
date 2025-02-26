// Package disablehdl provides a standardized way of disabling an API
package disablehdl

import (
	"errors"
	"net/http"

	"github.com/wingocard/braavos/internal/wlog"
	"github.com/wingocard/braavos/pkg/render"
)

const disabledAPIErrString = "api disabled"

// DisabledEndpoint will log the attempted request of a disabled endpoint
// and return an internal server error
// used to globally block clients from accessing an API
func DisabledEndpoint(wl wlog.Logger, api string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "disabled")

		wl.Debugf("disabled API requested: %s", api)

		render.InternalError(ctx, wl, w, errors.New(disabledAPIErrString))
	}
}
