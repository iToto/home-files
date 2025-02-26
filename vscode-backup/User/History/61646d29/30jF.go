// Package signalhdl is the handler that handles all signal HTTP requests
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
	AccountLeverage  int        `json:"account_leverage"`
	FixedTradeAmount null.Float `json:"fixed_trade_amount"`
	TradeStrategy    string     `json:"trade_strategy"`
	CurrencyPair     string     `json:"currency_pair"`
}

func (c createStrategyForm) Validate() error {
	return validation.ValidateStruct(&c,
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

type updateStrategyForm struct {
	ID               string     `json:"id"`
	SignalSourceID   string     `json:"signal_source_id"`
	Type             string     `json:"type"`
	Name             string     `json:"name"`
	Exchange         string     `json:"exchange"`
	Margin           string     `json:"margin"`
	Leverage         int        `json:"leverage"`
	AccountLeverage  int        `json:"account_leverage"`
	FixedTradeAmount null.Float `json:"fixed_trade_amount"`
	TradeStrategy    string     `json:"trade_strategy"`
	CurrencyPair     string     `json:"currency_pair"`
}

type createSignalForm struct {
	Enabled       bool   `json:"enabled"`
	Type          string `json:"type"`
	IP            string `json:"ip"`
	SignalVersion int    `json:"signal_version"`
}

func (c createSignalForm) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.IP,
			validation.Required.Error("missing required property: ip"),
		),
		validation.Field(
			&c.Type,
			validation.Required.Error("missing required property: type"),
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

		wl := wlog.WithServiceRequest(ctx, wl, "disable-strategy")
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

		wl := wlog.WithServiceRequest(ctx, wl, "enable-strategy")
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

		var al null.Int

		if req.AccountLeverage == 0 {
			al = null.IntFrom(0)
		} else {
			al = null.IntFrom(int64(req.AccountLeverage))
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
			Leverage:         req.Leverage,
			AccountLeverage:  al,
			FixedTradeAmount: req.FixedTradeAmount,
			TradeStrategy:    entities.TradeStrategyType(req.TradeStrategy),
			CurrencyPair:     entities.CurrencyPairType(req.CurrencyPair),
		}

		wl := wlog.WithServiceRequest(ctx, wl, "create-strategy")
		wl = wlog.WithStrategy(wl, strategy.Name)

		cs, err := signalService.CreateStrategy(ctx, wl, &strategy)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		render.JSON(ctx, wl, w, cs, http.StatusOK)
	}
}

func UpdateStrategy(
	wl wlog.Logger,
	signalService signalsvc.SVC,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := &updateStrategyForm{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.BadRequest(ctx, wl, w, render.ErrJSONDecode)
			return
		}

		// build entity from req
		strategy := entities.Strategy{
			ID:               req.ID,
			SignalSourceID:   req.SignalSourceID,
			Type:             entities.StrategyType(req.Type),
			Name:             req.Name,
			Exchange:         entities.ExchangeType(req.Exchange),
			Margin:           entities.MarginType(req.Margin),
			Leverage:         req.Leverage,
			FixedTradeAmount: req.FixedTradeAmount,
			TradeStrategy:    entities.TradeStrategyType(req.TradeStrategy),
			CurrencyPair:     entities.CurrencyPairType(req.CurrencyPair),
		}

		wl := wlog.WithServiceRequest(ctx, wl, "update-strategy")
		wl = wlog.WithStrategy(wl, strategy.Name)

		cs, err := signalService.UpdateStrategy(ctx, wl, &strategy)
		if err != nil {

			render.InternalError(ctx, wl, w, err)
			return
		}

		render.JSON(ctx, wl, w, cs, http.StatusOK)
	}
}

func GetStrategies(
	wl wlog.Logger,
	signalService signalsvc.SVC,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		wl := wlog.WithServiceRequest(ctx, wl, "get-strategies")

		// check for filters
		// enabled filter
		vars := r.URL.Query()
		enabledFilterPresent := false
		enabled := false

		if ev, err := strconv.ParseBool(vars.Get("enabled")); err == nil {
			enabledFilterPresent = true
			enabled = ev
		}

		wl.Debugf("enabled filter: %v", enabled)

		strats, err := signalService.GetStrategies(ctx, wl, enabledFilterPresent, enabled)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		render.JSON(ctx, wl, w, strats, http.StatusOK)
	}
}

func CreateSignal(
	wl wlog.Logger,
	signalService signalsvc.SVC,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := &createSignalForm{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.BadRequest(ctx, wl, w, render.ErrJSONDecode)
			return
		}
		if err := req.Validate(); err != nil {
			render.BadRequest(ctx, wl, w, err)
			return
		}

		// build entity from req
		signal := entities.SignalSource{
			Enabled:       req.Enabled,
			Type:          entities.ChainType(req.Type),
			IP:            req.IP,
			SignalVersion: int64(req.SignalVersion),
		}

		wl := wlog.WithServiceRequest(ctx, wl, "create-signal")
		wl = wlog.WithSignalSource(wl, signal.IP)

		cs, err := signalService.CreateSignal(ctx, wl, &signal)
		if err != nil {
			if errors.Is(err, signalsvc.ErrClientBadSignalType) {
				render.BadRequest(ctx, wl, w, err)
				return
			}

			if errors.Is(err, signalsvc.ErrClientBadVersion) {
				render.BadRequest(ctx, wl, w, err)
				return
			}

			render.InternalError(ctx, wl, w, err)
			return
		}

		render.JSON(ctx, wl, w, cs, http.StatusOK)
	}
}

func GetSignals(
	wl wlog.Logger,
	signalService signalsvc.SVC,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		wl := wlog.WithServiceRequest(ctx, wl, "get-signals")

		signals, err := signalService.GetSignals(ctx, wl)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		render.JSON(ctx, wl, w, signals, http.StatusOK)
	}
}
