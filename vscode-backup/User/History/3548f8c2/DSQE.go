// Package signalhdl is the handler that handles all signal HTTP requests
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/service/signalsvc"
	"yield-mvp/internal/wlog"
	"yield-mvp/internal/ycontext"
	"yield-mvp/pkg/render"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"github.com/guregu/null"
)

type createStrategyForm struct {
	Enabled          bool       `json:"enabled"`
	UserID           string     `json:"user_id"`
	SignalSourceID   string     `json:"signal_source_id"`
	Type             string     `json:"type"`
	Name             string     `json:"name"`
	Exchange         string     `json:"exchange"`
	Margin           string     `json:"margin"`
	Leverage         int        `json:"leverage"`
	FixedTradeAmount null.Float `json:"fixed_trade_amount"`
	TradeStrategy    string     `json:"trade_strategy"`
	CurrencyPair     string     `json:"currency_pair"`
}

func (c createStrategyForm) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.Enabled,
			validation.Required.Error("missing required property: enabled"),
		),
		validation.Field(
			&c.UserID,
			validation.Required.Error("missing required property: user_id"),
		),
		validation.Field(
			&c.SignalSourceID,
			validation.Required.Error("missing required property: signal_source_id"),
		),
		validation.Field(
			&c.Type,
			validation.Required.Error("missing required property: type"),
		),
		validation.Field(
			&c.Name,
			validation.Required.Error("missing required property: name"),
		),
		validation.Field(
			&c.Exchange,
			validation.Required.Error("missing required property: exchange"),
		),
		validation.Field(
			&c.Margin,
			validation.Required.Error("missing required property: margin"),
		),
		validation.Field(
			&c.Leverage,
			validation.Required.Error("missing required property: leverage"),
		),
		validation.Field(
			&c.TradeStrategy,
			validation.Required.Error("missing required property: trade_strategy"),
		),
		validation.Field(
			&c.CurrencyPair,
			validation.Required.Error("missing required property: currency_pair"),
		),
	)
}

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

		req := &createStrategyForm{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.BadRequest(ctx, wl, w, render.ErrJSONDecode)
			return
		}
		if err := req.Validate(); err != nil {
			render.BadRequest(ctx, wl, w, err)
			return
		}

		// build entity from req
		strategy := entities.Strategy{
			Enabled:          req.Enabled,
			UserID:           req.UserID,
			SignalSourceID:   req.SignalSourceID,
			Type:             entities.StrategyType(req.Type),
			Name:             req.Name,
			Exchange:         entities.ExchangeType(req.Exchange),
			Margin:           entities.MarginType(req.Margin),
			Leverage:         entities.LeverageAmount(req.Leverage),
			FixedTradeAmount: req.FixedTradeAmount,
			TradeStrategy:    entities.TradeStrategyType(req.TradeStrategy),
			CurrencyPair:     entities.CurrencyPairType(req.CurrencyPair),
		}

		wl := wlog.WithServiceRequest(ctx, wl, "signal")
		wl = wlog.WithStrategy(wl)

		strategy, err := signalService.EnableStrategy(ctx, wl, strategyName)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		render.JSON(ctx, wl, w, strategy, http.StatusOK)
	}
}
