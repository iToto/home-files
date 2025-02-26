package handler

import (
	"net/http"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
	"yield-mvp/pkg/render"
)

func GetCoinRoutesExchangeAccounts(wl wlog.Logger, cc *coinroutesapi.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "coinroutes")

		wl.Debug("requesting coin routes exchange accounts")

		ea, err := cc.GetExchangeAccounts(ctx)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
		}

		wl.Debugf("response: %+v", ea)
		render.JSON(ctx, wl, w, ea, http.StatusOK)
	}
}
