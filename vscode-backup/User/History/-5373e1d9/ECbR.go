package disablehdl

import (
	"errors"
	"net/http"

	"github.com/wingocard/braavos/internal/wlog"
	"github.com/wingocard/braavos/pkg/render"
)

const responseCode = 500
const disabledAPIErrString = "api disabled"

func DisabledEndpoint(wl wlog.Logger, api string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "disabled")

		wl.Debugf("disabled API requested: %s", api)

		render.InternalError(ctx, wl, w, errors.New(disabledAPIErrString))
	}
}
