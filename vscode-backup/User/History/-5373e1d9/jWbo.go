package disablehdl

import (
	"net/http"

	"github.com/wingocard/braavos/internal/wlog"
)

const responseCode = 500

func DisabledEndpoint(wl wlog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriteer, r *http.Request) {

	}
}
