package disablehdl

import (
	"net/http"

	"github.com/wingocard/braavos/internal/wlog"
	"github.com/wingocard/braavos/pkg/render"
)

const responseCode = 500

func DisabledEndpoint(wl wlog.Logger, api string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "disabled")

		wl.Debugf("disabled API requested: %s", api)

		render.ErrInternal()
	}
}
