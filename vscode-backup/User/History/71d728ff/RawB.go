package handler

import (
	"encoding/json"
	"net/http"
	"on-air/internal/entities"
	"on-air/internal/service/onair"
	"on-air/internal/wlog"
	"on-air/pkg/render"
)

type onAirStatusBody struct {
	IsOnAir bool `json:"is_on_air,omitempty"`
}

func GetOnAirStatus(wl wlog.Logger, onAirService onair.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		onAirStatus, err := onAirService.GetOnAirStatus(ctx, wl)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
		}

		render.JSON(ctx, wl, w, onAirStatus, http.StatusOK)
	}
}

func ToggleOnAirStatus(wl wlog.Logger, onAirService onair.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		onAir, err := onAirService.ToggleOnAirStatus(ctx, wl)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
		}

		render.JSON(ctx, wl, w, onAir, http.StatusOK)
	}
}

func SetOnAirStatus(wl wlog.Logger, onAirService onair.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// get the onAirStatus from the request body
		var onAirReq onAirStatusBody
		if err := json.NewDecoder(r.Body).Decode(&onAirReq); err != nil {
			render.BadRequest(ctx, wl, w, err)
			return
		}

		onAir := entities.OnAirStatus{IsOnAir: onAirReq.IsOnAir}

		err := onAirService.SetOnAirStatus(ctx, wl, onAir)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
		}

		render.JSON(ctx, wl, w, onAir, http.StatusOK)
	}
}
