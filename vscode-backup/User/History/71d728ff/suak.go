package handler

import (
	"net/http"
	"on-air/internal/service/onair"
	"on-air/internal/wlog"
	"on-air/pkg/render"
)

type oares struct {
	IsOnAir     bool   `json:"is_on_air,omitempty"`
	LastUpdated string `json:"last_updated,omitempty"`
	LastOnAir   string `json:"last_on_air,omitempty"`
}

func GetOnAirStatus(wl wlog.Logger, onAirService onair.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		onAirStatus, err := onAirService.GetOnAirStatus(ctx, wl)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
		}

		resp := oares{
			IsOnAir:     onAirStatus.IsOnAir,
			LastUpdated: onAirStatus.LastUpdated.Time.String(),
			LastOnAir:   onAirStatus.LastOnAir.Time.String(),
		}

		render.JSON(ctx, wl, w, resp, http.StatusOK)
	}
}
