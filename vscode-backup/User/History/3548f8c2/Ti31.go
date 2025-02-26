// Package signalhdl is the handler that handles all signal HTTP requests
package handler

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/tradelogger"
	"yield-mvp/internal/wlog"
	"yield-mvp/internal/ycontext"
	"yield-mvp/pkg/coinroutesapi"
	"yield-mvp/pkg/render"
	"yield-mvp/pkg/signalapi"

	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
)

func GetBTCSignal(
	wl wlog.Logger,
	sc *signalapi.Client,
	cc *coinroutesapi.Client,
	strats []string,
	db *sqlx.DB,
	dl *tradelogger.DataLogger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// FIXME: This should be added by a middleware
		chain := "btc"
		r = r.WithContext(ycontext.WithChain(r.Context(), chain))
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "signal")

		signalResp, err := sc.GetBTCSignal(ctx)
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

		// parse last signal to see if we have an operation to make
		signal, err := checkForSignal(ctx, wl, signalResp, db)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		if signal.Signal == entities.Null {
			// FIXME: This is just returning as a no-op. Should make this more relevant
			render.JSON(ctx, wl, w, signal, http.StatusOK)
			return
		}

		wl.Debugf("calculated signal: %+v", signal)

		// execute signal on all strategies
		for _, v := range strats {
			strategy := coinroutesapi.StrategyType(v)
			if err := strategy.Validate(); err != nil {
				wl.Info("skipping invalid strategy")
				continue
			}

			wl.Infof("handling strategy: %s", v)
			wl = wlog.WithStrategy(wl, v)
			var resp *coinroutesapi.ClientOrderCreateResponse

			// parse signal from signal API and create proper payload for CoinRoutes
			payload, err := parseSignalAndCreateTradePayload(
				ctx,
				wl,
				signal,
				coinroutesapi.BTC,
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
				positions, err := cc.GetPositions(ctx, strategy)
				if err != nil {
					wl.Error(err)
					continue
				}

				wl.Debugf("found positions: %+v", positions)
				if len(*positions) == 0 {
					wl.Infof("no-op: no positions found")
					continue
				}
				// parse positions to find same currency pair and set quantity
				for _, v := range *positions {
					wl.Debugf("cur: %s, search: %s", v.CurrencyPair, payload.CurrencyPair)
					if v.CurrencyPair == string(payload.CurrencyPair) {
						qty, err := strconv.ParseFloat(v.Quantity, 32)
						if err != nil {
							wl.Error(err)
							continue
						}

						payload.Quantity = strconv.FormatFloat(math.Abs(qty), 'f', 10, 32)
					}
				}

			case coinroutesapi.Sell:
				// short available balance (create contracts)
				balances, err := cc.GetBalances(ctx, strategy)
				if err != nil {
					wl.Error(err)
					continue
				}

				// parse balances to find chain and set quantity
				for _, v := range *balances {
					if strings.ToLower(v.Currency) == string(coinroutesapi.BTC) &&
						v.Exchange == "binancefutures" {
						amt, err := strconv.ParseFloat(v.Amount, 32)
						if err != nil {
							wl.Error(err)
							continue
						}

						if strategy == coinroutesapi.BTCLongShort {
							amt = 2 * amt // 2x leverage for long-short
						}

						payload.Quantity = strconv.FormatFloat(amt, 'f', 10, 32)
					}
				}
			}

			// check for zero amount
			orderFloat, err := strconv.ParseFloat(payload.Quantity, 32)
			if err != nil {
				render.InternalError(ctx, wl, w, err)
				return
			}

			if orderFloat == 0 {
				wl.Infof("no-op: quantity is zero value: %+v", payload.Quantity)
				render.JSON(ctx, wl, w, signal, http.StatusOK)
				return
			}

			// Make Trade
			wl.Debugf("about to create order with payload: %+v", payload)
			resp, err = cc.CreateClientOrders(ctx, payload)
			if err != nil {
				wl.Error(fmt.Errorf("unable to place order with error: %w", err))
				continue
			}

			//  record trade in log
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

			// log to BQ
			dl.Log(ctx, wl, string(strategy), signal, resp)
		}

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

		// parse last signal to see if we have an operation to make
		signal, err := checkForSignal(ctx, wl, signalResp, db)
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
				positions, err := cc.GetPositions(ctx, strategy)
				if err != nil {
					wl.Error(err)
					continue
				}

				wl.Debugf("found positions: %+v", positions)
				if len(*positions) == 0 {
					wl.Infof("no-op: no positions found")
					continue
				}
				// parse positions to find same currency pair and set quantity
				for _, v := range *positions {
					wl.Debugf("cur: %s, search: %s", v.CurrencyPair, payload.CurrencyPair)
					if v.CurrencyPair == string(payload.CurrencyPair) {
						qty, err := strconv.ParseFloat(v.Quantity, 32)
						if err != nil {
							wl.Error(err)
							continue
						}

						payload.Quantity = strconv.FormatFloat(math.Abs(qty), 'f', 10, 32)
					}
				}

			case coinroutesapi.Sell:
				// short available balance (create contracts)
				balances, err := cc.GetBalances(ctx, strategy)
				if err != nil {
					wl.Error(err)
					continue
				}

				// parse balances to find chain and set quantity
				for _, v := range *balances {
					if strings.ToLower(v.Currency) == string(coinroutesapi.ETH) &&
						v.Exchange == "binancefutures" {
						amt, err := strconv.ParseFloat(v.Amount, 32)
						if err != nil {
							wl.Error(err)
							continue
						}

						if strategy == coinroutesapi.ETHLongShort {
							amt = 2 * amt // 2x leverage for long-short
						}

						payload.Quantity = strconv.FormatFloat(amt, 'f', 10, 32)
					}
				}
			}

			// check for zero amount
			orderFloat, err := strconv.ParseFloat(payload.Quantity, 32)
			if err != nil {
				render.InternalError(ctx, wl, w, err)
				return
			}

			if orderFloat == 0 {
				wl.Infof("no-op: quantity is zero value: %+v", payload.Quantity)
				render.JSON(ctx, wl, w, signal, http.StatusOK)
				return
			}

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

			// log to BQ
			err = dl.Log(ctx, wl, string(strategy), signal, resp)
			if err != nil {
				render.InternalError(ctx, wl, w, err)
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
) (*entities.Signal, error) {
	signal := &entities.Signal{
		Chain:  sig.Chain,
		Signal: entities.SignalType(sig.Signal),
	}

	// get last chain trade from DB
	lastSignal, err := queryLastSignalForChain(ctx, wl, sig.Chain, db)
	if err != nil {
		// if no history found, use the last signal
		wl.Infof("no history found, using last signal %s on %s", sig.LastTrade, sig.LastTradeTime)
		signal.Signal = entities.SignalType(sig.LastTrade)
		signal.TradeTime = sig.LastTradeTime
		return signal, nil
	}

	wl.Debugf("lastSignal: %+v", lastSignal)

	// handle no history found
	if lastSignal == nil {
		// this shouldn't happen as sqlx returns an error for 404
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
		return signal, nil
	}

	wl.Debugf(
		"op: %s == %s && %s == %s",
		lastSignal.Signal,
		entities.SignalType(sig.LastTrade),
		lastSignal.TradeTime,
		sig.LastTradeTime,
	)

	// we need to create a new signal with the last_trade details and return it
	signal.Signal = entities.SignalType(sig.LastTrade)
	signal.TradeTime = sig.LastTradeTime

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

func parseSignalAndCreateTradePayload(
	ctx context.Context,
	wl wlog.Logger,
	sig *entities.Signal,
	cur coinroutesapi.CurrencyType,
	str coinroutesapi.StrategyType,
) (*coinroutesapi.ClientOrderCreateRequest, error) {
	var currencyPair coinroutesapi.CurrencyPairType

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
		}

		return &payload, nil

	} else if sig.Signal == entities.Null {
		wl.Debug("no-op: null signal received")
		return nil, nil
	} else {
		return nil, fmt.Errorf("invalid signal received: %s", sig.Signal)
	}
}
