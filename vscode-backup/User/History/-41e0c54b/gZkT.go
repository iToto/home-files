package handler

import (
	"net/http"
	"yield-mvp/internal/service/reportsvc"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/render"
)

func GenerateReport(
	wl wlog.Logger,
	reportService reportsvc.SVC,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "report")
		err := reportService.GenerateTradeReport(ctx, wl)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}
		render.JSON(ctx, wl, w, nil, http.StatusOK)
	}
}
