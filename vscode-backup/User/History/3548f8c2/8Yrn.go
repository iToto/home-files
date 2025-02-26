// Package signalhdl is the handler that handles all signal HTTP requests
package handler

import (
	"fmt"
	"net/http"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/service/signalsvc"
	"yield-mvp/internal/wlog"
	"yield-mvp/internal/ycontext"
	"yield-mvp/pkg/render"

	"github.com/gorilla/mux"
)

func ProcessSignalsForUser(
	wl wlog.Logger,
	signalService signalsvc.SVC,
	userID string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// FIXME: This should be added by a middleware
		r = r.WithContext(ycontext.WithUserID(r.Context(), userID))
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "user-signal")
		err := signalService.ProcessSignalsAndStratgiesForUser(ctx, wl, userID)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		render.JSON(ctx, wl, w, nil, http.StatusOK)
	}
}

func GetBTCSignal(
	wl wlog.Logger,
	signalService signalsvc.SVC,
	signalStrategies []entities.SignalStrategies,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// FIXME: This should be added by a middleware
		chain := entities.BTC
		r = r.WithContext(ycontext.WithChain(r.Context(), string(chain)))
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "signal")

		for _, signalStrategy := range signalStrategies {
			err := signalService.GetAndProcessSignal(ctx, wl, &signalStrategy)
			if err != nil {
				render.InternalError(ctx, wl, w, err)
				return
			}

		}

		render.JSON(ctx, wl, w, nil, http.StatusOK)
	}
}

func GetETHSignal(
	wl wlog.Logger,
	signalService signalsvc.SVC,
	signalStrategies []entities.SignalStrategies,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// FIXME: This should be added by a middleware
		chain := entities.ETH
		r = r.WithContext(ycontext.WithChain(r.Context(), string(chain)))
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "signal")

		for _, signalStrategy := range signalStrategies {
			err := signalService.GetAndProcessSignal(ctx, wl, &signalStrategy)
			if err != nil {
				render.InternalError(ctx, wl, w, err)
				return
			}
		}

		render.JSON(ctx, wl, w, nil, http.StatusOK)
	}
}

func DisableStrategy(
	wl wlog.Logger,
	signalService signalsvc.SVC,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		strategyName, ok := vars["name"]
		if !ok {
			render.BadRequest(ctx, wl, w, fmt.Errorf("missing required name in URL"))
			return
		}

		wl := wlog.WithServiceRequest(ctx, wl, "signal")
		wl = wlog.WithStrategy(wl, strategyName)

		strategy, err := signalService.DisableStrategyAndGoNeutral(ctx, wl, strategyName)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		render.JSON(ctx, wl, w, strategy, http.StatusOK)
	}
}

func EnableStrategy(
	wl wlog.Logger,
	signalService signalsvc.SVC,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		strategyName, ok := vars["name"]
		if !ok {
			render.BadRequest(ctx, wl, w, fmt.Errorf("missing required name in URL"))
			return
		}

		wl := wlog.WithServiceRequest(ctx, wl, "signal")
		wl = wlog.WithStrategy(wl, strategyName)

		strategy, err := signalService.EnableStrategy(ctx, wl, strategyName)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		render.JSON(ctx, wl, w, strategy, http.StatusOK)
	}
}

func CreateStrategy(
	wl wlog.Logger,
	signalService signalsvc.SVC,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		strategyName, ok := vars["name"]
		if !ok {
			render.BadRequest(ctx, wl, w, fmt.Errorf("missing required name in URL"))
			return
		}

		wl := wlog.WithServiceRequest(ctx, wl, "signal")
		wl = wlog.WithStrategy(wl, strategyName)

		strategy, err := signalService.EnableStrategy(ctx, wl, strategyName)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		render.JSON(ctx, wl, w, strategy, http.StatusOK)
	}
}
