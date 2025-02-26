package handler

import (
	"net/http"
	"yield-mvp/internal/service/reportsvc"
	"yield-mvp/internal/wlog"
)

func GenerateReport(
	wl wlog.Logger,
	reportService reportsvc.SVC,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
