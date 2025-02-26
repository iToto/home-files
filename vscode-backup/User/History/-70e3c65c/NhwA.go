package handler

import "yield-mvp/internal/wlog"

func GenerateExchangeReport(
	wl wlog.Logger,
	exchangeService exchangesvc.SVC,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "exchange")
		err := exchangeService.GenereateReport(ctx, wl)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}
		render.JSON(ctx, wl, w, nil, http.StatusOK)
	}
