package handler

import (
	"net/http"
	"yield-mvp/internal/service/exchangesvc"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/render"
)

func GenerateExchangeReport(
	wl wlog.Logger,
	exchangeService exchangesvc.SVC,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "exchange")
		data, err := exchangeService.GenereateReport(ctx, wl)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		// TODO: Add new render method to render a HTML table with the data
		render.JSON(ctx, wl, w, data, http.StatusOK)
	}
}
