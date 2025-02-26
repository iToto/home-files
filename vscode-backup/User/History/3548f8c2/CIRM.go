// Package signalhdl is the handler that handles all signal HTTP requests
package handler

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/signallogger"
	"yield-mvp/internal/tradelogger"
	"yield-mvp/internal/wlog"
	"yield-mvp/internal/ycontext"
	"yield-mvp/pkg/coinroutesapi"
	"yield-mvp/pkg/render"
	"yield-mvp/pkg/signalapi"

	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid"
)

const (
	twentyMinutes     int     = 1200
	fiveMinutes       int     = 300
	fiveSeconds       int     = 5
	tradeAmountBuffer float64 = 0.98
)

func GetBTCSignal(
	wl wlog.Logger,
	sc *signalapi.Client,
	cc *coinroutesapi.Client,
	strats []string,
	db *sqlx.DB,
	dl *tradelogger.DataLogger,
	sl *signallogger.DataLogger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(ycontext.WithChain(r.Context(), chain))
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "signal")

		render.JSON(ctx, wl, w, nil, http.StatusOK)
	}
}

func GetETHSignal(
	wl wlog.Logger,
	sc *signalapi.Client,
	cc *coinroutesapi.Client,
	strats []string,
	db *sqlx.DB,
	dl *tradelogger.DataLogger,
	sl *signallogger.DataLogger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// FIXME: This should be added by a middleware
		chain := "eth"
		r = r.WithContext(ycontext.WithChain(r.Context(), chain))
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "signal")

		signalResp, err := sc.GetETHSignal(ctx)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		if signalResp == nil {
			err = errors.New("received empty signal from api")
			render.InternalError(ctx, wl, w, err)
			return
		}

		wl.Infof("signal received: %+v", signalResp)

		// 	parse last signal to see if we have an operation to make
		signal, err := checkForSignal(ctx, wl, signalResp, db, sl)
		if err != nil {
			render.InternalError(
				ctx,
				wl,
				w,
				fmt.Errorf("could not find past signal for chain: %s with error: %w", chain, err),
			)
			return
		}

		// execute signal on all strategies
		for _, v := range strats {
			strategy := coinroutesapi.StrategyType(v)
			if err := strategy.Validate(); err != nil {
				wl.Infof("skipping invalid strategy: %s", v)
				continue
			}

			wl.Infof("handling strategy: %s", strategy)
			var resp *coinroutesapi.ClientOrderCreateResponse
			wl = wlog.WithStrategy(wl, v)

			// parse signal from signal API and create proper payload for CoinRoutes
			payload, err := parseSignalAndCreateTradePayload(
				ctx,
				wl,
				signal,
				coinroutesapi.ETH,
				strategy,
			)
			if err != nil {
				render.InternalError(ctx, wl, w, err)
			}

			// if we have a payload, then signal told us to make a trade
			if payload == nil {
				wl.Info("no payload, skipping")
				continue
			}

			// set amount
			switch payload.Side {
			case coinroutesapi.Buy:
				// drop contracts (positions)
				position, err := getCoinRoutesPositionForStrategy(ctx, wl, strategy, cc, payload.CurrencyPair)
				if err != nil {
					wl.Error(err)
					continue
				}
				payload.Quantity = strconv.FormatFloat(position, 'f', 10, 32)

			case coinroutesapi.Sell:
				// short available balance (create contracts)
				balance, err := getCoinRoutesBalanceForStrategy(
					ctx,
					wl,
					strategy,
					cc,
					coinroutesapi.ETH)
				if err != nil {
					wl.Error(err)
					continue
				}

				// add multiple for long strategy
				if strategy == coinroutesapi.ETHLongShort {
					wl.Debug("multiplying amount for LongShort")
					balance = 2 * balance // 2x leverage for long-short
				}

				payload.Quantity = strconv.FormatFloat(balance, 'f', 10, 32)
			}

			// check for zero amount
			orderFloat := 0.0
			if payload.Quantity != "" {
				orderFloat, err = strconv.ParseFloat(payload.Quantity, 32)
				if err != nil {
					render.InternalError(ctx, wl, w, err)
					return
				}
			}

			if orderFloat == 0 {
				wl.Infof("no-op: quantity is zero value: %+v", payload.Quantity)
				render.JSON(ctx, wl, w, signal, http.StatusOK)
				return
			}

			// apply buffer to trade amount and overwrite payload
			wl.Debugf(
				"applying buffer of %f. Original: %f Actual: %f",
				tradeAmountBuffer,
				orderFloat,
				orderFloat*tradeAmountBuffer,
			)
			orderFloat *= tradeAmountBuffer
			payload.Quantity = strconv.FormatFloat(orderFloat, 'f', 10, 32)

			// Make Trade
			wl.Debugf("about to create order with payload: %+v", payload)
			resp, err = cc.CreateClientOrders(ctx, payload)
			if err != nil {
				wl.Error(fmt.Errorf("unable to place order with error: %w", err))
				continue
			}

			// record trade in log
			err = insertLatestSignalTradedForChain(ctx, wl, chain, signal, db)
			if err != nil {
				render.InternalError(ctx, wl, w, err)
				return
			}

			if resp.ClientOrderId == "" {
				wl.Error(fmt.Errorf("no client_order_id found in response: %+v", resp))
				continue
			}

			wl.Infof("order placed with coinroutes: %+v", resp)

			// record trade in order table
			err = insertNewOrder(ctx, wl, resp, db)
			if err != nil {
				wl.Error(fmt.Errorf("unable to insert order into table: %w", err))
				continue
			}

			// log trade to BQ
			err = dl.Log(ctx, wl, string(strategy), signal, resp)
			if err != nil {
				wl.Error(fmt.Errorf("unable to log trade to bigquery: %w", err))
				continue
			}
		}
		render.JSON(ctx, wl, w, nil, http.StatusOK)
	}
}

func checkForSignal(
	ctx context.Context,
	wl wlog.Logger,
	sig *signalapi.SignalResponse,
	db *sqlx.DB,
	sl *signallogger.DataLogger,
) (*entities.Signal, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
	id, err := ulid.New(ulid.Now(), rng)
	if err != nil {
		return nil, fmt.Errorf("error generating id: %s", err)
	}

	signal := &entities.Signal{
		ID:     id.String(),
		Chain:  sig.Chain,
		Signal: entities.SignalType(sig.Signal),
	}

	delta := false

	// get last chain trade from DB
	lastSignal, err := queryLastSignalForChain(ctx, wl, sig.Chain, db)
	if err != nil {
		// if no history found, use the last signal
		wl.Infof("no history found, using last signal %s on %s", sig.LastTrade, sig.LastTradeTime)
		signal.Signal = entities.SignalType(sig.LastTrade)
		signal.TradeTime = sig.LastTradeTime

		// log signal to BQ
		err = sl.Log(ctx, wl, delta, sig, signal.ID) // FIXME: Should be passing internal entities
		if err != nil {
			wl.Error(fmt.Errorf("unable to log signal to bigquery: %w", err))
		}
		return signal, nil
	}

	wl.Debugf("lastSignal: %+v", lastSignal)

	// handle no history found
	if lastSignal == nil {
		// this shouldn't happen as sqlx returns an error for 404
		// log signal to BQ
		err = sl.Log(ctx, wl, delta, sig, signal.ID) // FIXME: Should be passing internal entities
		if err != nil {
			wl.Error(fmt.Errorf("unable to log signal to bigquery: %w", err))
		}
		return nil, fmt.Errorf("no trade history found for chain %s", sig.Chain)
	}

	// if last trade in DB == last_trade_time in signal, then we no-op
	if lastSignal.Signal == entities.SignalType(sig.LastTrade) &&
		lastSignal.TradeTime.Equal(sig.LastTradeTime) {
		wl.Debugf(
			"no-op: %s == %s && %s == %s",
			lastSignal.Signal,
			entities.SignalType(sig.LastTrade),
			lastSignal.TradeTime,
			sig.LastTradeTime,
		)
		signal.Signal = entities.Null
		signal.TradeTime = lastSignal.TradeTime
	} else { // new signal found!
		wl.Debugf(
			"op: %s == %s && %s == %s",
			lastSignal.Signal,
			entities.SignalType(sig.LastTrade),
			lastSignal.TradeTime,
			sig.LastTradeTime,
		)

		delta = true
		// we need to create a new signal with the last_trade details and return it
		signal.Signal = entities.SignalType(sig.LastTrade)
		signal.TradeTime = sig.LastTradeTime

	}

	// log signal to BQ
	wl.Debug("logging signal to BQ")
	err = sl.Log(ctx, wl, delta, sig, signal.ID) // FIXME: Should be passing internal entities
	if err != nil {
		wl.Error(fmt.Errorf("unable to log signal to bigquery: %w", err))
	}

	return signal, nil
}

func queryLastSignalForChain(
	ctx context.Context,
	wl wlog.Logger,
	chain string,
	db *sqlx.DB,
) (*entities.Signal, error) {
	var lastSignal entities.Signal

	query := "SELECT * FROM mvp_signal_log WHERE chain = $1 ORDER BY created_at DESC LIMIT 1"
	err := db.Get(&lastSignal, query, chain)
	if err != nil {
		wl.Debugf("error querying signal for chain %s with err: %s", chain, err)
		return nil, err
	}

	return &lastSignal, nil
}

func insertLatestSignalTradedForChain(
	ctx context.Context,
	wl wlog.Logger,
	chain string,
	signal *entities.Signal,
	db *sqlx.DB,
) error {
	signal.CreatedAt = null.NewTime(time.Now(), true)
	query := `INSERT INTO mvp_signal_log (chain, signal, trade_time, created_at) 
	VALUES (:chain, :signal, :trade_time, :created_at)`
	_, err := db.NamedQuery(query, signal)
	if err != nil {
		return err
	}

	return nil
}

func insertNewOrder(
	ctx context.Context,
	wl wlog.Logger,
	o *coinroutesapi.ClientOrderCreateResponse,
	db *sqlx.DB,
) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
	id, err := ulid.New(ulid.Now(), rng)
	if err != nil {
		return fmt.Errorf("error generating id: %s", err)
	}

	order := entities.Order{
		ID:            id.String(),
		ClientOrderId: o.ClientOrderId,
		Strategy:      o.Strategy,
		Status:        entities.OrderStatusType(o.OrderStatus),
		CurrencyPair:  o.CurrencyPair,
		AvgPrice:      o.AvgPrice,
		ExecutedQty:   o.ExecutedQty,
		CreatedAt:     null.NewTime(time.Now(), true),
	}

	var query string

	// check if finished_at is set and only insert it if it is
	if o.FinishedAt != "" {
		finishedAt, err := time.Parse(time.RFC3339Nano, o.FinishedAt)
		if err != nil {
			return err
		}

		order.FinishedAt = null.NewTime(finishedAt, true)
		query = `INSERT INTO mvp_order (
			id,
			client_order_id,
			strategy,
			status,
			currency_pair,
			avg_price,
			executed_qty,
			finished_at,
			created_at)
		VALUES (
			:id,
			:client_order_id,
			:strategy,
			:status,
			:currency_pair,
			:avg_price,
			:executed_qty,
			:finished_at,
			:created_at)`
	} else { // finished_at not set, therefore don't insert it
		query = `INSERT INTO mvp_order (
		id,
		client_order_id,
		strategy,
		status,
		currency_pair,
		avg_price,
		executed_qty,
		created_at)
	VALUES (
		:id,
		:client_order_id,
		:strategy,
		:status,
		:currency_pair,
		:avg_price,
		:executed_qty,
		:created_at)`

	}

	_, err = db.NamedQuery(query, order)
	if err != nil {
		return err
	}

	return nil
}

func parseSignalAndCreateTradePayload(
	ctx context.Context,
	wl wlog.Logger,
	sig *entities.Signal,
	cur coinroutesapi.CurrencyType,
	str coinroutesapi.StrategyType,
) (*coinroutesapi.ClientOrderCreateRequest, error) {
	var currencyPair coinroutesapi.CurrencyPairType
	tradeTTL := fiveMinutes
	intLength := 1

	if cur == coinroutesapi.BTC {
		currencyPair = coinroutesapi.BTCInversePerpetual
	} else {
		currencyPair = coinroutesapi.ETHInversePerpetual
	}

	// Parse signal and act on it
	if sig.Signal == entities.Buy || sig.Signal == entities.Long {
		// Handle Buy/Long
		payload := coinroutesapi.ClientOrderCreateRequest{
			OrderType:          coinroutesapi.SmartPost,
			OrderStatus:        coinroutesapi.Open,
			Aggression:         coinroutesapi.Neutral,
			CurrencyPair:       currencyPair,
			Quantity:           "", // this will be hydrated later
			Side:               coinroutesapi.Buy,
			Strategy:           str,
			UseFundingCurrency: true,
			EndOffset:          tradeTTL,
			IntervalLength:     intLength,
			IsTwap:             false,
		}

		// set twap to true for BTC
		if cur == coinroutesapi.BTC {
			payload.IsTwap = true
		}

		return &payload, nil

	} else if sig.Signal == entities.Sell || sig.Signal == entities.Short {
		// Handle Sell/Short
		payload := coinroutesapi.ClientOrderCreateRequest{
			OrderType:          coinroutesapi.SmartPost,
			OrderStatus:        coinroutesapi.Open,
			Aggression:         coinroutesapi.Neutral,
			CurrencyPair:       currencyPair,
			Quantity:           "", // this will be hydrated later
			Side:               coinroutesapi.Sell,
			Strategy:           str,
			UseFundingCurrency: false,
			EndOffset:          tradeTTL,
			IntervalLength:     intLength,
			IsTwap:             false,
		}

		// set twap to true for BTC
		if cur == coinroutesapi.BTC {
			payload.IsTwap = true
		}

		return &payload, nil

	} else if sig.Signal == entities.Null {
		wl.Debug("no-op: null signal received")
		return nil, nil
	} else {
		return nil, fmt.Errorf("invalid signal received: %s", sig.Signal)
	}
}

func getCoinRoutesBalanceForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	str coinroutesapi.StrategyType,
	cr *coinroutesapi.Client,
	chain coinroutesapi.CurrencyType,
) (float64, error) {
	var amt float64
	balances, err := cr.GetBalances(ctx, str)
	if err != nil {
		return amt, err
	}

	// parse balances to find chain and set quantity
	for _, v := range *balances {
		if strings.ToLower(v.Currency) == string(chain) &&
			v.Exchange == "binancefutures" {
			amt, err = strconv.ParseFloat(v.Amount, 32)
			if err != nil {
				return amt, err
			}
			return amt, nil
		}
	}

	return amt, fmt.Errorf("could not find balance")
}

func getCoinRoutesPositionForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	str coinroutesapi.StrategyType,
	cr *coinroutesapi.Client,
	currencyPair coinroutesapi.CurrencyPairType,
) (float64, error) {
	positions, err := cr.GetPositions(ctx, str)
	var qty float64
	if err != nil {
		return qty, err
	}

	wl.Debugf("found positions: %+v", positions)
	if len(*positions) == 0 {
		wl.Infof("no-op: no positions found")
		return qty, nil // FIXME: make custom error type here for no-op
	}

	// parse positions to find same currency pair and set quantity
	for _, v := range *positions {
		wl.Debugf("cur: %s, search: %s", v.CurrencyPair, currencyPair)
		if v.CurrencyPair == string(currencyPair) {
			qty, err := strconv.ParseFloat(v.Quantity, 32)
			if err != nil {
				wl.Error(err)
				continue
			}

			return math.Abs(qty), nil
		}
	}

	return qty, fmt.Errorf("unable to find positions")
}
